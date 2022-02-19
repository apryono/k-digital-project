package bootstrap

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"
	"github.com/k-digital-project/mini-api/usecase"
)

type Bootstrap struct {
	App        *fiber.App
	ContractUC usecase.ContractUC
	Validator  *validator.Validate
	Translator ut.Translator
}
