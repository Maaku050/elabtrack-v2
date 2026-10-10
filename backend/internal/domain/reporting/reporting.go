package reporting

import (
	"context"
	"errors"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/google/uuid"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

var ErrExportLimit = errors.New("report export exceeds5000 rows; narrow filters")

type Filter struct {
	DueToday                            bool
	Page, PerPage                       int
	Search, Status                      string
	From, To                            *time.Time
	EquipmentID, CategoryID, BorrowerID *uuid.UUID
}

func (f Filter) Valid() bool {
	return utf8.ValidString(f.Search) && utf8.ValidString(f.Status) && !strings.ContainsRune(f.Search, 0) && !strings.ContainsRune(f.Status, 0) && f.Page >= 1 && f.Page <= 100000 && f.PerPage >= 1 && f.PerPage <= 100 && utf8.RuneCountInString(f.Search) <= 200 && utf8.RuneCountInString(f.Status) <= 100 && (f.From == nil || f.To == nil || f.From.Before(*f.To))
}

type Definition struct {
	Key       string   `json:"key"`
	Label     string   `json:"label"`
	Columns   []string `json:"columns"`
	Equipment bool     `json:"equipment_filter"`
	Category  bool     `json:"category_filter"`
	Borrower  bool     `json:"borrower_filter"`
	AdminOnly bool     `json:"admin_only"`
}

var Definitions = []Definition{
	{"inventory", "Inventory", []string{"Equipment ID", "Name", "Category", "Status", "Available", "Reserved", "Checked out", "Damaged held", "Total tracked", "Created UTC"}, true, true, false, false},
	{"movements", "Stock movements", []string{"Movement ID", "Equipment", "Kind", "Available change", "Reserved change", "Checked-out change", "Damaged-held change", "Total change", "Available after", "Reserved after", "Checked out after", "Damaged held after", "Total after", "Actor ID", "Explanation", "Occurred UTC"}, true, true, true, false},
	{"requests", "Borrowing requests", loanColumns(), true, true, true, false}, {"issued", "Issued loans", loanColumns(), true, true, true, false}, {"active", "Active checkouts", loanColumns(), true, true, true, false}, {"overdue", "Overdue loans", loanColumns(), true, true, true, false},
	{"returns", "Return dispositions", []string{"Return ID", "Reference", "Equipment", "Good received", "Damaged received", "Lost confirmed", "Actor ID", "Explanation", "Occurred UTC"}, true, true, true, false},
	{"incidents", "Damage and loss", []string{"Return ID", "Reference", "Equipment", "Damaged received", "Lost confirmed", "Actor ID", "Explanation", "Occurred UTC"}, true, true, true, false},
	{"replacements", "Replacement obligations", []string{"Obligation ID", "Reference", "Equipment", "Incident", "Required", "Accepted", "Outstanding", "Recorded UTC"}, true, true, true, false},
	{"fines", "Fine assessment and resolution", []string{"Borrowing ID", "Reference", "Borrower", "Status", "Due UTC", "Completed UTC", "Assessed centavos", "Cleared centavos", "Outstanding centavos", "Assessment final", "Recent clearance methods (latest25)"}, true, true, true, false},
	{"fine-clearances", "Fine resolution history", []string{"Clearance ID", "Reference", "Borrower", "Assessed at clearance centavos", "Cleared centavos", "Method", "Resolution note", "Admin actor ID", "Occurred UTC"}, false, false, true, false},
	{"account-audit", "Account activity", []string{"Event ID", "Account ID", "Account name", "Actor ID", "Action", "Occurred UTC"}, false, false, false, true},
}

func loanColumns() []string {
	return []string{"Borrowing ID", "Reference", "Borrower", "Borrower ID", "Status", "Entry path", "Submitted UTC", "Issued UTC", "Due UTC", "Completed UTC"}
}
func Find(key string) (Definition, bool) {
	for _, d := range Definitions {
		if d.Key == key {
			return d, true
		}
	}
	return Definition{}, false
}
func (d Definition) Accepts(f Filter) bool {
	return (d.Equipment || f.EquipmentID == nil) && (d.Category || f.CategoryID == nil) && (d.Borrower || f.BorrowerID == nil)
}

type Page struct {
	Kind    string     `json:"kind"`
	Columns []string   `json:"columns"`
	Rows    [][]string `json:"rows"`
	Page    int        `json:"page"`
	PerPage int        `json:"per_page"`
	Total   int64      `json:"total"`
	AsOf    time.Time  `json:"as_of"`
}
type Activity struct {
	ID          uuid.UUID `json:"id"`
	BorrowingID uuid.UUID `json:"borrowing_id"`
	Reference   string    `json:"reference"`
	Kind        string    `json:"kind"`
	At          time.Time `json:"at"`
}
type Trend struct {
	Day     string `json:"day"`
	Issued  int64  `json:"issued"`
	Pending int64  `json:"pending"`
	Denied  int64  `json:"denied"`
}
type Dashboard struct {
	Trends  []Trend          `json:"trends"`
	Metrics map[string]int64 `json:"metrics"`
	Recent  []Activity       `json:"recent"`
	AsOf    time.Time        `json:"as_of"`
}
type Repository interface {
	Dashboard(context.Context, uuid.UUID, user.Role, int) (Dashboard, error)
	Report(context.Context, string, Filter, int) (Page, error)
}

// CSV textual cells must not be interpreted as spreadsheet formulas, even after leading control/space characters.
func SafeCell(s string) string {
	clean := strings.TrimLeftFunc(s, func(r rune) bool { return unicode.IsSpace(r) || r == '\uFEFF' })
	if clean != "" && strings.ContainsRune("=+-@", rune(clean[0])) {
		return "'" + s
	}
	return s
}
