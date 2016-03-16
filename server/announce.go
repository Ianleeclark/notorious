package server

import (
	"fmt"
	"gopkg.in/redis.v3"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (a *announceData) parseAnnounceData(req *http.Request) (err error) {
	query := req.URL.Query()
	a.info_hash = ParseInfoHash(query.Get("info_hash"))
	if a.info_hash == "" {
		err = fmt.Errorf("No info_hash provided.")
		return
	}
	if strings.Contains(req.RemoteAddr, ":") {
		a.ip = strings.Split(req.RemoteAddr, ":")[0]
	} else {
		a.ip = query.Get(req.RemoteAddr)
	}
	a.peer_id = query.Get("peer_id")

	a.port, err = GetInt(query, "port")
	if err != nil {
		return fmt.Errorf("Failed to get port")
	}
	a.downloaded, err = GetInt(query, "downloaded")
	if err != nil {
		err = fmt.Errorf("Failed to get downloaded byte count.")
		return
	}
	a.uploaded, err = GetInt(query, "uploaded")
	if err != nil {
		err = fmt.Errorf("Failed to get uploaded byte count.")
		return
	}
	a.left, err = GetInt(query, "left")
	if err != nil {
		err = fmt.Errorf("Failed to get remaining byte count.")
		return
	}
	a.numwant, err = GetInt(query, "numwant")
	if err != nil {
		a.numwant = 0
	}
	if x := query.Get("compact"); x != "" {
		a.compact, err = strconv.ParseBool(x)
		if err != nil {
			a.compact = false
		}
	}
	a.event = query.Get("event")

	return
}

func GetInt(u url.Values, key string) (ui uint64, err error) {
	if x := u.Get(key); x == "" {
		err = fmt.Errorf("Failed to locate the key in the url.")
	} else {
		ui, err = strconv.ParseUint(x, 10, 64)
		if err != nil {
			err = fmt.Errorf("Failed to parse uint from the key")
			return
		}
		return
	}
	return
}

func (data *announceData) StartedEventHandler(c *redis.Client) {
	// Called upon announce when a client starts a download or creates a new
	// torrent on the tracker. Adds a user to incomplete list in redis.

	if !data.infoHashExists(c) {
		data.createInfoHashKey(c)
	}

	keymember := ""
	ipport := ""

	if !(data.left == 0) {
		keymember = fmt.Sprintf("%s:incomplete", data.info_hash)
		ipport = fmt.Sprintf("%s:%d", data.ip, data.port)
	} else {
		keymember = fmt.Sprintf("%s:complete", data.info_hash)
		ipport = fmt.Sprintf("%s:%d", data.ip, data.port)
	}

	RedisSetKeyVal(c, keymember, ipport)
	if RedisSetKeyIfNotExists(c, keymember, ipport) {
		fmt.Printf("Adding host %s to %s\n", ipport, keymember)
	}
}

func (data *announceData) StoppedEventHandler(c *redis.Client) {
	// Called upon announce whenever a client attempts to shut-down gracefully.
	// Ensures that the client is removed from complete/incomplete lists.

	// TODO(ian): This is what happend whenever the torrent client shuts down
	// gracefully, so we need to call the mysql backend and store the info and
	// remove the ipport from completed/incomplete redis kvs

	if data.infoHashExists(c) {
		// TODO(ian): THis is not done!
		data.removeFromKVStorage(c, data.event)
	} else {
		return
	}
}

func (data *announceData) CompletedEventHandler(c *redis.Client) {
	// Called upon announce when a client finishes a download. Removes the
	// client from incomplete in redis and places their peer info into
	// complete.

	if !data.infoHashExists(c) {
		data.createInfoHashKey(c)
	} else {
		fmt.Printf("Removing host %s:%v to %s:incomplete\n", data.ip, data.port, data.info_hash)
		data.removeFromKVStorage(c, "incomplete")
	}

	keymember := fmt.Sprintf("%s:complete", data.info_hash)
	// TODO(ian): DRY!
	ipport := fmt.Sprintf("%s:%s", data.ip, data.port)
	if RedisSetKeyIfNotExists(c, keymember, ipport) {
		fmt.Printf("Adding host %s to %s:complete\n", ipport, data.info_hash)
	}
}

func (data *announceData) removeFromKVStorage(c *redis.Client, subkey string) {
	// Remove the subkey from the kv storage.

	ipport := fmt.Sprintf("%s:%s", data.ip, data.port)
	keymember := fmt.Sprintf("%s:%s", data.info_hash, subkey)
	RedisRemoveKeysValue(c, keymember, ipport)
}

func (data *announceData) infoHashExists(c *redis.Client) bool {
	return RedisGetBoolKeyVal(c, data.info_hash, data)
}

func (data *announceData) createInfoHashKey(c *redis.Client) {
	CreateNewTorrentKey(c, data.info_hash)
}

func ParseInfoHash(s string) string {
	return fmt.Sprintf("%x", s)
}

func decodeQueryURL(s string) url.Values {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	m, _ := url.ParseQuery(u.RawQuery)
	return m
}
