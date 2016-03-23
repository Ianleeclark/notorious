package main

import (
	"github.com/GrappigPanda/notorious/reaper"
	"github.com/GrappigPanda/notorious/server"
	"time"
)

func main() {
	go reaper.StartReapingScheduler(5 * 60 * time.Second)
	server.RunServer()
}
