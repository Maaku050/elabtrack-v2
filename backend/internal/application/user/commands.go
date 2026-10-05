package user

import (
	"context"

	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/google/uuid"
)

// Commands are mutation use-case inputs.
type UpdateProfileCommand struct {
	UserID uuid.UUID
	Name   *string
	Email  *string
}

// CommandResult is the output of a mutation use case.
type CommandResult struct {
	User UserDTO
}

// Queries are read use-case inputs.
type GetCurrentUserQuery struct {
	UserID uuid.UUID
}

type ListUsersQuery struct {
	Page shared.Page
}

// QueryResult is the output of a read use case.
type QueryResult struct {
	User UserDTO
}

// ListResult is the output of a paginated list query.
type ListResult struct {
	Items []UserDTO
	Meta  shared.PageMeta
}

// ctxKey is an unexported type used to attach request-scoped values.
type ctxKey struct{}

// WithUserID stores the authenticated user id in the context.
func WithUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, ctxKey{}, id)
}

// UserIDFromContext retrieves the authenticated user id, if any.
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(ctxKey{}).(uuid.UUID)
	return id, ok
}
