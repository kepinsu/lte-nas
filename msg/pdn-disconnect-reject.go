package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// PDNDisconnectReject represents a PDN DISCONNECT REJECT message as defined in TS
// 24.301 clause 8.3.21 (Table 8.3.21.1, direction network to UE).
//
// The positional mandatory part encodes ESM cause (V, 1 octet) without IEI tag. Tagged
// optional IEs: LLC SAPI (tag 0x32) is resolved via prototype in ie.ParseMultiIEs.
type PDNDisconnectReject struct {
	EsmHeader
	ESMCause *ie.IE // positional V type-3 (9.9.4.4), one octet
	LLCSAPI  *ie.IE // tagged optional type-3 TV, ambiguous tag 0x32 (9.9.4.7)
	IEs      []*ie.IE
}

// NewPDNDisconnectReject creates a PDN DISCONNECT REJECT message with the mandatory
// positional ESM cause from TS 24.301 clause 8.3.21, followed by optional tagged IEs
// classified by Kind. LLCSAPI (Kind 0x32) is resolved first; duplicates and unrecognized
// IEs go to the IEs field.
func NewPDNDisconnectReject(esmCause *ie.IE, ies ...*ie.IE) *PDNDisconnectReject {
	if esmCause == nil {
		panic("PDNDisconnectReject positional ESM cause must not be nil")
	}

	m := &PDNDisconnectReject{ESMCause: esmCause}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
		case ie.LLCServiceAccessPointIdentifier:
			if m.LLCSAPI == nil {
				m.LLCSAPI = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the PDN DISCONNECT REJECT message, or 0 for
// a nil receiver. The result includes the 3-octet ESM header, positional ESM cause
// and all named tagged IEs.
func (m *PDNDisconnectReject) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EsmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	if m.ESMCause != nil {
		n += m.ESMCause.MarshalLenWithoutTag()
	}
	if m.LLCSAPI != nil {
		n += m.LLCSAPI.MarshalLen()
	}
	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the PDN DISCONNECT REJECT message into b. The nil receiver
// returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too short.
// Header, positional ESM cause, named tagged IEs and extras follow TS 24.301 clause
// 8.3.21 table order.
func (m *PDNDisconnectReject) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: PDNDisconnectReject is nil", ie.ErrNilIE)
	}

	hlen := m.EsmHeader.MarshalLen()
	if hlen == 0 {
		return fmt.Errorf("%w: EsmHeader MarshalLen returned 0", ie.ErrNilIE)
	}

	if len(b) < hlen {
		return io.ErrShortBuffer
	}

	if err := m.EsmHeader.MarshalTo(b); err != nil {
		return fmt.Errorf("header: %w", err)
	}

	b = b[hlen:]
	n := 0

	if m.ESMCause == nil {
		return io.ErrUnexpectedEOF
	}
	if err := m.ESMCause.MarshalToWithoutTag(b[n:]); err != nil {
		return fmt.Errorf("ESMCause: %w", err)
	}
	n += m.ESMCause.MarshalLenWithoutTag()

	// Named tagged IEs in spec-table order.
	if m.LLCSAPI != nil {
		w, err := writeIE(b, n, m.LLCSAPI)
		if err != nil {
			return fmt.Errorf("LLCSAPI: %w", err)
		}
		n += w
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

// UnmarshalBinary decodes a PDN DISCONNECT REJECT message from b. The receiver is left
// unchanged on error (uses local tmp struct). Header, positional ESM cause via ie.ParseWithoutTag
// with prototype, tagged optional IEs resolved via ie.ParseMultiIEs.
func (m *PDNDisconnectReject) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: PDNDisconnectReject is nil", ie.ErrNilIE)
	}

	var tmp PDNDisconnectReject

	// Decode header (3 octets for plain ESM).
	pay, err := tmp.EsmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	pos := 0

	// Positional ESM cause (V, 1 octet).
	if len(pay[pos:]) < 1 {
		return io.ErrUnexpectedEOF
	}

	tmp.ESMCause, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalESMCause(0))
	if err != nil {
		return fmt.Errorf("ESMCause: %w", err)
	}
	pos += tmp.ESMCause.MarshalLenWithoutTag()

	// Tagged optional IEs — prototype resolves ambiguous tag 0x32 to LLC SAPI.
	if len(pay[pos:]) > 0 {
		tagged, err := ie.ParseMultiIEs(pay[pos:])
		if err != nil {
			return fmt.Errorf("optional: %w", err)
		}

		for _, i := range tagged {
			switch i.Kind {
			case ie.LLCServiceAccessPointIdentifier:
				if tmp.LLCSAPI == nil {
					tmp.LLCSAPI = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			default:
				tmp.IEs = append(tmp.IEs, i)
			}
		}
	}

	*m = tmp
	return nil
}

// MessageType returns the PDN DISCONNECT REJECT message type constant (0xD3).
func (m *PDNDisconnectReject) MessageType() uint8 {
	return ESMPDNDisconnectReject
}

// MessageTypeName returns the descriptive name from Table 9.8.2.
func (m *PDNDisconnectReject) MessageTypeName() string {
	return "PDN disconnect reject"
}

// IsRequest returns false — PDN DISCONNECT REJECT is sent network to UE.
func (m *PDNDisconnectReject) IsRequest() bool {
	return false
}

// PDNDisconnectReject implements the Message interface.
var _ Message = (*PDNDisconnectReject)(nil)

// ParsePDNDisconnectReject parses a byte slice as a PDN DISCONNECT REJECT message.
// Convenience wrapper around UnmarshalBinary on a fresh receiver.
func ParsePDNDisconnectReject(b []byte) (*PDNDisconnectReject, error) {
	m := &PDNDisconnectReject{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}
