package reaper

import (
	"fmt"
)

func reapInfoHash(c *redis.Client, infoHash string) (peersReaped uint) {
	// Fan-out method to reap peers who have outlived their TTL.
}

func reapPeers() {
	// Fans out each info in `keys *` from the Redis DB to the `reapInfoHash`
	// function.
	client = OpenClient()

	// For info_hash in keys, reapInfoHash(client, info_hash)
}

func StartReapingScheduler(waitTime time.Duration) {
	// The timer which sets off the peer reaping every `waitTime` seconds.
	go func() {
		for {
			time.Sleep(waitTime)
			reapPeers()
		}
	}()
}
