package ie

import (
	"errors"
	"testing"
)

func TestIE_EMMCause(t *testing.T) {
	t.Parallel()
	i := NewEMMCause(NetworkFailureMmCause)
	if got, err := i.EMMCause(); err != nil || got != NetworkFailureMmCause {
		t.Fatalf("EMMCause = %d, %v", got, err)
	}
	if _, err := NewESMCause(NetworkFailureSmCause).EMMCause(); !errors.Is(err, ErrUnexpectedIE) {
		t.Fatalf("wrong EMM kind = %v", err)
	}
	var nilIE *IE
	if _, err := nilIE.EMMCause(); !errors.Is(err, ErrNilIE) {
		t.Fatalf("nil EMM cause = %v", err)
	}
}
