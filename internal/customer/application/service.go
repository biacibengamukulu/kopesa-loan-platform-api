package application

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/biangacila/kopesa-loan-platform-api/internal/customer/domain"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/auth"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/httpx"
)

type Service struct {
	users    domain.UserRepository
	roles    domain.RoleRepository
	branches domain.BranchRepository
	areas    domain.AreaRepository
	auth     *auth.Manager
}

func NewService(
	users domain.UserRepository,
	roles domain.RoleRepository,
	branches domain.BranchRepository,
	areas domain.AreaRepository,
	authManager *auth.Manager,
) *Service {
	return &Service{
		users:    users,
		roles:    roles,
		branches: branches,
		areas:    areas,
		auth:     authManager,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string      `json:"accessToken"`
	RefreshToken string      `json:"refreshToken"`
	ExpiresIn    int         `json:"expiresIn"`
	User         domain.User `json:"user"`
}

type RegisterRequest struct {
	FullName     string   `json:"fullName"`
	Email        string   `json:"email"`
	Password     string   `json:"password"`
	Role         string   `json:"role"`
	AllowedRoles []string `json:"allowedRoles"`
	BranchID     *string  `json:"branchId"`
	AreaID       *string  `json:"areaId"`
	AvatarColor  string   `json:"avatarColor"`
}

func (s *Service) Login(req LoginRequest) (*LoginResponse, error) {
	user, err := s.users.FindByEmail(req.Email)
	if err != nil || user == nil {
		return nil, httpx.NewError(fiber.StatusUnauthorized, "AUTH_INVALID_CREDENTIALS", "invalid credentials")
	}
	if !user.Active {
		return nil, httpx.NewError(fiber.StatusForbidden, "AUTH_USER_INACTIVE", "user is inactive")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return nil, httpx.NewError(fiber.StatusUnauthorized, "AUTH_INVALID_CREDENTIALS", "invalid credentials")
	}

	accessToken, err := s.auth.Issue(user.ID, user.Email, user.Role, user.AllowedRoles, 15*time.Minute)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.auth.Issue(user.ID, user.Email, user.Role, user.AllowedRoles, 30*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900,
		User:         *user,
	}, nil
}

func (s *Service) Register(req RegisterRequest) (*domain.User, error) {
	if existing, err := s.users.FindByEmail(req.Email); err == nil && existing != nil {
		return nil, httpx.NewError(fiber.StatusConflict, "AUTH_EMAIL_ALREADY_EXISTS", "email already exists")
	}
	role := strings.TrimSpace(req.Role)
	if role == "" {
		role = "branch_agent"
	}
	roleDef, err := s.roles.FindByID(role)
	if err != nil {
		return nil, err
	}
	if roleDef == nil {
		return nil, httpx.NewError(fiber.StatusUnprocessableEntity, "IAM_ROLE_NOT_FOUND", "role not found")
	}
	if len(req.AllowedRoles) == 0 {
		req.AllowedRoles = []string{role}
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, err
	}
	user := domain.User{
		ID:           "u-" + strings.ReplaceAll(strings.ToLower(strings.TrimSpace(req.Email)), "@", "-"),
		FullName:     req.FullName,
		Email:        strings.ToLower(strings.TrimSpace(req.Email)),
		PasswordHash: string(passwordHash),
		Role:         role,
		AllowedRoles: req.AllowedRoles,
		BranchID:     req.BranchID,
		AreaID:       req.AreaID,
		AvatarColor:  req.AvatarColor,
		Active:       true,
	}
	if user.AvatarColor == "" {
		user.AvatarColor = "hsl(210 60% 45%)"
	}
	if err := s.users.Create(user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Service) Me(userID string) (*domain.User, error) {
	user, err := s.users.FindByID(userID)
	if err != nil || user == nil {
		return nil, httpx.NewError(fiber.StatusNotFound, "USER_NOT_FOUND", "user not found")
	}
	return user, nil
}

func (s *Service) ListUsers() ([]domain.User, error)      { return s.users.List() }
func (s *Service) ListRoles() ([]domain.Role, error)      { return s.roles.List() }
func (s *Service) ListBranches() ([]domain.Branch, error) { return s.branches.List() }
func (s *Service) ListAreas() ([]domain.Area, error)      { return s.areas.List() }
