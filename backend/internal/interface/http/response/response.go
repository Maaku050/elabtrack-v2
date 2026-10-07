package response

import (
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/shared/observability"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// Body retains the foundation success/message/data/meta envelope.
// Failures additionally carry code/message/requestId/optional fields in error.
// 204 is deliberately bodyless.
// Body is the single, consistent API response shape.
//
// Success:
//
//	{ "success": true, "message": "...", "data": {}, "meta": null }
//
// Error:
//
//	{ "success": false, "message": "...", "error": { "code": "...", "fields": {} } }
type Body struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    any              `json:"data"`
	Meta    *shared.PageMeta `json:"meta"`
	Error   *ErrorBody       `json:"error,omitempty"`
}

// ErrorBody is the error payload included when success is false.
type ErrorBody struct {
	Code      string            `json:"code"`
	Message   string            `json:"message"`
	RequestID string            `json:"requestId"`
	Fields    map[string]string `json:"fields,omitempty"`
}

// OK sends a 200 success response with data.
func OK(c fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusOK).JSON(Body{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a 201 success response with data.
func Created(c fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusCreated).JSON(Body{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Paginated sends a 200 success response with data and pagination meta.
func Paginated(c fiber.Ctx, message string, data any, meta shared.PageMeta) error {
	return c.Status(fiber.StatusOK).JSON(Body{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    &meta,
	})
}

// NoContent sends a 204 success response with no body.
func NoContent(c fiber.Ctx) error {
	return c.Status(fiber.StatusNoContent).SendString("")
}

// EnsureRequestID also covers parser/framework errors before normal middleware.
// Incoming values are deliberately ignored: correlation is server-owned.
func EnsureRequestID(c fiber.Ctx) string {
	id := observability.RequestID(c.Context())
	if id == "" {
		id = uuid.NewString()
		c.SetContext(observability.WithRequestID(c.Context(), id))
	}
	c.Set("X-Request-ID", id)
	return id
}
