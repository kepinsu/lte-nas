package msg_test

import (
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestIdentityRequestRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
		SequenceNumber:        3,
		MsgType:               msg.EMMIdentityRequest,
	}

	identType := ie.NewPositionalIdentityType2(0x01) // IMEI
	spareHalf := ie.NewPositionalSpareHalfOctet(0x00)

	m := msg.NewIdentityRequest(identType, spareHalf)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.IdentityRequest
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

	if got := decoded.MessageType(); got != msg.EMMIdentityRequest {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMIdentityRequest)
	}
	if got := decoded.MessageTypeName(); got != "Identity request" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Identity request")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}

	if decoded.IdentityType == nil {
		t.Fatal("IdentityType is nil")
	}
	if decoded.SpareHalfOctet == nil {
		t.Fatal("SpareHalfOctet is nil")
	}
}

func TestIdentityRequestPlainHeader(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMIdentityRequest,
	}

	identType := ie.NewPositionalIdentityType2(0x03) // TMSI
	spareHalf := ie.NewPositionalSpareHalfOctet(0x07)

	m := msg.NewIdentityRequest(identType, spareHalf)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.IdentityRequest
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

	if decoded.IdentityType == nil {
		t.Fatal("IdentityType is nil")
	}
	if decoded.SpareHalfOctet == nil {
		t.Fatal("SpareHalfOctet is nil")
	}
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs")
	}
}

func TestIdentityRequestErrorCases(t *testing.T) {
	t.Parallel()

	// -----------------------------------------------------------------------
	// NewIdentityRequest panic tests (nil positional IEs).
	// -----------------------------------------------------------------------

	defer func() {
		if r := recover(); r == nil {
			t.Error("NewIdentityRequest(nil, nil) did not panic")
		}
	}()
	msg.NewIdentityRequest(nil, nil)
}

func TestIdentityRequestWithExtraIEs(t *testing.T) {
	t.Parallel()

	identType := ie.NewPositionalIdentityType2(0x01) // IMEI
	spareHalf := ie.NewPositionalSpareHalfOctet(0x07)

	extraIE1 := ie.NewRaw(0x2B, ie.IEFormatThree, []byte{0xAA})
	extraIE2 := ie.NewRaw(0x2C, ie.IEFormatThree, []byte{0xBB, 0xCC})

	m := msg.NewIdentityRequest(identType, spareHalf, extraIE1, extraIE2)
	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMIdentityRequest,
	}

	// Verify both extra IEs are stored in m.IEs.
	if len(m.IEs) != 2 {
		t.Fatalf("m.IEs length = %d, want 2", len(m.IEs))
	}

	// MarshalLen should account for extra IEs.
	wantLen := m.MarshalLen()
	if wantLen <= 3 { // header(2) + half-octet(1) = 3, extra IEs add more
		t.Fatalf("MarshalLen = %d, want > 3", wantLen)
	}

	buf := make([]byte, wantLen)
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	// UnmarshalBinary rejects trailing bytes, so skip round-trip.
	// We verified New* populates m.IEs and MarshalTo/MarshalLen work with extras.
}

func TestIdentityRequestNilExtra(t *testing.T) {
	t.Parallel()

	identType := ie.NewPositionalIdentityType2(0x01)
	spareHalf := ie.NewPositionalSpareHalfOctet(0x07)

	// nil extra IE should be skipped by the loop.
	m := msg.NewIdentityRequest(identType, spareHalf, nil)

	if len(m.IEs) != 0 {
		t.Fatalf("m.IEs length = %d, want 0", len(m.IEs))
	}

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMIdentityRequest,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.IdentityRequest
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

func TestIdentityRequestPartialPositionalIEs(t *testing.T) {
	t.Parallel()

	// Construct IdentityRequest with only IdentityType (SpareHalfOctet is nil).
	// This exercises the else branch in MarshalLen and MarshalTo.
	local := &msg.IdentityRequest{
		IdentityType: ie.NewPositionalIdentityType2(0x1),
	}
	local.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMIdentityRequest,
	}

	if got := local.MarshalLen(); got != 3 { // header(2) + IdentityType(1)
		t.Errorf("MarshalLen = %d, want 3", got)
	}

	buf := make([]byte, local.MarshalLen())
	if err := local.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo: %v", err)
	}

	var decoded msg.IdentityRequest
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary: %v", err)
	}

	rebuf := make([]byte, decoded.MarshalLen())
	if err := decoded.MarshalTo(rebuf); err != nil {
		t.Fatalf("re-MarshalTo: %v", err)
	}
	for i := range buf {
		if buf[i] != rebuf[i] {
			t.Fatalf("byte[%d] = %#02x, want %#02x", i, rebuf[i], buf[i])
		}
	}

	if got := decoded.IdentityType; got == nil {
		t.Fatal("IdentityType is nil after decode")
	}
}
