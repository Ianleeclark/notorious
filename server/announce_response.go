package server

import (
	"bytes"
	"fmt"
	"github.com/GrappigPanda/notorious/bencode"
	"github.com/GrappigPanda/notorious/database"
	"github.com/GrappigPanda/notorious/server/peerStore"
	"net"
	"strconv"
	"strings"
)

// AnnounceResponseFailure Models the failure response sent on tracker
// failures.
type AnnounceResponseFailure struct {
	failure string
}

// AnnounceResponse models the response sent to peers
type AnnounceResponse struct {
	interval   int // Interval in seconds a client should wait |.| messages
	trackerID  string
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

	portCompact := []byte{byte(portInt >> 8), byte(portInt)}
	b.Write(portCompact)

	return
}

// CompactAllPeers Comapcts all of the peers according to BEP 23
func CompactAllPeers(ipport []string) []byte {
	var ret bytes.Buffer
	for i := range ipport {
		sz := strings.Split(ipport[i], ":")
		ip := sz[0]
		port := sz[1]

		err := compactIPPort(&ret, ip, port)
		if err != nil {
			panic(err)
		}
	}

	return ret.Bytes()
}

func formatResponseData(ips []string, data *announceData) string {
	return EncodeResponse(ips, data)
}

// EncodeResponse groups all of the peer-requested data into a nice bencoded
// string that we respond with.
func EncodeResponse(ipport []string, data *announceData) (resp string) {
	ret := ""
	completeCount := len(peerStore.RedisGetKeyVal(data.info_hash))
	incompleteCount := len(peerStore.RedisGetKeyVal(data.info_hash))
	ret += bencode.EncodeKV("complete", bencode.EncodeInt(completeCount))

	ret += bencode.EncodeKV("incomplete", bencode.EncodeInt(incompleteCount))
	if data.compact || !data.compact {
		ipstr := string(CompactAllPeers(ipport))
		ret += bencode.EncodeKV("peers", ipstr)
	} else {
		return bencode.EncodePeerList(ipport)
	}

	return fmt.Sprintf("d%se", ret)
}

func formatScrapeResponse(torrent *db.Torrent) string {
	subdir := fmt.Sprintf("d%s%s%s%s%s%se",
		bencode.EncodeByteString("complete"),
		bencode.EncodeInt(int(torrent.Seeders)),

		bencode.EncodeByteString("downloaded"),
		bencode.EncodeInt(int(torrent.Downloaded)),

		bencode.EncodeByteString("incomplete"),
		bencode.EncodeInt(int(torrent.Leechers)),
	)

	fileList := bencode.EncodeKV(
		bencode.EncodeByteString(torrent.InfoHash),
		subdir,
	)

	return fmt.Sprintf("d%se", fileList)
}

func createFailureMessage(msg string) string {
	return fmt.Sprintf("d%se", bencode.EncodeKV("failure reason", msg))
}
