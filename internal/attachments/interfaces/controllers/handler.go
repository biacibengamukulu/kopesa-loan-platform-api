package controllers

import (
	"mime/multipart"

	"github.com/gofiber/fiber/v2"

	"github.com/biangacila/kopesa-loan-platform-api/internal/attachments/application"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/auth"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/httpx"
	"github.com/biangacila/kopesa-loan-platform-api/internal/shared"
)

type handler struct{ service *application.Service }

func RegisterRoutes(router fiber.Router, service *application.Service) {
	h := &handler{service: service}
	router.Post("/attachments/presign", h.presign)
	router.Post("/attachments/upload", h.upload)
	router.Post("/attachments", h.finalize)
	router.Get("/attachments", h.list)
	router.Get("/attachments/:id", h.get)
}

func (h *handler) presign(c *fiber.Ctx) error {
	var req application.PresignRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.NewError(fiber.StatusBadRequest, "VALIDATION_FAILED", "invalid request body")
	}
	return httpx.Success(c, fiber.StatusOK, h.service.Presign(req), shared.Meta{})
}

func (h *handler) finalize(c *fiber.Ctx) error {
	var req application.FinalizeRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.NewError(fiber.StatusBadRequest, "VALIDATION_FAILED", "invalid request body")
	}
	if claims := auth.ClaimsFrom(c); claims != nil && req.CapturedBy == "" {
		req.CapturedBy = claims.UserID
	}
	data, err := h.service.Finalize(req)
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusCreated, data, shared.Meta{})
}

func (h *handler) upload(c *fiber.Ctx) error {
	formFile, err := c.FormFile("file")
	if err != nil {
		return httpx.NewError(fiber.StatusBadRequest, "VALIDATION_FAILED", "file is required")
	}
	file, err := formFile.Open()
	if err != nil {
		return httpx.NewError(fiber.StatusBadRequest, "VALIDATION_FAILED", "file cannot be opened")
	}
	defer file.Close()

	capturedBy := ""
	if claims := auth.ClaimsFrom(c); claims != nil {
		capturedBy = claims.UserID
	}
	data, err := h.service.Upload(application.UploadDirectRequest{
		ID:         c.Query("attachmentId"),
		Context:    c.Query("context"),
		EntityID:   c.Query("entityId"),
		FileName:   firstNonEmpty(c.Query("fileName"), formFile.Filename),
		MimeType:   formFile.Header.Get("Content-Type"),
		SizeBytes:  formFile.Size,
		CapturedBy: capturedBy,
		File:       file.(multipart.File),
	})
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusCreated, data, shared.Meta{})
}

func (h *handler) list(c *fiber.Ctx) error {
	data, err := h.service.List(c.Query("entityId"), c.Query("context"))
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, data, shared.Meta{})
}

func (h *handler) get(c *fiber.Ctx) error {
	data, err := h.service.Get(c.Params("id"))
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, data, shared.Meta{})
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
