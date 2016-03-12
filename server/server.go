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
	data.parseAnnounceData(req.URL)

	worker(client, data)
	x := RedisGetKeyVal(client, data.info_hash, data)
	fmt.Println(x)

	response := formatResponseData(x, data)
	fmt.Println(response)

	w.Write([]byte(response))
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
