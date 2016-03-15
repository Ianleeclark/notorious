package server

import (
	"fmt"
	"gopkg.in/redis.v3"
	"net/url"
	"strconv"
)

func (a *announceData) parseAnnounceData(u *url.URL) (err error) {
	query := u.Query()
	a.info_hash = ParseInfoHash(query.Get("info_hash"))
	if a.info_hash == "" {
		err = fmt.Errorf("No info_hash provided.")
		return
	}
	a.ip = query.Get("ip")
	if a.ip == "" {
		return fmt.Errorf("No info_hash provided.")
	}
	a.peer_id = query.Get("peer_id")
	if a.peer_id == "" {
		return fmt.Errorf("No info_hash provided.")
	}
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
		err = fmt.Errorf("Failed to get number of peers requested.")
		return
	}
	if x := query.Get("compact"); x != "" {
		a.compact, err = strconv.ParseBool(x)
		if err != nil {
			err = fmt.Errorf("Failed to parse a boolean value from `compact`.")
			return
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

	keymember := fmt.Sprintf("%s:incomplete", data.info_hash)
	RedisSetKeyVal(c, keymember, fmt.Sprintf("%s:%d", data.ip, data.port))

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
		data.removeFromKVStorage(c, "incomplete")
	}

	keymember := fmt.Sprintf("%s:complete", data.info_hash)
	// TODO(ian): DRY!
	ipport := fmt.Sprintf("%s:%s", data.ip, data.port)
	RedisSetKeyVal(c, keymember, ipport)
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
