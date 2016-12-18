package server

import (
	"fmt"
	a "github.com/GrappigPanda/notorious/announce"
	"github.com/GrappigPanda/notorious/announce/impl/rss"
	"github.com/GrappigPanda/notorious/config"
	"github.com/GrappigPanda/notorious/database/impl"
	"github.com/GrappigPanda/notorious/peerStore"
	"github.com/GrappigPanda/notorious/peerStore/impl"
	"log"
	"net/http"
)

// applicationContext houses data necessary for the handler to properly
// function as the application is desired.
type applicationContext struct {
	config          config.ConfigStruct
	trackerLevel    int
	peerStoreClient peerStore.PeerStore
	sqlObj          sqlStoreImpl.SQLStore
	rssNotifier     *rss.RSSNotifier
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
	interval    int
	minInterval int
	trackerID   string
	completed   int
	incomplete  int
	peers       interface{}
}

// AnnounceURL The announce path for the http calls to reach.
var AnnounceURL = "/announce"

// TODO(ian): Set this expireTime to a config-loaded value.
// expireTime := 5 * 60
// FIELDS The fields that we expect from a peer upon info hash lookup
var FIELDS = []string{"port", "uploaded", "downloaded", "left", "event", "compact"}

func (app *applicationContext) worker(data *a.AnnounceData) []string {
	if app.peerStoreClient.KeyExists(data.InfoHash) {
		x := app.peerStoreClient.GetKeyVal(data.InfoHash)

		app.peerStoreClient.SetIPMember(data.InfoHash, fmt.Sprintf("%s:%s", data.IP, data.Port))

		return x

	}

	app.peerStoreClient.CreateNewTorrentKey(data.InfoHash)
	return app.worker(data)
}

func (app *applicationContext) handleStatsTracking(data *a.AnnounceData) {
	app.sqlObj.UpdateStats(data.Uploaded, data.Downloaded)

	if app.trackerLevel > a.RATIOLESS {
		app.sqlObj.UpdatePeerStats(data.Uploaded, data.Downloaded, data.IP)
	}

	if data.Event == "completed" {
		app.sqlObj.UpdateTorrentStats(1, -1)
		return
	} else if data.Left == 0 {
		// TODO(ian): Don't assume the peer is already in the DB
		app.sqlObj.UpdateTorrentStats(1, -1)
		return
	} else if data.Event == "started" {
		app.sqlObj.UpdateTorrentStats(0, 1)
	}
}

func (app *applicationContext) requestHandler(w http.ResponseWriter, req *http.Request) {
	data := new(a.AnnounceData)
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
		app.worker(data)
		x := app.peerStoreClient.GetAllPeers(data.InfoHash)

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

func (app *applicationContext) scrapeHandler(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()

	infoHash := query.Get("InfoHash")
	if infoHash == "" {
		failMsg := fmt.Sprintf("Tracker does not support multiple entire DB scrapes.")
		writeErrorResponse(w, failMsg)
	} else {
		torrentData := app.sqlObj.ScrapeTorrent(infoHash)
		writeResponse(w, formatScrapeResponse(torrentData))
	}

	return
}

func (app *applicationContext) rssHandle(w http.ResponseWriter, req *http.Request) {
	rssData, err := app.rssNotifier.GetRSS()
	if err == nil {
		writeResponse(w, rssData)
	}
}

func writeErrorResponse(w http.ResponseWriter, failMsg string) {
	writeResponse(w, createFailureMessage(failMsg))
}

func writeResponse(w http.ResponseWriter, values string) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(values))
}

// RunServer spins up the server and muxes the routes.
func RunServer(rssNotifier *rss.RSSNotifier) {
	// Load the config and initiate a `SQLStore` implementation.
	sqlObj := sqlStoreImpl.InitSQLStoreByDBChoice()
	cfg := config.LoadConfig()

	app := applicationContext{
		config:          cfg,
		trackerLevel:    a.RATIOLESS,
		peerStoreClient: new(redisPeerStoreImpl.RedisStore),
		sqlObj:          sqlObj,
		rssNotifier:     rssNotifier,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/announce", app.requestHandler)
	mux.HandleFunc("/scrape", app.scrapeHandler)
	if cfg.UseRSS == true {
		log.Println("Starting RSS handler at http://localhost/rss/")
		mux.HandleFunc("/rss/", app.rssHandle)
	}
	http.ListenAndServe(":3000", mux)
}
