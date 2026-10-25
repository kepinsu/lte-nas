package ie_test

import (
	"bytes"
	"encoding"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

var (
	_ encoding.BinaryMarshaler   = (*ie.IE)(nil)
	_ encoding.BinaryUnmarshaler = (*ie.IE)(nil)
)

func TestIE_Marshal(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		ie   *ie.IE
		wire []byte
	}{
		{"type 1 tagged", ie.NewESMInformationTransferFlag(1), []byte{0xd1}},
		{"type 1 positional", ie.NewPositionalPDNType(3), []byte{3}},
		{"type 3 tagged", ie.NewEMMCause(ie.NetworkFailureMmCause), []byte{0x53, 0x11}},
		{"type 3 positional", ie.NewPositionalEMMCause(ie.NetworkFailureMmCause), []byte{0x11}},
		{"type 4 tagged", ie.NewEPSBearerContextStatus([]byte{0x20, 0x01}), []byte{0x57, 2, 0x20, 0x01}},
		{"type 4 positional", ie.NewPositionalEPSBearerContextStatus([]byte{0x20, 0x01}), []byte{2, 0x20, 0x01}},
		{"type 6 tagged", ie.NewESMMessageContainer([]byte{2, 1, 0xd0}), []byte{0x78, 0, 3, 2, 1, 0xd0}},
		{"type 6 positional", ie.NewPositionalESMMessageContainer([]byte{2, 1, 0xd0}), []byte{0, 3, 2, 1, 0xd0}},
		{"raw", ie.NewRaw(0x22, ie.IEFormatFour, []byte{1, 2}), []byte{0x22, 2, 1, 2}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			wire, err := tc.ie.Marshal()
			if err != nil || !bytes.Equal(wire, tc.wire) {
				t.Fatalf("Marshal = %x, %v; want %x", wire, err, tc.wire)
			}
			if tc.ie.MarshalLen() != len(wire) {
				t.Fatalf("MarshalLen = %d, want %d", tc.ie.MarshalLen(), len(wire))
			}
			for n := 0; n < len(wire); n++ {
				buffer := bytes.Repeat([]byte{0xff}, n)
				before := bytes.Clone(buffer)
				if err := tc.ie.MarshalTo(buffer); !errors.Is(err, io.ErrShortBuffer) || !bytes.Equal(buffer, before) {
					t.Fatalf("MarshalTo(%d) = %v, buffer = %x", n, err, buffer)
				}
			}
		})
	}
}

func TestIE_UnmarshalBinary(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name  string
		value *ie.IE
		wire  []byte
		kind  ie.Kind
	}{
		{"tagged", ie.NewT3448Value(0), []byte{0x6b, 1, 0x21}, ie.T3448Value},
		{"positional", ie.NewPositionalGPRSTimer2(0), []byte{1, 0x21}, ie.NoTag},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if err := tc.value.UnmarshalBinary(tc.wire); err != nil || tc.value.Kind != tc.kind || tc.value.Length != 1 || tc.value.Payload[0] != 0x21 {
				t.Fatalf("UnmarshalBinary = %+v, %v", tc.value, err)
			}
			if tc.kind != ie.NoTag {
				before := *tc.value
				if err := tc.value.UnmarshalBinary([]byte{0x5f, 1, 0x21}); !errors.Is(err, ie.ErrInvalidTag) || !reflect.DeepEqual(*tc.value, before) {
					t.Fatalf("wrong tag = %v", err)
				}
			}
		})
	}
}

func TestIE_UnmarshalBinaryAcceptsUnknownExtensions(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		value   *ie.IE
		wire    []byte
		payload []byte
	}{
		{
			name:    "type 4",
			value:   ie.NewEPSNetworkFeatureSupport(nil),
			wire:    []byte{byte(ie.EPSNetworkFeatureSupport), 4, 0, 0, 0, 0},
			payload: []byte{0, 0, 0, 0},
		},
		{
			name:    "type 6",
			value:   ie.NewCipheringKeyData(nil),
			wire:    append([]byte{byte(ie.CipheringKeyData), 0x08, 0xf1}, make([]byte, 2289)...),
			payload: make([]byte, 2289),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if err := tc.value.UnmarshalBinary(tc.wire); err != nil {
				t.Fatalf("UnmarshalBinary = %v", err)
			}
			if tc.value.Length != uint16(len(tc.payload)) || !bytes.Equal(tc.value.Payload, tc.payload) {
				t.Fatalf("decoded IE = %+v, want payload %x", tc.value, tc.payload)
			}
			if err := tc.value.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
				t.Fatalf("Validate = %v, want ErrInvalidLength", err)
			}
			if _, err := tc.value.Marshal(); !errors.Is(err, ie.ErrInvalidLength) {
				t.Fatalf("Marshal = %v, want ErrInvalidLength", err)
			}
		})
	}
}

func TestIE_MarshalWithoutTag(t *testing.T) {
	t.Parallel()
	i := ie.NewT3448Value(0x21)
	tag := i.Kind
	wire, err := i.MarshalWithoutTag()
	if err != nil || !bytes.Equal(wire, []byte{1, 0x21}) || i.Kind != tag {
		t.Fatalf("MarshalWithoutTag = %x, %v", wire, err)
	}
	if err := i.UnmarshalWithoutTag(wire); err != nil || i.Kind != ie.NoTag || i.Payload[0] != 0x21 {
		t.Fatalf("UnmarshalWithoutTag = %+v, %v", i, err)
	}
}

func TestIE_InvalidEnvelopes(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name  string
		value *ie.IE
		err   error
	}{
		{"nil", nil, ie.ErrNilIE},
		{"zero format", &ie.IE{}, ie.ErrInvalidFormat},
		{"type 1 overflow", ie.NewPositionalPDNType(0x10), ie.ErrInvalidValue},
		{"raw fixed unspecified length", &ie.IE{Kind: 0x22, Format: ie.IEFormatThree, Payload: []byte{1}}, ie.ErrInvalidLength},
		{"type 4 overflow", ie.NewRaw(0x22, ie.IEFormatFour, make([]byte, 256)), ie.ErrInvalidLength},
		{"incomplete APN AMBR", ie.NewAPNAggregateMaximumBitRate(make([]byte, 3)), ie.ErrInvalidLength},
		{"incomplete EPS QoS", ie.NewEPSQualityOfService(make([]byte, 4)), ie.ErrInvalidLength},
		{"incomplete PLMN", ie.NewPLMNList(make([]byte, 4)), ie.ErrInvalidLength},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if err := tc.value.Validate(); !errors.Is(err, tc.err) {
				t.Fatalf("Validate = %v, want %v", err, tc.err)
			}
		})
	}
}

func TestIE_UESecurityCapabilityLengths(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		length int
		valid  bool
	}{
		{name: "base", length: 2, valid: true},
		{name: "incomplete optional pair", length: 3},
		{name: "Iu mode", length: 4, valid: true},
		{name: "Iu and Gb modes", length: 5, valid: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := ie.NewPositionalUESecurityCapability(make([]byte, tc.length)).Validate()
			if tc.valid && err != nil {
				t.Fatalf("Validate = %v", err)
			}
			if !tc.valid && !errors.Is(err, ie.ErrInvalidLength) {
				t.Fatalf("Validate = %v, want ErrInvalidLength", err)
			}
		})
	}
}

func TestIE_NilAndRaw(t *testing.T) {
	t.Parallel()
	var nilIE *ie.IE
	if nilIE.MarshalLen() != 0 || nilIE.MarshalLenWithoutTag() != 0 || nilIE.Name() != "Unknown" {
		t.Fatal("nil receiver helpers")
	}
	if err := nilIE.UnmarshalBinary(nil); !errors.Is(err, ie.ErrNilIE) {
		t.Fatalf("nil UnmarshalBinary = %v", err)
	}
	raw := &ie.IE{Format: ie.IEFormatThree, Length: 2}
	if err := raw.UnmarshalWithoutTag([]byte{1, 2, 3}); err != nil || !bytes.Equal(raw.Payload, []byte{1, 2}) {
		t.Fatalf("raw fixed decoding = %+v, %v", raw, err)
	}
}

func FuzzIE(f *testing.F) {
	f.Add(true, []byte{0x57, 2, 0x20, 0})
	f.Add(false, []byte{1, 0x21})
	f.Fuzz(func(t *testing.T, tagged bool, wire []byte) {
		var (
			value *ie.IE
			err   error
		)
		if tagged {
			value, err = ie.Parse(wire, ie.EPSBearerContextStatus)
		} else {
			value, err = ie.ParseWithoutTag(wire, ie.NewPositionalGPRSTimer2(0))
		}
		if err != nil {
			return
		}
		encoded, err := value.Marshal()
		if err != nil {
			t.Fatalf("decoded IE cannot marshal: %v", err)
		}
		if !bytes.Equal(encoded, wire[:len(encoded)]) {
			t.Fatalf("round trip = %x, want %x", encoded, wire[:len(encoded)])
		}
	})
}
