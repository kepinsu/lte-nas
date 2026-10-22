package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestNASMessageContainerFields_RoundTrip(t *testing.T) {
	t.Parallel()
	payload := []byte{0x01, 0x02}
	fields := &ie.NASMessageContainerFields{Value: payload}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, payload) {
		t.Fatalf("Marshal = %x, want %x", encoded, payload)
	}

	value, err := ie.NewNASMessageContainerWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := value.NASMessageContainer()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("NASMessageContainer = %+v, want %+v", decoded, fields)
	}
}

func TestNASMessageContainerFields_PositionalRoundTrip(t *testing.T) {
	t.Parallel()
	payload := []byte{0xaa, 0xbb}
	fields := &ie.NASMessageContainerFields{Value: payload}

	value, err := ie.NewPositionalNASMessageContainerWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := value.NASMessageContainer()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("NASMessageContainer = %+v, want %+v", decoded, fields)
	}
}

func TestNASMessageContainerFields_BelowMinError(t *testing.T) {
	t.Parallel()
	fields := &ie.NASMessageContainerFields{Value: nil}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseNASMessageContainerFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse nil: error = %v, want ErrInvalidLength", err)
	}

	singleByte := []byte{0x01}
	fieldsSingle := &ie.NASMessageContainerFields{Value: singleByte}
	if err := fieldsSingle.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate 1 byte: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseNASMessageContainerFields(singleByte); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse 1 byte: error = %v, want ErrInvalidLength", err)
	}
}

func TestNASMessageContainerFields_AboveMaxError(t *testing.T) {
	t.Parallel()
	payload := make([]byte, 252)
	for i := range payload {
		payload[i] = byte(i)
	}
	fields := &ie.NASMessageContainerFields{Value: payload}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseNASMessageContainerFields(payload); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse: error = %v, want ErrInvalidLength", err)
	}
}

func TestNASMessageContainerFields_NilFields(t *testing.T) {
	t.Parallel()
	var nilFields *ie.NASMessageContainerFields
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

func TestNASMessageContainerFields_MarshalTo(t *testing.T) {
	t.Parallel()
	payload := []byte{0x10, 0x20}
	fields := &ie.NASMessageContainerFields{Value: payload}

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
