package ie_test

import (
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestMobileStationClassmark2Fields(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		payload [3]byte
	}{
		{
			name:    "zero payload",
			payload: [3]byte{0x00, 0x00, 0x00},
		},
		{
			name:    "all bits set",
			payload: [3]byte{0xff, 0xff, 0xff},
		},
		{
			name:    "classmark 2 typical",
			payload: [3]byte{0x4e, 0xa5, 0x02},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fields := &ie.MobileStationClassmark2Fields{Payload: tt.payload}

			// Marshal
			payload, err := fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(payload, tt.payload[:]) {
				t.Fatalf("Marshal = %x, want %x", payload, tt.payload)
			}

			// Unmarshal
			decoded, err := ie.ParseMobileStationClassmark2Fields(payload)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("Unmarshal = %+v, want %+v", decoded, fields)
			}

			// Positional constructor
			value, err := ie.NewPositionalMobileStationClassmark2WithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			rt, err := value.MobileStationClassmark2()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(rt, fields) {
				t.Fatalf("round-trip = %+v, want %+v", rt, fields)
			}

			// Tagged constructor
			tagged, err := ie.NewMobileStationClassmark2WithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			taggedRT, err := tagged.MobileStationClassmark2()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(taggedRT, fields) {
				t.Fatalf("tagged round-trip = %+v, want %+v", taggedRT, fields)
			}
		})
	}
}

func TestMobileStationClassmark2Fields_Invalid(t *testing.T) {
	t.Parallel()

	if _, err := ie.ParseMobileStationClassmark2Fields([]byte{0x00, 0x00}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse error = %v, want ErrInvalidLength", err)
	}

	if _, err := ie.ParseMobileStationClassmark2Fields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse nil error = %v, want ErrInvalidLength", err)
	}

	fields := &ie.MobileStationClassmark2Fields{Payload: [3]byte{0x4e, 0xa5, 0x02}}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo nil error = %v, want io.ErrShortBuffer", err)
	}

	var nilFields *ie.MobileStationClassmark2Fields
	if err := nilFields.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("nil Validate error = %v, want ErrInvalidValue", err)
	}
}

func TestMobileStationClassmark2Fields_MarshalLen(t *testing.T) {
	t.Parallel()

	var nilFields *ie.MobileStationClassmark2Fields
	if got := nilFields.MarshalLen(); got != 0 {
		t.Fatalf("MarshalLen(nil) = %d, want 0", got)
	}

	fields := &ie.MobileStationClassmark2Fields{Payload: [3]byte{0x42, 0x00, 0x00}}
	if got := fields.MarshalLen(); got != 3 {
		t.Fatalf("MarshalLen = %d, want 3", got)
	}
}
