package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// DownlinkNASTransport represents a DOWNLINK NAS TRANSPORT message as defined in
// TS 24.301 clause 8.2.12 (Table 8.2.12.1, direction network to UE).
//
// Positional mandatory IEs (in spec order): GUTI (type-3 V, one octet) and TMSI
// (type-4 LV). Tagged optional IEs: S-TMSI (pass-through), LAN Group ID
// (type-1 TV, tag 0xd7), and ESM Message Container (pass-through LV-E type-6).
type DownlinkNASTransport struct {
	EmmHeader
	GUTI                *ie.IE // positional mandatory type-3 (9.9.3.46), one-octet
	TMSI                *ie.IE // positional mandatory type-4, four-octet value
	STMSI               *ie.IE // tagged optional TLV type-4, pass-through
	LANGroupID          *ie.IE // tagged optional type-1 TV, tag 0xd7
	ESMMessageContainer *ie.IE // tagged optional TLV-E type-6, clause 9.9.3.15
	IEs                 []*ie.IE
}

// NewDownlinkNASTransport creates a DOWNLINK NAS TRANSPORT message with the
// mandatory positional GUTI and TMSI IEs from TS 24.301 clause 8.2.12,
// followed by optional tagged IEs classified by Kind. LAN Group ID tag is 0xd7;
// S-TMSI (tag 0xd1) and ESM Message Container (tag 0x78) are resolved. First
// occurrence wins for duplicates; unrecognized IEs go to the IEs field.
func NewDownlinkNASTransport(guti, tmsi *ie.IE, ies ...*ie.IE) *DownlinkNASTransport {
	if guti == nil || tmsi == nil {
		panic("DownlinkNASTransport GUTI and TMSI must not be nil")
	}

	m := &DownlinkNASTransport{GUTI: guti, TMSI: tmsi}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch {
		case i.Kind == ie.Kind(0xd1): // S-TMSI (tag 0xd1, type-4 TLV)
			if m.STMSI == nil {
				m.STMSI = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case i.Kind&0xf0 == 0xd0: // LAN Group ID (tag 0xd, type-1 TV)
			if m.LANGroupID == nil {
				m.LANGroupID = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case i.Kind == ie.ESMMessageContainer: // tag 0x78, type-6 TLV-E
			if m.ESMMessageContainer == nil {
				m.ESMMessageContainer = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the DOWNLINK NAS TRANSPORT message, or 0
// for a nil receiver. The result includes header plus positional GUTI (1 byte)
// and TMSI, then tagged IEs in spec-table order.
func (m *DownlinkNASTransport) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	n += m.GUTI.MarshalLenWithoutTag()
	n += m.TMSI.MarshalLenWithoutTag()

	if m.STMSI != nil {
		n += m.STMSI.MarshalLen()
	}
	if m.LANGroupID != nil {
		n += m.LANGroupID.MarshalLen()
	}
	if m.ESMMessageContainer != nil {
		n += m.ESMMessageContainer.MarshalLen()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the DOWNLINK NAS TRANSPORT message into b. The nil receiver
// returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too
// short. Writes header, positional GUTI and TMSI, then tagged IEs in order.
func (m *DownlinkNASTransport) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: DownlinkNASTransport is nil", ie.ErrNilIE)
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

	if m.GUTI == nil {
		return io.ErrUnexpectedEOF
	}
	if err := m.GUTI.MarshalToWithoutTag(b); err != nil {
		return fmt.Errorf("GUTI: %w", err)
	}
	b = b[m.GUTI.MarshalLenWithoutTag():]

	if m.TMSI == nil {
		return io.ErrUnexpectedEOF
	}
	tmsiLen := m.TMSI.MarshalLenWithoutTag()
	if err := m.TMSI.MarshalToWithoutTag(b); err != nil {
		return fmt.Errorf("TMSI: %w", err)
	}
	b = b[tmsiLen:]

	if m.STMSI != nil {
		w, err := writePassThroughTagged(b, 0, m.STMSI)
		if err != nil {
			return fmt.Errorf("STMSI: %w", err)
		}
		b = b[w:]
	}

	if m.LANGroupID != nil {
		w, err := writePassThroughTagged(b, 0, m.LANGroupID)
		if err != nil {
			return fmt.Errorf("LANGroupID: %w", err)
		}
		b = b[w:]
	}

	if m.ESMMessageContainer != nil {
		w, err := writeIE(b, 0, m.ESMMessageContainer)
		if err != nil {
			return fmt.Errorf("ESMMessageContainer: %w", err)
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

// UnmarshalBinary decodes a DOWNLINK NAS TRANSPORT message from b. The receiver
// is left unchanged on error (uses local tmp struct). Header via UnmarshalBinary,
// positional GUTI and TMSI via ParseWithoutTag with prototypes, tagged IEs
// resolved manually.
func (m *DownlinkNASTransport) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: DownlinkNASTransport is nil", ie.ErrNilIE)
	}

	var tmp DownlinkNASTransport

	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	pos := 0

	tmp.GUTI, err = ie.ParseWithoutTag(pay[pos:], ie.NewRaw(ie.NoTag, ie.IEFormatThree, []byte{0}))
	if err != nil {
		return fmt.Errorf("GUTI: %w", err)
	}
	pos += tmp.GUTI.MarshalLenWithoutTag()

	// TMSI is positional type-4 (LV format): wire has length byte + 4-octet value.
	// Prototype uses type-4 with explicit Length so ParseWithoutTag reads length byte from wire.
	tmp.TMSI, err = ie.ParseWithoutTag(pay[pos:], ie.NewRaw(ie.NoTag, ie.IEFormatFour, []byte{0}))
	if err != nil {
		return fmt.Errorf("TMSI: %w", err)
	}
	pos += tmp.TMSI.MarshalLenWithoutTag()

	if len(pay) > pos {
		tagged, err := ie.ParseMultiIEs(pay[pos:])
		if err != nil {
			return fmt.Errorf("optional: %w", err)
		}

		for _, i := range tagged {
			switch {
			case i.Kind == ie.Kind(0xd1): // S-TMSI (tag 0xd1, type-4 TLV)
				if tmp.STMSI == nil {
					tmp.STMSI = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case i.Kind&0xf0 == 0xd0: // LAN Group ID (tag 0xd, type-1 TV)
				if tmp.LANGroupID == nil {
					tmp.LANGroupID = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case i.Kind == ie.ESMMessageContainer: // tag 0x78, type-6 TLV-E
				if tmp.ESMMessageContainer == nil {
					tmp.ESMMessageContainer = i
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

// MessageType returns the DOWNLINK NAS TRANSPORT message type constant (0x62).
func (m *DownlinkNASTransport) MessageType() uint8 {
	return EMMDownlinkNASTransport
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *DownlinkNASTransport) MessageTypeName() string {
	return "Downlink NAS transport"
}

// IsRequest returns false — DOWNLINK NAS TRANSPORT is sent network to UE.
func (m *DownlinkNASTransport) IsRequest() bool {
	return false
}

// DownlinkNASTransport implements the Message interface.
var _ Message = (*DownlinkNASTransport)(nil)

// Header returns the EMM security header embedded in DownlinkNASTransport.
func (m *DownlinkNASTransport) Header() *EmmHeader {
	return &m.EmmHeader
}
