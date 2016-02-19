package redisManager

import (
    "gopkg.in/redis.v3"
)

func OpenClient() *redis.Client {
    client := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
        Password: "",
        DB: 0,
    })

    return client
}

func CreateNewTorrentKey(client *redis.Client, key string, value interface{}) {
    client.Set(key, "test", 0)
 }

func RedisSetKeyVal(client *redis.Client, key string, value interface{}) interface{} {
    val, err := client.Get(key).Result()
    if err != nil {
        CreateNewTorrentKey(client, key, value)
    }
    return val
}

func RedisGetKeyVal(client *redis.Client, key string, value interface{}) interface{} {
    val, _ := client.Get(key).Result()
    CreateNewTorrentKey(client, key, value)

    return val
}

func RedisGetBoolKeyVal(client *redis.Client, key string, value interface{}) bool {
    _, err := client.Get(key).Result()

    return err == nil
}



