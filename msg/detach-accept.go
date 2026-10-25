package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// DetachAccept represents a DETACH ACCEPT message as defined in TS 24.301
// clause 8.2.10 (Table 8.2.10.1, direction network to UE).
//
// The message uses an integrity-protected header (SHT 0x01..0x05, 7 octets).
// No mandatory positional payload — all content IEs are tagged optional
// (GUTI/EPS mobile identity tag 0x50, EMM cause tag 0x53, T3412 value
// tag 0x5A). See also Table 8.2.10.2 for the UE-terminated detach accept.
type DetachAccept struct {
	EmmHeader
	GUTI       *ie.IE // tagged optional TLV, 0x50 (9.9.3.12)
	EMMCause   *ie.IE // tagged optional TV type-3, 0x53 (9.9.3.9)
	T3412Value *ie.IE // tagged optional TV type-3, 0x5a (9.9.3.16)
	IEs        []*ie.IE
}

// NewDetachAccept creates a DETACH ACCEPT message with optional tagged IEs
// classified by their Kind. Fields for recognized tags are resolved first;
// duplicates and unrecognized IEs are appended to the IEs field.
func NewDetachAccept(ies ...*ie.IE) *DetachAccept {
	m := &DetachAccept{}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
		case ie.EPSMobileIdentity:
			if m.GUTI == nil {
				m.GUTI = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.EMMCause:
			if m.EMMCause == nil {
				m.EMMCause = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.T3412Value:
			if m.T3412Value == nil {
				m.T3412Value = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the DETACH ACCEPT message, or 0 for a
// nil receiver. The result includes header and all named tagged IEs in spec-table
// order. No positional mandatory body for this message.
func (m *DetachAccept) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	// Tagged IEs in spec-table order.
	if m.GUTI != nil {
		n += m.GUTI.MarshalLen()
	}
	if m.EMMCause != nil {
		n += m.EMMCause.MarshalLen()
	}
	if m.T3412Value != nil {
		n += m.T3412Value.MarshalLen()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the DETACH ACCEPT message into b. The nil receiver
// returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too
// short. Header and tagged IEs in spec-table order follow TS 24.301 clause
// 8.2.10.
func (m *DetachAccept) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: DetachAccept is nil", ie.ErrNilIE)
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

	n := hlen

	// Named tagged IEs in spec-table order.
	if m.GUTI != nil {
		w, err := writeIE(b, n, m.GUTI)
		if err != nil {
			return fmt.Errorf("GUTI: %w", err)
		}
		n += w
	}
	if m.EMMCause != nil {
		w, err := writeIE(b, n, m.EMMCause)
		if err != nil {
			return fmt.Errorf("EMMCause: %w", err)
		}
		n += w
	}
	if m.T3412Value != nil {
		w, err := writeIE(b, n, m.T3412Value)
		if err != nil {
			return fmt.Errorf("T3412Value: %w", err)
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

// UnmarshalBinary decodes a DETACH ACCEPT message from b. The receiver is left
// unchanged on error. Header via UnmarshalBinary; no positional part — all IEs
// are tagged optional and resolved via prototypes.
func (m *DetachAccept) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: DetachAccept is nil", ie.ErrNilIE)
	}

	var tmp DetachAccept

	// Decode header (7 octets for integrity-protected).
	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	// No mandatory positional body for Detach Accept (see TS 24.301 Table
	// 8.2.10.1.1 / 8.2.10.2.1). Tagged optional IEs only.
	if len(pay) > 0 {
		tagged, err := ie.ParseMultiIEs(pay)
		if err != nil {
			return fmt.Errorf("optional: %w", err)
		}

		for _, i := range tagged {
			switch i.Kind {
			case ie.EPSMobileIdentity:
				if tmp.GUTI == nil {
					tmp.GUTI = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case ie.EMMCause:
				if tmp.EMMCause == nil {
					tmp.EMMCause = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case ie.T3412Value:
				if tmp.T3412Value == nil {
					tmp.T3412Value = i
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

// MessageType returns the DETACH ACCEPT message type constant (0x46).
func (m *DetachAccept) MessageType() uint8 {
	return EMMDetachAccept
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *DetachAccept) MessageTypeName() string {
	return "Detach accept"
}

// IsRequest returns false — DETACH ACCEPT is sent network to UE.
func (m *DetachAccept) IsRequest() bool {
	return false
}

// DetachAccept implements the Message interface.
var _ Message = (*DetachAccept)(nil)

// Header returns the EMM security header embedded in DetachAccept.
func (m *DetachAccept) Header() *EmmHeader {
	return &m.EmmHeader
}
