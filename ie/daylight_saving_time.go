package ie

import (
	"fmt"
	"io"
)

// DaylightSavingTimeFields is the structured value of the Daylight saving time IE
// (TS 24.301 clause 9.9.3.6, see TS 24.008 clause 10.5.3.12).
type DaylightSavingTimeFields struct {
	// Offset is the daylight saving time offset in half-hours.
	Offset int8
}

// NewDaylightSavingTimeWithFields creates a tagged Daylight saving time IE from
// fields.
func NewDaylightSavingTimeWithFields(fields *DaylightSavingTimeFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewDaylightSavingTime(payload[0]), nil
}

// DaylightSavingTime returns the IE as structured fields.
func (i *IE) DaylightSavingTime() (*DaylightSavingTimeFields, error) {
	if err := i.requireDefinition(DaylightSavingTimeDefinition); err != nil {
		return nil, err
	}
	return ParseDaylightSavingTimeFields(i.Payload)
}

// ParseDaylightSavingTimeFields decodes a Daylight saving time value without its
// IEI.
func ParseDaylightSavingTimeFields(payload []byte) (*DaylightSavingTimeFields, error) {
	fields := &DaylightSavingTimeFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *DaylightSavingTimeFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: daylight saving time fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a Daylight saving time value, without its IEI.
func (f *DaylightSavingTimeFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	var encoded byte
	if f.Offset < 0 {
		encoded |= 0x80
		encoded |= byte(-f.Offset) & 0x7f
	} else {
		encoded |= byte(f.Offset) & 0x7f
	}
	return []byte{encoded}, nil
}

// MarshalTo encodes fields into b.
func (f *DaylightSavingTimeFields) MarshalTo(b []byte) error {
	payload, err := f.Marshal()
	if err != nil {
		return err
	}
	if len(b) < len(payload) {
		return io.ErrShortBuffer
	}
	b[0] = payload[0]
	return nil
}

// MarshalLen returns the size of the Daylight saving time value payload.
func (f *DaylightSavingTimeFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a Daylight saving time value without its IEI.
func (f *DaylightSavingTimeFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: daylight saving time fields are nil", ErrInvalidValue)
	}
	if len(payload) == 0 {
		return io.ErrShortBuffer
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: daylight saving time has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	*f = DaylightSavingTimeFields{
		Offset: int8(payload[0]) & 0x7f,
	}
	if payload[0]&0x80 != 0 {
		f.Offset = -f.Offset
	}
	return nil
}
