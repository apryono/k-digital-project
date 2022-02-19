package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/k-digital-project/mini-api/db/repository/models"
	"github.com/k-digital-project/mini-api/pkg/str"
)

type IUserRepository interface {
	Add(c context.Context, model *models.User) (string, error)
	FindByEmail(c context.Context, parameter models.UserParamater) (models.User, error)
}

type UserRepository struct {
	DB *sql.DB
	Tx *sql.Tx
}

// NewUserRepository ...
func NewUserRepository(DB *sql.DB, Tx *sql.Tx) IUserRepository {
	return &UserRepository{DB: DB, Tx: Tx}
}

func (repository UserRepository) scanRow(row *sql.Row) (res models.User, err error) {
	err = row.Scan(
		&res.ID, &res.Name, &res.Email, &res.Password,
		&res.Status, &res.RegisterType, &res.RegisterDetail,
		&res.EmailValidAt, &res.LastSeen, &res.CreatedAt, &res.UpdatedAt,
	)
	if err != nil {
		return res, err
	}
	return res, nil
}

//FindByEmail ...
func (repository UserRepository) FindByEmail(c context.Context, parameter models.UserParamater) (res models.User, err error) {
	var conditionString string
	if str.Contains(models.UserStatusWhitelist, parameter.Status) {
		conditionString += ` AND lower(def.status) = '` + strings.ToLower(parameter.Status) + `'`
	}
	statement := str.Spacing(models.UserSelectStatement, models.UserWhereStatement, ` AND lower(def.email) ilike $1`, conditionString, models.UserGroupByStatement)

	row := repository.DB.QueryRowContext(c, statement, strings.ToLower(parameter.Email))
	res, err = repository.scanRow(row)
	if err != nil {
		return res, err
	}

	return res, nil
}

//Add ...
func (repository UserRepository) Add(c context.Context, model *models.User) (res string, err error) {
	statement := `INSERT INTO "user" 
	(name, email, email_valid_at, password, status, register_type)
	VALUES ($1, $2, $3, $4, $5, $6) returning id`

	fmt.Println(statement)

	if repository.Tx != nil {
		err = repository.Tx.QueryRowContext(c, statement,
			model.Name, model.Email, model.EmailValidAt, model.Password, model.Status, model.RegisterType,
		).Scan(&res)
	} else {
		err = repository.DB.QueryRowContext(c, statement,
			model.Name, model.Email, model.EmailValidAt, model.Password, model.Status, model.RegisterType,
		).Scan(&res)
	}
	if err != nil {
		return res, err
	}

	return res, nil
}
