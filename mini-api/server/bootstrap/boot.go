package bootstrap

import (
	"github.com/gofiber/fiber/v2"
	"github.com/k-digital-project/mini-api/usecase"
)

type Bootstrap struct {
	App        *fiber.App
	ContractUC usecase.ContractUC
}
