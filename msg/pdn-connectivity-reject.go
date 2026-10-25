package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// PDNConnectivityReject represents a PDN CONNECTIVITY REJECT message as defined in
// TS 24.301 clause 8.3.19 (Table 8.3.19.1, direction network to UE).
//
// There are no positional mandatory IEs beyond the ESM header. Tagged optional IEs:
// LLC SAPI (tag 0x32) and ESM cause (tag 0x58) are resolved via prototypes.
type PDNConnectivityReject struct {
	EsmHeader
	LLCSAPI  *ie.IE // tagged optional type-3 TV, ambiguous tag 0x32 (9.9.4.7)
	ESMCause *ie.IE // tagged optional type-3 TV, ambiguous tag 0x58 (9.9.4.4)
	IEs      []*ie.IE
}

// NewPDNConnectivityReject creates a PDN CONNECTIVITY REJECT message with optional
// tagged IEs classified by Kind. LLCSAPI (Kind 0x32) and ESMCause (Kind 0x58) are
// resolved first; duplicates and unrecognized IEs are appended to the IEs field.
func NewPDNConnectivityReject(ies ...*ie.IE) *PDNConnectivityReject {
	m := &PDNConnectivityReject{}

	for _, i := range ies {
		if i == nil {
			continue
		}
		// Normalize collision tags (0x32, 0x58) to FormatFour TLV.
		if ie.IsAmbiguousTagged(i.Kind) {
			i.Format = ie.IEFormatFour
			i.Length = uint16(len(i.Payload))
		}
		switch i.Kind {
		case ie.LLCServiceAccessPointIdentifier:
			if m.LLCSAPI == nil {
				m.LLCSAPI = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.ESMCause:
			if m.ESMCause == nil {
				m.ESMCause = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the PDN CONNECTIVITY REJECT message, or 0 for
// a nil receiver. The result includes the 3-octet ESM header and all named tagged IEs.
func (m *PDNConnectivityReject) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EsmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	if m.LLCSAPI != nil {
		n += m.LLCSAPI.MarshalLen()
	}
	if m.ESMCause != nil {
		n += m.ESMCause.MarshalLen()
	}
	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the PDN CONNECTIVITY REJECT message into b. The nil receiver
// returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too short.
// Header, named tagged IEs in spec-table order and extras follow TS 24.301 clause
// 8.3.19 table order.
func (m *PDNConnectivityReject) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: PDNConnectivityReject is nil", ie.ErrNilIE)
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

	// Named tagged IEs in spec-table order.
	if m.LLCSAPI != nil {
		w, err := writeIE(b, n, m.LLCSAPI)
		if err != nil {
			return fmt.Errorf("LLCSAPI: %w", err)
		}
		n += w
	}

	if m.ESMCause != nil {
		w, err := writeIE(b, n, m.ESMCause)
		if err != nil {
			return fmt.Errorf("ESMCause: %w", err)
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

// UnmarshalBinary decodes a PDN CONNECTIVITY REJECT message from b. The receiver is
// left unchanged on error (uses local tmp struct). Header decoded, tagged optional IEs
// resolved via ie.ParseMultiIEs. Ambiguous tags 0x32 (LLC SAPI) and 0x58 (ESM cause)
// are resolved per prototypes.
func (m *PDNConnectivityReject) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: PDNConnectivityReject is nil", ie.ErrNilIE)
	}

	var tmp PDNConnectivityReject

	// Decode header (3 octets for plain ESM).
	pay, err := tmp.EsmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	pos := 0

	// Tagged optional IEs — prototypes resolve ambiguous tags 0x32 to LLC SAPI and
	// 0x58 to ESM cause.
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
			case ie.ESMCause:
				if tmp.ESMCause == nil {
					tmp.ESMCause = i
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

// MessageType returns the PDN CONNECTIVITY REJECT message type constant (0xD1).
func (m *PDNConnectivityReject) MessageType() uint8 {
	return ESMPDNConnectivityReject
}

// MessageTypeName returns the descriptive name from Table 9.8.2.
func (m *PDNConnectivityReject) MessageTypeName() string {
	return "PDN connectivity reject"
}

// IsRequest returns false — PDN CONNECTIVITY REJECT is sent network to UE.
func (m *PDNConnectivityReject) IsRequest() bool {
	return false
}

// PDNConnectivityReject implements the Message interface.
var _ Message = (*PDNConnectivityReject)(nil)

// ParsePDNConnectivityReject parses a byte slice as a PDN CONNECTIVITY REJECT
// message. Convenience wrapper around UnmarshalBinary on a fresh receiver.
func ParsePDNConnectivityReject(b []byte) (*PDNConnectivityReject, error) {
	m := &PDNConnectivityReject{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}
