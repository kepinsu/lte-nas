package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestTAURejectRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
		SequenceNumber:        1,
		MsgType:               msg.EMMTrackingAreaUpdateReject,
	}

	emmCause := ie.NewEMMCause(ie.MmCause(15)) // #15 "EPS services not allowed"
	t3346 := ie.NewT3346Value(0x1E)            // 180 minutes
	exCause := ie.NewExtendedEMMCause(0x05)    // extended cause 5
	lbt := ie.NewLowerBoundTimerValue(0x06)    // GPRS timer 3 value

	forbiddenTAI := ie.NewForbiddenTAIsForRoaming([]byte{0x23, 0x1A, 0xAB, 0xCD, 0xEF, 0x12})

	m := msg.NewTrackingAreaUpdateReject(
		emmCause, lbt, t3346, exCause, forbiddenTAI,
	)

	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.TrackingAreaUpdateReject
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

	if got := decoded.MessageType(); got != msg.EMMTrackingAreaUpdateReject {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMTrackingAreaUpdateReject)
	}
	if got := decoded.MessageTypeName(); got != "Tracking area update reject" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Tracking area update reject")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}

	if decoded.EMMCause == nil {
		t.Fatal("EMMCause is nil")
	}
	if decoded.T3346Value == nil {
		t.Error("T3346Value is nil")
	}
	if decoded.ExtendedEMMCause == nil {
		t.Error("ExtendedEMMCause is nil")
	}
	if decoded.LowerBoundTimer == nil {
		t.Error("LowerBoundTimer is nil")
	}
	if decoded.ForbiddenTAI == nil {
		t.Error("ForbiddenTAI is nil")
	}
}

func TestTAURejectEmptyOptional(t *testing.T) {
	t.Parallel()

	m := msg.NewTrackingAreaUpdateReject(ie.NewEMMCause(ie.MmCause(15)))

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0xAA, 0xBB, 0xCC, 0xDD},
		SequenceNumber:        2,
		MsgType:               msg.EMMTrackingAreaUpdateReject,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.TrackingAreaUpdateReject
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

	if decoded.T3346Value != nil {
		t.Error("unexpected T3346Value")
	}
	if decoded.ExtendedEMMCause != nil {
		t.Error("unexpected ExtendedEMMCause")
	}
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs")
	}
}

func TestTAURejectErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.TrackingAreaUpdateReject
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

	var decoded msg.TrackingAreaUpdateReject
	err = decoded.UnmarshalBinary([]byte{0x7E, 0xAA}) // header: SHT=7 invalid
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid SHT = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}

	err = decoded.UnmarshalBinary([]byte{0x17, 0xAA, 0xBB}) // truncated protected header (need 7 bytes)
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated protected = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Missing positional EMM cause after valid header (only half-octets)
	err = decoded.UnmarshalBinary([]byte{0x17, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0x4B})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated content (no EMM cause) = %v, want %T", err, io.ErrUnexpectedEOF)
	}
}

func TestTAURejectWithExtraIEs(t *testing.T) {
	t.Parallel()

	// Include T3346, a duplicate T3346 (goes to IEs), and an unrecognized IE.
	m := msg.NewTrackingAreaUpdateReject(
		ie.NewEMMCause(ie.MmCause(15)), // Tracking area not allowed
		ie.NewT3346Value(0x07),
		ie.NewT3346Value(0x08),                               // duplicate → goes to IEs
		ie.NewRaw(0x6f, ie.IEFormatFour, []byte{0xDE, 0xAD}), // unrecognized → goes to IEs
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0xAA, 0xBB, 0xCC, 0xDD},
		SequenceNumber:        5,
		MsgType:               msg.EMMTrackingAreaUpdateReject,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.TrackingAreaUpdateReject
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.T3346Value == nil {
		t.Error("T3346Value is nil")
	}
	if len(decoded.IEs) != 2 {
		t.Fatalf("IEs length = %d, want 2", len(decoded.IEs))
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

func TestTAURejectWithAllOptionalIEs(t *testing.T) {
	t.Parallel()

	m := msg.NewTrackingAreaUpdateReject(
		ie.NewEMMCause(ie.MmCause(15)), // Tracking area not allowed
		ie.NewT3346Value(0x07),
		ie.NewExtendedEMMCause(0x0A),
		ie.NewLowerBoundTimerValue(0x1A),
		ie.NewForbiddenTAIsForRoaming([]byte{0x23, 0x1A, 0xAB, 0xCD, 0xEF, 0x12, 0x34}),
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0xAA, 0xBB, 0xCC, 0xDD},
		SequenceNumber:        2,
		MsgType:               msg.EMMTrackingAreaUpdateReject,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.TrackingAreaUpdateReject
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.T3346Value == nil {
		t.Error("T3346Value is nil")
	}
	if decoded.ExtendedEMMCause == nil {
		t.Error("ExtendedEMMCause is nil")
	}
	if decoded.LowerBoundTimer == nil {
		t.Error("LowerBoundTimer is nil")
	}
	if decoded.ForbiddenTAI == nil {
		t.Error("ForbiddenTAI is nil")
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
