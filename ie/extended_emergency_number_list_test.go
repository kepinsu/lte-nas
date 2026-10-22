package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestExtendedEmergencyNumberListFields_RoundTrip(t *testing.T) {
	t.Parallel()
	payload := []byte{0x01, 0x02, 0x03, 0x04}
	fields := &ie.ExtendedEmergencyNumberListFields{Value: payload}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, payload) {
		t.Fatalf("Marshal = %x, want %x", encoded, payload)
	}

	value, err := ie.NewPositionalExtendedEmergencyNumberListWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := value.ExtendedEmergencyNumberList()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("ExtendedEmergencyNumberList = %+v, want %+v", decoded, fields)
	}
}

func TestExtendedEmergencyNumberListFields_MinLength(t *testing.T) {
	t.Parallel()
	payload := []byte{0x01, 0x02, 0x03, 0x04}
	fields := &ie.ExtendedEmergencyNumberListFields{Value: payload}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, payload) {
		t.Fatalf("Marshal = %x, want %x", encoded, payload)
	}

	parsed, err := ie.ParseExtendedEmergencyNumberListFields(payload)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(parsed, fields) {
		t.Fatalf("Parse = %+v, want %+v", parsed, fields)
	}
}

func TestExtendedEmergencyNumberListFields_AboveMaxError(t *testing.T) {
	t.Parallel()
	payload := []byte{0x01, 0x02, 0x03}
	fields := &ie.ExtendedEmergencyNumberListFields{Value: payload}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseExtendedEmergencyNumberListFields(payload); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse: error = %v, want ErrInvalidLength", err)
	}
}

func TestExtendedEmergencyNumberListFields_NilFields(t *testing.T) {
	t.Parallel()
	var nilFields *ie.ExtendedEmergencyNumberListFields
	if err := nilFields.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Validate nil: error = %v, want ErrInvalidValue", err)
	}
	_, err := nilFields.Marshal()
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Marshal nil: error = %v, want ErrInvalidValue", err)
	}
	if nilFields.MarshalLen() != 0 {
		t.Fatal("MarshalLen nil: expected 0")
	}
}

func TestExtendedEmergencyNumberListFields_MarshalTo(t *testing.T) {
	t.Parallel()
	payload := []byte{0x10, 0x20, 0x30, 0x40}
	fields := &ie.ExtendedEmergencyNumberListFields{Value: payload}

	buffer := make([]byte, 10)
	if err := fields.MarshalTo(buffer); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buffer[:len(payload)], payload) {
		t.Fatalf("MarshalTo = %x, want %x", buffer[:len(payload)], payload)
	}

	short := make([]byte, 3)
	if err := fields.MarshalTo(short); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo short buffer: error = %v, want io.ErrShortBuffer", err)
	}
}

func TestExtendedEmergencyNumberListFields_NilValue(t *testing.T) {
	t.Parallel()
	fields := &ie.ExtendedEmergencyNumberListFields{Value: nil}
	if fields.MarshalLen() != 0 {
		t.Fatalf("MarshalLen nil Value = %d, want 0", fields.MarshalLen())
	}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate nil Value: error = %v, want ErrInvalidLength", err)
	}
}
