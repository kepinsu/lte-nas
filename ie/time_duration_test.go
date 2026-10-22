package ie_test

import (
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestTimeDurationFields(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		seconds uint32
		wantBP  [3]byte
	}{
		{
			name:    "zero duration",
			seconds: 0x000000,
			wantBP:  [3]byte{0x00, 0x00, 0x00},
		},
		{
			name:    "max 24-bit",
			seconds: 0xffffff,
			wantBP:  [3]byte{0xff, 0xff, 0xff},
		},
		{
			name:    "typical duration 30s",
			seconds: 0x001e, // 30 seconds
			wantBP:  [3]byte{0x00, 0x00, 0x1e},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fields := &ie.TimeDurationFields{Seconds: tt.seconds}

			// Marshal
			payload, err := fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(payload, tt.wantBP[:]) {
				t.Fatalf("Marshal = %x, want %x", payload, tt.wantBP)
			}

			// Unmarshal
			decoded, err := ie.ParseTimeDurationFields(payload)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("Unmarshal = %+v, want %+v", decoded, fields)
			}

			// Positional constructor
			value, err := ie.NewPositionalTimeDurationWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			rt, err := value.TimeDuration()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(rt, fields) {
				t.Fatalf("round-trip = %+v, want %+v", rt, fields)
			}
		})
	}
}

func TestTimeDurationFields_Invalid(t *testing.T) {
	t.Parallel()

	if _, err := ie.ParseTimeDurationFields([]byte{0x00, 0x00}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse error = %v, want ErrInvalidLength", err)
	}

	if _, err := ie.ParseTimeDurationFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse nil error = %v, want ErrInvalidLength", err)
	}

	fields := &ie.TimeDurationFields{Seconds: 0x123456}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo nil error = %v, want io.ErrShortBuffer", err)
	}

	var nilFields *ie.TimeDurationFields
	if err := nilFields.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("nil Validate error = %v, want ErrInvalidValue", err)
	}

	// Test 24-bit overflow validation
	badFields := &ie.TimeDurationFields{Seconds: 0x1000000}
	if err := badFields.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Validate 25-bit value = %v, want ErrInvalidValue", err)
	}
	if _, err := badFields.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Marshal 25-bit value = %v, want ErrInvalidValue", err)
	}
}

func TestTimeDurationFields_MarshalLen(t *testing.T) {
	t.Parallel()

	var nilFields *ie.TimeDurationFields
	if got := nilFields.MarshalLen(); got != 0 {
		t.Fatalf("MarshalLen(nil) = %d, want 0", got)
	}

	fields := &ie.TimeDurationFields{Seconds: 3600}
	if got := fields.MarshalLen(); got != 3 {
		t.Fatalf("MarshalLen = %d, want 3", got)
	}
}
