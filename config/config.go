package config

import (
	"github.com/spf13/viper"
	"strconv"
)

// ConfigStruct holds the values that our config file holds
type ConfigStruct struct {
	MySQLHost string
	MySQLPort string
	MySQLUser string
	MySQLPass string
	MySQLDB   string
	Whitelist bool
}

// LoadConfig loads the config into the Config Struct and returns the
// ConfigStruct object. Will load from environmental variables (all caps) if we
// set a flag to true.
func LoadConfig() ConfigStruct {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("/etc/")
	viper.AddConfigPath("$GOPATH/src/github.com/GrappigPanda/notorious/")

	err := viper.ReadInConfig()
	if err != nil {
		panic("Failed to open config file")
	}

	if viper.GetBool("UseEnvVariables") == true {
		viper.AutomaticEnv()
		viper.BindEnv("mysqluser")
	}

	whitelist, err := strconv.ParseBool(viper.Get("whitelist").(string))
	if err != nil {
		whitelist = false
	}

	if viper.Get("MySQLPass").(string) != "" {
		return ConfigStruct{
			viper.Get("mysqlhost").(string),
			viper.Get("mysqlport").(string),
			viper.Get("mysqluser").(string),
			viper.Get("mysqlpass").(string),
			viper.Get("mysqldb").(string),
			whitelist,
		}
	} else {
		return ConfigStruct{
			viper.Get("mysqlhost").(string),
			viper.Get("mysqlport").(string),
			viper.Get("mysqluser").(string),
			"",
			viper.Get("mysqldb").(string),
			whitelist,
		}
	}

}
