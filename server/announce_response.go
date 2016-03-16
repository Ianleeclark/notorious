package server

import (
	"bytes"
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
}

func compactIPPort(b *bytes.Buffer, ip string, port string) (err error) {
	compactIP := net.ParseIP(ip).To4()
	if compactIP == nil {
		err = fmt.Errorf("Failed to compact peer %s", ip)
		return
	}

	b.Write(compactIP)

	portInt, err := strconv.Atoi(port)
	if err != nil {
		err = fmt.Errorf("Failed to format port (%s) as an integer.", port)
		return
	}
	// All credit to whatcd's ocelot tracker. I'm too dumb to figure this out
	// on my own.
	portCompact := []byte{byte(portInt >> 8), byte(portInt)}
	b.Write(portCompact)

	return
}

func CompactAllPeers(ipport []string) []byte {
	ret := bytes.NewBuffer(make([]byte, 0))
	for i := range ipport {
		sz := strings.Split(ipport[i], ":")
		ip := sz[0]
		port := sz[1]

		err := compactIPPort(ret, ip, port)
		if err != nil {
			panic(err)
		}
	}

	return ret.Bytes()
}

func formatResponseData(c *redis.Client, ips []string, data *announceData) string {
	return EncodeResponse(c, ips, data)
}

func EncodeResponse(c *redis.Client, ipport []string, data *announceData) (resp string) {
	ret := ""
	completeCount := len(RedisGetKeyVal(c, data.info_hash, data))
	incompleteCount := len(RedisGetKeyVal(c, data.info_hash, data))
	ret += bencode.EncodeKV("complete", bencode.EncodeInt(completeCount))

	ret += bencode.EncodeKV("incomplete", bencode.EncodeInt(incompleteCount))
	if data.compact {
		ipstr := string(CompactAllPeers(ipport))
		fmt.Println(ipport)
		ret += bencode.EncodeKV("peers", ipstr)
	} else {
		return bencode.EncodePeerList(ipport)
	}

	resp = fmt.Sprintf("d%se", ret)
	fmt.Printf("Response: %s\n", resp)

	return resp
}

func createFailureMessage(msg string) string {
	return fmt.Sprintf("d%se", bencode.EncodeKV("failure reason", msg))
}
