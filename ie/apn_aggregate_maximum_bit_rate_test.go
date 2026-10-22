package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestAPNAggregateMaximumBitRateFields(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		fields  *ie.APNAggregateMaximumBitRateFields
		payload []byte
	}{
		{
			name: "base rates",
			fields: &ie.APNAggregateMaximumBitRateFields{
				DownlinkBitsPerSecond: 64_000,
				UplinkBitsPerSecond:   8_640_000,
			},
			payload: []byte{0x40, 0xfe},
		},
		{
			name: "extended rates",
			fields: &ie.APNAggregateMaximumBitRateFields{
				DownlinkBitsPerSecond: 17_000_000,
				UplinkBitsPerSecond:   100_000_000,
			},
			payload: []byte{0xfe, 0xfe, 0x4b, 0x9e},
		},
		{
			name: "extended 2 rates",
			fields: &ie.APNAggregateMaximumBitRateFields{
				DownlinkBitsPerSecond: 300_000_000,
				UplinkBitsPerSecond:   512_000_000,
			},
			payload: []byte{0xfe, 0xfe, 0x66, 0xfa, 0x01, 0x01},
		},
		{
			name: "zero rates",
			fields: &ie.APNAggregateMaximumBitRateFields{
				DownlinkBitsPerSecond: 0,
				UplinkBitsPerSecond:   0,
			},
			payload: []byte{0xff, 0xff},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			payload, err := tc.fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(payload, tc.payload) {
				t.Fatalf("Marshal = %x, want %x", payload, tc.payload)
			}
			value, err := ie.NewAPNAggregateMaximumBitRateWithFields(tc.fields)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := value.APNAggregateMaximumBitRate()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, tc.fields) {
				t.Fatalf("APNAggregateMaximumBitRate = %+v, want %+v", decoded, tc.fields)
			}
		})
	}
}

func TestAPNAggregateMaximumBitRateFields_Invalid(t *testing.T) {
	t.Parallel()
	fields := &ie.APNAggregateMaximumBitRateFields{DownlinkBitsPerSecond: 16_500_000}
	if _, err := fields.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Marshal error = %v, want ErrInvalidValue", err)
	}
	if err := fields.MarshalTo(nil); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("MarshalTo error = %v, want ErrInvalidValue", err)
	}
	if _, err := ie.ParseAPNAggregateMaximumBitRateFields([]byte{0, 0}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseAPNAggregateMaximumBitRateFields error = %v, want ErrInvalidValue", err)
	}

	fields = &ie.APNAggregateMaximumBitRateFields{}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
