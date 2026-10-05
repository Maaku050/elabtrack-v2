package pagination

import (
	"strconv"

	"github.com/fullstacktemplate/backend/internal/domain/shared"
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
