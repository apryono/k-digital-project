package usecase

import (
	"context"
	"database/sql"
	"time"

	"github.com/k-digital-project/mini-api/db/repository"
	"github.com/k-digital-project/mini-api/db/repository/models"
	"github.com/k-digital-project/mini-api/pkg/functioncaller"
	"github.com/k-digital-project/mini-api/pkg/loggerpkg"
	"github.com/k-digital-project/mini-api/pkg/str"
	"github.com/k-digital-project/mini-api/usecase/requests"
)

type UserUC struct {
	*ContractUC
	*sql.Tx
}

func (uc UserUC) BuildBody(res *models.User, showPassword bool) {
}

func (uc UserUC) FindByEmail(c context.Context, parameter models.UserParamater, showPassword bool) (res models.User, err error) {

	repo := repository.NewUserRepository(uc.DB, uc.Tx)
	res, err = repo.FindByEmail(c, parameter)
	if err != nil {
		loggerpkg.Log(loggerpkg.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-find-email")
		return res, err
	}

	uc.BuildBody(&res, showPassword)

	return res, err
}

func (uc UserUC) Add(c context.Context, input *requests.UserRequest, isVerifiedEmail bool) (res models.User, err error) {
	_ = uc.checkDetail(c, input)

	res = models.User{
		Name:         input.Name,
		Email:        input.Email,
		Status:       input.Status,
		RegisterType: input.RegisterType,
		Password:     input.Password,
	}
	if isVerifiedEmail && str.Contains(models.UserSocialMediaWhitelist, input.RegisterType) {
		res.Status = models.UserStatusActive
		now := time.Now().UTC().Format(time.RFC3339)
		res.EmailValidAt = &now
	}

	repo := repository.NewUserRepository(uc.DB, uc.Tx)
	res.ID, err = repo.Add(c, &res)
	if err != nil {
		loggerpkg.Log(loggerpkg.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query-add")
		return res, err
	}

	return res, err
}

func (uc UserUC) checkDetail(c context.Context, input *requests.UserRequest) (err error) {
	if input.Password != "" {
		input.Password = uc.Aes.EncryptString(input.Password)
	}

	return err
}
