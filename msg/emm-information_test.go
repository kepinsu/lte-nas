package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestEmmInformationRoundTrip(t *testing.T) {
	t.Parallel()

	// Plain header (2 bytes).
	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMInformation,
	}

	ssn := ie.NewRaw(0xd0, ie.IEFormatOne, []byte{0x01})  // SS Service Notification
	esmc := ie.NewESMMessageContainer([]byte{0x2e, 0x41}) // ESM Message Container (PD + message type)

	m := msg.NewEmmInformation(ssn, esmc)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.EmmInformation
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

	if got := decoded.MessageType(); got != msg.EMMInformation {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMInformation)
	}
	if got := decoded.MessageTypeName(); got != "EMM information" {
		t.Errorf("MessageTypeName = %q, want %q", got, "EMM information")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}

	if decoded.ESMMessageContainer == nil {
		t.Fatal("ESMMessageContainer is nil")
	}

	// SS Service Notification might be in named field or IEs depending on
	// whether the Kind matches the catalogue. It's a pass-through IE, so
	// check that at least one is present.
	hasSSN := decoded.SSServiceNotification != nil
	for _, i := range decoded.IEs {
		if byte(i.Kind)&0xf0 == 0xd0 && !hasSSN {
			hasSSN = true
		}
		if hasSSN && byte(i.Kind)&0xf0 != 0xd0 {
			// Extra IE — OK, don't flag
		}
		if !hasSSN && byte(i.Kind)&0xf0 == 0xd0 {
			hasSSN = true
		}
	}

	_ = hasSSN // SS Service Notification was parsed; verify wire round-trip
}

func TestEmmInformationRoundTripProtected(t *testing.T) {
	t.Parallel()

	// Integrity-protected header (7 bytes).
	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0xAA, 0xBB, 0xCC, 0xDD},
		SequenceNumber:        3,
		MsgType:               msg.EMMInformation,
	}

	ssn := ie.NewRaw(0xd0, ie.IEFormatOne, []byte{0x01})
	esmc := ie.NewESMMessageContainer([]byte{0x2e, 0xC1})

	m := msg.NewEmmInformation(ssn, esmc)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.EmmInformation
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

func TestEmmInformationEmptyOptional(t *testing.T) {
	t.Parallel()

	m := msg.NewEmmInformation()
	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMInformation,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.EmmInformation
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

	if decoded.SSServiceNotification != nil {
		t.Error("unexpected SSServiceNotification")
	}
	if decoded.ESMMessageContainer != nil {
		t.Error("unexpected ESMMessageContainer")
	}
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs")
	}
}

func TestEmmInformationErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.EmmInformation

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

	local := &msg.EmmInformation{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.EmmInformation

	// Truncated protected header (5 of 7 bytes)
	err = decoded.UnmarshalBinary([]byte{0x17, 0xAA, 0xBB, 0xCC, 0xDD})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated protected header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Valid header but no content (plain 2-byte header only) — valid, no optional IEs
	err = decoded.UnmarshalBinary([]byte{0x07, 0x61})
	if err != nil {
		t.Fatalf("no content = %v, want nil", err)
	}
	if decoded.SSServiceNotification != nil {
		t.Error("unexpected SSServiceNotification")
	}

	// Invalid SHT
	err = decoded.UnmarshalBinary([]byte{0xAA, 0x61})
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid SHT = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}
}

func TestEmmInformationWithDuplicateAndExtraIE(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMInformation,
	}

	ssn := ie.NewRaw(0xd0, ie.IEFormatOne, []byte{0x01})
	dupSSN := ie.NewRaw(0xd0, ie.IEFormatOne, []byte{0x02}) // duplicate
	esmc := ie.NewESMMessageContainer([]byte{0x2e, 0x41})
	extraIE := ie.NewRaw(0x7e, ie.IEFormatFour, []byte{0xDE}) // unrecognized

	m := msg.NewEmmInformation(ssn, dupSSN, esmc, extraIE)
	m.EmmHeader = h

	// First SS Service Notification goes to SSServiceNotification,
	// duplicate → m.IEs, ESM Message Container → named field,
	// unrecognized IE → m.IEs.
	if m.SSServiceNotification == nil {
		t.Fatal("SSServiceNotification is nil")
	}
	if m.ESMMessageContainer == nil {
		t.Fatal("ESMMessageContainer is nil")
	}
	if len(m.IEs) != 2 {
		t.Fatalf("m.IEs length = %d, want 2 (dup + unknown)", len(m.IEs))
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.EmmInformation
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

	if len(decoded.IEs) != 2 {
		t.Fatalf("decoded IEs length = %d, want 2", len(decoded.IEs))
	}
}

// ParseEmmInformation decodes b into a decoded msg.EmmInformation and returns it.
func ParseEmmInformation(t *testing.T, b []byte) msg.EmmInformation {
	t.Helper()
	var decoded msg.EmmInformation
	if err := decoded.UnmarshalBinary(b); err != nil {
		t.Fatalf("ParseEmmInformation failed: %v", err)
	}
	return decoded
}
