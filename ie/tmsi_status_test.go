package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestTMSIStatusFields_RoundTrip(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		status uint8
	}{
		{name: "status clear", status: 0x00},
		{name: "status set", status: 0x02},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fields := &ie.TMSIStatusFields{Status: tc.status}

			encoded, err := fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(encoded, []byte{tc.status}) {
				t.Fatalf("Marshal = %x, want %02x", encoded, tc.status)
			}

			value, err := ie.NewTMSIStatusWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}

			decoded, err := value.TMSIStatus()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("TMSIStatus = %+v, want %+v", decoded, fields)
			}
		})
	}
}

func TestTMSIStatusFields_PositionalRoundTrip(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		status uint8
	}{
		{name: "status clear", status: 0x00},
		{name: "status set", status: 0x02},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fields := &ie.TMSIStatusFields{Status: tc.status}

			value, err := ie.NewPositionalTMSIStatusWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}

			decoded, err := value.TMSIStatus()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("TMSIStatus = %+v, want %+v", decoded, fields)
			}
		})
	}
}

func TestTMSIStatusFields_RejectSpareBits(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		payload []byte
	}{
		{name: "bit 0 set", payload: []byte{0x01}},
		{name: "bit 2 set", payload: []byte{0x04}},
		{name: "bit 7 set", payload: []byte{0x80}},
		{name: "all spare bits set", payload: []byte{0xfe}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fields := ie.TMSIStatusFields{Status: 0x02}
			before := fields
			if err := fields.UnmarshalBinary(tc.payload); !errors.Is(err, ie.ErrInvalidValue) {
				t.Fatalf("UnmarshalBinary = %v, want ErrInvalidValue", err)
			}
			if fields != before {
				t.Fatalf("UnmarshalBinary modified fields: got %+v, want %+v", fields, before)
			}
		})
	}

	fields := &ie.TMSIStatusFields{Status: 0x02}
	buffer := []byte{0xaa}
	if err := fields.MarshalTo(buffer[:0]); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo = %v, want io.ErrShortBuffer", err)
	}
	if buffer[0] != 0xaa {
		t.Fatalf("MarshalTo modified short buffer: got %x, want aa", buffer)
	}
}

func TestTMSIStatusFields_NilFields(t *testing.T) {
	t.Parallel()
	var nilFields *ie.TMSIStatusFields
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

func TestTMSIStatusFields_MarshalTo(t *testing.T) {
	t.Parallel()
	fields := &ie.TMSIStatusFields{Status: 0x02}

	buffer := make([]byte, 10)
	if err := fields.MarshalTo(buffer); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buffer[:1], []byte{0x02}) {
		t.Fatalf("MarshalTo = %x, want 02", buffer[:1])
	}

	short := make([]byte, 0)
	if err := fields.MarshalTo(short); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo short buffer: error = %v, want io.ErrShortBuffer", err)
	}
}
