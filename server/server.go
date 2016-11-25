package server

import (
	"fmt"
	"github.com/GrappigPanda/notorious/config"
	"github.com/GrappigPanda/notorious/database"
	. "github.com/GrappigPanda/notorious/announce"
	"github.com/GrappigPanda/notorious/server/peerStore"
	"net/http"
)

// applicationContext houses data necessary for the handler to properly
// function as the application is desired.
type applicationContext struct {
	config       config.ConfigStruct
	trackerLevel int
}

type scrapeData struct {
	infoHash string
}

// scrapeResponse is the data associated with a returned scrape.
type scrapeResponse struct {
	complete   uint64
	downloaded uint64
	incomplete uint64
}

// TorrentResponseData models what is sent back to the peer upon a succesful
// info hash lookup.
type TorrentResponseData struct {
	interval     int
	min_interval int
	tracker_id   string
	completed    int
	incomplete   int
	peers        interface{}
}

// ANNOUNCE_URL The announce path for the http calls to reach.
var ANNOUNCE_URL = "/announce"

// TODO(ian): Set this expireTime to a config-loaded value.
// expireTime := 5 * 60
// FIELDS The fields that we expect from a peer upon info hash lookup
var FIELDS = []string{"port", "uploaded", "downloaded", "left", "event", "compact"}

func worker(data *AnnounceData) []string {
	if peerStore.RedisGetBoolKeyVal(data.InfoHash) {
		x := peerStore.RedisGetKeyVal(data.InfoHash)

        peerStore.RedisSetIPMember(data.InfoHash, fmt.Sprintf("%s:%s", data.IP, data.Port))

		return x

	}

	peerStore.CreateNewTorrentKey(data.InfoHash)
	return worker(data)
}
func (app *applicationContext) handleStatsTracking(data *AnnounceData) {
	db.UpdateStats(data.Uploaded, data.Downloaded)

	if app.trackerLevel > RATIOLESS {
		db.UpdatePeerStats(data.Uploaded, data.Downloaded, data.IP)
	}

	if data.Event == "completed" {
		db.UpdateTorrentStats(1, -1)
		return
	} else if data.Left == 0 {
		// TODO(ian): Don't assume the peer is already in the DB
		db.UpdateTorrentStats(1, -1)
		return
	} else if data.Event == "started" {
		db.UpdateTorrentStats(0, 1)
	}
}

func (app *applicationContext) requestHandler(w http.ResponseWriter, req *http.Request) {
	data := new(AnnounceData)
	err := data.ParseAnnounceData(req)
	if err != nil {
		panic(err)

	}

    data.RequestContext.Whitelist = app.config.Whitelist

	fmt.Printf("Event: %s from host %s on port %v\n", data.Event, data.IP, data.Port)

	switch data.Event {
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

	if data.Event == "started" || data.Event == "completed" {
		worker(data)
		x := peerStore.RedisGetAllPeers(data.InfoHash)

		if len(x) > 0 {
			response := formatResponseData(x, data)
			writeResponse(w, response)

		} else {
			failMsg := fmt.Sprintf("No peers for torrent %s\n",
				data.InfoHash)
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

	infoHash := query.Get("InfoHash")
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
