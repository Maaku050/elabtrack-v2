package constants

// API version prefix.
const APIVersionPrefix = "/api/v1"

// Context keys used to attach request-scoped values.
const (
	CtxRequestID = "request_id"
)

// HTTP header names.
const (
	HeaderRequestID     = "X-Request-ID"
	HeaderAuthorization = "Authorization"
)

// Common error codes used in API responses.
const (
	CodeValidation         = "VALIDATION_ERROR"
	CodeUnauthorized       = "UNAUTHORIZED"
	CodeForbidden          = "FORBIDDEN"
	CodeNotFound           = "NOT_FOUND"
	CodeConflict           = "CONFLICT"
	CodeInvalidCredentials = "INVALID_CREDENTIALS"
	CodeTokenInvalid       = "TOKEN_INVALID"
	CodeRateLimited        = "RATE_LIMITED"
	CodeInternal           = "INTERNAL_ERROR"
	CodeBadRequest         = "BAD_REQUEST"
)
