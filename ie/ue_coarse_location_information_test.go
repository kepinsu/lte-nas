package ie_test

import (
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestUECoarseLocationInformationFields(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		payload [6]byte
	}{
		{
			name:    "zero payload",
			payload: [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			name:    "all bits set",
			payload: [6]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		},
		{
			name:    "typical location blob",
			payload: [6]byte{0x0a, 0xb2, 0xc3, 0xd4, 0xe5, 0xf6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fields := &ie.UECoarseLocationInformationFields{Payload: tt.payload}

			// Marshal
			payload, err := fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(payload, tt.payload[:]) {
				t.Fatalf("Marshal = %x, want %x", payload, tt.payload)
			}

			// Unmarshal
			decoded, err := ie.ParseUECoarseLocationInformationFields(payload)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("Unmarshal = %+v, want %+v", decoded, fields)
			}

			// Positional constructor
			value, err := ie.NewPositionalUECoarseLocationInformationWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			rt, err := value.UECoarseLocationInformation()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(rt, fields) {
				t.Fatalf("round-trip = %+v, want %+v", rt, fields)
			}

			// Tagged constructor
			tagged, err := ie.NewUECoarseLocationInformationWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			taggedRT, err := tagged.UECoarseLocationInformation()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(taggedRT, fields) {
				t.Fatalf("tagged round-trip = %+v, want %+v", taggedRT, fields)
			}
		})
	}
}

func TestUECoarseLocationInformationFields_Invalid(t *testing.T) {
	t.Parallel()

	if _, err := ie.ParseUECoarseLocationInformationFields([]byte{0x00, 0x00, 0x00, 0x00, 0x00}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse error = %v, want ErrInvalidLength", err)
	}

	if _, err := ie.ParseUECoarseLocationInformationFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse nil error = %v, want ErrInvalidLength", err)
	}

	fields := &ie.UECoarseLocationInformationFields{Payload: [6]byte{0x0a, 0xb2, 0xc3, 0xd4, 0xe5, 0xf6}}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo nil error = %v, want io.ErrShortBuffer", err)
	}

	var nilFields *ie.UECoarseLocationInformationFields
	if err := nilFields.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("nil Validate error = %v, want ErrInvalidValue", err)
	}
}

func TestUECoarseLocationInformationFields_MarshalLen(t *testing.T) {
	t.Parallel()

	var nilFields *ie.UECoarseLocationInformationFields
	if got := nilFields.MarshalLen(); got != 0 {
		t.Fatalf("MarshalLen(nil) = %d, want 0", got)
	}

	fields := &ie.UECoarseLocationInformationFields{Payload: [6]byte{0x42, 0x00, 0x00, 0x00, 0x00, 0x00}}
	if got := fields.MarshalLen(); got != 6 {
		t.Fatalf("MarshalLen = %d, want 6", got)
	}
}
