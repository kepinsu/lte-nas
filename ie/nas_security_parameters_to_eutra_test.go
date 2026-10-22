package ie_test

import (
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestNASSecurityParametersToEUTRAFields(t *testing.T) {
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
			name:    "typical security parameters",
			payload: [6]byte{0x82, 0xf4, 0xc3, 0xb1, 0x2d, 0x67},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fields := &ie.NASSecurityParametersToEUTRAFields{Payload: tt.payload}

			// Marshal
			payload, err := fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(payload, tt.payload[:]) {
				t.Fatalf("Marshal = %x, want %x", payload, tt.payload)
			}

			// Unmarshal
			decoded, err := ie.ParseNASSecurityParametersToEUTRAFields(payload)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("Unmarshal = %+v, want %+v", decoded, fields)
			}

			// Positional constructor
			value, err := ie.NewPositionalNASSecurityParametersToEUTRAWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			rt, err := value.NASSecurityParametersToEUTRA()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(rt, fields) {
				t.Fatalf("round-trip = %+v, want %+v", rt, fields)
			}

			// Tagged constructor
			tagged, err := ie.NewNASSecurityParametersToEUTRAWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			taggedRT, err := tagged.NASSecurityParametersToEUTRA()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(taggedRT, fields) {
				t.Fatalf("tagged round-trip = %+v, want %+v", taggedRT, fields)
			}
		})
	}
}

func TestNASSecurityParametersToEUTRAFields_Invalid(t *testing.T) {
	t.Parallel()

	if _, err := ie.ParseNASSecurityParametersToEUTRAFields([]byte{0x00, 0x00, 0x00, 0x00, 0x00}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse error = %v, want ErrInvalidLength", err)
	}

	if _, err := ie.ParseNASSecurityParametersToEUTRAFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse nil error = %v, want ErrInvalidLength", err)
	}

	fields := &ie.NASSecurityParametersToEUTRAFields{Payload: [6]byte{0x82, 0xf4, 0xc3, 0xb1, 0x2d, 0x67}}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo nil error = %v, want io.ErrShortBuffer", err)
	}

	var nilFields *ie.NASSecurityParametersToEUTRAFields
	if err := nilFields.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("nil Validate error = %v, want ErrInvalidValue", err)
	}
}

func TestNASSecurityParametersToEUTRAFields_MarshalLen(t *testing.T) {
	t.Parallel()

	var nilFields *ie.NASSecurityParametersToEUTRAFields
	if got := nilFields.MarshalLen(); got != 0 {
		t.Fatalf("MarshalLen(nil) = %d, want 0", got)
	}

	fields := &ie.NASSecurityParametersToEUTRAFields{Payload: [6]byte{0x42, 0x00, 0x00, 0x00, 0x00, 0x00}}
	if got := fields.MarshalLen(); got != 6 {
		t.Fatalf("MarshalLen = %d, want 6", got)
	}
}
