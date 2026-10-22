package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestProtocolConfigurationOptionsFields_RoundTrip(t *testing.T) {
	t.Parallel()
	payload := []byte{0x01, 0x02, 0x03}
	fields := &ie.ProtocolConfigurationOptionsFields{Value: payload}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, payload) {
		t.Fatalf("Marshal = %x, want %x", encoded, payload)
	}

	value, err := ie.NewProtocolConfigurationOptionsWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := value.ProtocolConfigurationOptions()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("ProtocolConfigurationOptions = %+v, want %+v", decoded, fields)
	}
}

func TestProtocolConfigurationOptionsFields_PositionalRoundTrip(t *testing.T) {
	t.Parallel()
	payload := []byte{0xaa, 0xbb, 0xcc}
	fields := &ie.ProtocolConfigurationOptionsFields{Value: payload}

	value, err := ie.NewPositionalProtocolConfigurationOptionsWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := value.ProtocolConfigurationOptions()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("ProtocolConfigurationOptions = %+v, want %+v", decoded, fields)
	}
}

func TestProtocolConfigurationOptionsFields_MinLength(t *testing.T) {
	t.Parallel()
	payload := []byte{0x01}
	fields := &ie.ProtocolConfigurationOptionsFields{Value: payload}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, payload) {
		t.Fatalf("Marshal = %x, want %x", encoded, payload)
	}

	parsed, err := ie.ParseProtocolConfigurationOptionsFields(payload)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(parsed, fields) {
		t.Fatalf("Parse = %+v, want %+v", parsed, fields)
	}
}

func TestProtocolConfigurationOptionsFields_MaxLength(t *testing.T) {
	t.Parallel()
	payload := make([]byte, 251)
	for i := range payload {
		payload[i] = byte(i)
	}
	fields := &ie.ProtocolConfigurationOptionsFields{Value: payload}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if len(encoded) != 251 {
		t.Fatalf("Marshal length = %d, want 251", len(encoded))
	}

	parsed, err := ie.ParseProtocolConfigurationOptionsFields(payload)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(parsed, fields) {
		t.Fatalf("Parse = %+v, want %+v", parsed, fields)
	}
}

func TestProtocolConfigurationOptionsFields_BelowMinError(t *testing.T) {
	t.Parallel()
	fields := &ie.ProtocolConfigurationOptionsFields{Value: nil}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseProtocolConfigurationOptionsFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse nil: error = %v, want ErrInvalidLength", err)
	}
}

func TestProtocolConfigurationOptionsFields_AboveMaxError(t *testing.T) {
	t.Parallel()
	payload := make([]byte, 252)
	fields := &ie.ProtocolConfigurationOptionsFields{Value: payload}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseProtocolConfigurationOptionsFields(payload); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse: error = %v, want ErrInvalidLength", err)
	}
}

func TestProtocolConfigurationOptionsFields_NilFields(t *testing.T) {
	t.Parallel()
	var nilFields *ie.ProtocolConfigurationOptionsFields
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

func TestProtocolConfigurationOptionsFields_NilValue(t *testing.T) {
	t.Parallel()
	fields := &ie.ProtocolConfigurationOptionsFields{Value: nil}
	if fields.MarshalLen() != 0 {
		t.Fatalf("MarshalLen nil Value = %d, want 0", fields.MarshalLen())
	}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate nil Value: error = %v, want ErrInvalidLength", err)
	}
}

func TestProtocolConfigurationOptionsFields_MarshalTo(t *testing.T) {
	t.Parallel()
	payload := []byte{0x10, 0x20, 0x30}
	fields := &ie.ProtocolConfigurationOptionsFields{Value: payload}

	buffer := make([]byte, 10)
	if err := fields.MarshalTo(buffer); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buffer[:len(payload)], payload) {
		t.Fatalf("MarshalTo = %x, want %x", buffer[:len(payload)], payload)
	}

	short := make([]byte, 2)
	if err := fields.MarshalTo(short); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo short buffer: error = %v, want io.ErrShortBuffer", err)
	}
}
