package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// Notification represents a NOTIFICATION message as defined in TS 24.301
// clause 8.3.18A (Table 8.3.18A.1, direction network to UE).
//
// The positional mandatory part carries a single Notification Indicator (LV, 2 octets:
// length=0x01 followed by one indicator byte) per TS 24.301 clause 9.9.4.7A.
type Notification struct {
	EsmHeader
	NotificationIndicator *ie.IE // positional LV (9.9.4.7A), length=1, 2 octets on wire
	IEs                   []*ie.IE
}

// NewNotification creates a NOTIFICATION message with the mandatory positional
// Notification Indicator IE from TS 24.301 clause 8.3.18A.
func NewNotification(notificationIndicator *ie.IE, ies ...*ie.IE) *Notification {
	if notificationIndicator == nil {
		panic("Notification NotificationIndicator must not be nil")
	}

	m := &Notification{NotificationIndicator: notificationIndicator}

	for _, i := range ies {
		if i == nil {
			continue
		}
		m.IEs = append(m.IEs, i)
	}

	return m
}

// MarshalLen returns the wire size of the NOTIFICATION message, or 0 for a nil
// receiver. The result includes header (3 octets) plus the positional LV IE
// (2 octets: length=0x01 + value). No tagged IEs are emitted.
func (m *Notification) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EsmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	n += m.NotificationIndicator.MarshalLenWithoutTag()
	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the NOTIFICATION message into b. The nil receiver returns a
// wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too short.
func (m *Notification) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: Notification is nil", ie.ErrNilIE)
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

	if m.NotificationIndicator == nil {
		return io.ErrUnexpectedEOF
	}
	if err := m.NotificationIndicator.MarshalToWithoutTag(b); err != nil {
		return fmt.Errorf("NotificationIndicator: %w", err)
	}

	return nil
}

// UnmarshalBinary decodes a NOTIFICATION message from b. The receiver is left
// unchanged on error (uses local tmp struct). Header, positional Notification
// Indicator via ParseWithoutTag. Rejects trailing bytes after the indicator.
func (m *Notification) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: Notification is nil", ie.ErrNilIE)
	}

	var tmp Notification

	// Decode header (3 octets for plain ESM).
	pay, err := tmp.EsmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	if len(pay) < 2 {
		return io.ErrUnexpectedEOF
	}

	// Positional Notification Indicator (LV, no tag).
	tmp.NotificationIndicator, err = ie.ParseWithoutTag(pay[:2], ie.NewPositionalNotificationIndicator(0))
	if err != nil {
		return fmt.Errorf("NotificationIndicator: %w", err)
	}

	if len(pay) > 2 {
		return fmt.Errorf("%w: trailing bytes after Notification Indicator", ie.ErrUnexpectedIE)
	}

	*m = tmp
	return nil
}

// MessageType returns the NOTIFICATION message type constant (0xDB).
func (m *Notification) MessageType() uint8 {
	return ESMNotification
}

// MessageTypeName returns the descriptive name from Table 9.8.2.
func (m *Notification) MessageTypeName() string {
	return "Notification"
}

// IsRequest returns false — NOTIFICATION is network to UE.
func (m *Notification) IsRequest() bool {
	return false
}

// Notification implements the Message interface.
var _ Message = (*Notification)(nil)

// ParseNotification parses a byte slice as a NOTIFICATION message. Convenience
// wrapper around UnmarshalBinary on a fresh receiver.
func ParseNotification(b []byte) (*Notification, error) {
	m := &Notification{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}
