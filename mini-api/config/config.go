package config

import (
	"database/sql"
	"log"
	"time"

	"github.com/joho/godotenv"
	postgrespkg "github.com/k-digital-project/mini-api/pkg/postgresql"
	"github.com/k-digital-project/mini-api/pkg/str"
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

	dbConn := postgrespkg.Connection{
		Host:                    res.EnvConfig["DATABASE_HOST"],
		DbName:                  res.EnvConfig["DATABASE_DB"],
		User:                    res.EnvConfig["DATABASE_USER"],
		Password:                res.EnvConfig["DATABASE_PASSWORD"],
		Port:                    str.StringToInt(res.EnvConfig["DATABASE_PORT"]),
		SslMode:                 res.EnvConfig["DATABASE_SSL_MODE"],
		DBMaxConnection:         str.StringToInt(res.EnvConfig["DATABASE_MAX_CONNECTION"]),
		DBMAxIdleConnection:     str.StringToInt(res.EnvConfig["DATABASE_MAX_IDLE_CONNECTION"]),
		DBMaxLifeTimeConnection: str.StringToInt(res.EnvConfig["DATABASE_MAX_LIFETIME_CONNECTION"]),
	}
	res.DB, err = dbConn.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
	res.DB.SetMaxOpenConns(dbConn.DBMaxConnection)
	res.DB.SetMaxIdleConns(dbConn.DBMAxIdleConnection)
	res.DB.SetConnMaxLifetime(time.Duration(dbConn.DBMaxLifeTimeConnection) * time.Second)

	return res, err
}
