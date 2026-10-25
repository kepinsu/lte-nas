package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// ActivateDedicatedEPSBearerContextRequest represents an ACTIVATE DEDICATED EPS
// BEARER CONTEXT REQUEST message as defined in TS 24.301 clause 8.3.3
// (Table 8.3.3.1, direction network to UE).
//
// The positional mandatory part encodes linked EPS bearer identity and spare
// half octet as a packed type-1 pair (one wire byte), followed by EPS quality
// of service (LV positional) and traffic flow template (LV-E positional). Tagged
// optional LLC SAPI (tag 0x32) is resolved via prototype in ie.ParseMultiIEs.
type ActivateDedicatedEPSBearerContextRequest struct {
	EsmHeader
	LinkedEPSBearerIdentity *ie.IE // positional type-1 pair, high nibble (9.9.3.5a)
	SpareHalfOctet          *ie.IE // positional type-1 pair, low nibble (9.9.2.9)
	EPSQoS                  *ie.IE // positional LV type-4 (9.9.4.3)
	TrafficFlowTemplate     *ie.IE // positional LV-E type-6 (9.9.4.16)
	LLCSAPI                 *ie.IE // tagged optional type-3 TV, ambiguous tag 0x32 (9.9.4.7)
	IEs                     []*ie.IE
}

// NewActivateDedicatedEPSBearerContextRequest creates an ACTIVATE DEDICATED EPS
// BEARER CONTEXT REQUEST message with the mandatory positional IEs from TS
// 24.301 clause 8.3.3, followed by optional tagged IEs classified by Kind.
// LLCSAPI (tag 0x32) is resolved first; duplicates and unrecognized IEs are
// appended to the IEs field.
func NewActivateDedicatedEPSBearerContextRequest(linkedBearerIdentity, spareHalfOctet,
	ePSQoS, trafficFlowTemplate *ie.IE, ies ...*ie.IE,
) *ActivateDedicatedEPSBearerContextRequest {
	if linkedBearerIdentity == nil || spareHalfOctet == nil ||
		ePSQoS == nil || trafficFlowTemplate == nil {
		panic("ActivateDedicatedEPSBearerContextRequest positional IEs must not be nil")
	}

	m := &ActivateDedicatedEPSBearerContextRequest{
		LinkedEPSBearerIdentity: linkedBearerIdentity,
		SpareHalfOctet:          spareHalfOctet,
		EPSQoS:                  ePSQoS,
		TrafficFlowTemplate:     trafficFlowTemplate,
	}

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

// MarshalLen returns the wire size of the ACTIVATE DEDICATED EPS BEARER CONTEXT
// REQUEST message, or 0 for a nil receiver. The result includes the 3-octet ESM
// header, packed half-octets (1 byte), positional IEs and all named tagged IEs
// in spec-table order.
func (m *ActivateDedicatedEPSBearerContextRequest) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EsmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	n += 1 // packed half-octets (1 byte)

	if m.EPSQoS != nil {
		n += m.EPSQoS.MarshalLenWithoutTag()
	}
	if m.TrafficFlowTemplate != nil {
		n += m.TrafficFlowTemplate.MarshalLenWithoutTag()
	}

	if m.LLCSAPI != nil {
		n += m.LLCSAPI.MarshalLen()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the ACTIVATE DEDICATED EPS BEARER CONTEXT REQUEST message
// into b. The nil receiver returns a wrapped ie.ErrNilIE; io.ErrShortBuffer
// without writing if b is too short. Header, packed half-octets, positional IEs
// and tagged optional IEs follow TS 24.301 clause 8.3.3 table order.
func (m *ActivateDedicatedEPSBearerContextRequest) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: ActivateDedicatedEPSBearerContextRequest is nil", ie.ErrNilIE)
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

	// Packed half-octets: spare half octet high nibble, linked EPS bearer identity low
	// nibble (wire format: [7:4]=spare, [3:0]=linked_bearer_identity).
	hb, err := ie.MarshalHalfOctets(m.LinkedEPSBearerIdentity, m.SpareHalfOctet)
	if err != nil {
		return fmt.Errorf("half-octets: %w", err)
	}
	if len(b[n:]) < len(hb) {
		return io.ErrShortBuffer
	}
	copy(b[n:], hb)
	n += len(hb)

	if m.EPSQoS == nil {
		return io.ErrUnexpectedEOF
	}
	if err := m.EPSQoS.MarshalToWithoutTag(b[n:]); err != nil {
		return fmt.Errorf("EPSQoS: %w", err)
	}
	n += m.EPSQoS.MarshalLenWithoutTag()

	if m.TrafficFlowTemplate == nil {
		return io.ErrUnexpectedEOF
	}
	if err := m.TrafficFlowTemplate.MarshalToWithoutTag(b[n:]); err != nil {
		return fmt.Errorf("TrafficFlowTemplate: %w", err)
	}
	n += m.TrafficFlowTemplate.MarshalLenWithoutTag()

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

// UnmarshalBinary decodes an ACTIVATE DEDICATED EPS BEARER CONTEXT REQUEST
// message from b. The receiver is left unchanged on error (uses local tmp
// struct). Header, packed half-octets via ParseHalfOctets, positional IEs via
// ie.ParseWithoutTag with prototypes, tagged optional LLC SAPI resolved via
// prototype in ie.ParseMultiIEs.
func (m *ActivateDedicatedEPSBearerContextRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: ActivateDedicatedEPSBearerContextRequest is nil", ie.ErrNilIE)
	}

	var tmp ActivateDedicatedEPSBearerContextRequest

	// Decode header (3 octets for plain ESM).
	pay, err := tmp.EsmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	pos := 0

	// Packed half-octets (linked bearer identity + spare bits in one wire byte).
	if len(pay[pos:]) < 1 {
		return io.ErrUnexpectedEOF
	}

	lProto := ie.NewPositionalLinkedEPSBearerIdentity(0)
	hProto := ie.NewPositionalSpareHalfOctet(0)
	tmp.LinkedEPSBearerIdentity, tmp.SpareHalfOctet, err = ie.ParseHalfOctets(pay[pos:], lProto, hProto)
	if err != nil {
		return fmt.Errorf("half-octets: %w", err)
	}
	pos++

	// Positional EPS QoS (LV, no tag).
	if len(pay[pos:]) < 2 {
		return io.ErrUnexpectedEOF
	}

	tmp.EPSQoS, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalEPSQualityOfService([]byte{0, 0}))
	if err != nil {
		return fmt.Errorf("EPSQoS: %w", err)
	}
	pos += tmp.EPSQoS.MarshalLenWithoutTag()

	// Positional Traffic Flow Template (LV-E, no tag).
	if len(pay[pos:]) < 2 {
		return io.ErrUnexpectedEOF
	}

	tmp.TrafficFlowTemplate, err = ie.ParseWithoutTag(pay[pos:], ie.NewTrafficFlowTemplate([]byte{0, 0}))
	if err != nil {
		return fmt.Errorf("TrafficFlowTemplate: %w", err)
	}
	pos += tmp.TrafficFlowTemplate.MarshalLenWithoutTag()

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

// MessageType returns the ACTIVATE DEDICATED EPS BEARER CONTEXT REQUEST message
// type constant (0xC5).
func (m *ActivateDedicatedEPSBearerContextRequest) MessageType() uint8 {
	return ESMActivateDedicatedEPSBearerContextRequest
}

// MessageTypeName returns the descriptive name from Table 9.8.2.
func (m *ActivateDedicatedEPSBearerContextRequest) MessageTypeName() string {
	return "Activate dedicated EPS bearer context request"
}

// IsRequest returns false — ACTIVATE DEDICATED EPS BEARER CONTEXT REQUEST is
// sent network to UE.
func (m *ActivateDedicatedEPSBearerContextRequest) IsRequest() bool {
	return false
}

// ActivateDedicatedEPSBearerContextRequest implements the Message interface.
var _ Message = (*ActivateDedicatedEPSBearerContextRequest)(nil)

// ParseActivateDedicatedEPSBearerContextRequest parses a byte slice as an ACTIVATE
// DEDICATED EPS BEARER CONTEXT REQUEST message. Convenience wrapper around
// UnmarshalBinary on a fresh receiver.
func ParseActivateDedicatedEPSBearerContextRequest(b []byte) (*ActivateDedicatedEPSBearerContextRequest, error) {
	m := &ActivateDedicatedEPSBearerContextRequest{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}
