package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestAdditionalInformationFields_RoundTrip(t *testing.T) {
	t.Parallel()
	payload := []byte{0x01, 0x02, 0xab}
	fields := &ie.AdditionalInformationFields{Value: payload}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, payload) {
		t.Fatalf("Marshal = %x, want %x", encoded, payload)
	}

	value, err := ie.NewAdditionalInformationWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := value.AdditionalInformation()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("AdditionalInformation = %+v, want %+v", decoded, fields)
	}
}

func TestAdditionalInformationFields_PositionalRoundTrip(t *testing.T) {
	t.Parallel()
	payload := []byte{0xaa, 0xbb}
	fields := &ie.AdditionalInformationFields{Value: payload}

	value, err := ie.NewPositionalAdditionalInformationWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := value.AdditionalInformation()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("AdditionalInformation = %+v, want %+v", decoded, fields)
	}
}

func TestAdditionalInformationFields_BelowMinError(t *testing.T) {
	t.Parallel()
	fields := &ie.AdditionalInformationFields{Value: nil}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseAdditionalInformationFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse nil: error = %v, want ErrInvalidLength", err)
	}
}

func TestAdditionalInformationFields_AboveMaxError(t *testing.T) {
	t.Parallel()
	payload := make([]byte, 256)
	for i := range payload {
		payload[i] = byte(i)
	}
	fields := &ie.AdditionalInformationFields{Value: payload}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseAdditionalInformationFields(payload); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse: error = %v, want ErrInvalidLength", err)
	}
}

func TestAdditionalInformationFields_NilFields(t *testing.T) {
	t.Parallel()
	var nilFields *ie.AdditionalInformationFields
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

func TestAdditionalInformationFields_NilValue(t *testing.T) {
	t.Parallel()
	fields := &ie.AdditionalInformationFields{Value: nil}
	if fields.MarshalLen() != 0 {
		t.Fatalf("MarshalLen nil Value = %d, want 0", fields.MarshalLen())
	}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate nil Value: error = %v, want ErrInvalidLength", err)
	}
}

func TestAdditionalInformationFields_MarshalTo(t *testing.T) {
	t.Parallel()
	payload := []byte{0x10, 0x20}
	fields := &ie.AdditionalInformationFields{Value: payload}

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
