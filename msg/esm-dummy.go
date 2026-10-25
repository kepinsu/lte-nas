package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// DummyMessage represents an ESM DUMMY message as defined in TS 24.301
// clause 8.3.12A (Table 8.3.12A.1, direction network to UE).
//
// This message carries no body IEs beyond the 3-octet ESM header. Its
// wire size is always exactly 3 octets regardless of header content.
type DummyMessage struct {
	EsmHeader
}

// NewESMDummy creates a pre-configured ESM DUMMY message. It takes no
// arguments because the ESM DUMMY message carries no body IEs — only the
// 3-octet header is valid. Set BearerIdentity, PTI and MsgType on the
// returned struct after construction.
func NewESMDummy() *DummyMessage {
	return &DummyMessage{}
}

// MarshalLen always returns 3 for the fixed ESM header length. The nil
// receiver returns 0.
func (m *DummyMessage) MarshalLen() int {
	if m == nil {
		return 0
	}
	return 3
}

// MarshalTo writes exactly 3 bytes (the ESM header) into b. The nil receiver
// returns a wrapped ie.ErrNilIE; io.ErrShortBuffer if b is too short.
func (m *DummyMessage) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: DummyMessage is nil", ie.ErrNilIE)
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

	return nil
}

// UnmarshalBinary decodes an ESM DUMMY message from b. The receiver is left
// unchanged on error (uses local tmp struct). Header via UnmarshalBinary;
// any trailing payload bytes after the 3-octet header cause io.ErrUnexpectedEOF.
func (m *DummyMessage) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: DummyMessage is nil", ie.ErrNilIE)
	}

	var tmp DummyMessage

	pay, err := tmp.EsmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	// ESM DUMMY carries no body IEs; reject any trailing bytes.
	if len(pay) > 0 {
		return io.ErrUnexpectedEOF
	}

	*m = tmp
	return nil
}

// MessageType returns the ESM DUMMY message type constant (0xDC).
func (m *DummyMessage) MessageType() uint8 {
	return ESMDummy
}

// MessageTypeName returns the descriptive name from Table 9.8.2.
func (m *DummyMessage) MessageTypeName() string {
	return "ESM dummy message"
}

// IsRequest returns false — ESM DUMMY is network to UE.
func (m *DummyMessage) IsRequest() bool {
	return false
}

// DummyMessage implements the Message interface.
var _ Message = (*DummyMessage)(nil)

// ParseESMDummy parses a byte slice as an ESM DUMMY message. Convenience wrapper
// around UnmarshalBinary on a fresh receiver.
func ParseESMDummy(b []byte) (*DummyMessage, error) {
	m := &DummyMessage{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}
