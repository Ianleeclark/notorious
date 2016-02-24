package redisManager

import (
	"gopkg.in/redis.v3"
)

func OpenClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return client
}

func CreateNewTorrentKey(client *redis.Client, key string, value interface{}) {
	// TODO(ian): You might want to set this explicitly in parameters
	// value := *TorrentRequestData

	// Here the key is the info_hash for the torrent and value is
	// the newest peer for the torrent
	client.SAdd(key, "ip")
	RedisSetKeyVal(client, key, "ip", "127.0.0.1")
}

func RedisSetKeyVal(client *redis.Client, key string, member string, value interface{}) interface{} {
	client.SAdd("12345:ip", "127.0.0.1")
	return 1
}

func RedisGetKeyVal(client *redis.Client, key string, value interface{}) interface{} {
	val, err := client.Get(key).Result()
	if err != nil {
		CreateNewTorrentKey(client, key, value)
	}

	return val
}

func RedisGetBoolKeyVal(client *redis.Client, key string, value interface{}) bool {
	_, err := client.Get(key).Result()

	return err != nil
}
