package peerStore

import (
	"bytes"
	"fmt"
	"gopkg.in/redis.v3"
	"strings"
	"time"
)

type RedisStore struct {
    client *redis.Client
}

// EXPIRETIME signifies how long a peer will live under the specified info_hash
// until the reaper removes it.
var EXPIRETIME int64 = 5 * 60

// OpenClient opens a connection to redis.
func OpenClient() (client *redis.Client) {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return
}

func (p *RedisStore) SetKeyIfNotExists(key, value string) (retval bool) {
    return redisSetKeyIfNotExists(p.client, key, value)
}

func (p *RedisStore) SetKV(key, value string) {
    redisSetKeyVal(p.client, key, value)
}

func (p *RedisStore) RemoveKV(key, value string) {
    // TODO(ian): Refactor this so we don't have to delete a value from a key
    if value != "" || value == "" {
        redisRemoveKeysValue(p.client, key, value)
    }
}

func (p *RedisStore) KeyExists(key string) (retval bool) {
    return redisGetBoolKeyVal(key)
}

func (p *RedisStore) GetKeyVal(key string) []string {
    return redisGetKeyVal(key)
}

// RedisSetIPMember sets a key as a member of an infohash and sets a timeout.
func redisSetIPMember(infoHash, ipPort string) (retval int) {
	c := OpenClient()

	keymember := concatenateKeyMember(infoHash, "ip")

	currTime := int64(time.Now().UTC().AddDate(0, 0, 1).Unix())

    key := fmt.Sprintf("%s:%v", ipPort, currTime)

	if err := c.SAdd(keymember, key).Err(); err != nil {
		retval = 0
		panic("Failed to add key")

	} else {
		retval = 1
	}

	return
}

// RedisSetKeyVal Sets a key to the specified value. Used mostly with adding a
// peer into an info_hash
func redisSetKeyVal(c *redis.Client, keymember string, value string) {
	// RedisSetKeyVal sets a key:member's value to value. Returns nothing as of
	// yet.
	currTime := int64(time.Now().UTC().Unix())
	currTime += EXPIRETIME
	value = fmt.Sprintf("%v:%v", value, currTime)

	if sz := strings.Split(value, ":"); len(sz) >= 1 {
		// If the value being added can be converted to an int, it is a ip:port key
		// and we can set an expiration on it.
		c.SAdd(keymember, value)
	}
}

// RedisGetKeyVal Lookup a peer in the specified infohash at `key`
func redisGetKeyVal(key string) []string {
	c := OpenClient()

	// RedisGetKeyVal retrieves a value from the Redis store by looking up the
	// provided key. If the key does not yet exist, we create the key in the KV
	// storage or if the value is empty, we add the current requester to the
	// list.
	keymember := concatenateKeyMember(key, "complete")

	val, err := c.SMembers(keymember).Result()
	if err != nil {
		// Fail because the key doesn't exist in the KV storage.
		createNewTorrentKey(keymember)
	}

	return val
}

// RedisGetAllPeers fetches all peers from the info_hash at `key`
func redisGetAllPeers(key string) []string {
	c := OpenClient()

	keymember := concatenateKeyMember(key, "complete")

	val, err := c.SRandMemberN(keymember, 30).Result()
	if err != nil {
		// Fail because the key doesn't exist in the KV storage.
		createNewTorrentKey(keymember)
	}

	if len(val) == 30 {
		return val
	}

	keymember = concatenateKeyMember(key, "incomplete")

	val2, err := c.SRandMemberN(keymember, int64(30-len(val))).Result()
	if err != nil {
		panic("Failed to get incomplete peers for")
	} else {
		val = append(val, val2...)
	}

	return val
}

// RedisGetCount counts all of the peers at `info_hash`
func redisGetCount(c *redis.Client, info_hash string, member string) (retval int, err error) {
	// A generic function which is used to retrieve either the complete count
	// or the incomplete count for a specified `info_hash`.
	keymember := concatenateKeyMember(info_hash, member)

	x, err := c.SMembers(keymember).Result()
	if err != nil {
		// TODO(ian): Add actual error checking here.
		err = fmt.Errorf("The info hash %s with member %s doesn't exist", info_hash, member)
	}

	retval = len(x)
	return
}

// RedisGetBoolKeyVal Checks if a `key` exists
func redisGetBoolKeyVal(key string) bool {
	c := OpenClient()
	ret, _ := c.Exists(key).Result()

	return ret
}

// RedisSetKeyIfNotExists Set a key if it doesn't exist.
func redisSetKeyIfNotExists(c *redis.Client, keymember string, value string) (rv bool) {
	rv = redisGetBoolKeyVal(keymember)
	if !rv {
		redisSetKeyVal(c, keymember, value)
	}
	return
}

// RedisRemoveKeysValue Remove a `value` from `key` in the redis kv storage. `key` is typically
// a keymember of info_hash:(in)complete and the value is typically the
// ip:port concatenated.
func redisRemoveKeysValue(c *redis.Client, key string, value string) {
	c.SRem(key, value)
}

// CreateNewTorrentKey creates a new key. By default, it adds a member
// ":ip". I don't think this ought to ever be generalized, as I just want
// Redis to function in one specific way in notorious.
func createNewTorrentKey(key string) {
	c := OpenClient()
	c.SAdd(key, "complete", "incomplete")

}

// concatenateKeyMember concatenates the key and the member delimited by the
// character ":"
func concatenateKeyMember(key string, member string) string {
	var buffer bytes.Buffer

	buffer.WriteString(key)
	buffer.WriteString(":")
	buffer.WriteString(member)

	return buffer.String()
}

// createIPPortPair creates a string formatted ("%s:%s", value.ip,
// value.port) looking like so: "127.0.0.1:6886" and returns this value.
func createIPPortPair(ip, port string) string {
	return fmt.Sprintf("%v:%v", ip, port)
}
