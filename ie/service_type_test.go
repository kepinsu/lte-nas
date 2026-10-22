package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestServiceTypeFields_RoundTrip(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name  string
		value ie.ServiceType
	}{
		{name: "voice", value: ie.ServiceTypeVoice},
		{name: "ss", value: ie.ServiceTypeSS},
		{name: "sms", value: ie.ServiceTypeSMS},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fields := &ie.ServiceTypeFields{Type: tc.value}

			encoded, err := fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(encoded, []byte{byte(tc.value)}) {
				t.Fatalf("Marshal = %x, want %02x", encoded, tc.value)
			}

			value, err := ie.NewPositionalServiceTypeWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}

			decoded, err := value.ServiceType()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("ServiceType = %+v, want %+v", decoded, fields)
			}
		})
	}
}

func TestServiceTypeFields_RejectReserved(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		payload []byte
	}{
		{name: "empty", payload: nil},
		{name: "reserved value 0x04", payload: []byte{0x04}},
		{name: "reserved full nibble", payload: []byte{0x0f}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fields := ie.ServiceTypeFields{Type: ie.ServiceTypeVoice}
			before := fields
			if err := fields.UnmarshalBinary(tc.payload); !errors.Is(err, ie.ErrInvalidValue) && !errors.Is(err, ie.ErrInvalidLength) {
				t.Fatalf("UnmarshalBinary = %v, want invalid value or length", err)
			}
			if fields != before {
				t.Fatalf("UnmarshalBinary modified fields: got %+v, want %+v", fields, before)
			}
		})
	}

	fields := &ie.ServiceTypeFields{Type: ie.ServiceTypeVoice}
	buffer := []byte{0xaa}
	if err := fields.MarshalTo(buffer[:0]); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo = %v, want io.ErrShortBuffer", err)
	}
	if buffer[0] != 0xaa {
		t.Fatalf("MarshalTo modified short buffer: got %x, want aa", buffer)
	}
}

func TestServiceTypeFields_NilFields(t *testing.T) {
	t.Parallel()
	var nilFields *ie.ServiceTypeFields
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

func TestServiceTypeFields_MarshalTo(t *testing.T) {
	t.Parallel()
	fields := &ie.ServiceTypeFields{Type: ie.ServiceTypeSMS}

	buffer := make([]byte, 10)
	if err := fields.MarshalTo(buffer); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buffer[:1], []byte{0x03}) {
		t.Fatalf("MarshalTo = %x, want 03", buffer[:1])
	}

	short := make([]byte, 0)
	if err := fields.MarshalTo(short); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo short buffer: error = %v, want io.ErrShortBuffer", err)
	}
}
