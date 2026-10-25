package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// GUTIReallocationComplete represents a GUTI REALLOCATION COMPLETE message as
// defined in TS 24.301 clause 8.2.17 (Table 8.2.17.1, direction UE to network).
//
// This message carries no mandatory content IEs — it is essentially empty
// after the header. Plain EMM header (SHT=0x00, 2 octets).
type GUTIReallocationComplete struct {
	EmmHeader
	IEs []*ie.IE
}

// NewGUTIReallocationComplete creates a GUTI REALLOCATION COMPLETE message.
// This message has no positional mandatory IEs; optional tagged IEs may be
// provided as variadic arguments.
func NewGUTIReallocationComplete(ies ...*ie.IE) *GUTIReallocationComplete {
	m := &GUTIReallocationComplete{}

	for _, i := range ies {
		if i == nil {
			continue
		}
		m.IEs = append(m.IEs, i)
	}

	return m
}

// MarshalLen returns the wire size of the GUTI REALLOCATION COMPLETE message, or
// 0 for a nil receiver. The result includes only the header and any extra IEs.
func (m *GUTIReallocationComplete) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the GUTI REALLOCATION COMPLETE message into b. The nil
// receiver returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b
// is too short. Header only (no content IEs) per TS 24.301 clause 8.2.17.
func (m *GUTIReallocationComplete) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: GUTIReallocationComplete is nil", ie.ErrNilIE)
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

	for _, i := range m.IEs {
		w, err := writeIE(b, n, i)
		if err != nil {
			return fmt.Errorf("extra IE: %w", err)
		}
		n += w
	}

	return nil
}

// UnmarshalBinary decodes a GUTI REALLOCATION COMPLETE message from b. The
// receiver is left unchanged on error. Header only — the spec table has no
// content IEs beyond the header. Extra bytes in the payload are classified as
// optional tagged IEs.
func (m *GUTIReallocationComplete) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: GUTIReallocationComplete is nil", ie.ErrNilIE)
	}

	var tmp GUTIReallocationComplete

	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	if len(pay) == 0 {
		// Empty optional part is valid for GUTI REALLOCATION COMPLETE.
		*m = tmp
		return nil
	}

	tagged, err := ie.ParseMultiIEs(pay)
	if err != nil {
		return fmt.Errorf("optional: %w", err)
	}

	tmp.IEs = append(tmp.IEs, tagged...)

	*m = tmp
	return nil
}

// MessageType returns the GUTI REALLOCATION COMPLETE message type constant (0x51).
func (m *GUTIReallocationComplete) MessageType() uint8 {
	return EMMGUTIReallocationComplete
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *GUTIReallocationComplete) MessageTypeName() string {
	return "Guti reallocation complete"
}

// IsRequest returns true — GUTI REALLOCATION COMPLETE is sent UE to network.
func (m *GUTIReallocationComplete) IsRequest() bool {
	return true
}

// GUTIReallocationComplete implements the Message interface.
var _ Message = (*GUTIReallocationComplete)(nil)

// Header returns the EMM security header embedded in GUTIReallocationComplete.
func (m *GUTIReallocationComplete) Header() *EmmHeader {
	return &m.EmmHeader
}
