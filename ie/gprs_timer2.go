package ie

import (
	"fmt"
	"io"
	"time"
)

// GPRSTimer2Fields is the structured value of the GPRS timer 2 IE (TS 24.301
// clause 9.9.3.16A). Its unit and value coding is identical to GPRS timer.
type GPRSTimer2Fields struct {
	Duration time.Duration
}

// NewT3402ValueForAttachRejectWithFields creates a tagged T3402 value for an
// attach reject IE from fields.
func NewT3402ValueForAttachRejectWithFields(fields *GPRSTimer2Fields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewT3402ValueForAttachReject(payload[0]), nil
}

// NewT3346ValueWithFields creates a tagged T3346 value IE from fields.
func NewT3346ValueWithFields(fields *GPRSTimer2Fields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewT3346Value(payload[0]), nil
}

// NewT3324ValueWithFields creates a tagged T3324 value IE from fields.
func NewT3324ValueWithFields(fields *GPRSTimer2Fields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewT3324Value(payload[0]), nil
}

// NewT3448ValueWithFields creates a tagged T3448 value IE from fields.
func NewT3448ValueWithFields(fields *GPRSTimer2Fields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewT3448Value(payload[0]), nil
}

// NewPositionalGPRSTimer2WithFields creates a positional GPRS timer 2 IE from
// fields.
func NewPositionalGPRSTimer2WithFields(fields *GPRSTimer2Fields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalGPRSTimer2(payload[0]), nil
}

// GPRSTimer2 returns the GPRS timer 2 IE as structured fields.
func (i *IE) GPRSTimer2() (*GPRSTimer2Fields, error) {
	if err := i.requireDefinition(GPRSTimer2Definition); err != nil {
		return nil, err
	}
	return ParseGPRSTimer2Fields(i.Payload)
}

// ParseGPRSTimer2Fields decodes a GPRS timer 2 value without its IEI and
// length octet.
func ParseGPRSTimer2Fields(payload []byte) (*GPRSTimer2Fields, error) {
	fields := &GPRSTimer2Fields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Duration is representable by the GPRS timer 2 coding.
func (f *GPRSTimer2Fields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: GPRS timer 2 fields are nil", ErrInvalidValue)
	}
	_, err := encodeGPRSTimerDuration(f.Duration, "GPRS timer 2")
	return err
}

// Marshal encodes fields into a GPRS timer 2 value without its IEI.
func (f *GPRSTimer2Fields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, fmt.Errorf("%w: GPRS timer 2 fields are nil", ErrInvalidValue)
	}
	encoded, err := encodeGPRSTimerDuration(f.Duration, "GPRS timer 2")
	if err != nil {
		return nil, err
	}
	return []byte{encoded}, nil
}

// MarshalTo encodes fields into b.
func (f *GPRSTimer2Fields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the GPRS timer 2 value payload.
func (f *GPRSTimer2Fields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a GPRS timer 2 value without its IEI and length
// octet.
func (f *GPRSTimer2Fields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: GPRS timer 2 fields are nil", ErrInvalidValue)
	}
	duration, err := decodeGPRSTimerDuration(payload, "GPRS timer 2")
	if err != nil {
		return err
	}
	*f = GPRSTimer2Fields{Duration: duration}
	return nil
}
