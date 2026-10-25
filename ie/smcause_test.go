package ie

import (
	"errors"
	"testing"
)

func TestIE_ESMCause(t *testing.T) {
	t.Parallel()
	i := NewESMCause(NetworkFailureSmCause)
	if got, err := i.ESMCause(); err != nil || got != NetworkFailureSmCause {
		t.Fatalf("ESMCause = %d, %v", got, err)
	}
	if _, err := NewEMMCause(NetworkFailureMmCause).ESMCause(); !errors.Is(err, ErrUnexpectedIE) {
		t.Fatalf("wrong ESM kind = %v", err)
	}
	var nilIE *IE
	if _, err := nilIE.ESMCause(); !errors.Is(err, ErrNilIE) {
		t.Fatalf("nil ESM cause = %v", err)
	}
}
