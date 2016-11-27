package announce

import (
	"fmt"
	r "github.com/GrappigPanda/notorious/kvStoreInterfaces"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// ParseAnnounceData handles getting the annunce data from a remote client and
// parses it into an acceptable data structure.
func (a *AnnounceData) ParseAnnounceData(req *http.Request) (err error) {
	query := req.URL.Query()

	a.RequestContext = requestAppContext{
		dbConn:    nil,
		Whitelist: false,
	}

	a.InfoHash = ParseInfoHash(query.Get("InfoHash"))
	if a.InfoHash == "" {
		err = fmt.Errorf("No InfoHash provided.")
		return
	}
	if strings.Contains(req.RemoteAddr, ":") {
		a.IP = strings.Split(req.RemoteAddr, ":")[0]
	} else {
		a.IP = query.Get(req.RemoteAddr)
	}
	a.PeerID = query.Get("peer_id")

	a.Port, err = GetInt(query, "port")
	if err != nil {
		return fmt.Errorf("Failed to get port")
	}
	a.Downloaded, err = GetInt(query, "downloaded")
	if err != nil {
		err = fmt.Errorf("Failed to get downloaded byte count.")
		return
	}
	a.Uploaded, err = GetInt(query, "uploaded")
	if err != nil {
		err = fmt.Errorf("Failed to get uploaded byte count.")
		return
	}
	a.Left, err = GetInt(query, "left")
	if err != nil {
		err = fmt.Errorf("Failed to get remaining byte count.")
		return
	}
	a.Numwant, err = GetInt(query, "numwant")
	if err != nil {
		a.Numwant = 0
	}
	if x := query.Get("compact"); x != "" {
		a.Compact, err = strconv.ParseBool(x)
		if err != nil {
			a.Compact = false
		}
	}
	a.Event = query.Get("event")
	if a.Event == " " || a.Event == "" {
		a.Event = "started"
	}

	a.RequestContext.redisClient = r.OpenClient()

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
func (a *AnnounceData) StartedEventHandler() (err error) {
	// Called upon announce when a client starts a download or creates a new
	// torrent on the tracker. Adds a user to incomplete list in redis.
	err = nil

	if !a.infoHashExists() && a.RequestContext.Whitelist {
		err = fmt.Errorf("Torrent not authorized for use")
		return
	} else if !a.infoHashExists() && !a.RequestContext.Whitelist {
		// If the info hash isn't in redis and we're not Whitelisting, add it
		// to Redis.
		a.createInfoHashKey()
	}

	keymember := ""
	ipport := ""

	if !(a.Left == 0) {
		keymember = fmt.Sprintf("%s:incomplete", a.InfoHash)
		ipport = fmt.Sprintf("%s:%d", a.IP, a.Port)
	} else {
		keymember = fmt.Sprintf("%s:complete", a.InfoHash)
		ipport = fmt.Sprintf("%s:%d", a.IP, a.Port)
	}

	r.RedisSetKeyVal(a.RequestContext.redisClient, keymember, ipport)
	if r.RedisSetKeyIfNotExists(a.RequestContext.redisClient, keymember, ipport) {
		fmt.Printf("Adding host %s to %s\n", ipport, keymember)
	}

	return
}

// StoppedEventHandler Called upon announce whenever a client attempts to shut-down gracefully.
// Ensures that the client is removed from complete/incomplete lists.
// TODO(ian): This is what happened whenever the torrent client shuts down
// gracefully, so we need to call the mysql backend and store the info and
// remove the ipport from completed/incomplete redis kvs
func (a *AnnounceData) StoppedEventHandler() {

	if a.infoHashExists() {
		a.removeFromKVStorage(a.Event)
	} else {
		return
	}
}

// CompletedEventHandler Called upon announce when a client finishes a download. Removes the
// client from incomplete in redis and places their peer info into
// complete.
func (a *AnnounceData) CompletedEventHandler() {

	if !a.infoHashExists() {
		a.createInfoHashKey()
	} else {
		a.removeFromKVStorage("incomplete")
	}

	keymember := fmt.Sprintf("%s:complete", a.InfoHash)
	// TODO(ian): DRY!
	ipport := fmt.Sprintf("%s:%s", a.IP, a.Port)
	if r.RedisSetKeyIfNotExists(a.RequestContext.redisClient, keymember, ipport) {
		fmt.Printf("Adding host %s to %s:complete\n", ipport, a.InfoHash)
	}
}

func (a *AnnounceData) removeFromKVStorage(subkey string) {
	// Remove the subkey from the kv storage.
	ipport := fmt.Sprintf("%s:%d", a.IP, a.Port)
	keymember := fmt.Sprintf("%s:%s", a.InfoHash, subkey)

	fmt.Printf("Removing host %s from %v\n", ipport, keymember)
	r.RedisRemoveKeysValue(a.RequestContext.redisClient, keymember, ipport)
}

func (a *AnnounceData) infoHashExists() bool {
	return r.RedisGetBoolKeyVal(nil, a.InfoHash)
}

func (a *AnnounceData) createInfoHashKey() {
	r.CreateNewTorrentKey(nil, a.InfoHash)
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
