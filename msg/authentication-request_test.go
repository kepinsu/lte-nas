package msg_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestAuthenticationRequestRoundTrip(t *testing.T) {
	t.Parallel()

	// Full IE set: all four positional mandatory fields.
	nkid := ie.NewPositionalNASKeySetIdentifier(0x1)
	spare := ie.NewPositionalSpareHalfOctet(0x2)
	rand := ie.NewPositionalAuthenticationParameterRAND(make([]byte, 16))
	copy(rand.Payload, []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA})
	autn := ie.NewPositionalAuthenticationParameterAUTN(make([]byte, 16))
	copy(autn.Payload, []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
		0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10})

	m := msg.NewAuthenticationRequest(nkid, spare, rand, autn)
	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMAuthenticationRequest,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	if len(buf) != m.MarshalLen() {
		t.Fatalf("marshal length = %d, want %d", len(buf), m.MarshalLen())
	}

	var decoded msg.AuthenticationRequest
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

	if got := decoded.MessageType(); got != msg.EMMAuthenticationRequest {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMAuthenticationRequest)
	}
	if got := decoded.MessageTypeName(); got != "Authentication request" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Authentication request")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}

	if decoded.NASKeySetIdentifier == nil {
		t.Fatal("NASKeySetIdentifier is nil")
	}
	if decoded.AuthenticationParameterRAND == nil {
		t.Fatal("AuthenticationParameterRAND is nil")
	}
	if decoded.AuthenticationParameterAUTN == nil {
		t.Fatal("AuthenticationParameterAUTN is nil")
	}
}

func TestAuthenticationRequestProtectedHeader(t *testing.T) {
	t.Parallel()

	nkid := ie.NewPositionalNASKeySetIdentifier(0x1)
	spare := ie.NewPositionalSpareHalfOctet(0x2)
	rand := ie.NewPositionalAuthenticationParameterRAND(make([]byte, 16))
	copy(rand.Payload, []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA})
	autn := ie.NewPositionalAuthenticationParameterAUTN(make([]byte, 16))
	copy(autn.Payload, []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
		0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10})

	m := msg.NewAuthenticationRequest(nkid, spare, rand, autn)
	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
		SequenceNumber:        5,
		MsgType:               msg.EMMAuthenticationRequest,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.AuthenticationRequest
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
	if decoded.EmmHeader.SequenceNumber != m.EmmHeader.SequenceNumber {
		t.Errorf("SeqNum = %d, want %d", decoded.EmmHeader.SequenceNumber, m.EmmHeader.SequenceNumber)
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

func TestAuthenticationRequestParseHelper(t *testing.T) {
	t.Parallel()

	nkid := ie.NewPositionalNASKeySetIdentifier(0x1)
	spare := ie.NewPositionalSpareHalfOctet(0x2)
	rand := ie.NewPositionalAuthenticationParameterRAND(make([]byte, 16))
	copy(rand.Payload, []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA})
	autn := ie.NewPositionalAuthenticationParameterAUTN(make([]byte, 16))
	copy(autn.Payload, []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
		0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10})

	m := msg.NewAuthenticationRequest(nkid, spare, rand, autn)
	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMAuthenticationRequest,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.AuthenticationRequest
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	// Verify field values.
	nkidVal, err := decoded.NASKeySetIdentifier.ValueAsUint8()
	if err != nil {
		t.Fatalf("NASKeySetIdentifier.ValueAsUint8: %v", err)
	}
	if nkidVal != 0x1 {
		t.Errorf("NASKeySetIdentifier = %#02x, want 0x01", nkidVal)
	}

	randFields, err := decoded.AuthenticationParameterRAND.AuthenticationParameterRAND()
	if err != nil {
		t.Fatalf("AuthenticationParameterRAND: %v", err)
	}
	if !bytes.Equal(randFields.RandomChallenge[:], []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA}) {
		t.Error("RAND value mismatch")
	}

	autnFields, err := decoded.AuthenticationParameterAUTN.AuthenticationParameterAUTN()
	if err != nil {
		t.Fatalf("AuthenticationParameterAUTN: %v", err)
	}
	if !bytes.Equal(autnFields.MessageAuthenticationCode[:], []byte{0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10}) {
		t.Error("AUTN value mismatch")
	}
}

func TestAuthenticationRequestErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.AuthenticationRequest

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

	local := &msg.AuthenticationRequest{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.AuthenticationRequest
	err = decoded.UnmarshalBinary([]byte{0x07, 0x52}) // header only - no payload
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header-only = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0x07, 0x52, 0xAB}) // only half-octets
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated content = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0x07, 0x52, 0xAB,
		0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22,
		0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA}) // no AUTN
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated auth request = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0x07, 0x52, 0xAB,
		0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22,
		0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F}) // truncated AUTN (17 bytes: len+15)
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated AUTN = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0xAA, 0x52}) // invalid SHT and PD
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid SHT = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}
}

func TestAuthenticationRequestExtraIEs(t *testing.T) {
	t.Parallel()

	nkid := ie.NewPositionalNASKeySetIdentifier(0x1)
	spare := ie.NewPositionalSpareHalfOctet(0x2)
	rand := ie.NewPositionalAuthenticationParameterRAND(make([]byte, 16))
	autn := ie.NewPositionalAuthenticationParameterAUTN(make([]byte, 16))

	// Extra IE that NewAuthenticationRequest appends to m.IEs.
	extra := ie.NewRaw(0x7e, ie.IEFormatFour, []byte{0x99})

	m := msg.NewAuthenticationRequest(nkid, spare, rand, autn, extra)
	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMAuthenticationRequest,
	}

	if len(m.IEs) != 1 {
		t.Fatalf("m.IEs length = %d, want 1", len(m.IEs))
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.AuthenticationRequest
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.AuthenticationParameterRAND == nil {
		t.Fatal("decoded AuthenticationParameterRAND is nil")
	}

	// Test nil in variadic (loop should skip).
	m2 := msg.NewAuthenticationRequest(nkid, spare, rand, autn, nil)
	if len(m2.IEs) != 0 {
		t.Fatalf("m2 IEs length = %d, want 0", len(m2.IEs))
	}
}
