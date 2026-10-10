package reporting

import (
	"testing"
	"time"
)

func TestFormulaSafety(t *testing.T) {
	for _, s := range []string{"=1+1", "+SUM(A1)", "-2+3", "@cmd", " \t=HYPERLINK(\"x\")", "\r\n+formula", "\uFEFF=1", "\u00a0=1"} {
		if got := SafeCell(s); got != "'"+s {
			t.Fatalf("unsafe prefix %q", s)
		}
	}
	for _, s := range []string{"Spoon", "28366", "20", "PHP 10.00", "a,b\nc"} {
		if SafeCell(s) != s {
			t.Fatal("modified harmless cell")
		}
	}
}
func TestFilterBounds(t *testing.T) {
	valid := Filter{Page: 1, PerPage: 25}
	if !valid.Valid() {
		t.Fatal("default")
	}
	for _, bad := range []Filter{{Page: 0, PerPage: 25}, {Page: 100001, PerPage: 25}, {Page: 1, PerPage: 101}} {
		if bad.Valid() {
			t.Fatal("unbounded")
		}
	}
	a, b := time.Now(), time.Now().Add(-time.Hour)
	valid.From = &a
	valid.To = &b
	if valid.Valid() {
		t.Fatal("reversed dates")
	}
}
