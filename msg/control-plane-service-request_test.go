package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestControlPlaneServiceRequestRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtectedAndPartiallyCiphered,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0xAA, 0xBB, 0xCC, 0xDD},
		SequenceNumber:        7,
		MsgType:               msg.EMMControlPlaneServiceRequest,
	}

	m := msg.NewControlPlaneServiceRequest(
		ie.NewPositionalControlPlaneServiceType(0x02),
		ie.NewNASKeySetIdentifier(0x01),
		ie.NewESMMessageContainer([]byte{0xD0, 0x01, 0xC1}),
		ie.NewNASMessageContainer([]byte{0x67, 0x02}),
		ie.NewEPSBearerContextStatus([]byte{0x01, 0xff}),
		ie.NewDeviceProperties(0x01),
	)

	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ControlPlaneServiceRequest
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

	if got := decoded.MessageType(); got != msg.EMMControlPlaneServiceRequest {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMControlPlaneServiceRequest)
	}
	if got := decoded.MessageTypeName(); got != "Control plane service request" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Control plane service request")
	}
	if !decoded.IsRequest() {
		t.Error("IsRequest = false, want true")
	}
	if decoded.ControlPlaneServiceType == nil {
		t.Fatal("ControlPlaneServiceType is nil")
	}
	if decoded.NASKeySetIdentifier == nil {
		t.Fatal("NASKeySetIdentifier is nil")
	}
}

func TestControlPlaneServiceRequestPlainHeader(t *testing.T) {
	t.Parallel()

	m := msg.NewControlPlaneServiceRequest(
		ie.NewPositionalControlPlaneServiceType(0x01),
		ie.NewNASKeySetIdentifier(0x01),
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMControlPlaneServiceRequest,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ControlPlaneServiceRequest
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

func TestControlPlaneServiceRequestEmptyOptional(t *testing.T) {
	t.Parallel()

	m := msg.NewControlPlaneServiceRequest(
		ie.NewPositionalControlPlaneServiceType(0x01),
		ie.NewNASKeySetIdentifier(0x01),
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMControlPlaneServiceRequest,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ControlPlaneServiceRequest
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.EsmMessageContainer != nil {
		t.Error("unexpected EsmMessageContainer")
	}
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs")
	}
}

func TestControlPlaneServiceRequestErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.ControlPlaneServiceRequest

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

	local := &msg.ControlPlaneServiceRequest{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.ControlPlaneServiceRequest
	err = decoded.UnmarshalBinary([]byte{0x07, 0x4D}) // header only - no payload
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header-only = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0x17, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0x4D}) // truncated protected
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated protected = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0xAA, 0x4D}) // invalid SHT
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid SHT = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}
}

func TestControlPlaneServiceRequestAllOptionalIEs(t *testing.T) {
	t.Parallel()

	m := msg.NewControlPlaneServiceRequest(
		ie.NewPositionalControlPlaneServiceType(0x02),
		ie.NewNASKeySetIdentifier(0x01),
		ie.NewESMMessageContainer([]byte{0xD0, 0x01, 0xC1}),
		ie.NewNASMessageContainer([]byte{0x67, 0x02}),
		ie.NewEPSBearerContextStatus([]byte{0x01, 0xFF}),
		ie.NewDeviceProperties(0x01),
		ie.NewUERequestType(0x01),
		ie.NewPagingRestriction([]byte{0x11}),
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMControlPlaneServiceRequest,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ControlPlaneServiceRequest
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.EsmMessageContainer == nil {
		t.Error("EsmMessageContainer is nil")
	}
	if decoded.NasMessageContainer == nil {
		t.Error("NasMessageContainer is nil")
	}
	if decoded.EPSBearerContextStatus == nil {
		t.Error("EPSBearerContextStatus is nil")
	}
	if decoded.DeviceProperties == nil {
		t.Error("DeviceProperties is nil")
	}
	if decoded.UeRequestType == nil {
		t.Error("UeRequestType is nil")
	}
	if decoded.PagingRestriction == nil {
		t.Error("PagingRestriction is nil")
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

func TestControlPlaneServiceRequestWithExtraIEs(t *testing.T) {
	t.Parallel()

	// Send an unrecognized IE tag (0x6f works with format 4 but not recognized here).
	m := msg.NewControlPlaneServiceRequest(
		ie.NewPositionalControlPlaneServiceType(0x01),
		ie.NewNASKeySetIdentifier(0x01),
		ie.NewRaw(0x6f, ie.IEFormatFour, []byte{0xDE, 0xAD}),
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMControlPlaneServiceRequest,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ControlPlaneServiceRequest
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.EsmMessageContainer != nil {
		t.Error("unexpected EsmMessageContainer")
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
