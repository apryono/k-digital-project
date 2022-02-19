package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/k-digital-project/mini-api/pkg/functioncaller"
	"github.com/k-digital-project/mini-api/pkg/loggerpkg"
	"github.com/k-digital-project/mini-api/usecase/viewmodel"
)

func InternalServer(ctx *fiber.Ctx, err error) error {
	// Statuscode defaults to 500
	code := fiber.StatusInternalServerError

	// Retreive the custom statuscode if it's an fiber.*Error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	loggerpkg.Log(loggerpkg.ErrorLevel, err.Error(), functioncaller.PrintFuncName(), "internal-server")
	return ctx.Status(code).JSON([]interface{}{viewmodel.ResponseErrorVM{Messages: err.Error()}})
}
