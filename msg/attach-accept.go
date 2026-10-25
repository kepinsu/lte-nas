package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// AttachAccept represents an ATTACH ACCEPT message as defined in TS 24.301
// clause 8.2.1 (Table 8.2.1.1, direction network to UE).
//
// The positional mandatory part encodes EPS attach result and spare half
// octet as a packed half-octet, followed by T3412 value, tracking area
// identity list and ESM message container without IEI tags. Tagged optional
// IEs (e.g. GUTI, EMM cause) are resolved via message context prototypes.
type AttachAccept struct {
	EmmHeader
	EPSAttachResult     *ie.IE // positional type-1 (9.9.3.10), half-octet with SpareHalfOctet
	SpareHalfOctet      *ie.IE // positional type-1 (9.9.2.9), half-octet with EPSAttachResult
	T3412Value          *ie.IE // positional type-1 (9.9.3.16)
	TAIList             *ie.IE // positional LV type-4 (9.9.3.33)
	ESMMessageContainer *ie.IE // positional LV-E type-6 (9.9.3.15), pass-through
	// Tagged optional IEs resolved via prototypes in the spec table.
	GUTI          *ie.IE // tagged optional TLV, 0x50 MobileIdentity (9.9.3.12)
	EMMCause      *ie.IE // tagged optional TV type-3, 0x53 (9.9.3.9)
	T3412Extended *ie.IE // tagged optional TLV-E, 0x5e (9.9.3.16B)
	EMMResult     *ie.IE // tagged optional, placeholder for any non-prototype tag
	IEs           []*ie.IE
}

// NewAttachAccept creates an ATTACH ACCEPT message with the mandatory positional
// IEs from TS 24.301 clause 8.2.1, followed by optional tagged IEs classified
// by their Kind. Fields for recognized tags are resolved first; duplicates and
// unrecognized IEs are appended to the IEs field.
func NewAttachAccept(eppResult, spareHalfOctet, t3412, taiList,
	esmMsgContainer *ie.IE, ies ...*ie.IE,
) *AttachAccept {
	m := &AttachAccept{
		EPSAttachResult:     eppResult,
		SpareHalfOctet:      spareHalfOctet,
		T3412Value:          t3412,
		TAIList:             taiList,
		ESMMessageContainer: esmMsgContainer,
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
		case ie.MobileIdentity:
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
		case ie.T3412ExtendedValue:
			if m.T3412Extended == nil {
				m.T3412Extended = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the ATTACH ACCEPT message, or 0 for a
// nil receiver. The result includes header, packed half-octets (1 byte),
// T3412 value, positional IEs and all tagged IEs in spec-table order.
func (m *AttachAccept) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	// Packed half-octets (EPS attach result low nibble, spare high nibble).
	if m.EPSAttachResult != nil && m.SpareHalfOctet != nil {
		n++ // packed half-octet
	} else {
		if m.EPSAttachResult != nil {
			n++
		}
		if m.SpareHalfOctet != nil {
			n++
		}
	}

	if m.T3412Value != nil {
		n++ // type-1 positional
	}
	if m.TAIList != nil {
		n += m.TAIList.MarshalLen()
	}
	if m.ESMMessageContainer != nil {
		n += m.ESMMessageContainer.MarshalLen()
	}

	// Tagged IEs in spec-table order.
	if m.GUTI != nil {
		n += m.GUTI.MarshalLen()
	}
	if m.EMMCause != nil {
		n += m.EMMCause.MarshalLen()
	}
	if m.T3412Extended != nil {
		n += m.T3412Extended.MarshalLen()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the ATTACH ACCEPT message into b. The nil receiver
// returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too
// short. Header, packed half-octets, positional IEs, named tagged IEs and
// extras follow TS 24.301 clause 8.2.1 table order.
func (m *AttachAccept) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: AttachAccept is nil", ie.ErrNilIE)
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

	// Packed half-octets (EPS attach result low nibble, spare high nibble).
	if m.EPSAttachResult != nil && m.SpareHalfOctet != nil {
		hb, err := ie.MarshalHalfOctets(m.EPSAttachResult, m.SpareHalfOctet)
		if err != nil {
			return fmt.Errorf("half-octets: %w", err)
		}
		if len(b) < n+len(hb) {
			return io.ErrShortBuffer
		}
		copy(b[n:n+len(hb)], hb)
		n += len(hb)
	} else {
		if m.EPSAttachResult != nil {
			if err := m.EPSAttachResult.MarshalToWithoutTag(b[n:]); err != nil {
				return fmt.Errorf("EPSAttachResult: %w", err)
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

	// Positional IEs in spec-table order.
	if m.T3412Value != nil {
		if err := m.T3412Value.MarshalToWithoutTag(b[n:]); err != nil {
			return fmt.Errorf("T3412Value: %w", err)
		}
		n++
	}
	if m.TAIList != nil {
		if err := m.TAIList.MarshalToWithoutTag(b[n:]); err != nil {
			return fmt.Errorf("TAIList: %w", err)
		}
		n += m.TAIList.MarshalLenWithoutTag()
	}
	if m.ESMMessageContainer != nil {
		if err := m.ESMMessageContainer.MarshalToWithoutTag(b[n:]); err != nil {
			return fmt.Errorf("ESMMessageContainer: %w", err)
		}
		n += m.ESMMessageContainer.MarshalLenWithoutTag()
	}

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
	if m.T3412Extended != nil {
		w, err := writeIE(b, n, m.T3412Extended)
		if err != nil {
			return fmt.Errorf("T3412Extended: %w", err)
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

// UnmarshalBinary decodes an ATTACH ACCEPT message from b. The receiver is
// left unchanged on error. Header, positional IEs via ParseHalfOctets +
// ParseWithoutTag prototypes, tagged optional IEs resolved via prototypes.
func (m *AttachAccept) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: AttachAccept is nil", ie.ErrNilIE)
	}

	var tmp AttachAccept

	// Decode header (7 octets for integrity-protected).
	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	// Positional mandatory IEs. Half-octets (EPSAttachResult + SpareHalfOctet)
	// packed in one byte.
	pos := 0

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	lProto := ie.NewPositionalEPSAttachResult(0)
	hProto := ie.NewPositionalSpareHalfOctet(0)
	lIE, hIE, err := ie.ParseHalfOctets(pay[pos:], lProto, hProto)
	if err != nil {
		return fmt.Errorf("half-octets: %w", err)
	}
	tmp.EPSAttachResult = lIE
	tmp.SpareHalfOctet = hIE

	pos++ // consumed 1 byte for packed half-octets

	// Remaining positional IEs via ParseWithoutTag with prototypes.
	tmp.T3412Value, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalGPRSTimer(0))
	if err != nil {
		return fmt.Errorf("T3412Value: %w", err)
	}
	pos++

	tmp.TAIList, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalTrackingAreaIdentityList([]byte{0, 0, 0, 0, 0, 0}))
	if err != nil {
		return fmt.Errorf("TAIList: %w", err)
	}
	pos += tmp.TAIList.MarshalLenWithoutTag()

	tmp.ESMMessageContainer, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalESMMessageContainer([]byte{0, 0, 0}))
	if err != nil {
		return fmt.Errorf("ESMMessageContainer: %w", err)
	}
	pos += tmp.ESMMessageContainer.MarshalLenWithoutTag()

	// Tagged optional IEs — prototypes resolve ambiguous tags.
	tagged, err := ie.ParseMultiIEs(pay[pos:])
	if err != nil {
		return fmt.Errorf("optional: %w", err)
	}

	for _, i := range tagged {
		switch i.Kind {
		case ie.MobileIdentity:
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
		case ie.T3412ExtendedValue:
			if tmp.T3412Extended == nil {
				tmp.T3412Extended = i
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

// MessageType returns the ATTACH ACCEPT message type constant (0x42).
func (m *AttachAccept) MessageType() uint8 {
	return EMMAttachAccept
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *AttachAccept) MessageTypeName() string {
	return "Attach accept"
}

// IsRequest returns false — ATTACH ACCEPT is sent network to UE.
func (m *AttachAccept) IsRequest() bool {
	return false
}

// AttachAccept implements the Message interface.
var _ Message = (*AttachAccept)(nil)

// Header returns the EMM security header embedded in AttachAccept.
func (m *AttachAccept) Header() *EmmHeader {
	return &m.EmmHeader
}
