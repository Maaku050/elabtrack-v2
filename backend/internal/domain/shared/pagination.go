package shared

// Page describes a paginated query request.
type Page struct {
	Page    int
	PerPage int
	Search  string
	Sort    string
	Order   string
}

// PageMeta describes the pagination metadata returned to clients.
type PageMeta struct {
	Page     int `json:"page"`
	PerPage  int `json:"per_page"`
	Total    int `json:"total"`
	LastPage int `json:"last_page"`
}

// Normalize clamps page/per_page to safe defaults.
func (p *Page) Normalize(maxPerPage int) {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PerPage < 1 {
		p.PerPage = 20
	}
	if maxPerPage > 0 && p.PerPage > maxPerPage {
		p.PerPage = maxPerPage
	}
	if p.Sort == "" {
		p.Sort = "created_at"
	}
	if p.Order != "asc" && p.Order != "desc" {
		p.Order = "desc"
	}
}

// Offset returns the SQL offset for this page.
func (p Page) Offset() int { return (p.Page - 1) * p.PerPage }

// LastPage computes the last page number given a total.
func LastPage(total, perPage int) int {
	if perPage <= 0 {
		return 1
	}
	last := total / perPage
	if total%perPage != 0 {
		last++
	}
	if last < 1 {
		last = 1
	}
	return last
}
