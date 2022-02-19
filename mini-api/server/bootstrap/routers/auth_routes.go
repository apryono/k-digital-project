package routers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/k-digital-project/mini-api/middlewares"
	"github.com/k-digital-project/mini-api/pkg/str"
	"github.com/k-digital-project/mini-api/server/handlers"
)

//AuthRoutes ...
type AuthRoutes struct {
	RouterGroup fiber.Router
	Handler     handlers.Handler
}

//RegisterRoute register auth routes
func (route AuthRoutes) RegisterRoute() {
	handler := handlers.AuthHandler{Handler: route.Handler}

	r := route.RouterGroup.Group("/api/auth")
	r.Use(middlewares.SavingContextValue(time.Duration(str.StringToInt(route.Handler.ContractUC.EnvConfig["APP_TIMEOUT"])) * time.Second))
	r.Post("/register/email", handler.RegisterByEmail)
}
