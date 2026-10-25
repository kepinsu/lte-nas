package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestDetachAcceptRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
		SequenceNumber:        1,
		MsgType:               msg.EMMDetachAccept,
	}

	guti := ie.NewEPSMobileIdentity([]byte{0x40, 0x05, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE})
	emmCause := ie.NewEMMCause(ie.MmCause(15)) // #15 "MS identity info not provided"
	t3412 := ie.NewT3412Value(0x06)            // 60 seconds

	m := msg.NewDetachAccept(guti, emmCause, t3412)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.DetachAccept
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

	if got := decoded.MessageType(); got != msg.EMMDetachAccept {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMDetachAccept)
	}
	if got := decoded.MessageTypeName(); got != "Detach accept" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Detach accept")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}

	if decoded.GUTI == nil {
		t.Error("GUTI is nil")
	}
	if decoded.EMMCause == nil {
		t.Error("EMMCause is nil")
	}
	if decoded.T3412Value == nil {
		t.Error("T3412Value is nil")
	}
}

func TestDetachAcceptEmpty(t *testing.T) {
	t.Parallel()

	m := msg.NewDetachAccept()

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0xAA, 0xBB, 0xCC, 0xDD},
		SequenceNumber:        2,
		MsgType:               msg.EMMDetachAccept,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.DetachAccept
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

	if decoded.GUTI != nil {
		t.Error("unexpected GUTI")
	}
	if decoded.EMMCause != nil {
		t.Error("unexpected EMMCause")
	}
	if decoded.T3412Value != nil {
		t.Error("unexpected T3412Value")
	}
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs")
	}
}

func TestDetachAcceptErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.DetachAccept
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

	var decoded msg.DetachAccept
	err = decoded.UnmarshalBinary([]byte{0x1E, 0x46}) // invalid PD (E instead of 7)
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid PD = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}

	err = decoded.UnmarshalBinary([]byte{0x17, 0xAA, 0xBB}) // truncated protected header (need 7 bytes)
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated protected = %v, want %T", err, io.ErrUnexpectedEOF)
	}
}
