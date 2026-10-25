package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestServiceAcceptRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x55, 0x66, 0x77, 0x88},
		SequenceNumber:        2,
		MsgType:               msg.EMMServiceAccept,
	}

	m := msg.NewServiceAccept(
		ie.NewEPSBearerContextStatus([]byte{0x01, 0xff}),
		ie.NewT3448Value(0xAA),
	)

	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ServiceAccept
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

	if got := decoded.MessageType(); got != msg.EMMServiceAccept {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMServiceAccept)
	}
	if got := decoded.MessageTypeName(); got != "Service accept" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Service accept")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if decoded.EPSBearerContextStatus == nil {
		t.Fatal("EpsBearerContextStatus is nil")
	}
}

func TestServiceAcceptEmptyContent(t *testing.T) {
	t.Parallel()

	m := msg.NewServiceAccept()

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0xAA, 0xBB, 0xCC, 0xDD},
		SequenceNumber:        1,
		MsgType:               msg.EMMServiceAccept,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ServiceAccept
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

	if decoded.EPSBearerContextStatus != nil {
		t.Error("unexpected EpsBearerContextStatus")
	}
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs")
	}
}

func TestServiceAcceptProtectedHeaders(t *testing.T) {
	t.Parallel()

	shts := []uint8{
		msg.SecurityHeaderTypeIntegrityProtected,
		msg.SecurityHeaderTypeIntegrityProtectedAndCiphered,
		msg.SecurityHeaderTypeIntegrityProtectedWithNewEPSContext,
		msg.SecurityHeaderTypeIntegrityProtectedAndCipheredWithNewEPSContext,
	}

	for _, sht := range shts {
		t.Run("", func(t *testing.T) {
			m := msg.NewServiceAccept(ie.NewEPSBearerContextStatus([]byte{0x01, 0xff}))
			m.EmmHeader = msg.EmmHeader{
				SecurityHeaderType:    sht,
				ProtocolDiscriminator: byte(msg.PdEmm),
				Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
				SequenceNumber:        5,
				MsgType:               msg.EMMServiceAccept,
			}

			buf := make([]byte, m.MarshalLen())
			if err := m.MarshalTo(buf); err != nil {
				t.Fatalf("MarshalTo failed: %v", err)
			}

			var decoded msg.ServiceAccept
			if err := decoded.UnmarshalBinary(buf); err != nil {
				t.Fatalf("UnmarshalBinary failed: %v", err)
			}

			if decoded.EmmHeader.SecurityHeaderType != sht {
				t.Errorf("SHT = %#02x, want %#02x", decoded.EmmHeader.SecurityHeaderType, sht)
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
		})
	}
}

func TestServiceAcceptWithDuplicateAndExtraIE(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0xAA, 0xBB, 0xCC, 0xDD},
		SequenceNumber:        1,
		MsgType:               msg.EMMServiceAccept,
	}

	firstStatus := ie.NewEPSBearerContextStatus([]byte{0x01, 0xff})
	dupStatus := ie.NewEPSBearerContextStatus([]byte{0x02, 0xee})   // duplicate
	extraIE := ie.NewRaw(0x6f, ie.IEFormatFour, []byte{0xDE, 0xAD}) // unrecognized

	m := msg.NewServiceAccept(firstStatus, dupStatus, extraIE)
	m.EmmHeader = h

	// First goes to EPSBearerContextStatus, second is a duplicate → m.IEs,
	// unrecognized IE also goes to m.IEs.
	if m.EPSBearerContextStatus == nil {
		t.Fatal("EPSBearerContextStatus is nil")
	}
	if len(m.IEs) != 2 {
		t.Fatalf("m.IEs length = %d, want 2 (dup + unknown)", len(m.IEs))
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	// UnmarshalBinary decodes the extra IEs back into m.IEs.
	var decoded msg.ServiceAccept
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

func TestServiceAcceptNilExtraIE(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMServiceAccept,
	}

	nilIE := (*ie.IE)(nil)
	extraIE := ie.NewRaw(0x6f, ie.IEFormatFour, []byte{0xDE})

	m := msg.NewServiceAccept(nilIE, extraIE)
	m.EmmHeader = h

	if len(m.IEs) != 1 {
		t.Fatalf("m.IEs length = %d, want 1 (nil skipped)", len(m.IEs))
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ServiceAccept
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

func TestServiceAcceptErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.ServiceAccept

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

	local := &msg.ServiceAccept{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.ServiceAccept
	err = decoded.UnmarshalBinary([]byte{0x17, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE}) // truncated header - too few bytes for protected
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header-only = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0xAA, 0x4F}) // invalid SHT
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid SHT = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}
}
