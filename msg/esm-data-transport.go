package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// EsmDataTransport represents an ESM DATA TRANSPORT message as defined in TS
// 24.301 clause 8.3.25 (Table 8.3.25.1, dual-direction; treated as network to
// UE for this implementation).
//
// The message carries an ESM header followed by tagged optional IEs: user data
// container (LV-E, tag 0x78) and control plane only indication (TV-type-1, tag
// 0x90). Both are resolved via prototypes in ie.ParseMultiIEs; duplicates and
// unrecognized IEs are appended to the IEs field.
type EsmDataTransport struct {
	EsmHeader
	UserDataContainer     *ie.IE // tagged optional LV-E (tag 0x78, clause 9.9.4.24), pass-through
	ControlPlaneOnlyIndic *ie.IE // tagged optional TV-type-1 (tag 0x90, clause 9.9.4.23)
	IEs                   []*ie.IE
}

// NewEsmDataTransport creates an ESM DATA TRANSPORT message with optional tagged
// IEs. First occurrence of UserDataContainer (tag 0x78) and ControlPlaneOnlyIndic
// (tag 0x90) wins; duplicates and unrecognized IEs are appended to the IEs field.
func NewEsmDataTransport(ies ...*ie.IE) *EsmDataTransport {
	m := &EsmDataTransport{}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
		case ie.ESMMessageContainer: // user data container (same tag 0x78)
			if m.UserDataContainer == nil {
				m.UserDataContainer = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.ControlPlaneOnlyIndication:
			if m.ControlPlaneOnlyIndic == nil {
				m.ControlPlaneOnlyIndic = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the ESM DATA TRANSPORT message, or 0 for a
// nil receiver. The result includes the 3-octet ESM header plus all named tagged
// IEs in spec-table order.
func (m *EsmDataTransport) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EsmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	if m.UserDataContainer != nil {
		n += m.UserDataContainer.MarshalLen()
	}

	if m.ControlPlaneOnlyIndic != nil {
		n += m.ControlPlaneOnlyIndic.MarshalLen()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the ESM DATA TRANSPORT message into b. The nil receiver
// returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too
// short. Writes header, then named tagged IEs in spec-table order.
func (m *EsmDataTransport) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: EsmDataTransport is nil", ie.ErrNilIE)
	}

	hlen := m.EsmHeader.MarshalLen()
	if hlen == 0 {
		return fmt.Errorf("%w: EsmHeader MarshalLen returned 0", ie.ErrNilIE)
	}

	if len(b) < hlen {
		return io.ErrShortBuffer
	}

	if err := m.EsmHeader.MarshalTo(b); err != nil {
		return fmt.Errorf("header: %w", err)
	}

	b = b[hlen:]
	n := 0

	// Named tagged IEs in spec-table order.
	if m.UserDataContainer != nil {
		w, err := writeIE(b, n, m.UserDataContainer)
		if err != nil {
			return fmt.Errorf("user data container: %w", err)
		}
		n += w
	}

	if m.ControlPlaneOnlyIndic != nil {
		w, err := writeIE(b, n, m.ControlPlaneOnlyIndic)
		if err != nil {
			return fmt.Errorf("control plane only indication: %w", err)
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

// UnmarshalBinary decodes an ESM DATA TRANSPORT message from b. The receiver is
// left unchanged on error (uses local tmp struct). Header via EsmHeader
// UnmarshalBinary; tagged optional IEs resolved via ie.ParseMultiIEs.
func (m *EsmDataTransport) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: EsmDataTransport is nil", ie.ErrNilIE)
	}

	var tmp EsmDataTransport

	// Decode header (3 octets for plain ESM).
	pay, err := tmp.EsmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	// Tagged optional IEs — prototypes resolve tag 0x78 to user data container
	// (LV-E, clause 9.9.4.24) and tag 0x90 to control plane only indication
	// (TV-type-1, clause 9.9.4.23).
	tagged, err := ie.ParseMultiIEs(pay)
	if err != nil {
		return fmt.Errorf("optional: %w", err)
	}

	for _, i := range tagged {
		switch i.Kind {
		case ie.ESMMessageContainer: // user data container (same tag 0x78)
			if tmp.UserDataContainer == nil {
				tmp.UserDataContainer = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.ControlPlaneOnlyIndication:
			if tmp.ControlPlaneOnlyIndic == nil {
				tmp.ControlPlaneOnlyIndic = i
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

// MessageType returns the ESM DATA TRANSPORT message type constant (0xEB).
func (m *EsmDataTransport) MessageType() uint8 {
	return ESMDataTransport
}

// MessageTypeName returns the descriptive name from Table 9.8.2.
func (m *EsmDataTransport) MessageTypeName() string {
	return "ESM data transport"
}

// IsRequest returns false — ESM DATA TRANSPORT is treated as network to UE.
func (m *EsmDataTransport) IsRequest() bool {
	return false
}

// EsmDataTransport implements the Message interface.
var _ Message = (*EsmDataTransport)(nil)

// ParseEsmDataTransport parses a byte slice as an ESM DATA TRANSPORT message.
// Convenience wrapper around UnmarshalBinary on a fresh receiver.
func ParseEsmDataTransport(b []byte) (*EsmDataTransport, error) {
	m := &EsmDataTransport{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}
