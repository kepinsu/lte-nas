package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestTAURequestRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMTrackingAreaUpdateRequest,
	}

	drx := ie.NewRaw(0x5c, ie.IEFormatFour, []byte{0x78, 0x01})
	guti := ie.NewEPSMobileIdentity([]byte{0x40, 0x05, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE})
	oldGUTI := ie.NewPositionalEPSMobileIdentity([]byte{0x40, 0x05, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE})
	ueNC := ie.NewUENetworkCapability([]byte{0xAA, 0xBB})
	taiList := ie.NewTAIList([]byte{0x23, 0x1A, 0xAB, 0xCD, 0xEF, 0x12, 0x34})

	m := msg.NewTrackingAreaUpdateRequest(
		ie.NewPositionalEPSAttachType(0x03), // EPS update type
		ie.NewPositionalNASKeySetIdentifier(0x01),
		oldGUTI,
		drx, guti, ueNC, taiList,
	)

	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.TrackingAreaUpdateRequest
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

	if got := decoded.MessageType(); got != msg.EMMTrackingAreaUpdateRequest {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMTrackingAreaUpdateRequest)
	}
	if got := decoded.MessageTypeName(); got != "Tracking area update request" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Tracking area update request")
	}
	if !decoded.IsRequest() {
		t.Error("IsRequest = false, want true")
	}

	if decoded.EPSUpdateType == nil {
		t.Fatal("EPSUpdateType is nil")
	}
	if decoded.NASKeySetIdentifier == nil {
		t.Fatal("NASKeySetIdentifier is nil")
	}
	if decoded.OldGUTI == nil {
		t.Fatal("OldGUTI is nil")
	}
	if decoded.DRXParameter == nil {
		t.Error("DRXParameter is nil")
	}
	if decoded.GUTI == nil {
		t.Error("GUTI is nil")
	}
	if decoded.UENetworkCapability == nil {
		t.Error("UENetworkCapability is nil")
	}
	if decoded.TAIList == nil {
		t.Error("TAIList is nil")
	}
}

func TestTAURequestEmptyOptional(t *testing.T) {
	t.Parallel()

	m := msg.NewTrackingAreaUpdateRequest(
		ie.NewPositionalEPSAttachType(0x03),
		ie.NewPositionalNASKeySetIdentifier(0x01),
		ie.NewPositionalEPSMobileIdentity([]byte{0x23, 0x05, 0x12, 0x34, 0xAB}),
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMTrackingAreaUpdateRequest,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.TrackingAreaUpdateRequest
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

	if decoded.DRXParameter != nil {
		t.Error("unexpected DRXParameter in empty optional")
	}
	if decoded.UENetworkCapability != nil {
		t.Error("unexpected UENetworkCapability in empty optional")
	}
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs in empty optional")
	}
}

func TestTAURequestProtectedHeader(t *testing.T) {
	t.Parallel()

	m := msg.NewTrackingAreaUpdateRequest(
		ie.NewPositionalEPSAttachType(0x03),
		ie.NewPositionalNASKeySetIdentifier(0x01),
		ie.NewPositionalEPSMobileIdentity([]byte{0x23, 0x05, 0x12, 0x34}),
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
		SequenceNumber:        5,
		MsgType:               msg.EMMTrackingAreaUpdateRequest,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.TrackingAreaUpdateRequest
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

func TestTAURequestErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.TrackingAreaUpdateRequest
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

	var decoded msg.TrackingAreaUpdateRequest
	err = decoded.UnmarshalBinary([]byte{0x07, 0x48}) // header only (plain, PD=7) - no payload
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header-only = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0x17, 0xAA, 0xBB}) // truncated protected header (need 7 bytes)
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated protected = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0x1E, 0x48}) // invalid PD
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid PD = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}
}

func TestTAURequestAllOptionalIEs(t *testing.T) {
	t.Parallel()

	m := msg.NewTrackingAreaUpdateRequest(
		ie.NewPositionalEPSAttachType(0x03), // EPS update type
		ie.NewPositionalNASKeySetIdentifier(0x01),
		ie.NewPositionalEPSMobileIdentity([]byte{0x40, 0x05, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE}),
		ie.NewRaw(0x5c, ie.IEFormatFour, []byte{0x78, 0x01}),
		ie.NewEPSMobileIdentity([]byte{0x40, 0x05, 0xAA, 0xBB}),
		ie.NewTAIList([]byte{0x23, 0x1A, 0xAB, 0xCD, 0xEF, 0x12, 0x34}),
		ie.NewEMMCause(0x16),
		ie.NewEPSBearerContextStatus([]byte{0x01, 0xFF}),
		ie.NewESMMessageContainer([]byte{0xD0, 0x01, 0xC1}),
		ie.NewTrackingAreaIdentity([]byte{0x04, 0x23, 0x1A, 0xAB, 0xCD}),
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMTrackingAreaUpdateRequest,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.TrackingAreaUpdateRequest
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.DRXParameter == nil {
		t.Error("DRXParameter is nil")
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
	if decoded.EPSBearerContextStatus == nil {
		t.Error("EPSBearerContextStatus is nil")
	}
	if decoded.ESMMessageContainer == nil {
		t.Error("ESMMessageContainer is nil")
	}
	if decoded.LastVisitedRegisteredTA == nil {
		t.Error("LastVisitedRegisteredTA is nil")
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

func TestTAURequestWithExtraIEs(t *testing.T) {
	t.Parallel()

	// Send an unrecognized IE tag (0x6f works with format 4 but not recognized here).
	m := msg.NewTrackingAreaUpdateRequest(
		ie.NewPositionalEPSAttachType(0x03),
		ie.NewPositionalNASKeySetIdentifier(0x01),
		ie.NewPositionalEPSMobileIdentity([]byte{0x40, 0x05, 0xAA, 0xBB}),
		ie.NewRaw(0x6f, ie.IEFormatFour, []byte{0xDE, 0xAD}),
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMTrackingAreaUpdateRequest,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.TrackingAreaUpdateRequest
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.GUTI != nil {
		t.Error("unexpected GUTI")
	}
	if len(decoded.IEs) != 1 {
		t.Fatalf("IEs length = %d, want 1", len(decoded.IEs))
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
