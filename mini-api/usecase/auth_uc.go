package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/k-digital-project/mini-api/db/repository/models"
	"github.com/k-digital-project/mini-api/helper"
	"github.com/k-digital-project/mini-api/pkg/functioncaller"
	"github.com/k-digital-project/mini-api/pkg/loggerpkg"
	"github.com/k-digital-project/mini-api/pkg/str"
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
	user, err := userUc.FindByEmail(c, models.UserParamater{Email: data.Email}, false)
	if err != nil {
		loggerpkg.Log(loggerpkg.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-find")
	}

	if user.ID == "" {
		userRequest := requests.UserRequest{
			Name:         data.Name,
			Email:        data.Email,
			Password:     data.Password,
			Status:       models.UserStatusActive,
			RegisterType: models.UserRegisterTypeEmail,
		}

		user, err = userUc.Add(c, &userRequest, true)
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

// UserLogin ...
func (uc AuthUC) UserLogin(c context.Context, input *requests.UserLoginRequest) (res models.User, err error) {
	input.Password, err = uc.AesFront.Decrypt(input.Password)
	if err != nil {
		loggerpkg.Log(loggerpkg.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "decrypt_password")
		return res, err
	}

	userUc := UserUC{ContractUC: uc.ContractUC}
	res, err = userUc.FindByEmail(c, models.UserParamater{Email: input.User}, true)
	if err != nil && err != sql.ErrNoRows {
		loggerpkg.Log(loggerpkg.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-find")
	}

	if res.Password != input.Password {
		loggerpkg.Log(loggerpkg.WarnLevel, helper.InvalidPassword, functioncaller.PrintFuncName(), "check_password")
		return res, errors.New(helper.InvalidPassword)
	}

	return res, err
}

func (uc AuthUC) AuthLogin(c context.Context, input *requests.UserLoginRequest) (res viewmodel.JwtVM, err error) {
	if !str.CheckEmail(input.User) {
		loggerpkg.Log(loggerpkg.WarnLevel, input.User, functioncaller.PrintFuncName(), "invalid_email")
		return res, errors.New(helper.InvalidEmail)
	}

	user := models.User{}

	user, err = uc.UserLogin(c, input)
	if err != nil {
		loggerpkg.Log(loggerpkg.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "user-login")
		return res, err
	}

	// Jwe the payload & Generate jwt token
	payload := map[string]interface{}{
		"user_id": user.ID,
	}
	jwtUc := JwtUC{ContractUC: uc.ContractUC}
	err = jwtUc.GenerateToken(c, payload, &res)
	if err != nil {
		loggerpkg.Log(loggerpkg.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "generate_token")
		return res, err
	}

	res.UserID = user.ID

	return res, err
}
