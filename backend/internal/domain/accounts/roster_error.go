package accounts

import "github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"

// RosterIssue carries bounded coordinates and a trusted classification, never
// workbook text, filenames, XML contents or parser internals.
type RosterIssueKind uint8

const (
	RosterSize RosterIssueKind = iota
	RosterArchive
	RosterUnsafeArchive
	RosterExternalContent
	RosterFormula
	RosterXML
	RosterWorksheets
	RosterHeaders
	RosterColumns
	RosterRows
	RosterCellLength
	RosterEmpty
)

type RosterIssue struct {
	Kind        RosterIssueKind
	Row, Column int
}

func (e *RosterIssue) Error() string { return "invalid Student workbook" }
func (e *RosterIssue) Unwrap() error { return shared.ErrInvalidInput }
