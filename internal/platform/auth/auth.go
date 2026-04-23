package auth

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/apperr"
)

type Claims struct {
	UserID       string   `json:"userId"`
	Email        string   `json:"email"`
	Role         string   `json:"role"`
	AllowedRoles []string `json:"allowedRoles"`
	jwt.RegisteredClaims
}

type Manager struct {
	secret []byte
}

func NewManager(secret string) *Manager {
	return &Manager{secret: []byte(secret)}
}

func (m *Manager) Issue(userID, email, role string, allowedRoles []string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID:       userID,
		Email:        email,
		Role:         role,
		AllowedRoles: allowedRoles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	})
	return token.SignedString(m.secret)
}

func (m *Manager) Parse(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fiber.ErrUnauthorized
	}
	return claims, nil
}

func Middleware(m *Manager, optional bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		raw := strings.TrimSpace(c.Get("Authorization"))
		if raw == "" {
			if optional {
				return c.Next()
			}
			return apperr.New(fiber.StatusUnauthorized, "AUTH_TOKEN_MISSING", "authorization token required")
		}
		tokenString := strings.TrimSpace(strings.TrimPrefix(raw, "Bearer"))
		claims, err := m.Parse(tokenString)
		if err != nil {
			return apperr.New(fiber.StatusUnauthorized, "AUTH_TOKEN_EXPIRED", "invalid or expired token")
		}
		c.Locals("claims", claims)
		return c.Next()
	}
}

func ClaimsFrom(c *fiber.Ctx) *Claims {
	claims, _ := c.Locals("claims").(*Claims)
	return claims
}
