package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// BearerResourceAllocationReject represents a BEARER RESOURCE ALLOCATION REJECT
// message as defined in TS 24.301 clause 8.3.7 (Table 8.3.7.1, direction
// network to UE).
//
// Positional mandatory: ESM cause (type-1, 1 octet). Tagged optional IEs:
// LLC SAPI (tag 0x32) and extended QoS (tag 0x5C).
type BearerResourceAllocationReject struct {
	EsmHeader
	ESMCause    *ie.IE // positional type-1 (9.9.4.4)
	LLCSAPI     *ie.IE // tagged optional TV type-3, ambiguous tag 0x32 (9.9.4.7)
	ExtendedQoS *ie.IE // tagged optional TLV type-4, ambiguous tag 0x5C (9.9.4.30)
	IEs         []*ie.IE
}

// NewBearerResourceAllocationReject creates a BEARER RESOURCE ALLOCATION REJECT
// message with the mandatory positional ESM cause IE, followed by optional
// tagged IEs classified by Kind. ESMCause is the positional mandatory field;
// duplicates and unrecognized IEs are appended to the IEs field.
func NewBearerResourceAllocationReject(esmCause *ie.IE, ies ...*ie.IE) *BearerResourceAllocationReject {
	if esmCause == nil {
		panic("BearerResourceAllocationReject positional ESM cause must not be nil")
	}

	m := &BearerResourceAllocationReject{
		ESMCause: esmCause,
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		// Normalize collision tags (0x32, 0x58, 0x5b, 0x5c) to FormatFour TLV.
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
		case ie.ExtendedQualityOfService:
			if m.ExtendedQoS == nil {
				m.ExtendedQoS = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the BEARER RESOURCE ALLOCATION REJECT
// message, or 0 for a nil receiver. The result includes the 3-octet ESM header,
// positional ESM cause and all named tagged IEs.
func (m *BearerResourceAllocationReject) MarshalLen() int {
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
	if m.ExtendedQoS != nil {
		n += m.ExtendedQoS.MarshalLen()
	}
	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the BEARER RESOURCE ALLOCATION REJECT message into b.
// The nil receiver returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without
// writing if b is too short. Header, positional ESM cause, named tagged IEs
// and extras follow TS 24.301 clause 8.3.7 table order.
func (m *BearerResourceAllocationReject) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: BearerResourceAllocationReject is nil", ie.ErrNilIE)
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

	// Positional ESM cause.
	if m.ESMCause != nil {
		n0 := n
		if err := m.ESMCause.MarshalToWithoutTag(b[n:]); err != nil {
			return fmt.Errorf("ESMCause: %w", err)
		}
		n += m.ESMCause.MarshalLenWithoutTag()
		_ = n0 // unused; prevents unused-variable lint.
	}

	// Named tagged IEs in spec-table order.
	if m.LLCSAPI != nil {
		w, err := writeIE(b, n, m.LLCSAPI)
		if err != nil {
			return fmt.Errorf("LLCSAPI: %w", err)
		}
		n += w
	}
	if m.ExtendedQoS != nil {
		w, err := writeIE(b, n, m.ExtendedQoS)
		if err != nil {
			return fmt.Errorf("ExtendedQoS: %w", err)
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

// UnmarshalBinary decodes a BEARER RESOURCE ALLOCATION REJECT message from b.
// The receiver is left unchanged on error (uses local tmp struct). Header via
// UnmarshalBinary, positional ESM cause via ParseWithoutTag with prototype,
// tagged optional IEs resolved via ie.ParseMultiIEs.
func (m *BearerResourceAllocationReject) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: BearerResourceAllocationReject is nil", ie.ErrNilIE)
	}

	var tmp BearerResourceAllocationReject

	// Decode header (3 octets for plain ESM).
	pay, err := tmp.EsmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	pos := 0

	// Positional mandatory ESM cause (type-1, no tag).
	if len(pay[pos:]) == 0 {
		return io.ErrUnexpectedEOF
	}

	tmp.ESMCause, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalESMCause(ie.NetworkFailureSmCause))
	if err != nil {
		return fmt.Errorf("ESMCause: %w", err)
	}
	pos += tmp.ESMCause.MarshalLenWithoutTag()

	// Tagged optional IEs — prototypes resolve ambiguous tags 0x32 to
	// LLC SAPI and 0x5C to extended QoS.
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
			case ie.ExtendedQualityOfService:
				if tmp.ExtendedQoS == nil {
					tmp.ExtendedQoS = i
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

// MessageType returns the BEARER RESOURCE ALLOCATION REJECT message type
// constant (0xD5).
func (m *BearerResourceAllocationReject) MessageType() uint8 {
	return ESMBearerResourceAllocationReject
}

// MessageTypeName returns the descriptive name from Table 9.8.2.
func (m *BearerResourceAllocationReject) MessageTypeName() string {
	return "Bearer resource allocation reject"
}

// IsRequest returns false — BEARER RESOURCE ALLOCATION REJECT is sent network
// to UE.
func (m *BearerResourceAllocationReject) IsRequest() bool {
	return false
}

// BearerResourceAllocationReject implements the Message interface.
var _ Message = (*BearerResourceAllocationReject)(nil)

// ParseBearerResourceAllocationReject parses a byte slice as a BEARER RESOURCE
// ALLOCATION REJECT message. Convenience wrapper around UnmarshalBinary on a
// fresh receiver.
func ParseBearerResourceAllocationReject(b []byte) (*BearerResourceAllocationReject, error) {
	m := &BearerResourceAllocationReject{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}
