package config

import (
	"database/sql"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"github.com/k-digital-project/mini-api/pkg/aes"
	"github.com/k-digital-project/mini-api/pkg/aesfront"
	"github.com/k-digital-project/mini-api/pkg/jwe"
	"github.com/k-digital-project/mini-api/pkg/jwt"
	postgrespkg "github.com/k-digital-project/mini-api/pkg/postgresql"
	redisPkg "github.com/k-digital-project/mini-api/pkg/redis"
	"github.com/k-digital-project/mini-api/pkg/str"
)

//Configs ...
type Configs struct {
	EnvConfig   map[string]string
	DB          *sql.DB
	RedisClient redisPkg.RedisClient
	JweCred     jwe.Credential
	JwtCred     jwt.Credential
	Aes         aes.Credential
	AesFront    aesfront.Credential
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

	// redis conn
	redisOption := &redis.Options{
		Addr:     res.EnvConfig["REDIS_HOST"],
		Password: res.EnvConfig["REDIS_PASSWORD"],
		DB:       0,
	}
	res.RedisClient = redisPkg.RedisClient{Client: redis.NewClient(redisOption)}

	// jwe
	res.JweCred = jwe.Credential{
		KeyLocation: res.EnvConfig["APP_PRIVATE_KEY_LOCATION"],
		Passphrase:  res.EnvConfig["APP_PRIVATE_KEY_PASSPHRASE"],
	}

	// jwt
	res.JwtCred = jwt.Credential{
		Secret:           res.EnvConfig["TOKEN_SECRET"],
		ExpSecret:        str.StringToInt(res.EnvConfig["TOKEN_EXP_SECRET"]),
		RefreshSecret:    res.EnvConfig["TOKEN_REFRESH_SECRET"],
		RefreshExpSecret: str.StringToInt(res.EnvConfig["TOKEN_EXP_REFRESH_SECRET"]),
	}

	// aes
	res.Aes = aes.Credential{
		Key: res.EnvConfig["AES_KEY"],
	}

	// aes front
	res.AesFront = aesfront.Credential{
		Key: res.EnvConfig["AES_FRONT_KEY"],
		Iv:  res.EnvConfig["AES_FRONT_IV"],
	}

	return res, err
}
