package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// DownlinkGenericNASTransport represents a DOWNLINK GENERIC NAS TRANSPORT
// message as defined in TS 24.301 clause 8.2.31 (Table 8.2.31.1, direction
// network to UE).
//
// The tagged optional part carries a single pass-through IE: generic NAS
// message container (LV type-4, no catalogue Kind). No positional mandatory
// part beyond the header.
type DownlinkGenericNASTransport struct {
	EmmHeader
	GenericNASContainer *ie.IE // tagged optional LV type-4, pass-through (no catalogue Kind)
	IEs                 []*ie.IE
}

// NewDownlinkGenericNASTransport creates a DOWNLINK GENERIC NAS TRANSPORT
// message with the optional tagged generic NAS message container IE as the
// first variadic parameter, followed by additional tagged IEs. First occurrence
// wins; duplicates and unrecognized IEs are appended to the IEs field.
func NewDownlinkGenericNASTransport(container *ie.IE, ies ...*ie.IE) *DownlinkGenericNASTransport {
	if container == nil {
		panic("DownlinkGenericNASTransport generic NAS container must not be nil")
	}

	m := &DownlinkGenericNASTransport{GenericNASContainer: container}

	for _, i := range ies {
		if i == nil {
			continue
		}
		m.IEs = append(m.IEs, i)
	}

	return m
}

// MarshalLen returns the wire size of the DOWNLINK GENERIC NAS TRANSPORT
// message, or 0 for a nil receiver. The result includes header and the tagged
// generic NAS message container IE.
func (m *DownlinkGenericNASTransport) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	n += m.GenericNASContainer.MarshalLen()

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the DOWNLINK GENERIC NAS TRANSPORT message into b. The nil
// receiver returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if
// b is too short. Writes header and tagged generic NAS message container IE.
func (m *DownlinkGenericNASTransport) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: DownlinkGenericNASTransport is nil", ie.ErrNilIE)
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

	if m.GenericNASContainer == nil {
		return io.ErrUnexpectedEOF
	}
	w, err := writePassThroughTagged(b, 0, m.GenericNASContainer)
	if err != nil {
		return fmt.Errorf("GenericNASContainer: %w", err)
	}
	b = b[w:]

	for _, i := range m.IEs {
		w, err := writeIE(b, 0, i)
		if err != nil {
			return fmt.Errorf("extra IE: %w", err)
		}
		b = b[w:]
	}

	return nil
}

// UnmarshalBinary decodes a DOWNLINK GENERIC NAS TRANSPORT message from b. The
// receiver is left unchanged on error (uses local tmp struct). Header via
// UnmarshalBinary, tagged generic NAS message container parsed as pass-through.
func (m *DownlinkGenericNASTransport) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: DownlinkGenericNASTransport is nil", ie.ErrNilIE)
	}

	var tmp DownlinkGenericNASTransport

	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	container, extras, err := parseTaggedContainerIEs(pay)
	if err != nil {
		return fmt.Errorf("optional: %w", err)
	}

	tmp.GenericNASContainer = container
	tmp.IEs = extras

	*m = tmp
	return nil
}

// MessageType returns the DOWNLINK GENERIC NAS TRANSPORT message type constant
// (0x68).
func (m *DownlinkGenericNASTransport) MessageType() uint8 {
	return EMMDownlinkGenericNASTransport
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *DownlinkGenericNASTransport) MessageTypeName() string {
	return "Downlink generic NAS transport"
}

// IsRequest returns false — DOWNLINK GENERIC NAS TRANSPORT is sent network to
// UE.
func (m *DownlinkGenericNASTransport) IsRequest() bool {
	return false
}

// DownlinkGenericNASTransport implements the Message interface.
var _ Message = (*DownlinkGenericNASTransport)(nil)

// parseTaggedContainerIEs parses tagged IEs from the payload, returning the
// first IE as a container (for generic NAS transport) and the rest as extras.
func parseTaggedContainerIEs(b []byte) (*ie.IE, []*ie.IE, error) {
	if len(b) == 0 {
		return nil, nil, io.ErrUnexpectedEOF
	}

	container, remaining, err := parseFirstTaggedIE(b)
	if err != nil {
		return nil, nil, err
	}

	var extras []*ie.IE
	for len(remaining) > 0 {
		exIE, rest, err := parseFirstTaggedIE(remaining)
		if err != nil {
			return nil, nil, err
		}
		extras = append(extras, exIE)
		remaining = rest
	}

	return container, extras, nil
}

// parseFirstTaggedIE decodes the first tagged IE from b, returning the IE and
// remaining bytes. Known catalogue tags (ESM Message Container tag 0x78) use
// ie.Parse. For non-catalogue IEs, it first attempts type-4 TLV (tag + length
// + value); if the length byte indicates more data than is available, it falls
// back to type-1 TV (single wire byte where payload & 0x0f | Kind = tag).
func parseFirstTaggedIE(b []byte) (*ie.IE, []byte, error) {
	if len(b) == 0 {
		return nil, nil, io.ErrUnexpectedEOF
	}

	tag := b[0]

	switch tag {
	case 0x78: // ESM Message Container — catalogue type-6 TLV-E
		if i, err := ie.Parse(b, ie.Kind(tag)); err == nil {
			return i, b[i.MarshalLen():], nil
		}
		fallthrough
	default:
		// Non-catalogue IE — try type-4 TLV first (requires length byte and payload)
		if len(b) >= 2 {
			length := int(b[1])
			if length > 0 && len(b) >= 2+length {
				// Looks like type-4 TLV: tag + length + value
				payload := make([]byte, length)
				copy(payload, b[2:2+length])
				return &ie.IE{Kind: ie.Kind(tag), Payload: payload, Format: ie.IEFormatFour}, b[2+length:], nil
			}
			// Length indicates no payload or truncated — treat as type-1 TV
		} else if len(b) >= 1 {
			// Type-4 TLV impossible with less than 2 bytes, try type-1 TV
		} else {
			return nil, nil, io.ErrUnexpectedEOF
		}
		// Type-1 TV: single wire byte = payload & 0x0f | tag_nibble (Kind)
		wire := b[0]
		return &ie.IE{Kind: ie.Kind(wire & 0xf0), Payload: []byte{wire & 0x0f}, Format: ie.IEFormatOne}, b[1:], nil
	}
}

// Header returns the EMM security header embedded in DownlinkGenericNASTransport.
func (m *DownlinkGenericNASTransport) Header() *EmmHeader {
	return &m.EmmHeader
}
