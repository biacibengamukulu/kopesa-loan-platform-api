package cassandra

import (
	"strings"

	"github.com/gocql/gocql"

	"github.com/biangacila/kopesa-loan-platform-api/internal/customer/domain"
)

type UserRepository struct{ session *gocql.Session }
type RoleRepository struct{ session *gocql.Session }
type BranchRepository struct{ session *gocql.Session }
type AreaRepository struct{ session *gocql.Session }

func NewUserRepository(session *gocql.Session) *UserRepository {
	return &UserRepository{session: session}
}
func NewRoleRepository(session *gocql.Session) *RoleRepository {
	return &RoleRepository{session: session}
}
func NewBranchRepository(session *gocql.Session) *BranchRepository {
	return &BranchRepository{session: session}
}
func NewAreaRepository(session *gocql.Session) *AreaRepository {
	return &AreaRepository{session: session}
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	iter := r.session.Query(`SELECT id, full_name, email, password_hash, role, allowed_roles, branch_id, area_id, avatar_color, active FROM iam_users_by_email WHERE email = ? LIMIT 1`, strings.ToLower(email)).Iter()
	defer iter.Close()

	var user domain.User
	var branchID, areaID *string
	if !iter.Scan(&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.Role, &user.AllowedRoles, &branchID, &areaID, &user.AvatarColor, &user.Active) {
		return nil, nil
	}
	user.BranchID = branchID
	user.AreaID = areaID
	return &user, iter.Close()
}

func (r *UserRepository) FindByID(id string) (*domain.User, error) {
	iter := r.session.Query(`SELECT id, full_name, email, password_hash, role, allowed_roles, branch_id, area_id, avatar_color, active FROM iam_users WHERE id = ? LIMIT 1`, id).Iter()
	defer iter.Close()
	var user domain.User
	var branchID, areaID *string
	if !iter.Scan(&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.Role, &user.AllowedRoles, &branchID, &areaID, &user.AvatarColor, &user.Active) {
		return nil, nil
	}
	user.BranchID = branchID
	user.AreaID = areaID
	return &user, iter.Close()
}

func (r *UserRepository) List() ([]domain.User, error) {
	iter := r.session.Query(`SELECT id, full_name, email, password_hash, role, allowed_roles, branch_id, area_id, avatar_color, active FROM iam_users`).Iter()
	defer iter.Close()
	out := make([]domain.User, 0)
	for {
		var user domain.User
		var branchID, areaID *string
		if !iter.Scan(&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.Role, &user.AllowedRoles, &branchID, &areaID, &user.AvatarColor, &user.Active) {
			break
		}
		user.BranchID = branchID
		user.AreaID = areaID
		out = append(out, user)
	}
	return out, iter.Close()
}

func (r *UserRepository) Create(user domain.User) error {
	if err := r.session.Query(`INSERT INTO iam_users (id, full_name, email, password_hash, role, allowed_roles, branch_id, area_id, avatar_color, active) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ID, user.FullName, strings.ToLower(user.Email), user.PasswordHash, user.Role, user.AllowedRoles, user.BranchID, user.AreaID, user.AvatarColor, user.Active,
	).Exec(); err != nil {
		return err
	}
	return r.session.Query(`INSERT INTO iam_users_by_email (email, id, full_name, password_hash, role, allowed_roles, branch_id, area_id, avatar_color, active) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		strings.ToLower(user.Email), user.ID, user.FullName, user.PasswordHash, user.Role, user.AllowedRoles, user.BranchID, user.AreaID, user.AvatarColor, user.Active,
	).Exec()
}

func (r *UserRepository) Update(user domain.User) error {
	return r.session.Query(`UPDATE iam_users SET full_name = ?, email = ?, password_hash = ?, role = ?, allowed_roles = ?, branch_id = ?, area_id = ?, avatar_color = ?, active = ? WHERE id = ?`,
		user.FullName, strings.ToLower(user.Email), user.PasswordHash, user.Role, user.AllowedRoles, user.BranchID, user.AreaID, user.AvatarColor, user.Active, user.ID,
	).Exec()
}

func (r *RoleRepository) List() ([]domain.Role, error) {
	iter := r.session.Query(`SELECT id, label, scope, permissions FROM iam_roles`).Iter()
	defer iter.Close()
	out := make([]domain.Role, 0)
	for {
		var role domain.Role
		if !iter.Scan(&role.ID, &role.Label, &role.Scope, &role.Permissions) {
			break
		}
		out = append(out, role)
	}
	return out, iter.Close()
}

func (r *RoleRepository) FindByID(id string) (*domain.Role, error) {
	iter := r.session.Query(`SELECT id, label, scope, permissions FROM iam_roles WHERE id = ? LIMIT 1`, id).Iter()
	defer iter.Close()
	var role domain.Role
	if !iter.Scan(&role.ID, &role.Label, &role.Scope, &role.Permissions) {
		return nil, nil
	}
	return &role, iter.Close()
}

func (r *BranchRepository) List() ([]domain.Branch, error) {
	iter := r.session.Query(`SELECT id, code, name, area_id FROM iam_branches`).Iter()
	defer iter.Close()
	out := make([]domain.Branch, 0)
	for {
		var branch domain.Branch
		if !iter.Scan(&branch.ID, &branch.Code, &branch.Name, &branch.AreaID) {
			break
		}
		out = append(out, branch)
	}
	return out, iter.Close()
}

func (r *AreaRepository) List() ([]domain.Area, error) {
	iter := r.session.Query(`SELECT id, name FROM iam_areas`).Iter()
	defer iter.Close()
	out := make([]domain.Area, 0)
	for {
		var area domain.Area
		if !iter.Scan(&area.ID, &area.Name) {
			break
		}
		out = append(out, area)
	}
	return out, iter.Close()
}
