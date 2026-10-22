package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestUserDataContainerFields_RoundTrip(t *testing.T) {
	t.Parallel()
	payload := []byte{0x01, 0x02}
	fields := &ie.UserDataContainerFields{Value: payload}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, payload) {
		t.Fatalf("Marshal = %x, want %x", encoded, payload)
	}

	value, err := ie.NewPositionalUserDataContainerWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := value.UserDataContainer()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("UserDataContainer = %+v, want %+v", decoded, fields)
	}
}

func TestUserDataContainerFields_EmptyPayload(t *testing.T) {
	t.Parallel()
	fields := &ie.UserDataContainerFields{Value: nil}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if len(encoded) != 0 {
		t.Fatalf("Marshal empty = %d bytes, want 0", len(encoded))
	}

	parsed, err := ie.ParseUserDataContainerFields(nil)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(parsed, fields) {
		t.Fatalf("Parse = %+v, want %+v", parsed, fields)
	}
}

func TestUserDataContainerFields_AboveMaxError(t *testing.T) {
	t.Parallel()
	payload := make([]byte, 65536)
	fields := &ie.UserDataContainerFields{Value: payload}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseUserDataContainerFields(payload); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse: error = %v, want ErrInvalidLength", err)
	}
}

func TestUserDataContainerFields_NilFields(t *testing.T) {
	t.Parallel()
	var nilFields *ie.UserDataContainerFields
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

func TestUserDataContainerFields_MarshalTo(t *testing.T) {
	t.Parallel()
	payload := []byte{0x10, 0x20}
	fields := &ie.UserDataContainerFields{Value: payload}

	buffer := make([]byte, 10)
	if err := fields.MarshalTo(buffer); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buffer[:len(payload)], payload) {
		t.Fatalf("MarshalTo = %x, want %x", buffer[:len(payload)], payload)
	}

	short := make([]byte, 1)
	if err := fields.MarshalTo(short); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo short buffer: error = %v, want io.ErrShortBuffer", err)
	}
}
