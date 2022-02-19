package bootstrap

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/k-digital-project/mini-api/server/bootstrap/routers"
	"github.com/k-digital-project/mini-api/server/handlers"
)

//RegisterRouters regiter all routers
func (boot Bootstrap) RegisterRouters() {
	handler := handlers.Handler{
		FiberApp:   boot.App,
		ContractUC: &boot.ContractUC,
		Validator:  boot.Validator,
		Translator: boot.Translator,
	}

	boot.App.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusOK).JSON("success")
	})

	apiV1 := boot.App.Group("/v1")

	// Auth Routes
	authRoutes := routers.AuthRoutes{RouterGroup: apiV1, Handler: handler}
	authRoutes.RegisterRoute()

}
