package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// UplinkGenericNASTransport represents an UPLINK GENERIC NAS TRANSPORT message
// as defined in TS 24.301 clause 8.2.32 (Table 8.2.32.1, direction UE to
// network).
//
// The tagged optional part carries a single pass-through IE: generic NAS
// message container (LV type-4, no catalogue Kind). No positional mandatory
// part beyond the header.
type UplinkGenericNASTransport struct {
	EmmHeader
	GenericNASContainer *ie.IE // tagged optional LV type-4, pass-through (no catalogue Kind)
	IEs                 []*ie.IE
}

// NewUplinkGenericNASTransport creates an UPLINK GENERIC NAS TRANSPORT message
// with the optional tagged generic NAS message container IE as the first
// variadic parameter, followed by additional tagged IEs. First occurrence wins;
// duplicates and unrecognized IEs are appended to the IEs field.
func NewUplinkGenericNASTransport(container *ie.IE, ies ...*ie.IE) *UplinkGenericNASTransport {
	if container == nil {
		panic("UplinkGenericNASTransport generic NAS container must not be nil")
	}

	m := &UplinkGenericNASTransport{GenericNASContainer: container}

	for _, i := range ies {
		if i == nil {
			continue
		}
		m.IEs = append(m.IEs, i)
	}

	return m
}

// MarshalLen returns the wire size of the UPLINK GENERIC NAS TRANSPORT message,
// or 0 for a nil receiver. The result includes header and the tagged generic
// NAS message container IE.
func (m *UplinkGenericNASTransport) MarshalLen() int {
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

// MarshalTo encodes the UPLINK GENERIC NAS TRANSPORT message into b. The nil
// receiver returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if
// b is too short. Writes header and tagged generic NAS message container IE.
func (m *UplinkGenericNASTransport) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: UplinkGenericNASTransport is nil", ie.ErrNilIE)
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

// UnmarshalBinary decodes an UPLINK GENERIC NAS TRANSPORT message from b. The
// receiver is left unchanged on error (uses local tmp struct). Header via
// UnmarshalBinary, tagged generic NAS message container parsed as pass-through.
func (m *UplinkGenericNASTransport) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: UplinkGenericNASTransport is nil", ie.ErrNilIE)
	}

	var tmp UplinkGenericNASTransport

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

// MessageType returns the UPLINK GENERIC NAS TRANSPORT message type constant
// (0x69).
func (m *UplinkGenericNASTransport) MessageType() uint8 {
	return EMMUplinkGenericNASTransport
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *UplinkGenericNASTransport) MessageTypeName() string {
	return "Uplink generic NAS transport"
}

// IsRequest returns true — UPLINK GENERIC NAS TRANSPORT is sent UE to network.
func (m *UplinkGenericNASTransport) IsRequest() bool {
	return true
}

// UplinkGenericNASTransport implements the Message interface.
var _ Message = (*UplinkGenericNASTransport)(nil)

// Header returns the EMM security header embedded in UplinkGenericNASTransport.
func (m *UplinkGenericNASTransport) Header() *EmmHeader {
	return &m.EmmHeader
}
