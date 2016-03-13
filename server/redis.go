package server

import (
	"bytes"
	"fmt"
	"gopkg.in/redis.v3"
)

func RedisSetKeyVal(client *redis.Client, keymember string, value string) {
	// RedisSetKeyVal sets a key:member's value to value. Returns nothing as of
	// yet.
	client.SAdd(keymember, value)
}

func RedisGetKeyVal(client *redis.Client, key string, value *announceData) []string {
	// RedisGetKeyVal retrieves a value from the Redis store by looking up the
	// provided key. If the key does not yet exist, we create the key in the KV
	// storage or if the value is empty, we add the current requester to the
	// list.
	keymember := concatenateKeyMember(key, "ip")

	val, err := client.SMembers(keymember).Result()
	if err != nil {
		// Fail because the key doesn't exist in the KV storage.
		CreateNewTorrentKey(client, keymember)
	}

	// If no keys yet exist in the KV storage.
	if len(val) == 0 {
		RedisSetKeyVal(client, keymember, createIpPortPair(value))
	}

	return val
}

func RedisGetCount(c *redis.Client, info_hash string, member string) (retval []string, err error) {
	// A generic function which is used to retrieve either the complete count
	// or the incomplete count for a specified `info_hash`.
	keymember := concatenateKeyMember(info_hash, member)

	retval, err = c.SMembers(keymember).Result()
	if err != nil {
		// TODO(ian): Add actual error checking here.
		err = fmt.Errorf("The info hash %s with member %s doesn't exist", info_hash, member)
	}

	return
}

func RedisGetBoolKeyVal(client *redis.Client, key string, value interface{}) bool {
	_, err := client.Get(key).Result()

	return err != nil
}

func CreateNewTorrentKey(client *redis.Client, key string) {
	// CreateNewTorrentKey creates a new key. By default, it adds a member
	// ":ip". I don't think this ought to ever be generalized, as I just want
	// Redis to function in one specific way in notorious.

	// TODO(ian): You might want to set this explicitly in parameters
	// value := *TorrentRequestData
	client.SAdd(key, "ip")
}

func concatenateKeyMember(key string, member string) string {
	var buffer bytes.Buffer

	buffer.WriteString(key)
	buffer.WriteString(":")
	buffer.WriteString(member)

	return buffer.String()
}

func createIpPortPair(value *announceData) string {
	// createIpPortPair creates a string formatted ("%s:%s", value.ip,
	// value.port) looking like so: "127.0.0.1:6886" and returns this value.
	return fmt.Sprintf("%s:%s", value.ip, value.port)
}
