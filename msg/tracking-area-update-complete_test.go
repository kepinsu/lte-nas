package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestTAUCompleteRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMTrackingAreaUpdateComplete,
	}

	m := msg.NewTrackingAreaUpdateComplete()
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.TrackingAreaUpdateComplete
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

	if got := decoded.MessageType(); got != msg.EMMTrackingAreaUpdateComplete {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMTrackingAreaUpdateComplete)
	}
	if got := decoded.MessageTypeName(); got != "Tracking area update complete" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Tracking area update complete")
	}
	if !decoded.IsRequest() {
		t.Error("IsRequest = false, want true")
	}

	if len(decoded.IEs) != 0 {
		t.Error("unexpected IEs")
	}
}

func TestTAUCompleteWithOptionalTags(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMTrackingAreaUpdateComplete,
	}

	ta := ie.NewTrackingAreaIdentity([]byte{0x23, 0x1A, 0xAB, 0xCD, 0xEF})

	m := msg.NewTrackingAreaUpdateComplete(ta)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.TrackingAreaUpdateComplete
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

	if len(decoded.IEs) != 1 {
		t.Fatalf("expected 1 IE, got %d", len(decoded.IEs))
	}
	if decoded.IEs[0].Kind != ie.TrackingAreaIdentity {
		t.Errorf("IE[0] kind = %v, want TrackingAreaIdentity", decoded.IEs[0].Kind)
	}
}

func TestTAUCompleteProtectedHeader(t *testing.T) {
	t.Parallel()

	m := msg.NewTrackingAreaUpdateComplete()

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
		SequenceNumber:        3,
		MsgType:               msg.EMMTrackingAreaUpdateComplete,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.TrackingAreaUpdateComplete
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.EmmHeader.SecurityHeaderType != msg.SecurityHeaderTypeIntegrityProtected {
		t.Errorf("SHT = %#02x, want %#02x", decoded.EmmHeader.SecurityHeaderType, msg.SecurityHeaderTypeIntegrityProtected)
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

func TestTAUCompleteErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.TrackingAreaUpdateComplete
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

	var decoded msg.TrackingAreaUpdateComplete
	err = decoded.UnmarshalBinary([]byte{0x17, 0xAA, 0xBB}) // truncated protected header (need 7 bytes)
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated protected = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0x1E, 0x4A}) // invalid PD
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid PD = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}

	// Plain header-only is valid (no mandatory body for TAU Complete)
	err = decoded.UnmarshalBinary([]byte{0x07, 0x4A})
	if err != nil {
		t.Fatalf("header-only = %v, want nil", err)
	}
	if decoded.EmmHeader.MsgType != msg.EMMTrackingAreaUpdateComplete {
		t.Errorf("MsgType = %#02x, want %#02x", decoded.EmmHeader.MsgType, msg.EMMTrackingAreaUpdateComplete)
	}
	if len(decoded.IEs) != 0 {
		t.Errorf("IEs = %d, want 0", len(decoded.IEs))
	}
}
