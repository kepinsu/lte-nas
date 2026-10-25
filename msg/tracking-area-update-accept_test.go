package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestTAUAcceptRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
		SequenceNumber:        1,
		MsgType:               msg.EMMTrackingAreaUpdateAccept,
	}

	t3412 := ie.NewT3412Value(0x06)
	guti := ie.NewEPSMobileIdentity([]byte{0x40, 0x05, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE})
	taiList := ie.NewTAIList([]byte{0x23, 0x1A, 0xAB, 0xCD, 0xEF, 0x12, 0x34})
	bearerCtx := ie.NewEPSBearerContextStatus([]byte{0x01, 0x02})
	emmCause := ie.NewEMMCause(ie.MmCause(15))
	locAreaID := ie.NewLocationAreaIdentification([]byte{0x23, 0x1A, 0xAB, 0xCD, 0xEF})

	m := msg.NewTrackingAreaUpdateAccept(
		ie.NewPositionalEPSAttachResult(0x01),
		ie.NewPositionalSpareHalfOctet(0x00),
		t3412, guti, taiList, bearerCtx, emmCause, locAreaID,
	)

	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.TrackingAreaUpdateAccept
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

	if got := decoded.MessageType(); got != msg.EMMTrackingAreaUpdateAccept {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMTrackingAreaUpdateAccept)
	}
	if got := decoded.MessageTypeName(); got != "Tracking area update accept" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Tracking area update accept")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}

	if decoded.EPSUpdateResult == nil {
		t.Fatal("EPSUpdateResult is nil")
	}
	if decoded.SpareHalfOctet == nil {
		t.Fatal("SpareHalfOctet is nil")
	}
	if decoded.T3412Value == nil {
		t.Error("T3412Value is nil")
	}
	if decoded.GUTI == nil {
		t.Error("GUTI is nil")
	}
	if decoded.TAIList == nil {
		t.Error("TAIList is nil")
	}
	if decoded.EMMCause == nil {
		t.Error("EMMCause is nil")
	}
	if decoded.LocationAreaIdentification == nil {
		t.Error("LocationAreaIdentification is nil")
	}
	if decoded.EPSBearerContextStatus == nil {
		t.Error("EPSBearerContextStatus is nil")
	}
}

func TestTAUAcceptEmptyOptional(t *testing.T) {
	t.Parallel()

	m := msg.NewTrackingAreaUpdateAccept(
		ie.NewPositionalEPSAttachResult(0x01),
		ie.NewPositionalSpareHalfOctet(0x00),
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0xAA, 0xBB, 0xCC, 0xDD},
		SequenceNumber:        2,
		MsgType:               msg.EMMTrackingAreaUpdateAccept,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.TrackingAreaUpdateAccept
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

	if decoded.T3412Value != nil {
		t.Error("unexpected T3412Value")
	}
	if decoded.GUTI != nil {
		t.Error("unexpected GUTI")
	}
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs")
	}
}

func TestTAUAcceptErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.TrackingAreaUpdateAccept
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

	var decoded msg.TrackingAreaUpdateAccept
	err = decoded.UnmarshalBinary([]byte{0x1E, 0x49}) // invalid PD
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid PD = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}

	err = decoded.UnmarshalBinary([]byte{0x17, 0xAA, 0xBB}) // truncated protected header (need 7 bytes)
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated protected = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Truncated content after valid header
	err = decoded.UnmarshalBinary([]byte{0x17, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0x49})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated content = %v, want %T", err, io.ErrUnexpectedEOF)
	}
}

func TestTAUAcceptWithExtraIEs(t *testing.T) {
	t.Parallel()

	// Include an unrecognized tag and a duplicate T3412Value.
	t3412 := ie.NewT3412Value(0x06)
	dupT3412 := ie.NewT3412Value(0x07)
	unknown := ie.NewRaw(0x6f, ie.IEFormatFour, []byte{0xDE, 0xAD})

	m := msg.NewTrackingAreaUpdateAccept(
		ie.NewPositionalEPSAttachResult(0x01),
		ie.NewPositionalSpareHalfOctet(0x00),
		t3412, dupT3412, unknown,
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0xAA, 0xBB, 0xCC, 0xDD},
		SequenceNumber:        3,
		MsgType:               msg.EMMTrackingAreaUpdateAccept,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.TrackingAreaUpdateAccept
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	// First T3412 should be stored in named field.
	if decoded.T3412Value == nil {
		t.Error("T3412Value is nil")
	}

	// Duplicate and unknown IEs should be in decoded.IEs.
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
