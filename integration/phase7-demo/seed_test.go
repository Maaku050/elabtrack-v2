package main

import (
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/catalogimage"
	"testing"
)

func TestLocalIllustrationsAreCanonicalBoundedAndDistinct(t *testing.T) {
	seen := map[string]bool{}
	for n := 0; n < 8; n++ {
		v, e := (catalogimage.Validator{}).Validate(illustration(n))
		if e != nil || v.Width != 320 || v.Height != 240 || seen[v.Hash] {
			t.Fatalf("invalid illustration%d", n)
		}
		seen[v.Hash] = true
	}
}
