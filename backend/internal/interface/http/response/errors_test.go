package response

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	domainauth "github.com/Maaku050/elabtrack-v2/backend/internal/domain/auth"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	domainuser "github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/shared/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestFoundationErrorMapping(t *testing.T) {
	cases := []struct {
		err    error
		status int
		code   string
	}{
		{domainuser.ErrUserNotFound, 404, "NOT_FOUND"}, {shared.ErrNotFound, 404, "NOT_FOUND"},
		{domainuser.ErrEmailAlreadyExists, 409, "CONFLICT"}, {shared.ErrConflict, 409, "CONFLICT"},
		{domainuser.ErrInvalidCredentials, 401, "INVALID_CREDENTIALS"}, {shared.ErrUnauthorized, 401, "UNAUTHORIZED"},
		{domainauth.ErrTokenNotFound, 401, "TOKEN_INVALID"}, {domainauth.ErrTokenInvalid, 401, "TOKEN_INVALID"}, {domainauth.ErrTokenExpired, 401, "TOKEN_INVALID"}, {domainauth.ErrTokenRevoked, 401, "TOKEN_INVALID"},
		{shared.ErrForbidden, 403, "FORBIDDEN"}, {domainuser.ErrUserInactive, 403, "FORBIDDEN"},
		{domainuser.ErrInvalidEmail, 400, "VALIDATION_ERROR"}, {domainuser.ErrInvalidName, 400, "VALIDATION_ERROR"}, {domainuser.ErrPasswordTooShort, 400, "VALIDATION_ERROR"}, {shared.ErrValidation, 400, "VALIDATION_ERROR"}, {shared.ErrInvalidInput, 400, "VALIDATION_ERROR"},
		{validator.FieldErrors{"email": "is required"}, 400, "VALIDATION_ERROR"},
		{shared.ErrInternal, 500, "INTERNAL_ERROR"}, {errors.New("private-SQL-secret"), 500, "INTERNAL_ERROR"},
		{&pgconn.PgError{Message: "private-SQL-secret", Detail: "private-SQL-secret"}, 500, "INTERNAL_ERROR"},
		{fiber.NewError(400, "private-SQL-secret"), 400, "BAD_REQUEST"}, {fiber.ErrUnauthorized, 401, "UNAUTHORIZED"}, {fiber.ErrForbidden, 403, "FORBIDDEN"}, {fiber.ErrNotFound, 404, "NOT_FOUND"}, {fiber.ErrMethodNotAllowed, 405, "METHOD_NOT_ALLOWED"},
		{fiber.ErrRequestEntityTooLarge, 413, "PAYLOAD_TOO_LARGE"}, {fiber.ErrUnsupportedMediaType, 415, "UNSUPPORTED_MEDIA_TYPE"}, {fiber.ErrUnprocessableEntity, 400, "VALIDATION_ERROR"}, {fiber.ErrTooManyRequests, 429, "RATE_LIMITED"}, {fiber.ErrRequestHeaderFieldsTooLarge, 431, "BAD_REQUEST"}, {fiber.ErrServiceUnavailable, 503, "SERVICE_UNAVAILABLE"},
		{fiber.NewError(200, "private-SQL-secret"), 500, "INTERNAL_ERROR"}, {fiber.NewError(599, "private-SQL-secret"), 500, "INTERNAL_ERROR"},
		{errors.Join(shared.ErrInternal, shared.ErrForbidden), 500, "INTERNAL_ERROR"},
	}
	for i, tt := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			for _, err := range []error{tt.err, fmt.Errorf("private-SQL-secret: %w", tt.err)} {
				m := Map(err)
				if m.Status != tt.status || m.Code != tt.code || m.Message == "" || strings.Contains(m.Message, "private") {
					t.Fatalf("unsafe/inconsistent mapping: %+v", m)
				}
			}
		})
	}
	if Map(nil).Status != 500 {
		t.Fatal("nil error must fail safely")
	}
	for _, err := range []error{domainuser.ErrInvalidEmail, domainuser.ErrInvalidName, domainuser.ErrPasswordTooShort, validator.FieldErrors{"email": "is required"}} {
		if len(Map(fmt.Errorf("wrapped: %w", err)).Fields) == 0 {
			t.Fatal("structured validation fields lost")
		}
	}
}
