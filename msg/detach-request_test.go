package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestDetachRequestRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMDetachRequest,
	}

	esmContainer := ie.NewESMMessageContainer([]byte{0xD0, 0x01, 0xC1})
	ueNC := ie.NewUENetworkCapability([]byte{0xAA, 0xBB})
	ueNCExtra := ie.NewUENetworkCapability([]byte{0xCC, 0xDD})

	m := msg.NewDetachRequest(
		ie.NewPositionalDetachType(0x01),
		ie.NewPositionalNASKeySetIdentifier(0x01),
		ie.NewPositionalEPSMobileIdentity([]byte{0x23, 0x05, 0x12, 0x34, 0xAB, 0xCD}),
		esmContainer, ueNC, ueNCExtra,
	)

	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	if len(buf) != m.MarshalLen() {
		t.Fatalf("marshal length = %d, want %d", len(buf), m.MarshalLen())
	}

	var decoded msg.DetachRequest
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

	if got := decoded.MessageType(); got != msg.EMMDetachRequest {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMDetachRequest)
	}
	if got := decoded.MessageTypeName(); got != "Detach request" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Detach request")
	}
	if !decoded.IsRequest() {
		t.Error("IsRequest = false, want true")
	}

	if decoded.DetachType == nil {
		t.Fatal("DetachType is nil")
	}
	if decoded.NASKeySetIdentifier == nil {
		t.Fatal("NASKeySetIdentifier is nil")
	}
	if decoded.EPSMobileIdentity == nil {
		t.Fatal("EPSMobileIdentity is nil")
	}
	if decoded.ESMMessageContainer == nil {
		t.Error("ESMMessageContainer is nil")
	}
	if decoded.UENetworkCapability == nil {
		t.Error("UENetworkCapability is nil")
	}
	if len(decoded.IEs) == 0 {
		t.Error("unexpected empty IEs (UE NC duplicate should be there)")
	} else {
		if decoded.IEs[0].Kind != ie.UENetworkCapability {
			t.Errorf("IEs[0] kind = %v, want UENetworkCapability", decoded.IEs[0].Kind)
		}
	}
}

func TestDetachRequestEmptyOptional(t *testing.T) {
	t.Parallel()

	m := msg.NewDetachRequest(
		ie.NewPositionalDetachType(0x00),
		ie.NewPositionalNASKeySetIdentifier(0x01),
		ie.NewPositionalEPSMobileIdentity([]byte{0x23, 0x05, 0x12, 0x34, 0xAB}),
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMDetachRequest,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.DetachRequest
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

	if decoded.ESMMessageContainer != nil {
		t.Error("unexpected ESMMessageContainer in empty optional")
	}
	if decoded.UENetworkCapability != nil {
		t.Error("unexpected UENetworkCapability in empty optional")
	}
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs in empty optional")
	}
}

func TestDetachRequestProtectedHeader(t *testing.T) {
	t.Parallel()

	m := msg.NewDetachRequest(
		ie.NewPositionalDetachType(0x01),
		ie.NewPositionalNASKeySetIdentifier(0x02),
		ie.NewPositionalEPSMobileIdentity([]byte{0x23, 0x05, 0x12, 0x34}),
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
		SequenceNumber:        5,
		MsgType:               msg.EMMDetachRequest,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.DetachRequest
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.EmmHeader.SecurityHeaderType != msg.SecurityHeaderTypeIntegrityProtected {
		t.Errorf("SHT = %#02x, want %#02x", decoded.EmmHeader.SecurityHeaderType, msg.SecurityHeaderTypeIntegrityProtected)
	}
	for i := range 4 {
		if decoded.EmmHeader.Mac[i] != m.EmmHeader.Mac[i] {
			t.Errorf("Mac[%d] = %#02x, want %#02x", i, decoded.EmmHeader.Mac[i], m.EmmHeader.Mac[i])
		}
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

func TestDetachRequestErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.DetachRequest

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

	local := &msg.DetachRequest{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.DetachRequest
	err = decoded.UnmarshalBinary([]byte{0x07, 0x45}) // header only (plain, PD=7) - no payload
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header-only = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0x17, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0x45}) // truncated protected
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated protected = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0x1E, 0x45}) // invalid PD (E instead of 7)
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid PD = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}

	err = decoded.UnmarshalBinary([]byte{0x7E, 0xAA}) // SHT=7 (invalid for plain)
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid SHT = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}
}
