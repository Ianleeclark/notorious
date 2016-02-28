package server

import (
	"bytes"
	"fmt"
	"github.com/GrappigPanda/notorious/bencode"
	"gopkg.in/redis.v3"
	"net/http"
	"net/url"
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
	torrentdata["ip"] = strings.Split(ip, ":")[0]
	fmt.Println(torrentdata["ip"])
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

		RedisSetKeyVal(client,
			concatenateKeyMember(torrdata.info_hash, "ip"),
			createIpPortPair(torrdata))

		return x

	} else {
		CreateNewTorrentKey(client, torrdata.info_hash, torrdata)
		return worker(client, torrdata)
	}
}

func formatResponseData(ipport []string, torrentdata *TorrentRequestData) string {
	if torrentdata.compact {
		return EncodeResponse(ipport)
	} else {
		// TODO(ian): Support bep-23
		return EncodeResponse(ipport)
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

func encodeKV(key string, value string) string {
	return fmt.Sprintf("%s%s", bencode.EncodeByteString(key), bencode.EncodeByteString(value))
}

func EncodeResponse(ipport []string) string {
	ret := "d"

	ret += encodeKV("interval", "30") // Interval
	ret += encodeKV("tracker_id", "1234")
	ret += encodeKV("complete", "1")
	ret += encodeKV("incomplete", "111111111111111111111")
	ret += "5:peersd"

	for i := range ipport {
		data := strings.Split(ipport[i], ":")

		ret += encodeKV("peer_id", "1")
		ret += encodeKV("ip", data[0])
		ret += encodeKV("port", data[1])
	}

	ret += "e"

	ret += "e"

	return ret
}

func requestHandler(w http.ResponseWriter, req *http.Request) {
	client := OpenClient()

	torrentdata := parseTorrentGetRequestURI(req.RequestURI, req.RemoteAddr)
	if torrentdata != nil {
		data := fillEmptyMapValues(torrentdata)

		worker(client, data)
		x := RedisGetKeyVal(client, data.info_hash, data)

		response := formatResponseData(x, data)
		fmt.Println(response)

		w.Write([]byte(response))
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
	// CreateNewTorrentKey creates a new key. By default, it adds a member
	// ":ip". I don't think this ought to ever be generalized, as I just want
	// Redis to function in one specific way in notorious.

	// TODO(ian): You might want to set this explicitly in parameters
	// value := *TorrentRequestData
	client.SAdd(key, "ip")
}

func createIpPortPair(value *TorrentRequestData) string {
	// createIpPortPair creates a string formatted ("%s:%s", value.ip,
	// value.port) looking like so: "127.0.0.1:6886" and returns this value.
	fmt.Println(value.ip, value.port)
	return fmt.Sprintf("%s:%s", value.ip, value.port)
}

func RedisSetKeyVal(client *redis.Client, keymember string, value string) {
	// RedisSetKeyVal sets a key:member's value to value. Returns nothing as of
	// yet.
	client.SAdd(keymember, value)
}

func RedisGetKeyVal(client *redis.Client, key string, value *TorrentRequestData) []string {
	// RedisGetKeyVal retrieves a value from the Redis store by looking up the
	// provided key. If the key does not yet exist, we create the key in the KV
	// storage or if the value is empty, we add the current requester to the
	// list.
	keymember := concatenateKeyMember(key, "ip")

	val, err := client.SMembers(keymember).Result()
	if err != nil {
		// Fail because the key doesn't exist in the KV storage.
		CreateNewTorrentKey(client, keymember, value)
	}

	// If no keys yet exist in the KV storage.
	if len(val) == 0 {
		RedisSetKeyVal(client, keymember, createIpPortPair(value))
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
