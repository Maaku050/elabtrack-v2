package handlers

import (
	"errors"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/accounts"
	"strings"
	"testing"
)

func TestRosterDiagnosticsAreSafeAndLocated(t *testing.T) {
	for _, tc := range []struct {
		kind d.RosterIssueKind
		want string
	}{
		{d.RosterFormula, "Row 2, Name: Formulas"}, {d.RosterHeaders, "Header does not match"}, {d.RosterColumns, "five template columns"}, {d.RosterWorksheets, "exactly one worksheet"}, {d.RosterRows, "500 Student"}, {d.RosterSize, "768 KiB"}, {d.RosterExternalContent, "External links"}, {d.RosterEmpty, "at least one Student"},
	} {
		got := rosterErrorMessage(&d.RosterIssue{Kind: tc.kind, Row: 2, Column: 2})
		if !strings.Contains(got, tc.want) || len(got) > 256 {
			t.Fatalf("diagnostic %d: %s", tc.kind, got)
		}
	}
	if strings.Contains(rosterErrorMessage(errors.New("private-secret-workbook-name")), "private-secret") {
		t.Fatal("raw parser error leaked")
	}
	if strings.Contains(rosterErrorMessage(&d.RosterIssue{Kind: d.RosterFormula, Row: 99999, Column: 99}), "99999") {
		t.Fatal("unsafe location accepted")
	}
}
