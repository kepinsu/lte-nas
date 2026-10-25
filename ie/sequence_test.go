package ie_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestParseMultiIEs(t *testing.T) {
	t.Parallel()
	got, err := ie.ParseMultiIEs([]byte{
		0x57, 2, 0x20, 0x01,
		0x53, 0x11,
		0xd1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 || got[0].Kind != ie.EPSBearerContextStatus || got[1].Kind != ie.EMMCause || got[2].Kind != ie.ESMInformationTransferFlag {
		t.Fatalf("Kind values = %+v", got)
	}
	if !bytes.Equal(got[0].Payload, []byte{0x20, 0x01}) || !bytes.Equal(got[1].Payload, []byte{0x11}) || !bytes.Equal(got[2].Payload, []byte{1}) {
		t.Fatalf("payloads = %+v", got)
	}
}

func TestParseMultiIEsAmbiguousTag(t *testing.T) {
	t.Parallel()
	if _, err := ie.ParseMultiIEs([]byte{0x58, 1}); !errors.Is(err, ie.ErrAmbiguousIE) {
		t.Fatalf("ParseMultiIEs = %v, want ErrAmbiguousIE", err)
	}
}

func TestParseMultiIEsSkipsUnknownNonComprehensionRequired(t *testing.T) {
	t.Parallel()
	got, err := ie.ParseMultiIEs([]byte{0x22, 1, 0xff, 0x53, 0x11})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Kind != ie.EMMCause {
		t.Fatalf("ParseMultiIEs = %+v", got)
	}
}

func TestParseMultiIEsRejectsUnknownComprehensionRequired(t *testing.T) {
	t.Parallel()
	// Type-1 unknown tag 0x01 requires comprehension — return ErrUnknownIE.
	if _, err := ie.ParseMultiIEs([]byte{0x01, 0}); !errors.Is(err, ie.ErrUnknownIE) {
		t.Fatalf("ParseMultiIEs = %v, want ErrUnknownIE", err)
	}
	// Type-6 Raw bytes 0x7E/0x7F are decoded as opaque payload when well-formed.
	// This is per the requirement that raw type-6 tags are returned by ParseMultiIEs.
	got, err := ie.ParseMultiIEs([]byte{0x7e, 0, 0})
	if err != nil {
		t.Fatalf("ParseMultiIEs(7e 00 00) = error: %v, want decoded IE", err)
	}
	if len(got) != 1 || got[0].Kind != 0x7e {
		t.Fatalf("ParseMultiIEs(7e 00 00) = %v, want single IE with kind 0x7e", got)
	}
}

func TestParseMultiIEsTruncated(t *testing.T) {
	t.Parallel()
	for _, wire := range [][]byte{
		{0x57},
		{0x57, 2, 0x20},
		{0x22},
		{0x7d, 0},
	} {
		if _, err := ie.ParseMultiIEs(wire); !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("ParseMultiIEs(%x) = %v, want io.ErrUnexpectedEOF", wire, err)
		}
	}
}

func TestHalfOctets(t *testing.T) {
	t.Parallel()
	low := ie.NewPositionalPDNType(3)
	high := ie.NewPositionalRequestType(5)
	wire, err := ie.MarshalHalfOctets(low, high)
	if err != nil || !bytes.Equal(wire, []byte{0x53}) {
		t.Fatalf("MarshalHalfOctets = %x, %v", wire, err)
	}
	parsedLow, parsedHigh, err := ie.ParseHalfOctets(wire, ie.NewPositionalPDNType(0), ie.NewPositionalRequestType(0))
	if err != nil || parsedLow.Payload[0] != 3 || parsedHigh.Payload[0] != 5 || parsedLow.Kind != ie.NoTag || parsedHigh.Kind != ie.NoTag {
		t.Fatalf("ParseHalfOctets = %+v, %+v, %v", parsedLow, parsedHigh, err)
	}
}
