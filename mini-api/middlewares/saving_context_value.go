package middlewares

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
)

// SavingContextValue Save all values saved in Fiber Locals to Context
// So we can access this values from anywhere in the application flow
// For every incoming request will always have different values
// It is important to register this middleware in the end of middleware stack
// So we can capture all the fiber's Locals value
func SavingContextValue(timeout time.Duration) fiber.Handler {
	return func(f *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		ctx = context.WithValue(ctx, "requestid", f.Locals("requestid"))
		ctx = context.WithValue(ctx, "user_id", cast.ToString(f.Locals("user_id")))

		f.Locals("ctx", ctx)
		if err := f.Next(); err != nil {
			return err
		}

		return nil
	}
}
