package config

import (
	"github.com/spf13/viper"
)

// ConfigStruct holds the values that our config file holds
type ConfigStruct struct {
	MySQLHost string
	MySQLPort int
	MySQLUser string
	MySQLPass string
	MySQLDB   string
}

// LoadConfig loads the config into the Config Struct and returns the
// ConfigStruct object
func LoadConfig() ConfigStruct {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
    viper.AddConfigPath("../")

	err := viper.ReadInConfig()
	if err != nil {
		panic("Failed to open config file")
	}

    if(viper.GetBool("MySQLPass")) {
        return ConfigStruct{
            viper.Get("MySQLHost").(string),
            viper.Get("MySQLPort").(int),
            viper.Get("MySQLUser").(string),
            viper.Get("MySQLPass").(string),
            viper.Get("MySQLDB").(string),
        }
    } else {
        return ConfigStruct{
            viper.Get("MySQLHost").(string),
            viper.Get("MySQLPort").(int),
            viper.Get("MySQLUser").(string),
            "",
            viper.Get("MySQLDB").(string),
        }
    }

}
