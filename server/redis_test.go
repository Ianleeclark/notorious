package server

import (
    "testing"
)

var DATA = announceData{
    info_hash: "12345123451234512345",
    peer_id: "12345123451234512345",
    ip: "127.0.0.1",
    event: "STARTED",
    port: 6767,
    uploaded: 1024,
    downloaded: 512,
    left: 0,
    numwant: 30,
    compact: true,
    redisClient: OpenClient(),
}

func TestRedisSetIPMember(t *testing.T) {
    ret := RedisSetIPMember(&DATA)

    expectedReturn := 1

    if ret != expectedReturn {
        t.Fatalf("Expected %v, got %v", expectedReturn, ret)
    }
}

func TestRedisSetKeyVal(t *testing.T) {
    RedisSetKeyVal(DATA.redisClient, "test:1234", "1024")

    ret, _ := DATA.redisClient.SMembers("test:1234").Result()

    expectedReturn := ">1"

    if len(ret) == 0 {
        t.Fatalf("Expected %v, got %v", expectedReturn, len(ret))
    }
}

func TestRedisGetKeyVal(t *testing.T) {
    DATA.redisClient.SAdd("RedisGetKeyValTest:1024:complete", "1024")
    ret := RedisGetKeyVal(DATA.redisClient, "RedisGetKeyValTest:1024", &DATA)
    expectedReturn := ">1"

    if len(ret) == 0 {
        t.Fatalf("Expected %v, got %v", expectedReturn, len(ret))
    }
}

func TestRedisGetKeyValNoPreexistKey(t *testing.T) {
    DATA.redisClient.SAdd("RedisGetKeyValTest:1025", "1024")
    ret := RedisGetKeyVal(DATA.redisClient, "RedisGetKeyValTest:1025", &DATA)
    expectedReturn := 0

    if len(ret) != expectedReturn {
        t.Fatalf("Expected %v, got %v", expectedReturn, len(ret))
    }
}

