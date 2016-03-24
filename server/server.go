package server

import (
	"fmt"
	"net/http"
)

var FIELDS = []string{"port", "uploaded", "downloaded", "left", "event", "compact"}

func worker(data *announceData) []string {
	if RedisGetBoolKeyVal(data.redisClient, data.info_hash, data) {
		x := RedisGetKeyVal(data.redisClient, data.info_hash, data)

		RedisSetIPMember(data)

		return x

	} else {
		CreateNewTorrentKey(data.redisClient, data.info_hash)
		return worker(data)
	}
}

func requestHandler(w http.ResponseWriter, req *http.Request) {
	data := new(announceData)
	err := data.parseAnnounceData(req)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Event: %s from host %s on port %v\n", data.event, data.ip, data.port)

	switch data.event {

	case "started":
		data.StartedEventHandler()

	case "stopped":
		data.StoppedEventHandler()

	case "completed":
		data.CompletedEventHandler()
	default:
		panic(fmt.Errorf("We're somehow getting this strange error..."))
	}

	if data.event == "started" || data.event == "completed" {
		worker(data)
		x := RedisGetKeyVal(data.redisClient, data.info_hash, data)
		// TODO(ian): Move this into a seperate function.
		// TODO(ian): Remove this magic number and use data.numwant, but limit it
		// to 30 max, as that's the bittorrent protocol suggested limit.
		if len(x) >= 30 {
			x = x[0:30]
		} else {
			x = x[0:len(x)]
		}

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

func RunServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/announce", requestHandler)
	http.ListenAndServe(":3000", mux)
}
