package usecase

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	"github.com/k-digital-project/mini-api/pkg/aes"
	"github.com/k-digital-project/mini-api/pkg/aesfront"
	"github.com/k-digital-project/mini-api/pkg/jwe"
	"github.com/k-digital-project/mini-api/pkg/jwt"
	"github.com/k-digital-project/mini-api/pkg/loggerpkg"
	"github.com/k-digital-project/mini-api/pkg/redis"
)

type ContractUC struct {
	EnvConfig   map[string]string
	DB          *sql.DB
	Translator  ut.Translator
	Validate    *validator.Validate
	RedisClient redis.RedisClient
	JweCred     jwe.Credential
	JwtCred     jwt.Credential
	Aes         aes.Credential
	AesFront    aesfront.Credential
}

// StoreToRedisExp save data to redis with key and exp time
func (uc ContractUC) StoreToRedisExp(key string, val interface{}, duration string) error {
	ctx := "ContractUC.StoreToRedisExp"

	dur, err := time.ParseDuration(duration)
	if err != nil {
		loggerpkg.Log(loggerpkg.WarnLevel, err.Error(), ctx, "parse_duration")
		return err
	}

	b, err := json.Marshal(val)
	if err != nil {
		loggerpkg.Log(loggerpkg.WarnLevel, err.Error(), ctx, "json_marshal")
		return err
	}

	err = uc.RedisClient.Client.Set(key, string(b), dur).Err()
	if err != nil {
		loggerpkg.Log(loggerpkg.WarnLevel, err.Error(), ctx, "redis_set")
		return err
	}

	return err
}

// GetFromRedis get value from redis by key
func (uc ContractUC) GetFromRedis(key string, cb interface{}) error {
	ctx := "ContractUC.GetFromRedis"

	res, err := uc.RedisClient.Client.Get(key).Result()
	if err != nil {
		loggerpkg.Log(loggerpkg.WarnLevel, err.Error(), ctx, "redis_get")
		return err
	}

	if res == "" {
		loggerpkg.Log(loggerpkg.WarnLevel, "", ctx, "redis_empty")
		return errors.New("[Redis] Value of " + key + " is empty.")
	}

	err = json.Unmarshal([]byte(res), &cb)
	if err != nil {
		loggerpkg.Log(loggerpkg.WarnLevel, err.Error(), ctx, "json_unmarshal")
		return err
	}

	return err
}

// RemoveFromRedis remove a key from redis
func (uc ContractUC) RemoveFromRedis(key string) error {
	ctx := "ContractUC.RemoveFromRedis"

	err := uc.RedisClient.Client.Del(key).Err()
	if err != nil {
		loggerpkg.Log(loggerpkg.WarnLevel, err.Error(), ctx, "redis_delete")
		return err
	}

	return err
}
