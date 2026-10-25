package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// AuthenticationRequest represents an AUTHENTICATION REQUEST message as defined in TS 24.301
// clause 8.2.7 (Table 8.2.7.1, direction network to UE).
//
// The positional mandatory part encodes NAS key set identifier and spare half octet as a type-1
// pair (packed in a single wire byte), followed by authentication parameter RAND (16 bytes, V)
// and authentication parameter AUTN (LV type-4, 16 value bytes). Per TS 24.301 V18.9.0 Table
// 8.2.7.1 these are all positional — no tagged optional IEs are defined.
type AuthenticationRequest struct {
	EmmHeader
	NASKeySetIdentifier         *ie.IE   // positional type-1 (9.9.3.21), half-octet with SpareHalfOctet
	SpareHalfOctet              *ie.IE   // positional type-1 (9.9.2.9), half-octet with NASKeySetIdentifier
	AuthenticationParameterRAND *ie.IE   // positional V type-3 (9.9.3.3), 16 octets
	AuthenticationParameterAUTN *ie.IE   // positional LV type-4 (9.9.3.2), 16 value octets
	IEs                         []*ie.IE // trailing unrecognized IEs
}

// NewAuthenticationRequest creates an AUTHENTICATION REQUEST message with the mandatory
// positional IEs from TS 24.301 clause 8.2.7 (Table 8.2.7.1).
func NewAuthenticationRequest(nasKeySetID, spareHalfOctet, rand, autn *ie.IE, ies ...*ie.IE) *AuthenticationRequest {
	if nasKeySetID == nil || spareHalfOctet == nil || rand == nil || autn == nil {
		panic("AuthenticationRequest positional IEs must not be nil")
	}

	m := &AuthenticationRequest{
		NASKeySetIdentifier:         nasKeySetID,
		SpareHalfOctet:              spareHalfOctet,
		AuthenticationParameterRAND: rand,
		AuthenticationParameterAUTN: autn,
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		m.IEs = append(m.IEs, i)
	}

	return m
}

// MarshalLen returns the wire size of the AUTHENTICATION REQUEST message, or 0 for a nil
// receiver. The result includes header, packed half-octets (1 byte), positional IEs and
// extras.
func (m *AuthenticationRequest) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	// Half-octets: packed (1 byte) or individual.
	if m.NASKeySetIdentifier != nil && m.SpareHalfOctet != nil {
		n++
	} else {
		if m.NASKeySetIdentifier != nil {
			n++
		}
		if m.SpareHalfOctet != nil {
			n++
		}
	}

	if m.AuthenticationParameterRAND != nil {
		n += m.AuthenticationParameterRAND.MarshalLenWithoutTag()
	}
	if m.AuthenticationParameterAUTN != nil {
		n += m.AuthenticationParameterAUTN.MarshalLenWithoutTag()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the AUTHENTICATION REQUEST message into b. The nil receiver returns
// a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too short. Header, packed
// half-octets, positional IEs and extras follow TS 24.301 clause 8.2.7 table order.
func (m *AuthenticationRequest) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: AuthenticationRequest is nil", ie.ErrNilIE)
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

	// Packed half-octets (NAS key set ID high nibble, spare low nibble).
	if m.NASKeySetIdentifier != nil && m.SpareHalfOctet != nil {
		hb, err := ie.MarshalHalfOctets(m.NASKeySetIdentifier, m.SpareHalfOctet)
		if err != nil {
			return fmt.Errorf("half-octets: %w", err)
		}
		if len(b) < n+len(hb) {
			return io.ErrShortBuffer
		}
		copy(b[n:n+len(hb)], hb)
		n += len(hb)
	} else {
		if m.NASKeySetIdentifier != nil {
			if err := m.NASKeySetIdentifier.MarshalToWithoutTag(b[n:]); err != nil {
				return fmt.Errorf("NASKeySetIdentifier: %w", err)
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

	// Positional mandatory IEs in spec-table order.
	if m.AuthenticationParameterRAND != nil {
		n0 := n
		if err := m.AuthenticationParameterRAND.MarshalToWithoutTag(b[n:]); err != nil {
			return fmt.Errorf("AuthenticationParameterRAND: %w", err)
		}
		n += m.AuthenticationParameterRAND.MarshalLenWithoutTag()
		_ = n0
	}

	if m.AuthenticationParameterAUTN != nil {
		n0 := n
		if err := m.AuthenticationParameterAUTN.MarshalToWithoutTag(b[n:]); err != nil {
			return fmt.Errorf("AuthenticationParameterAUTN: %w", err)
		}
		n += m.AuthenticationParameterAUTN.MarshalLenWithoutTag()
		_ = n0
	}

	// Extra tagged IEs.
	for _, i := range m.IEs {
		w, err := writeIE(b, n, i)
		if err != nil {
			return fmt.Errorf("extra IE: %w", err)
		}
		n += w
	}

	return nil
}

// UnmarshalBinary decodes an AUTHENTICATION REQUEST message from b. The receiver is left
// unchanged on error. Header, half-octets via ParseHalfOctets, positional IEs via
// ie.ParseWithoutTag with prototypes. Per TS 24.301 Table 8.2.7.1 the body has four
// positional mandatory IEs (two half-octets, RAND, AUTN). No tagged optional part.
func (m *AuthenticationRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: AuthenticationRequest is nil", ie.ErrNilIE)
	}

	var tmp AuthenticationRequest

	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	pos := 0

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	// Positional mandatory half-octets (NAS key set ID high nibble, spare low nibble).
	lProto := ie.NewPositionalNASKeySetIdentifier(0)
	hProto := ie.NewPositionalSpareHalfOctet(0)
	nkidIE, spareIE, err := ie.ParseHalfOctets(pay[pos:], lProto, hProto)
	if err != nil {
		return fmt.Errorf("half-octets: %w", err)
	}
	tmp.NASKeySetIdentifier = nkidIE
	tmp.SpareHalfOctet = spareIE

	pos++ // consumed 1 byte for packed half-octets

	// Authentication parameter RAND (V, 16 bytes).
	tmp.AuthenticationParameterRAND, err = ie.ParseWithoutTag(
		pay[pos:], ie.NewPositionalAuthenticationParameterRAND(make([]byte, 16)),
	)
	if err != nil {
		return fmt.Errorf("AuthenticationParameterRAND: %w", err)
	}
	pos += tmp.AuthenticationParameterRAND.MarshalLenWithoutTag()

	// Authentication parameter AUTN (LV, 16 value bytes).
	tmp.AuthenticationParameterAUTN, err = ie.ParseWithoutTag(
		pay[pos:], ie.NewPositionalAuthenticationParameterAUTN(make([]byte, 16)),
	)
	if err != nil {
		return fmt.Errorf("AuthenticationParameterAUTN: %w", err)
	}

	*m = tmp
	return nil
}

// MessageType returns the AUTHENTICATION REQUEST message type constant (0x52).
func (m *AuthenticationRequest) MessageType() uint8 {
	return EMMAuthenticationRequest
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *AuthenticationRequest) MessageTypeName() string {
	return "Authentication request"
}

// IsRequest returns false — AUTHENTICATION REQUEST is sent network to UE.
func (m *AuthenticationRequest) IsRequest() bool {
	return false
}

// AuthenticationRequest implements the Message interface.
var _ Message = (*AuthenticationRequest)(nil)

// Header returns the EMM security header embedded in AuthenticationRequest.
func (m *AuthenticationRequest) Header() *EmmHeader {
	return &m.EmmHeader
}
