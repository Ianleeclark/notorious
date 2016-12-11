package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strconv"
)

// ConfigStruct holds the values that our config file holds
type ConfigStruct struct {
	DBChoice  string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	Whitelist bool
}

// LoadConfig loads the config into the Config Struct and returns the // ConfigStruct object. Will load from environmental variables (all caps) if we
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

	return loadMySQLOptions(whitelist)
}

func loadMySQLOptions(whitelist bool) ConfigStruct {
	var sqlDeployOption string
	if viper.GetBool("UseMySQL") {
		sqlDeployOption = "mysql"
	} else {
		sqlDeployOption = "postgres"
	}

	if viper.Get(fmt.Sprintf("dbpass", sqlDeployOption)).(string) != "" {
		return ConfigStruct{
			sqlDeployOption,
			viper.Get("dbhost").(string),
			viper.Get("dbport").(string),
			viper.Get("dbuser").(string),
			viper.Get("dbpass").(string),
			viper.Get("dbname").(string),
			whitelist,
		}
	} else {
		return ConfigStruct{
			sqlDeployOption,
			viper.Get("dbhost").(string),
			viper.Get("dbport").(string),
			viper.Get("dbuser").(string),
			"",
			viper.Get("dbname").(string),
			whitelist,
		}
	}
}
