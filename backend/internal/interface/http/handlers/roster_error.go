package handlers

import (
	"errors"
	"fmt"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/accounts"
)

func rosterErrorMessage(err error) string {
	var issue *d.RosterIssue
	if !errors.As(err, &issue) {
		return "Invalid Student workbook. Download the template and use one worksheet, text Student IDs and at most 500 rows; formulas are not accepted."
	}
	message := "Invalid Student workbook. Download a fresh template."
	switch issue.Kind {
	case d.RosterSize:
		message = "Choose a nonempty .xlsx workbook no larger than 768 KiB."
	case d.RosterArchive:
		message = "The file is not a readable .xlsx workbook. Save it as .xlsx or download a fresh template."
	case d.RosterUnsafeArchive:
		message = "The workbook exceeds safe archive limits or contains invalid archive entries. Use the downloaded template."
	case d.RosterExternalContent:
		message = "External links, macros and embedded files are not accepted. Use a plain .xlsx Student template."
	case d.RosterFormula:
		message = "Formulas are not accepted. Replace the formula with a literal text value."
	case d.RosterXML:
		message = "The workbook contains unreadable worksheet data. Save a fresh .xlsx copy using the template."
	case d.RosterWorksheets:
		message = "Use exactly one worksheet in the Student workbook."
	case d.RosterHeaders:
		message = "Header does not match the template. Use studentId, name, email, course, contactNumber in that order."
	case d.RosterColumns:
		message = "Only the five template columns are accepted. Remove extra columns, including passwords and roles."
	case d.RosterRows:
		message = "The workbook extends beyond row 501. Use at most 500 Student rows below the header."
	case d.RosterCellLength:
		message = "This cell exceeds the 512-byte import limit. Shorten the value."
	case d.RosterEmpty:
		message = "Add at least one Student row below the template header."
	}
	if issue.Row >= 1 && issue.Row <= 501 {
		location := fmt.Sprintf("Row %d", issue.Row)
		if issue.Column >= 1 && issue.Column <= 5 {
			location += ", " + []string{"Student ID", "Name", "Email", "Course", "Contact number"}[issue.Column-1]
		}
		message = location + ": " + message
	}
	return message
}
