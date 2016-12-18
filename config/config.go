package config

import (
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
	IRCCfg    *IRCConfig
	UseRSS    bool
}

type IRCConfig struct {
	Nick   string
	Pass   string
	User   string
	Name   string
	Server string
	Port   int
	Chan   string
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
		viper.BindEnv("dbuser")
	}

	whitelist, err := strconv.ParseBool(viper.Get("whitelist").(string))
	if err != nil {
		whitelist = false
	}

	return loadSQLOptions(whitelist)
}

func loadSQLOptions(whitelist bool) ConfigStruct {
	var sqlDeployOption string
	if viper.GetBool("UseMySQL") {
		sqlDeployOption = "mysql"
	} else {
		sqlDeployOption = "postgres"
	}

	var ircCfg *IRCConfig
	if viper.GetBool("UseIRCNotify") {
		ircCfg = loadIRCOptions()
	} else {
		ircCfg = nil
	}

	var useRSS bool
	if viper.GetBool("UseRSSNotify") {
		useRSS = true
	} else {
		useRSS = false
	}

	if viper.Get("dbpass").(string) != "" {
		return ConfigStruct{
			sqlDeployOption,
			viper.Get("dbhost").(string),
			viper.Get("dbport").(string),
			viper.Get("dbuser").(string),
			viper.Get("dbpass").(string),
			viper.Get("dbname").(string),
			whitelist,
			ircCfg,
			useRSS,
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
			ircCfg,
			useRSS,
		}
	}
}

func loadIRCOptions() *IRCConfig {
	return &IRCConfig{
		Nick:   viper.Get("ircnick").(string),
		Pass:   viper.Get("ircpass").(string),
		User:   viper.Get("ircnick").(string),
		Name:   viper.Get("ircnick").(string),
		Server: viper.Get("ircserver").(string),
		Port:   viper.GetInt("ircport"),
		Chan:   viper.Get("ircchan").(string),
	}
}
