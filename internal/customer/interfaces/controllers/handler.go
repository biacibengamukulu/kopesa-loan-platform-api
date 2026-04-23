package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/biangacila/kopesa-loan-platform-api/internal/customer/application"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/auth"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/httpx"
	"github.com/biangacila/kopesa-loan-platform-api/internal/shared"
)

type Handler struct {
	service *application.Service
	auth    *auth.Manager
}

func RegisterRoutes(router fiber.Router, service *application.Service, authManager *auth.Manager) {
	h := &Handler{service: service, auth: authManager}
	router.Post("/auth/login", h.Login)
	router.Post("/auth/register", h.Register)

	protected := router.Group("", auth.Middleware(authManager, false))
	protected.Get("/auth/me", h.Me)
	protected.Get("/users", h.ListUsers)
	protected.Post("/users", h.CreateUser)
	protected.Get("/roles", h.ListRoles)
	protected.Get("/branches", h.ListBranches)
	protected.Get("/areas", h.ListAreas)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var req application.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.NewError(fiber.StatusBadRequest, "VALIDATION_FAILED", "invalid request body")
	}
	resp, err := h.service.Login(req)
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, resp, shared.Meta{})
}

func (h *Handler) Me(c *fiber.Ctx) error {
	claims := auth.ClaimsFrom(c)
	user, err := h.service.Me(claims.UserID)
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, user, shared.Meta{})
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var req application.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.NewError(fiber.StatusBadRequest, "VALIDATION_FAILED", "invalid request body")
	}
	user, err := h.service.Register(req)
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusCreated, user, shared.Meta{})
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	var req application.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.NewError(fiber.StatusBadRequest, "VALIDATION_FAILED", "invalid request body")
	}
	user, err := h.service.Register(req)
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusCreated, user, shared.Meta{})
}

func (h *Handler) ListUsers(c *fiber.Ctx) error {
	data, err := h.service.ListUsers()
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, data, shared.Meta{})
}

func (h *Handler) ListRoles(c *fiber.Ctx) error {
	data, err := h.service.ListRoles()
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, data, shared.Meta{})
}

func (h *Handler) ListBranches(c *fiber.Ctx) error {
	data, err := h.service.ListBranches()
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, data, shared.Meta{})
}

func (h *Handler) ListAreas(c *fiber.Ctx) error {
	data, err := h.service.ListAreas()
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, data, shared.Meta{})
}
