package ie_test

import (
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestHashMMEFields(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		payload [8]byte
	}{
		{
			name:    "zero hash",
			payload: [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			name:    "all bits set",
			payload: [8]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		},
		{
			name:    "typical hash",
			payload: [8]byte{0x4a, 0xe3, 0xc7, 0xb1, 0x2d, 0xf8, 0x96, 0x50},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fields := &ie.HashMMEFields{Payload: tt.payload}

			// Marshal
			payload, err := fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(payload, tt.payload[:]) {
				t.Fatalf("Marshal = %x, want %x", payload, tt.payload)
			}

			// Unmarshal
			decoded, err := ie.ParseHashMMEFields(payload)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("Unmarshal = %+v, want %+v", decoded, fields)
			}

			// Positional constructor
			value, err := ie.NewPositionalHashMMEWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			rt, err := value.HashMME()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(rt, fields) {
				t.Fatalf("round-trip = %+v, want %+v", rt, fields)
			}

			// Tagged constructor
			tagged, err := ie.NewHashMMEWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			taggedRT, err := tagged.HashMME()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(taggedRT, fields) {
				t.Fatalf("tagged round-trip = %+v, want %+v", taggedRT, fields)
			}
		})
	}
}

func TestHashMMEFields_Invalid(t *testing.T) {
	t.Parallel()

	if _, err := ie.ParseHashMMEFields([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse error = %v, want ErrInvalidLength", err)
	}

	if _, err := ie.ParseHashMMEFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse nil error = %v, want ErrInvalidLength", err)
	}

	fields := &ie.HashMMEFields{Payload: [8]byte{0x4a, 0xe3, 0xc7, 0xb1, 0x2d, 0xf8, 0x96, 0x50}}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo nil error = %v, want io.ErrShortBuffer", err)
	}

	var nilFields *ie.HashMMEFields
	if err := nilFields.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("nil Validate error = %v, want ErrInvalidValue", err)
	}
}

func TestHashMMEFields_MarshalLen(t *testing.T) {
	t.Parallel()

	var nilFields *ie.HashMMEFields
	if got := nilFields.MarshalLen(); got != 0 {
		t.Fatalf("MarshalLen(nil) = %d, want 0", got)
	}

	fields := &ie.HashMMEFields{Payload: [8]byte{0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}
	if got := fields.MarshalLen(); got != 8 {
		t.Fatalf("MarshalLen = %d, want 8", got)
	}
}
