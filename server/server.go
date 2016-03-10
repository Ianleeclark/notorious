package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/GrappigPanda/notorious/bencode"
	"gopkg.in/redis.v3"
	"net"
	"net/http"
	"strconv"
	"strings"
)

var FIELDS = []string{"port", "uploaded", "downloaded", "left", "event", "compact"}

func compactIPPort(ip string, port string) []byte {
	res := bytes.NewBuffer(make([]byte, 0))

	intPort, err := strconv.Atoi(port)
	if err != nil {
		panic("failure1")
	}

	if err := binary.Write(res, binary.BigEndian, binary.BigEndian.Uint32(net.ParseIP(ip).To4())); err != nil {
		panic("failure0")
	}

	err = binary.Write(res, binary.BigEndian, uint16(intPort))
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

		ret.Write(compactIPPort(ip, port))
	}

	return ret.Bytes()
}

func formatResponseData(ips []string, data *announceData) string {
	compactPeerList := CompactAllPeers(ips)
	return EncodeResponse(compactPeerList, data)
}

func encodeKV(key string, value string) string {
	if value[0] == 'i' {
		return fmt.Sprintf("%s%s", bencode.EncodeByteString(key), value)
	}
	return fmt.Sprintf("%s%s", bencode.EncodeByteString(key), bencode.EncodeByteString(value))
}

func EncodeResponse(ipport []byte, data *announceData) string {
	ret := "d"
	ret += encodeKV("complete", bencode.EncodeInt(1))
	ipstr := ""

	count := 0

	ipstr += string(ipport)

	ret += encodeKV("incomplete", bencode.EncodeInt(count))
	ret += "5:peers"
	ret += strconv.Itoa(count)
	ret += ":"
	ret += ipstr
	ret += "e"
	return ret
}

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

func CreateNewTorrentKey(client *redis.Client, key string) {
	// CreateNewTorrentKey creates a new key. By default, it adds a member
	// ":ip". I don't think this ought to ever be generalized, as I just want
	// Redis to function in one specific way in notorious.

	// TODO(ian): You might want to set this explicitly in parameters
	// value := *TorrentRequestData
	client.SAdd(key, "ip")
}

func createIpPortPair(value *announceData) string {
	// createIpPortPair creates a string formatted ("%s:%s", value.ip,
	// value.port) looking like so: "127.0.0.1:6886" and returns this value.
	return fmt.Sprintf("%s:%s", value.ip, value.port)
}
