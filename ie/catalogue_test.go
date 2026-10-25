package ie_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestTaggedKindsAreWireTags(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		kind ie.Kind
		want byte
	}{
		{"T3412 extended value", ie.T3412ExtendedValue, 0x5e},
		{"full network name", ie.FullNetworkName, 0x43},
		{"short network name", ie.ShortNetworkName, 0x45},
		{"nonce UE", ie.NonceUE, 0x55},
		{"nonce MME", ie.NonceMME, 0x56},
		{"ESM cause", ie.ESMCause, 0x58},
		{"UE network capability", ie.UENetworkCapability, 0x58},
		{"device properties", ie.DeviceProperties, 0xc0},
		{"ESM device properties", ie.ESMDeviceProperties, 0xd0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if byte(tc.kind) != tc.want {
				t.Fatalf("Kind = %#02x, want %#02x", tc.kind, tc.want)
			}
		})
	}
}

func TestTaggedConstructorUsesItsWireTag(t *testing.T) {
	t.Parallel()
	i := ie.NewT3412ExtendedValue(0x21)
	if i.Kind != ie.T3412ExtendedValue {
		t.Fatalf("Kind = %#02x, want %#02x", i.Kind, ie.T3412ExtendedValue)
	}
	wire, err := i.Marshal()
	if err != nil || !bytes.Equal(wire, []byte{0x5e, 1, 0x21}) {
		t.Fatalf("Marshal = %x, %v", wire, err)
	}

	decoded := ie.NewT3412ExtendedValue(0)
	if err := decoded.UnmarshalBinary(wire); err != nil || decoded.Kind != ie.T3412ExtendedValue || !bytes.Equal(decoded.Payload, []byte{0x21}) {
		t.Fatalf("UnmarshalBinary = %+v, %v", decoded, err)
	}
}

func TestPositionalConstructorStaysTagless(t *testing.T) {
	t.Parallel()
	i := ie.NewPositionalGPRSTimer(0x21)
	if i.Kind != ie.NoTag {
		t.Fatalf("Kind = %#02x, want NoTag", i.Kind)
	}
	wire, err := i.Marshal()
	if err != nil || !bytes.Equal(wire, []byte{0x21}) {
		t.Fatalf("Marshal = %x, %v", wire, err)
	}
	decoded, err := ie.ParseWithoutTag(wire, ie.NewPositionalGPRSTimer(0))
	if err != nil || decoded.Kind != ie.NoTag || !bytes.Equal(decoded.Payload, []byte{0x21}) {
		t.Fatalf("ParseWithoutTag = %+v, %v", decoded, err)
	}
}

func TestParseRequiresUnambiguousTaggedFormat(t *testing.T) {
	t.Parallel()
	parsed, err := ie.Parse([]byte{0x57, 2, 0x20, 0x01}, ie.EPSBearerContextStatus)
	if err != nil || parsed.Kind != ie.EPSBearerContextStatus || !bytes.Equal(parsed.Payload, []byte{0x20, 0x01}) {
		t.Fatalf("Parse = %+v, %v", parsed, err)
	}
	if _, err := ie.Parse([]byte{0x58, 1}, ie.ESMCause); !errors.Is(err, ie.ErrAmbiguousIE) {
		t.Fatalf("ambiguous Parse = %v", err)
	}
}
