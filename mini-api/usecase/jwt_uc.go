package usecase

import (
	"context"
	"errors"

	"github.com/k-digital-project/mini-api/helper"
	"github.com/k-digital-project/mini-api/pkg/loggerpkg"
	"github.com/k-digital-project/mini-api/usecase/viewmodel"
	"github.com/rs/xid"
)

type JwtUC struct {
	*ContractUC
}

//GenerateToken use to generate token
func (uc JwtUC) GenerateToken(c context.Context, payload map[string]interface{}, res *viewmodel.JwtVM) (err error) {
	ctx := "JwtUC.GenerateToken"

	deviceID := xid.New().String()
	payload["device_id"] = deviceID
	err = uc.StoreToRedisExp("userDeviceID"+payload["user_id"].(string), deviceID, uc.EnvConfig["TOKEN_EXP_SECRET"]+"m")
	if err != nil {
		loggerpkg.Log(loggerpkg.WarnLevel, err.Error(), ctx, "device_id", c.Value("requestid"))
		return errors.New(helper.InternalServer)
	}

	jwePayload, err := uc.ContractUC.JweCred.Generate(payload)
	if err != nil {
		loggerpkg.Log(loggerpkg.WarnLevel, err.Error(), ctx, "jwe", c.Value("requestid"))
		return errors.New(helper.JWT)
	}

	res.Token, res.ExpiredDate, err = uc.ContractUC.JwtCred.GetToken(jwePayload)
	if err != nil {
		loggerpkg.Log(loggerpkg.WarnLevel, err.Error(), ctx, "jwt", c.Value("requestid"))
		return errors.New(helper.JWT)
	}
	res.RefreshToken, res.RefreshExpiredDate, err = uc.ContractUC.JwtCred.GetRefreshToken(jwePayload)
	if err != nil {
		loggerpkg.Log(loggerpkg.WarnLevel, err.Error(), ctx, "refresh_jwt", c.Value("requestid"))
		return errors.New(helper.JWT)
	}

	return err
}
