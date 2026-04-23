package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/biangacila/kopesa-loan-platform-api/internal/campaign/application"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/auth"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/httpx"
	"github.com/biangacila/kopesa-loan-platform-api/internal/shared"
)

type handler struct{ service *application.Service }

func RegisterRoutes(router fiber.Router, service *application.Service) {
	h := &handler{service: service}
	router.Get("/campaigns", h.listCampaigns)
	router.Post("/campaigns", h.createCampaign)
	router.Get("/campaigns/:id/routes", h.listRoutes)
	router.Post("/campaigns/:id/routes", h.createRoute)
	router.Get("/leads", h.listLeads)
	router.Post("/campaigns/:id/leads", h.createLead)
	router.Post("/leads/:id/qualify", h.qualifyLead)
}

func (h *handler) listCampaigns(c *fiber.Ctx) error {
	data, err := h.service.ListCampaigns()
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, data, shared.Meta{})
}

func (h *handler) createCampaign(c *fiber.Ctx) error {
	var req application.CreateCampaignRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.NewError(fiber.StatusBadRequest, "VALIDATION_FAILED", "invalid request body")
	}
	data, err := h.service.CreateCampaign(req)
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusCreated, data, shared.Meta{})
}

func (h *handler) listRoutes(c *fiber.Ctx) error {
	data, err := h.service.ListRoutes(c.Params("id"))
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, data, shared.Meta{})
}

func (h *handler) createRoute(c *fiber.Ctx) error {
	var req application.CreateRouteRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.NewError(fiber.StatusBadRequest, "VALIDATION_FAILED", "invalid request body")
	}
	data, err := h.service.CreateRoute(c.Params("id"), req)
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusCreated, data, shared.Meta{})
}

func (h *handler) listLeads(c *fiber.Ctx) error {
	data, err := h.service.ListLeads()
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, data, shared.Meta{})
}

func (h *handler) createLead(c *fiber.Ctx) error {
	var req application.CreateLeadRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.NewError(fiber.StatusBadRequest, "VALIDATION_FAILED", "invalid request body")
	}
	if claims := auth.ClaimsFrom(c); claims != nil && req.CapturedBy == "" {
		req.CapturedBy = claims.UserID
	}
	data, err := h.service.CreateLead(c.Params("id"), req)
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusCreated, data, shared.Meta{})
}

func (h *handler) qualifyLead(c *fiber.Ctx) error {
	var req struct {
		Qualified string `json:"qualified"`
	}
	if err := c.BodyParser(&req); err != nil {
		return httpx.NewError(fiber.StatusBadRequest, "VALIDATION_FAILED", "invalid request body")
	}
	data, err := h.service.QualifyLead(c.Params("id"), req.Qualified)
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, data, shared.Meta{})
}
