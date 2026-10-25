package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// AuthenticationReject represents an AUTHENTICATION REJECT message as defined in TS 24.301
// clause 8.2.6 (Table 8.2.6.1, direction network to UE).
//
// This message has no body IEs beyond the header and message type. The plain EMM header
// (SHT=0x00, 2 octets) plus message type (1 octet) produce a 3-octet wire format for
// unprotect messages. Per TS 24.301 Table 8.2.6.1 no optional IEs are defined — any
// bytes after the header and message type produce an error.
type AuthenticationReject struct {
	EmmHeader
	IEs []*ie.IE // trailing unrecognized IEs (reserved for future use)
}

// NewAuthenticationReject creates an empty AUTHENTICATION REJECT message with no body IEs.
func NewAuthenticationReject(ies ...*ie.IE) *AuthenticationReject {
	m := &AuthenticationReject{}

	for _, i := range ies {
		if i == nil {
			continue
		}
		m.IEs = append(m.IEs, i)
	}

	return m
}

// MarshalLen returns the wire size of the AUTHENTICATION REJECT message, or 0 for a nil
// receiver. The result includes header only — no body IEs are present per TS 24.301
// Table 8.2.6.1.
func (m *AuthenticationReject) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the AUTHENTICATION REJECT message into b. The nil receiver returns
// a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too short. Header only
// per TS 24.301 clause 8.2.6 table order.
func (m *AuthenticationReject) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: AuthenticationReject is nil", ie.ErrNilIE)
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

	for _, i := range m.IEs {
		w, err := writeIE(b, n, i)
		if err != nil {
			return fmt.Errorf("extra IE: %w", err)
		}
		n += w
	}

	return nil
}

// UnmarshalBinary decodes an AUTHENTICATION REJECT message from b. The receiver is left
// unchanged on error. Header only — per TS 24.301 Table 8.2.6.1 there are no body IEs.
// Any remaining bytes after the header produce an error.
func (m *AuthenticationReject) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: AuthenticationReject is nil", ie.ErrNilIE)
	}

	var tmp AuthenticationReject

	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	// Per TS 24.301 Table 8.2.6.1 there are no body IEs — any remaining bytes
	// are unexpected trailing content and produce an error.
	if len(pay) > 0 {
		return fmt.Errorf("%w: unexpected %d bytes after header in AuthenticationReject", io.ErrUnexpectedEOF, len(pay))
	}

	*m = tmp
	return nil
}

// MessageType returns the AUTHENTICATION REJECT message type constant (0x54).
func (m *AuthenticationReject) MessageType() uint8 {
	return EMMAuthenticationReject
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *AuthenticationReject) MessageTypeName() string {
	return "Authentication reject"
}

// IsRequest returns false — AUTHENTICATION REJECT is sent network to UE.
func (m *AuthenticationReject) IsRequest() bool {
	return false
}

// AuthenticationReject implements the Message interface.
var _ Message = (*AuthenticationReject)(nil)

// Header returns the EMM security header embedded in AuthenticationReject.
func (m *AuthenticationReject) Header() *EmmHeader {
	return &m.EmmHeader
}
