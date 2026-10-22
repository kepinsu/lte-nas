package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestUENetworkCapabilityFields_RoundTrip(t *testing.T) {
	t.Parallel()
	// 0x81: bit7=CN type, bits5-6=sparse(0), bit0=data. Valid.
	payload := []byte{0x81, 0xff, 0xaa}
	fields := &ie.UENetworkCapabilityFields{Value: payload}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, payload) {
		t.Fatalf("Marshal = %x, want %x", encoded, payload)
	}

	value, err := ie.NewUENetworkCapabilityWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := value.UENetworkCapability()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("UENetworkCapability = %+v, want %+v", decoded, fields)
	}
}

func TestUENetworkCapabilityFields_PositionalRoundTrip(t *testing.T) {
	t.Parallel()
	payload := []byte{0x81, 0xff}
	fields := &ie.UENetworkCapabilityFields{Value: payload}

	value, err := ie.NewPositionalUENetworkCapabilityWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := value.UENetworkCapability()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("UENetworkCapability = %+v, want %+v", decoded, fields)
	}
}

func TestUENetworkCapabilityFields_Parse(t *testing.T) {
	t.Parallel()
	payload := []byte{0x81, 0x0a}
	fields, err := ie.ParseUENetworkCapabilityFields(payload)
	if err != nil {
		t.Fatal(err)
	}
	want := &ie.UENetworkCapabilityFields{Value: payload}
	if !reflect.DeepEqual(fields, want) {
		t.Fatalf("Parse = %+v, want %+v", fields, want)
	}
	reencoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(reencoded, payload) {
		t.Fatalf("Marshal = %x, want %x", reencoded, payload)
	}
}

func TestUENetworkCapabilityFields_MinLength(t *testing.T) {
	t.Parallel()
	payload := []byte{0x81, 0x02}
	fields := &ie.UENetworkCapabilityFields{Value: payload}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, payload) {
		t.Fatalf("Marshal = %x, want %x", encoded, payload)
	}

	parsed, err := ie.ParseUENetworkCapabilityFields(payload)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(parsed, fields) {
		t.Fatalf("Parse = %+v, want %+v", parsed, fields)
	}
}

func TestUENetworkCapabilityFields_MaxLength(t *testing.T) {
	t.Parallel()
	payload := make([]byte, 13)
	for i := range payload {
		if i == 0 {
			payload[i] = byte(i) & 0x9f // clear spare bits 5-6
		} else {
			payload[i] = byte(i)
		}
	}
	fields := &ie.UENetworkCapabilityFields{Value: payload}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if len(encoded) != 13 {
		t.Fatalf("Marshal length = %d, want 13", len(encoded))
	}

	parsed, err := ie.ParseUENetworkCapabilityFields(payload)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(parsed, fields) {
		t.Fatalf("Parse = %+v, want %+v", parsed, fields)
	}
}

func TestUENetworkCapabilityFields_BelowMinError(t *testing.T) {
	t.Parallel()
	fields := &ie.UENetworkCapabilityFields{Value: nil}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseUENetworkCapabilityFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse nil: error = %v, want ErrInvalidLength", err)
	}
}

func TestUENetworkCapabilityFields_AboveMaxError(t *testing.T) {
	t.Parallel()
	payload := make([]byte, 14)
	fields := &ie.UENetworkCapabilityFields{Value: payload}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseUENetworkCapabilityFields(payload); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse: error = %v, want ErrInvalidLength", err)
	}
}

func TestUENetworkCapabilityFields_NilFields(t *testing.T) {
	t.Parallel()
	var nilFields *ie.UENetworkCapabilityFields
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

func TestUENetworkCapabilityFields_NilValue(t *testing.T) {
	t.Parallel()
	fields := &ie.UENetworkCapabilityFields{Value: nil}
	if fields.MarshalLen() != 0 {
		t.Fatalf("MarshalLen nil Value = %d, want 0", fields.MarshalLen())
	}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate nil Value: error = %v, want ErrInvalidLength", err)
	}
}

func TestUENetworkCapabilityFields_MarshalTo(t *testing.T) {
	t.Parallel()
	payload := []byte{0x81, 0xff, 0xaa}
	fields := &ie.UENetworkCapabilityFields{Value: payload}

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

func TestUENetworkCapabilityFields_SpareBitsError(t *testing.T) {
	t.Parallel()
	// Byte 1 with spare bits 5-6 set (0x60)
	payload := []byte{0x60, 0x01}
	if _, err := ie.ParseUENetworkCapabilityFields(payload); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Parse spare bits: error = %v, want ErrInvalidValue", err)
	}
	fields := &ie.UENetworkCapabilityFields{Value: payload}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Validate spare bits: error = %v, want ErrInvalidValue", err)
	}
}
