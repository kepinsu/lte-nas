package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// UplinkNASTransport represents an UPLINK NAS TRANSPORT message as defined in
// TS 24.301 clause 8.2.30 (Table 8.2.30.1, direction UE to network).
//
// Tagged optional IEs: S-TMSI (pass-through), ESM Message Container
// (pass-through LV-E type-6, tag 0x78), and EMM Capabilities (type-1 TV,
// tag 0xc1). No positional mandatory part beyond the header.
type UplinkNASTransport struct {
	EmmHeader
	STMSI               *ie.IE // tagged optional TLV type-4, S-TMSI (pass-through)
	ESMMessageContainer *ie.IE // tagged optional TLV-E type-6, clause 9.9.3.15
	EMMCapabilities     *ie.IE // tagged optional type-1 TV, tag 0xc1
	IEs                 []*ie.IE
}

// NewUplinkNASTransport creates an UPLINK NAS TRANSPORT message with optional
// tagged IEs classified by Kind. First occurrence wins; duplicates and
// unrecognized IEs are appended to the IEs field. EMM Capabilities tag is 0xc1
// (type-1 TV); S-TMSI is type-4 LV; ESM Message Container tag 0x78.
func NewUplinkNASTransport(ies ...*ie.IE) *UplinkNASTransport {
	m := &UplinkNASTransport{}

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
		case i.Kind&0xf0 == 0xc0: // EMM Capabilities (tag 0xc, type-1 TV)
			if m.EMMCapabilities == nil {
				m.EMMCapabilities = i
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

// MarshalLen returns the wire size of the UPLINK NAS TRANSPORT message, or 0
// for a nil receiver. The result includes header and all tagged IEs in spec-
// table order.
func (m *UplinkNASTransport) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	if m.STMSI != nil {
		n += m.STMSI.MarshalLen()
	}
	if m.ESMMessageContainer != nil {
		n += m.ESMMessageContainer.MarshalLen()
	}
	if m.EMMCapabilities != nil {
		n += m.EMMCapabilities.MarshalLen()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the UPLINK NAS TRANSPORT message into b. The nil receiver
// returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too
// short. Writes header then named tagged IEs in spec-table order.
func (m *UplinkNASTransport) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: UplinkNASTransport is nil", ie.ErrNilIE)
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

	if m.STMSI != nil {
		w, err := writePassThroughTagged(b, 0, m.STMSI)
		if err != nil {
			return fmt.Errorf("STMSI: %w", err)
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

	if m.EMMCapabilities != nil {
		w, err := writePassThroughTagged(b, 0, m.EMMCapabilities)
		if err != nil {
			return fmt.Errorf("EMMCapabilities: %w", err)
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

// UnmarshalBinary decodes an UPLINK NAS TRANSPORT message from b. The receiver
// is left unchanged on error (uses local tmp struct). Header via UnmarshalBinary,
// tagged optional IEs resolved manually from payload.
func (m *UplinkNASTransport) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: UplinkNASTransport is nil", ie.ErrNilIE)
	}

	var tmp UplinkNASTransport

	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	// All IEs in Uplink NAS Transport are optional tagged IEs; empty payload is valid.
	tagged, err := ie.ParseMultiIEs(pay)
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
		case i.Kind&0xf0 == 0xc0: // EMM Capabilities (tag 0xc, type-1 TV)
			if tmp.EMMCapabilities == nil {
				tmp.EMMCapabilities = i
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

	*m = tmp
	return nil
}

// MessageType returns the UPLINK NAS TRANSPORT message type constant (0x63).
func (m *UplinkNASTransport) MessageType() uint8 {
	return EMMUplinkNASTransport
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *UplinkNASTransport) MessageTypeName() string {
	return "Uplink NAS transport"
}

// IsRequest returns true — UPLINK NAS TRANSPORT is sent UE to network.
func (m *UplinkNASTransport) IsRequest() bool {
	return true
}

// UplinkNASTransport implements the Message interface.
var _ Message = (*UplinkNASTransport)(nil)

// Header returns the EMM security header embedded in UplinkNASTransport.
func (m *UplinkNASTransport) Header() *EmmHeader {
	return &m.EmmHeader
}
