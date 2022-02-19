package config

import (
	"database/sql"
	"log"

	"github.com/joho/godotenv"
)

//Configs ...
type Configs struct {
	EnvConfig map[string]string
	DB        *sql.DB
}

//LoadConfigs load all configurations
func LoadConfigs() (res Configs, err error) {
	res.EnvConfig, err = godotenv.Read("../.env")
	if err != nil {
		log.Fatal("Error loading ..env file")
	}

	return res, err
}
