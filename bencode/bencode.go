package bencode

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func EncodeInt(x int) string {
	// Encode an an int to a bencoded integer: "i<x>e"
	return fmt.Sprintf("i%de", x)
}

func EncodeList(items []string) string {
	// Encode a list of items (`items`) to a bencoded list:
	// "l<item1><item2><...><itemN>e"
	tmp := "l"
	for i := range items {
		if items[i][0] == 'l' || items[i][0] == 'd' {
			tmp = fmt.Sprintf("%s%s", tmp, items[i])
		} else {
			tmp = fmt.Sprintf("%s%s", tmp, EncodeByteString(items[i]))
		}
	}
	tmp = fmt.Sprintf("%se", tmp)
	return tmp
}

func EncodeDictionary(kvpairs []string) (retdict string) {
	// Take a list of bencoded KVpairs and return a bencoded dictionary.

	retdict = "d"
	for i := range kvpairs {
		retdict += kvpairs[i]
	}
	retdict += "e"

	return
}

func EncodeByteString(key string) string {
	// Encode a string to <key length>:<key>
	return fmt.Sprintf("%d:%s", utf8.RuneCountInString(key), key)
}

func EncodePeerList(peers []string) (retlist string) {
	// Handles peer list creation for non-compact responses. Mostly deprecated
	// for most torrent clients nowadays as compact is the default. Returns a
	// bencoded list of bencoded dictionaries containing "peer id", "ip",
	// "port": "ld7:peer id20:<peer id>2:ip9:<127.0.0.1>4:port4:7878ee"
	// peers contains a ip:port

	var tmpDict []string

	for i := range peers {
		var tmp []string
		peerSplit := strings.Split(peers[i], ":")

		// TODO(ian): Figure out an actual way to do peer id.
		tmp = append(tmp, EncodeKV("peer id", "11111111111111111111"))
		tmp = append(tmp, EncodeKV("ip", peerSplit[0]))
		tmp = append(tmp, EncodeKV("port", peerSplit[1]))

		tmpDict = append(tmpDict, EncodeDictionary(tmp))
	}

	peerList := EncodeList(tmpDict)
	peerList = EncodeKV("peers", peerList)
	retlist = fmt.Sprintf("d%se", peerList)

	return
}

func EncodeKV(key string, value string) string {
	key = EncodeByteString(key)
	if value[0] == 'i' || value[0] == 'l' || value[0] == 'd' {
		value = value
	} else {
		value = EncodeByteString(value)
	}
	return fmt.Sprintf("%s%s", key, value)
}
