package httpx

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		reqID := c.Get("X-Request-Id")
		if strings.TrimSpace(reqID) == "" {
			reqID = "req_" + uuid.NewString()
		}
		c.Locals("requestId", reqID)
		return c.Next()
	}
}

func ErrorMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.Next(); err != nil {
			return Fail(c, err)
		}
		return nil
	}
}
