package server

import (
	"bytes"
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

func RedisGetBoolKeyVal(client *redis.Client, key string, value interface{}) bool {
	_, err := client.Get(key).Result()

	return err != nil
}

func concatenateKeyMember(key string, member string) string {
	var buffer bytes.Buffer

	buffer.WriteString(key)
	buffer.WriteString(":")
	buffer.WriteString(member)

	return buffer.String()
}
