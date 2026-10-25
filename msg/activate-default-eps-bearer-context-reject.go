package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// ActivateDefaultEPSBearerContextReject represents an ACTIVATE DEFAULT EPS
// BEARER CONTEXT REJECT message as defined in TS 24.301 clause 8.3.5
// (Table 8.3.5.1, direction network to UE).
//
// The positional mandatory part encodes ESM cause (V, one octet) without IEI
// tag. PTI is carried in the ESM header (octet 2) and is not a separate field.
type ActivateDefaultEPSBearerContextReject struct {
	EsmHeader
	ESMCause *ie.IE // positional mandatory V (9.9.4.4)
	IEs      []*ie.IE
}

// NewActivateDefaultEPSBearerContextReject creates an ACTIVATE DEFAULT EPS
// BEARER CONTEXT REJECT message with the mandatory positional ESM cause IE.
func NewActivateDefaultEPSBearerContextReject(esmCause *ie.IE) *ActivateDefaultEPSBearerContextReject {
	if esmCause == nil {
		panic("ActivateDefaultEPSBearerContextReject positional IE must not be nil")
	}

	return &ActivateDefaultEPSBearerContextReject{ESMCause: esmCause}
}

// MarshalLen returns the wire size of the ACTIVATE DEFAULT EPS BEARER CONTEXT
// REJECT message, or 0 for a nil receiver. The result includes the 3-octet
// ESM header plus positional ESM cause.
func (m *ActivateDefaultEPSBearerContextReject) MarshalLen() int {
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

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the ACTIVATE DEFAULT EPS BEARER CONTEXT REJECT message
// into b. The nil receiver returns a wrapped ie.ErrNilIE; io.ErrShortBuffer
// without writing if b is too short. Header then positional ESM cause follow.
func (m *ActivateDefaultEPSBearerContextReject) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: ActivateDefaultEPSBearerContextReject is nil", ie.ErrNilIE)
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

	if m.ESMCause == nil {
		return io.ErrUnexpectedEOF
	}
	n := 0
	if err := m.ESMCause.MarshalToWithoutTag(b[n:]); err != nil {
		return fmt.Errorf("ESMCause: %w", err)
	}
	_ = m.ESMCause.MarshalLenWithoutTag()

	return nil
}

// UnmarshalBinary decodes an ACTIVATE DEFAULT EPS BEARER CONTEXT REJECT
// message from b. The receiver is left unchanged on error (uses local tmp
// struct). Header then positional ESM cause via ie.ParseWithoutTag.
func (m *ActivateDefaultEPSBearerContextReject) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: ActivateDefaultEPSBearerContextReject is nil", ie.ErrNilIE)
	}

	var tmp ActivateDefaultEPSBearerContextReject

	// Decode header (3 octets for plain ESM).
	pay, err := tmp.EsmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	pos := 0

	tmp.ESMCause, err = ie.ParseWithoutTag(pay[pos:],
		ie.NewPositionalESMCause(0))
	if err != nil {
		return fmt.Errorf("ESMCause: %w", err)
	}
	pos += tmp.ESMCause.MarshalLenWithoutTag()

	// Tagged optional IEs — prototype resolves ambiguous ESM cause tag.
	tagged, err := ie.ParseMultiIEs(pay[pos:])
	if err != nil {
		return fmt.Errorf("optional: %w", err)
	}

	for _, i := range tagged {
		switch i.Kind {
		case ie.ESMCause:
			if tmp.IEs == nil {
				tmp.IEs = append(tmp.IEs, i)
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

// MessageType returns the ACTIVATE DEFAULT EPS BEARER CONTEXT REJECT message
// type constant (0xC3).
func (m *ActivateDefaultEPSBearerContextReject) MessageType() uint8 {
	return ESMActivateDefaultEPSBearerContextReject
}

// MessageTypeName returns the descriptive name from Table 9.8.2.
func (m *ActivateDefaultEPSBearerContextReject) MessageTypeName() string {
	return "Activate default EPS bearer context reject"
}

// IsRequest returns false — ACTIVATE DEFAULT EPS BEARER CONTEXT REJECT is
// sent network to UE.
func (m *ActivateDefaultEPSBearerContextReject) IsRequest() bool {
	return false
}

// ActivateDefaultEPSBearerContextReject implements the Message interface.
var _ Message = (*ActivateDefaultEPSBearerContextReject)(nil)

// ParseActivateDefaultEPSBearerContextReject parses a byte slice as an ACTIVATE
// DEFAULT EPS BEARER CONTEXT REJECT message. Convenience wrapper around
// UnmarshalBinary on a fresh receiver.
func ParseActivateDefaultEPSBearerContextReject(b []byte) (*ActivateDefaultEPSBearerContextReject, error) {
	m := &ActivateDefaultEPSBearerContextReject{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}
