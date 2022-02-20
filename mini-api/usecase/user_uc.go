package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/k-digital-project/mini-api/db/repository"
	"github.com/k-digital-project/mini-api/db/repository/models"
	"github.com/k-digital-project/mini-api/helper"
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
	res.Password = str.ShowString(showPassword, uc.ContractUC.Aes.DecryptString(res.Password))
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
	if input.Email != "" {
		if !str.CheckEmail(input.Email) {
			loggerpkg.Log(loggerpkg.WarnLevel, helper.InvalidEmail, functioncaller.PrintFuncName(), "invalid_email")
			return errors.New(helper.InvalidEmail)
		}
	}

	return err
}

func (uc UserUC) Edit(c context.Context, id string, input *requests.UserRequest) (res models.User, err error) {
	err = uc.checkDetail(c, input)
	if err != nil {
		loggerpkg.Log(loggerpkg.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "check-detail", c.Value("requestid"))
		return res, err
	}

	user, _ := uc.FindByEmail(c, models.UserParamater{Email: input.Email}, false)
	if user.Email == input.Email {
		loggerpkg.Log(loggerpkg.WarnLevel, helper.DuplicateEmail, functioncaller.PrintFuncName(), "find_user", c.Value("requestid"))
		return res, errors.New(helper.DuplicateEmail)
	}

	res = models.User{
		ID:       id,
		Name:     input.Name,
		Email:    input.Email,
		Status:   input.Status,
		Password: input.Password,
	}

	repo := repository.NewUserRepository(uc.DB, uc.Tx)
	res.ID, err = repo.Edit(c, &res)
	if err != nil {
		loggerpkg.Log(loggerpkg.WarnLevel, helper.DuplicateEmail, functioncaller.PrintFuncName(), "edit-user", c.Value("requestid"))
		return res, err
	}

	return res, err
}

func (uc UserUC) UpdateLastSeen(c context.Context, userID string) (res string, err error) {
	now := time.Now().UTC()
	repo := repository.NewUserRepository(uc.ContractUC.DB, uc.Tx)
	_, err = repo.EditLastSeen(c, models.User{
		ID:       userID,
		LastSeen: now.Format(time.RFC3339),
	})

	return now.Format(time.RFC3339), err
}
