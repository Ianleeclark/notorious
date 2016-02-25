package server

import (
	"bytes"
	"fmt"
	"github.com/GrappigPanda/notorious/bencode"
	"gopkg.in/redis.v3"
	"net/http"
	"strings"
)

const (
	STARTED = iota
	COMPLETED
	STOPPED
)

type TorrentEvent struct {
	started   int
	completed int
	stopped   int
}

type TorrentRequestData struct {
	info_hash  string //20 byte sha1 hash
	peer_id    string //max len 20
	ip         string //optional
	port       int    // port number the peer is listening on
	uploaded   int    // base10 ascii amount uploaded so far
	downloaded int    // base10 ascii amount downloaded so far
	left       int    // # of bytes left to download (base 10 ascii)
	event      int
	compact    bool
}

var ANNOUNCE_URL = "/announce"

func parseTorrentGetRequestURI(s string) map[string]interface{} {
	tmp := strings.Split(s, "?")
	tmp = strings.Split(tmp[1], "%26")
	result := make(map[string]interface{})
	for i := range tmp {
		if tmp[i] != ANNOUNCE_URL {
			data := strings.Split(tmp[i], "=")
			result[data[0]] = data[1]
		}
	}
	return result
}
func fillEmptyMapValues(torrentMap map[string]interface{}) *TorrentRequestData {
	// TODO(ian): DRY.
	_, ok := torrentMap["port"]
	if !ok {
		torrentMap["port"] = 0
	}
	_, ok = torrentMap["uploaded"]
	if !ok {
		torrentMap["uploaded"] = 0
	}
	_, ok = torrentMap["downloaded"]
	if !ok {
		torrentMap["downloaded"] = 0
	}
	_, ok = torrentMap["left"]
	if !ok {
		torrentMap["left"] = 0
	}
	_, ok = torrentMap["event"]
	if !ok {
		torrentMap["event"] = STOPPED
	}
	_, ok = torrentMap["compact"]
	if !ok {
		torrentMap["compact"] = true
	}

	x := TorrentRequestData{
		torrentMap["info_hash"].(string),
		torrentMap["peer_id"].(string),
		torrentMap["ip"].(string),
		torrentMap["port"].(int),
		torrentMap["uploaded"].(int),
		torrentMap["downloaded"].(int),
		torrentMap["left"].(int),
		torrentMap["event"].(int),
		torrentMap["compact"].(bool),
	}
	return &x
}

func worker(client *redis.Client, torrdata *TorrentRequestData) []string {
	if RedisGetBoolKeyVal(client, torrdata.info_hash, torrdata) {
		return RedisGetKeyVal(client, torrdata.info_hash, torrdata)
	} else {
		CreateNewTorrentKey(client, torrdata.info_hash, torrdata)
		return worker(client, torrdata)
	}
}

func formatIpData(ips []string, compact bool) string {
	encodedList := bencode.EncodeList(ips)

	if compact {
		return encodedList
	} else {
		// TODO(ian): Support non bep-23
		return encodedList
	}
}

func requestHandler(w http.ResponseWriter, req *http.Request) {
	client := OpenClient()
	fmt.Printf("%v", req)

	torrentdata := parseTorrentGetRequestURI(req.RequestURI)
	fmt.Printf("%v", torrentdata)
	if torrentdata != nil {
		data := fillEmptyMapValues(torrentdata)

		ipData := formatIpData(worker(client, data), data.compact)

		w.Write([]byte(ipData))
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

func CreateNewTorrentKey(client *redis.Client, key string, value *TorrentRequestData) {
	// TODO(ian): You might want to set this explicitly in parameters
	// value := *TorrentRequestData

	// Here the key is the info_hash for the torrent and value is
	// the newest peer for the torrent
	client.SAdd(key, "ip")
	RedisSetKeyVal(client, key, "ip", value.ip)
}

func RedisSetKeyVal(client *redis.Client, key string, member string, value string) interface{} {
	keymember := concatenateKeyMember(key, member)
	client.SAdd(keymember, value)
	return 1
}

func RedisGetKeyVal(client *redis.Client, key string, value *TorrentRequestData) []string {
	keymember := concatenateKeyMember(key, "ip")

	val, err := client.SMembers(keymember).Result()
	if err != nil {
		CreateNewTorrentKey(client, key, value)
	}

	return val
}

func RedisGetBoolKeyVal(client *redis.Client, key string, value interface{}) bool {
	_, err := client.Get(key).Result()

	return err != nil
}

func concatenateKeyMember(key string, member string) string {
	var buffer bytes.Buffer
	buffer.WriteString(key)
	buffer.WriteString(":")
	buffer.WriteString(member)
	return buffer.String()
}
