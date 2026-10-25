package ie

import (
	"fmt"
	"io"
	"time"
)

// GPRSTimerFields is the structured value of the GPRS timer IE (TS 24.301
// clause 9.9.3.16). Duration is encoded with the most granular exact unit
// available in TS 24.008. A zero duration represents a deactivated timer.
type GPRSTimerFields struct {
	Duration time.Duration
}

// NewT3402ValueWithFields creates a tagged T3402 value IE from fields.
func NewT3402ValueWithFields(fields *GPRSTimerFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewT3402Value(payload[0]), nil
}

// NewT3423ValueWithFields creates a tagged T3423 value IE from fields.
func NewT3423ValueWithFields(fields *GPRSTimerFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewT3423Value(payload[0]), nil
}

// NewT3412ValueWithFields creates a tagged T3412 value IE from fields.
func NewT3412ValueWithFields(fields *GPRSTimerFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewT3412Value(payload[0]), nil
}

// NewT3442ValueWithFields creates a tagged T3442 value IE from fields.
func NewT3442ValueWithFields(fields *GPRSTimerFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewT3442Value(payload[0]), nil
}

// NewPositionalGPRSTimerWithFields creates a positional GPRS timer IE from
// fields.
func NewPositionalGPRSTimerWithFields(fields *GPRSTimerFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalGPRSTimer(payload[0]), nil
}

// GPRSTimer returns the GPRS timer IE as structured fields.
func (i *IE) GPRSTimer() (*GPRSTimerFields, error) {
	if err := i.requireDefinition(GPRSTimerDefinition); err != nil {
		return nil, err
	}
	return ParseGPRSTimerFields(i.Payload)
}

// ParseGPRSTimerFields decodes a GPRS timer value without its IEI.
func ParseGPRSTimerFields(payload []byte) (*GPRSTimerFields, error) {
	fields := &GPRSTimerFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Duration is representable by the GPRS timer coding.
func (f *GPRSTimerFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: GPRS timer fields are nil", ErrInvalidValue)
	}
	_, err := encodeGPRSTimerDuration(f.Duration, "GPRS timer")
	return err
}

// Marshal encodes fields into a GPRS timer value without its IEI.
func (f *GPRSTimerFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, fmt.Errorf("%w: GPRS timer fields are nil", ErrInvalidValue)
	}
	encoded, err := encodeGPRSTimerDuration(f.Duration, "GPRS timer")
	if err != nil {
		return nil, err
	}
	return []byte{encoded}, nil
}

// MarshalTo encodes fields into b.
func (f *GPRSTimerFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the GPRS timer value payload.
func (f *GPRSTimerFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a GPRS timer value without its IEI.
func (f *GPRSTimerFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: GPRS timer fields are nil", ErrInvalidValue)
	}
	duration, err := decodeGPRSTimerDuration(payload, "GPRS timer")
	if err != nil {
		return err
	}
	*f = GPRSTimerFields{Duration: duration}
	return nil
}

func decodeGPRSTimerValue(payload []byte, name string) (uint8, uint8, error) {
	if len(payload) != 1 {
		return 0, 0, fmt.Errorf("%w: %s has %d value octets, want 1", ErrInvalidLength, name, len(payload))
	}
	return payload[0] >> 5, payload[0] & 0x1f, nil
}

func encodeGPRSTimerDuration(duration time.Duration, name string) (byte, error) {
	if duration < 0 {
		return 0, fmt.Errorf("%w: %s duration cannot be negative", ErrInvalidValue, name)
	}
	if duration == 0 {
		return 0xe0, nil
	}
	for _, option := range [...]struct {
		unit     byte
		duration time.Duration
	}{
		{unit: 0, duration: 2 * time.Second},
		{unit: 1, duration: time.Minute},
		{unit: 2, duration: 6 * time.Minute},
	} {
		if duration%option.duration != 0 {
			continue
		}
		value := duration / option.duration
		if value <= 31 {
			return option.unit<<5 | byte(value), nil
		}
	}
	return 0, fmt.Errorf("%w: %s duration %s is not representable", ErrInvalidValue, name, duration)
}

func decodeGPRSTimerDuration(payload []byte, name string) (time.Duration, error) {
	unit, value, err := decodeGPRSTimerValue(payload, name)
	if err != nil {
		return 0, err
	}
	switch unit {
	case 7:
		return 0, nil
	case 0:
		return time.Duration(value) * 2 * time.Second, nil
	case 2:
		return time.Duration(value) * 6 * time.Minute, nil
	default:
		return time.Duration(value) * time.Minute, nil
	}
}
