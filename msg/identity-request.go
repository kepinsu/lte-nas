package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// IdentityRequest represents an IDENTITY REQUEST message as defined in TS 24.301
// clause 8.2.18 (Table 8.2.18.1, direction network to UE).
//
// The positional mandatory part encodes identity type and spare half octet as
// packed half-octets (wire byte), matching the spec table format. Integrity-
// protected header (SHT 0x01..0x05, 7 octets). No tagged optional IEs per
// TS 24.301 Table 8.2.18.1.
type IdentityRequest struct {
	EmmHeader
	IdentityType   *ie.IE // positional type-1 (9.9.3.17), half-octet with SpareHalfOctet
	SpareHalfOctet *ie.IE // positional type-1 (9.9.2.9), half-octet with IdentityType
	IEs            []*ie.IE
}

// NewIdentityRequest creates an IDENTITY REQUEST message with the mandatory
// positional half-octets from TS 24.301 clause 8.2.18, followed by optional
// tagged IEs classified by their Kind. (No optional IEs per spec, but the
// variadic parameter is provided for API consistency.)
func NewIdentityRequest(identityType, spareHalfOctet *ie.IE, ies ...*ie.IE) *IdentityRequest {
	if identityType == nil || spareHalfOctet == nil {
		panic("IdentityRequest positional IEs must not be nil")
	}

	m := &IdentityRequest{
		IdentityType:   identityType,
		SpareHalfOctet: spareHalfOctet,
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		m.IEs = append(m.IEs, i)
	}

	return m
}

// MarshalLen returns the wire size of the IDENTITY REQUEST message, or 0 for a
// nil receiver. The result includes header and packed half-octets (1 byte).
func (m *IdentityRequest) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	// Packed half-octets (identity type low nibble, spare high nibble).
	if m.IdentityType != nil && m.SpareHalfOctet != nil {
		n++ // packed half-octet
	} else {
		if m.IdentityType != nil {
			n++
		}
		if m.SpareHalfOctet != nil {
			n++
		}
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the IDENTITY REQUEST message into b. The nil receiver
// returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too
// short. Header, packed half-octets in spec-table order per TS 24.301 clause
// 8.2.18 table order.
func (m *IdentityRequest) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: IdentityRequest is nil", ie.ErrNilIE)
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

	// Packed half-octets (identity type low nibble, spare high nibble).
	if m.IdentityType != nil && m.SpareHalfOctet != nil {
		hb, err := ie.MarshalHalfOctets(m.IdentityType, m.SpareHalfOctet)
		if err != nil {
			return fmt.Errorf("half-octets: %w", err)
		}
		if len(b) < n+len(hb) {
			return io.ErrShortBuffer
		}
		copy(b[n:n+len(hb)], hb)
		n += len(hb)
	} else {
		if m.IdentityType != nil {
			if err := m.IdentityType.MarshalToWithoutTag(b[n:]); err != nil {
				return fmt.Errorf("IdentityType: %w", err)
			}
			n++
		}
		if m.SpareHalfOctet != nil {
			if err := m.SpareHalfOctet.MarshalToWithoutTag(b[n:]); err != nil {
				return fmt.Errorf("SpareHalfOctet: %w", err)
			}
			n++
		}
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

// UnmarshalBinary decodes an IDENTITY REQUEST message from b. The receiver is
// left unchanged on error. Header, half-octets via ParseHalfOctets. No tagged
// optional IEs per TS 24.301 Table 8.2.18.1.
func (m *IdentityRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: IdentityRequest is nil", ie.ErrNilIE)
	}

	var tmp IdentityRequest

	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	pos := 0

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	// Positional mandatory: half-octets (IdentityType low nibble,
	// SpareHalfOctet high nibble) packed in one wire byte.
	lProto := ie.NewPositionalIdentityType2(0)
	hProto := ie.NewPositionalSpareHalfOctet(0)
	lIE, hIE, err := ie.ParseHalfOctets(pay[pos:], lProto, hProto)
	if err != nil {
		return fmt.Errorf("half-octets: %w", err)
	}
	tmp.IdentityType = lIE
	tmp.SpareHalfOctet = hIE

	*m = tmp
	return nil
}

// MessageType returns the IDENTITY REQUEST message type constant (0x55).
func (m *IdentityRequest) MessageType() uint8 {
	return EMMIdentityRequest
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *IdentityRequest) MessageTypeName() string {
	return "Identity request"
}

// IsRequest returns false — IDENTITY REQUEST is sent network to UE.
func (m *IdentityRequest) IsRequest() bool {
	return false
}

// IdentityRequest implements the Message interface.
var _ Message = (*IdentityRequest)(nil)

// Header returns the EMM security header embedded in IdentityRequest.
func (m *IdentityRequest) Header() *EmmHeader {
	return &m.EmmHeader
}
