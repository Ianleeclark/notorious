package main

import (
	"github.com/GrappigPanda/notorious/reaper"
	"github.com/GrappigPanda/notorious/server"
	"github.com/GrappigPanda/notorious/database"
	"time"
)

func Init() {
	dbConn, err := db.OpenConnection()
	if err != nil {
		panic("Failed to open connection to remote database.")
	}
	db.InitDB(dbConn)

	c := server.OpenClient()
	_, err = c.Ping().Result()
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
        server.CreateNewTorrentKey(c, *infoHash)
    }
}

func main() {
	Init()
	go reaper.StartReapingScheduler(5 * time.Minute)
	server.RunServer()
}
