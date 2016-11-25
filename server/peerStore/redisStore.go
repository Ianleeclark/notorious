package peerStore

import (
    r "github.com/GrappigPanda/notorious/kvStoreInterfaces"
	"gopkg.in/redis.v3"
)

type RedisStore struct {
    client *redis.Client
}

func (p *RedisStore) SetKeyIfNotExists(key, value string) (retval bool) {
    return r.RedisSetKeyIfNotExists(p.client, key, value)
}

func (p *RedisStore) SetKV(key, value string) {
    r.RedisSetKeyVal(p.client, key, value)
}

func (p *RedisStore) RemoveKV(key, value string) {
    // TODO(ian): Refactor this so we don't have to delete a value from a key
    if value != "" || value == "" {
        r.RedisRemoveKeysValue(p.client, key, value)
    }
}

func (p *RedisStore) KeyExists(key string) (retval bool) {
    return r.RedisGetBoolKeyVal(key)
}

func (p *RedisStore) GetKeyVal(key string) []string {
    return r.RedisGetKeyVal(key)
}

func (p *RedisStore) GetAllPeers(key string) []string {
    return r.RedisGetAllPeers(key)
}

func (p *RedisStore) SetIPMember(infoHash, ipPort string) (retval int) {
    return r.RedisSetIPMember(infoHash, ipPort)
}
