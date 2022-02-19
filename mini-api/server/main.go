package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	conf "github.com/k-digital-project/mini-api/config"
	"github.com/k-digital-project/mini-api/middlewares"
	"github.com/k-digital-project/mini-api/server/bootstrap"
	"github.com/k-digital-project/mini-api/usecase"
)

var (
	logFormat = `{"host":"${host}","pid":"${pid}","time":"${time}","req_id":"${locals:requestid}","status":"${status}","method":"${method}","latency":"${latency}","path":"${path}",` +
		`"user_agent":"${ua}","in":"${bytesReceived}", "req_body":"", "out":"${bytesSent}","res_body":"${resBody}"}`
)

func main() {
	configs, err := conf.LoadConfigs()
	if err != nil {
		log.Fatal(err.Error())
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: middlewares.InternalServer,
	})

	ContractUC := usecase.ContractUC{
		EnvConfig: configs.EnvConfig,
		DB:        configs.DB,
	}

	boot := bootstrap.Bootstrap{
		App:        app,
		ContractUC: ContractUC,
	}
	boot.App.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
	}))
	boot.App.Use(recover.New())
	boot.App.Use(requestid.New())
	boot.App.Use(cors.New(cors.Config{
		AllowOrigins:     configs.EnvConfig["APP_CORS_DOMAIN"],
		AllowMethods:     http.MethodHead + "," + http.MethodGet + "," + http.MethodPost + "," + http.MethodPut + "," + http.MethodPatch + "," + http.MethodDelete,
		AllowHeaders:     "*",
		AllowCredentials: false,
	}))
	boot.App.Use(logger.New(logger.Config{
		Format:     logFormat + "\n",
		TimeFormat: time.RFC3339,
		TimeZone:   "Asia/Jakarta",
	}))
	boot.RegisterRouters()
	log.Fatal(boot.App.Listen(configs.EnvConfig["APP_HOST"]))
	defer configs.DB.Close()
}
