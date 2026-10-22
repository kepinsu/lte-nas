package ie

import (
	"fmt"
	"io"
)

// ServingPLMNRateControlFields is the structured value of the serving PLMN rate
// control IE (TS 24.301 clause 9.9.4.28). The two-octet payload encodes a rate
// control period in milliseconds as a big-endian uint16.
type ServingPLMNRateControlFields struct {
	// RateControlPeriodMS is the rate control period in milliseconds.
	RateControlPeriodMS uint16
}

// NewServingPLMNRateControlWithFields creates a tagged serving PLMN rate control
// IE from fields.
func NewServingPLMNRateControlWithFields(fields *ServingPLMNRateControlFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewServingPLMNRateControl(payload), nil
}

// NewPositionalServingPLMNRateControlWithFields creates a positional serving PLMN
// rate control IE from fields.
func NewPositionalServingPLMNRateControlWithFields(fields *ServingPLMNRateControlFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalServingPLMNRateControl(payload), nil
}

// ServingPLMNRateControl returns the IE as structured fields.
func (i *IE) ServingPLMNRateControl() (*ServingPLMNRateControlFields, error) {
	if err := i.requireDefinition(ServingPLMNRateControlDefinition); err != nil {
		return nil, err
	}
	return ParseServingPLMNRateControlFields(i.Payload)
}

// ParseServingPLMNRateControlFields decodes a serving PLMN rate control value
// without its IEI and length octet.
func ParseServingPLMNRateControlFields(payload []byte) (*ServingPLMNRateControlFields, error) {
	fields := &ServingPLMNRateControlFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *ServingPLMNRateControlFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: serving PLMN rate control fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a 2-octet serving PLMN rate control value.
func (f *ServingPLMNRateControlFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, 2)
	payload[0] = byte(f.RateControlPeriodMS >> 8)
	payload[1] = byte(f.RateControlPeriodMS)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *ServingPLMNRateControlFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 2 {
		return io.ErrShortBuffer
	}
	b[0] = byte(f.RateControlPeriodMS >> 8)
	b[1] = byte(f.RateControlPeriodMS)
	return nil
}

// MarshalLen returns the size of the serving PLMN rate control value payload.
func (f *ServingPLMNRateControlFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 2
}

// UnmarshalBinary decodes a 2-octet serving PLMN rate control value.
func (f *ServingPLMNRateControlFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: serving PLMN rate control fields are nil", ErrInvalidValue)
	}
	if len(payload) != 2 {
		return fmt.Errorf("%w: serving PLMN rate control has %d value octets, want 2", ErrInvalidLength, len(payload))
	}
	var decoded ServingPLMNRateControlFields
	decoded.RateControlPeriodMS = uint16(payload[0])<<8 | uint16(payload[1])
	*f = decoded
	return nil
}
