package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/GrappigPanda/notorious/bencode"
	"gopkg.in/redis.v3"
	"net"
	"strconv"
	"strings"
)

// TODO(ian): Finish crafting a response.
type AnnounceResponseFailure struct {
	failure string
}

type AnnounceResponse struct {
	interval   int // Interval in seconds a client should wait |.| messages
	trackerId  string
	complete   uint
	incomplete uint
	peers      PeerList
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

func formatResponseData(c *redis.Client, ips []string, data *announceData) string {
	compactPeerList := CompactAllPeers(ips)
	return EncodeResponse(c, compactPeerList, data)
}

func encodeKV(key string, value string) string {
	if value[0] == 'i' {
		return fmt.Sprintf("%s%s", bencode.EncodeByteString(key), value)
	}
	return fmt.Sprintf("%s%s", bencode.EncodeByteString(key), bencode.EncodeByteString(value))
}

func EncodeResponse(c *redis.Client, ipport []byte, data *announceData) (resp string) {
	ret := ""
	completeCount := len(RedisGetKeyVal(c, data.info_hash, data))
	incompleteCount := len(RedisGetKeyVal(c, data.info_hash, data))
	ret += encodeKV("complete", bencode.EncodeInt(completeCount))

	ipstr := string(ipport)

	ret += encodeKV("incomplete", bencode.EncodeInt(incompleteCount))
	if data.compact {
		ret += encodeKV("peers", ipstr)
	} else {
		// TODO(ian): Add an option if compact = 0
		return ""
	}

	resp = bencode.EncodeDictionary(ret, "")

	return
}
