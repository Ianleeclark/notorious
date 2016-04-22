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

func TestCreateIpPortPair(t *testing.T) {
    expectedReturn := "127.0.0.1:6767"
    ret := createIpPortPair(&DATA)

    if expectedReturn != ret {
        t.Fatalf("Expected %v, got %v", expectedReturn, ret)
    }
}

func TestConcatenateKeyMember(t *testing.T) {
    expectedReturn := "127.0.0.1:1024"
    ret := concatenateKeyMember("127.0.0.1", "1024")

    if expectedReturn != ret {
        t.Fatalf("Expected %v, got %v", expectedReturn, ret)
    }
}

func TestCreateNewTorrentKey(t *testing.T) {
    CreateNewTorrentKey(DATA.redisClient, "testTestCreateNewTorrentKey")

    ret, err := DATA.redisClient.Exists("testTestCreateNewTorrentKey").Result()
    if err != nil {
        t.Fatalf("%v", err)
    }
    if ret != true {
        t.Fatalf("CreateNewTorrentKey:complete failed to create")
    }

    ret, err = DATA.redisClient.SIsMember("testTestCreateNewTorrentKey", "complete").Result()
    if ret != true {
        t.Fatalf("testTestCreateNewTorrentKey:complete is not a member")
    }

    ret, err = DATA.redisClient.SIsMember("testTestCreateNewTorrentKey", "incomplete").Result()
    if ret != true {
        t.Fatalf("testTestCreateNewTorrentKey:incomplete is not a member")
    }

}

func TestRedisRemoveKeyValues(t *testing.T) {
    DATA.redisClient.SAdd("TestRedisRemoveKeyVal", "Test1")
    ret, err := DATA.redisClient.SIsMember("TestRedisRemoveKeyVal", "Test1").Result()
    if err != nil {
        t.Fatalf("%v", err)
    }
    if ret != true {
        t.Fatalf("Failed in setup of TestRedisRemoveKeyValues to add a key")
    }

    RedisRemoveKeysValue(DATA.redisClient, "TestRedisRemoveKeyVal", "Test1")
    ret, err = DATA.redisClient.SIsMember("TestRedisRemoveKeyVal", "Test1").Result()
    if err != nil {
        t.Fatalf("%v", err)
    }
    if ret == true {
        t.Fatalf("RedisRemoveKeyVal failed to remove the key")
    }

}

func TestRedisGetBoolKeyVal(t *testing.T) {
    RedisSetKeyVal(DATA.redisClient, "TestRedisGetBoolKeyVal", "1024")

    expectedReturn := true
    ret := RedisGetBoolKeyVal(DATA.redisClient, "TestRedisGetBoolKeyVal", "1024")

    if ret != expectedReturn {
        t.Fatalf("Expected %v, got %v", expectedReturn, ret)
    }
}

func TestRedisSetKeyIfNotExists(t *testing.T) {
    expectedReturn := true
    ret := RedisSetKeyIfNotExists(DATA.redisClient, "TestRedisSetKeyIfNotExists", "1024")

    if ret != expectedReturn {
        t.Fatalf("Expected %v, got %v", expectedReturn, ret)
    }
}

func TestRedisSetKeyIfNotExistsPreExistingKey(t *testing.T) {
    expectedReturn := true
    RedisSetKeyVal(DATA.redisClient, "TestRedisSetKeyIfNotExists", "1024")
    ret := RedisSetKeyIfNotExists(DATA.redisClient, "TestRedisSetKeyIfNotExists", "1024")

    if ret != expectedReturn {
        t.Fatalf("Expected %v, got %v", expectedReturn, ret)
    }
}

