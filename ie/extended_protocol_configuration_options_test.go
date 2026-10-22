package ie_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestExtendedProtocolConfigurationOptionsFields_RoundTrip(t *testing.T) {
	t.Parallel()
	payload := []byte{0x01, 0x02, 0x03}
	fields := &ie.ExtendedProtocolConfigurationOptionsFields{Value: payload}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, payload) {
		t.Fatalf("Marshal = %x, want %x", encoded, payload)
	}

	value, err := ie.NewExtendedProtocolConfigurationOptionsWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := value.ExtendedProtocolConfigurationOptions()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("ExtendedProtocolConfigurationOptions = %+v, want %+v", decoded, fields)
	}
}

func TestExtendedProtocolConfigurationOptionsFields_PositionalRoundTrip(t *testing.T) {
	t.Parallel()
	payload := []byte{0xaa, 0xbb, 0xcc}
	fields := &ie.ExtendedProtocolConfigurationOptionsFields{Value: payload}

	value, err := ie.NewPositionalExtendedProtocolConfigurationOptionsWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := value.ExtendedProtocolConfigurationOptions()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("ExtendedProtocolConfigurationOptions = %+v, want %+v", decoded, fields)
	}
}

func TestExtendedProtocolConfigurationOptionsFields_MinLength(t *testing.T) {
	t.Parallel()
	payload := []byte{0x01}
	fields := &ie.ExtendedProtocolConfigurationOptionsFields{Value: payload}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, payload) {
		t.Fatalf("Marshal = %x, want %x", encoded, payload)
	}

	parsed, err := ie.ParseExtendedProtocolConfigurationOptionsFields(payload)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(parsed, fields) {
		t.Fatalf("Parse = %+v, want %+v", parsed, fields)
	}
}

func TestExtendedProtocolConfigurationOptionsFields_MaxLength(t *testing.T) {
	t.Parallel()
	payload := make([]byte, 65535)
	for i := range payload {
		payload[i] = byte(i & 0xff)
	}
	fields := &ie.ExtendedProtocolConfigurationOptionsFields{Value: payload}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if len(encoded) != 65535 {
		t.Fatalf("Marshal length = %d, want 65535", len(encoded))
	}

	parsed, err := ie.ParseExtendedProtocolConfigurationOptionsFields(payload)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(parsed, fields) {
		t.Fatalf("Parse = %+v, want %+v", parsed, fields)
	}
}

func TestExtendedProtocolConfigurationOptionsFields_BelowMinError(t *testing.T) {
	t.Parallel()
	fields := &ie.ExtendedProtocolConfigurationOptionsFields{Value: nil}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseExtendedProtocolConfigurationOptionsFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse nil: error = %v, want ErrInvalidLength", err)
	}
}

func TestExtendedProtocolConfigurationOptionsFields_AboveMaxError(t *testing.T) {
	t.Parallel()
	payload := make([]byte, 65536)
	fields := &ie.ExtendedProtocolConfigurationOptionsFields{Value: payload}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseExtendedProtocolConfigurationOptionsFields(payload); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse: error = %v, want ErrInvalidLength", err)
	}
}

func TestExtendedProtocolConfigurationOptionsFields_NilFields(t *testing.T) {
	t.Parallel()
	var nilFields *ie.ExtendedProtocolConfigurationOptionsFields
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

func TestExtendedProtocolConfigurationOptionsFields_NilValue(t *testing.T) {
	t.Parallel()
	fields := &ie.ExtendedProtocolConfigurationOptionsFields{Value: nil}
	if fields.MarshalLen() != 0 {
		t.Fatalf("MarshalLen nil Value = %d, want 0", fields.MarshalLen())
	}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate nil Value: error = %v, want ErrInvalidLength", err)
	}
}

func TestExtendedProtocolConfigurationOptionsFields_MarshalTo(t *testing.T) {
	t.Parallel()
	payload := []byte{0x10, 0x20, 0x30}
	fields := &ie.ExtendedProtocolConfigurationOptionsFields{Value: payload}

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

func TestExtendedProtocolConfigurationOptionsFields_Tagged(t *testing.T) {
	t.Parallel()
	payload := []byte{0x01, 0x02, 0x03}
	fields := &ie.ExtendedProtocolConfigurationOptionsFields{Value: payload}

	encoded, err := fields.MarshalTagged()
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x7b, 0x00, 0x03, 0x01, 0x02, 0x03}
	if !bytes.Equal(encoded, want) {
		t.Fatalf("MarshalTagged = %x, want %x", encoded, want)
	}

	var decoded ie.ExtendedProtocolConfigurationOptionsFields
	if err := decoded.UnmarshalBinaryTagged(encoded); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(&decoded, fields) {
		t.Fatalf("UnmarshalBinaryTagged = %+v, want %+v", &decoded, fields)
	}
}

func TestExtendedProtocolConfigurationOptionsFields_MarshalTaggedLong(t *testing.T) {
	t.Parallel()
	payload := make([]byte, 65535)
	fields := &ie.ExtendedProtocolConfigurationOptionsFields{Value: payload}

	encoded, err := fields.MarshalTagged()
	if err != nil {
		t.Fatal(err)
	}
	if len(encoded) != 3+65535 {
		t.Fatalf("MarshalTagged length = %d, want %d", len(encoded), 3+65535)
	}

	short := make([]byte, 2)
	if err := fields.MarshalTo(short); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo short buffer: error = %v, want io.ErrShortBuffer", err)
	}

	var decoded ie.ExtendedProtocolConfigurationOptionsFields
	if err := decoded.UnmarshalBinaryTagged(encoded); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(&decoded, fields) {
		t.Fatalf("UnmarshalBinaryTagged long = %+v, want %+v", &decoded, fields)
	}

	if len(short) < len(fields.Value) {
		_ = binary.BigEndian.Uint16(short[:2]) // ensure no compile warning for unused variable
	}
}
