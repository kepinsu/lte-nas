package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestEmergencyNumberListFields_RoundTrip(t *testing.T) {
	t.Parallel()
	payload := []byte{0x01, 0x02, 0x03}
	fields := &ie.EmergencyNumberListFields{Value: payload}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, payload) {
		t.Fatalf("Marshal = %x, want %x", encoded, payload)
	}

	value, err := ie.NewEmergencyNumberListWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := value.EmergencyNumberList()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("EmergencyNumberList = %+v, want %+v", decoded, fields)
	}
}

func TestEmergencyNumberListFields_PositionalRoundTrip(t *testing.T) {
	t.Parallel()
	payload := []byte{0xaa, 0xbb, 0xcc}
	fields := &ie.EmergencyNumberListFields{Value: payload}

	value, err := ie.NewPositionalEmergencyNumberListWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := value.EmergencyNumberList()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("EmergencyNumberList = %+v, want %+v", decoded, fields)
	}
}

func TestEmergencyNumberListFields_BelowMinError(t *testing.T) {
	t.Parallel()
	fields := &ie.EmergencyNumberListFields{Value: nil}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseEmergencyNumberListFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse nil: error = %v, want ErrInvalidLength", err)
	}

	twoBytes := []byte{0x01, 0x02}
	fieldsTwo := &ie.EmergencyNumberListFields{Value: twoBytes}
	if err := fieldsTwo.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate 2 bytes: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseEmergencyNumberListFields(twoBytes); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse 2 bytes: error = %v, want ErrInvalidLength", err)
	}
}

func TestEmergencyNumberListFields_AboveMaxError(t *testing.T) {
	t.Parallel()
	payload := make([]byte, 49)
	for i := range payload {
		payload[i] = byte(i)
	}
	fields := &ie.EmergencyNumberListFields{Value: payload}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseEmergencyNumberListFields(payload); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse: error = %v, want ErrInvalidLength", err)
	}
}

func TestEmergencyNumberListFields_NilFields(t *testing.T) {
	t.Parallel()
	var nilFields *ie.EmergencyNumberListFields
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

func TestEmergencyNumberListFields_MarshalTo(t *testing.T) {
	t.Parallel()
	payload := []byte{0x10, 0x20, 0x30}
	fields := &ie.EmergencyNumberListFields{Value: payload}

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
