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
	db.UpdateStats(data.uploaded, data.downloaded)

	if app.trackerLevel > RATIOLESS {
		db.UpdatePeerStats(data.uploaded, data.downloaded, data.ip)
	}

	if data.event == "completed" {
		db.UpdateTorrentStats(1, -1)
		return
	} else if data.left == 0 {
		// TODO(ian): Don't assume the peer is already in the DB
		db.UpdateTorrentStats(1, -1)
		return
	} else if data.event == "started" {
		db.UpdateTorrentStats(0, 1)
	}
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
			writeErrorResponse(w, err.Error())

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
			response := formatResponseData(x, data)
			writeResponse(w, response)

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
	if infoHash == "" {
		failMsg := fmt.Sprintf("Tracker does not support multiple entire DB scrapes.")
		writeErrorResponse(w, failMsg)
	} else {
		torrentData := db.ScrapeTorrent(dbConn, infoHash)
		writeResponse(w, formatScrapeResponse(torrentData))
	}

	return
}

func writeErrorResponse(w http.ResponseWriter, failMsg string) {
	writeResponse(w, createFailureMessage(failMsg))
}

func writeResponse(w http.ResponseWriter, values string) {
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
