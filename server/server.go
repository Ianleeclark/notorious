package server

import (
	"bytes"
	"fmt"
	"github.com/GrappigPanda/notorious/bencode"
	"gopkg.in/redis.v3"
	"net"
	"net/http"
	"net/url"
	"strconv"
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
	event      int    // TorrentEvent
	numwant    int    // Number of peers requested by client.
	compact    bool   // Bep23 peer list compression decision: True -> compress bep23
}

var ANNOUNCE_URL = "/announce"

var FIELDS = []string{"port", "uploaded", "downloaded", "left", "event", "compact"}

func parseTorrentGetRequestURI(request string, ip string) map[string]interface{} {
	torrentdata := make(map[string]interface{})
	querydata := decodeQueryURL(request)

	lookupvalues := []string{
		"info_hash",
		"peer_id",
		// ip intentially left out
		"port",
		"uploaded",
		"downloaded",
		"left",
		"event",
		"numwant",
		"compact",
	}

	lookupints := []string{
		"uploaded",
		"downloaded",
		"left",
		"event",
		"numwant",
	}

	for i := range lookupvalues {
		index := lookupvalues[i]

		// TODO(ian): Delegate responsibility to another function.
		x := querydata[index]
		if len(x) <= 0 {
			continue
		}

		if index == "info_hash" {
			torrentdata[index] = ParseInfoHash(x[0])
		} else if index == "compact" {
			if querydata[index][0] == "1" {
				torrentdata[index] = true
			} else {
				torrentdata[index] = false
			}
		} else {
			torrentdata[index] = querydata[index][0]
		}

		// TODO(ian): Delegate responsibility here, as well.
		for j := range lookupints {
			if lookupints[j] == index {
				val, _ := strconv.Atoi(index)
				torrentdata[index] = val
			}
		}
	}

	// TODO(ian): Add a generic size lookup function and assert > 0
	x := strings.Split(ip, ":")
	if len(x) > 0 {
		torrentdata["ip"] = x[0]
	}

	fmt.Println(torrentdata)
	return torrentdata
}

func ParseInfoHash(s string) string {
	return fmt.Sprintf("%x", s)
}

func fillEmptyMapValues(torrentMap map[string]interface{}) *TorrentRequestData {
	_, ok := torrentMap["uploaded"]
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

	_, ok = torrentMap["numwant"]
	if !ok {
		torrentMap["numwant"] = 30
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
		torrentMap["numwant"].(int),
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
func compactIPPort(ip string, port string) []byte {
    res := bytes.NewBuffer(make([]byte, 0))

    intPort, err := strconv.Atoi(port)
    if err != nil {
        panic("failure1")
    }


    if err := binary.Write(res, binary.BigEndian, binary.BigEndian.Uint32(net.ParseIP(ip).To4())); err != nil {
        panic("failure0")
    }

    err = binary.Write(res, binary.BigEndian, uint16(intPort)); 
        if err != nil {
        panic("failure2")
    }

    return res.Bytes()
}

func CompactAllPeers(ipport []string) []byte {
    ret := bytes.NewBuffer(make([]byte, 0))
    for i := range ipport {
        sz := strings.Split(ipport[i], ":")
        ip := sz[0]
        port := sz[1]

        ret += compactIPPort(ip, port)
    }

    return ret.Bytes()
}

func formatResponseData(ips []string, torrentdata *TorrentRequestData) string {
	for i := range ips {
		ips[i] = compactifyIpPort(ips[i])
	}

	if torrentdata.compact {
		// TODO(ian): Support bep-23
		return EncodeResponse(ips, torrentdata)
	} else {
		return EncodeResponse(ips, torrentdata)
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

func EncodeResponse(ipport []string, torrentdata *TorrentRequestData) string {
	ret := ""
	ret += encodeKV("interval", "30")
	ret += encodeKV("complete", "1")
	ipstr := ""
	
	if torrentdata.compact {
		count := 0
		
		for i := range ipport {
			ipstr += fmt.Sprintf("%s", ipport[i])
			count += 1
		}
		
		ret += encodeKV("incomplete", fmt.Sprintf("%d", count))
		ret += bencode.EncodeDictionary("peers", ipstr)
	} else {

		for i := range ipport {
			data := strings.Split(ipport[i], ":")

			ipstr += encodeKV("peer_id", "1")
			ipstr += encodeKV("ip", data[0])
			ipstr += encodeKV("port", data[1])
		}

		ret += bencode.EncodeDictionary("peers", ipstr)
	}

	return bencode.EncodeDictionary("", ret)
}

func requestHandler(w http.ResponseWriter, req *http.Request) {
	client := OpenClient()

	torrentdata := parseTorrentGetRequestURI(req.RequestURI, req.RemoteAddr)
	if torrentdata != nil {
		data := fillEmptyMapValues(torrentdata)

		worker(client, data)
		x := RedisGetKeyVal(client, data.info_hash, data)
		fmt.Println(x)

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
