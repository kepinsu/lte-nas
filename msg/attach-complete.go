package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// AttachComplete represents an ATTACH COMPLETE message as defined in TS 24.301
// clause 8.2.2 (Table 8.2.2.1, direction UE to network).
//
// There are no positional mandatory IEs after the header. The ESM message
// container (tag 0x78) is mandatory and all other tagged IEs are optional,
// resolved via message context prototypes.
type AttachComplete struct {
	EmmHeader
	ESMMessageContainer  *ie.IE // mandatory TLV-E, 0x78 (9.9.3.15), pass-through
	S1UETNBGlobalID      *ie.IE // optional TLV, 0x39 (9.9.3.24)
	UES6ACapability      *ie.IE // optional TV type-1, 0x3c (9.9.3.42)
	EMMCause             *ie.IE // optional TV type-1, 0x53 (9.9.3.9)
	PTI                  *ie.IE // optional TV type-1, 0x59 (9.4)
	T3402Value           *ie.IE // optional TLV, 0x16 (9.9.3.16A)
	TrackingAreaIdentity *ie.IE // optional TLV, 0x1d (9.9.3.32)
	IEs                  []*ie.IE
}

// NewAttachComplete creates an ATTACH COMPLETE message with optional tagged
// IEs classified by their Kind. Fields for recognized tags are resolved first;
// duplicates and unrecognized IEs are appended to the IEs field.
func NewAttachComplete(ies ...*ie.IE) *AttachComplete {
	m := &AttachComplete{}

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
		case ie.NetworkResourceIdentifierContainer:
			if m.S1UETNBGlobalID == nil {
				m.S1UETNBGlobalID = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.UEAdditionalSecurityCapability:
			if m.UES6ACapability == nil {
				m.UES6ACapability = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.EMMCause:
			if m.EMMCause == nil {
				m.EMMCause = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.T3402ValueForAttachReject:
			if m.T3402Value == nil {
				m.T3402Value = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.T3402Value:
			if m.T3402Value == nil {
				m.T3402Value = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the ATTACH COMPLETE message, or 0 for a
// nil receiver. The result includes header and all tagged IEs in spec-table order.
func (m *AttachComplete) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	if m.ESMMessageContainer != nil {
		n += m.ESMMessageContainer.MarshalLen()
	}

	// Named tagged IEs in spec-table order.
	if m.S1UETNBGlobalID != nil {
		n += m.S1UETNBGlobalID.MarshalLen()
	}
	if m.UES6ACapability != nil {
		n += m.UES6ACapability.MarshalLen()
	}
	if m.EMMCause != nil {
		n += m.EMMCause.MarshalLen()
	}
	if m.PTI != nil {
		n += m.PTI.MarshalLen()
	}
	if m.T3402Value != nil {
		n += m.T3402Value.MarshalLen()
	}
	if m.TrackingAreaIdentity != nil {
		n += m.TrackingAreaIdentity.MarshalLen()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the ATTACH COMPLETE message into b. The nil receiver
// returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too
// short. Header, mandatory tagged IE and all optional tagged IEs follow spec-table order.
func (m *AttachComplete) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: AttachComplete is nil", ie.ErrNilIE)
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

	// Mandatory tagged IE first.
	if m.ESMMessageContainer != nil {
		w, err := writeIE(b, n, m.ESMMessageContainer)
		if err != nil {
			return fmt.Errorf("ESMMessageContainer: %w", err)
		}
		n += w
	}

	// Optional tagged IEs in spec-table order.
	if m.S1UETNBGlobalID != nil {
		w, err := writeIE(b, n, m.S1UETNBGlobalID)
		if err != nil {
			return fmt.Errorf("S1UETNBGlobalID: %w", err)
		}
		n += w
	}
	if m.UES6ACapability != nil {
		w, err := writeIE(b, n, m.UES6ACapability)
		if err != nil {
			return fmt.Errorf("UES6ACapability: %w", err)
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
	if m.PTI != nil {
		w, err := writeIE(b, n, m.PTI)
		if err != nil {
			return fmt.Errorf("PTI: %w", err)
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
	if m.TrackingAreaIdentity != nil {
		w, err := writeIE(b, n, m.TrackingAreaIdentity)
		if err != nil {
			return fmt.Errorf("TrackingAreaIdentity: %w", err)
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

// UnmarshalBinary decodes an ATTACH COMPLETE message from b. The receiver is
// left unchanged on error. Header first, then tagged optional IEs resolved
// via prototypes including the mandatory ESM message container.
func (m *AttachComplete) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: AttachComplete is nil", ie.ErrNilIE)
	}

	var tmp AttachComplete

	// Decode header (2 octets for plain).
	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	// Tagged optional IEs — prototype for mandatory ESM message container.
	tagged, err := ie.ParseMultiIEs(pay)
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
		case ie.NetworkResourceIdentifierContainer:
			if tmp.S1UETNBGlobalID == nil {
				tmp.S1UETNBGlobalID = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.UEAdditionalSecurityCapability:
			if tmp.UES6ACapability == nil {
				tmp.UES6ACapability = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.EMMCause:
			if tmp.EMMCause == nil {
				tmp.EMMCause = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.T3402ValueForAttachReject:
			if tmp.T3402Value == nil {
				tmp.T3402Value = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.T3402Value:
			if tmp.T3402Value == nil {
				tmp.T3402Value = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		default:
			tmp.IEs = append(tmp.IEs, i)
		}
	}

	if tmp.ESMMessageContainer == nil {
		return fmt.Errorf("%w: ESM message container is mandatory", ie.ErrInvalidLength)
	}

	*m = tmp
	return nil
}

// MessageType returns the ATTACH COMPLETE message type constant (0x43).
func (m *AttachComplete) MessageType() uint8 {
	return EMMAttachComplete
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *AttachComplete) MessageTypeName() string {
	return "Attach complete"
}

// IsRequest returns true — ATTACH COMPLETE is sent UE to network.
func (m *AttachComplete) IsRequest() bool {
	return true
}

// AttachComplete implements the Message interface.
var _ Message = (*AttachComplete)(nil)

// Header returns the EMM security header embedded in AttachComplete.
func (m *AttachComplete) Header() *EmmHeader {
	return &m.EmmHeader
}
