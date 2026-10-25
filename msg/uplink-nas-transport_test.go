package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestUplinkNASTransportRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMUplinkNASTransport,
	}

	stmsi := ie.NewRaw(0xd1, ie.IEFormatFour, []byte{0x55, 0x66})
	esmc := ie.NewESMMessageContainer([]byte{0x2e, 0xD0})
	caps := ie.NewRaw(0xc0, ie.IEFormatOne, []byte{0x03})

	m := msg.NewUplinkNASTransport(stmsi, esmc, caps)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.UplinkNASTransport
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

	if got := decoded.MessageType(); got != msg.EMMUplinkNASTransport {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMUplinkNASTransport)
	}
	if got := decoded.MessageTypeName(); got != "Uplink NAS transport" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Uplink NAS transport")
	}
	if !decoded.IsRequest() {
		t.Error("IsRequest = false, want true")
	}
	if decoded.ESMMessageContainer == nil {
		t.Error("missing ESMMessageContainer after decode")
	}
}

func TestUplinkNASTransportProtected(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
		SequenceNumber:        1,
		MsgType:               msg.EMMUplinkNASTransport,
	}

	stmsi := ie.NewRaw(0xd1, ie.IEFormatFour, []byte{0xAA, 0xBB})
	esmc := ie.NewESMMessageContainer([]byte{0x2e, 0xC1})
	caps := ie.NewRaw(0xc0, ie.IEFormatOne, []byte{0x01})

	m := msg.NewUplinkNASTransport(stmsi, esmc, caps)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.UplinkNASTransport
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
}

func TestUplinkNASTransportEmptyOptional(t *testing.T) {
	t.Parallel()

	m := msg.NewUplinkNASTransport()
	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMUplinkNASTransport,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.UplinkNASTransport
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

	if decoded.STMSI != nil {
		t.Error("unexpected STMSI")
	}
	if decoded.ESMMessageContainer != nil {
		t.Error("unexpected ESMMessageContainer")
	}
	if decoded.EMMCapabilities != nil {
		t.Error("unexpected EMMCapabilities")
	}
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs")
	}
}

func TestUplinkNASTransportErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.UplinkNASTransport

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

	local := &msg.UplinkNASTransport{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.UplinkNASTransport

	// Truncated protected header
	err = decoded.UnmarshalBinary([]byte{0x37, 0xAA, 0xBB, 0xCC, 0xDD})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Plain header only (no content) — valid for Uplink NAS Transport
	err = decoded.UnmarshalBinary([]byte{0x07, 0x63})
	if err != nil {
		t.Fatalf("no content = %v, want nil", err)
	}
	if decoded.STMSI != nil {
		t.Error("unexpected STMSI")
	}

	// Invalid SHT
	err = decoded.UnmarshalBinary([]byte{0xAA, 0x63})
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid SHT = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}
}

func TestUplinkNASTransportParseHelper(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMUplinkNASTransport,
	}

	stmsi := ie.NewRaw(0xd1, ie.IEFormatFour, []byte{0xAA, 0xBB})
	caps := ie.NewRaw(0xc0, ie.IEFormatOne, []byte{0x03})

	m := msg.NewUplinkNASTransport(stmsi, caps)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.UplinkNASTransport
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

	if decoded.STMSI == nil {
		t.Fatal("decoded STMSI is nil")
	}

	// Verify Message interface methods.
	if decoded.MessageType() != msg.EMMUplinkNASTransport {
		t.Errorf("MessageType = %#02x, want %#02x", decoded.MessageType(), msg.EMMUplinkNASTransport)
	}
	if !decoded.IsRequest() {
		t.Error("IsRequest = false, want true")
	}
	if decoded.MessageTypeName() != "Uplink NAS transport" {
		t.Errorf("MessageTypeName = %q, want %q", decoded.MessageTypeName(), "Uplink NAS transport")
	}
}

// parseUplinkNASTags — indirect tests through UplinkNASTransport marshal/unmarshal
// that exercise unknown-type tag handling.

func TestUplinkNASTransportParseUnknownType4(t *testing.T) {
	t.Parallel()

	// Build a wire message with an unknown type-4 TLV IE (tag 0xAA, length 2).
	// The decoder's default case should treat it as a pass-through type-4 TLV.
	// EMM Uplink NAS Transport = 0x63

	wire := []byte{
		0x07, 0x63, // plain header: PTI=0, msg type 0x63
		0x78, 0x00, 0x01, 0xC1, // ESM message container (type-6 TLV-E): len=1
		0xAA, 0x02, 0x11, 0x22, // unknown type-4 TLV: tag=0xAA, len=2
	}

	var decoded msg.UplinkNASTransport
	err := decoded.UnmarshalBinary(wire)
	if err != nil {
		t.Fatalf("UnmarshalBinary = error: %v", err)
	}

	if decoded.ESMMessageContainer == nil {
		t.Error("missing ESMMessageContainer")
	}

	// The unknown IE should be in decoded.IEs (as a pass-through)
	found := false
	for _, gotIE := range decoded.IEs {
		if gotIE.Kind == 0xAA && len(gotIE.Payload) == 2 && gotIE.Payload[0] == 0x11 {
			found = true
		}
	}
	if !found {
		t.Errorf("unknown type-4 IE not found in decoded.IEs (got %d IEs)", len(decoded.IEs))
	}
}

func TestUplinkNASTransportParseUnknownType1(t *testing.T) {
	t.Parallel()

	// Build a wire message with an unknown type-1 TV IE.
	// Tag 0xA5: type = 0xA5 & 0x07 = 5, falls to default.
	// In parseUplinkNASTags default case, if len < 2 it falls to type-1 TV branch.
	// wire byte: payload_nibble (low) | tag_nibble (high)
	// If we write byte 0xA5 directly, parseUplinkNASTags creates:
	//   Kind = ie.Kind(0xA5 & 0xf0) = 0xA0, Payload = []byte{0x05}
	// EMM Uplink NAS Transport = 0x63

	// Plain header + unknown type-1 IE (wire byte 0xB3)
	wire := []byte{
		0x07, 0x63, // plain header
		0xB3, // unknown type-1: Kind=0xB0, Payload=[0x03]
	}

	var decoded msg.UplinkNASTransport
	err := decoded.UnmarshalBinary(wire)
	if err != nil {
		t.Fatalf("UnmarshalBinary = error: %v", err)
	}

	if len(decoded.IEs) != 1 {
		t.Fatalf("expected 1 IE in decoded.IEs, got %d", len(decoded.IEs))
	}

	gotIE := decoded.IEs[0]
	if gotIE.Kind != 0xB0 {
		t.Errorf("IE Kind = %#02x, want %#02x", gotIE.Kind, 0xB0)
	}
	if len(gotIE.Payload) != 1 || gotIE.Payload[0] != 0x03 {
		t.Errorf("IE Payload = %v, want [0x03]", gotIE.Payload)
	}
	if gotIE.Format != ie.IEFormatOne {
		t.Errorf("IE Format = %d, want %d", gotIE.Format, ie.IEFormatOne)
	}
}

func TestUplinkNASTransportParseMultipleUnknown(t *testing.T) {
	t.Parallel()

	// Mix known and unknown IEs.
	// S-TMSI (0xd1) + unknown type-4 + EMM Capabilities (type-1)
	// EMM Uplink NAS Transport = 0x63
	wire := []byte{
		0x07, 0x63, // plain header: PTI=0, msg type 0x63
		0xD1, 0x02, 0x55, 0x66, // S-TMSI (type-4): len=2, value=55 66
		0xAA, 0x01, 0x42, // unknown type-4 TLV: tag=0xAA, len=1, payload=0x42
		0xC3, // EMM Capabilities (type-1 TV): tag=0xc, value=0x3
	}

	var decoded msg.UplinkNASTransport
	err := decoded.UnmarshalBinary(wire)
	if err != nil {
		t.Fatalf("UnmarshalBinary = error: %v", err)
	}

	if decoded.STMSI == nil {
		t.Error("missing STMSI")
	}

	if decoded.EMMCapabilities == nil {
		t.Error("missing EMMCapabilities")
	}

	// Should have 1 unknown IE in decoded.IEs (the type-4 one)
	var foundUnknown bool
	for _, gotIE := range decoded.IEs {
		if gotIE.Kind == 0xAA {
			foundUnknown = true
			break
		}
	}
	if !foundUnknown {
		t.Errorf("unknown type-4 IE not found in decoded.IEs (got %d IEs: %+v)", len(decoded.IEs), decoded.IEs)
	}
}
