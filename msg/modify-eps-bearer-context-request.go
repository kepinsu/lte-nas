package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// ModifyEPSBearerContextRequest represents a MODIFY EPS BEARER CONTEXT REQUEST
// message as defined in TS 24.301 clause 8.3.18 (Table 8.3.18.1, direction
// network to UE).
//
// The positional mandatory part encodes EPS quality of service (LV, no tag) and
// traffic flow template (LV-E, positional, where applicable). Tagged optional
// IEs: Extended quality of service (tag 0x5C) and other IEs from the spec table.
type ModifyEPSBearerContextRequest struct {
	EsmHeader
	EPSQoS              *ie.IE // positional LV type-4 (9.9.4.3)
	TrafficFlowTemplate *ie.IE // positional LV-E type-6 (9.9.4.16), where applicable
	ExtendedQoS         *ie.IE // tagged optional TV type-4, ambiguous tag 0x5C (9.9.4.30)
	IEs                 []*ie.IE
}

// NewModifyEPSBearerContextRequest creates a MODIFY EPS BEARER CONTEXT REQUEST
// message with the mandatory positional IEs from TS 24.301 clause 8.3.18,
// followed by optional tagged IEs classified by Kind. TrafficFlowTemplate is
// variadic (TFT where applicable per the spec). ExtendedQoS (Kind 0x5C) is
// resolved first; duplicates and unrecognized IEs go to the IEs field.
func NewModifyEPSBearerContextRequest(ePSQoS *ie.IE, ies ...*ie.IE) *ModifyEPSBearerContextRequest {
	if ePSQoS == nil {
		panic("ModifyEPSBearerContextRequest positional EPS QoS must not be nil")
	}

	m := &ModifyEPSBearerContextRequest{
		EPSQoS: ePSQoS,
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
		case ie.TrafficFlowTemplate:
			if m.TrafficFlowTemplate == nil {
				m.TrafficFlowTemplate = i
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

// MarshalLen returns the wire size of the MODIFY EPS BEARER CONTEXT REQUEST
// message, or 0 for a nil receiver. The result includes the 3-octet ESM header,
// positional EPS QoS (LV), optional TFT (LV-E), tagged IEs and extras.
func (m *ModifyEPSBearerContextRequest) MarshalLen() int {
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
	if m.TrafficFlowTemplate != nil {
		n += m.TrafficFlowTemplate.MarshalLenWithoutTag()
	}
	if m.ExtendedQoS != nil {
		n += m.ExtendedQoS.MarshalLen()
	}
	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the MODIFY EPS BEARER CONTEXT REQUEST message into b. The
// nil receiver returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing
// if b is too short. Header, positional IEs, named tagged IEs and extras follow
// TS 24.301 clause 8.3.18 table order.
func (m *ModifyEPSBearerContextRequest) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: ModifyEPSBearerContextRequest is nil", ie.ErrNilIE)
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

	if m.TrafficFlowTemplate != nil {
		if err := m.TrafficFlowTemplate.MarshalToWithoutTag(b[n:]); err != nil {
			return fmt.Errorf("TrafficFlowTemplate: %w", err)
		}
		n += m.TrafficFlowTemplate.MarshalLenWithoutTag()
	}

	// Named tagged IEs in spec-table order.
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

// UnmarshalBinary decodes a MODIFY EPS BEARER CONTEXT REQUEST message from b.
// The receiver is left unchanged on error (uses local tmp struct). Header via
// UnmarshalBinary, positional EPS QoS and optional TFT via ParseWithoutTag with
// prototypes, tagged IEs resolved via ie.ParseMultiIEs with 0x5C = extended QoS.
func (m *ModifyEPSBearerContextRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: ModifyEPSBearerContextRequest is nil", ie.ErrNilIE)
	}

	var tmp ModifyEPSBearerContextRequest

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

	// Positional optional TFT (LV-E, no tag, present where applicable).
	if len(pay[pos:]) > 0 {
		tmp.TrafficFlowTemplate, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalTrafficFlowTemplate([]byte{0, 0}))
		if err != nil {
			return fmt.Errorf("TrafficFlowTemplate: %w", err)
		}
		pos += tmp.TrafficFlowTemplate.MarshalLenWithoutTag()
	}

	// Tagged optional IEs — prototype resolves ambiguous tag 0x5C to extended QoS.
	if len(pay[pos:]) > 0 {
		tagged, err := ie.ParseMultiIEs(pay[pos:])
		if err != nil {
			return fmt.Errorf("optional: %w", err)
		}

		for _, i := range tagged {
			switch i.Kind {
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

// MessageType returns the MODIFY EPS BEARER CONTEXT REQUEST message type constant (0xC9).
func (m *ModifyEPSBearerContextRequest) MessageType() uint8 {
	return ESMModifyEPSBearerContextRequest
}

// MessageTypeName returns the descriptive name from Table 9.8.2.
func (m *ModifyEPSBearerContextRequest) MessageTypeName() string {
	return "Modify EPS bearer context request"
}

// IsRequest returns false — MODIFY EPS BEARER CONTEXT REQUEST is sent network to UE.
func (m *ModifyEPSBearerContextRequest) IsRequest() bool {
	return false
}

// ModifyEPSBearerContextRequest implements the Message interface.
var _ Message = (*ModifyEPSBearerContextRequest)(nil)

// ParseModifyEPSBearerContextRequest parses a byte slice as a MODIFY EPS BEARER
// CONTEXT REQUEST message. Convenience wrapper around UnmarshalBinary on a fresh
// receiver.
func ParseModifyEPSBearerContextRequest(b []byte) (*ModifyEPSBearerContextRequest, error) {
	m := &ModifyEPSBearerContextRequest{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}
