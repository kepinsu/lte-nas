package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// SecurityModeComplete represents a SECURITY MODE COMPLETE message as defined in
// TS 24.301 clause 8.2.21 (Table 8.2.21.1, direction UE to network).
//
// There are no positional mandatory IEs beyond the header. All content IEs are
// tagged optional in the order shown in TS 24.301 Table 8.2.21.1. The security
// header type is 0x00 (plain), yielding a 2-octet header per TS 24.301 clause
// 9.3.1 NOTE 1.
type SecurityModeComplete struct {
	EmmHeader
	IMEISV                      *ie.IE // tagged optional TLV, 0x23 (9.9.2.3)
	ReplayedNASMessageContainer *ie.IE // tagged optional TLV-E, 0x79 (9.9.3.51), pass-through
	UERadioCapabilityID         *ie.IE // tagged optional TLV, 0x66 (9.9.3.60)
	UECoarseLocationInformation *ie.IE // tagged optional TLV, 0x67 (9.9.3.72)
	IEs                         []*ie.IE
}

// NewSecurityModeComplete creates a SECURITY MODE COMPLETE message with optional
// tagged IEs classified by their Kind. Since the spec table has no positional
// mandatory content IEs (TS 24.301 clause 8.2.21 Table 8.2.21.1), the
// constructor takes only variadic tagged IEs. First occurrence of each Kind wins;
// duplicates and unrecognized IEs are appended to the IEs field.
func NewSecurityModeComplete(ies ...*ie.IE) *SecurityModeComplete {
	m := &SecurityModeComplete{}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
		case ie.MobileIdentity:
			if m.IMEISV == nil {
				m.IMEISV = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.ReplayedNASMessageContainer:
			if m.ReplayedNASMessageContainer == nil {
				m.ReplayedNASMessageContainer = i
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
			if m.IEs == nil {
				m.IEs = append(m.IEs, i)
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.UECoarseLocationInformation:
			if m.UECoarseLocationInformation == nil {
				m.UECoarseLocationInformation = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the SECURITY MODE COMPLETE message, or 0
// for a nil receiver. The result includes header and all named tagged IEs in
// spec-table order.
func (m *SecurityModeComplete) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	if m.IMEISV != nil {
		n += m.IMEISV.MarshalLen()
	}
	if m.ReplayedNASMessageContainer != nil {
		n += m.ReplayedNASMessageContainer.MarshalLen()
	}
	if m.UERadioCapabilityID != nil {
		n += m.UERadioCapabilityID.MarshalLen()
	}
	if m.UECoarseLocationInformation != nil {
		n += m.UECoarseLocationInformation.MarshalLen()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the SECURITY MODE COMPLETE message into b. The nil receiver
// returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too
// short. Header and all tagged IEs in spec-table order follow TS 24.301 clause
// 8.2.21 table order.
func (m *SecurityModeComplete) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: SecurityModeComplete is nil", ie.ErrNilIE)
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

	// Named tagged IEs in spec-table order.
	if m.IMEISV != nil {
		w, err := writeIE(b, n, m.IMEISV)
		if err != nil {
			return fmt.Errorf("IMEISV: %w", err)
		}
		n += w
	}
	if m.ReplayedNASMessageContainer != nil {
		w, err := writeIE(b, n, m.ReplayedNASMessageContainer)
		if err != nil {
			return fmt.Errorf("replayed NAS message container: %w", err)
		}
		n += w
	}
	if m.UERadioCapabilityID != nil {
		w, err := writeIE(b, n, m.UERadioCapabilityID)
		if err != nil {
			return fmt.Errorf("UE radio capability ID: %w", err)
		}
		n += w
	}
	if m.UECoarseLocationInformation != nil {
		w, err := writeIE(b, n, m.UECoarseLocationInformation)
		if err != nil {
			return fmt.Errorf("UE coarse location information: %w", err)
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

// UnmarshalBinary decodes a SECURITY MODE COMPLETE message from b. The receiver
// is left unchanged on error. Header only, then all tagged optional IEs resolved
// via prototypes for recognized tags. No positional mandatory part per TS 24.301
// Table 8.2.21.1.
func (m *SecurityModeComplete) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: SecurityModeComplete is nil", ie.ErrNilIE)
	}

	var tmp SecurityModeComplete

	// Decode header (2 octets for plain SHT 0x00).
	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	// Tagged optional IEs — prototypes resolve all recognized tags.
	if len(pay) > 0 {
		tagged, err := ie.ParseMultiIEs(pay)
		if err != nil {
			return fmt.Errorf("optional: %w", err)
		}

		for _, i := range tagged {
			switch i.Kind {
			case ie.MobileIdentity:
				if tmp.IMEISV == nil {
					tmp.IMEISV = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case ie.ReplayedNASMessageContainer:
				if tmp.ReplayedNASMessageContainer == nil {
					tmp.ReplayedNASMessageContainer = i
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
				if tmp.IEs == nil {
					tmp.IEs = append(tmp.IEs, i)
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case ie.UECoarseLocationInformation:
				if tmp.UECoarseLocationInformation == nil {
					tmp.UECoarseLocationInformation = i
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

// MessageType returns the SECURITY MODE COMPLETE message type constant (0x5E).
func (m *SecurityModeComplete) MessageType() uint8 {
	return EMMSecurityModeComplete
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *SecurityModeComplete) MessageTypeName() string {
	return "Security mode complete"
}

// IsRequest returns true — SECURITY MODE COMPLETE is sent UE to network.
func (m *SecurityModeComplete) IsRequest() bool {
	return true
}

// SecurityModeComplete implements the Message interface.
var _ Message = (*SecurityModeComplete)(nil)

// Header returns the EMM security header embedded in SecurityModeComplete.
func (m *SecurityModeComplete) Header() *EmmHeader {
	return &m.EmmHeader
}
