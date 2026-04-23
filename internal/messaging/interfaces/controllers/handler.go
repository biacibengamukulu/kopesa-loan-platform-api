package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/biangacila/kopesa-loan-platform-api/internal/messaging/application"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/auth"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/httpx"
	"github.com/biangacila/kopesa-loan-platform-api/internal/shared"
)

type handler struct{ service *application.Service }

func RegisterRoutes(router fiber.Router, service *application.Service) {
	h := &handler{service: service}
	router.Get("/messaging/templates", h.listTemplates)
	router.Get("/messaging/log", h.listLogs)
	router.Post("/messaging/send", h.send)
}

func (h *handler) listTemplates(c *fiber.Ctx) error {
	data, err := h.service.ListTemplates()
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, data, shared.Meta{})
}

func (h *handler) listLogs(c *fiber.Ctx) error {
	data, err := h.service.ListLogs(c.Query("context"), c.Query("entityId"))
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, data, shared.Meta{})
}

func (h *handler) send(c *fiber.Ctx) error {
	var req application.SendRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.NewError(fiber.StatusBadRequest, "VALIDATION_FAILED", "invalid request body")
	}
	if claims := auth.ClaimsFrom(c); claims != nil && req.SentBy == "" {
		req.SentBy = claims.UserID
	}
	data, err := h.service.Send(req)
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusAccepted, data, shared.Meta{})
}
