package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestExtendedServiceRequestRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMExtendedServiceRequest,
	}

	m := msg.NewExtendedServiceRequest(
		ie.NewPositionalServiceType(0x02),
		ie.NewNASKeySetIdentifier(0x01),
		ie.NewMobileIdentity([]byte{0xf0, 0x12, 0x34, 0xab, 0xcd, 0xef}),
		ie.NewCSFBResponse(0x01),
		ie.NewEPSBearerContextStatus([]byte{0x01, 0xff}),
		ie.NewDeviceProperties(0x01),
		ie.NewUERequestType(0x01),
		ie.NewPagingRestriction([]byte{0xff}),
	)

	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ExtendedServiceRequest
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

	if got := decoded.MessageType(); got != msg.EMMExtendedServiceRequest {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMExtendedServiceRequest)
	}
	if got := decoded.MessageTypeName(); got != "Extended service request" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Extended service request")
	}
	if !decoded.IsRequest() {
		t.Error("IsRequest = false, want true")
	}
	if decoded.ServiceType == nil {
		t.Fatal("ServiceType is nil")
	}
	if decoded.NASKeySetIdentifier == nil {
		t.Fatal("NASKeySetIdentifier is nil")
	}
	if decoded.MobileIdentity == nil {
		t.Fatal("MobileIdentity is nil")
	}
	if decoded.CSFBResponse == nil {
		t.Error("CSFBResponse is nil")
	}
	if decoded.EPSBearerContextStatus == nil {
		t.Error("EpsBearerContextStatus is nil")
	}
}

func TestExtendedServiceRequestEmptyOptional(t *testing.T) {
	t.Parallel()

	m := msg.NewExtendedServiceRequest(
		ie.NewPositionalServiceType(0x02),
		ie.NewNASKeySetIdentifier(0x01),
		ie.NewMobileIdentity([]byte{0xf0, 0x12, 0x34, 0xab, 0xcd}),
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMExtendedServiceRequest,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ExtendedServiceRequest
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

	if decoded.CSFBResponse != nil {
		t.Error("unexpected CSFBResponse in empty optional")
	}
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs in empty optional")
	}
}

func TestExtendedServiceRequestProtectedHeader(t *testing.T) {
	t.Parallel()

	m := msg.NewExtendedServiceRequest(
		ie.NewPositionalServiceType(0x02),
		ie.NewNASKeySetIdentifier(0x01),
		ie.NewMobileIdentity([]byte{0xf0, 0x12, 0x34, 0xab}),
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
		SequenceNumber:        5,
		MsgType:               msg.EMMExtendedServiceRequest,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ExtendedServiceRequest
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

func TestExtendedServiceRequestErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.ExtendedServiceRequest

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

	local := &msg.ExtendedServiceRequest{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.ExtendedServiceRequest
	err = decoded.UnmarshalBinary([]byte{0x07, 0x4C}) // header only - no payload

	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header-only = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0x17, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0x4C}) // truncated protected
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated protected = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0xAA, 0x4C}) // invalid SHT
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid SHT = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}
}

func TestExtendedServiceRequestExtraIes(t *testing.T) {
	t.Parallel()

	m := msg.NewExtendedServiceRequest(
		ie.NewPositionalServiceType(0x02),
		ie.NewNASKeySetIdentifier(0x01),
		ie.NewMobileIdentity([]byte{0xf0, 0x12, 0x34}),
		ie.NewCSFBResponse(0x01),
		ie.NewCSFBResponse(0x02), // duplicate
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMExtendedServiceRequest,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ExtendedServiceRequest
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.CSFBResponse == nil {
		t.Fatal("first CSFBResponse is nil")
	}

	if got := len(decoded.IEs); got != 1 {
		t.Fatalf("IEs length = %d, want 1", got)
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
