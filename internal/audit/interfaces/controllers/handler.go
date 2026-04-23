package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/biangacila/kopesa-loan-platform-api/internal/audit/application"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/httpx"
	"github.com/biangacila/kopesa-loan-platform-api/internal/shared"
)

type Handler struct {
	service *application.Service
}

func RegisterRoutes(router fiber.Router, service *application.Service) {
	h := &Handler{service: service}
	router.Get("/audit/events", h.List)
}

func (h *Handler) List(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "100"))
	events, err := h.service.List(limit)
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, events, shared.Meta{})
}
