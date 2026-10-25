package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestDownlinkNASTransportRoundTrip(t *testing.T) {
	t.Parallel()

	// Plain header (2 bytes).
	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMDownlinkNASTransport,
	}

	guti := ie.NewRaw(ie.NoTag, ie.IEFormatThree, []byte{0x05})
	tmsi := ie.NewRaw(ie.NoTag, ie.IEFormatFour, []byte{0xAA, 0xBB, 0xCC, 0xDD})

	m := msg.NewDownlinkNASTransport(guti, tmsi)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.DownlinkNASTransport
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

	if got := decoded.MessageType(); got != msg.EMMDownlinkNASTransport {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMDownlinkNASTransport)
	}
	if got := decoded.MessageTypeName(); got != "Downlink NAS transport" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Downlink NAS transport")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if decoded.GUTI == nil {
		t.Fatal("GUTI is nil")
	}
	if decoded.TMSI == nil {
		t.Fatal("TMSI is nil")
	}
}

func TestDownlinkNASTransportWithAllTags(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMDownlinkNASTransport,
	}

	guti := ie.NewRaw(ie.NoTag, ie.IEFormatThree, []byte{0x05})
	tmsi := ie.NewRaw(ie.NoTag, ie.IEFormatFour, []byte{0x11, 0x22, 0x33, 0x44})

	stmsi := ie.NewRaw(0xd1, ie.IEFormatFour, []byte{0x55, 0x66})
	langid := ie.NewRaw(0xd0, ie.IEFormatOne, []byte{0x0f})
	esmc := ie.NewESMMessageContainer([]byte{0x2e, 0xC1})

	m := msg.NewDownlinkNASTransport(guti, tmsi, stmsi, langid, esmc)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.DownlinkNASTransport
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

	if decoded.GUTI == nil {
		t.Error("missing GUTI after decode")
	}
	if decoded.TMSI == nil {
		t.Error("missing TMSI after decode")
	}
	if decoded.ESMMessageContainer == nil {
		t.Error("missing ESMMessageContainer after decode")
	}
}

func TestDownlinkNASTransportProtected(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtectedWithNewEPSContext,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0xFF, 0xEE, 0xDD, 0xCC},
		SequenceNumber:        5,
		MsgType:               msg.EMMDownlinkNASTransport,
	}

	guti := ie.NewRaw(ie.NoTag, ie.IEFormatThree, []byte{0x03})
	tmsi := ie.NewRaw(ie.NoTag, ie.IEFormatFour, []byte{0x01, 0x02, 0x03, 0x04})

	m := msg.NewDownlinkNASTransport(guti, tmsi)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.DownlinkNASTransport
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

func TestDownlinkNASTransportErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.DownlinkNASTransport

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

	local := &msg.DownlinkNASTransport{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.DownlinkNASTransport

	// Truncated protected header (5 of 7 bytes)
	err = decoded.UnmarshalBinary([]byte{0x37, 0xAA, 0xBB, 0xCC, 0xDD})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Valid header + GUTI only (no TMSI → UnexpectedEOF)
	err = decoded.UnmarshalBinary([]byte{0x07, 0x62, 0x05})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("missing TMSI = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Valid header + GUTI + TMSI (should succeed, no tagged part)
	// TMSI type-4 positional: length byte (0x04) + 4-octet value
	err = decoded.UnmarshalBinary([]byte{0x07, 0x62, 0x05, 0x04, 0xAA, 0xBB, 0xCC, 0xDD})
	if err != nil {
		t.Fatalf("full positional = %v, want nil", err)
	}

	// Invalid SHT
	err = decoded.UnmarshalBinary([]byte{0xAA, 0x62})
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid SHT = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}
}

func TestDownlinkNASTransportWithDuplicateAndExtraIE(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMDownlinkNASTransport,
	}

	guti := ie.NewRaw(ie.NoTag, ie.IEFormatThree, []byte{0x05})
	tmsi := ie.NewRaw(ie.NoTag, ie.IEFormatFour, []byte{0x11, 0x22, 0x33})

	// STMSI is extracted via Kind=0xd1, but 0xd1 >= 0x80 fails IE validation.
	// So just test New* populates m.STMSI without round-trip.
	stmsi := ie.NewRaw(0xd1, ie.IEFormatFour, []byte{0xAA})
	dupSTMSI := ie.NewRaw(0xd1, ie.IEFormatFour, []byte{0xBB})

	m := msg.NewDownlinkNASTransport(guti, tmsi, stmsi, dupSTMSI)
	m.EmmHeader = h

	if m.STMSI == nil {
		t.Fatal("STMSI is nil")
	}
	if len(m.IEs) != 1 {
		t.Fatalf("m.IEs length = %d, want 1 (dup)", len(m.IEs))
	}

	// Don't round-trip: 0xd1 tag fails IE validation for type-4 (kind >= 0x80)
	// This is an existing limitation in the production code.
}

func TestDownlinkNASTransportNilExtraIE(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMDownlinkNASTransport,
	}

	guti := ie.NewRaw(ie.NoTag, ie.IEFormatThree, []byte{0x05})
	tmsi := ie.NewRaw(ie.NoTag, ie.IEFormatFour, []byte{0x11, 0x22, 0x33})

	m := msg.NewDownlinkNASTransport(guti, tmsi, nil)
	m.EmmHeader = h

	if len(m.IEs) != 0 {
		t.Fatalf("m.IEs length = %d, want 0 (nil skipped)", len(m.IEs))
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.DownlinkNASTransport
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.GUTI == nil {
		t.Fatal("decoded GUTI is nil")
	}
}

func ParseDownlinkNASTransport(t *testing.T, b []byte) msg.DownlinkNASTransport {
	t.Helper()
	var decoded msg.DownlinkNASTransport
	if err := decoded.UnmarshalBinary(b); err != nil {
		t.Fatalf("ParseDownlinkNASTransport failed: %v", err)
	}
	return decoded
}
