package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// BearerResourceAllocationRequest represents a BEARER RESOURCE ALLOCATION
// REQUEST message as defined in TS 24.301 clause 8.3.8 (Table 8.3.8.1,
// direction UE to network).
//
// The positional mandatory part encodes the linked EPS bearer identity and
// spare half octet (packed type-1 pair, 1 wire byte), followed by EPS quality
// of service (LV positional) and traffic flow template (LV-E positional).
// Tagged optional IEs: LLC SAPI (tag 0x32) and extended QoS (tag 0x5C).
type BearerResourceAllocationRequest struct {
	EsmHeader
	BearerIdentity      *ie.IE // positional type-1 packed with SpareHalfOctet, Linked EPS bearer identity (9.9.4.6)
	SpareHalfOctet      *ie.IE // positional type-1 packed with BearerIdentity (9.9.2.9)
	EPSQoS              *ie.IE // positional LV type-4 (9.9.4.3)
	TrafficFlowTemplate *ie.IE // positional LV-E type-6 (9.9.4.16)
	LLCSAPI             *ie.IE // tagged optional TV type-3, ambiguous tag 0x32 (9.9.4.7)
	ExtendedQoS         *ie.IE // tagged optional TLV type-4, ambiguous tag 0x5C (9.9.4.30)
	IEs                 []*ie.IE
}

// NewBearerResourceAllocationRequest creates a BEARER RESOURCE ALLOCATION
// REQUEST message with the mandatory positional IEs from TS 24.301 clause
// 8.3.8, followed by optional tagged IEs classified by Kind.
func NewBearerResourceAllocationRequest(bearerIdentity, spareHalfOctet, ePSQoS,
	trafficFlowTemplate *ie.IE, ies ...*ie.IE,
) *BearerResourceAllocationRequest {
	if bearerIdentity == nil || spareHalfOctet == nil || ePSQoS == nil ||
		trafficFlowTemplate == nil {
		panic("BearerResourceAllocationRequest positional IEs must not be nil")
	}

	m := &BearerResourceAllocationRequest{
		BearerIdentity:      bearerIdentity,
		SpareHalfOctet:      spareHalfOctet,
		EPSQoS:              ePSQoS,
		TrafficFlowTemplate: trafficFlowTemplate,
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

// MarshalLen returns the wire size of the BEARER RESOURCE ALLOCATION REQUEST
// message, or 0 for a nil receiver. The result includes the 3-octet ESM header,
// packed half-octets (1 byte), positional IEs and all named tagged IEs.
func (m *BearerResourceAllocationRequest) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EsmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	// Packed half-octets: combined (1 byte) or individual.
	if m.BearerIdentity != nil && m.SpareHalfOctet != nil {
		n++ // packed half-octet
	} else {
		if m.BearerIdentity != nil {
			n++
		}
		if m.SpareHalfOctet != nil {
			n++
		}
	}

	if m.EPSQoS != nil {
		n += m.EPSQoS.MarshalLenWithoutTag()
	}
	if m.TrafficFlowTemplate != nil {
		n += m.TrafficFlowTemplate.MarshalLenWithoutTag()
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

// MarshalTo encodes the BEARER RESOURCE ALLOCATION REQUEST message into b.
// The nil receiver returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without
// writing if b is too short. Header, packed half-octets (or individual),
// positional IEs, named tagged IEs and extras follow TS 24.301 clause 8.3.8
// table order.
func (m *BearerResourceAllocationRequest) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: BearerResourceAllocationRequest is nil", ie.ErrNilIE)
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

	// Packed half-octets (linked EPS bearer identity low nibble, spare
	// half octet high nibble).
	if m.BearerIdentity != nil && m.SpareHalfOctet != nil {
		hb, err := ie.MarshalHalfOctets(m.BearerIdentity, m.SpareHalfOctet)
		if err != nil {
			return fmt.Errorf("half-octets: %w", err)
		}
		if len(b)-n < len(hb) {
			return io.ErrShortBuffer
		}
		copy(b[n:n+len(hb)], hb)
		n += len(hb)
	} else {
		if m.BearerIdentity != nil {
			if err := m.BearerIdentity.MarshalToWithoutTag(b[n:]); err != nil {
				return fmt.Errorf("BearerIdentity: %w", err)
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

	// Positional IEs in spec-table order.
	if m.EPSQoS != nil {
		n0 := n
		if err := m.EPSQoS.MarshalToWithoutTag(b[n:]); err != nil {
			return fmt.Errorf("EPSQoS: %w", err)
		}
		n += m.EPSQoS.MarshalLenWithoutTag()
		_ = n0 // unused; prevents unused-variable lint.
	}
	if m.TrafficFlowTemplate != nil {
		n0 := n
		if err := m.TrafficFlowTemplate.MarshalToWithoutTag(b[n:]); err != nil {
			return fmt.Errorf("TrafficFlowTemplate: %w", err)
		}
		n += m.TrafficFlowTemplate.MarshalLenWithoutTag()
		_ = n0
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

// UnmarshalBinary decodes a BEARER RESOURCE ALLOCATION REQUEST message from b.
// The receiver is left unchanged on error (uses local tmp struct). Header via
// UnmarshalBinary, packed half-octets via ParseHalfOctets with positional
// prototypes, remaining positional IEs via ie.ParseWithoutTag, tagged optional
// IEs resolved via ie.ParseMultiIEs with prototypes.
func (m *BearerResourceAllocationRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: BearerResourceAllocationRequest is nil", ie.ErrNilIE)
	}

	var tmp BearerResourceAllocationRequest

	// Decode header (3 octets for plain ESM).
	pay, err := tmp.EsmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	pos := 0

	// Packed half-octets (linked EPS bearer identity low nibble, spare
	// half octet high nibble).
	if len(pay[pos:]) == 0 {
		return io.ErrUnexpectedEOF
	}

	lProto := ie.NewPositionalLinkedEPSBearerIdentity(0)
	hProto := ie.NewPositionalSpareHalfOctet(0)
	lIE, hIE, err := ie.ParseHalfOctets(pay[pos:], lProto, hProto)
	if err != nil {
		return fmt.Errorf("half-octets: %w", err)
	}
	tmp.BearerIdentity = lIE
	tmp.SpareHalfOctet = hIE

	pos++ // consumed 1 byte for packed half-octets

	// Positional EPS QoS (LV, no tag).
	tmp.EPSQoS, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalEPSQualityOfService([]byte{0, 0}))
	if err != nil {
		return fmt.Errorf("EPSQoS: %w", err)
	}
	pos += tmp.EPSQoS.MarshalLenWithoutTag()

	// Positional optional TFT (LV-E, no tag).
	if len(pay[pos:]) > 0 {
		tmp.TrafficFlowTemplate, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalTrafficFlowTemplate([]byte{0, 0}))
		if err != nil {
			return fmt.Errorf("TrafficFlowTemplate: %w", err)
		}
		pos += tmp.TrafficFlowTemplate.MarshalLenWithoutTag()
	}

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

// MessageType returns the BEARER RESOURCE ALLOCATION REQUEST message type
// constant (0xD4).
func (m *BearerResourceAllocationRequest) MessageType() uint8 {
	return ESMBearerResourceAllocationRequest
}

// MessageTypeName returns the descriptive name from Table 9.8.2.
func (m *BearerResourceAllocationRequest) MessageTypeName() string {
	return "Bearer resource allocation request"
}

// IsRequest returns true — BEARER RESOURCE ALLOCATION REQUEST is sent UE to
// network.
func (m *BearerResourceAllocationRequest) IsRequest() bool {
	return true
}

// BearerResourceAllocationRequest implements the Message interface.
var _ Message = (*BearerResourceAllocationRequest)(nil)

// ParseBearerResourceAllocationRequest parses a byte slice as a BEARER RESOURCE
// ALLOCATION REQUEST message. Convenience wrapper around UnmarshalBinary on a
// fresh receiver.
func ParseBearerResourceAllocationRequest(b []byte) (*BearerResourceAllocationRequest, error) {
	m := &BearerResourceAllocationRequest{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}
