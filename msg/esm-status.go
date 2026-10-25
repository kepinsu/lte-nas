package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// ESMStatusMsg represents an ESM STATUS message as defined in TS 24.301
// clause 8.3.15 (Table 8.3.15.1, direction network to UE).
//
// The message carries a single positional mandatory IE: ESM cause (type-3 V,
// one octet) per TS 24.301 clause 9.9.4.4. There is no tagged optional part.
type ESMStatusMsg struct {
	EsmHeader
	ESMCause *ie.IE // positional mandatory type-3 (9.9.4.4), one-octet cause
	IEs      []*ie.IE
}

// NewESMStatusMsg creates an ESM STATUS message with the mandatory positional
// ESM cause IE from TS 24.301 clause 8.3.15.
func NewESMStatusMsg(cause *ie.IE, _ ...*ie.IE) *ESMStatusMsg {
	if cause == nil {
		panic("ESMStatusMsg ESM cause must not be nil")
	}

	return &ESMStatusMsg{ESMCause: cause}
}

// MarshalLen returns the wire size of the ESM STATUS message, or 0 for a nil
// receiver. The result includes header (3 octets) plus the positional ESM cause
// IE (one octet). No tagged IEs are emitted.
func (m *ESMStatusMsg) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EsmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	n += m.ESMCause.MarshalLenWithoutTag()
	return n
}

// MarshalTo encodes the ESM STATUS message into b. The nil receiver returns a
// wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too short.
func (m *ESMStatusMsg) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: ESMStatusMsg is nil", ie.ErrNilIE)
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

	if m.ESMCause == nil {
		return io.ErrUnexpectedEOF
	}
	if err := m.ESMCause.MarshalToWithoutTag(b); err != nil {
		return fmt.Errorf("ESMCause: %w", err)
	}

	return nil
}

// UnmarshalBinary decodes an ESM STATUS message from b. The receiver is left
// unchanged on error (uses local tmp struct). Header via UnmarshalBinary,
// positional ESM cause via ParseWithoutTag. Rejects trailing bytes after cause.
func (m *ESMStatusMsg) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: ESMStatusMsg is nil", ie.ErrNilIE)
	}

	var tmp ESMStatusMsg

	pay, err := tmp.EsmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	tmp.ESMCause, err = ie.ParseWithoutTag(pay[:1], ie.NewPositionalESMCause(0))
	if err != nil {
		return fmt.Errorf("ESMCause: %w", err)
	}

	if len(pay) > 1 {
		return fmt.Errorf("%w: trailing bytes after ESM cause", ie.ErrUnexpectedIE)
	}

	*m = tmp
	return nil
}

// MessageType returns the ESM STATUS message type constant (0xE8).
func (m *ESMStatusMsg) MessageType() uint8 {
	return ESMStatus
}

// MessageTypeName returns the descriptive name from Table 9.8.2.
func (m *ESMStatusMsg) MessageTypeName() string {
	return "ESM status"
}

// IsRequest returns false — ESM STATUS is network to UE.
func (m *ESMStatusMsg) IsRequest() bool {
	return false
}

// ESMStatusMsg implements the Message interface.
var _ Message = (*ESMStatusMsg)(nil)

// ParseESMStatusMsg parses a byte slice as an ESM STATUS message. Convenience
// wrapper around UnmarshalBinary on a fresh receiver.
func ParseESMStatusMsg(b []byte) (*ESMStatusMsg, error) {
	m := &ESMStatusMsg{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}
