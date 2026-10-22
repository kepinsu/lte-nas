package ie_test

import (
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestUEAdditionalSecurityCapabilityFields(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		payload [4]byte
	}{
		{
			name:    "zero capability",
			payload: [4]byte{0x00, 0x00, 0x00, 0x00},
		},
		{
			name:    "all bits set",
			payload: [4]byte{0xff, 0xff, 0xff, 0xff},
		},
		{
			name:    "EA support only",
			payload: [4]byte{0x0f, 0x00, 0x00, 0x00},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fields := &ie.UEAdditionalSecurityCapabilityFields{Payload: tt.payload}

			// Marshal
			payload, err := fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(payload, tt.payload[:]) {
				t.Fatalf("Marshal = %x, want %x", payload, tt.payload)
			}

			// Unmarshal
			decoded, err := ie.ParseUEAdditionalSecurityCapabilityFields(payload)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("Unmarshal = %+v, want %+v", decoded, fields)
			}

			// Positional constructor
			value, err := ie.NewPositionalUEAdditionalSecurityCapabilityWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			rt, err := value.UEAdditionalSecurityCapability()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(rt, fields) {
				t.Fatalf("round-trip = %+v, want %+v", rt, fields)
			}

			// Tagged constructor
			tagged, err := ie.NewUEAdditionalSecurityCapabilityWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			taggedRT, err := tagged.UEAdditionalSecurityCapability()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(taggedRT, fields) {
				t.Fatalf("tagged round-trip = %+v, want %+v", taggedRT, fields)
			}
		})
	}
}

func TestUEAdditionalSecurityCapabilityFields_Invalid(t *testing.T) {
	t.Parallel()

	if _, err := ie.ParseUEAdditionalSecurityCapabilityFields([]byte{0x00, 0x00, 0x00}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse error = %v, want ErrInvalidLength", err)
	}

	if _, err := ie.ParseUEAdditionalSecurityCapabilityFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse nil error = %v, want ErrInvalidLength", err)
	}

	fields := &ie.UEAdditionalSecurityCapabilityFields{Payload: [4]byte{0x0f, 0xff, 0x00, 0x01}}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo nil error = %v, want io.ErrShortBuffer", err)
	}

	var nilFields *ie.UEAdditionalSecurityCapabilityFields
	if err := nilFields.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("nil Validate error = %v, want ErrInvalidValue", err)
	}
}

func TestUEAdditionalSecurityCapabilityFields_MarshalLen(t *testing.T) {
	t.Parallel()

	var nilFields *ie.UEAdditionalSecurityCapabilityFields
	if got := nilFields.MarshalLen(); got != 0 {
		t.Fatalf("MarshalLen(nil) = %d, want 0", got)
	}

	fields := &ie.UEAdditionalSecurityCapabilityFields{Payload: [4]byte{0x42, 0x00, 0x00, 0x00}}
	if got := fields.MarshalLen(); got != 4 {
		t.Fatalf("MarshalLen = %d, want 4", got)
	}
}
