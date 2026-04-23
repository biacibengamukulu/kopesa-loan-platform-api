package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/biangacila/kopesa-loan-platform-api/internal/loan/application"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/auth"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/httpx"
	"github.com/biangacila/kopesa-loan-platform-api/internal/shared"
)

func RegisterRoutes(router fiber.Router, service *application.Service) {
	h := &handler{service: service}
	router.Get("/loans", h.listLoans)
	router.Get("/loans/:id", h.getLoan)
	router.Get("/loans/applications", h.listApplications)
	router.Get("/loans/applications/:id", h.getApplication)
	router.Post("/loans/applications", h.createApplication)
	router.Post("/loans/applications/:id/assess", h.assessApplication)
	router.Post("/loans/applications/:id/approve", h.approveApplication)
	router.Post("/loans/:id/disburse", h.disburseLoan)
}

type handler struct{ service *application.Service }

func (h *handler) listLoans(c *fiber.Ctx) error {
	data, err := h.service.ListLoans(c.Query("status"), c.Query("branchId"))
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, data, shared.Meta{})
}

func (h *handler) getLoan(c *fiber.Ctx) error {
	data, err := h.service.GetLoan(c.Params("id"))
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, data, shared.Meta{})
}

func (h *handler) listApplications(c *fiber.Ctx) error {
	data, err := h.service.ListApplications(c.Query("status"))
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, data, shared.Meta{})
}

func (h *handler) getApplication(c *fiber.Ctx) error {
	data, err := h.service.GetApplication(c.Params("id"))
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, data, shared.Meta{})
}

func (h *handler) createApplication(c *fiber.Ctx) error {
	var req application.CreateApplicationRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.NewError(fiber.StatusBadRequest, "VALIDATION_FAILED", "invalid request body")
	}
	if claims := auth.ClaimsFrom(c); claims != nil && req.CreatedBy == "" {
		req.CreatedBy = claims.UserID
	}
	data, err := h.service.CreateApplication(req)
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusCreated, data, shared.Meta{})
}

func (h *handler) assessApplication(c *fiber.Ctx) error {
	var req application.AssessApplicationRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.NewError(fiber.StatusBadRequest, "VALIDATION_FAILED", "invalid request body")
	}
	data, err := h.service.AssessApplication(c.Params("id"), req)
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, data, shared.Meta{})
}

func (h *handler) approveApplication(c *fiber.Ctx) error {
	var req application.ApproveApplicationRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.NewError(fiber.StatusBadRequest, "VALIDATION_FAILED", "invalid request body")
	}
	if claims := auth.ClaimsFrom(c); claims != nil && req.ApprovedBy == "" {
		req.ApprovedBy = claims.UserID
	}
	data, err := h.service.ApproveApplication(c.Params("id"), req)
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, data, shared.Meta{})
}

func (h *handler) disburseLoan(c *fiber.Ctx) error {
	var req application.DisburseLoanRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.NewError(fiber.StatusBadRequest, "VALIDATION_FAILED", "invalid request body")
	}
	if claims := auth.ClaimsFrom(c); claims != nil && req.DisbursedBy == "" {
		req.DisbursedBy = claims.UserID
	}
	data, err := h.service.DisburseLoan(c.Params("id"), req)
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, data, shared.Meta{})
}
