package reaper

import (
	"fmt"
	"gopkg.in/redis.v3"
)

func OpenClient() (client *redis.Client) {
	// Opens a connection to the redis connection.
	// TODO(ian): Add a config option for redis host:port
	// TODO(ian): Add error checking here.
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return
}

func getAllKeys(c *redis.Client, keymember string) (keys []string, err error) {
	// getAllKeys gets all the keys for a specified `keymember`

	allKeys := c.Keys(keymember)
	keys, err = allKeys.Result()
	if err != nil {
		fmt.Errorf("Failed to reap peers")
	}

	return
}
