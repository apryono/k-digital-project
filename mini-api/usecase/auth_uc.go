package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/k-digital-project/mini-api/db/repository/models"
	"github.com/k-digital-project/mini-api/helper"
	"github.com/k-digital-project/mini-api/pkg/functioncaller"
	"github.com/k-digital-project/mini-api/pkg/loggerpkg"
	"github.com/k-digital-project/mini-api/usecase/requests"
	"github.com/k-digital-project/mini-api/usecase/viewmodel"
)

type AuthUC struct {
	*ContractUC
	Tx *sql.Tx
}

//RegisterByEmail register by email
func (uc AuthUC) RegisterByEmail(c context.Context, data *requests.RegisterByEmailRequest) (res viewmodel.JwtVM, err error) {

	userUc := UserUC{ContractUC: uc.ContractUC}
	user, _ := userUc.FindByEmail(c, models.UserParamater{Email: data.Email}, false)

	if user.ID == "" {
		userRequest := requests.UserRequest{
			Name:         data.Name,
			Email:        data.Email,
			Status:       models.UserStatusPending,
			RegisterType: models.UserRegisterTypeEmail,
		}

		user, err = userUc.Add(c, &userRequest, false)
		if err != nil {
			loggerpkg.Log(loggerpkg.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-add")
			return res, err
		}
	} else if user.EmailValidAt == nil || user.Status == models.UserStatusActive {
		loggerpkg.Log(loggerpkg.WarnLevel, helper.DuplicateEmail, functioncaller.PrintFuncName(), "duplicate_email")
		return res, errors.New(helper.DuplicateEmail)
	}

	res, err = uc.GenerateToken(c, user.ID)
	if err != nil {
		loggerpkg.Log(loggerpkg.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "generate_token")
		return res, err
	}

	return res, err
}

// GenerateToken ...
func (uc AuthUC) GenerateToken(c context.Context, id string) (res viewmodel.JwtVM, err error) {
	payload := map[string]interface{}{
		"user_id": id,
	}
	jwtUc := JwtUC{ContractUC: uc.ContractUC}
	err = jwtUc.GenerateToken(c, payload, &res)
	if err != nil {
		loggerpkg.Log(loggerpkg.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "generate_token")
		return res, err
	}

	return res, err
}
