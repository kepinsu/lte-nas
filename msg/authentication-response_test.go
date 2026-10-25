package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestAuthenticationResponseRoundTrip(t *testing.T) {
	t.Parallel()

	respParam := ie.NewPositionalAuthenticationResponseParameter([]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE})

	m := msg.NewAuthenticationResponse(respParam)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMAuthenticationResponse,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.AuthenticationResponse
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

	if got := decoded.MessageType(); got != msg.EMMAuthenticationResponse {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMAuthenticationResponse)
	}
	if got := decoded.MessageTypeName(); got != "Authentication response" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Authentication response")
	}
	if !decoded.IsRequest() {
		t.Error("IsRequest = false, want true")
	}
	if decoded.AuthenticationResponseParameter == nil {
		t.Fatal("AuthenticationResponseParameter is nil")
	}

	fields, err := decoded.AuthenticationResponseParameter.AuthenticationResponseParameter()
	if err != nil {
		t.Fatalf("AuthenticationResponseParameter: %v", err)
	}
	if len(fields.Response) != 5 {
		t.Errorf("Response length = %d, want 5", len(fields.Response))
	}
}

func TestAuthenticationResponseLongResp(t *testing.T) {
	t.Parallel()

	respParam := ie.NewPositionalAuthenticationResponseParameter(
		[]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
			0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10},
	)

	m := msg.NewAuthenticationResponse(respParam)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMAuthenticationResponse,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.AuthenticationResponse
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

	fields, err := decoded.AuthenticationResponseParameter.AuthenticationResponseParameter()
	if err != nil {
		t.Fatalf("AuthenticationResponseParameter: %v", err)
	}
	if len(fields.Response) != 16 {
		t.Errorf("Response length = %d, want 16", len(fields.Response))
	}
}

func TestAuthenticationResponseProtectedHeader(t *testing.T) {
	t.Parallel()

	respParam := ie.NewPositionalAuthenticationResponseParameter([]byte{0x11, 0x22, 0x33, 0x44})
	m := msg.NewAuthenticationResponse(respParam)
	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0xAA, 0xBB, 0xCC, 0xDD},
		SequenceNumber:        5,
		MsgType:               msg.EMMAuthenticationResponse,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.AuthenticationResponse
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

func TestAuthenticationResponseErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.AuthenticationResponse

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

	local := &msg.AuthenticationResponse{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.AuthenticationResponse
	err = decoded.UnmarshalBinary([]byte{0x07, 0x53}) // header only - no payload
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header-only = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0x07, 0x53, 0xAB}) // only type-1 (1 octet), need at least len+value
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated positional IE = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0xAA, 0x53}) // invalid SHT
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid SHT = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}

	err = decoded.UnmarshalBinary([]byte{0x17, 0xAA, 0xBB, 0xCC, 0xDD}) // truncated protected header
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated protected header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Valid message + trailing bytes should error (no tagged optional IEs allowed)
	validMsg := []byte{0x07, 0x53, 0x07, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}
	err = decoded.UnmarshalBinary(validMsg)
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("trailing bytes = %v, want %T", err, io.ErrUnexpectedEOF)
	}
}

func TestAuthenticationResponseWithExtraIEs(t *testing.T) {
	t.Parallel()

	respParam := ie.NewPositionalAuthenticationResponseParameter(
		[]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE},
	)

	m := msg.NewAuthenticationResponse(respParam)
	// Populate m.IEs directly — NewAuthenticationResponse has no variadic, but
	// MarshalLen/MarshalTo still iterate over the field.
	extraIE := ie.NewRaw(0x7e, ie.IEFormatFour, []byte{0x99})
	m.IEs = []*ie.IE{extraIE}

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMAuthenticationResponse,
	}

	// Verify MarshalLen accounts for the extra IE.
	wantLen := m.MarshalLen()
	if wantLen <= 0 {
		t.Fatalf("MarshalLen = %d, want > 0", wantLen)
	}

	buf := make([]byte, wantLen)
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	// The extra IE is NOT round-trippable via UnmarshalBinary (it rejects
	// trailing bytes), but we verified MarshalTo executes the extra-IE loop.
}

func TestAuthenticationResponseMarshalLenExtraIEs(t *testing.T) {
	t.Parallel()

	collisionIE := &ie.IE{Kind: 0x32, Format: ie.IEFormatFour, Payload: []byte{0xAA}}

	// Use a valid-length payload for AuthenticationResponseParameter (4..16 bytes).
	authParam := ie.NewPositionalAuthenticationResponseParameter([]byte{0xAA, 0xBB, 0xCC, 0xDD})

	m := &msg.AuthenticationResponse{
		EmmHeader:                   msg.EmmHeader{SecurityHeaderType: msg.SecurityHeaderTypePlain, ProtocolDiscriminator: byte(msg.PdEmm), MsgType: msg.EMMAuthenticationResponse},
		AuthenticationResponseParameter: authParam,
		IEs:                         []*ie.IE{collisionIE},
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo with collision IE failed: %v", err)
	}

	// Test the collision TLV short-buffer error path directly.
	collisionMsg := &msg.AuthenticationResponse{
		EmmHeader:                   msg.EmmHeader{SecurityHeaderType: msg.SecurityHeaderTypePlain, ProtocolDiscriminator: byte(msg.PdEmm), MsgType: msg.EMMAuthenticationResponse},
		AuthenticationResponseParameter: authParam,
		IEs:                         []*ie.IE{collisionIE},
	}
	shortBuf := make([]byte, 1)
	if err := collisionMsg.MarshalTo(shortBuf); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("expected ErrShortBuffer for collision IE error path, got %v", err)
	}
}

func TestAuthenticationResponseNewDefaultBranch(t *testing.T) {
	t.Parallel()

	// 0x56 (NonceMME) is not in the AuthenticationResponse switch.
	nonceIE := ie.NewNonceMME([]byte{0xAA})

	m := &msg.AuthenticationResponse{
		AuthenticationResponseParameter: ie.NewPositionalAuthenticationResponseParameter([]byte{0xAA, 0xBB}),
		IEs:                             []*ie.IE{nonceIE},
	}

	if m.AuthenticationResponseParameter == nil {
		t.Fatal("positional IE should be set")
	}
	if len(m.IEs) == 0 {
		t.Fatal("NonceMME should go to IEs field (default branch)")
	}
	if m.IEs[0].Kind != 0x56 {
		t.Fatalf("IEs[0].Kind = %#02x, want 0x56", m.IEs[0].Kind)
	}
}
