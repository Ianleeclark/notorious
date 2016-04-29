package server

import (
	"fmt"
	"github.com/GrappigPanda/notorious/database"
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
	if a.event == " " || a.event == "" {
		a.event = "started"
	}

	a.requestContext.redisClient = OpenClient()

	return
}

// GetInt converts the `key` from url.Values to a uint64
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

// StartedEventHandler handles whenever a peer sends the STARTED event to the
// tracker.
func (a *announceData) StartedEventHandler() (err error) {
	// Called upon announce when a client starts a download or creates a new
	// torrent on the tracker. Adds a user to incomplete list in redis.
	err = nil

	if !a.infoHashExists() && a.requestContext.whitelist {
		_, err := db.GetWhitelistedTorrent(a.info_hash)
		if err != nil {
			err = fmt.Errorf("Info hash %s not authorized for use", a.info_hash)
		}
	} else if !a.infoHashExists() && !a.requestContext.whitelist {
		// If the info hash isn't in redis and we're not whitelisting, add it
		// to Redis.
		a.createInfoHashKey()
	}

	keymember := ""
	ipport := ""

	if !(a.left == 0) {
		keymember = fmt.Sprintf("%s:incomplete", a.info_hash)
		ipport = fmt.Sprintf("%s:%d", a.ip, a.port)
	} else {
		keymember = fmt.Sprintf("%s:complete", a.info_hash)
		ipport = fmt.Sprintf("%s:%d", a.ip, a.port)
	}

	RedisSetKeyVal(a.requestContext.redisClient, keymember, ipport)
	if RedisSetKeyIfNotExists(a.requestContext.redisClient, keymember, ipport) {
		fmt.Printf("Adding host %s to %s\n", ipport, keymember)
	}

	return
}

// StoppedEventHandler Called upon announce whenever a client attempts to shut-down gracefully.
// Ensures that the client is removed from complete/incomplete lists.
// TODO(ian): This is what happened whenever the torrent client shuts down
// gracefully, so we need to call the mysql backend and store the info and
// remove the ipport from completed/incomplete redis kvs
func (a *announceData) StoppedEventHandler() {

	if a.infoHashExists() {
		a.removeFromKVStorage(a.event)
	} else {
		return
	}
}

// CompletedEventHandler Called upon announce when a client finishes a download. Removes the
// client from incomplete in redis and places their peer info into
// complete.
func (a *announceData) CompletedEventHandler() {

	if !a.infoHashExists() {
		a.createInfoHashKey()
	} else {
		a.removeFromKVStorage("incomplete")
	}

	keymember := fmt.Sprintf("%s:complete", a.info_hash)
	// TODO(ian): DRY!
	ipport := fmt.Sprintf("%s:%s", a.ip, a.port)
	if RedisSetKeyIfNotExists(a.requestContext.redisClient, keymember, ipport) {
		fmt.Printf("Adding host %s to %s:complete\n", ipport, a.info_hash)
	}
}

func (a *announceData) removeFromKVStorage(subkey string) {
	// Remove the subkey from the kv storage.
	ipport := fmt.Sprintf("%s:%d", a.ip, a.port)
	keymember := fmt.Sprintf("%s:%s", a.info_hash, subkey)

	fmt.Printf("Removing host %s from %v\n", ipport, keymember)
	RedisRemoveKeysValue(a.requestContext.redisClient, keymember, ipport)
}

func (a *announceData) infoHashExists() bool {
	return RedisGetBoolKeyVal(a.requestContext.redisClient, a.info_hash)
}

func (a *announceData) createInfoHashKey() {
	CreateNewTorrentKey(a.requestContext.redisClient, a.info_hash)
}

// ParseInfoHash parses the encoded info hash. Such a simple solution for a
// problem I hate more than koalas.
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
