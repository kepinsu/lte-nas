package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestHeaderCompressionConfigurationFields_RoundTrip(t *testing.T) {
	t.Parallel()
	payload := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	fields := &ie.HeaderCompressionConfigurationFields{Value: payload}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, payload) {
		t.Fatalf("Marshal = %x, want %x", encoded, payload)
	}

	value, err := ie.NewHeaderCompressionConfigurationWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := value.HeaderCompressionConfiguration()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("HeaderCompressionConfiguration = %+v, want %+v", decoded, fields)
	}
}

func TestHeaderCompressionConfigurationFields_PositionalRoundTrip(t *testing.T) {
	t.Parallel()
	payload := []byte{0xaa, 0xbb, 0xcc}
	fields := &ie.HeaderCompressionConfigurationFields{Value: payload}

	value, err := ie.NewPositionalHeaderCompressionConfigurationWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := value.HeaderCompressionConfiguration()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("HeaderCompressionConfiguration = %+v, want %+v", decoded, fields)
	}
}

func TestHeaderCompressionConfigurationFields_MinLength(t *testing.T) {
	t.Parallel()
	payload := []byte{0x01, 0x02, 0x03}
	fields := &ie.HeaderCompressionConfigurationFields{Value: payload}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, payload) {
		t.Fatalf("Marshal = %x, want %x", encoded, payload)
	}

	parsed, err := ie.ParseHeaderCompressionConfigurationFields(payload)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(parsed, fields) {
		t.Fatalf("Parse = %+v, want %+v", parsed, fields)
	}
}

func TestHeaderCompressionConfigurationFields_MaxLength(t *testing.T) {
	t.Parallel()
	payload := make([]byte, 255)
	for i := range payload {
		payload[i] = byte(i)
	}
	fields := &ie.HeaderCompressionConfigurationFields{Value: payload}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if len(encoded) != 255 {
		t.Fatalf("Marshal length = %d, want 255", len(encoded))
	}

	parsed, err := ie.ParseHeaderCompressionConfigurationFields(payload)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(parsed, fields) {
		t.Fatalf("Parse = %+v, want %+v", parsed, fields)
	}
}

func TestHeaderCompressionConfigurationFields_BelowMinError(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"too short 2", "too short 1", "empty"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			payload := make([]byte, 2)
			if name == "empty" {
				payload = nil
			}
			fields := &ie.HeaderCompressionConfigurationFields{Value: payload}
			if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
				t.Fatalf("%s: Validate error = %v, want ErrInvalidLength", name, err)
			}
			if _, err := ie.ParseHeaderCompressionConfigurationFields(payload); !errors.Is(err, ie.ErrInvalidLength) {
				t.Fatalf("%s: Parse error = %v, want ErrInvalidLength", name, err)
			}
		})
	}
}

func TestHeaderCompressionConfigurationFields_AboveMaxError(t *testing.T) {
	t.Parallel()
	payload := make([]byte, 256)
	fields := &ie.HeaderCompressionConfigurationFields{Value: payload}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseHeaderCompressionConfigurationFields(payload); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse: error = %v, want ErrInvalidLength", err)
	}
}

func TestHeaderCompressionConfigurationFields_NilFields(t *testing.T) {
	t.Parallel()
	var nilFields *ie.HeaderCompressionConfigurationFields
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

func TestHeaderCompressionConfigurationFields_NilValue(t *testing.T) {
	t.Parallel()
	fields := &ie.HeaderCompressionConfigurationFields{Value: nil}
	if fields.MarshalLen() != 0 {
		t.Fatalf("MarshalLen nil Value = %d, want 0", fields.MarshalLen())
	}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate nil Value: error = %v, want ErrInvalidLength", err)
	}
}

func TestHeaderCompressionConfigurationFields_MarshalTo(t *testing.T) {
	t.Parallel()
	payload := []byte{0x10, 0x20, 0x30, 0x40}
	fields := &ie.HeaderCompressionConfigurationFields{Value: payload}

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
