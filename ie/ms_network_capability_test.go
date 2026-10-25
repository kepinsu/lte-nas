package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestMSNetworkCapabilityFields_RoundTrip(t *testing.T) {
	t.Parallel()
	// 0x83: bit7=CN type, bits2-6=sparse(0), bits1+0=caps. Valid.
	payload := []byte{0x83, 0xff, 0xaa}
	fields := &ie.MSNetworkCapabilityFields{Value: payload}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, payload) {
		t.Fatalf("Marshal = %x, want %x", encoded, payload)
	}

	value, err := ie.NewMSNetworkCapabilityWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := value.MSNetworkCapability()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("MSNetworkCapability = %+v, want %+v", decoded, fields)
	}
}

func TestMSNetworkCapabilityFields_PositionalRoundTrip(t *testing.T) {
	t.Parallel()
	payload := []byte{0x83, 0xff}
	fields := &ie.MSNetworkCapabilityFields{Value: payload}

	value, err := ie.NewPositionalMSNetworkCapabilityWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := value.MSNetworkCapability()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("MSNetworkCapability = %+v, want %+v", decoded, fields)
	}
}

func TestMSNetworkCapabilityFields_Parse(t *testing.T) {
	t.Parallel()
	payload := []byte{0x83, 0x0a}
	fields, err := ie.ParseMSNetworkCapabilityFields(payload)
	if err != nil {
		t.Fatal(err)
	}
	want := &ie.MSNetworkCapabilityFields{Value: payload}
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

func TestMSNetworkCapabilityFields_MinLength(t *testing.T) {
	t.Parallel()
	payload := []byte{0x83, 0x02}
	fields := &ie.MSNetworkCapabilityFields{Value: payload}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, payload) {
		t.Fatalf("Marshal = %x, want %x", encoded, payload)
	}

	parsed, err := ie.ParseMSNetworkCapabilityFields(payload)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(parsed, fields) {
		t.Fatalf("Parse = %+v, want %+v", parsed, fields)
	}
}

func TestMSNetworkCapabilityFields_MaxLength(t *testing.T) {
	t.Parallel()
	payload := make([]byte, 8)
	for i := range payload {
		if i == 0 {
			payload[i] = byte(i) & 0x83 // clear spare bits 2-6
		} else {
			payload[i] = byte(i)
		}
	}
	fields := &ie.MSNetworkCapabilityFields{Value: payload}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if len(encoded) != 8 {
		t.Fatalf("Marshal length = %d, want 8", len(encoded))
	}

	parsed, err := ie.ParseMSNetworkCapabilityFields(payload)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(parsed, fields) {
		t.Fatalf("Parse = %+v, want %+v", parsed, fields)
	}
}

func TestMSNetworkCapabilityFields_BelowMinError(t *testing.T) {
	t.Parallel()
	fields := &ie.MSNetworkCapabilityFields{Value: nil}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseMSNetworkCapabilityFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse nil: error = %v, want ErrInvalidLength", err)
	}
}

func TestMSNetworkCapabilityFields_AboveMaxError(t *testing.T) {
	t.Parallel()
	payload := make([]byte, 9)
	fields := &ie.MSNetworkCapabilityFields{Value: payload}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseMSNetworkCapabilityFields(payload); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse: error = %v, want ErrInvalidLength", err)
	}
}

func TestMSNetworkCapabilityFields_NilFields(t *testing.T) {
	t.Parallel()
	var nilFields *ie.MSNetworkCapabilityFields
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

func TestMSNetworkCapabilityFields_NilValue(t *testing.T) {
	t.Parallel()
	fields := &ie.MSNetworkCapabilityFields{Value: nil}
	if fields.MarshalLen() != 0 {
		t.Fatalf("MarshalLen nil Value = %d, want 0", fields.MarshalLen())
	}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate nil Value: error = %v, want ErrInvalidLength", err)
	}
}

func TestMSNetworkCapabilityFields_MarshalTo(t *testing.T) {
	t.Parallel()
	payload := []byte{0x83, 0xff, 0xaa}
	fields := &ie.MSNetworkCapabilityFields{Value: payload}

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

func TestMSNetworkCapabilityFields_SpareBitsError(t *testing.T) {
	t.Parallel()
	// Byte 1 with spare bits 2-6 set (0x7c)
	payload := []byte{0x7c, 0x01}
	if _, err := ie.ParseMSNetworkCapabilityFields(payload); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Parse spare bits: error = %v, want ErrInvalidValue", err)
	}
	fields := &ie.MSNetworkCapabilityFields{Value: payload}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Validate spare bits: error = %v, want ErrInvalidValue", err)
	}
}
