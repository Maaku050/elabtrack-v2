package response

import (
	"errors"
	"net/http"

	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/auth"
	domainshared "github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	domainuser "github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/shared/constants"
	"github.com/Maaku050/elabtrack-v2/backend/internal/shared/validator"
	"github.com/gofiber/fiber/v3"
)

// Error sends a structured error response. The HTTP status is derived from
// the error type when possible; otherwise 500. Production responses never
// leak SQL, stack traces, or internal paths.
func Error(c fiber.Ctx, err error) error {
	if err == nil {
		return Internal(c, "unexpected error")
	}

	// Validation field errors from the request validator.
	if fe, ok := err.(validator.FieldErrors); ok && fe != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(Body{
			Success: false,
			Message: "Validation failed.",
			Error: &ErrorBody{
				Code:   constants.CodeValidation,
				Fields: fe,
			},
		})
	}

	// Domain error -> HTTP status mapping.
	switch {
	case errors.Is(err, domainuser.ErrUserNotFound),
		errors.Is(err, auth.ErrTokenNotFound),
		errors.Is(err, domainshared.ErrNotFound):
		return Fail(c, http.StatusNotFound, "Resource not found.", constants.CodeNotFound)

	case errors.Is(err, domainuser.ErrEmailAlreadyExists),
		errors.Is(err, domainshared.ErrConflict):
		return Fail(c, http.StatusConflict, "Resource already exists.", constants.CodeConflict)

	case errors.Is(err, domainuser.ErrInvalidCredentials):
		return Fail(c, http.StatusUnauthorized, "Invalid email or password.", constants.CodeInvalidCredentials)

	case errors.Is(err, auth.ErrTokenInvalid),
		errors.Is(err, auth.ErrTokenExpired),
		errors.Is(err, auth.ErrTokenRevoked),
		errors.Is(err, domainshared.ErrUnauthorized):
		return Fail(c, http.StatusUnauthorized, "Authentication required.", constants.CodeUnauthorized)

	case errors.Is(err, domainshared.ErrForbidden):
		return Fail(c, http.StatusForbidden, "You do not have access to this resource.", constants.CodeForbidden)

	case errors.Is(err, domainuser.ErrUserInactive):
		return Fail(c, http.StatusForbidden, "Account is inactive.", constants.CodeForbidden)

	case errors.Is(err, domainuser.ErrInvalidEmail),
		errors.Is(err, domainuser.ErrInvalidName),
		errors.Is(err, domainuser.ErrPasswordTooShort),
		errors.Is(err, domainshared.ErrValidation),
		errors.Is(err, domainshared.ErrInvalidInput):
		return Fail(c, http.StatusBadRequest, err.Error(), constants.CodeValidation)
	}

	// Fallback: internal server error. The original error is logged upstream;
	// the client only sees a generic message.
	return Internal(c, "Something went wrong. Please try again later.")
}

// Fail sends a structured failure response with a specific status + code.
func Fail(c fiber.Ctx, status int, message, code string) error {
	return c.Status(status).JSON(Body{
		Success: false,
		Message: message,
		Error:   &ErrorBody{Code: code},
	})
}

// Internal sends a generic 500 response.
func Internal(c fiber.Ctx, message string) error {
	return c.Status(http.StatusInternalServerError).JSON(Body{
		Success: false,
		Message: message,
		Error:   &ErrorBody{Code: constants.CodeInternal},
	})
}

// BadRequest sends a 400 response with an optional code.
func BadRequest(c fiber.Ctx, message string) error {
	return Fail(c, http.StatusBadRequest, message, constants.CodeBadRequest)
}

// Unauthorized sends a 401 response.
func Unauthorized(c fiber.Ctx, message string) error {
	return Fail(c, http.StatusUnauthorized, message, constants.CodeUnauthorized)
}

// Forbidden sends a 403 response.
func Forbidden(c fiber.Ctx, message string) error {
	return Fail(c, http.StatusForbidden, message, constants.CodeForbidden)
}
