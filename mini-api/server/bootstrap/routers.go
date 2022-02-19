package bootstrap

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/k-digital-project/mini-api/server/handlers"
)

//RegisterRouters regiter all routers
func (boot Bootstrap) RegisterRouters() {
	_ = handlers.Handlers{
		FiberApp:   boot.App,
		ContractUC: &boot.ContractUC,
	}

	boot.App.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusOK).JSON("success")
	})

}
