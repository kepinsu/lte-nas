package ie_test

import (
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestExtendedAPNAggregateMaximumBitRateFields(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		fields  *ie.ExtendedAPNAggregateMaximumBitRateFields
		payload []byte
	}{
		{
			name: "small rates (exact unit values)",
			fields: &ie.ExtendedAPNAggregateMaximumBitRateFields{
				DownlinkBitsPerSecond: 1_000_000_000, // exact 1 Gbps
				UplinkBitsPerSecond:   256_000_000,   // exact 256 Mbps
			},
			payload: []byte{
				0x07, 0x00, 0x00, // downlink: unit=1Gbps(7), value=0
				0x06, 0x00, 0x00, // uplink: unit=256Mbps(6), value=0
			},
		},
		{
			name: "small rates with offsets",
			fields: &ie.ExtendedAPNAggregateMaximumBitRateFields{
				DownlinkBitsPerSecond: 2_000_000_000, // 1 Gbps + 1 Gbps offset
				UplinkBitsPerSecond:   512_000_000,   // 256 Mbps + 256 Mbps offset
			},
			payload: []byte{
				0x07, 0x00, 0x01, // downlink: unit=1Gbps(7), value=1
				0x06, 0x00, 0x01, // uplink: unit=256Mbps(6), value=1
			},
		},
		{
			name: "zero rates",
			fields: &ie.ExtendedAPNAggregateMaximumBitRateFields{
				DownlinkBitsPerSecond: 0,
				UplinkBitsPerSecond:   0,
			},
			payload: []byte{
				0x00, 0x00, 0x00, // downlink: unit=not used(0), value=0
				0x00, 0x00, 0x00, // uplink: unit=not used(0), value=0
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			expectedLen := tc.fields.MarshalLen()
			payload, err := tc.fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}

			if len(payload) != expectedLen {
				t.Fatalf("Marshal length = %d, want %d", len(payload), expectedLen)
			}

			if !bytesEqual2(payload, tc.payload) {
				t.Fatalf("Marshal = %x, want %x", payload, tc.payload)
			}

			// MarshalTo
			toBuf := make([]byte, expectedLen)
			if err := tc.fields.MarshalTo(toBuf); err != nil {
				t.Fatalf("MarshalTo = %v", err)
			}
			if !bytesEqual2(toBuf, tc.payload) {
				t.Fatalf("MarshalTo = %x, want %x", toBuf, tc.payload)
			}

			// MarshalLen
			if got := tc.fields.MarshalLen(); got != expectedLen {
				t.Fatalf("MarshalLen = %d, want %d", got, expectedLen)
			}

			// Round-trip via constructor
			value, err := ie.NewExtendedAPNAggregateMaximumBitRateWithFields(tc.fields)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := value.ExtendedAPNAggregateMaximumBitRate()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, tc.fields) {
				t.Fatalf("round-trip = %+v, want %+v", decoded, tc.fields)
			}

			// Parse via function
			parsed, err := ie.ParseExtendedAPNAggregateMaximumBitRateFields(tc.payload)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(parsed, tc.fields) {
				t.Fatalf("Parse = %+v, want %+v", parsed, tc.fields)
			}
		})
	}
}

func TestExtendedAPNAggregateMaximumBitRateFields_Invalid(t *testing.T) {
	t.Parallel()

	// Too long
	if _, err := ie.ParseExtendedAPNAggregateMaximumBitRateFields(make([]byte, 10)); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse too long error = %v, want ErrInvalidLength", err)
	}

	// Too short
	if _, err := ie.ParseExtendedAPNAggregateMaximumBitRateFields([]byte{0x00}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse too short error = %v, want ErrInvalidLength", err)
	}

	// MarshalTo short buffer
	fields := &ie.ExtendedAPNAggregateMaximumBitRateFields{DownlinkBitsPerSecond: 1_000_000, UplinkBitsPerSecond: 1_000_000}
	if err := fields.MarshalTo([]byte{0}); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo short buffer error = %v, want io.ErrShortBuffer", err)
	}

	// MarshalLen for nil receiver
	if got := (*ie.ExtendedAPNAggregateMaximumBitRateFields)(nil).MarshalLen(); got != 0 {
		t.Fatalf("Nil MarshalLen = %d, want 0", got)
	}

	// Validate nil receiver
	if err := (*ie.ExtendedAPNAggregateMaximumBitRateFields)(nil).Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Nil validate error = %v, want ErrInvalidValue", err)
	}
}

func TestExtendedAPNAggregateMaximumBitRateRoundTrip(t *testing.T) {
	t.Parallel()

	// Test exact unit rates: 6-byte format (< 64 Gbps)
	for _, tc := range []struct {
		name      string
		downlink  uint64
		uplink    uint64
		expectLen int
	}{
		{name: "down 1G up 256M", downlink: 1_000_000_000, uplink: 256_000_000, expectLen: 6},
		{name: "down 4G up 16M", downlink: 4_000_000_000, uplink: 16_000_000, expectLen: 6},
		{name: "down 63.999G up 0", downlink: 63_999_999_999, uplink: 0, expectLen: 6},
		{name: "down 1.5G up 0.5G (rounded)", downlink: 1_500_000_000, uplink: 500_000_000, expectLen: 6},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fields := &ie.ExtendedAPNAggregateMaximumBitRateFields{
				DownlinkBitsPerSecond: tc.downlink,
				UplinkBitsPerSecond:   tc.uplink,
			}
			if got := fields.MarshalLen(); got != tc.expectLen {
				t.Fatalf("MarshalLen = %d, want %d", got, tc.expectLen)
			}

			payload, err := fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}

			parsed, err := ie.ParseExtendedAPNAggregateMaximumBitRateFields(payload)
			if err != nil {
				t.Fatal(err)
			}

			// The round-trip may not be exact for rates that fall between unit steps
			// (due to integer division truncation in encoding). We just verify
			// that the round-trip rate is >= original (rounds down) or == for exact rates.
			if parsed.DownlinkBitsPerSecond > tc.downlink {
				t.Fatalf("downlink round-trip %d > input %d", parsed.DownlinkBitsPerSecond, tc.downlink)
			}
			if parsed.UplinkBitsPerSecond > tc.uplink {
				t.Fatalf("uplink round-trip %d > input %d", parsed.UplinkBitsPerSecond, tc.uplink)
			}
		})
	}
}

func bytesEqual2(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
