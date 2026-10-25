package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// TrackingAreaUpdateReject represents a TRACKING AREA UPDATE REJECT message as
// defined in TS 24.301 clause 8.2.28 (Table 8.2.28.1, direction network to UE).
//
// The positional mandatory part encodes EMM cause as a single-octet type (V).
// All other content IEs are tagged optional (T3346 value 0x5f, Extended EMM
// cause 0xa0, lower bound timer value 0x1c, forbidden TAI lists 0x1d/0x1e).
// Integrity-protected header (SHT 0x01..0x05, 7 octets).
type TrackingAreaUpdateReject struct {
	EmmHeader
	EMMCause *ie.IE // positional type-3 (9.9.3.9)
	// Tagged optional IEs resolved via prototypes in the spec table.
	T3346Value       *ie.IE // tagged optional TLV, 0x5f (9.9.3.16A)
	ExtendedEMMCause *ie.IE // tagged optional TV type-3, 0xa0 (9.9.3.26A)
	LowerBoundTimer  *ie.IE // tagged optional TLV, 0x1c (9.9.3.16B)
	ForbiddenTAI     *ie.IE // tagged optional TLV, 0x1d/0x1e (9.9.3.33)
	IEs              []*ie.IE
}

// NewTrackingAreaUpdateReject creates a TRACKING AREA UPDATE REJECT message with
// the mandatory positional IEs from TS 24.301 clause 8.2.28, followed by optional
// tagged IEs classified by their Kind. Fields for recognized tags are resolved
// first; duplicates and unrecognized IEs are appended to the IEs field.
func NewTrackingAreaUpdateReject(emmCause *ie.IE, ies ...*ie.IE) *TrackingAreaUpdateReject {
	if emmCause == nil {
		panic("TrackingAreaUpdateReject positional IE EMMCause must not be nil")
	}

	m := &TrackingAreaUpdateReject{
		EMMCause: emmCause,
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
		case ie.T3346Value:
			if m.T3346Value == nil {
				m.T3346Value = i
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
		case ie.ForbiddenTAIsForRoaming:
			if m.ForbiddenTAI == nil {
				m.ForbiddenTAI = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.ForbiddenTAIsForRegionalProvisionOfService:
			if m.ForbiddenTAI == nil {
				m.ForbiddenTAI = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the TRACKING AREA UPDATE REJECT message, or
// 0 for a nil receiver. The result includes header, positional mandatory EMM cause,
// and all named tagged IEs in spec-table order.
func (m *TrackingAreaUpdateReject) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	if m.EMMCause != nil {
		n += m.EMMCause.MarshalLenWithoutTag()
	}

	// Tagged IEs in spec-table order.
	if m.T3346Value != nil {
		n += m.T3346Value.MarshalLen()
	}
	if m.ExtendedEMMCause != nil {
		n += m.ExtendedEMMCause.MarshalLen()
	}
	if m.LowerBoundTimer != nil {
		n += m.LowerBoundTimer.MarshalLen()
	}
	if m.ForbiddenTAI != nil {
		n += m.ForbiddenTAI.MarshalLen()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the TRACKING AREA UPDATE REJECT message into b. The nil
// receiver returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b
// is too short. Header, positional EMM cause, tagged IEs in spec-table order
// and extras follow TS 24.301 clause 8.2.28 table order.
func (m *TrackingAreaUpdateReject) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: TrackingAreaUpdateReject is nil", ie.ErrNilIE)
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

	// Positional mandatory IEs.
	if m.EMMCause != nil {
		if err := m.EMMCause.MarshalToWithoutTag(b[n:]); err != nil {
			return fmt.Errorf("EMMCause: %w", err)
		}
		n += m.EMMCause.MarshalLenWithoutTag()
	}

	// Named tagged IEs in spec-table order.
	if m.T3346Value != nil {
		w, err := writeIE(b, n, m.T3346Value)
		if err != nil {
			return fmt.Errorf("T3346Value: %w", err)
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
	if m.ForbiddenTAI != nil {
		w, err := writeIE(b, n, m.ForbiddenTAI)
		if err != nil {
			return fmt.Errorf("ForbiddenTAI: %w", err)
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

// UnmarshalBinary decodes a TRACKING AREA UPDATE REJECT message from b. The
// receiver is left unchanged on error. Header via UnmarshalBinary, positional
// EMM cause via ParseWithoutTag, tagged optional IEs resolved via prototypes.
func (m *TrackingAreaUpdateReject) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: TrackingAreaUpdateReject is nil", ie.ErrNilIE)
	}

	var tmp TrackingAreaUpdateReject

	// Decode header (7 octets for integrity-protected).
	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	pos := 0

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	// Positional mandatory: EMM cause (1 octet).
	tmp.EMMCause, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalEMMCause(0))
	if err != nil {
		return fmt.Errorf("EMMCause: %w", err)
	}
	pos += tmp.EMMCause.MarshalLenWithoutTag()

	// Tagged optional IEs — no ambiguous tags for TAU reject.
	if len(pay[pos:]) > 0 {
		tagged, err := ie.ParseMultiIEs(pay[pos:])
		if err != nil {
			return fmt.Errorf("optional: %w", err)
		}

		for _, i := range tagged {
			switch i.Kind {
			case ie.T3346Value:
				if tmp.T3346Value == nil {
					tmp.T3346Value = i
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
			case ie.ForbiddenTAIsForRoaming:
				if tmp.ForbiddenTAI == nil {
					tmp.ForbiddenTAI = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case ie.ForbiddenTAIsForRegionalProvisionOfService:
				if tmp.ForbiddenTAI == nil {
					tmp.ForbiddenTAI = i
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

// MessageType returns the TRACKING AREA UPDATE REJECT message type constant (0x4B).
func (m *TrackingAreaUpdateReject) MessageType() uint8 {
	return EMMTrackingAreaUpdateReject
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *TrackingAreaUpdateReject) MessageTypeName() string {
	return "Tracking area update reject"
}

// IsRequest returns false — TRACKING AREA UPDATE REJECT is sent network to UE.
func (m *TrackingAreaUpdateReject) IsRequest() bool {
	return false
}

// TrackingAreaUpdateReject implements the Message interface.
var _ Message = (*TrackingAreaUpdateReject)(nil)

// Header returns the EMM security header embedded in TrackingAreaUpdateReject.
func (m *TrackingAreaUpdateReject) Header() *EmmHeader {
	return &m.EmmHeader
}
