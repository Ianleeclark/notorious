package server

import (
	"fmt"
	"gopkg.in/redis.v3"
	"net/http"
)

var FIELDS = []string{"port", "uploaded", "downloaded", "left", "event", "compact"}

func worker(client *redis.Client, data *announceData) []string {
	if RedisGetBoolKeyVal(client, data.info_hash, data) {
		x := RedisGetKeyVal(client, data.info_hash, data)

		RedisSetKeyVal(client,
			concatenateKeyMember(data.info_hash, "ip"),
			createIpPortPair(data))

		return x

	} else {
		CreateNewTorrentKey(client, data.info_hash)
		return worker(client, data)
	}
}

func requestHandler(w http.ResponseWriter, req *http.Request) {
	client := OpenClient()

	data := new(announceData)
	err := data.parseAnnounceData(req)
	if err != nil {
		panic(err)
	}

	switch data.event {
	case "started":
		data.StartedEventHandler(client)
	case "stopped":
		data.StoppedEventHandler(client)
	case "completed":
		data.CompletedEventHandler(client)
	}
	fmt.Printf("Event: %s from host %s on port %v\n", data.event, data.ip, data.port)

	worker(client, data)
	x := RedisGetKeyVal(client, data.info_hash, data)
	if len(x) >= 30 {
		x = x[0:30]
	} else {
		x = x[0:len(x)]
	}

	if len(x) > 0 {
		response := formatResponseData(client, x, data)
		fmt.Printf("Resp: %s\n", response)

		w.Write([]byte(response))
	} else {
		failMsg := fmt.Sprintf("No peers for torrent %s\n", data.info_hash)
		w.Write([]byte(createFailureMessage(failMsg)))
	}
}

func RunServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/announce", requestHandler)
	http.ListenAndServe(":3000", mux)
}

func OpenClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return client
}
