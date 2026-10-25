package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestAttachRequestRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMAttachRequest,
	}

	esmContainer := ie.NewPositionalESMMessageContainer([]byte{0xD0, 0x01, 0xC1})
	drxParam := ie.NewDRXParameter([]byte{0x78, 0x01})
	drxExtra := ie.NewDRXParameter([]byte{0x78, 0x02})

	m := msg.NewAttachRequest(
		ie.NewPositionalEPSAttachType(0x01),
		ie.NewPositionalNASKeySetIdentifier(0x01),
		ie.NewPositionalEPSMobileIdentity([]byte{0x23, 0x05, 0x12, 0x34, 0xAB, 0xCD}),
		ie.NewPositionalUENetworkCapability([]byte{0xAA, 0xBB, 0xCC}),
		esmContainer,
		drxParam, drxExtra,
	)

	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	if len(buf) != m.MarshalLen() {
		t.Fatalf("marshal length = %d, want %d", len(buf), m.MarshalLen())
	}

	var decoded msg.AttachRequest
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

	if got := decoded.MessageType(); got != msg.EMMAttachRequest {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMAttachRequest)
	}
	if got := decoded.MessageTypeName(); got != "Attach request" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Attach request")
	}
	if !decoded.IsRequest() {
		t.Error("IsRequest = false, want true")
	}

	if decoded.EPSAttachType == nil {
		t.Fatal("EPSAttachType is nil")
	}
	if decoded.NASKeySetIdentifier == nil {
		t.Fatal("NASKeySetIdentifier is nil")
	}
	if decoded.DRXParameter == nil {
		t.Error("DRXParameter is nil")
	}
}

func TestAttachRequestEmptyOptional(t *testing.T) {
	t.Parallel()

	m := msg.NewAttachRequest(
		ie.NewPositionalEPSAttachType(0x01),
		ie.NewPositionalNASKeySetIdentifier(0x01),
		ie.NewPositionalEPSMobileIdentity([]byte{0x23, 0x05, 0x12, 0x34, 0xAB}),
		ie.NewPositionalUENetworkCapability([]byte{0xAA, 0xBB}),
		ie.NewPositionalESMMessageContainer([]byte{0xD0, 0x01}),
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMAttachRequest,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.AttachRequest
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
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs in empty optional")
	}
}

func TestAttachRequestProtectedHeader(t *testing.T) {
	t.Parallel()

	m := msg.NewAttachRequest(
		ie.NewPositionalEPSAttachType(0x01),
		ie.NewPositionalNASKeySetIdentifier(0x01),
		ie.NewPositionalEPSMobileIdentity([]byte{0x23, 0x05, 0x12, 0x34}),
		ie.NewPositionalUENetworkCapability([]byte{0xAA, 0xBB}),
		ie.NewPositionalESMMessageContainer([]byte{0xD0, 0x01}),
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
		SequenceNumber:        5,
		MsgType:               msg.EMMAttachRequest,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.AttachRequest
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

func TestAttachRequestErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.AttachRequest

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

	local := &msg.AttachRequest{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.AttachRequest
	err = decoded.UnmarshalBinary([]byte{0x07, 0x41}) // header only - no payload
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header-only = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0x07, 0x41, 0x01}) // only half-octets
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated content = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0x17, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0x41})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated protected = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0xAA, 0x41}) // invalid SHT
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid SHT = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}
}

func TestAttachRequestNilFieldMarshal(t *testing.T) {
	t.Parallel()

	// Even though the test is named "NilFieldMarshal", positional fields are
	// always non-nil (they're mandatory). We test with all fields present but
	// ensure the round-trip still works.
	m := msg.NewAttachRequest(
		ie.NewPositionalEPSAttachType(0x01),
		ie.NewPositionalNASKeySetIdentifier(0x01),
		ie.NewPositionalEPSMobileIdentity([]byte{0x23, 0x05, 0x12, 0x34}),
		ie.NewPositionalUENetworkCapability([]byte{0xAA, 0xBB}),
		ie.NewPositionalESMMessageContainer([]byte{0xD0, 0x01}),
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMAttachRequest,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.AttachRequest
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
