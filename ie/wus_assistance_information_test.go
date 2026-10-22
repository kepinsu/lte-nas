package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestWUSAssistanceInformationFields_RoundTrip(t *testing.T) {
	t.Parallel()
	payload := []byte{0x01, 0x02}
	fields := &ie.WUSAssistanceInformationFields{Value: payload}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, payload) {
		t.Fatalf("Marshal = %x, want %x", encoded, payload)
	}

	value, err := ie.NewWUSAssistanceInformationWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := value.WUSAssistanceInformation()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("WUSAssistanceInformation = %+v, want %+v", decoded, fields)
	}
}

func TestWUSAssistanceInformationFields_PositionalRoundTrip(t *testing.T) {
	t.Parallel()
	payload := []byte{0xaa, 0xbb}
	fields := &ie.WUSAssistanceInformationFields{Value: payload}

	value, err := ie.NewPositionalWUSAssistanceInformationWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := value.WUSAssistanceInformation()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("WUSAssistanceInformation = %+v, want %+v", decoded, fields)
	}
}

func TestWUSAssistanceInformationFields_BelowMinError(t *testing.T) {
	t.Parallel()
	fields := &ie.WUSAssistanceInformationFields{Value: nil}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseWUSAssistanceInformationFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse nil: error = %v, want ErrInvalidLength", err)
	}
}

func TestWUSAssistanceInformationFields_AboveMaxError(t *testing.T) {
	t.Parallel()
	payload := make([]byte, 256)
	for i := range payload {
		payload[i] = byte(i)
	}
	fields := &ie.WUSAssistanceInformationFields{Value: payload}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseWUSAssistanceInformationFields(payload); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse: error = %v, want ErrInvalidLength", err)
	}
}

func TestWUSAssistanceInformationFields_NilFields(t *testing.T) {
	t.Parallel()
	var nilFields *ie.WUSAssistanceInformationFields
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

func TestWUSAssistanceInformationFields_MarshalTo(t *testing.T) {
	t.Parallel()
	payload := []byte{0x10, 0x20}
	fields := &ie.WUSAssistanceInformationFields{Value: payload}

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
