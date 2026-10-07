package pagination

import (
	"strconv"

	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/shared/validator"
	"github.com/gofiber/fiber/v3"
)

// MaxPerPage caps per_page to prevent unbounded queries.
const MaxPerPage = 100

// FromContext extracts pagination query params from a Fiber request and
// returns a normalized shared.Page.
func FromContext(c fiber.Ctx) shared.Page {
	q := shared.Page{
		Page:    atoi(c.Query("page"), 1),
		PerPage: atoi(c.Query("per_page"), 20),
		Search:  c.Query("search"),
		Sort:    c.Query("sort"),
		Order:   c.Query("order"),
	}
	q.Normalize(MaxPerPage)
	return q
}

func atoi(s string, fallback int) int {
	if s == "" {
		return fallback
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 {
		return fallback
	}
	return n
}

// ValidateQuery rejects supplied invalid values instead of silently hiding them.
func ValidateQuery(c fiber.Ctx) validator.FieldErrors {
	fields := validator.FieldErrors{}
	for _, key := range []string{"page", "per_page"} {
		if c.Request().URI().QueryArgs().Has(key) {
			n, err := strconv.Atoi(c.Query(key))
			if err != nil || n < 1 || n > 2147483647 {
				fields[key] = "must be a positive integer at most 2147483647"
			}
		}
	}
	// PostgreSQL integer OFFSET must remain representable even on 64-bit hosts.
	page := FromContext(c)
	if page.Page > 2147483647/page.PerPage {
		fields["page"] = "is too large for this page size"
	}
	if c.Request().URI().QueryArgs().Has("order") && c.Query("order") != "asc" && c.Query("order") != "desc" {
		fields["order"] = "must be asc or desc"
	}
	if c.Request().URI().QueryArgs().Has("sort") {
		switch c.Query("sort") {
		case "created_at", "updated_at", "email", "name":
		default:
			fields["sort"] = "is not a supported sort field"
		}
	}
	if len(c.Query("search")) > 256 {
		fields["search"] = "must be at most 256 bytes"
	}
	if len(fields) == 0 {
		return nil
	}
	return fields
}
