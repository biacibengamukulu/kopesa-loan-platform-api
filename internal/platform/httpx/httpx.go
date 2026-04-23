package httpx

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/apperr"
	"github.com/biangacila/kopesa-loan-platform-api/internal/shared"
)

func NewError(status int, code, message string, details ...shared.FieldError) *apperr.Error {
	return apperr.New(status, code, message, details...)
}

func Success(c *fiber.Ctx, status int, data any, meta shared.Meta) error {
	if meta.RequestID == "" {
		meta.RequestID = RequestID(c)
	}
	c.Set("X-Request-Id", meta.RequestID)
	return c.Status(status).JSON(shared.ResponseEnvelope{
		Data:  data,
		Meta:  meta,
		Error: nil,
	})
}

func Fail(c *fiber.Ctx, err error) error {
	var appErr *apperr.Error
	if !errors.As(err, &appErr) {
		appErr = apperr.New(fiber.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
	}

	reqID := RequestID(c)
	c.Set("X-Request-Id", reqID)
	return c.Status(appErr.Status).JSON(shared.ResponseEnvelope{
		Data: nil,
		Meta: shared.Meta{RequestID: reqID},
		Error: &shared.ErrorBody{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		},
	})
}

func RequestID(c *fiber.Ctx) string {
	return c.Locals("requestId").(string)
}
