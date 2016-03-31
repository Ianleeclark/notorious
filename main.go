package main

import (
	"github.com/GrappigPanda/notorious/reaper"
	"github.com/GrappigPanda/notorious/server"
	"time"
)

func main() {
	c := server.OpenClient()
	_, err := c.Ping().Result()
	if err != nil {
		panic("No Redis instance detected. If deploying without Docker, install redis-server")
	}

	go reaper.StartReapingScheduler(1 * time.Minute)
	server.RunServer()
}
