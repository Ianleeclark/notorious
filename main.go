package main

import (
	"github.com/GrappigPanda/notorious/announce/impl"
	"github.com/GrappigPanda/notorious/config"
	"github.com/GrappigPanda/notorious/database"
	"github.com/GrappigPanda/notorious/reaper"
	"github.com/GrappigPanda/notorious/server"
	"time"
)

// Init handles initialziation of the server.
func init() {
	config := config.LoadConfig()
	db.InitDB(&config)
	go reaper.StartReapingScheduler(1 * time.Minute)
}

func main() {
	config := config.LoadConfig()

	postgresCatcher := catcherImpl.NewPostgresCatcher(config)
	postgresCatcher.HandleNewTorrent()
	server.RunServer(postgresCatcher.GetRSSNotifier())
}
