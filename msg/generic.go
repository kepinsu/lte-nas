package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// Generic wraps an unparseable or unknown NAS message as raw bytes. It implements
// the Message interface so it can stand in for any message type during dispatch or
// inspection. MessageType() makes a best-effort extraction using PD/SHT rules; it
// returns 0 when the type cannot be determined.
type Generic struct {
	Raw []byte
}

// NewGeneric returns a new Generic that stores a copy of raw (never the slice
// header itself). This ensures callers can reuse their input buffer after passing it.
func NewGeneric(raw []byte) *Generic {
	if raw == nil {
		return &Generic{}
	}
	cp := make([]byte, len(raw))
	copy(cp, raw)
	return &Generic{Raw: cp}
}

// MarshalTo copies g.Raw into dst. The nil receiver returns a wrapped ie.ErrNilIE
// error; if dst is too short, io.ErrShortBuffer is returned without writing.
func (g *Generic) MarshalTo(dst []byte) error {
	if g == nil {
		return fmt.Errorf("%w: Generic is nil", ie.ErrNilIE)
	}
	if len(g.Raw) == 0 {
		return fmt.Errorf("%w: Generic Raw is nil or empty", ie.ErrNilIE)
	}
	if len(dst) < len(g.Raw) {
		return io.ErrShortBuffer
	}
	copy(dst, g.Raw)
	return nil
}

// MarshalLen returns len(Raw). The nil receiver returns 0.
func (g *Generic) MarshalLen() int {
	if g == nil {
		return 0
	}
	return len(g.Raw)
}

// UnmarshalBinary stores a copy of b in g.Raw.
func (g *Generic) UnmarshalBinary(b []byte) error {
	if b == nil {
		g.Raw = nil
		return nil
	}
	cp := make([]byte, len(b))
	copy(cp, b)
	g.Raw = cp
	return nil
}

// MessageType returns a best-effort message type from g.Raw. It applies PD/SHT
// header extraction rules: for 2+ bytes with low nibble 0x7E it treats byte[1]
// as an EMM type; for 3+ bytes with low nibble 0x2E it treats byte[2] as an ESM
// type; otherwise it returns 0.
func (g *Generic) MessageType() uint8 {
	if g == nil || len(g.Raw) < 2 {
		return 0
	}

	switch g.Raw[0] & 0x0f {
	case PdEmm:
		return g.Raw[1]
	case PdEsm:
		if len(g.Raw) >= 3 {
			return g.Raw[2]
		}
	}

	return 0
}

// MessageTypeName returns a string like "Generic (0x41)" based on the best-effort
// MessageType(). When no type can be determined, it falls back to "Generic (0x00)".
func (g *Generic) MessageTypeName() string {
	mt := g.MessageType()
	return fmt.Sprintf("Generic (%#02x)", mt)
}

// IsRequest returns false because Generic represents an unknown type whose
// request/response direction cannot be determined.
func (g *Generic) IsRequest() bool {
	return false
}
