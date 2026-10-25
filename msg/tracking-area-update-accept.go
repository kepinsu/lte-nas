package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// TrackingAreaUpdateAccept represents a TRACKING AREA UPDATE ACCEPT message as
// defined in TS 24.301 clause 8.2.26 (Table 8.2.26.1, direction network to UE).
//
// The positional mandatory part encodes EPS update result and spare half octet as
// packed half-octets (wire byte). No additional mandatory positional IEs — all
// remaining fields are tagged optional (T3412 value 0x5a, GUTI/EPS mobile
// identity 0x50, TAI list 0x54, etc.). Integrity-protected header (SHT 0x01..0x05,
// 7 octets).
type TrackingAreaUpdateAccept struct {
	EmmHeader
	EPSUpdateResult *ie.IE // positional type-1 (9.9.3.13), half-octet with SpareHalfOctet
	SpareHalfOctet  *ie.IE // positional type-1 (9.9.2.9), half-octet with EPSUpdateResult
	// Tagged optional IEs resolved via prototypes in the spec table.
	T3412Value                 *ie.IE // tagged optional TV type-3, 0x5a (9.9.3.16)
	GUTI                       *ie.IE // tagged optional TLV, 0x50 (9.9.3.12)
	TAIList                    *ie.IE // tagged optional TLV, 0x54 (9.9.3.33)
	EPSBearerContextStatus     *ie.IE // tagged optional TLV, 0x57 (9.9.2.1)
	EMMCause                   *ie.IE // tagged optional TV type-3, 0x53 (9.9.3.9)
	LocationAreaIdentification *ie.IE // tagged optional TV type-3, 0x13 (9.9.2.2)
	MSIdentity                 *ie.IE // tagged optional TLV, 0x23 (9.9.2.3)
	IEs                        []*ie.IE
}

// NewTrackingAreaUpdateAccept creates a TRACKING AREA UPDATE ACCEPT message with
// the mandatory positional IEs from TS 24.301 clause 8.2.26, followed by optional
// tagged IEs classified by their Kind. Fields for recognized tags are resolved
// first; duplicates and unrecognized IEs are appended to the IEs field.
func NewTrackingAreaUpdateAccept(epsUpdateResult, spareHalfOctet *ie.IE,
	ies ...*ie.IE,
) *TrackingAreaUpdateAccept {
	if epsUpdateResult == nil || spareHalfOctet == nil {
		panic("TrackingAreaUpdateAccept positional IEs must not be nil")
	}

	m := &TrackingAreaUpdateAccept{
		EPSUpdateResult: epsUpdateResult,
		SpareHalfOctet:  spareHalfOctet,
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
		case ie.T3412Value:
			if m.T3412Value == nil {
				m.T3412Value = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.EPSMobileIdentity:
			if m.GUTI == nil {
				m.GUTI = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.TAIList:
			if m.TAIList == nil {
				m.TAIList = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.EPSBearerContextStatus:
			if m.EPSBearerContextStatus == nil {
				m.EPSBearerContextStatus = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.EMMCause:
			if m.EMMCause == nil {
				m.EMMCause = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.LocationAreaIdentification:
			if m.LocationAreaIdentification == nil {
				m.LocationAreaIdentification = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.MobileIdentity:
			if m.MSIdentity == nil {
				m.MSIdentity = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the TRACKING AREA UPDATE ACCEPT message, or
// 0 for a nil receiver. The result includes header, packed half-octets (1 byte),
// and all named tagged IEs in spec-table order.
func (m *TrackingAreaUpdateAccept) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	// Packed half-octets (EPS update result low nibble, spare high nibble).
	if m.EPSUpdateResult != nil && m.SpareHalfOctet != nil {
		n++ // packed half-octet
	} else {
		if m.EPSUpdateResult != nil {
			n++
		}
		if m.SpareHalfOctet != nil {
			n++
		}
	}

	// Tagged IEs in spec-table order.
	if m.T3412Value != nil {
		n += m.T3412Value.MarshalLen()
	}
	if m.GUTI != nil {
		n += m.GUTI.MarshalLen()
	}
	if m.TAIList != nil {
		n += m.TAIList.MarshalLen()
	}
	if m.EPSBearerContextStatus != nil {
		n += m.EPSBearerContextStatus.MarshalLen()
	}
	if m.EMMCause != nil {
		n += m.EMMCause.MarshalLen()
	}
	if m.LocationAreaIdentification != nil {
		n += m.LocationAreaIdentification.MarshalLen()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the TRACKING AREA UPDATE ACCEPT message into b. The nil
// receiver returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b
// is too short. Header, packed half-octets, and tagged IEs in spec-table order
// follow TS 24.301 clause 8.2.26 table order.
func (m *TrackingAreaUpdateAccept) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: TrackingAreaUpdateAccept is nil", ie.ErrNilIE)
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

	// Packed half-octets (EPS update result low nibble, spare high nibble).
	if m.EPSUpdateResult != nil && m.SpareHalfOctet != nil {
		hb, err := ie.MarshalHalfOctets(m.EPSUpdateResult, m.SpareHalfOctet)
		if err != nil {
			return fmt.Errorf("half-octets: %w", err)
		}
		if len(b) < n+len(hb) {
			return io.ErrShortBuffer
		}
		copy(b[n:n+len(hb)], hb)
		n += len(hb)
	} else {
		if m.EPSUpdateResult != nil {
			if err := m.EPSUpdateResult.MarshalToWithoutTag(b[n:]); err != nil {
				return fmt.Errorf("EPSUpdateResult: %w", err)
			}
			n++
		}
		if m.SpareHalfOctet != nil {
			if err := m.SpareHalfOctet.MarshalToWithoutTag(b[n:]); err != nil {
				return fmt.Errorf("SpareHalfOctet: %w", err)
			}
			n++
		}
	}

	// Named tagged IEs in spec-table order.
	if m.T3412Value != nil {
		w, err := writeIE(b, n, m.T3412Value)
		if err != nil {
			return fmt.Errorf("T3412Value: %w", err)
		}
		n += w
	}
	if m.GUTI != nil {
		w, err := writeIE(b, n, m.GUTI)
		if err != nil {
			return fmt.Errorf("GUTI: %w", err)
		}
		n += w
	}
	if m.TAIList != nil {
		w, err := writeIE(b, n, m.TAIList)
		if err != nil {
			return fmt.Errorf("TAIList: %w", err)
		}
		n += w
	}
	if m.EPSBearerContextStatus != nil {
		w, err := writeIE(b, n, m.EPSBearerContextStatus)
		if err != nil {
			return fmt.Errorf("EPSBearerContextStatus: %w", err)
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
	if m.LocationAreaIdentification != nil {
		w, err := writeIE(b, n, m.LocationAreaIdentification)
		if err != nil {
			return fmt.Errorf("LocationAreaIdentification: %w", err)
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

// UnmarshalBinary decodes a TRACKING AREA UPDATE ACCEPT message from b. The
// receiver is left unchanged on error. Header, half-octets via ParseHalfOctets,
// tagged optional IEs resolved via prototypes. No positional mandatory part
// beyond the half-octets (see TS 24.301 Table 8.2.26.1).
func (m *TrackingAreaUpdateAccept) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: TrackingAreaUpdateAccept is nil", ie.ErrNilIE)
	}

	var tmp TrackingAreaUpdateAccept

	// Decode header (7 octets for integrity-protected).
	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	pos := 0

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	// Positional mandatory: half-octets (EPSUpdateResult low nibble,
	// SpareHalfOctet high nibble) packed in one wire byte.
	lProto := ie.NewPositionalEPSAttachResult(0)
	hProto := ie.NewPositionalSpareHalfOctet(0)
	lIE, hIE, err := ie.ParseHalfOctets(pay[pos:], lProto, hProto)
	if err != nil {
		return fmt.Errorf("half-octets: %w", err)
	}
	tmp.EPSUpdateResult = lIE
	tmp.SpareHalfOctet = hIE

	pos++ // consumed 1 byte for packed half-octets

	// Tagged optional IEs — no positional mandatory IEs after the header per
	// TS 24.301 Table 8.2.26.1 (all content IEs are tagged optional).
	if len(pay[pos:]) > 0 {
		tagged, err := ie.ParseMultiIEs(pay[pos:])
		if err != nil {
			return fmt.Errorf("optional: %w", err)
		}

		for _, i := range tagged {
			switch i.Kind {
			case ie.T3412Value:
				if tmp.T3412Value == nil {
					tmp.T3412Value = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case ie.EPSMobileIdentity:
				if tmp.GUTI == nil {
					tmp.GUTI = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case ie.TAIList:
				if tmp.TAIList == nil {
					tmp.TAIList = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case ie.EPSBearerContextStatus:
				if tmp.EPSBearerContextStatus == nil {
					tmp.EPSBearerContextStatus = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case ie.EMMCause:
				if tmp.EMMCause == nil {
					tmp.EMMCause = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case ie.LocationAreaIdentification:
				if tmp.LocationAreaIdentification == nil {
					tmp.LocationAreaIdentification = i
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

// MessageType returns the TRACKING AREA UPDATE ACCEPT message type constant (0x49).
func (m *TrackingAreaUpdateAccept) MessageType() uint8 {
	return EMMTrackingAreaUpdateAccept
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *TrackingAreaUpdateAccept) MessageTypeName() string {
	return "Tracking area update accept"
}

// IsRequest returns false — TRACKING AREA UPDATE ACCEPT is sent network to UE.
func (m *TrackingAreaUpdateAccept) IsRequest() bool {
	return false
}

// TrackingAreaUpdateAccept implements the Message interface.
var _ Message = (*TrackingAreaUpdateAccept)(nil)

// Header returns the EMM security header embedded in TrackingAreaUpdateAccept.
func (m *TrackingAreaUpdateAccept) Header() *EmmHeader {
	return &m.EmmHeader
}
