package ie_test

import (
	"bytes"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestCipheringKeySequenceNumberFields(t *testing.T) {
	t.Parallel()
	fields := &ie.CipheringKeySequenceNumberFields{KeySequenceNumber: 6}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x06}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewCipheringKeySequenceNumberWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.CipheringKeySequenceNumber()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("CipheringKeySequenceNumber = %+v, want %+v", decoded, fields)
	}
}

func TestPagingIdentityFields(t *testing.T) {
	t.Parallel()
	fields := &ie.PagingIdentityFields{Identity: ie.PagingIdentityTMSI}
	value, err := ie.NewPositionalPagingIdentityWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.PagingIdentity()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("PagingIdentity = %+v, want %+v", decoded, fields)
	}
}

func TestExtendedEMMCauseFields(t *testing.T) {
	t.Parallel()
	fields := &ie.ExtendedEMMCauseFields{
		EUTRANNotAllowed:                true,
		EPSOptimizationInfoNotSupported: true,
		NBIoTNotAllowed:                 true,
	}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x07}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewExtendedEMMCauseWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.ExtendedEMMCause()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("ExtendedEMMCause = %+v, want %+v", decoded, fields)
	}
}
