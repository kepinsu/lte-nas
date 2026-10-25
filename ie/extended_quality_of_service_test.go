package ie_test

import (
	"encoding/binary"
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestExtendedQualityOfServiceFields(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		fields  *ie.ExtendedQualityOfServiceFields
		payload []byte
	}{
		{
			name: "all fields set",
			fields: &ie.ExtendedQualityOfServiceFields{
				QCI:           0x45,
				TrafficClass:  0x1234,
				PriorityLevel: 7,
				ARPPreEmption: 0xc0,
				GBR:           100_000, // fits in 24 bits
				MBR:           300_000, // upper 16 bits stored
				TOSField:      0x1e5c,  // no wire bytes — not round-tripped
			},
			payload: []byte{
				0x45,       // QCI
				0x12, 0x34, // TrafficClass big-endian
				0x07,             // PriorityLevel
				0xc0,             // ARPPreEmption
				0x01, 0x86, 0xa0, // GBR 24-bit big-endian = 100000
				0x04, 0x93, // MBR upper 16 bits of 300000 = 0x493
			},
		},
		{
			name: "minimal fields",
			fields: &ie.ExtendedQualityOfServiceFields{
				QCI:           0x00,
				TrafficClass:  0x0000,
				PriorityLevel: 0,
				ARPPreEmption: 0x00,
				GBR:           0x00,
				MBR:           0x00,
			},
			payload: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			name: "maximum field values",
			fields: &ie.ExtendedQualityOfServiceFields{
				QCI:           0xff,
				TrafficClass:  0xffff,
				PriorityLevel: 0x7f,
				ARPPreEmption: 0xff,
				GBR:           0xffffff, // max 24-bit = 16777215
				MBR:           0xffff00, // max upper-16 bits = 16776960
			},
			payload: []byte{
				0xff,
				0xff, 0xff,
				0x7f,
				0xff,
				0xff, 0xff, 0xff,
				0xff, 0xff,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			payload, err := tc.fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if len(payload) != ie.ExtendedQualityOfServiceLength {
				t.Fatalf("Marshal length = %d, want %d", len(payload), ie.ExtendedQualityOfServiceLength)
			}
			if !bytesEqual3(payload, tc.payload) {
				t.Fatalf("Marshal = %x, want %x", payload, tc.payload)
			}

			// MarshalTo
			toBuf := make([]byte, ie.ExtendedQualityOfServiceLength)
			if err := tc.fields.MarshalTo(toBuf); err != nil {
				t.Fatalf("MarshalTo = %v", err)
			}
			if !bytesEqual3(toBuf, tc.payload) {
				t.Fatalf("MarshalTo = %x, want %x", toBuf, tc.payload)
			}

			// MarshalLen
			if got := tc.fields.MarshalLen(); got != ie.ExtendedQualityOfServiceLength {
				t.Fatalf("MarshalLen = %d, want %d", got, ie.ExtendedQualityOfServiceLength)
			}

			// Round-trip via constructor (TOSField not on wire)
			ieObj, err := ie.NewExtendedQualityOfServiceWithFields(tc.fields)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := ieObj.ExtendedQualityOfService()
			if err != nil {
				t.Fatal(err)
			}
			// Verify key fields, ignoring TOSField and MBR lower byte (not on wire)
			if decoded.QCI != tc.fields.QCI ||
				decoded.TrafficClass != tc.fields.TrafficClass ||
				decoded.PriorityLevel != tc.fields.PriorityLevel ||
				decoded.ARPPreEmption != tc.fields.ARPPreEmption ||
				decoded.GBR != tc.fields.GBR {
				t.Fatalf("round-trip key fields mismatch: got %+v, want %+v", decoded, tc.fields)
			}
			// MBR: only upper 16 bits survive
			gotMBR := decoded.MBR & (uint32(0xffff) << 8)
			wantMBR := tc.fields.MBR & (uint32(0xffff) << 8)
			if gotMBR != wantMBR {
				t.Fatalf("round-trip MBR upper mismatch: got 0x%06x, want 0x%06x", gotMBR, wantMBR)
			}

			// Positional round-trip
			posIE, err := ie.NewPositionalExtendedQualityOfServiceWithFields(tc.fields)
			if err != nil {
				t.Fatal(err)
			}
			posDecoded, err := posIE.ExtendedQualityOfService()
			if err != nil {
				t.Fatal(err)
			}
			if posDecoded.QCI != tc.fields.QCI ||
				posDecoded.TrafficClass != tc.fields.TrafficClass ||
				posDecoded.GBR != tc.fields.GBR {
				t.Fatalf("positional round-trip mismatch: got %+v, want %+v", posDecoded, tc.fields)
			}

			// Parse directly
			parsed, err := ie.ParseExtendedQualityOfServiceFields(tc.payload)
			if err != nil {
				t.Fatal(err)
			}
			if parsed.QCI != tc.fields.QCI ||
				parsed.TrafficClass != tc.fields.TrafficClass ||
				parsed.GBR != tc.fields.GBR {
				t.Fatalf("Parse mismatch: got %+v, want %+v", parsed, tc.fields)
			}
		})
	}
}

func TestExtendedQualityOfServiceFields_Invalid(t *testing.T) {
	t.Parallel()

	// Validate nil
	if err := (*ie.ExtendedQualityOfServiceFields)(nil).Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Validate(nil) = %v, want ErrInvalidValue", err)
	}

	// Marshal nil
	var nilFields *ie.ExtendedQualityOfServiceFields
	if _, err := nilFields.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Marshal(nil) = %v, want ErrInvalidValue", err)
	}

	// MarshalTo nil receiver
	if _, err := (*ie.ExtendedQualityOfServiceFields)(nil).Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Marshal(nil) = %v, want ErrInvalidValue", err)
	}

	// MarshalTo with short buffer
	fields := &ie.ExtendedQualityOfServiceFields{QCI: 0x45}
	if err := fields.MarshalTo(make([]byte, 9)); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo(short buffer) = %v, want io.ErrShortBuffer", err)
	}

	// Unmarshal nil
	if err := (*ie.ExtendedQualityOfServiceFields)(nil).UnmarshalBinary([]byte{0}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("UnmarshalBinary(nil) = %v, want ErrInvalidValue", err)
	}

	// Unmarshal wrong length (need exactly 10 bytes)
	fields2 := &ie.ExtendedQualityOfServiceFields{}
	if err := fields2.UnmarshalBinary([]byte{0, 1}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("UnmarshalBinary(wrong length) = %v, want ErrInvalidLength", err)
	}

	// Parse with wrong length
	if _, err := ie.ParseExtendedQualityOfServiceFields([]byte{0, 1}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse(wrong length) = %v, want ErrInvalidLength", err)
	}

	// MarshalLen on nil
	if got := (*ie.ExtendedQualityOfServiceFields)(nil).MarshalLen(); got != 0 {
		t.Fatalf("MarshalLen(nil) = %d, want 0", got)
	}

	// Test constructor with short payload
	ieObj := ie.NewExtendedQualityOfService([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8})
	if _, err := ieObj.ExtendedQualityOfService(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("IE.ExtendedQualityOfService on short payload = %v, want ErrInvalidLength", err)
	}

	posIE := ie.NewPositionalExtendedQualityOfService([]byte{0, 1, 2, 3, 4, 5})
	if _, err := posIE.ExtendedQualityOfService(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("positional IE on short payload = %v, want ErrInvalidLength", err)
	}

	// Test IE accessor with wrong IE type
	wrongIE := ie.NewQualityOfService([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 0})
	if _, err := wrongIE.ExtendedQualityOfService(); !errors.Is(err, ie.ErrUnexpectedIE) {
		t.Fatalf("IE.ExtendedQualityOfService on wrong IE type = %v, want ErrUnexpectedIE", err)
	}
}

func TestExtendedQualityOfServiceRoundTrip(t *testing.T) {
	t.Parallel()

	// Full round-trip with meaningful values
	fields := &ie.ExtendedQualityOfServiceFields{
		QCI:           0x46,
		TrafficClass:  0x2a50,
		PriorityLevel: 4,
		ARPPreEmption: 0xa5,
		GBR:           1_234_567 & 0xffffff,
		MBR:           8_765_432 & 0xffffff,
		TOSField:      0x1e5c,
	}

	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	// Verify length is exactly 10
	if len(payload) != 10 {
		t.Fatalf("payload length = %d, want 10", len(payload))
	}

	// Verify QCI in byte 0
	if payload[0] != 0x46 {
		t.Fatalf("QCI byte = 0x%02x, want 0x46", payload[0])
	}

	// Verify TrafficClass (big-endian in bytes 1-2)
	tcBytes := binary.BigEndian.Uint16(payload[1:3])
	if tcBytes != fields.TrafficClass {
		t.Fatalf("TrafficClass bytes = 0x%04x, want 0x%04x", tcBytes, fields.TrafficClass)
	}

	// Verify GBR (bytes 5-7 as 24-bit big-endian)
	gbrDecoded := uint32(payload[5])<<16 | uint32(payload[6])<<8 | uint32(payload[7])
	if gbrDecoded != fields.GBR {
		t.Fatalf("GBR decoded = 0x%06x, want 0x%06x", gbrDecoded, fields.GBR)
	}

	// Verify MBR upper 16 bits (bytes 8-9)
	mbrUpper := uint32(payload[8])<<16 | uint32(payload[9])<<8
	wantMBRUpper := fields.MBR & (uint32(0xffff) << 8)
	if mbrUpper != wantMBRUpper {
		t.Fatalf("MBR upper bytes = 0x%06x, want 0x%06x", mbrUpper, wantMBRUpper)
	}

	// Round-trip via Parse — MBR lower byte is lost (only 16 bits on wire)
	parsed, err := ie.ParseExtendedQualityOfServiceFields(payload)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.QCI != fields.QCI ||
		parsed.TrafficClass != fields.TrafficClass ||
		parsed.PriorityLevel != fields.PriorityLevel ||
		parsed.ARPPreEmption != fields.ARPPreEmption ||
		parsed.GBR != fields.GBR {
		t.Fatalf("Parse round-trip mismatch for non-MBR fields")
	}
	if (parsed.MBR & 0xffff00) != wantMBRUpper {
		t.Fatalf("Parse round-trip MBR upper mismatch: got 0x%06x, want 0x%06x", parsed.MBR, wantMBRUpper)
	}

	// Round-trip via IE — same as above, MBR lower byte lost
	ieObj, err := ie.NewExtendedQualityOfServiceWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	if !bytesEqual3(ieObj.Payload, payload) {
		t.Fatalf("IE payload = %x, want %x", ieObj.Payload, payload)
	}
	recovered, err := ieObj.ExtendedQualityOfService()
	if err != nil {
		t.Fatal(err)
	}
	if recovered.QCI != fields.QCI ||
		recovered.TrafficClass != fields.TrafficClass ||
		recovered.GBR != fields.GBR {
		t.Fatalf("IE round-trip mismatch")
	}
}

func bytesEqual3(a, b []byte) bool {
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
