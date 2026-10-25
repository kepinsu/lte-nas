package ie

import (
	"fmt"
	"io"
)

// AdditionalUpdateResultValue identifies the result of a combined attach or tracking
// area updating procedure (TS 24.301 clause 9.9.3.0A).
type AdditionalUpdateResultValue uint8

const (
	AdditionalUpdateResultNoAdditionalInfo AdditionalUpdateResultValue = iota
	AdditionalUpdateResultCSFallbackNotPref
	AdditionalUpdateResultSMSOnly
)

// AdditionalUpdateResultFields is the structured value of the Additional update
// result IE (TS 24.301 clause 9.9.3.0A).
type AdditionalUpdateResultFields struct {
	Value AdditionalUpdateResultValue
}

// NewAdditionalUpdateResultWithFields creates a tagged Additional update result
// IE from fields.
func NewAdditionalUpdateResultWithFields(fields *AdditionalUpdateResultFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewAdditionalUpdateResult(payload[0]), nil
}

// NewPositionalAdditionalUpdateResultWithFields creates a positional Additional
// update result IE from fields.
func NewPositionalAdditionalUpdateResultWithFields(fields *AdditionalUpdateResultFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalAdditionalUpdateResult(payload[0]), nil
}

// AdditionalUpdateResult returns the IE as structured fields.
func (i *IE) AdditionalUpdateResult() (*AdditionalUpdateResultFields, error) {
	if err := i.requireDefinition(AdditionalUpdateResultDefinition); err != nil {
		return nil, err
	}
	return ParseAdditionalUpdateResultFields(i.Payload)
}

// ParseAdditionalUpdateResultFields decodes an Additional update result value
// without its IEI.
func ParseAdditionalUpdateResultFields(payload []byte) (*AdditionalUpdateResultFields, error) {
	fields := &AdditionalUpdateResultFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Value is a defined result.
func (f *AdditionalUpdateResultFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: additional update result fields are nil", ErrInvalidValue)
	}
	if f.Value > AdditionalUpdateResultSMSOnly {
		return fmt.Errorf("%w: additional update result %d is reserved", ErrInvalidValue, f.Value)
	}
	return nil
}

// Marshal encodes fields into an Additional update result value, without its IEI.
func (f *AdditionalUpdateResultFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.Value)}, nil
}

// MarshalTo encodes fields into b.
func (f *AdditionalUpdateResultFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the Additional update result value payload.
func (f *AdditionalUpdateResultFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes an Additional update result value without its IEI.
func (f *AdditionalUpdateResultFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: additional update result fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: additional update result has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xfc != 0 {
		return fmt.Errorf("%w: additional update result spare bits are nonzero", ErrInvalidValue)
	}
	*f = AdditionalUpdateResultFields{
		Value: AdditionalUpdateResultValue(payload[0] & 0x03),
	}
	return nil
}
