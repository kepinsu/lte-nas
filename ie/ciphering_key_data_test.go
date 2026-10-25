package ie_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestCipheringKeyDataFields_RoundTrip(t *testing.T) {
	t.Parallel()
	payload := make([]byte, 32)
	for i := range payload {
		payload[i] = byte(i)
	}
	fields := &ie.CipheringKeyDataFields{Value: payload}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, payload) {
		t.Fatalf("Marshal = %x, want %x", encoded, payload)
	}

	value, err := ie.NewPositionalCipheringKeyDataWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := value.CipheringKeyData()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("CipheringKeyData = %+v, want %+v", decoded, fields)
	}
}

func TestCipheringKeyDataFields_BelowMinError(t *testing.T) {
	t.Parallel()
	payload := make([]byte, 31)
	for i := range payload {
		payload[i] = byte(i)
	}
	fields := &ie.CipheringKeyDataFields{Value: payload}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseCipheringKeyDataFields(payload); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse: error = %v, want ErrInvalidLength", err)
	}
}

func TestCipheringKeyDataFields_AboveMaxError(t *testing.T) {
	t.Parallel()
	payload := make([]byte, 2289)
	for i := range payload {
		payload[i] = byte(i & 0xff)
	}
	fields := &ie.CipheringKeyDataFields{Value: payload}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseCipheringKeyDataFields(payload); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse: error = %v, want ErrInvalidLength", err)
	}
}

func TestCipheringKeyDataFields_NilFields(t *testing.T) {
	t.Parallel()
	var nilFields *ie.CipheringKeyDataFields
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

func TestCipheringKeyDataFields_MarshalTo(t *testing.T) {
	t.Parallel()
	payload := []byte{0x10, 0x20, 0x30}
	fields := &ie.CipheringKeyDataFields{Value: payload}

	// Should fail because value is below minimum (32 bytes)
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate short payload: error = %v, want ErrInvalidLength", err)
	}
}

func TestCipheringKeyDataFields_NilValue(t *testing.T) {
	t.Parallel()
	fields := &ie.CipheringKeyDataFields{Value: nil}
	if fields.MarshalLen() != 0 {
		t.Fatalf("MarshalLen nil Value = %d, want 0", fields.MarshalLen())
	}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate nil Value: error = %v, want ErrInvalidLength", err)
	}
}
