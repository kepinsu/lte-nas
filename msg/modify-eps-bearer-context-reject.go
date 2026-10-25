package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// ModifyEPSBearerContextReject represents a MODIFY EPS BEARER CONTEXT REJECT
// message as defined in TS 24.301 clause 8.3.17 (Table 8.3.17.1, direction UE
// to network).
//
// The positional mandatory part encodes ESM cause (V, 1 octet). Tagged optional
// IEs: Extended quality of service (tag 0x5C).
type ModifyEPSBearerContextReject struct {
	EsmHeader
	ESMCause    *ie.IE // positional V type-3 (9.9.4.4), one octet
	ExtendedQoS *ie.IE // tagged optional TV type-4, ambiguous tag 0x5C (9.9.4.30)
	IEs         []*ie.IE
}

// NewModifyEPSBearerContextReject creates a MODIFY EPS BEARER CONTEXT REJECT
// message with the mandatory positional ESM cause from TS 24.301 clause 8.3.17,
// followed by optional tagged IEs classified by Kind. ExtendedQoS (Kind 0x5C)
// is resolved first; duplicates and unrecognized IEs go to the IEs field.
func NewModifyEPSBearerContextReject(esmCause *ie.IE, ies ...*ie.IE) *ModifyEPSBearerContextReject {
	if esmCause == nil {
		panic("ModifyEPSBearerContextReject positional ESM cause must not be nil")
	}

	m := &ModifyEPSBearerContextReject{
		ESMCause: esmCause,
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
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the MODIFY EPS BEARER CONTEXT REJECT message,
// or 0 for a nil receiver. The result includes the 3-octet ESM header, positional
// ESM cause and all named tagged IEs.
func (m *ModifyEPSBearerContextReject) MarshalLen() int {
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
	if m.ExtendedQoS != nil {
		n += m.ExtendedQoS.MarshalLen()
	}
	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the MODIFY EPS BEARER CONTEXT REJECT message into b. The nil
// receiver returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b
// is too short. Header, positional ESM cause, named tagged IEs and extras follow
// TS 24.301 clause 8.3.17 table order.
func (m *ModifyEPSBearerContextReject) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: ModifyEPSBearerContextReject is nil", ie.ErrNilIE)
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

// UnmarshalBinary decodes a MODIFY EPS BEARER CONTEXT REJECT message from b.
// The receiver is left unchanged on error (uses local tmp struct). Header via
// UnmarshalBinary, positional ESM cause via ParseWithoutTag with prototypes,
// tagged IEs resolved via ie.ParseMultiIEs (0x5C = extended QoS).
func (m *ModifyEPSBearerContextReject) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: ModifyEPSBearerContextReject is nil", ie.ErrNilIE)
	}

	var tmp ModifyEPSBearerContextReject

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

// MessageType returns the MODIFY EPS BEARER CONTEXT REJECT message type constant (0xCB).
func (m *ModifyEPSBearerContextReject) MessageType() uint8 {
	return ESMModifyEPSBearerContextReject
}

// MessageTypeName returns the descriptive name from Table 9.8.2.
func (m *ModifyEPSBearerContextReject) MessageTypeName() string {
	return "Modify EPS bearer context reject"
}

// IsRequest returns false — MODIFY EPS BEARER CONTEXT REJECT is sent network to UE.
func (m *ModifyEPSBearerContextReject) IsRequest() bool {
	return false
}

// ModifyEPSBearerContextReject implements the Message interface.
var _ Message = (*ModifyEPSBearerContextReject)(nil)

// ParseModifyEPSBearerContextReject parses a byte slice as a MODIFY EPS BEARER
// CONTEXT REJECT message. Convenience wrapper around UnmarshalBinary on a fresh
// receiver.
func ParseModifyEPSBearerContextReject(b []byte) (*ModifyEPSBearerContextReject, error) {
	m := &ModifyEPSBearerContextReject{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}
