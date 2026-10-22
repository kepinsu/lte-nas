package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestDetachTypeFields_RoundTrip(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name  string
		value ie.DetachType
	}{
		{name: "EPS only", value: ie.DetachTypeEPSOnly},
		{name: "combined EPS IMSI", value: ie.DetachTypeCombinedEPSIMSI},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fields := &ie.DetachTypeFields{Type: tc.value}

			encoded, err := fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(encoded, []byte{byte(tc.value)}) {
				t.Fatalf("Marshal = %x, want %02x", encoded, tc.value)
			}

			value, err := ie.NewDetachTypeWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}

			decoded, err := value.DetachType()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("DetachType = %+v, want %+v", decoded, fields)
			}
		})
	}
}

func TestDetachTypeFields_PositionalRoundTrip(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name  string
		value ie.DetachType
	}{
		{name: "EPS only", value: ie.DetachTypeEPSOnly},
		{name: "combined EPS IMSI", value: ie.DetachTypeCombinedEPSIMSI},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fields := &ie.DetachTypeFields{Type: tc.value}

			value, err := ie.NewPositionalDetachTypeWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}

			decoded, err := value.DetachType()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("DetachType = %+v, want %+v", decoded, fields)
			}
		})
	}
}

func TestDetachTypeFields_RejectReserved(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		payload []byte
	}{
		{name: "empty", payload: nil},
		{name: "spare bit nonzero", payload: []byte{0x11}},
		{name: "unused value", payload: []byte{0x04}},
		{name: "reserved full nibble", payload: []byte{0x0f}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fields := ie.DetachTypeFields{Type: ie.DetachTypeEPSOnly}
			before := fields
			if err := fields.UnmarshalBinary(tc.payload); !errors.Is(err, ie.ErrInvalidValue) && !errors.Is(err, ie.ErrInvalidLength) {
				t.Fatalf("UnmarshalBinary = %v, want invalid value or length", err)
			}
			if fields != before {
				t.Fatalf("UnmarshalBinary modified fields: got %+v, want %+v", fields, before)
			}
		})
	}

	fields := &ie.DetachTypeFields{Type: ie.DetachTypeEPSOnly}
	buffer := []byte{0xaa}
	if err := fields.MarshalTo(buffer[:0]); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo = %v, want io.ErrShortBuffer", err)
	}
	if buffer[0] != 0xaa {
		t.Fatalf("MarshalTo modified short buffer: got %x, want aa", buffer)
	}
}

func TestDetachTypeFields_NilFields(t *testing.T) {
	t.Parallel()
	var nilFields *ie.DetachTypeFields
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

func TestDetachTypeFields_MarshalTo(t *testing.T) {
	t.Parallel()
	fields := &ie.DetachTypeFields{Type: ie.DetachTypeEPSOnly}

	buffer := make([]byte, 10)
	if err := fields.MarshalTo(buffer); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buffer[:1], []byte{0x01}) {
		t.Fatalf("MarshalTo = %x, want 01", buffer[:1])
	}

	short := make([]byte, 0)
	if err := fields.MarshalTo(short); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo short buffer: error = %v, want io.ErrShortBuffer", err)
	}
}
