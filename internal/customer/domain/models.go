package domain

type User struct {
	ID           string   `json:"id"`
	FullName     string   `json:"fullName"`
	Email        string   `json:"email"`
	PasswordHash string   `json:"-"`
	Role         string   `json:"role"`
	AllowedRoles []string `json:"allowedRoles"`
	BranchID     *string  `json:"branchId,omitempty"`
	AreaID       *string  `json:"areaId,omitempty"`
	AvatarColor  string   `json:"avatarColor"`
	Active       bool     `json:"active"`
}

type Role struct {
	ID          string   `json:"id"`
	Label       string   `json:"label"`
	Scope       string   `json:"scope"`
	Permissions []string `json:"permissions"`
}

type Branch struct {
	ID     string `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	AreaID string `json:"areaId"`
}

type Area struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserRepository interface {
	FindByEmail(email string) (*User, error)
	FindByID(id string) (*User, error)
	List() ([]User, error)
	Create(user User) error
	Update(user User) error
}

type RoleRepository interface {
	List() ([]Role, error)
	FindByID(id string) (*Role, error)
}

type BranchRepository interface {
	List() ([]Branch, error)
}

type AreaRepository interface {
	List() ([]Area, error)
}
