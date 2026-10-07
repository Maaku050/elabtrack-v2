package shared

import (
	"errors"
	"testing"
)

type credentialError struct{}
type safeWrapper struct{ cause error }

func (e safeWrapper) Error() string { return "secret wrapper" }
func (e safeWrapper) Unwrap() error { return e.cause }

func (credentialError) Error() string {
	panic("raw credential Error() must never be used for diagnostics")
}
func TestSafeInternalFailureMetadata(t *testing.T) {
	err := Internal("synthetic.repository", safeWrapper{credentialError{}})
	class, operation := FailureDetails(err)
	if !errors.Is(err, ErrInternal) || err.Error() != "internal server error" || class != "shared.credentialError" || operation != "synthetic.repository" {
		t.Fatal("safe failure metadata lost")
	}
	class, _ = FailureDetails(safeWrapper{err})
	if class != "shared.credentialError" {
		t.Fatal("wrapped type diagnostics lost")
	}
}
