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
    expectedReturn := false
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

func TestRedisGetCount(t *testing.T) {
    DATA.redisClient.SAdd("TestRedisGetCount", "Test")
    DATA.redisClient.SAdd("TestRedisGetCount:Test", "1235")
    DATA.redisClient.SAdd("TestRedisGetCount:Test", "1236")
    DATA.redisClient.SAdd("TestRedisGetCount:Test", "1237")
    DATA.redisClient.SAdd("TestRedisGetCount:Test", "1238")

    expectedReturn := 4
    ret, err := RedisGetCount(DATA.redisClient, "TestRedisGetCount", "Test")
    if err != nil {
        t.Fatalf("%v", err)
    }

    if ret != expectedReturn {
        t.Fatalf("Expected %v, got %v", expectedReturn, ret)
    }
}

func TestRedisGetAllPeers(t *testing.T) {
    DATA.redisClient.SAdd("TestRedisGetAllPeers", "complete")
    DATA.redisClient.SAdd("TestRedisGetAllPeers:complete", "1235")
    DATA.redisClient.SAdd("TestRedisGetAllPeers:complete", "1236")
    DATA.redisClient.SAdd("TestRedisGetAllPeers:complete", "1237")
    DATA.redisClient.SAdd("TestRedisGetAllPeers:complete", "1238")

    ret := RedisGetAllPeers(DATA.redisClient, "TestRedisGetAllPeers", &DATA)
    x := len(ret)

    if x != 4 {
        t.Fatalf("Expected 4 peers, got %v", x)
    }
}

func TestRedisGetAllPeersValGT30(t *testing.T) {
    DATA.redisClient.SAdd("TestRedisGetAllPeers1", "complete")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1201")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1202")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1203")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1204")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1205")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1206")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1207")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1208")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1209")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1210")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1211")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1212")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1213")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1214")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1215")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1216")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1217")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1218")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1209")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1200")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1201")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1202")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1203")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1204")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1205")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1216")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1217")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1218")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1219")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1220")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1221")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1222")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1221")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1222")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1223")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1224")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1225")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1226")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1227")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1228")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1229")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1230")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1231")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1232")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1233")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1234")
    DATA.redisClient.SAdd("TestRedisGetAllPeers1:complete", "1235")

    ret := RedisGetAllPeers(DATA.redisClient, "TestRedisGetAllPeers1", &DATA)
    x := len(ret)

    if x != 30 {
        t.Fatalf("Expected 30 peers, got %v", x)
    }
}

func TestRedisGetAllPeersValLT30(t *testing.T) {
    DATA.redisClient.SAdd("TestRedisGetAllPeers2", "complete")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:complete", "1201")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:complete", "1202")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:complete", "1203")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:complete", "1204")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:complete", "1205")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:complete", "1206")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:complete", "1207")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:complete", "1208")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:complete", "1209")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:complete", "1210")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:complete", "1211")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:complete", "1212")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:complete", "1213")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:complete", "1214")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:complete", "1215")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2", "incomplete")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1216")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1217")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1218")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1209")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1200")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1201")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1202")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1203")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1204")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1205")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1216")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1217")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1218")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1219")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1220")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1221")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1222")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1221")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1222")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1223")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1224")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1225")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1226")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1227")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1228")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1229")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1230")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1231")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1232")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1233")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1234")
    DATA.redisClient.SAdd("TestRedisGetAllPeers2:incomplete", "1235")

    ret := RedisGetAllPeers(DATA.redisClient, "TestRedisGetAllPeers2", &DATA)
    x := len(ret)

    if x != 30 {
        t.Fatalf("Expected 30 peers, got %v", x)
    }
}
