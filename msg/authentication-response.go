package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// AuthenticationResponse represents an AUTHENTICATION RESPONSE message as defined in TS
// 24.301 clause 8.2.8 (Table 8.2.8.1, direction UE to network).
//
// The positional mandatory part encodes the authentication response parameter (RES* / RES)
// as a type-4 LV IE with value length 5–17 octets. No tagged optional IEs are defined per
// TS 24.301 Table 8.2.8.1 — any remaining bytes after parsing the positional IE must be
// empty, otherwise UnmarshalBinary returns an error.
type AuthenticationResponse struct {
	EmmHeader
	AuthenticationResponseParameter *ie.IE   // positional LV type-4 (9.9.3.4), 5–17 value octets
	IEs                             []*ie.IE // trailing unrecognized IEs (reserved for future use)
}

// NewAuthenticationResponse creates an AUTHENTICATION RESPONSE message with the mandatory
// positional authentication response parameter from TS 24.301 clause 8.2.8.
func NewAuthenticationResponse(authRespParam *ie.IE) *AuthenticationResponse {
	if authRespParam == nil {
		panic("AuthenticationResponse positional IEs must not be nil")
	}

	return &AuthenticationResponse{
		AuthenticationResponseParameter: authRespParam,
	}
}

// MarshalLen returns the wire size of the AUTHENTICATION RESPONSE message, or 0 for a nil
// receiver. The result includes header and positional authentication response parameter.
func (m *AuthenticationResponse) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	if m.AuthenticationResponseParameter != nil {
		n += m.AuthenticationResponseParameter.MarshalLenWithoutTag()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the AUTHENTICATION RESPONSE message into b. The nil receiver returns
// a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too short. Header,
// positional authentication response parameter follow TS 24.301 clause 8.2.8 table order.
func (m *AuthenticationResponse) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: AuthenticationResponse is nil", ie.ErrNilIE)
	}

	hlen := m.EmmHeader.MarshalLen()
	if hlen == 0 {
		return fmt.Errorf("%w: EmmHeader MarshalLen returned 0", ie.ErrNilIE)
	}

	if len(b) < hlen {
		return io.ErrShortBuffer
	}

	if err := m.EmmHeader.MarshalTo(b); err != nil {
		return fmt.Errorf("header: %w", err)
	}

	b = b[hlen:]
	n := 0

	if m.AuthenticationResponseParameter != nil {
		if err := m.AuthenticationResponseParameter.MarshalToWithoutTag(b[n:]); err != nil {
			return fmt.Errorf("AuthenticationResponseParameter: %w", err)
		}
		n += m.AuthenticationResponseParameter.MarshalLenWithoutTag()
	}

	for _, i := range m.IEs {
		w, err := writeIE(b, n, i)
		if err != nil {
			return fmt.Errorf("extra IE: %w", err)
		}
		n += w
	}

	return nil
}

// UnmarshalBinary decodes an AUTHENTICATION RESPONSE message from b. The receiver is left
// unchanged on error. Header, positional authentication response parameter via
// ie.ParseWithoutTag. Per TS 24.301 Table 8.2.8.1 there is no tagged optional part —
// any remaining bytes after the positional IE are treated as trailing content and produce
// an io.ErrUnexpectedEOF error.
func (m *AuthenticationResponse) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: AuthenticationResponse is nil", ie.ErrNilIE)
	}

	var tmp AuthenticationResponse

	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	pos := 0

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	// Positional mandatory authentication response parameter (RES / RES*, LV type-4, 5–17 value bytes).
	tmp.AuthenticationResponseParameter, err = ie.ParseWithoutTag(
		pay[pos:], ie.NewPositionalAuthenticationResponseParameter(make([]byte, 5)),
	)
	if err != nil {
		return fmt.Errorf("AuthenticationResponseParameter: %w", err)
	}
	pos += tmp.AuthenticationResponseParameter.MarshalLenWithoutTag()

	// No tagged optional IEs per TS 24.301 Table 8.2.8.1 — any remaining bytes are an error.
	if pos < len(pay) {
		trailing := len(pay) - pos
		return fmt.Errorf(
			"%w: unexpected trailing %d bytes after AuthenticationResponseParameter",
			io.ErrUnexpectedEOF, trailing,
		)
	}

	*m = tmp
	return nil
}

// MessageType returns the AUTHENTICATION RESPONSE message type constant (0x53).
func (m *AuthenticationResponse) MessageType() uint8 {
	return EMMAuthenticationResponse
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *AuthenticationResponse) MessageTypeName() string {
	return "Authentication response"
}

// IsRequest returns true — AUTHENTICATION RESPONSE is sent UE to network.
func (m *AuthenticationResponse) IsRequest() bool {
	return true
}

// AuthenticationResponse implements the Message interface.
var _ Message = (*AuthenticationResponse)(nil)

// Header returns the EMM security header embedded in AuthenticationResponse.
func (m *AuthenticationResponse) Header() *EmmHeader {
	return &m.EmmHeader
}
