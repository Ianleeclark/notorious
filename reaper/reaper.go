package reaper

import (
	"fmt"
	"gopkg.in/redis.v3"
	"time"
)

func reapInfoHash(c *redis.Client, infoHash string, out chan int) {
	// Fan-out method to reap peers who have outlived their TTL.
	keys, err := c.SMembers(infoHash).Result()
	if err != nil {
		panic(err)
	}

	count := 0
	for i := range keys {
		if c.TTL(keys[i]).Val() <= 0 {
			if c.SRem(infoHash, keys[i]).Val() == 1 {
				count += 1
			}
		}
	}

	out <- count
}

func reapPeers() (peersReaped int) {
	// Fans out each info in `keys *` from the Redis DB to the `reapInfoHash`
	// function.
	client := OpenClient()

	keys, err := getAllKeys(client, "*")
	if err != nil {
		panic(err)
	}

	out := make(chan int)
	for i := range keys {
		go reapInfoHash(client, keys[i], out)
		peersReaped += <-out
	}

	return
}

func StartReapingScheduler(waitTime time.Duration) {
	// The timer which sets off the peer reaping every `waitTime` seconds.
	reapedPeers := 0
	go func() {
		for {
			time.Sleep(waitTime)
			fmt.Println("Starting peer reaper")
			reapedPeers += reapPeers()
			fmt.Printf("%v peers reaped total\n", reapedPeers)
		}
	}()
}
