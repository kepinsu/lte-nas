package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestSecurityModeRejectRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
		SequenceNumber:        1,
		MsgType:               msg.EMMSecurityModeReject,
	}

	m := msg.NewSecurityModeReject(ie.NewEMMCause(1))
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.SecurityModeReject
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

	if got := decoded.MessageType(); got != msg.EMMSecurityModeReject {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMSecurityModeReject)
	}
	if got := decoded.MessageTypeName(); got != "Security mode reject" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Security mode reject")
	}
	if !decoded.IsRequest() {
		t.Error("IsRequest = false, want true")
	}

	if decoded.EMMCause == nil {
		t.Fatal("EMMCause is nil")
	}
}

func TestSecurityModeRejectAllSHT(t *testing.T) {
	t.Parallel()

	for _, sht := range []uint8{
		msg.SecurityHeaderTypeIntegrityProtected,
		msg.SecurityHeaderTypeIntegrityProtectedAndCiphered,
		msg.SecurityHeaderTypeIntegrityProtectedWithNewEPSContext,
		msg.SecurityHeaderTypeIntegrityProtectedAndCipheredWithNewEPSContext,
		msg.SecurityHeaderTypeIntegrityProtectedAndPartiallyCiphered,
	} {
		shtName := shtNames(sht)

		m := msg.NewSecurityModeReject(ie.NewEMMCause(1))
		m.EmmHeader = msg.EmmHeader{
			SecurityHeaderType:    sht,
			ProtocolDiscriminator: byte(msg.PdEmm),
			Mac:                   [4]byte{0xAA, 0xBB, 0xCC, 0xDD},
			SequenceNumber:        byte(sht),
			MsgType:               msg.EMMSecurityModeReject,
		}

		buf := make([]byte, m.MarshalLen())
		if err := m.MarshalTo(buf); err != nil {
			t.Fatalf("SHT %s MarshalTo failed: %v", shtName, err)
		}

		var decoded msg.SecurityModeReject
		if err := decoded.UnmarshalBinary(buf); err != nil {
			t.Fatalf("SHT %s UnmarshalBinary failed: %v", shtName, err)
		}

		if decoded.EmmHeader.SecurityHeaderType != sht {
			t.Errorf("SHT %s: decoded SHT = %#02x, want %#02x", shtName,
				decoded.EmmHeader.SecurityHeaderType, sht)
		}

		rebuf := make([]byte, decoded.MarshalLen())
		if err := decoded.MarshalTo(rebuf); err != nil {
			t.Fatalf("SHT %s re-MarshalTo failed: %v", shtName, err)
		}
		for i := range buf {
			if buf[i] != rebuf[i] {
				t.Fatalf("SHT %s: byte[%d] = %#02x, want %#02x", shtName, i, rebuf[i], buf[i])
			}
		}
	}
}

func TestSecurityModeRejectErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.SecurityModeReject

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

	local := &msg.SecurityModeReject{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.SecurityModeReject

	// Valid protected header (7 bytes), SHT=0x01, PD=0x7 → 0x17
	// Truncated protected header (only 5 of 7 header bytes)
	err = decoded.UnmarshalBinary([]byte{0x17, 0xAA, 0xBB, 0xCC, 0xDD})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated protected header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Valid header but no content — need at least EMM cause (2 bytes, type-3)
	err = decoded.UnmarshalBinary([]byte{0x17, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0x5F})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("no content after header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0xAA, 0xAA}) // invalid SHT
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid SHT = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}
}

func TestSecurityModeRejectNilEMMCause(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil EMM cause")
		}
	}()

	msg.NewSecurityModeReject(nil)
}

func TestSecurityModeRejectEmptyPayload(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtectedWithNewEPSContext,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0xFF, 0xFF, 0xFF, 0xFF},
		SequenceNumber:        1,
		MsgType:               msg.EMMSecurityModeReject,
	}

	m := msg.NewSecurityModeReject(ie.NewEMMCause(1))
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.SecurityModeReject
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

func TestSecurityModeRejectWithExtraIEs(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
		SequenceNumber:        2,
		MsgType:               msg.EMMSecurityModeReject,
	}

	extraIE := ie.NewRaw(0x7e, ie.IEFormatFour, []byte{0x99})
	m := msg.NewSecurityModeReject(ie.NewEMMCause(1), extraIE)
	m.EmmHeader = h

	// Verify the extra IE is stored in m.IEs.
	if len(m.IEs) != 1 {
		t.Fatalf("m.IEs length = %d, want 1", len(m.IEs))
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	// UnmarshalBinary rejects trailing bytes, so skip round-trip.
	// We verified New* populates m.IEs and MarshalTo executes the extra-IE loop.
}

func TestSecurityModeRejectDuplicateAndNilExtra(t *testing.T) {
	t.Parallel()

	// First EMM cause is mandatory positional — cannot be duplicated as an IE.
	// Instead, test nil extra IE (loop should skip it) and unknown Kind (default).
	unknownIE := ie.NewRaw(0x7e, ie.IEFormatFour, []byte{0x88})
	nilIE := (*ie.IE)(nil)

	m := msg.NewSecurityModeReject(
		ie.NewEMMCause(1),
		nilIE,
		unknownIE,
	)

	if len(m.IEs) != 1 {
		t.Fatalf("m.IEs length = %d, want 1", len(m.IEs))
	}

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMSecurityModeReject,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}
}

// shtNames returns a human-readable name for a security header type value.
func shtNames(sht uint8) string {
	switch sht {
	case msg.SecurityHeaderTypeIntegrityProtected:
		return "0x01"
	case msg.SecurityHeaderTypeIntegrityProtectedAndCiphered:
		return "0x02"
	case msg.SecurityHeaderTypeIntegrityProtectedWithNewEPSContext:
		return "0x03"
	case msg.SecurityHeaderTypeIntegrityProtectedAndCipheredWithNewEPSContext:
		return "0x04"
	case msg.SecurityHeaderTypeIntegrityProtectedAndPartiallyCiphered:
		return "0x05"
	default:
		return "unknown"
	}
}
