package routers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/k-digital-project/mini-api/middlewares"
	"github.com/k-digital-project/mini-api/pkg/str"
	"github.com/k-digital-project/mini-api/server/handlers"
)

// UserRoutes ...
type UserRoutes struct {
	RouterGroup fiber.Router
	Handler     handlers.Handler
}

//RegisterRoute register auth routes
func (route UserRoutes) RegisterRoute() {
	handler := handlers.UserHandler{Handler: route.Handler}
	jwtMiddleware := middlewares.JwtMiddleware{ContractUC: handler.ContractUC}

	r := route.RouterGroup.Group("/api/user")
	r.Use(jwtMiddleware.VerifyUser)
	r.Use(middlewares.SavingContextValue(time.Duration(str.StringToInt(route.Handler.ContractUC.EnvConfig["APP_TIMEOUT"])) * time.Second))
	r.Put("/edit", handler.EditUser)
}
