package reaper

import (
	"fmt"
	"gopkg.in/redis.v3"
)

func getAllKeys(c *redis.Client, keymember string) (keys []string, err error) {
	// getAllKeys gets all the keys for a specified `keymember`

	allKeys := c.Keys(keymember)
	keys, err = allKeys.Result()
	if err != nil {
		fmt.Errorf("Failed to reap peers")
	}

	return
}
