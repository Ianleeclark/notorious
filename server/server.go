package server

import (
	"fmt"
	"net/http"
)

var FIELDS = []string{"port", "uploaded", "downloaded", "left", "event", "compact"}

func worker(data *announceData) []string {
	if RedisGetBoolKeyVal(data.redisClient, data.info_hash, data) {
		x := RedisGetKeyVal(data.redisClient, data.info_hash, data)

		RedisSetKeyVal(data.redisClient,
			concatenateKeyMember(data.info_hash, "ip"),
			createIpPortPair(data))

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

	switch data.event {
	case "started":
		data.StartedEventHandler()
	case "stopped":
		data.StoppedEventHandler()
	case "completed":
		data.CompletedEventHandler()
	default:
		data.StartedEventHandler()
	}

	fmt.Printf("Event: %s from host %s on port %v\n", data.event, data.ip, data.port)

	if data.event == "started" || data.event == "completed" || data.event == "" || data.event == " " {
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
			response := formatResponseData(x, data)
			fmt.Printf("Resp: %s\n", response)

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
