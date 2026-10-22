package ie_test

import (
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestHeaderCompressionConfigurationStatusFields(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		payload [2]byte
	}{
		{
			name:    "zero bitmap",
			payload: [2]byte{0x00, 0x00},
		},
		{
			name:    "all bits set",
			payload: [2]byte{0xff, 0xff},
		},
		{
			name:    "mixed",
			payload: [2]byte{0x5a, 0xa5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fields := &ie.HeaderCompressionConfigurationStatusFields{Payload: tt.payload}

			// Marshal
			payload, err := fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(payload, tt.payload[:]) {
				t.Fatalf("Marshal = %x, want %x", payload, tt.payload)
			}

			// Unmarshal
			decoded, err := ie.ParseHeaderCompressionConfigurationStatusFields(payload)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("Unmarshal = %+v, want %+v", decoded, fields)
			}

			// Positional constructor
			value, err := ie.NewPositionalHeaderCompressionConfigurationStatusWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			rt, err := value.HeaderCompressionConfigurationStatus()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(rt, fields) {
				t.Fatalf("round-trip = %+v, want %+v", rt, fields)
			}

			// Tagged constructor
			tagged, err := ie.NewHeaderCompressionConfigurationStatusWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			taggedRT, err := tagged.HeaderCompressionConfigurationStatus()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(taggedRT, fields) {
				t.Fatalf("tagged round-trip = %+v, want %+v", taggedRT, fields)
			}
		})
	}
}

func TestHeaderCompressionConfigurationStatusFields_Invalid(t *testing.T) {
	t.Parallel()

	if _, err := ie.ParseHeaderCompressionConfigurationStatusFields([]byte{0x00}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse error = %v, want ErrInvalidLength", err)
	}

	if _, err := ie.ParseHeaderCompressionConfigurationStatusFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse nil error = %v, want ErrInvalidLength", err)
	}

	fields := &ie.HeaderCompressionConfigurationStatusFields{Payload: [2]byte{0x12, 0x34}}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo nil error = %v, want io.ErrShortBuffer", err)
	}

	var nilFields *ie.HeaderCompressionConfigurationStatusFields
	if err := nilFields.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("nil Validate error = %v, want ErrInvalidValue", err)
	}
}

func TestHeaderCompressionConfigurationStatusFields_MarshalLen(t *testing.T) {
	t.Parallel()

	var nilFields *ie.HeaderCompressionConfigurationStatusFields
	if got := nilFields.MarshalLen(); got != 0 {
		t.Fatalf("MarshalLen(nil) = %d, want 0", got)
	}

	fields := &ie.HeaderCompressionConfigurationStatusFields{Payload: [2]byte{0x42, 0x00}}
	if got := fields.MarshalLen(); got != 2 {
		t.Fatalf("MarshalLen = %d, want 2", got)
	}
}
