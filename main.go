package main

import (
	"github.com/GrappigPanda/notorious/reaper"
	"github.com/GrappigPanda/notorious/server"
	"github.com/spf13/viper"
	"time"
)

func loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic("Failed to open config file")
	}
}

func main() {
	go reaper.StartReapingScheduler(5 * 60 * time.Second)
	server.RunServer()
}
