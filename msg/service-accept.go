package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// ServiceAccept represents a SERVICE ACCEPT message as defined in TS 24.301
// clause 8.2.34 (Table 8.2.34.1, direction network to UE).
//
// This message has no mandatory content IEs — only the header is required. Tagged
// optional IEs include EPS bearer context status (0x57), T3448 value (0x6B),
// forbidden tracking area identity lists (0x1D / 0x1E) and EPS additional
// request result (0x37). The message is integrity-protected (SHT 0x01..0x05,
// 7-octet header).
type ServiceAccept struct {
	EmmHeader
	EPSBearerContextStatus *ie.IE // tagged optional TLV type-4 (9.9.2.1), tag 0x57
	T3448Value             *ie.IE // tagged optional TLV type-4 (9.9.3.16A), tag 0x6B
	IEs                    []*ie.IE
}

// NewServiceAccept creates a SERVICE ACCEPT message. There are no mandatory
// positional IEs per TS 24.301 Table 8.2.34.1; all content IEs are tagged
// optional and classified by their Kind. First occurrence wins; duplicates and
// unrecognized IEs are appended to the IEs field.
func NewServiceAccept(ies ...*ie.IE) *ServiceAccept {
	m := &ServiceAccept{}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
		case ie.EPSBearerContextStatus:
			if m.EPSBearerContextStatus == nil {
				m.EPSBearerContextStatus = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.T3448Value:
			if m.T3448Value == nil {
				m.T3448Value = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the SERVICE ACCEPT message, or 0 for a nil
// receiver. The result includes header (7 octets) and all named tagged IEs in
// spec-table order.
func (m *ServiceAccept) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	if m.EPSBearerContextStatus != nil {
		n += m.EPSBearerContextStatus.MarshalLen()
	}
	if m.T3448Value != nil {
		n += m.T3448Value.MarshalLen()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the SERVICE ACCEPT message into b. The nil receiver returns
// a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too short.
// Header and named tagged IEs follow TS 24.301 clause 8.2.34 table order.
func (m *ServiceAccept) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: ServiceAccept is nil", ie.ErrNilIE)
	}

	hlen := m.EmmHeader.MarshalLen()
	if hlen == 0 {
		return fmt.Errorf("%w: EmmHeader MarshalLen returned 0", ie.ErrNilIE)
	}

	if len(b) < hlen {
		return io.ErrShortBuffer
	}

	if err := m.EmmHeader.MarshalTo(b); err != nil {
		return fmt.Errorf("header: %w", err)
	}

	b = b[hlen:]

	if m.EPSBearerContextStatus != nil {
		w, err := writeIE(b, 0, m.EPSBearerContextStatus)
		if err != nil {
			return fmt.Errorf("EPSBearerContextStatus: %w", err)
		}
		b = b[w:]
	}
	if m.T3448Value != nil {
		w, err := writeIE(b, 0, m.T3448Value)
		if err != nil {
			return fmt.Errorf("T3448Value: %w", err)
		}
		b = b[w:]
	}

	for _, i := range m.IEs {
		w, err := writeIE(b, 0, i)
		if err != nil {
			return fmt.Errorf("extra IE: %w", err)
		}
		b = b[w:]
	}

	return nil
}

// UnmarshalBinary decodes a SERVICE ACCEPT message from b. The receiver is left
// unchanged on error. Header decoded, no positional mandatory part, tagged
// optional IEs resolved via expected prototypes.
func (m *ServiceAccept) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: ServiceAccept is nil", ie.ErrNilIE)
	}

	var tmp ServiceAccept

	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	if len(pay) == 0 {
		// No content IEs to decode — all fields remain nil, which is valid.
		*m = tmp
		return nil
	}

	// No positional mandatory IEs — all content IEs are tagged optional.
	tagged, err := ie.ParseMultiIEs(pay)
	if err != nil {
		return fmt.Errorf("optional: %w", err)
	}

	for _, i := range tagged {
		switch i.Kind {
		case ie.EPSBearerContextStatus:
			if tmp.EPSBearerContextStatus == nil {
				tmp.EPSBearerContextStatus = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.T3448Value:
			if tmp.T3448Value == nil {
				tmp.T3448Value = i
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

// MessageType returns the SERVICE ACCEPT message type constant (0x4F).
func (m *ServiceAccept) MessageType() uint8 {
	return EMMServiceAccept
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *ServiceAccept) MessageTypeName() string {
	return "Service accept"
}

// IsRequest returns false — SERVICE ACCEPT is sent network to UE.
func (m *ServiceAccept) IsRequest() bool {
	return false
}

// ServiceAccept implements the Message interface.
var _ Message = (*ServiceAccept)(nil)

// Header returns the EMM security header embedded in ServiceAccept.
func (m *ServiceAccept) Header() *EmmHeader {
	return &m.EmmHeader
}
