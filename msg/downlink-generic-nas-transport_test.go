package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestDownlinkGenericNASTransportRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMDownlinkGenericNASTransport,
	}

	container := ie.NewRaw(0x67, ie.IEFormatFour, []byte{0xAA, 0xBB, 0xCC})

	m := msg.NewDownlinkGenericNASTransport(container)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.DownlinkGenericNASTransport
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	rebuf := make([]byte, decoded.MarshalLen())
	if err := decoded.MarshalTo(rebuf); err != nil {
		t.Fatalf("re-MarshalTo failed: %v", err)
	}
	for i := range buf {
		if buf[i] != rebuf[i] {
			t.Fatalf("byte[%d] = %#02x, want %#02x", i, rebuf[i], buf[i])
		}
	}

	if got := decoded.MessageType(); got != msg.EMMDownlinkGenericNASTransport {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMDownlinkGenericNASTransport)
	}
	if got := decoded.MessageTypeName(); got != "Downlink generic NAS transport" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Downlink generic NAS transport")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if decoded.GenericNASContainer == nil {
		t.Fatal("GenericNASContainer is nil")
	}
}

func TestDownlinkGenericNASTransportProtected(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtectedAndCipheredWithNewEPSContext,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0xFF, 0xEE, 0xDD, 0xCC},
		SequenceNumber:        4,
		MsgType:               msg.EMMDownlinkGenericNASTransport,
	}

	container := ie.NewRaw(0x67, ie.IEFormatFour, []byte{0x11, 0x22, 0x33, 0x44})

	m := msg.NewDownlinkGenericNASTransport(container)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.DownlinkGenericNASTransport
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	rebuf := make([]byte, decoded.MarshalLen())
	if err := decoded.MarshalTo(rebuf); err != nil {
		t.Fatalf("re-MarshalTo failed: %v", err)
	}
	for i := range buf {
		if buf[i] != rebuf[i] {
			t.Fatalf("byte[%d] = %#02x, want %#02x", i, rebuf[i], buf[i])
		}
	}

	if decoded.GenericNASContainer == nil {
		t.Fatal("GenericNASContainer is nil after decode")
	}
}

func TestDownlinkGenericNASTransportErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.DownlinkGenericNASTransport

	if got := nilMsg.MarshalLen(); got != 0 {
		t.Errorf("nil MarshalLen = %d, want 0", got)
	}

	err := nilMsg.MarshalTo(nil)
	if err == nil {
		t.Fatal("nil MarshalTo = nil, want error")
	}
	if !errors.Is(err, ie.ErrNilIE) {
		t.Errorf("nil MarshalTo = %v, want wrapped %T", err, ie.ErrNilIE)
	}

	local := &msg.DownlinkGenericNASTransport{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.DownlinkGenericNASTransport

	// Truncated protected header (5 of 7 bytes)
	err = decoded.UnmarshalBinary([]byte{0x37, 0xAA, 0xBB, 0xCC, 0xDD})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Plain header only (no content)
	err = decoded.UnmarshalBinary([]byte{0x07, 0x68})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("no content = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Invalid SHT
	err = decoded.UnmarshalBinary([]byte{0xAA, 0x68})
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid SHT = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}

	// Type-4 container IE with valid header + tagged content
	gutiContainer := ie.NewRaw(0x67, ie.IEFormatFour, []byte{0xAA, 0xBB})
	m := msg.NewDownlinkGenericNASTransport(gutiContainer)
	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMDownlinkGenericNASTransport,
	}

	buf = make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded2 msg.DownlinkGenericNASTransport
	if err := decoded2.UnmarshalBinary(buf); err != nil {
		t.Fatalf("decode with container = %v, want nil", err)
	}

	rebuf := make([]byte, decoded2.MarshalLen())
	if err := decoded2.MarshalTo(rebuf); err != nil {
		t.Fatalf("re-MarshalTo = %v, want nil", err)
	}
	for i := range buf {
		if buf[i] != rebuf[i] {
			t.Fatalf("byte[%d] = %#02x, want %#02x", i, rebuf[i], buf[i])
		}
	}
}

func TestDownlinkGenericNASTransportType1IE(t *testing.T) {
	t.Parallel()

	// Craft a raw message with a type-1 TV IE (single wire byte).
	// Wire: 0x78 (ESM container, type-6 TLV) + extra IE with tag 0xA8 length=0 (type-1 fallback).
	rawData := []byte{
		0x07, 0x68, // header (PD=EMM, SHT=plain, type=0x68)
		0x78, 0x00, 0x01, 0xAA, // ESM container (type-6 TLV)
		0xA8, // extra IE: tag=0xA8 (tag nibble 0xA), type-1 TV → payload 0x8
	}

	var decoded msg.DownlinkGenericNASTransport
	if err := decoded.UnmarshalBinary(rawData); err != nil {
		t.Fatalf("UnmarshalBinary: %v", err)
	}

	if decoded.GenericNASContainer == nil {
		t.Fatal("GenericNASContainer is nil")
	}
	if len(decoded.IEs) != 1 {
		t.Fatalf("expected 1 extra IE, got %d", len(decoded.IEs))
	}
}

// TestDownlinkGenericNASTransportExtraIEs tests the extra IE loop in
// parseTaggedContainerIEs by providing a message with multiple tagged IEs.
func TestDownlinkGenericNASTransportExtraIEs(t *testing.T) {
	t.Parallel()

	// Create a container with type-4 TLV and additional type-4 extra IEs
	container := ie.NewRaw(0x67, ie.IEFormatFour, []byte{0xAA, 0xBB})
	extra1 := ie.NewRaw(0x6A, ie.IEFormatFour, []byte{0xCC})
	extra2 := ie.NewRaw(0x6B, ie.IEFormatFour, []byte{0xDD})

	m := msg.NewDownlinkGenericNASTransport(container, extra1, extra2)
	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMDownlinkGenericNASTransport,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.DownlinkGenericNASTransport
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary: %v", err)
	}

	rebuf := make([]byte, decoded.MarshalLen())
	if err := decoded.MarshalTo(rebuf); err != nil {
		t.Fatalf("re-MarshalTo: %v", err)
	}
	for i := range buf {
		if buf[i] != rebuf[i] {
			t.Fatalf("byte[%d] = %#02x, want %#02x", i, rebuf[i], buf[i])
		}
	}

	if decoded.GenericNASContainer == nil {
		t.Fatal("GenericNASContainer is nil")
	}
	if len(decoded.IEs) != 2 {
		t.Fatalf("expected 2 extra IEs, got %d", len(decoded.IEs))
	}
}

// ParseDownlinkGenericNASTransport decodes b into a decoded msg.DownlinkGenericNASTransport.
func ParseDownlinkGenericNASTransport(t *testing.T, b []byte) msg.DownlinkGenericNASTransport {
	t.Helper()
	var decoded msg.DownlinkGenericNASTransport
	if err := decoded.UnmarshalBinary(b); err != nil {
		t.Fatalf("ParseDownlinkGenericNASTransport failed: %v", err)
	}
	return decoded
}
