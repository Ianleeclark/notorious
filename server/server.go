package server

import (
	"fmt"
	"github.com/GrappigPanda/notorious/config"
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

func (app *applicationContext) requestHandler(w http.ResponseWriter, req *http.Request) {
	data := new(announceData)
	data.requestContext = requestAppContext{
		dbConn: nil,
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
			failMsg := fmt.Sprintf("No peers for torrent %s\n", data.info_hash)
			w.Write([]byte(createFailureMessage(failMsg)))
		}
	}
}

func scrapeHandler(w http.ResponseWriter, req *http.Request) interface{} {
	//query := req.URL.Query()
	//infoHash := ParseInfoHash(query.Get("info_hash"))

	//data := db.ScrapeTorrent(db.OpenConnection(), infoHash)
	//return data
	return "TODO"
}

// RunServer spins up the server and muxes the url
func RunServer() {
	app := applicationContext{
		config: config.LoadConfig(),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/announce", app.requestHandler)
	//mux.HandleFunc("/scrape", scrapeHandler)
	http.ListenAndServe(":3000", mux)
}
