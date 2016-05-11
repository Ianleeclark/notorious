package main

import (
	"github.com/GrappigPanda/notorious/database"
	"github.com/GrappigPanda/notorious/reaper"
	"github.com/GrappigPanda/notorious/server"
	"time"
)

func Init() {
	dbConn, err := db.OpenConnection()
	if err != nil {
		panic("Failed to open connection to remote database.")
	}
	db.InitDB(dbConn)

	go reaper.StartReapingScheduler(1 * time.Minute)
}

func main() {
	Init()
	server.RunServer()
}
