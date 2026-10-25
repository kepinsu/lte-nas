package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestEPSAttachTypeFields(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name  string
		value ie.EPSAttachType
	}{
		{name: "EPS", value: ie.EPSAttachTypeEPS},
		{name: "combined EPS and IMSI", value: ie.EPSAttachTypeCombinedEPSAndIMSI},
		{name: "EPS RLOS", value: ie.EPSAttachTypeEPSRLOS},
		{name: "emergency", value: ie.EPSAttachTypeEPSEmergency},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fields := &ie.EPSAttachTypeFields{Type: tc.value}
			payload, err := fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(payload, []byte{byte(tc.value)}) {
				t.Fatalf("Marshal = %x, want %02x", payload, tc.value)
			}

			value, err := ie.NewPositionalEPSAttachTypeWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			if value.Kind != ie.NoTag {
				t.Fatalf("Kind = %#02x, want NoTag", value.Kind)
			}
			decoded, err := value.EPSAttachType()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("EPSAttachType = %+v, want %+v", decoded, fields)
			}
		})
	}
}

func TestEPSAttachTypeFields_RejectInvalid(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		payload []byte
	}{
		{name: "empty", payload: nil},
		{name: "spare bit", payload: []byte{0x11}},
		{name: "unused value", payload: []byte{0x04}},
		{name: "reserved value", payload: []byte{0x07}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fields := ie.EPSAttachTypeFields{Type: ie.EPSAttachTypeEPS}
			before := fields
			if err := fields.UnmarshalBinary(tc.payload); !errors.Is(err, ie.ErrInvalidValue) && !errors.Is(err, ie.ErrInvalidLength) {
				t.Fatalf("UnmarshalBinary = %v, want invalid value or length", err)
			}
			if fields != before {
				t.Fatalf("UnmarshalBinary modified fields: got %+v, want %+v", fields, before)
			}
		})
	}

	fields := &ie.EPSAttachTypeFields{Type: ie.EPSAttachTypeEPS}
	buffer := []byte{0xaa}
	if err := fields.MarshalTo(buffer[:0]); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo = %v, want io.ErrShortBuffer", err)
	}
	if buffer[0] != 0xaa {
		t.Fatalf("MarshalTo modified short buffer: got %x, want aa", buffer)
	}
}
