package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// ActivateDedicatedEPSBearerContextAccept represents an ACTIVATE DEDICATED EPS
// BEARER CONTEXT ACCEPT message as defined in TS 24.301 clause 8.3.1
// (Table 8.3.1.1, direction network to UE).
//
// The positional mandatory part encodes EPS quality of service (LV positional,
// no tag). Tagged optional IEs — LLC SAPI (tag 0x32), ESM cause (tag 0x58) and
// Traffic Flow Template (tag 0x36, type-4 TLV) — are resolved via prototypes in
// ie.ParseMultiIEs.
type ActivateDedicatedEPSBearerContextAccept struct {
	EsmHeader
	EPSQoS              *ie.IE // positional LV type-4 (9.9.4.3)
	LLCSAPI             *ie.IE // tagged optional type-3 TV, ambiguous tag 0x32 (9.9.4.7)
	ESMCause            *ie.IE // tagged optional type-3 TV, ambiguous tag 0x58 (9.9.4.4)
	TrafficFlowTemplate *ie.IE // tagged optional type-4 TLV, 0x36 (9.9.4.16)
	IEs                 []*ie.IE
}

// NewActivateDedicatedEPSBearerContextAccept creates an ACTIVATE DEDICATED EPS
// BEARER CONTEXT ACCEPT message with the mandatory positional IE from TS 24.301
// clause 8.3.1, followed by optional tagged IEs classified by Kind. LLCSAPI
// (tag 0x32), ESMCause (tag 0x58) and TrafficFlowTemplate (tag 0x36) are
// resolved first; duplicates and unrecognized IEs are appended to the IEs field.
func NewActivateDedicatedEPSBearerContextAccept(ePSQoS *ie.IE,
	ies ...*ie.IE,
) *ActivateDedicatedEPSBearerContextAccept {
	if ePSQoS == nil {
		panic("ActivateDedicatedEPSBearerContextAccept positional IE must not be nil")
	}

	m := &ActivateDedicatedEPSBearerContextAccept{
		EPSQoS: ePSQoS,
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
		case ie.ESMCause:
			if m.ESMCause == nil {
				m.ESMCause = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.TrafficFlowTemplate:
			if m.TrafficFlowTemplate == nil {
				m.TrafficFlowTemplate = i
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
// ACCEPT message, or 0 for a nil receiver. The result includes the 3-octet ESM
// header, positional IEs and all named tagged IEs in spec-table order.
func (m *ActivateDedicatedEPSBearerContextAccept) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EsmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	if m.EPSQoS != nil {
		n += m.EPSQoS.MarshalLenWithoutTag()
	}

	if m.LLCSAPI != nil {
		n += m.LLCSAPI.MarshalLen()
	}
	if m.ESMCause != nil {
		n += m.ESMCause.MarshalLen()
	}
	if m.TrafficFlowTemplate != nil {
		n += m.TrafficFlowTemplate.MarshalLen()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the ACTIVATE DEDICATED EPS BEARER CONTEXT ACCEPT message
// into b. The nil receiver returns a wrapped ie.ErrNilIE; io.ErrShortBuffer
// without writing if b is too short. Header, positional IEs and tagged optional
// IEs in spec-table order follow TS 24.301 clause 8.3.1 table order.
func (m *ActivateDedicatedEPSBearerContextAccept) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: ActivateDedicatedEPSBearerContextAccept is nil", ie.ErrNilIE)
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

	if m.EPSQoS == nil {
		return io.ErrUnexpectedEOF
	}
	if err := m.EPSQoS.MarshalToWithoutTag(b[n:]); err != nil {
		return fmt.Errorf("EPSQoS: %w", err)
	}
	n += m.EPSQoS.MarshalLenWithoutTag()

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

	if m.TrafficFlowTemplate != nil {
		w, err := writeIE(b, n, m.TrafficFlowTemplate)
		if err != nil {
			return fmt.Errorf("TrafficFlowTemplate: %w", err)
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

// UnmarshalBinary decodes an ACTIVATE DEDICATED EPS BEARER CONTEXT ACCEPT
// message from b. The receiver is left unchanged on error (uses local tmp
// struct). Header, positional EPS QoS via ie.ParseWithoutTag with prototype,
// tagged optional IEs resolved via prototypes in ie.ParseMultiIEs.
func (m *ActivateDedicatedEPSBearerContextAccept) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: ActivateDedicatedEPSBearerContextAccept is nil", ie.ErrNilIE)
	}

	var tmp ActivateDedicatedEPSBearerContextAccept

	// Decode header (3 octets for plain ESM).
	pay, err := tmp.EsmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	pos := 0

	// Positional EPS QoS (LV, no tag).
	if len(pay[pos:]) < 2 {
		return io.ErrUnexpectedEOF
	}

	tmp.EPSQoS, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalEPSQualityOfService([]byte{0, 0}))
	if err != nil {
		return fmt.Errorf("EPSQoS: %w", err)
	}
	pos += tmp.EPSQoS.MarshalLenWithoutTag()

	// Tagged optional IEs — prototypes resolve ambiguous tags 0x32 (LLC SAPI) and
	// 0x58 (ESM cause). Traffic Flow Template tag 0x36 is unambiguous.
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
			case ie.TrafficFlowTemplate:
				if tmp.TrafficFlowTemplate == nil {
					tmp.TrafficFlowTemplate = i
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

// MessageType returns the ACTIVATE DEDICATED EPS BEARER CONTEXT ACCEPT message
// type constant (0xC6).
func (m *ActivateDedicatedEPSBearerContextAccept) MessageType() uint8 {
	return ESMActivateDedicatedEPSBearerContextAccept
}

// MessageTypeName returns the descriptive name from Table 9.8.2.
func (m *ActivateDedicatedEPSBearerContextAccept) MessageTypeName() string {
	return "Activate dedicated EPS bearer context accept"
}

// IsRequest returns false — ACTIVATE DEDICATED EPS BEARER CONTEXT ACCEPT is
// sent network to UE.
func (m *ActivateDedicatedEPSBearerContextAccept) IsRequest() bool {
	return false
}

// ActivateDedicatedEPSBearerContextAccept implements the Message interface.
var _ Message = (*ActivateDedicatedEPSBearerContextAccept)(nil)

// ParseActivateDedicatedEPSBearerContextAccept parses a byte slice as an ACTIVATE
// DEDICATED EPS BEARER CONTEXT ACCEPT message. Convenience wrapper around
// UnmarshalBinary on a fresh receiver.
func ParseActivateDedicatedEPSBearerContextAccept(b []byte) (*ActivateDedicatedEPSBearerContextAccept, error) {
	m := &ActivateDedicatedEPSBearerContextAccept{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}
