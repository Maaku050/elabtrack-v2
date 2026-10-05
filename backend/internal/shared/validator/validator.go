package validator

import (
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

// Validator wraps go-playground/validator with a singleton instance.
type Validator struct {
	v *validator.Validate
}

var (
	once     sync.Once
	shared   *Validator
	sharedMu sync.Mutex
)

// New constructs a Validator with sensible defaults.
func New() *Validator {
	return &Validator{v: validator.New(validator.WithRequiredStructEnabled())}
}

// Shared returns a process-wide singleton Validator.
func Shared() *Validator {
	sharedMu.Lock()
	defer sharedMu.Unlock()
	if shared == nil {
		once.Do(func() { shared = New() })
	}
	return shared
}

// Validate validates a struct and returns a FieldErrors map keyed by JSON tag
// when validation fails. Returns nil if validation passes.
func (v *Validator) Validate(s any) FieldErrors {
	err := v.v.Struct(s)
	if err == nil {
		return nil
	}
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return FieldErrors{"_": err.Error()}
	}
	out := FieldErrors{}
	for _, e := range errs {
		field := e.Field()
		// Use the JSON tag if present (go-playground/validator exposes it
		// via StructField; we fall back to the Go field name).
		tag := jsonTag(s, field)
		if tag == "" {
			tag = strings.ToLower(field)
		}
		out[tag] = msgForTag(e)
	}
	return out
}

// FieldErrors maps field name -> human-readable error message.
type FieldErrors map[string]string

// Error makes FieldErrors satisfy the error interface so it can be returned
// from handlers and detected by the centralized error handler.
func (f FieldErrors) Error() string {
	if len(f) == 0 {
		return "validation error"
	}
	parts := make([]string, 0, len(f))
	for k, v := range f {
		parts = append(parts, k+": "+v)
	}
	return "validation error: " + strings.Join(parts, "; ")
}

func msgForTag(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return "must be at least " + e.Param() + " characters"
	case "max":
		return "must be at most " + e.Param() + " characters"
	case "len":
		return "must be exactly " + e.Param() + " characters"
	case "oneof":
		return "must be one of: " + e.Param()
	case "uuid":
		return "must be a valid UUID"
	default:
		return "is invalid"
	}
}

// jsonTag attempts to resolve the JSON tag for a struct field via the
// validate engine's StructField. Falls back to "" when not resolvable.
func jsonTag(s any, field string) string {
	// Lightweight approach: go-playground/validator does not expose the JSON
	// tag directly without reflection on the type. We use reflect here only
	// to read the struct tag, which is cheap and only runs on validation
	// failure.
	return reflectJSONTag(s, field)
}
