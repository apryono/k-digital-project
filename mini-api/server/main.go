package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	idTranslations "github.com/go-playground/validator/v10/translations/id"
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
	validatorDriver *validator.Validate
	uni             *ut.UniversalTranslator
	translator      ut.Translator
	logFormat       = `{"host":"${host}","pid":"${pid}","time":"${time}","req_id":"${locals:requestid}","status":"${status}","method":"${method}","latency":"${latency}","path":"${path}",` +
		`"user_agent":"${ua}","in":"${bytesReceived}", "req_body":"", "out":"${bytesSent}","res_body":"${resBody}"}`
)

func main() {
	configs, err := conf.LoadConfigs()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer configs.DB.Close()

	// init validation driver
	validatorInit(&configs)

	app := fiber.New(fiber.Config{
		ErrorHandler: middlewares.InternalServer,
	})

	ContractUC := usecase.ContractUC{
		EnvConfig:   configs.EnvConfig,
		DB:          configs.DB,
		Validate:    validatorDriver,
		Translator:  translator,
		RedisClient: configs.RedisClient,
		JweCred:     configs.JweCred,
		JwtCred:     configs.JwtCred,
	}

	boot := bootstrap.Bootstrap{
		App:        app,
		ContractUC: ContractUC,
		Validator:  validatorDriver,
		Translator: translator,
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

}

func validatorInit(configs *conf.Configs) {
	en := en.New()
	id := id.New()
	uni = ut.New(en, id)

	transEN, _ := uni.GetTranslator("en")
	transID, _ := uni.GetTranslator("id")

	validatorDriver = validator.New()

	enTranslations.RegisterDefaultTranslations(validatorDriver, transEN)
	idTranslations.RegisterDefaultTranslations(validatorDriver, transID)

	switch configs.EnvConfig["APP_LOCALE"] {
	case "en":
		translator = transEN
	case "id":
		translator = transID
	}
}
