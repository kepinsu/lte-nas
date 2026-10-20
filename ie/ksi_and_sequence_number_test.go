package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestKSIAndSequenceNumberFields(t *testing.T) {
	t.Parallel()
	fields := &ie.KSIAndSequenceNumberFields{KeySetIdentifier: 5, SequenceNumber: 17}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0xb1}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewPositionalKSIAndSequenceNumberWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.KSIAndSequenceNumber()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("KSIAndSequenceNumber = %+v, want %+v", decoded, fields)
	}
}

func TestKSIAndSequenceNumberFields_Invalid(t *testing.T) {
	t.Parallel()
	fields := &ie.KSIAndSequenceNumberFields{KeySetIdentifier: 8}
	if _, err := fields.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Marshal error = %v, want ErrInvalidValue", err)
	}
	if err := fields.MarshalTo(nil); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("MarshalTo error = %v, want ErrInvalidValue", err)
	}

	fields = &ie.KSIAndSequenceNumberFields{}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
