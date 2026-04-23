package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/biangacila/kopesa-loan-platform-api/internal/arrears/application"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/auth"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/httpx"
	"github.com/biangacila/kopesa-loan-platform-api/internal/shared"
)

type handler struct{ service *application.Service }

func RegisterRoutes(router fiber.Router, service *application.Service) {
	h := &handler{service: service}
	router.Get("/arrears/cases", h.list)
	router.Get("/arrears/cases/:id", h.get)
	router.Post("/arrears/cases/:id/allocate", h.allocate)
	router.Post("/arrears/cases/:id/ptps", h.ptp)
	router.Post("/arrears/cases/:id/payments", h.payment)
}

func (h *handler) list(c *fiber.Ctx) error {
	data, err := h.service.ListCases()
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, data, shared.Meta{})
}

func (h *handler) get(c *fiber.Ctx) error {
	data, err := h.service.GetCase(c.Params("id"))
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, data, shared.Meta{})
}

func (h *handler) allocate(c *fiber.Ctx) error {
	var req application.AllocateRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.NewError(fiber.StatusBadRequest, "VALIDATION_FAILED", "invalid request body")
	}
	data, err := h.service.Allocate(c.Params("id"), req)
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, data, shared.Meta{})
}

func (h *handler) ptp(c *fiber.Ctx) error {
	var req application.PTPRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.NewError(fiber.StatusBadRequest, "VALIDATION_FAILED", "invalid request body")
	}
	if claims := auth.ClaimsFrom(c); claims != nil && req.CapturedBy == "" {
		req.CapturedBy = claims.UserID
	}
	data, err := h.service.CreatePTP(c.Params("id"), req)
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusCreated, data, shared.Meta{})
}

func (h *handler) payment(c *fiber.Ctx) error {
	var req application.PaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.NewError(fiber.StatusBadRequest, "VALIDATION_FAILED", "invalid request body")
	}
	if claims := auth.ClaimsFrom(c); claims != nil && req.CapturedBy == "" {
		req.CapturedBy = claims.UserID
	}
	if req.CapturedAt == "" {
		req.CapturedAt = time.Now().UTC().Format(time.RFC3339)
	}
	data, err := h.service.CapturePayment(c.Params("id"), req)
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusCreated, data, shared.Meta{})
}
