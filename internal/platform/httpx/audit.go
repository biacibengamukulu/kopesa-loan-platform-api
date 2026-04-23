package httpx

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	auditapp "github.com/biangacila/kopesa-loan-platform-api/internal/audit/application"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/auth"
)

func AuditMiddleware(auditService *auditapp.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			return err
		}
		if auditService == nil {
			return nil
		}
		if c.Response().StatusCode() >= 400 {
			return nil
		}
		switch c.Method() {
		case fiber.MethodPost, fiber.MethodPatch, fiber.MethodPut, fiber.MethodDelete:
		default:
			return nil
		}

		claims := auth.ClaimsFrom(c)
		actor := ""
		role := ""
		if claims != nil {
			actor = claims.UserID
			role = claims.Role
		}

		entityID := c.Params("id")
		entityType := routeEntity(c.Route().Path)
		action := strings.ToLower(strings.TrimPrefix(c.Method(), ""))
		action = action + " " + c.Route().Path

		_ = auditService.Record(auditapp.RecordInput{
			Actor:      actor,
			ActorRole:  role,
			Action:     action,
			EntityType: entityType,
			EntityID:   entityID,
			IP:         c.IP(),
			UserAgent:  c.Get("User-Agent"),
			RequestID:  RequestID(c),
			Metadata: map[string]string{
				"method": c.Method(),
				"path":   c.Path(),
			},
		})
		return nil
	}
}

func routeEntity(path string) string {
	path = strings.Trim(path, "/")
	if path == "" {
		return "root"
	}
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return "unknown"
}
