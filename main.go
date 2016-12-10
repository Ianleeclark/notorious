package main

import (
	"github.com/GrappigPanda/notorious/database/mysql"
	"github.com/GrappigPanda/notorious/reaper"
	"github.com/GrappigPanda/notorious/server"
	"time"
)

// Init handles initialziation of the server.
func Init() {
	dbConn, err := mysql.OpenConnection()
	if err != nil {
		panic("Failed to open connection to remote database.")
	}
	mysql.InitDB(dbConn)

	go reaper.StartReapingScheduler(1 * time.Minute)
}

func main() {
	server.RunServer()
}
