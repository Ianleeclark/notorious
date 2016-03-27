package main

import (
	"github.com/GrappigPanda/notorious/reaper"
	"github.com/GrappigPanda/notorious/server"
	"github.com/spf13/viper"
	"time"
)

type configStruct struct {
	MySQLHost string
	MySQLPort string
	MySQLUser string
	MySQLPass string
}

func loadConfig() configStruct {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic("Failed to open config file")
	}

	return configStruct{
		viper.Get("MySQLHost").(string),
		viper.Get("MySQLPort").(string),
		viper.Get("MySQLUser").(string),
		viper.Get("MySQLPass").(string),
	}
}

func main() {
	c := server.OpenClient()
	_, err := c.Ping().Result()
	if err != nil {
		panic("No Redis instance detected. If deploying without Docker, install redis-server")
	}

	go reaper.StartReapingScheduler(1 * time.Minute)
	server.RunServer()
}
