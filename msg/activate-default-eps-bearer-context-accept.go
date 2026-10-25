package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// ActivateDefaultEPSBearerContextAccept represents an ACTIVATE DEFAULT EPS
// BEARER CONTEXT ACCEPT message as defined in TS 24.301 clause 8.3.4
// (Table 8.3.4.1, direction network to UE).
//
// The positional mandatory part encodes EPS quality of service and protocol
// configuration options without IEI tags. Tagged optional IEs (Extended PCO
// tag 0x7b, Traffic Flow Template tag 0x36/0x85 depending on format) are
// resolved via prototypes in ie.ParseMultiIEs.
type ActivateDefaultEPSBearerContextAccept struct {
	EsmHeader
	EPSQoS                    *ie.IE // positional LV (9.9.4.3)
	ProtocolConfigurationOpts *ie.IE // positional LV (9.9.4.11)
	ExtendedPCO               *ie.IE // tagged optional type-6 TLV-E (tag 0x7b, 9.9.4.26)
	TrafficFlowTemplate       *ie.IE // tagged optional type-4 TLV (tag 0x36, 9.9.4.16)
	OperatorDeterminedBearer  *ie.IE // tagged optional type-4 TLV (tag 0x7e)
	IEs                       []*ie.IE
}

// NewActivateDefaultEPSBearerContextAccept creates an ACTIVATE DEFAULT EPS
// BEARER CONTEXT ACCEPT message with the mandatory positional IEs from TS
// 24.301 clause 8.3.4, followed by optional tagged IEs classified by Kind.
// ExtendedPCO (tag 0x7b), TrafficFlowTemplate (tag 0x36) and
// OperatorDeterminedBearer (tag 0x7e) are resolved first; duplicates and
// unrecognized IEs are appended to the IEs field.
func NewActivateDefaultEPSBearerContextAccept(ePSQoS, protocolConfigurationOpts *ie.IE,
	ies ...*ie.IE,
) *ActivateDefaultEPSBearerContextAccept {
	if ePSQoS == nil || protocolConfigurationOpts == nil {
		panic("ActivateDefaultEPSBearerContextAccept positional IEs must not be nil")
	}

	m := &ActivateDefaultEPSBearerContextAccept{
		EPSQoS:                    ePSQoS,
		ProtocolConfigurationOpts: protocolConfigurationOpts,
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
		case ie.ExtendedProtocolConfigurationOptions:
			if m.ExtendedPCO == nil {
				m.ExtendedPCO = i
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

// MarshalLen returns the wire size of the ACTIVATE DEFAULT EPS BEARER CONTEXT
// ACCEPT message, or 0 for a nil receiver. The result includes the 3-octet
// ESM header plus positional IEs and all named tagged IEs in spec-table order.
func (m *ActivateDefaultEPSBearerContextAccept) MarshalLen() int {
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
	if m.ProtocolConfigurationOpts != nil {
		n += m.ProtocolConfigurationOpts.MarshalLenWithoutTag()
	}

	if m.ExtendedPCO != nil {
		n += m.ExtendedPCO.MarshalLen()
	}
	if m.TrafficFlowTemplate != nil {
		n += m.TrafficFlowTemplate.MarshalLen()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the ACTIVATE DEFAULT EPS BEARER CONTEXT ACCEPT message
// into b. The nil receiver returns a wrapped ie.ErrNilIE; io.ErrShortBuffer
// without writing if b is too short. Header, positional IEs in spec-table
// order, then named tagged IEs follow.
func (m *ActivateDefaultEPSBearerContextAccept) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: ActivateDefaultEPSBearerContextAccept is nil", ie.ErrNilIE)
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

	if m.EPSQoS == nil {
		return io.ErrUnexpectedEOF
	}
	n := 0
	if err := m.EPSQoS.MarshalToWithoutTag(b[n:]); err != nil {
		return fmt.Errorf("EPSQoS: %w", err)
	}
	n += m.EPSQoS.MarshalLenWithoutTag()

	if m.ProtocolConfigurationOpts == nil {
		return io.ErrUnexpectedEOF
	}
	if err := m.ProtocolConfigurationOpts.MarshalToWithoutTag(b[n:]); err != nil {
		return fmt.Errorf("ProtocolConfigurationOpts: %w", err)
	}
	n += m.ProtocolConfigurationOpts.MarshalLenWithoutTag()

	// Named tagged IEs in spec-table order.
	if m.ExtendedPCO != nil {
		w, err := writeIE(b, n, m.ExtendedPCO)
		if err != nil {
			return fmt.Errorf("ExtendedPCO: %w", err)
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

// UnmarshalBinary decodes an ACTIVATE DEFAULT EPS BEARER CONTEXT ACCEPT
// message from b. The receiver is left unchanged on error (uses local tmp
// struct). Header, positional IEs via ie.ParseWithoutTag with prototypes,
// tagged optional IEs resolved via ie.ParseMultiIEs.
func (m *ActivateDefaultEPSBearerContextAccept) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: ActivateDefaultEPSBearerContextAccept is nil", ie.ErrNilIE)
	}

	var tmp ActivateDefaultEPSBearerContextAccept

	// Decode header (3 octets for plain ESM).
	pay, err := tmp.EsmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	pos := 0

	tmp.EPSQoS, err = ie.ParseWithoutTag(pay[pos:],
		ie.NewPositionalEPSQualityOfService([]byte{0, 0}))
	if err != nil {
		return fmt.Errorf("EPSQoS: %w", err)
	}
	pos += tmp.EPSQoS.MarshalLenWithoutTag()

	tmp.ProtocolConfigurationOpts, err = ie.ParseWithoutTag(pay[pos:],
		ie.NewPositionalProtocolConfigurationOptions([]byte{0, 0}))
	if err != nil {
		return fmt.Errorf("ProtocolConfigurationOpts: %w", err)
	}
	pos += tmp.ProtocolConfigurationOpts.MarshalLenWithoutTag()

	// Tagged optional IEs — prototypes resolve ambiguous tags.
	tagged, err := ie.ParseMultiIEs(pay[pos:])
	if err != nil {
		return fmt.Errorf("optional: %w", err)
	}

	for _, i := range tagged {
		switch i.Kind {
		case ie.ExtendedProtocolConfigurationOptions:
			if tmp.ExtendedPCO == nil {
				tmp.ExtendedPCO = i
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

	*m = tmp
	return nil
}

// MessageType returns the ACTIVATE DEFAULT EPS BEARER CONTEXT ACCEPT message
// type constant (0xC2).
func (m *ActivateDefaultEPSBearerContextAccept) MessageType() uint8 {
	return ESMActivateDefaultEPSBearerContextAccept
}

// MessageTypeName returns the descriptive name from Table 9.8.2.
func (m *ActivateDefaultEPSBearerContextAccept) MessageTypeName() string {
	return "Activate default EPS bearer context accept"
}

// IsRequest returns false — ACTIVATE DEFAULT EPS BEARER CONTEXT ACCEPT is
// sent network to UE.
func (m *ActivateDefaultEPSBearerContextAccept) IsRequest() bool {
	return false
}

// ActivateDefaultEPSBearerContextAccept implements the Message interface.
var _ Message = (*ActivateDefaultEPSBearerContextAccept)(nil)

// ParseActivateDefaultEPSBearerContextAccept parses a byte slice as an ACTIVATE
// DEFAULT EPS BEARER CONTEXT ACCEPT message. Convenience wrapper around
// UnmarshalBinary on a fresh receiver.
func ParseActivateDefaultEPSBearerContextAccept(b []byte) (*ActivateDefaultEPSBearerContextAccept, error) {
	m := &ActivateDefaultEPSBearerContextAccept{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}
