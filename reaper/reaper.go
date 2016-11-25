package reaper

import (
	"fmt"
	"github.com/GrappigPanda/notorious/database"
	r "github.com/GrappigPanda/notorious/kvStoreInterfaces"
	"gopkg.in/redis.v3"
	"strconv"
	"strings"
	"time"
)

func reapInfoHash(c *redis.Client, infoHash string, out chan int) {
	// Fan-out method to reap peers who have outlived their TTL.
	keys, err := c.SMembers(infoHash).Result()
	if err != nil {
		panic(err)
	}

	count := 0
	currTime := int64(time.Now().UTC().Unix())

	for i := range keys {
		if x := strings.Split(keys[i], ":"); len(x) != 3 {
			c.SRem(infoHash, keys[i])

		} else {
			endTime := convertTimeToUnixTimeStamp(x[2])
			if currTime >= endTime {
				c.SRem(infoHash, keys[i])
				count++
			}
		}
	}

	out <- count
}

func convertTimeToUnixTimeStamp(time string) (endTime int64) {
	endTime, err := strconv.ParseInt(time, 10, 64)
	if err != nil {
		panic(err)
	}

	return
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

// StartReapingScheduler hoists a seperate timer process and at the end of each
// time countdown, calls the reaper to reap old peers
func StartReapingScheduler(waitTime time.Duration) {
	// The timer which sets off the peer reaping every `waitTime` seconds.
	reapedPeers := 0
	go func() {
		for {
			// Handle any other cleanup or Notorious-related functions
			c := r.OpenClient()
			_, err := c.Ping().Result()
			if err != nil {
				panic("No Redis instance detected. If deploying without Docker, install redis-server")
			}

			infoHash := new(string)
			name := new(string)
			addedBy := new(string)
			dateAdded := new(int64)

			x, err := db.GetWhitelistedTorrents()
			for x.Next() {
				x.Scan(infoHash, name, addedBy, dateAdded)
				r.CreateNewTorrentKey(*infoHash)
			}

			// Start the actual peer reaper.
			time.Sleep(waitTime)
			fmt.Println("Starting peer reaper")
			reapedPeers += reapPeers()
			fmt.Printf("%v peers reaped total\n", reapedPeers)
		}
	}()
}
