package ie

import (
	"fmt"
	"io"
	"time"
)

// TimeZoneAndTimeFields is the structured value of the Time zone and time IE
// (TS 24.301 clause 9.9.3.30). Time is the universal network time, represented
// in UTC; Offset is the local time-zone offset.
type TimeZoneAndTimeFields struct {
	Time   time.Time
	Offset time.Duration
}

// NewTimeZoneAndTimeWithFields creates a tagged Time zone and time IE from
// fields.
func NewTimeZoneAndTimeWithFields(fields *TimeZoneAndTimeFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewTimeZoneAndTime(payload), nil
}

// NewPositionalTimeZoneAndTimeWithFields creates a positional Time zone and
// time IE from fields.
func NewPositionalTimeZoneAndTimeWithFields(fields *TimeZoneAndTimeFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalTimeZoneAndTime(payload), nil
}

// TimeZoneAndTime returns the IE as structured fields.
func (i *IE) TimeZoneAndTime() (*TimeZoneAndTimeFields, error) {
	if err := i.requireDefinition(TimeZoneAndTimeDefinition); err != nil {
		return nil, err
	}
	return ParseTimeZoneAndTimeFields(i.Payload)
}

// ParseTimeZoneAndTimeFields decodes a Time zone and time value without its
// IEI.
func ParseTimeZoneAndTimeFields(payload []byte) (*TimeZoneAndTimeFields, error) {
	fields := &TimeZoneAndTimeFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Time and Offset are representable by the IE.
func (f *TimeZoneAndTimeFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: time zone and time fields are nil", ErrInvalidValue)
	}
	if f.Time.IsZero() {
		return fmt.Errorf("%w: universal network time is zero", ErrInvalidValue)
	}
	utc := f.Time.UTC()
	if utc.Year() < 2000 || utc.Year() > 2099 {
		return fmt.Errorf("%w: universal network time year %d is outside 2000..2099", ErrInvalidValue, utc.Year())
	}
	_, err := encodeTimeZone(f.Offset)
	return err
}

// Marshal encodes fields into a Time zone and time value without its IEI.
func (f *TimeZoneAndTimeFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	utc := f.Time.UTC()
	components := []struct {
		value   int
		maximum int
	}{
		{value: utc.Year() - 2000, maximum: 99},
		{value: int(utc.Month()), maximum: 12},
		{value: utc.Day(), maximum: 31},
		{value: utc.Hour(), maximum: 23},
		{value: utc.Minute(), maximum: 59},
		{value: utc.Second(), maximum: 59},
	}
	payload := make([]byte, 7)
	for index, component := range components {
		encoded, err := encodeTimeZoneTimeComponent(component.value, component.maximum)
		if err != nil {
			return nil, err
		}
		payload[index] = encoded
	}
	offset, err := encodeTimeZone(f.Offset)
	if err != nil {
		return nil, err
	}
	payload[6] = offset
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *TimeZoneAndTimeFields) MarshalTo(b []byte) error {
	payload, err := f.Marshal()
	if err != nil {
		return err
	}
	if len(b) < len(payload) {
		return io.ErrShortBuffer
	}
	copy(b, payload)
	return nil
}

// MarshalLen returns the size of the Time zone and time value payload.
func (f *TimeZoneAndTimeFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 7
}

// UnmarshalBinary decodes a Time zone and time value without its IEI.
func (f *TimeZoneAndTimeFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: time zone and time fields are nil", ErrInvalidValue)
	}
	if len(payload) != 7 {
		return fmt.Errorf("%w: time zone and time has %d value octets, want 7", ErrInvalidLength, len(payload))
	}
	maximums := [...]int{99, 12, 31, 23, 59, 59}
	values := make([]int, len(maximums))
	for index, maximum := range maximums {
		value, err := decodeTimeZoneTimeComponent(payload[index], maximum)
		if err != nil {
			return err
		}
		values[index] = value
	}
	offset, err := decodeTimeZone(payload[6])
	if err != nil {
		return err
	}
	decodedTime := time.Date(2000+values[0], time.Month(values[1]), values[2], values[3], values[4], values[5], 0, time.UTC)
	if decodedTime.Year() != 2000+values[0] || int(decodedTime.Month()) != values[1] || decodedTime.Day() != values[2] {
		return fmt.Errorf("%w: time zone and time contains an invalid calendar date", ErrInvalidValue)
	}
	*f = TimeZoneAndTimeFields{Time: decodedTime, Offset: offset}
	return nil
}
