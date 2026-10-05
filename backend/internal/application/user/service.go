package user

import (
	"context"
	"errors"
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

// UpdateProfile applies a partial update to the current user's profile.
func (s *Service) UpdateProfile(ctx context.Context, cmd UpdateProfileCommand) (UserDTO, error) {
	u, err := s.repo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return UserDTO{}, fmt.Errorf("find user: %w", err)
	}

	if cmd.Name != nil {
		name, err := domainuser.ParseName(*cmd.Name)
		if err != nil {
			return UserDTO{}, err
		}
		u.Name = name.String()
	}
	if cmd.Email != nil {
		email, err := domainuser.ParseEmail(*cmd.Email)
		if err != nil {
			return UserDTO{}, err
		}
		// Check uniqueness only when the email actually changes.
		if email.String() != u.Email {
			existing, err := s.repo.FindByEmail(ctx, email.String())
			if err == nil && existing.ID != u.ID {
				return UserDTO{}, domainuser.ErrEmailAlreadyExists
			}
			if err != nil && !errors.Is(err, domainuser.ErrUserNotFound) && !errors.Is(err, shared.ErrNotFound) {
				return UserDTO{}, fmt.Errorf("check email uniqueness: %w", err)
			}
			u.Email = email.String()
		}
	}

	if err := s.repo.Update(ctx, u); err != nil {
		return UserDTO{}, fmt.Errorf("update user: %w", err)
	}
	return toDTO(u), nil
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
