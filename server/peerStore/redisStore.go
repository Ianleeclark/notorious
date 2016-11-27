package peerStore

import (
	r "github.com/GrappigPanda/notorious/kvStoreInterfaces"
	"gopkg.in/redis.v3"
)

// RedisStore represents the implementation of a `PeerStore` object.
type RedisStore struct {
	client *redis.Client
}

// SetKeyIfNotExists wraps around the generic RedisSetKeyIfNotExists function
func (p *RedisStore) SetKeyIfNotExists(key, value string) (retval bool) {
	return r.RedisSetKeyIfNotExists(p.client, key, value)
}

// SetKV wraps around the generic `RedisSetKeyVal` function
func (p *RedisStore) SetKV(key, value string) {
	r.RedisSetKeyVal(p.client, key, value)
}

// RemoveKV wraps around the specific `RedisRemoveKeysValue` function
func (p *RedisStore) RemoveKV(key, value string) {
	// TODO(ian): Refactor this so we don't have to delete a value from a key
	if value != "" || value == "" {
		r.RedisRemoveKeysValue(p.client, key, value)
	}
}

// KeyExists wraps around the specific `RedisGetBoolKeyVal` function
func (p *RedisStore) KeyExists(key string) (retval bool) {
	return r.RedisGetBoolKeyVal(p.client, key)
}

// GetKeyVal wraps around the specific `RedisGetKeyVal` function
func (p *RedisStore) GetKeyVal(key string) []string {
	return r.RedisGetKeyVal(p.client, key)
}

// GetAllPeers wraps around the specific `RedisGetAllPeers` function
func (p *RedisStore) GetAllPeers(key string) []string {
	return r.RedisGetAllPeers(p.client, key)
}

// SetIPMember wraps around the specific `RedisSetIPMember` function
func (p *RedisStore) SetIPMember(infoHash, ipPort string) (retval int) {
	return r.RedisSetIPMember(p.client, infoHash, ipPort)
}
