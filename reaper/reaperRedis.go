package reaper

import (
	"gopkg.in/redis.v3"
)

func OpenClient() (client *redis.Client) {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return
}
