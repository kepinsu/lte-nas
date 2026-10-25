package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// ModifyEPSBearerContextAccept represents a MODIFY EPS BEARER CONTEXT ACCEPT
// message as defined in TS 24.301 clause 8.3.16 (Table 8.3.16.1, direction UE
// to network).
//
// The positional mandatory part encodes EPS quality of service (LV, no tag). Tagged
// optional IEs: Extended quality of service (tag 0x5C), Protocol configuration options
// (type-4 TLV, tag 0x7a).
type ModifyEPSBearerContextAccept struct {
	EsmHeader
	EPSQoS             *ie.IE // positional LV type-4 (9.9.4.3)
	ExtendedQoS        *ie.IE // tagged optional TV type-4, ambiguous tag 0x5C (9.9.4.30)
	ProtocolConfigOpts *ie.IE // tagged optional TLV type-4, tag 0x7a (9.9.4.11)
	IEs                []*ie.IE
}

// NewModifyEPSBearerContextAccept creates a MODIFY EPS BEARER CONTEXT ACCEPT
// message with the mandatory positional IEs from TS 24.301 clause 8.3.16,
// followed by optional tagged IEs classified by Kind. ExtendedQoS (Kind 0x5C)
// and ProtocolConfigOpts (Kind 0x7a) are resolved first; duplicates and
// unrecognized IEs go to the IEs field.
func NewModifyEPSBearerContextAccept(ePSQoS *ie.IE, ies ...*ie.IE) *ModifyEPSBearerContextAccept {
	if ePSQoS == nil {
		panic("ModifyEPSBearerContextAccept positional EPS QoS must not be nil")
	}

	m := &ModifyEPSBearerContextAccept{
		EPSQoS: ePSQoS,
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
		case ie.ExtendedQualityOfService:
			if m.ExtendedQoS == nil {
				m.ExtendedQoS = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.ProtocolConfigurationOptions:
			if m.ProtocolConfigOpts == nil {
				m.ProtocolConfigOpts = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the MODIFY EPS BEARER CONTEXT ACCEPT message,
// or 0 for a nil receiver. The result includes the 3-octet ESM header, positional
// EPS QoS and all named tagged IEs in spec-table order.
func (m *ModifyEPSBearerContextAccept) MarshalLen() int {
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
	if m.ExtendedQoS != nil {
		n += m.ExtendedQoS.MarshalLen()
	}
	if m.ProtocolConfigOpts != nil {
		n += m.ProtocolConfigOpts.MarshalLen()
	}
	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the MODIFY EPS BEARER CONTEXT ACCEPT message into b. The nil
// receiver returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b
// is too short. Header, positional EPS QoS, named tagged IEs and extras follow
// TS 24.301 clause 8.3.16 table order.
func (m *ModifyEPSBearerContextAccept) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: ModifyEPSBearerContextAccept is nil", ie.ErrNilIE)
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
	if m.ExtendedQoS != nil {
		w, err := writeIE(b, n, m.ExtendedQoS)
		if err != nil {
			return fmt.Errorf("ExtendedQoS: %w", err)
		}
		n += w
	}
	if m.ProtocolConfigOpts != nil {
		w, err := writeIE(b, n, m.ProtocolConfigOpts)
		if err != nil {
			return fmt.Errorf("ProtocolConfigOpts: %w", err)
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

// UnmarshalBinary decodes a MODIFY EPS BEARER CONTEXT ACCEPT message from b.
// The receiver is left unchanged on error (uses local tmp struct). Header via
// UnmarshalBinary, positional EPS QoS via ParseWithoutTag with prototypes, tagged
// IEs resolved via ie.ParseMultiIEs (0x5C = extended QoS, 0x7a = PCO).
func (m *ModifyEPSBearerContextAccept) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: ModifyEPSBearerContextAccept is nil", ie.ErrNilIE)
	}

	var tmp ModifyEPSBearerContextAccept

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

	// Tagged optional IEs — prototypes resolve ambiguous tag 0x5C to extended QoS,
	// and tag 0x7a to protocol configuration options.
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
			case ie.ProtocolConfigurationOptions:
				if tmp.ProtocolConfigOpts == nil {
					tmp.ProtocolConfigOpts = i
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

// MessageType returns the MODIFY EPS BEARER CONTEXT ACCEPT message type constant (0xCA).
func (m *ModifyEPSBearerContextAccept) MessageType() uint8 {
	return ESMModifyEPSBearerContextAccept
}

// MessageTypeName returns the descriptive name from Table 9.8.2.
func (m *ModifyEPSBearerContextAccept) MessageTypeName() string {
	return "Modify EPS bearer context accept"
}

// IsRequest returns false — MODIFY EPS BEARER CONTEXT ACCEPT is sent network to UE.
func (m *ModifyEPSBearerContextAccept) IsRequest() bool {
	return false
}

// ModifyEPSBearerContextAccept implements the Message interface.
var _ Message = (*ModifyEPSBearerContextAccept)(nil)

// ParseModifyEPSBearerContextAccept parses a byte slice as a MODIFY EPS BEARER
// CONTEXT ACCEPT message. Convenience wrapper around UnmarshalBinary on a fresh
// receiver.
func ParseModifyEPSBearerContextAccept(b []byte) (*ModifyEPSBearerContextAccept, error) {
	m := &ModifyEPSBearerContextAccept{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}
