package usecase

import "database/sql"

type ContractUC struct {
	EnvConfig map[string]string
	DB        *sql.DB
}
