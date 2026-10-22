package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestMobileStationClassmark3Fields_RoundTrip(t *testing.T) {
	t.Parallel()
	payload := []byte{0x01, 0x02}
	fields := &ie.MobileStationClassmark3Fields{Value: payload}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, payload) {
		t.Fatalf("Marshal = %x, want %x", encoded, payload)
	}

	value, err := ie.NewMobileStationClassmark3WithFields(fields)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := value.MobileStationClassmark3()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("MobileStationClassmark3 = %+v, want %+v", decoded, fields)
	}
}

func TestMobileStationClassmark3Fields_PositionalRoundTrip(t *testing.T) {
	t.Parallel()
	payload := []byte{0xaa, 0xbb}
	fields := &ie.MobileStationClassmark3Fields{Value: payload}

	value, err := ie.NewPositionalMobileStationClassmark3WithFields(fields)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := value.MobileStationClassmark3()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("MobileStationClassmark3 = %+v, want %+v", decoded, fields)
	}
}

func TestMobileStationClassmark3Fields_EmptyPayload(t *testing.T) {
	t.Parallel()
	fields := &ie.MobileStationClassmark3Fields{Value: nil}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if len(encoded) != 0 {
		t.Fatalf("Marshal empty = %d bytes, want 0", len(encoded))
	}

	parsed, err := ie.ParseMobileStationClassmark3Fields(nil)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(parsed, fields) {
		t.Fatalf("Parse = %+v, want %+v", parsed, fields)
	}
}

func TestMobileStationClassmark3Fields_MaxLength(t *testing.T) {
	t.Parallel()
	payload := make([]byte, 32)
	for i := range payload {
		payload[i] = byte(i)
	}
	fields := &ie.MobileStationClassmark3Fields{Value: payload}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, payload) {
		t.Fatalf("Marshal = %x, want %x", encoded, payload)
	}

	parsed, err := ie.ParseMobileStationClassmark3Fields(payload)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(parsed, fields) {
		t.Fatalf("Parse = %+v, want %+v", parsed, fields)
	}
}

func TestMobileStationClassmark3Fields_AboveMaxError(t *testing.T) {
	t.Parallel()
	payload := make([]byte, 33)
	fields := &ie.MobileStationClassmark3Fields{Value: payload}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseMobileStationClassmark3Fields(payload); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse: error = %v, want ErrInvalidLength", err)
	}
}

func TestMobileStationClassmark3Fields_NilFields(t *testing.T) {
	t.Parallel()
	var nilFields *ie.MobileStationClassmark3Fields
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

func TestMobileStationClassmark3Fields_MarshalTo(t *testing.T) {
	t.Parallel()
	payload := []byte{0x10, 0x20}
	fields := &ie.MobileStationClassmark3Fields{Value: payload}

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
