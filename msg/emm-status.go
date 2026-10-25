package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// EmmStatus represents an EMM STATUS message as defined in TS 24.301
// clause 8.2.14 (Table 8.2.14.1, dual-direction).
//
// The message carries a single positional mandatory IE: EMM cause (type-3 V,
// one octet) per TS 24.301 clause 9.9.3.9. There is no tagged optional part.
type EmmStatus struct {
	EmmHeader
	EMMCause *ie.IE // positional mandatory type-3 (9.9.3.9), one-octet cause
	IEs      []*ie.IE
}

// NewEmmStatus creates an EMM STATUS message with the mandatory positional
// EMM cause IE from TS 24.301 clause 8.2.14. The variadic ies are ignored
// because EMM STATUS has no tagged optional IEs; any supplied IEs are stored
// in the IEs field (unused for this message type).
func NewEmmStatus(cause *ie.IE, _ ...*ie.IE) *EmmStatus {
	if cause == nil {
		panic("EmmStatus EMM cause must not be nil")
	}

	return &EmmStatus{EMMCause: cause}
}

// MarshalLen returns the wire size of the EMM STATUS message, or 0 for a nil
// receiver. The result includes header (2 plain / 7 integrity-protected) plus
// the positional EMM cause IE (one octet). No tagged IEs are emitted.
func (m *EmmStatus) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	n += m.EMMCause.MarshalLenWithoutTag()
	return n
}

// MarshalTo encodes the EMM STATUS message into b. The nil receiver returns a
// wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too short.
// Writes header and positional EMM cause IE only.
func (m *EmmStatus) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: EmmStatus is nil", ie.ErrNilIE)
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

	if m.EMMCause == nil {
		return io.ErrUnexpectedEOF
	}
	if err := m.EMMCause.MarshalToWithoutTag(b); err != nil {
		return fmt.Errorf("EMMCause: %w", err)
	}

	return nil
}

// UnmarshalBinary decodes an EMM STATUS message from b. The receiver is left
// unchanged on error (uses local tmp struct). Header via UnmarshalBinary,
// positional EMM cause via ParseWithoutTag. Rejects trailing bytes after EMM
// cause per TS 24.301 clause 8.2.14.
func (m *EmmStatus) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: EmmStatus is nil", ie.ErrNilIE)
	}

	var tmp EmmStatus

	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	tmp.EMMCause, err = ie.ParseWithoutTag(pay[:1], ie.NewPositionalEMMCause(0))
	if err != nil {
		return fmt.Errorf("EMMCause: %w", err)
	}

	if len(pay) > 1 {
		return fmt.Errorf("%w: trailing bytes after EMM cause", ie.ErrUnexpectedIE)
	}

	*m = tmp
	return nil
}

// MessageType returns the EMM STATUS message type constant (0x60).
func (m *EmmStatus) MessageType() uint8 {
	return EMMStatus
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *EmmStatus) MessageTypeName() string {
	return "EMM status"
}

// IsRequest returns false — EMM STATUS is dual-direction but classified as
// network-originated in the Message interface contract.
func (m *EmmStatus) IsRequest() bool {
	return false
}

// EmmStatus implements the Message interface.
var _ Message = (*EmmStatus)(nil)

// Header returns the EMM security header embedded in EmmStatus.
func (m *EmmStatus) Header() *EmmHeader {
	return &m.EmmHeader
}
