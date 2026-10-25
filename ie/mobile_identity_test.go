package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestMobileIdentityFields_IMSIRoundTrip(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		digits  string
		wantLen int
	}{
		{name: "odd digits (15)", digits: "123456789012345", wantLen: 8},
		{name: "even digits (14)", digits: "12345678901234", wantLen: 8},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fields := &ie.MobileIdentityFields{Type: ie.MobileIdentityTypeIMSI, Digits: tc.digits}

			encoded, err := fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if len(encoded) != tc.wantLen {
				t.Fatalf("Marshal = %x (len=%d), want len %d", encoded, len(encoded), tc.wantLen)
			}

			value, err := ie.NewMobileIdentityWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}

			decoded, err := value.MobileIdentity()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("MobileIdentity = %+v, want %+v", decoded, fields)
			}
		})
	}
}

func TestMobileIdentityFields_IMEIRoundTrip(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		digits  string
		wantLen int
	}{
		{name: "odd digits", digits: "352099001234567", wantLen: 8},
		{name: "even digits (14)", digits: "35209900123456", wantLen: 8},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fields := &ie.MobileIdentityFields{Type: ie.MobileIdentityTypeIMEI, Digits: tc.digits}

			encoded, err := fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if len(encoded) != tc.wantLen {
				t.Fatalf("Marshal = %x (len=%d), want len %d", encoded, len(encoded), tc.wantLen)
			}

			value, err := ie.NewPositionalMobileIdentityWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}

			decoded, err := value.MobileIdentity()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("MobileIdentity = %+v, want %+v", decoded, fields)
			}
		})
	}
}

func TestMobileIdentityFields_TMSIRoundTrip(t *testing.T) {
	t.Parallel()
	fields := &ie.MobileIdentityFields{Type: ie.MobileIdentityTypeTMSI, Digits: "000000000123456"}

	encoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	value, err := ie.NewMobileIdentityWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := value.MobileIdentity()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("MobileIdentity = %+v, want %+v", decoded, fields)
	}

	// Check the first byte encodes type 0x05 with odd indicator
	if encoded[0]&0x07 != 0x05 {
		t.Fatalf("First byte type field = %02x, want type 05", encoded[0])
	}
}

func TestMobileIdentityFields_RejectReservedType(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		payload []byte
	}{
		{name: "empty", payload: nil},
		{name: "GUTI type 06 (reserved)", payload: []byte{0xf6, 0x02, 0xf8}},
		{name: "unknown type 07", payload: []byte{0x31}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fields := ie.MobileIdentityFields{Type: ie.MobileIdentityTypeIMSI, Digits: "123"}
			before := fields
			if err := fields.UnmarshalBinary(tc.payload); !errors.Is(err, ie.ErrInvalidValue) && !errors.Is(err, ie.ErrInvalidLength) {
				t.Fatalf("UnmarshalBinary = %v, want invalid value or length", err)
			}
			if fields != before {
				t.Fatalf("UnmarshalBinary modified fields: got %+v, want %+v", fields, before)
			}
		})
	}

	fields := &ie.MobileIdentityFields{Type: ie.MobileIdentityTypeIMSI, Digits: "123"}
	buffer := []byte{0xaa}
	if err := fields.MarshalTo(buffer[:0]); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo = %v, want io.ErrShortBuffer", err)
	}
	if buffer[0] != 0xaa {
		t.Fatalf("MarshalTo modified short buffer: got %x, want aa", buffer)
	}
}

func TestMobileIdentityFields_NilFields(t *testing.T) {
	t.Parallel()
	var nilFields *ie.MobileIdentityFields
	if err := nilFields.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Validate nil: error = %v, want ErrInvalidValue", err)
	}
	_, err := nilFields.Marshal()
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Marshal nil: error = %v, want ErrInvalidValue", err)
	}
	if nilFields.MarshalLen() != 0 {
		t.Fatal("MarshalLen nil: expected 0")
	}
}

func TestMobileIdentityFields_MarshalTo(t *testing.T) {
	t.Parallel()
	fields := &ie.MobileIdentityFields{Type: ie.MobileIdentityTypeIMSI, Digits: "12345"}

	buffer := make([]byte, 10)
	if err := fields.MarshalTo(buffer); err != nil {
		t.Fatal(err)
	}
	want := []byte{0x19, 0x32, 0x54}
	if !bytes.Equal(buffer[:len(want)], want) {
		t.Fatalf("MarshalTo = %x, want %x", buffer[:len(want)], want)
	}

	short := make([]byte, 0)
	if err := fields.MarshalTo(short); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo short buffer: error = %v, want io.ErrShortBuffer", err)
	}
}
