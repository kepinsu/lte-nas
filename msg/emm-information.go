package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// EMMInformation represents an EMM INFORMATION message as defined in TS 24.301
// clause 8.2.13 (Table 8.2.13.1, dual-direction).
//
// No positional mandatory part beyond the header. Tagged optional IEs include
// SS Service Notification (type-1 TV, tag 0xd9) and ESM Message Container
// (pass-through LV-E type-6, tag 0x78).
type EmmInformation struct {
	EmmHeader
	SSServiceNotification *ie.IE // tagged optional type-1 TV, tag 0xd9
	ESMMessageContainer   *ie.IE // tagged optional TLV-E type-6, clause 9.9.3.15
	IEs                   []*ie.IE
}

// NewEmmInformation creates an EMM INFORMATION message with optional tagged
// IEs classified by their Kind. First occurrence wins; duplicates and
// unrecognized IEs are appended to the IEs field. ESM Message Container tag is
// ie.ESMMessageContainer (0x78); SS Service Notification (tag 0xd9) falls
// through to the generic IEs list since it has no catalogue Kind.
func NewEmmInformation(ies ...*ie.IE) *EmmInformation {
	m := &EmmInformation{}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch {
		case i.Kind&0xf0 == 0xd0: // SS Service Notification (tag 0xd, type-1 TV)
			if m.SSServiceNotification == nil {
				m.SSServiceNotification = i
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

// MarshalLen returns the wire size of the EMM INFORMATION message, or 0 for a
// nil receiver. The result includes header and all tagged IEs in spec-table
// order (SS Service Notification first, then ESM Message Container).
func (m *EmmInformation) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	if m.SSServiceNotification != nil {
		n += m.SSServiceNotification.MarshalLen()
	}
	if m.ESMMessageContainer != nil {
		n += m.ESMMessageContainer.MarshalLen()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the EMM INFORMATION message into b. The nil receiver
// returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too
// short. Writes header then named tagged IEs in spec-table order followed by
// extras.
func (m *EmmInformation) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: EmmInformation is nil", ie.ErrNilIE)
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

	if m.SSServiceNotification != nil {
		w, err := writePassThroughTagged(b, 0, m.SSServiceNotification)
		if err != nil {
			return fmt.Errorf("SSServiceNotification: %w", err)
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

// UnmarshalBinary decodes an EMM INFORMATION message from b. The receiver is
// left unchanged on error (uses local tmp struct). Header via UnmarshalBinary,
// tagged optional IEs parsed from payload using manual tag inspection for
// non-catalogue IEs and ie.Parse for catalogue ones.
func (m *EmmInformation) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: EmmInformation is nil", ie.ErrNilIE)
	}

	var tmp EmmInformation

	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	// All IEs in EMM INFORMATION are optional tagged IEs; empty payload is valid.
	tagged, err := parseEMMInformationTags(pay)
	if err != nil {
		return fmt.Errorf("optional: %w", err)
	}

	for _, i := range tagged {
		switch {
		case i.Kind == ie.ESMMessageContainer: // tag 0x78, type-6 TLV-E
			if tmp.ESMMessageContainer == nil {
				tmp.ESMMessageContainer = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case i.Kind&0xf0 == 0xd0: // SS Service Notification (tag 0xd, type-1 TV)
			if tmp.SSServiceNotification == nil {
				tmp.SSServiceNotification = i
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

// MessageType returns the EMM INFORMATION message type constant (0x61).
func (m *EmmInformation) MessageType() uint8 {
	return EMMInformation
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *EmmInformation) MessageTypeName() string {
	return "EMM information"
}

// IsRequest returns false — EMM INFORMATION is dual-direction but classified
// as network-originated in the Message interface contract.
func (m *EmmInformation) IsRequest() bool {
	return false
}

// EmmInformation implements the Message interface.
var _ Message = (*EmmInformation)(nil)

// parseEMMInformationTags parses tagged optional IEs from the payload of an
// EMM INFORMATION message. Catalogue IE ESM Message Container (tag 0x78) is
// decoded via ie.Parse; non-catalogue IE SS Service Notification (tag 0xd9,
// type-1 TV) is parsed as pass-through.
func parseEMMInformationTags(b []byte) ([]*ie.IE, error) {
	var ies []*ie.IE

	for offset := 0; offset < len(b); {
		tag := b[offset]

		switch tag {
		case 0x78: // ESM Message Container — catalogue type-6 TLV-E
			if i, err := ie.Parse(b[offset:], ie.Kind(tag)); err == nil {
				ies = append(ies, i)
				offset += i.MarshalLen()
				continue
			}
			fallthrough
		default:
			// Unknown tag — try type-4 TLV first, then fall back to type-1 TV
			if len(b[offset:]) >= 2 {
				length := int(b[offset+1])
				if length > 0 && len(b[offset:]) >= 2+length {
					// Looks like type-4 TLV
					payload := make([]byte, length)
					copy(payload, b[offset+2:])
					ies = append(ies, &ie.IE{Kind: ie.Kind(tag), Payload: payload, Format: ie.IEFormatFour})
					offset += 2 + length
					continue
				}
			}
			if len(b[offset:]) >= 1 {
				wire := b[offset]
				ies = append(ies, &ie.IE{Kind: ie.Kind(wire & 0xf0), Payload: []byte{wire & 0x0f}, Format: ie.IEFormatOne})
				offset++
			} else {
				return nil, io.ErrUnexpectedEOF
			}
		}
	}

	return ies, nil
}

// Header returns the EMM security header embedded in EmmInformation.
func (m *EmmInformation) Header() *EmmHeader {
	return &m.EmmHeader
}
