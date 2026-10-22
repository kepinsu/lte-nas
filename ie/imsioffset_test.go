package ie_test

import (
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestImsioffsetFields(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		fields  *ie.ImsioffsetFields
		wantErr bool
	}{
		{
			name:    "round-trip zero",
			fields:  &ie.ImsioffsetFields{IMSIOffset: 0x0000},
			wantErr: false,
		},
		{
			name:    "round-trip max",
			fields:  &ie.ImsioffsetFields{IMSIOffset: 0xffff},
			wantErr: false,
		},
		{
			name:    "round-trip value",
			fields:  &ie.ImsioffsetFields{IMSIOffset: 0x1234},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Marshal
			payload, err := tt.fields.Marshal()
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if uint16(payload[0])<<8|uint16(payload[1]) != tt.fields.IMSIOffset {
				t.Fatalf("Marshal payload = %x, want offset 0x%04x", payload, tt.fields.IMSIOffset)
			}

			// Unmarshal
			decoded, err := ie.ParseIMSIOffsetFields(payload)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, tt.fields) {
				t.Fatalf("Unmarshal = %+v, want %+v", decoded, tt.fields)
			}

			// Positional constructor
			value, err := ie.NewPositionalIMSIOffsetWithFields(tt.fields)
			if err != nil {
				t.Fatal(err)
			}
			roundTrip, err := value.IMSIOffset()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(roundTrip, tt.fields) {
				t.Fatalf("round-trip = %+v, want %+v", roundTrip, tt.fields)
			}

			// Tagged constructor
			tagged, err := ie.NewIMSIOffsetWithFields(tt.fields)
			if err != nil {
				t.Fatal(err)
			}
			taggedRT, err := tagged.IMSIOffset()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(taggedRT, tt.fields) {
				t.Fatalf("tagged round-trip = %+v, want %+v", taggedRT, tt.fields)
			}
		})
	}
}

func TestImsioffsetFields_Invalid(t *testing.T) {
	t.Parallel()

	if _, err := ie.ParseIMSIOffsetFields([]byte{0x00}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseIMSIOffsetFields error = %v, want ErrInvalidLength", err)
	}

	if _, err := ie.ParseIMSIOffsetFields([]byte{0x00, 0x00, 0x00}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseIMSIOffsetFields 3-byte error = %v, want ErrInvalidLength", err)
	}

	if _, err := ie.ParseIMSIOffsetFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseIMSIOffsetFields nil error = %v, want ErrInvalidLength", err)
	}

	fields := &ie.ImsioffsetFields{IMSIOffset: 0x1234}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo nil error = %v, want io.ErrShortBuffer", err)
	}

	if err := fields.MarshalTo([]byte{0}); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo short buffer error = %v, want io.ErrShortBuffer", err)
	}

	var nilFields *ie.ImsioffsetFields
	if err := nilFields.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("nil Validate error = %v, want ErrInvalidValue", err)
	}

	if _, err := nilFields.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("nil Marshal error = %v, want ErrInvalidValue", err)
	}
}

func TestImsioffsetFields_MarshalLen(t *testing.T) {
	t.Parallel()

	var nilFields *ie.ImsioffsetFields
	if got := nilFields.MarshalLen(); got != 0 {
		t.Fatalf("MarshalLen(nil) = %d, want 0", got)
	}

	fields := &ie.ImsioffsetFields{IMSIOffset: 0x42}
	if got := fields.MarshalLen(); got != 2 {
		t.Fatalf("MarshalLen = %d, want 2", got)
	}
}
