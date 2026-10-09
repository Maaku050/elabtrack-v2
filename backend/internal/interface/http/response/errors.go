package response

import (
	"errors"

	domainaccounts "github.com/Maaku050/elabtrack-v2/backend/internal/domain/accounts"
	domainauth "github.com/Maaku050/elabtrack-v2/backend/internal/domain/auth"
	domainborrowing "github.com/Maaku050/elabtrack-v2/backend/internal/domain/borrowing"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	domainterms "github.com/Maaku050/elabtrack-v2/backend/internal/domain/terms"
	domainuser "github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/shared/constants"
	"github.com/Maaku050/elabtrack-v2/backend/internal/shared/validator"
	"github.com/gofiber/fiber/v3"
)

// Mapping contains only trusted public text, never err.Error().
type Mapping struct {
	Status        int
	Code, Message string
	Fields        map[string]string
}

// Map is the single transport mapping for wrapped foundation errors.
func Map(err error) Mapping {
	result := func(status int, code, message string) Mapping {
		return Mapping{Status: status, Code: code, Message: message}
	}
	// Internal classification takes priority over any joined lower-level cause.
	if err == nil || errors.Is(err, shared.ErrInternal) {
		return result(500, constants.CodeInternal, "Something went wrong. Please try again later.")
	}
	var fe validator.FieldErrors
	if errors.As(err, &fe) && len(fe) > 0 {
		return Mapping{400, constants.CodeValidation, "Validation failed.", fe}
	}
	switch {
	case errors.Is(err, domainborrowing.ErrExpired):
		return result(409, "BORROWING_EXPIRED", "The request expired and its reservation was released.")
	case errors.Is(err, domainborrowing.ErrStock):
		return result(409, "EQUIPMENT_NOT_AVAILABLE", "Equipment availability changed. Review current stock.")
	case errors.Is(err, domainborrowing.ErrState):
		return result(409, "BORROWING_STATE_CONFLICT", "The borrowing state changed. Refresh its details.")
	case errors.Is(err, domainborrowing.ErrKey):
		return result(409, "IDEMPOTENCY_CONFLICT", "This command key belongs to a different request.")
	case errors.Is(err, domainborrowing.ErrEligibility):
		return result(403, "BORROWER_NOT_ELIGIBLE", "The borrower must be active and complete account activation.")
	case errors.Is(err, domainaccounts.ErrDomainsMissing):
		return result(503, "STUDENT_DOMAINS_NOT_CONFIGURED", "Approved Student email domains are not configured.")
	case errors.Is(err, domainaccounts.ErrStudentDomain):
		return Mapping{400, constants.CodeValidation, "Validation failed.", map[string]string{"email": "must use a configured institutional domain"}}
	case errors.Is(err, domainaccounts.ErrStudentIDExists):
		return result(409, "STUDENT_ID_EXISTS", "This Student ID is already assigned.")
	case errors.Is(err, domainaccounts.ErrActivationInvalid):
		return result(400, "ACTIVATION_INVALID", "Activation link is invalid or expired. Request another from an administrator.")
	case errors.Is(err, domainaccounts.ErrCooldown):
		return result(429, "ACTIVATION_COOLDOWN", "Wait before sending another activation link.")
	case errors.Is(err, domainterms.ErrNotPublished):
		return result(503, "TERMS_NOT_PUBLISHED", "FSMO has not published borrowing terms yet.")
	case errors.Is(err, domainterms.ErrVersionNotFound):
		return result(404, "TERMS_VERSION_NOT_FOUND", "Terms version was not found.")
	case errors.Is(err, domainterms.ErrVersionChanged):
		return result(409, "TERMS_VERSION_CHANGED", "Terms have changed. Review the current version before accepting.")
	case errors.Is(err, domainterms.ErrAcceptanceRequired):
		return result(409, "TERMS_ACCEPTANCE_REQUIRED", "Accept the current borrowing terms before initiating a new borrowing.")
	case errors.Is(err, domainterms.ErrVersionExists):
		return result(409, "TERMS_VERSION_EXISTS", "This terms version identifier has already been published.")
	case errors.Is(err, domainterms.ErrPublicationChanged):
		return result(409, "TERMS_PUBLICATION_CHANGED", "The current terms version changed. Review it before publishing.")
	case errors.Is(err, domainuser.ErrUserNotFound), errors.Is(err, shared.ErrNotFound):
		return result(404, constants.CodeNotFound, "Resource not found.")
	case errors.Is(err, domainuser.ErrEmailAlreadyExists), errors.Is(err, shared.ErrConflict):
		return result(409, constants.CodeConflict, "The request conflicts with current state.")
	case errors.Is(err, domainuser.ErrInvalidCredentials):
		return result(401, constants.CodeInvalidCredentials, "Invalid email or password.")
	case errors.Is(err, domainauth.ErrTokenNotFound), errors.Is(err, domainauth.ErrTokenInvalid), errors.Is(err, domainauth.ErrTokenExpired), errors.Is(err, domainauth.ErrTokenRevoked):
		return result(401, constants.CodeTokenInvalid, "Session is invalid or expired. Please sign in again.")
	case errors.Is(err, shared.ErrUnauthorized):
		return result(401, constants.CodeUnauthorized, "Authentication required.")
	case errors.Is(err, domainuser.ErrUserInactive), errors.Is(err, shared.ErrForbidden):
		return result(403, constants.CodeForbidden, "You do not have access to this resource.")
	case errors.Is(err, domainuser.ErrInvalidEmail):
		return Mapping{400, constants.CodeValidation, "Validation failed.", map[string]string{"email": "must be a valid email address"}}
	case errors.Is(err, domainuser.ErrInvalidName):
		return Mapping{400, constants.CodeValidation, "Validation failed.", map[string]string{"name": "is invalid"}}
	case errors.Is(err, domainuser.ErrPasswordTooShort):
		return Mapping{400, constants.CodeValidation, "Validation failed.", map[string]string{"password": "does not meet password requirements"}}
	case errors.Is(err, shared.ErrValidation), errors.Is(err, shared.ErrInvalidInput):
		return result(400, constants.CodeValidation, "Invalid request.")
	}
	var framework *fiber.Error
	if errors.As(err, &framework) && framework != nil {
		switch framework.Code {
		case 400:
			return result(400, constants.CodeBadRequest, "Invalid request.")
		case 401:
			return Map(shared.ErrUnauthorized)
		case 403:
			return Map(shared.ErrForbidden)
		case 404:
			return Map(shared.ErrNotFound)
		case 405:
			return result(405, constants.CodeMethodNotAllowed, "Method is not allowed.")
		case 409:
			return Map(shared.ErrConflict)
		case 413:
			return result(413, constants.CodePayloadTooLarge, "Request body is too large.")
		case 415:
			return result(415, constants.CodeUnsupportedMedia, "JSON content type required.")
		case 422:
			return Map(shared.ErrValidation) // one 400 policy for input validation
		case 429:
			return result(429, constants.CodeRateLimited, "Too many requests. Please slow down.")
		case 431:
			return result(431, constants.CodeBadRequest, "Request headers are too large.")
		case 503:
			return result(503, constants.CodeUnavailable, "Service is not ready. Please try again later.")
		}
	}
	return Map(shared.ErrInternal)
}

type outcomeKey struct{}

// Outcome is safe metadata for the completion logger. It never retains raw errors.
type Outcome struct {
	Code, ErrorClass, Operation string
	Unexpected                  bool
}

func OutcomeFromContext(c fiber.Ctx) Outcome {
	value, _ := c.Locals(outcomeKey{}).(Outcome)
	return value
}
func record(c fiber.Ctx, m Mapping, err error) {
	// Only the trusted unavailable-policy mapping is expected. Joined internal
	// failures still map to INTERNAL_ERROR and retain error-level visibility.
	out := Outcome{Code: m.Code, Unexpected: m.Status >= 500 && m.Code != "TERMS_NOT_PUBLISHED"}
	if err != nil {
		out.ErrorClass, out.Operation = shared.FailureDetails(err)
	}
	c.Locals(outcomeKey{}, out)
}
func send(c fiber.Ctx, m Mapping) error {
	id := EnsureRequestID(c)
	return c.Status(m.Status).JSON(Body{Success: false, Message: m.Message, Error: &ErrorBody{Code: m.Code, Message: m.Message, RequestID: id, Fields: m.Fields}})
}
func Error(c fiber.Ctx, err error) error { m := Map(err); record(c, m, err); return send(c, m) }

// LoginError keeps Phase 1E account-enumeration containment, including a
// failed credential lookup. Unexpected lookup failures are still ERROR events.
func LoginError(c fiber.Ctx, err error) error {
	m := Map(err)
	_, operation := shared.FailureDetails(err)
	if errors.Is(err, domainuser.ErrUserInactive) || operation == "auth.login_lookup" {
		public := Map(domainuser.ErrInvalidCredentials)
		record(c, public, err)
		out := OutcomeFromContext(c)
		out.Unexpected = m.Status >= 500
		c.Locals(outcomeKey{}, out)
		return send(c, public)
	}
	return Error(c, err)
}

// Fail accepts code-owned public text only; never pass request/error values.
func Fail(c fiber.Ctx, status int, message, code string) error {
	if status >= 500 && code != constants.CodeUnavailable {
		return Error(c, shared.ErrInternal)
	}
	m := Mapping{Status: status, Code: code, Message: message}
	record(c, m, nil)
	return send(c, m)
}

// Internal deliberately ignores caller text and emits the one safe 500 message.
func Internal(c fiber.Ctx, _ string) error { return Error(c, shared.ErrInternal) }
func BadRequest(c fiber.Ctx, message string) error {
	return Fail(c, 400, message, constants.CodeBadRequest)
}
func Unauthorized(c fiber.Ctx, message string) error {
	return Fail(c, 401, message, constants.CodeUnauthorized)
}
func Forbidden(c fiber.Ctx, message string) error {
	return Fail(c, 403, message, constants.CodeForbidden)
}

// Unavailable records a dependency failure while exposing only readiness state.
func Unavailable(c fiber.Ctx, err error) error {
	m := Map(fiber.ErrServiceUnavailable)
	record(c, m, err)
	return send(c, m)
}
