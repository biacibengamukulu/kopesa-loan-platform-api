package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/httpx"
	"github.com/biangacila/kopesa-loan-platform-api/internal/reporting/application"
	"github.com/biangacila/kopesa-loan-platform-api/internal/shared"
)

type handler struct{ service *application.Service }

func RegisterRoutes(router fiber.Router, service *application.Service) {
	h := &handler{service: service}
	router.Get("/reports/exec/overview", h.execOverview)
}

func (h *handler) execOverview(c *fiber.Ctx) error {
	data, err := h.service.ExecOverview(c.Query("period"))
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, data, shared.Meta{})
}
