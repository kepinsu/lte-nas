package ie_test

import (
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestServingPLMNRateControlFields(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		value  uint16
		wantBP [2]byte
	}{
		{
			name:   "zero period",
			value:  0x0000,
			wantBP: [2]byte{0x00, 0x00},
		},
		{
			name:   "max period",
			value:  0xffff,
			wantBP: [2]byte{0xff, 0xff},
		},
		{
			name:   "typical period",
			value:  0x1f40, // 8000 ms
			wantBP: [2]byte{0x1f, 0x40},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fields := &ie.ServingPLMNRateControlFields{RateControlPeriodMS: tt.value}

			// Marshal
			payload, err := fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(payload, tt.wantBP[:]) {
				t.Fatalf("Marshal = %x, want %x", payload, tt.wantBP)
			}

			// Unmarshal
			decoded, err := ie.ParseServingPLMNRateControlFields(payload)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("Unmarshal = %+v, want %+v", decoded, fields)
			}

			// Positional constructor
			value, err := ie.NewPositionalServingPLMNRateControlWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			rt, err := value.ServingPLMNRateControl()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(rt, fields) {
				t.Fatalf("round-trip = %+v, want %+v", rt, fields)
			}

			// Tagged constructor
			tagged, err := ie.NewServingPLMNRateControlWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			taggedRT, err := tagged.ServingPLMNRateControl()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(taggedRT, fields) {
				t.Fatalf("tagged round-trip = %+v, want %+v", taggedRT, fields)
			}
		})
	}
}

func TestServingPLMNRateControlFields_Invalid(t *testing.T) {
	t.Parallel()

	if _, err := ie.ParseServingPLMNRateControlFields([]byte{0x00}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse error = %v, want ErrInvalidLength", err)
	}

	if _, err := ie.ParseServingPLMNRateControlFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse nil error = %v, want ErrInvalidLength", err)
	}

	fields := &ie.ServingPLMNRateControlFields{RateControlPeriodMS: 0x1234}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo nil error = %v, want io.ErrShortBuffer", err)
	}

	var nilFields *ie.ServingPLMNRateControlFields
	if err := nilFields.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("nil Validate error = %v, want ErrInvalidValue", err)
	}
}

func TestServingPLMNRateControlFields_MarshalLen(t *testing.T) {
	t.Parallel()

	var nilFields *ie.ServingPLMNRateControlFields
	if got := nilFields.MarshalLen(); got != 0 {
		t.Fatalf("MarshalLen(nil) = %d, want 0", got)
	}

	fields := &ie.ServingPLMNRateControlFields{RateControlPeriodMS: 0x42}
	if got := fields.MarshalLen(); got != 2 {
		t.Fatalf("MarshalLen = %d, want 2", got)
	}
}
