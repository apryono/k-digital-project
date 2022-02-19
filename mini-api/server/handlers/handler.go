package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/k-digital-project/mini-api/usecase"
)

type Handlers struct {
	FiberApp   *fiber.App
	ContractUC *usecase.ContractUC
	DB         *sql.DB
}
