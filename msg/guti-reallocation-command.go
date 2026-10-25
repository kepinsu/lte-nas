package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// GUTIReallocationCommand represents a GUTI REALLOCATION COMMAND message as
// defined in TS 24.301 clause 8.2.16 (Table 8.2.16.1, direction network to UE).
//
// The positional mandatory part encodes GUTI as a plain LV (no IEI). Tagged
// optional IEs (TAI list, DCN-ID, UE radio capability ID, UE radio capability
// ID deletion indication) are resolved via prototypes. Integrity-protected
// header (SHT 0x01..0x05, 7 octets).
type GUTIReallocationCommand struct {
	EmmHeader
	GUTI                           *ie.IE // positional LV type-4 (9.9.3.12)
	TAIList                        *ie.IE // tagged optional TLV, 0x54 (9.9.3.33)
	DCNID                          *ie.IE // tagged optional TLV, 0x65 (9.9.3.48)
	UERadioCapabilityID            *ie.IE // tagged optional TLV, 0x66 (9.9.3.60)
	UERadioCapabilityIDDeletionInd *ie.IE // tagged optional TV type-1, 0xB (9.9.3.61)
	IEs                            []*ie.IE
}

// NewGUTIReallocationCommand creates a GUTI REALLOCATION COMMAND message with
// the mandatory positional IE from TS 24.301 clause 8.2.16, followed by optional
// tagged IEs classified by their Kind. Fields for recognized tags are resolved
// first; duplicates and unrecognized IEs are appended to the IEs field.
func NewGUTIReallocationCommand(guti *ie.IE, ies ...*ie.IE) *GUTIReallocationCommand {
	if guti == nil {
		panic("GUTIReallocationCommand positional IEs must not be nil")
	}

	m := &GUTIReallocationCommand{
		GUTI: guti,
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
		case ie.TAIList:
			if m.TAIList == nil {
				m.TAIList = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.DCNID:
			if m.DCNID == nil {
				m.DCNID = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.UERadioCapabilityID:
			if m.UERadioCapabilityID == nil {
				m.UERadioCapabilityID = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.UERadioCapabilityIDDeletionIndication:
			if m.UERadioCapabilityIDDeletionInd == nil {
				m.UERadioCapabilityIDDeletionInd = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the GUTI REALLOCATION COMMAND message, or
// 0 for a nil receiver. The result includes header, positional GUTI (length+value)
// and all named tagged IEs in spec-table order.
func (m *GUTIReallocationCommand) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	if m.GUTI != nil {
		n += m.GUTI.MarshalLenWithoutTag()
	}
	if m.TAIList != nil {
		n += m.TAIList.MarshalLen()
	}
	if m.DCNID != nil {
		n += m.DCNID.MarshalLen()
	}
	if m.UERadioCapabilityID != nil {
		n += m.UERadioCapabilityID.MarshalLen()
	}
	if m.UERadioCapabilityIDDeletionInd != nil {
		n += m.UERadioCapabilityIDDeletionInd.MarshalLen()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the GUTI REALLOCATION COMMAND message into b. The nil
// receiver returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b
// is too short. Header, positional GUTI and tagged IEs in spec-table order
// follow TS 24.301 clause 8.2.16 table order.
func (m *GUTIReallocationCommand) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: GUTIReallocationCommand is nil", ie.ErrNilIE)
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

	if m.GUTI != nil {
		if err := m.GUTI.MarshalToWithoutTag(b[n:]); err != nil {
			return fmt.Errorf("GUTI: %w", err)
		}
		n += m.GUTI.MarshalLenWithoutTag()
	}

	if m.TAIList != nil {
		w, err := writeIE(b, n, m.TAIList)
		if err != nil {
			return fmt.Errorf("TAIList: %w", err)
		}
		n += w
	}

	if m.DCNID != nil {
		w, err := writeIE(b, n, m.DCNID)
		if err != nil {
			return fmt.Errorf("DCNID: %w", err)
		}
		n += w
	}

	if m.UERadioCapabilityID != nil {
		w, err := writeIE(b, n, m.UERadioCapabilityID)
		if err != nil {
			return fmt.Errorf("UERadioCapabilityID: %w", err)
		}
		n += w
	}

	if m.UERadioCapabilityIDDeletionInd != nil {
		w, err := writeIE(b, n, m.UERadioCapabilityIDDeletionInd)
		if err != nil {
			return fmt.Errorf("UERadioCapabilityIDDeletionInd: %w", err)
		}
		n += w
	}

	for _, i := range m.IEs {
		w, err := writeIE(b, n, i)
		if err != nil {
			return fmt.Errorf("extra IE: %w", err)
		}
		n += w
	}

	return nil
}

// UnmarshalBinary decodes a GUTI REALLOCATION COMMAND message from b. The
// receiver is left unchanged on error. Header, positional GUTI via
// ParseWithoutTag, tagged optional IEs resolved via prototypes for ambiguous
// tags in TS 24.301 Table 8.2.16.1 (none are ambiguous — all have distinct IEI).
func (m *GUTIReallocationCommand) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: GUTIReallocationCommand is nil", ie.ErrNilIE)
	}

	var tmp GUTIReallocationCommand

	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	pos := 0

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	tmp.GUTI, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalEPSMobileIdentity([]byte{0, 0}))
	if err != nil {
		return fmt.Errorf("GUTI: %w", err)
	}
	pos += tmp.GUTI.MarshalLenWithoutTag()

	if len(pay[pos:]) > 0 {
		tagged, err := ie.ParseMultiIEs(pay[pos:])
		if err != nil {
			return fmt.Errorf("optional: %w", err)
		}

		for _, i := range tagged {
			switch i.Kind {
			case ie.TAIList:
				if tmp.TAIList == nil {
					tmp.TAIList = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case ie.DCNID:
				if tmp.DCNID == nil {
					tmp.DCNID = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case ie.UERadioCapabilityID:
				if tmp.UERadioCapabilityID == nil {
					tmp.UERadioCapabilityID = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case ie.UERadioCapabilityIDDeletionIndication:
				if tmp.UERadioCapabilityIDDeletionInd == nil {
					tmp.UERadioCapabilityIDDeletionInd = i
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

// MessageType returns the GUTI REALLOCATION COMMAND message type constant (0x50).
func (m *GUTIReallocationCommand) MessageType() uint8 {
	return EMMGUTIReallocationCommand
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *GUTIReallocationCommand) MessageTypeName() string {
	return "Guti reallocation command"
}

// IsRequest returns false — GUTI REALLOCATION COMMAND is sent network to UE.
func (m *GUTIReallocationCommand) IsRequest() bool {
	return false
}

// GUTIReallocationCommand implements the Message interface.
var _ Message = (*GUTIReallocationCommand)(nil)

// Header returns the EMM security header embedded in GUTIReallocationCommand.
func (m *GUTIReallocationCommand) Header() *EmmHeader {
	return &m.EmmHeader
}
