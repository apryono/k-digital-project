package routers

import (
	"github.com/gofiber/fiber/v2"
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
	r.Post("/register/email", handler.RegisterByEmail)
}
