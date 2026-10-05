package user

import (
	"context"
	"fmt"

	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	domainuser "github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/google/uuid"
)

// Service is the user application service. It coordinates domain logic
// through the user repository port. It contains no HTTP or persistence code.
type Service struct {
	repo domainuser.Repository
}

// NewService constructs a user application service.
func NewService(repo domainuser.Repository) *Service {
	return &Service{repo: repo}
}

// GetByID returns a single user by id.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (UserDTO, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return UserDTO{}, fmt.Errorf("find user by id: %w", err)
	}
	return toDTO(u), nil
}

// GetByEmail returns a single user by email.
func (s *Service) GetByEmail(ctx context.Context, email string) (UserDTO, error) {
	u, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return UserDTO{}, fmt.Errorf("find user by email: %w", err)
	}
	return toDTO(u), nil
}

// UpdateProfile changes only the authenticated account's display name.
// The repository must not rewrite role, status, password or email from a snapshot.
func (s *Service) UpdateProfile(ctx context.Context, cmd UpdateProfileCommand) (UserDTO, error) {
	name, err := domainuser.ParseName(cmd.Name)
	if err != nil {
		return UserDTO{}, err
	}
	account, err := s.repo.UpdateProfile(ctx, cmd.UserID, name.String())
	if err != nil {
		return UserDTO{}, fmt.Errorf("update profile: %w", err)
	}
	return FromAccount(*account), nil
}

// FromAccount maps a safe current-account snapshot to the existing profile DTO.
func FromAccount(a domainuser.Account) UserDTO {
	return UserDTO{ID: a.ID, Email: a.Email, Name: a.Name, Role: string(a.Role), IsActive: a.IsActive, CreatedAt: a.CreatedAt, UpdatedAt: a.UpdatedAt}
}

// List returns a paginated list of users.
func (s *Service) List(ctx context.Context, q ListUsersQuery) (ListResult, error) {
	q.Page.Normalize(100)
	items, total, err := s.repo.List(ctx, q.Page)
	if err != nil {
		return ListResult{}, fmt.Errorf("list users: %w", err)
	}
	dtos := make([]UserDTO, 0, len(items))
	for _, u := range items {
		dtos = append(dtos, toDTO(u))
	}
	return ListResult{
		Items: dtos,
		Meta: shared.PageMeta{
			Page:     q.Page.Page,
			PerPage:  q.Page.PerPage,
			Total:    total,
			LastPage: shared.LastPage(total, q.Page.PerPage),
		},
	}, nil
}

func toDTO(u *domainuser.User) UserDTO {
	return UserDTO{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Role:      string(u.Role),
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
