package config

import (
	"github.com/NotoriousTracker/viper"
)

type ConfigStruct struct {
	MySQLHost string
	MySQLPort string
	MySQLUser string
	MySQLPass string
	MySQLDB   string
}

func LoadConfig() ConfigStruct {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic("Failed to open config file")
	}

	return ConfigStruct{
		viper.Get("MySQLHost").(string),
		viper.Get("MySQLPort").(string),
		viper.Get("MySQLUser").(string),
		viper.Get("MySQLPass").(string),
		viper.Get("MySQLDB").(string),
	}
}
