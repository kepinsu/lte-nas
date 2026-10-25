package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestUplinkGenericNASTransportRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMUplinkGenericNASTransport,
	}

	container := ie.NewRaw(0x67, ie.IEFormatFour, []byte{0xDD, 0xEE, 0xFF})

	m := msg.NewUplinkGenericNASTransport(container)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.UplinkGenericNASTransport
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

	if got := decoded.MessageType(); got != msg.EMMUplinkGenericNASTransport {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMUplinkGenericNASTransport)
	}
	if got := decoded.MessageTypeName(); got != "Uplink generic NAS transport" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Uplink generic NAS transport")
	}
	if !decoded.IsRequest() {
		t.Error("IsRequest = false, want true")
	}
	if decoded.GenericNASContainer == nil {
		t.Fatal("GenericNASContainer is nil")
	}
}

func TestUplinkGenericNASTransportProtected(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x99, 0x88, 0x77, 0x66},
		SequenceNumber:        2,
		MsgType:               msg.EMMUplinkGenericNASTransport,
	}

	container := ie.NewRaw(0x67, ie.IEFormatFour, []byte{0xAA, 0xBB})

	m := msg.NewUplinkGenericNASTransport(container)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.UplinkGenericNASTransport
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

func TestUplinkGenericNASTransportWithExtraIEs(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMUplinkGenericNASTransport,
	}

	container := ie.NewRaw(0x67, ie.IEFormatFour, []byte{0xDD, 0xEE})
	extraIE := ie.NewRaw(0x68, ie.IEFormatFour, []byte{0xAA, 0xBB})

	m := msg.NewUplinkGenericNASTransport(container, extraIE)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.UplinkGenericNASTransport
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	rebuf := make([]byte, decoded.MarshalLen())
	if err := decoded.MarshalTo(rebuf); err != nil {
		t.Fatalf("re-MarshalTo failed: %v", err)
	}
	if len(buf) != len(rebuf) {
		t.Fatalf("re-marshal length = %d, want %d", len(rebuf), len(buf))
	}
	for i := range buf {
		if buf[i] != rebuf[i] {
			t.Fatalf("byte[%d] = %#02x, want %#02x", i, rebuf[i], buf[i])
		}
	}

	if len(decoded.IEs) != 1 {
		t.Fatalf("decoded IEs length = %d, want 1", len(decoded.IEs))
	}
	if decoded.IEs[0].Kind != 0x68 {
		t.Errorf("IE[0] Kind = 0x%02x, want 0x68", decoded.IEs[0].Kind)
	}
}

func TestUplinkGenericNASTransportErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.UplinkGenericNASTransport

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

	local := &msg.UplinkGenericNASTransport{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.UplinkGenericNASTransport

	// Truncated protected header (5 of 7 bytes)
	err = decoded.UnmarshalBinary([]byte{0x37, 0xAA, 0xBB, 0xCC, 0xDD})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Plain header only (no content)
	err = decoded.UnmarshalBinary([]byte{0x07, 0x69})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("no content = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Invalid SHT
	err = decoded.UnmarshalBinary([]byte{0xAA, 0x69})
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid SHT = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}
}

// ParseUplinkGenericNASTransport decodes b into a decoded msg.UplinkGenericNASTransport.
func ParseUplinkGenericNASTransport(t *testing.T, b []byte) msg.UplinkGenericNASTransport {
	t.Helper()
	var decoded msg.UplinkGenericNASTransport
	if err := decoded.UnmarshalBinary(b); err != nil {
		t.Fatalf("ParseUplinkGenericNASTransport failed: %v", err)
	}
	return decoded
}
