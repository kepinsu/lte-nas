package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestIdentityResponseRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMIdentityResponse,
	}

	mobileId := ie.NewPositionalMobileIdentity([]byte{0x05, 0x12, 0x34, 0xAB, 0xCD})

	m := msg.NewIdentityResponse(mobileId)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.IdentityResponse
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

	if got := decoded.MessageType(); got != msg.EMMIdentityResponse {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMIdentityResponse)
	}
	if got := decoded.MessageTypeName(); got != "Identity response" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Identity response")
	}
	if !decoded.IsRequest() {
		t.Error("IsRequest = false, want true")
	}

	if decoded.MobileIdentity == nil {
		t.Fatal("MobileIdentity is nil")
	}
}

func TestIdentityResponseProtectedHeader(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
		SequenceNumber:        2,
		MsgType:               msg.EMMIdentityResponse,
	}

	mobileId := ie.NewPositionalMobileIdentity([]byte{0x03, 0xAB, 0xCD, 0xEF})

	m := msg.NewIdentityResponse(mobileId)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.IdentityResponse
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

	if decoded.MobileIdentity == nil {
		t.Fatal("MobileIdentity is nil")
	}
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs")
	}
}

func TestIdentityResponseWithExtraIEs(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMIdentityResponse,
	}

	mobileId := ie.NewPositionalMobileIdentity([]byte{0x05, 0x12, 0x34, 0xAB, 0xCD})
	// Extra tagged IEs that NewIdentityResponse will append to m.IEs.
	extra1 := ie.NewRaw(0x2B, ie.IEFormatThree, []byte{0xAA})
	extra2 := ie.NewRaw(0x2C, ie.IEFormatThree, []byte{0xBB, 0xCC})

	m := msg.NewIdentityResponse(mobileId, extra1, extra2)
	m.EmmHeader = h

	// New* populates m.IEs even though UnmarshalBinary does not parse them back.
	if len(m.IEs) != 2 {
		t.Fatalf("m.IEs length = %d, want 2", len(m.IEs))
	}
	if m.IEs[0].Kind != 0x2B {
		t.Errorf("IE[0] Kind = 0x%02x, want 0x2B", m.IEs[0].Kind)
	}

	// Marshal to verify the extra-IE loops in MarshalTo and MarshalLen execute.
	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.IdentityResponse
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.MobileIdentity == nil {
		t.Fatal("MobileIdentity is nil")
	}
}

func TestIdentityResponseErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.IdentityResponse
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

	local := &msg.IdentityResponse{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.IdentityResponse
	err = decoded.UnmarshalBinary([]byte{0xAA, 0x56}) // invalid SHT
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid SHT = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}

	err = decoded.UnmarshalBinary([]byte{0xAA, 0xAA}) // valid header, invalid PD
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid PD = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}

	err = decoded.UnmarshalBinary([]byte{0x17, 0x56}) // plain header, no payload
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated content = %v, want %T", err, io.ErrUnexpectedEOF)
	}
}
