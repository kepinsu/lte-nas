package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// AttachReject represents an ATTACH REJECT message as defined in TS 24.301
// clause 8.2.3 (Table 8.2.3.1, direction network to UE).
//
// The single positional mandatory IE (EMM cause, type-1) follows the header.
// Tagged optional IEs (ESM message container, timers, extended cause) are
// resolved via ie.ParseMultiIEs with no prototypes — all tagged optional IEs
// for this message are globally unambiguous.
type AttachReject struct {
	EmmHeader
	EMMCause            *ie.IE // positional type-1 (9.9.3.9)
	ESMMessageContainer *ie.IE // optional TLV-E, 0x78 (9.9.3.15), pass-through
	T3346Value          *ie.IE // optional TLV, 0x5f (9.9.3.16A)
	T3402Value          *ie.IE // optional TLV, 0x16 (9.9.3.16A)
	ExtendedEMMCause    *ie.IE // optional TV type-1, 0xa0 (9.9.3.26A)
	LowerBoundTimer     *ie.IE // optional TLV, 0x1c (9.9.3.16B)
	IEs                 []*ie.IE
}

// NewAttachReject creates an ATTACH REJECT message with the mandatory EMM cause
// IE followed by optional tagged IEs classified by their Kind. Duplicate or
// unrecognized IEs are appended to the IEs field.
func NewAttachReject(cause *ie.IE, ies ...*ie.IE) *AttachReject {
	m := &AttachReject{EMMCause: cause}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
		case ie.ESMMessageContainer:
			if m.ESMMessageContainer == nil {
				m.ESMMessageContainer = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.T3346Value:
			if m.T3346Value == nil {
				m.T3346Value = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.T3402ValueForAttachReject:
			if m.T3402Value == nil {
				m.T3402Value = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.ExtendedEMMCause:
			if m.ExtendedEMMCause == nil {
				m.ExtendedEMMCause = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.LowerBoundTimerValue:
			if m.LowerBoundTimer == nil {
				m.LowerBoundTimer = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the ATTACH REJECT message, or 0 for a
// nil receiver. The result includes header, positional EMM cause IE and all
// tagged IEs in spec-table order.
func (m *AttachReject) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	if m.EMMCause != nil {
		n++ // type-1 positional = 1 byte
	}

	// Tagged IEs in spec-table order.
	if m.ESMMessageContainer != nil {
		n += m.ESMMessageContainer.MarshalLen()
	}
	if m.T3346Value != nil {
		n += m.T3346Value.MarshalLen()
	}
	if m.T3402Value != nil {
		n += m.T3402Value.MarshalLen()
	}
	if m.ExtendedEMMCause != nil {
		n += m.ExtendedEMMCause.MarshalLen()
	}
	if m.LowerBoundTimer != nil {
		n += m.LowerBoundTimer.MarshalLen()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the ATTACH REJECT message into b. The nil receiver
// returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too
// short. Header, positional EMM cause IE and all tagged IEs follow spec-table order.
func (m *AttachReject) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: AttachReject is nil", ie.ErrNilIE)
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
	n := 0

	// Positional EMM cause IE.
	if m.EMMCause != nil {
		if err := m.EMMCause.MarshalToWithoutTag(b[n:]); err != nil {
			return fmt.Errorf("EMMCause: %w", err)
		}
		n++
	}

	// Tagged IEs in spec-table order.
	if m.ESMMessageContainer != nil {
		w, err := writeIE(b, n, m.ESMMessageContainer)
		if err != nil {
			return fmt.Errorf("ESMMessageContainer: %w", err)
		}
		n += w
	}
	if m.T3346Value != nil {
		w, err := writeIE(b, n, m.T3346Value)
		if err != nil {
			return fmt.Errorf("T3346Value: %w", err)
		}
		n += w
	}
	if m.T3402Value != nil {
		w, err := writeIE(b, n, m.T3402Value)
		if err != nil {
			return fmt.Errorf("T3402Value: %w", err)
		}
		n += w
	}
	if m.ExtendedEMMCause != nil {
		w, err := writeIE(b, n, m.ExtendedEMMCause)
		if err != nil {
			return fmt.Errorf("ExtendedEMMCause: %w", err)
		}
		n += w
	}
	if m.LowerBoundTimer != nil {
		w, err := writeIE(b, n, m.LowerBoundTimer)
		if err != nil {
			return fmt.Errorf("LowerBoundTimer: %w", err)
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

// UnmarshalBinary decodes an ATTACH REJECT message from b. The receiver is
// left unchanged on error. Header first, positional EMM cause IE via
// ParseWithoutTag prototype, then tagged optional IEs (all unambiguous for this message).
func (m *AttachReject) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: AttachReject is nil", ie.ErrNilIE)
	}

	var tmp AttachReject

	// Decode header (7 octets for integrity-protected).
	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	// Positional mandatory IE via ParseWithoutTag prototype.
	pos := 0

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	tmp.EMMCause, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalEMMCause(0))
	if err != nil {
		return fmt.Errorf("EMMCause: %w", err)
	}
	pos++

	// Tagged optional IEs — no prototypes needed (all tags are unambiguous).
	tagged, err := ie.ParseMultiIEs(pay[pos:])
	if err != nil {
		return fmt.Errorf("optional: %w", err)
	}

	for _, i := range tagged {
		switch i.Kind {
		case ie.ESMMessageContainer:
			if tmp.ESMMessageContainer == nil {
				tmp.ESMMessageContainer = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.T3346Value:
			if tmp.T3346Value == nil {
				tmp.T3346Value = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.T3402ValueForAttachReject:
			if tmp.T3402Value == nil {
				tmp.T3402Value = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.ExtendedEMMCause:
			if tmp.ExtendedEMMCause == nil {
				tmp.ExtendedEMMCause = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.LowerBoundTimerValue:
			if tmp.LowerBoundTimer == nil {
				tmp.LowerBoundTimer = i
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

// MessageType returns the ATTACH REJECT message type constant (0x44).
func (m *AttachReject) MessageType() uint8 {
	return EMMAttachReject
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *AttachReject) MessageTypeName() string {
	return "Attach reject"
}

// IsRequest returns false — ATTACH REJECT is sent network to UE.
func (m *AttachReject) IsRequest() bool {
	return false
}

// AttachReject implements the Message interface.
var _ Message = (*AttachReject)(nil)

// Header returns the EMM security header embedded in AttachReject.
func (m *AttachReject) Header() *EmmHeader {
	return &m.EmmHeader
}
