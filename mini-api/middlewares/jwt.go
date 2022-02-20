package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/k-digital-project/mini-api/helper"
	"github.com/k-digital-project/mini-api/pkg/functioncaller"
	"github.com/k-digital-project/mini-api/pkg/interfacepkg"
	"github.com/k-digital-project/mini-api/pkg/loggerpkg"
	"github.com/k-digital-project/mini-api/server/handlers"
	"github.com/k-digital-project/mini-api/usecase"
)

// JwtMiddleware ...
type JwtMiddleware struct {
	*usecase.ContractUC
}

// verify jwt middleware
func (jwtMiddleware JwtMiddleware) verify(ctx *fiber.Ctx, role string) (res map[string]interface{}, err error) {
	claims := &jwt.StandardClaims{}

	header := ctx.Get("Authorization")
	if !strings.Contains(header, "Bearer") {
		loggerpkg.Log(loggerpkg.WarnLevel, helper.HeaderNotPresent, functioncaller.PrintFuncName(), "middleware-jwt-header")
		return res, errors.New(helper.HeaderNotPresent)
	}

	//check claims and signing method
	token := strings.Replace(header, "Bearer ", "", -1)
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if jwt.SigningMethodHS256 != token.Method {
			loggerpkg.Log(loggerpkg.WarnLevel, helper.UnexpectedSigningMethod, functioncaller.PrintFuncName(), "middleware-jwt-checkSigningMethod")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		secret := jwtMiddleware.EnvConfig["TOKEN_SECRET"]
		return []byte(secret), nil
	})
	if err != nil {
		loggerpkg.Log(loggerpkg.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "middleware-jwt-checkClaims")
		return res, errors.New(helper.UnexpectedClaims)
	}

	//check token live time
	if claims.ExpiresAt < time.Now().Unix() {
		loggerpkg.Log(loggerpkg.WarnLevel, helper.ExpiredToken, functioncaller.PrintFuncName(), "middleware-jwt-checkTokenLiveTime")
		return res, errors.New(helper.ExpiredToken)
	}

	//jwe roll back encrypted id
	res, err = jwtMiddleware.JweCred.Rollback(claims.Id)
	if err != nil {
		loggerpkg.Log(loggerpkg.WarnLevel, helper.Unauthorized, functioncaller.PrintFuncName(), "pkg-jwe-rollback")
		return res, errors.New(helper.Unauthorized)
	}
	if res == nil {
		loggerpkg.Log(loggerpkg.WarnLevel, helper.Unauthorized, functioncaller.PrintFuncName(), "pkg-jwe-resultNil")
		return res, errors.New(helper.Unauthorized)
	}

	if role != "" && fmt.Sprintf("%v", res["role"]) != role {
		loggerpkg.Log(loggerpkg.WarnLevel, helper.InvalidRole, functioncaller.PrintFuncName(), "pkg-jwe-resultNil")
		return res, errors.New(helper.InvalidRole)
	}

	loggerpkg.Log(loggerpkg.InfoLevel, interfacepkg.Marshal(res), functioncaller.PrintFuncName(), "user", ctx.Locals("requestid"))

	return res, nil
}

// VerifyUser jwt middleware
func (jwtMiddleware JwtMiddleware) VerifyUser(ctx *fiber.Ctx) (err error) {
	c := context.Background()
	handler := handlers.Handler{ContractUC: jwtMiddleware.ContractUC}

	jweRes, err := jwtMiddleware.verify(ctx, "")
	if err != nil {
		loggerpkg.Log(loggerpkg.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "verify")
		return handler.SendResponse(ctx, nil, nil, err.Error(), http.StatusUnauthorized)
	}

	userID := jweRes["user_id"].(string)
	var lastSeen string
	err = jwtMiddleware.GetFromRedis("latestSeen"+userID, &lastSeen)
	if err != nil {
		loggerpkg.Log(loggerpkg.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "get_lastSeen_from_redis")

		userUc := usecase.UserUC{ContractUC: jwtMiddleware.ContractUC}
		now, err := userUc.UpdateLastSeen(c, userID)
		if err != nil {
			loggerpkg.Log(loggerpkg.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "update_last_seen")
			return handler.SendResponse(ctx, nil, nil, err.Error(), http.StatusUnauthorized)
		}

		err = jwtMiddleware.ContractUC.StoreToRedisExp("latestSeen"+userID, now, "1m")
		if err != nil {
			loggerpkg.Log(loggerpkg.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "store_to_redis")
			return errors.New(helper.InternalServer)
		}
	}

	// set id to use case contract
	ctx.Locals("user_id", userID)

	return ctx.Next()
}
