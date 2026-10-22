package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestUERadioCapabilityIDFields_RoundTrip(t *testing.T) {
	t.Parallel()
	payload := []byte{0x01, 0x02}
	fields := &ie.UERadioCapabilityIDFields{Value: payload}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, payload) {
		t.Fatalf("Marshal = %x, want %x", encoded, payload)
	}

	value, err := ie.NewUERadioCapabilityIDWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := value.UERadioCapabilityID()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("UERadioCapabilityID = %+v, want %+v", decoded, fields)
	}
}

func TestUERadioCapabilityIDFields_PositionalRoundTrip(t *testing.T) {
	t.Parallel()
	payload := []byte{0xaa, 0xbb}
	fields := &ie.UERadioCapabilityIDFields{Value: payload}

	value, err := ie.NewPositionalUERadioCapabilityIDWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := value.UERadioCapabilityID()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("UERadioCapabilityID = %+v, want %+v", decoded, fields)
	}
}

func TestUERadioCapabilityIDFields_BelowMinError(t *testing.T) {
	t.Parallel()
	fields := &ie.UERadioCapabilityIDFields{Value: nil}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseUERadioCapabilityIDFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse nil: error = %v, want ErrInvalidLength", err)
	}
}

func TestUERadioCapabilityIDFields_AboveMaxError(t *testing.T) {
	t.Parallel()
	payload := make([]byte, 256)
	for i := range payload {
		payload[i] = byte(i)
	}
	fields := &ie.UERadioCapabilityIDFields{Value: payload}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseUERadioCapabilityIDFields(payload); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse: error = %v, want ErrInvalidLength", err)
	}
}

func TestUERadioCapabilityIDFields_NilFields(t *testing.T) {
	t.Parallel()
	var nilFields *ie.UERadioCapabilityIDFields
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

func TestUERadioCapabilityIDFields_MarshalTo(t *testing.T) {
	t.Parallel()
	payload := []byte{0x10, 0x20}
	fields := &ie.UERadioCapabilityIDFields{Value: payload}

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
