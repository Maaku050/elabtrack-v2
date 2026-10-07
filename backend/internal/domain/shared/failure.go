package shared

import (
	"errors"
	"fmt"
)

// InternalFailure retains safe operation/type diagnostics, discarding raw
// infrastructure text and values at the boundary. It is transport independent.
type InternalFailure struct{ Operation, CauseType string }

func (e *InternalFailure) Error() string { return "internal server error" }
func (e *InternalFailure) Unwrap() error { return ErrInternal }
func Internal(operation string, cause error) error {
	class, _ := FailureDetails(cause)
	return &InternalFailure{Operation: operation, CauseType: class}
}

// FailureDetails inspects type metadata only; never invokes Error or logs values.
func FailureDetails(err error) (class, operation string) {
	var failure *InternalFailure
	if errors.As(err, &failure) && failure != nil {
		return failure.CauseType, failure.Operation
	}
	for i := 0; i < 8 && err != nil; i++ {
		next := errors.Unwrap(err)
		if next == nil {
			break
		}
		err = next
	}
	return fmt.Sprintf("%T", err), ""
}
