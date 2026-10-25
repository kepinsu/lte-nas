package ie

import (
	"fmt"
	"io"
)

// TimeDurationFields is the structured value of the time duration IE (TS 24.301
// clause 9.9.3.68). The three-octet value encodes a duration in seconds as a
// big-endian 24-bit unsigned integer.
type TimeDurationFields struct {
	// Seconds is the time duration in seconds (24-bit big-endian).
	Seconds uint32
}

// NewPositionalTimeDurationWithFields creates a positional time duration IE from fields.
func NewPositionalTimeDurationWithFields(fields *TimeDurationFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalTimeDuration(payload), nil
}

// TimeDuration returns the IE as structured fields.
func (i *IE) TimeDuration() (*TimeDurationFields, error) {
	if err := i.requireDefinition(TimeDurationDefinition); err != nil {
		return nil, err
	}
	return ParseTimeDurationFields(i.Payload)
}

// ParseTimeDurationFields decodes a time duration value without its IEI and length octet.
func ParseTimeDurationFields(payload []byte) (*TimeDurationFields, error) {
	fields := &TimeDurationFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *TimeDurationFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: time duration fields are nil", ErrInvalidValue)
	}
	if f.Seconds > 0xffffff {
		return fmt.Errorf("%w: time duration seconds %d exceeds 24 bits", ErrInvalidValue, f.Seconds)
	}
	return nil
}

// Marshal encodes fields into a 3-octet time duration value.
func (f *TimeDurationFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, 3)
	payload[0] = byte(f.Seconds >> 16)
	payload[1] = byte(f.Seconds >> 8)
	payload[2] = byte(f.Seconds)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *TimeDurationFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 3 {
		return io.ErrShortBuffer
	}
	b[0] = byte(f.Seconds >> 16)
	b[1] = byte(f.Seconds >> 8)
	b[2] = byte(f.Seconds)
	return nil
}

// MarshalLen returns the size of the time duration value payload.
func (f *TimeDurationFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 3
}

// UnmarshalBinary decodes a 3-octet time duration value.
func (f *TimeDurationFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: time duration fields are nil", ErrInvalidValue)
	}
	if len(payload) != 3 {
		return fmt.Errorf("%w: time duration has %d value octets, want 3", ErrInvalidLength, len(payload))
	}
	var decoded TimeDurationFields
	decoded.Seconds = uint32(payload[0])<<16 | uint32(payload[1])<<8 | uint32(payload[2])
	*f = decoded
	return nil
}
