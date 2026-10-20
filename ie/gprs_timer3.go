package ie

import (
	"fmt"
	"io"
	"time"
)

// GPRSTimer3Fields is the structured value of the GPRS timer 3 IE (TS 24.301
// clause 9.9.3.16B). Duration is encoded with the most granular exact unit
// available in TS 24.008. A zero duration represents a deactivated timer.
type GPRSTimer3Fields struct {
	Duration time.Duration
}

// NewLowerBoundTimerValueWithFields creates a tagged lower bound timer value
// IE from fields.
func NewLowerBoundTimerValueWithFields(fields *GPRSTimer3Fields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewLowerBoundTimerValue(payload[0]), nil
}

// NewBackOffTimerValueWithFields creates a tagged back-off timer value IE from
// fields.
func NewBackOffTimerValueWithFields(fields *GPRSTimer3Fields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewBackOffTimerValue(payload[0]), nil
}

// NewMaximumTimeOffsetWithFields creates a tagged maximum time offset IE from
// fields.
func NewMaximumTimeOffsetWithFields(fields *GPRSTimer3Fields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewMaximumTimeOffset(payload[0]), nil
}

// NewT3412ExtendedValueWithFields creates a tagged T3412 extended value IE
// from fields.
func NewT3412ExtendedValueWithFields(fields *GPRSTimer3Fields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewT3412ExtendedValue(payload[0]), nil
}

// NewT3447ValueWithFields creates a tagged T3447 value IE from fields.
func NewT3447ValueWithFields(fields *GPRSTimer3Fields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewT3447Value(payload[0]), nil
}

// NewPositionalGPRSTimer3WithFields creates a positional GPRS timer 3 IE from
// fields.
func NewPositionalGPRSTimer3WithFields(fields *GPRSTimer3Fields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalGPRSTimer3(payload[0]), nil
}

// GPRSTimer3 returns the GPRS timer 3 IE as structured fields.
func (i *IE) GPRSTimer3() (*GPRSTimer3Fields, error) {
	if err := i.requireDefinition(GPRSTimer3Definition); err != nil {
		return nil, err
	}
	return ParseGPRSTimer3Fields(i.Payload)
}

// ParseGPRSTimer3Fields decodes a GPRS timer 3 value without its IEI and
// length octet.
func ParseGPRSTimer3Fields(payload []byte) (*GPRSTimer3Fields, error) {
	fields := &GPRSTimer3Fields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Duration is representable by the GPRS timer 3 coding.
func (f *GPRSTimer3Fields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: GPRS timer 3 fields are nil", ErrInvalidValue)
	}
	_, err := encodeGPRSTimer3Duration(f.Duration)
	return err
}

// Marshal encodes fields into a GPRS timer 3 value without its IEI and length
// octet.
func (f *GPRSTimer3Fields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, fmt.Errorf("%w: GPRS timer 3 fields are nil", ErrInvalidValue)
	}
	encoded, err := encodeGPRSTimer3Duration(f.Duration)
	if err != nil {
		return nil, err
	}
	return []byte{encoded}, nil
}

// MarshalTo encodes fields into b.
func (f *GPRSTimer3Fields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the GPRS timer 3 value payload.
func (f *GPRSTimer3Fields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a GPRS timer 3 value without its IEI and length
// octet.
func (f *GPRSTimer3Fields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: GPRS timer 3 fields are nil", ErrInvalidValue)
	}
	duration, err := decodeGPRSTimer3Duration(payload)
	if err != nil {
		return err
	}
	*f = GPRSTimer3Fields{Duration: duration}
	return nil
}

func encodeGPRSTimer3Duration(duration time.Duration) (byte, error) {
	if duration < 0 {
		return 0, fmt.Errorf("%w: GPRS timer 3 duration cannot be negative", ErrInvalidValue)
	}
	if duration == 0 {
		return 0xe0, nil
	}
	for _, option := range [...]struct {
		unit     byte
		duration time.Duration
	}{
		{unit: 3, duration: 2 * time.Second},
		{unit: 4, duration: 30 * time.Second},
		{unit: 5, duration: time.Minute},
		{unit: 0, duration: 10 * time.Minute},
		{unit: 1, duration: time.Hour},
		{unit: 2, duration: 10 * time.Hour},
		{unit: 6, duration: 320 * time.Hour},
	} {
		if duration%option.duration != 0 {
			continue
		}
		value := duration / option.duration
		if value <= 31 {
			return option.unit<<5 | byte(value), nil
		}
	}
	return 0, fmt.Errorf("%w: GPRS timer 3 duration %s is not representable", ErrInvalidValue, duration)
}

func decodeGPRSTimer3Duration(payload []byte) (time.Duration, error) {
	unit, value, err := decodeGPRSTimerValue(payload, "GPRS timer 3")
	if err != nil {
		return 0, err
	}
	switch unit {
	case 0:
		return time.Duration(value) * 10 * time.Minute, nil
	case 1:
		return time.Duration(value) * time.Hour, nil
	case 2:
		return time.Duration(value) * 10 * time.Hour, nil
	case 3:
		return time.Duration(value) * 2 * time.Second, nil
	case 4:
		return time.Duration(value) * 30 * time.Second, nil
	case 5:
		return time.Duration(value) * time.Minute, nil
	case 6:
		return time.Duration(value) * 320 * time.Hour, nil
	case 7:
		return 0, nil
	default:
		return 0, fmt.Errorf("%w: GPRS timer 3 has an invalid unit", ErrInvalidValue)
	}
}
