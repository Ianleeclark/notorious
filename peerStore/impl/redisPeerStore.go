package redisPeerStoreImpl

import (
	"github.com/GrappigPanda/notorious/peerStore/redis"
	"gopkg.in/redis.v3"
)

// RedisStore represents the implementation of a `PeerStore` object.
type RedisStore struct {
	client *redis.Client
}

// SetKeyIfNotExists wraps around the generic RedisSetKeyIfNotExists function
func (p *RedisStore) SetKeyIfNotExists(key, value string) (retval bool) {
	return redisPeerStore.SetKeyIfNotExists(p.client, key, value)
}

// SetKV wraps around the generic `SetKeyVal` function
func (p *RedisStore) SetKV(key, value string) {
	redisPeerStore.SetKeyVal(p.client, key, value)
}

// RemoveKV wraps around the specific `RemoveKeysValue` function
func (p *RedisStore) RemoveKV(key, value string) {
	// TODO(ian): Refactor this so we don't have to delete a value from a key
	if value != "" || value == "" {
		redisPeerStore.RemoveKeysValue(p.client, key, value)
	}
}

// KeyExists wraps around the specific `GetBoolKeyVal` function
func (p *RedisStore) KeyExists(key string) (retval bool) {
	return redisPeerStore.GetBoolKeyVal(p.client, key)
}

// GetKeyVal wraps around the specific `GetKeyVal` function
func (p *RedisStore) GetKeyVal(key string) []string {
	return redisPeerStore.GetKeyVal(p.client, key)
}

// GetAllPeers wraps around the specific `GetAllPeers` function
func (p *RedisStore) GetAllPeers(key string) []string {
	return redisPeerStore.GetAllPeers(p.client, key)
}

// SetIPMember wraps around the specific `SetIPMember` function
func (p *RedisStore) SetIPMember(infoHash, ipPort string) (retval int) {
	return redisPeerStore.SetIPMember(p.client, infoHash, ipPort)
}

// CreateNewTorrentKey wraps around the specific `CreateNewTorrentKey` function
func (p *RedisStore) CreateNewTorrentKey(infoHash string) {
	redisPeerStore.CreateNewTorrentKey(p.client, infoHash)
}
