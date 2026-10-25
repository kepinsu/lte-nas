package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// PDNDisconnectRequest represents a PDN DISCONNECT REQUEST message as defined in
// TS 24.301 clause 8.3.22 (Table 8.3.22.1, direction UE to network).
//
// There are no positional mandatory IEs beyond the ESM header. The only tagged optional
// IE is the ESM message container (pass-through LV-E). Tagged IEs without ambiguous
// tags are resolved via ie.ParseMultiIEs and classified by Kind.
type PDNDisconnectRequest struct {
	EsmHeader
	ESMMessageContainer *ie.IE // tagged optional LV-E (9.9.3.15), pass-through
	IEs                 []*ie.IE
}

// NewPDNDisconnectRequest creates a PDN DISCONNECT REQUEST message with optional
// tagged IEs classified by Kind. ESMMessageContainer (Kind 0x78) is resolved first;
// duplicates and unrecognized IEs are appended to the IEs field.
func NewPDNDisconnectRequest(ies ...*ie.IE) *PDNDisconnectRequest {
	m := &PDNDisconnectRequest{}

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
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the PDN DISCONNECT REQUEST message, or 0 for
// a nil receiver. The result includes the 3-octet ESM header and all tagged IEs.
func (m *PDNDisconnectRequest) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EsmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	if m.ESMMessageContainer != nil {
		n += m.ESMMessageContainer.MarshalLen()
	}
	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the PDN DISCONNECT REQUEST message into b. The nil receiver
// returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too short.
// Header, named tagged IEs in spec-table order and extras follow TS 24.301 clause
// 8.3.22 table order.
func (m *PDNDisconnectRequest) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: PDNDisconnectRequest is nil", ie.ErrNilIE)
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
	if m.ESMMessageContainer != nil {
		w, err := writeIE(b, n, m.ESMMessageContainer)
		if err != nil {
			return fmt.Errorf("ESMMessageContainer: %w", err)
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

// UnmarshalBinary decodes a PDN DISCONNECT REQUEST message from b. The receiver is
// left unchanged on error (uses local tmp struct). Header decoded, tagged optional IEs
// resolved via ie.ParseMultiIEs. The ESM message container (Kind 0x78) is handled by
// ie.Parse and classified into the named field.
func (m *PDNDisconnectRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: PDNDisconnectRequest is nil", ie.ErrNilIE)
	}

	var tmp PDNDisconnectRequest

	// Decode header (3 octets for plain ESM).
	pay, err := tmp.EsmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	pos := 0

	// Tagged optional IEs — ambiguous tag 0x78 (ESM message container) is not
	// in the prototypes map since it is unambiguous; ie.Parse resolves it.
	if len(pay[pos:]) > 0 {
		tagged, err := ie.ParseMultiIEs(pay[pos:])
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
			default:
				tmp.IEs = append(tmp.IEs, i)
			}
		}
	}

	*m = tmp
	return nil
}

// MessageType returns the PDN DISCONNECT REQUEST message type constant (0xD2).
func (m *PDNDisconnectRequest) MessageType() uint8 {
	return ESMPDNDisconnectRequest
}

// MessageTypeName returns the descriptive name from Table 9.8.2.
func (m *PDNDisconnectRequest) MessageTypeName() string {
	return "PDN disconnect request"
}

// IsRequest returns true — PDN DISCONNECT REQUEST is sent UE to network.
func (m *PDNDisconnectRequest) IsRequest() bool {
	return true
}

// PDNDisconnectRequest implements the Message interface.
var _ Message = (*PDNDisconnectRequest)(nil)

// ParsePDNDisconnectRequest parses a byte slice as a PDN DISCONNECT REQUEST
// message. Convenience wrapper around UnmarshalBinary on a fresh receiver.
func ParsePDNDisconnectRequest(b []byte) (*PDNDisconnectRequest, error) {
	m := &PDNDisconnectRequest{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}
