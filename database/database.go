package db

import (
	"github.com/GrappigPanda/notorious/config"
	"github.com/GrappigPanda/notorious/database/mysql"
	"github.com/GrappigPanda/notorious/database/postgres"
	"github.com/jinzhu/gorm"
	// We use a blank import here because I'm afraid of breaking anything
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func InitDB(config *config.ConfigStruct) {
	if (*config).DBChoice == "mysql" {
		conn, err := mysql.OpenConnection()
		if err != nil {
			panic("Unable to open connection to remote server")
		}
		mysql.InitDB(conn)
	} else if (*config).DBChoice == "postgres" {
		conn, err := postgres.OpenConnection()
		if err != nil {
			panic("Unable to open connection to remote server")
		}
		postgres.InitDB(conn)
	} else {
		panic("Invalid Config choice for DBChoice. Set either `UsePostgres` or `UseMySQL`.")
	}
}

func OpenDBChoiceConnection() (*gorm.DB, error) {
	cfg := config.LoadConfig()

	if cfg.DBChoice == "mysql" {
		return mysql.OpenConnection()
	} else if cfg.DBChoice == "postgres" {
		return postgres.OpenConnection()
	} else {
		panic("Invalid database choice found for `OpenDBChoiceConnection`.")
	}
}
