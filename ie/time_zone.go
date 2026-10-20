package ie

import (
	"fmt"
	"io"
	"time"
)

const maximumTimeZoneQuarterHours = 79

// TimeZoneFields is the structured value of the Time zone IE (TS 24.301
// clause 9.9.3.29). Offset must be a signed multiple of 15 minutes.
type TimeZoneFields struct {
	Offset time.Duration
}

// NewTimeZoneWithFields creates a tagged Time zone IE from fields.
func NewTimeZoneWithFields(fields *TimeZoneFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewTimeZone(payload[0]), nil
}

// NewPositionalTimeZoneWithFields creates a positional Time zone IE from
// fields.
func NewPositionalTimeZoneWithFields(fields *TimeZoneFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalTimeZone(payload[0]), nil
}

// TimeZone returns the IE as structured fields.
func (i *IE) TimeZone() (*TimeZoneFields, error) {
	if err := i.requireDefinition(TimeZoneDefinition); err != nil {
		return nil, err
	}
	return ParseTimeZoneFields(i.Payload)
}

// ParseTimeZoneFields decodes a Time zone value without its IEI.
func ParseTimeZoneFields(payload []byte) (*TimeZoneFields, error) {
	fields := &TimeZoneFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Offset is representable by the Time zone IE.
func (f *TimeZoneFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: time zone fields are nil", ErrInvalidValue)
	}
	_, err := encodeTimeZone(f.Offset)
	return err
}

// Marshal encodes fields into a Time zone value without its IEI.
func (f *TimeZoneFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, fmt.Errorf("%w: time zone fields are nil", ErrInvalidValue)
	}
	encoded, err := encodeTimeZone(f.Offset)
	if err != nil {
		return nil, err
	}
	return []byte{encoded}, nil
}

// MarshalTo encodes fields into b.
func (f *TimeZoneFields) MarshalTo(b []byte) error {
	payload, err := f.Marshal()
	if err != nil {
		return err
	}
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	b[0] = payload[0]
	return nil
}

// MarshalLen returns the size of the Time zone value payload.
func (f *TimeZoneFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a Time zone value without its IEI.
func (f *TimeZoneFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: time zone fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: time zone has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	offset, err := decodeTimeZone(payload[0])
	if err != nil {
		return err
	}
	*f = TimeZoneFields{Offset: offset}
	return nil
}

func encodeTimeZone(offset time.Duration) (byte, error) {
	if offset%(15*time.Minute) != 0 {
		return 0, fmt.Errorf("%w: time zone offset %s is not a multiple of 15 minutes", ErrInvalidValue, offset)
	}
	negative := offset < 0
	if negative {
		offset = -offset
	}
	quarterHours := int(offset / (15 * time.Minute))
	if quarterHours > maximumTimeZoneQuarterHours {
		return 0, fmt.Errorf("%w: time zone offset %s exceeds 79 quarter-hours", ErrInvalidValue, offset)
	}
	encoded := byte(quarterHours%10)<<4 | byte(quarterHours/10)
	if negative {
		encoded |= 0x08
	}
	return encoded, nil
}

func decodeTimeZone(encoded byte) (time.Duration, error) {
	units := encoded >> 4
	tens := encoded & 0x07
	if units > 9 {
		return 0, fmt.Errorf("%w: time zone contains an invalid BCD digit", ErrInvalidValue)
	}
	quarterHours := int(tens)*10 + int(units)
	if quarterHours > maximumTimeZoneQuarterHours {
		return 0, fmt.Errorf("%w: time zone has an invalid quarter-hour value", ErrInvalidValue)
	}
	offset := time.Duration(quarterHours) * 15 * time.Minute
	if encoded&0x08 != 0 {
		offset = -offset
	}
	return offset, nil
}

func encodeTimeZoneTimeComponent(value, maximum int) (byte, error) {
	if value < 0 || value > maximum {
		return 0, fmt.Errorf("%w: time component %d is outside 0..%d", ErrInvalidValue, value, maximum)
	}
	return byte(value%10)<<4 | byte(value/10), nil
}

func decodeTimeZoneTimeComponent(encoded byte, maximum int) (int, error) {
	units := encoded >> 4
	tens := encoded & 0x0f
	if units > 9 || tens > 9 {
		return 0, fmt.Errorf("%w: time contains an invalid BCD digit", ErrInvalidValue)
	}
	value := int(tens)*10 + int(units)
	if value > maximum {
		return 0, fmt.Errorf("%w: time component %d is outside 0..%d", ErrInvalidValue, value, maximum)
	}
	return value, nil
}
