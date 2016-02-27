package server

import (
	"bytes"
	"fmt"
	"github.com/GrappigPanda/notorious/bencode"
	"gopkg.in/redis.v3"
	"net/http"
	"net/url"
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
	port       string // port number the peer is listening on
	uploaded   int    // base10 ascii amount uploaded so far
	downloaded int    // base10 ascii amount downloaded so far
	left       int    // # of bytes left to download (base 10 ascii)
	event      int
	compact    bool
}

var ANNOUNCE_URL = "/announce"

var FIELDS = []string{"port", "uploaded", "downloaded", "left", "event", "compact"}

func parseTorrentGetRequestURI(request string, ip string) map[string]interface{} {
	torrentdata := make(map[string]interface{})
	querydata := decodeQueryURL(request)

	torrentdata["info_hash"] = ParseInfoHash(querydata["info_hash"][0])
	fmt.Println(torrentdata["info_hash"])
	torrentdata["ip"] = ip
	torrentdata["port"] = querydata["port"][0]
	torrentdata["peer_id"] = querydata["peer_id"][0]

	return torrentdata
}

func ParseInfoHash(s string) string {
	return fmt.Sprintf("%x", s)
}

func fillEmptyMapValues(torrentMap map[string]interface{}) *TorrentRequestData {
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
		torrentMap["port"].(string),
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
		x := RedisGetKeyVal(client, torrdata.info_hash, torrdata)
		return x
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

func decodeQueryURL(s string) url.Values {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	m, _ := url.ParseQuery(u.RawQuery)
	return m
}

func requestHandler(w http.ResponseWriter, req *http.Request) {
	client := OpenClient()

	torrentdata := parseTorrentGetRequestURI(req.RequestURI, req.RemoteAddr)
	if torrentdata != nil {
		data := fillEmptyMapValues(torrentdata)

		worker(client, data)
		ipData := bencode.EncodeResponse()

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
	fmt.Println(keymember)

	val, err := client.SMembers(keymember).Result()
	if err != nil {
		panic("Invalid lookup")
	}
	if len(val) == 0 {
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
