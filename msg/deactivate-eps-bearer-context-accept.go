package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// DeactivateEPSBearerContextAccept represents a DEACTIVATE EPS BEARER CONTEXT
// ACCEPT message as defined in TS 24.301 clause 8.3.11 (Table 8.3.11.1,
// direction UE to network).
//
// The positional mandatory part encodes EPS quality of service (LV, no tag). Tagged
// optional IEs: Protocol configuration options (type-4 TLV, tag 0x7a), Traffic flow
// template (type-4 TLV, tag 0x36).
type DeactivateEPSBearerContextAccept struct {
	EsmHeader
	EPSQoS              *ie.IE // positional LV type-4 (9.9.4.3)
	ProtocolConfigOpts  *ie.IE // tagged optional TLV type-4, tag 0x7a (9.9.4.11)
	TrafficFlowTemplate *ie.IE // tagged optional TLV type-4, tag 0x36 (9.9.4.16)
	IEs                 []*ie.IE
}

// NewDeactivateEPSBearerContextAccept creates a DEACTIVATE EPS BEARER CONTEXT
// ACCEPT message with the mandatory positional IEs from TS 24.301 clause 8.3.11,
// followed by optional tagged IEs classified by Kind. ProtocolConfigOpts (Kind 0x7a)
// and TrafficFlowTemplate (Kind 0x36) are resolved first; duplicates and
// unrecognized IEs go to the IEs field.
func NewDeactivateEPSBearerContextAccept(ePSQoS *ie.IE, ies ...*ie.IE) *DeactivateEPSBearerContextAccept {
	if ePSQoS == nil {
		panic("DeactivateEPSBearerContextAccept positional EPS QoS must not be nil")
	}

	m := &DeactivateEPSBearerContextAccept{
		EPSQoS: ePSQoS,
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
		case ie.ProtocolConfigurationOptions:
			if m.ProtocolConfigOpts == nil {
				m.ProtocolConfigOpts = i
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

// MarshalLen returns the wire size of the DEACTIVATE EPS BEARER CONTEXT ACCEPT
// message, or 0 for a nil receiver. The result includes the 3-octet ESM header,
// positional EPS QoS and all named tagged IEs.
func (m *DeactivateEPSBearerContextAccept) MarshalLen() int {
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
	if m.ProtocolConfigOpts != nil {
		n += m.ProtocolConfigOpts.MarshalLen()
	}
	if m.TrafficFlowTemplate != nil {
		n += m.TrafficFlowTemplate.MarshalLen()
	}
	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the DEACTIVATE EPS BEARER CONTEXT ACCEPT message into b. The
// nil receiver returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if
// b is too short. Header, positional EPS QoS, named tagged IEs and extras follow
// TS 24.301 clause 8.3.11 table order.
func (m *DeactivateEPSBearerContextAccept) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: DeactivateEPSBearerContextAccept is nil", ie.ErrNilIE)
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
	if m.ProtocolConfigOpts != nil {
		w, err := writeIE(b, n, m.ProtocolConfigOpts)
		if err != nil {
			return fmt.Errorf("ProtocolConfigOpts: %w", err)
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

// UnmarshalBinary decodes a DEACTIVATE EPS BEARER CONTEXT ACCEPT message from b.
// The receiver is left unchanged on error (uses local tmp struct). Header via
// UnmarshalBinary, positional EPS QoS via ParseWithoutTag with prototypes, tagged
// IEs resolved via ie.ParseMultiIEs (0x7a = PCO, 0x36 = TFT).
func (m *DeactivateEPSBearerContextAccept) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: DeactivateEPSBearerContextAccept is nil", ie.ErrNilIE)
	}

	var tmp DeactivateEPSBearerContextAccept

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

	// Tagged optional IEs — prototypes resolve tag 0x7a to protocol configuration
	// options and tag 0x36 to traffic flow template.
	if len(pay[pos:]) > 0 {
		tagged, err := ie.ParseMultiIEs(pay[pos:])
		if err != nil {
			return fmt.Errorf("optional: %w", err)
		}

		for _, i := range tagged {
			switch i.Kind {
			case ie.ProtocolConfigurationOptions:
				if tmp.ProtocolConfigOpts == nil {
					tmp.ProtocolConfigOpts = i
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

// MessageType returns the DEACTIVATE EPS BEARER CONTEXT ACCEPT message type constant (0xCE).
func (m *DeactivateEPSBearerContextAccept) MessageType() uint8 {
	return ESMDeactivateEPSBearerContextAccept
}

// MessageTypeName returns the descriptive name from Table 9.8.2.
func (m *DeactivateEPSBearerContextAccept) MessageTypeName() string {
	return "Deactivate EPS bearer context accept"
}

// IsRequest returns false — DEACTIVATE EPS BEARER CONTEXT ACCEPT is sent network to UE.
func (m *DeactivateEPSBearerContextAccept) IsRequest() bool {
	return false
}

// DeactivateEPSBearerContextAccept implements the Message interface.
var _ Message = (*DeactivateEPSBearerContextAccept)(nil)

// ParseDeactivateEPSBearerContextAccept parses a byte slice as a DEACTIVATE EPS
// BEARER CONTEXT ACCEPT message. Convenience wrapper around UnmarshalBinary on a
// fresh receiver.
func ParseDeactivateEPSBearerContextAccept(b []byte) (*DeactivateEPSBearerContextAccept, error) {
	m := &DeactivateEPSBearerContextAccept{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}
