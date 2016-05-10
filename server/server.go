package server

import (
	"fmt"
	"github.com/GrappigPanda/notorious/config"
	"github.com/GrappigPanda/notorious/database"
	"net/http"
)

// FIELDS The fields that we expect from a peer upon info hash lookup
var FIELDS = []string{"port", "uploaded", "downloaded", "left", "event", "compact"}

func worker(data *announceData) []string {
	if RedisGetBoolKeyVal(data.requestContext.redisClient, data.info_hash) {
		x := RedisGetKeyVal(data, data.info_hash)

		RedisSetIPMember(data)

		return x

	}

	CreateNewTorrentKey(data.requestContext.redisClient, data.info_hash)
	return worker(data)
}
func (app *applicationContext) handleStatsTracking(data *announceData) {
	if app.trackerLevel > RATIOLESS {
		db.UpdatePeerStats(data.uploaded, data.downloaded, data.ip)
	}
	db.UpdateStats(data.uploaded, data.downloaded)
}

func (app *applicationContext) requestHandler(w http.ResponseWriter, req *http.Request) {
	data := new(announceData)
	data.requestContext = requestAppContext{
		dbConn:    nil,
		whitelist: app.config.Whitelist,
	}

	err := data.parseAnnounceData(req)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Event: %s from host %s on port %v\n", data.event, data.ip, data.port)

	switch data.event {
	case "started":
		err := data.StartedEventHandler()
		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(createFailureMessage(err.Error())))

			return
		}
	case "stopped":
		data.StoppedEventHandler()

	case "completed":
		data.CompletedEventHandler()
	default:
		panic(fmt.Errorf("We're somehow getting this strange error..."))
	}

	if data.event == "started" || data.event == "completed" {
		worker(data)
		x := RedisGetAllPeers(data, data.info_hash)

		if len(x) > 0 {
			w.Header().Set("Content-Type", "text/plain")
			response := formatResponseData(x, data)

			w.Write([]byte(response))

		} else {
			failMsg := fmt.Sprintf("No peers for torrent %s\n",
				data.info_hash)
			writeErrorResponse(w, failMsg)
		}
	}

	app.handleStatsTracking(data)
}

func scrapeHandler(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	dbConn, err := db.OpenConnection()
	if err != nil {
		panic(err)
	}

	infoHash := query.Get("info_hash")
	if infoHash != "" {
		values, err := db.ScrapeTorrentFromInfoHash(
			dbConn,
			ParseInfoHash(infoHash))
		if err != nil {
			failMsg := fmt.Sprintf("Torrent not found.")
			writeErrorResponse(w, failMsg)
		}

		writeResponse(w, req, values)
	} else {
		writeResponse(w, req, db.ScrapeTorrent(dbConn))
	}

	return
}

func writeErrorResponse(w http.ResponseWriter, failMsg string) {
	w.Write([]byte(createFailureMessage(failMsg)))
}

func writeResponse(w http.ResponseWriter, req *http.Request, values string) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(values))
}

// RunServer spins up the server and muxes the routes.
func RunServer() {
	app := applicationContext{
		config:       config.LoadConfig(),
		trackerLevel: RATIOLESS,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/announce", app.requestHandler)
	mux.HandleFunc("/scrape", scrapeHandler)
	http.ListenAndServe(":3000", mux)
}
