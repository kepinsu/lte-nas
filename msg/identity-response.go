package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// IdentityResponse represents an IDENTITY RESPONSE message as defined in TS
// 24.301 clause 8.2.19 (Table 8.2.19.1, direction UE to network).
//
// The positional mandatory part encodes mobile identity as a plain LV (no IEI).
// Plain EMM header (SHT=0x00, 2 octets). No tagged optional IEs per TS 24.301
// Table 8.2.19.1.
type IdentityResponse struct {
	EmmHeader
	MobileIdentity *ie.IE // positional LV type-4 (9.9.2.3)
	IEs            []*ie.IE
}

// NewIdentityResponse creates an IDENTITY RESPONSE message with the mandatory
// positional IE from TS 24.301 clause 8.2.19, followed by optional tagged IEs
// classified by their Kind. (No optional IEs per spec, but the variadic
// parameter is provided for API consistency.)
func NewIdentityResponse(mobileIdentity *ie.IE, ies ...*ie.IE) *IdentityResponse {
	if mobileIdentity == nil {
		panic("IdentityResponse positional IEs must not be nil")
	}

	m := &IdentityResponse{
		MobileIdentity: mobileIdentity,
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		m.IEs = append(m.IEs, i)
	}

	return m
}

// MarshalLen returns the wire size of the IDENTITY RESPONSE message, or 0 for a
// nil receiver. The result includes header and positional mobile identity (length+value).
func (m *IdentityResponse) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	if m.MobileIdentity != nil {
		n += m.MobileIdentity.MarshalLenWithoutTag()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the IDENTITY RESPONSE message into b. The nil receiver
// returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too
// short. Header, positional mobile identity per TS 24.301 clause 8.2.19 table
// order.
func (m *IdentityResponse) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: IdentityResponse is nil", ie.ErrNilIE)
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

	if m.MobileIdentity != nil {
		if err := m.MobileIdentity.MarshalToWithoutTag(b[n:]); err != nil {
			return fmt.Errorf("MobileIdentity: %w", err)
		}
		n += m.MobileIdentity.MarshalLenWithoutTag()
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

// UnmarshalBinary decodes an IDENTITY RESPONSE message from b. The receiver is
// left unchanged on error. Header, positional mobile identity via
// ParseWithoutTag. No tagged optional IEs per TS 24.301 Table 8.2.19.1.
func (m *IdentityResponse) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: IdentityResponse is nil", ie.ErrNilIE)
	}

	var tmp IdentityResponse

	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	pos := 0

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	tmp.MobileIdentity, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalMobileIdentity([]byte{0, 0}))
	if err != nil {
		return fmt.Errorf("MobileIdentity: %w", err)
	}

	*m = tmp
	return nil
}

// MessageType returns the IDENTITY RESPONSE message type constant (0x56).
func (m *IdentityResponse) MessageType() uint8 {
	return EMMIdentityResponse
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *IdentityResponse) MessageTypeName() string {
	return "Identity response"
}

// IsRequest returns true — IDENTITY RESPONSE is sent UE to network.
func (m *IdentityResponse) IsRequest() bool {
	return true
}

// IdentityResponse implements the Message interface.
var _ Message = (*IdentityResponse)(nil)

// Header returns the EMM security header embedded in IdentityResponse.
func (m *IdentityResponse) Header() *EmmHeader {
	return &m.EmmHeader
}
