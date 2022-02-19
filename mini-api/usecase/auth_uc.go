package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/k-digital-project/mini-api/db/helper"
	"github.com/k-digital-project/mini-api/db/repository/models"
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

	//generate token

	return res, err
}
