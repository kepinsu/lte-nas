package ie

import (
	"fmt"
	"io"
)

// NotificationIndicatorFields is the structured value of the Notification
// indicator IE (TS 24.301 clause 9.9.4.7A).
type NotificationIndicatorFields struct {
	Value uint8
}

// NewNotificationIndicatorWithFields creates a tagged Notification indicator IE
// from fields.
func NewNotificationIndicatorWithFields(fields *NotificationIndicatorFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalNotificationIndicator(payload[0]), nil
}

// NewPositionalNotificationIndicatorWithFields creates a positional Notification
// indicator IE from fields.
func NewPositionalNotificationIndicatorWithFields(fields *NotificationIndicatorFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalNotificationIndicator(payload[0]), nil
}

// NotificationIndicator returns the IE as structured fields.
func (i *IE) NotificationIndicator() (*NotificationIndicatorFields, error) {
	if err := i.requireDefinition(NotificationIndicatorDefinition); err != nil {
		return nil, err
	}
	return ParseNotificationIndicatorFields(i.Payload)
}

// ParseNotificationIndicatorFields decodes a Notification indicator value without
// its IEI.
func ParseNotificationIndicatorFields(payload []byte) (*NotificationIndicatorFields, error) {
	fields := &NotificationIndicatorFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Value is a defined pattern.
func (f *NotificationIndicatorFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: notification indicator fields are nil", ErrInvalidValue)
	}
	v := f.Value
	if v != 0x01 && v != 0x60 && v != 0x7f {
		return fmt.Errorf("%w: notification indicator value %d is reserved", ErrInvalidValue, v)
	}
	return nil
}

// Marshal encodes fields into a Notification indicator value, without its IEI.
func (f *NotificationIndicatorFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{f.Value}, nil
}

// MarshalTo encodes fields into b.
func (f *NotificationIndicatorFields) MarshalTo(b []byte) error {
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	payload, err := f.Marshal()
	if err != nil {
		return err
	}
	b[0] = payload[0]
	return nil
}

// MarshalLen returns the size of the Notification indicator value payload.
func (f *NotificationIndicatorFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a Notification indicator value without its IEI.
func (f *NotificationIndicatorFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: notification indicator fields are nil", ErrInvalidValue)
	}
	if len(payload) == 0 {
		return io.ErrShortBuffer
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: notification indicator has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	v := payload[0]
	if v != 0x01 && v != 0x60 && v != 0x7f {
		return fmt.Errorf("%w: notification indicator value %d is reserved", ErrInvalidValue, v)
	}
	*f = NotificationIndicatorFields{
		Value: v,
	}
	return nil
}
