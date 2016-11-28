package reaper

import (
	"fmt"
	"github.com/GrappigPanda/notorious/database/mysql"
	"github.com/GrappigPanda/notorious/peerStore/redis"
	"gopkg.in/redis.v3"
	"strconv"
	"strings"
	"sync"
	"time"
)

func reapInfoHash(c *redis.Client, infoHash string, out chan int, currTime int64) {
	// Fan-out method to reap peers who have outlived their TTL.
	keys, err := c.SMembers(infoHash).Result()
	if err != nil {
		return
	}

	count := 0

	for i := range keys {
		if x := strings.Split(keys[i], ":"); len(x) != 3 {
			c.SRem(infoHash, keys[i])
		} else {
			endTime, err := convertTimeToUnixTimeStamp(x[2])
			if err != nil {
				panic("alskdj")
			}

			if currTime >= endTime {
				c.SRem(infoHash, keys[i])
				count++
			}
		}

	}

	out <- count
}

func convertTimeToUnixTimeStamp(time string) (endTime int64, err error) {
	return strconv.ParseInt(time, 10, 64)
}

func reapPeers(currTime int64) int {
	// Fans out each info in `keys *` from the Redis DB to the `reapInfoHash`
	// function.
	client := redisPeerStore.OpenClient()

	keys, err := getAllKeys(client, "*")
	if err != nil {
		return 0
	}

	out := make(chan int, 100000)
	timeout := make(chan bool, 5)

	peersReaped := 0

	go func() {
		select {
		case count := <-out:
			peersReaped += count
		}
	}()

	var wg sync.WaitGroup

	for i := range keys {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			reapInfoHash(client, keys[i], out, currTime)
			wg.Done()
		}(&wg)
	}

	wg.Wait()
	timeout <- true

	return peersReaped
}

// StartReapingScheduler hoists a seperate timer process and at the end of each
// time countdown, calls the reaper to reap old peers
func StartReapingScheduler(waitTime time.Duration) {
	// The timer which sets off the peer reaping every `waitTime` seconds.
	reapedPeers := 0
	go func() {
		for {
			// Handle any other cleanup or Notorious-related functions
			c := redisPeerStore.OpenClient()
			_, err := c.Ping().Result()
			if err != nil {
				panic("No Redis instance detected. If deploying without Docker, install redis-server")
			}

			infoHash := new(string)
			name := new(string)
			addedBy := new(string)
			dateAdded := new(int64)

			x, err := mysql.GetWhitelistedTorrents(nil)
			for x.Next() {
				x.Scan(infoHash, name, addedBy, dateAdded)
				redisPeerStore.CreateNewTorrentKey(nil, *infoHash)
			}

			// Start the actual peer reaper.
			time.Sleep(waitTime)
			fmt.Println("Starting peer reaper")
			reapedPeers += reapPeers(int64(time.Now().UTC().Unix()))
			fmt.Printf("%v peers reaped total\n", reapedPeers)
		}
	}()
}
