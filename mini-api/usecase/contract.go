package usecase

import (
	"database/sql"

	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
)

type ContractUC struct {
	EnvConfig  map[string]string
	DB         *sql.DB
	Translator ut.Translator
	Validate   *validator.Validate
}
