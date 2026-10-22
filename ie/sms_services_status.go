package ie

import (
	"fmt"
	"io"
)

// SMSServicesStatusValue identifies the SMS services availability status (TS 24.301
// clause 9.9.3.4B).
type SMSServicesStatusValue uint8

const (
	SMSServicesStatusNotAvailable           SMSServicesStatusValue = iota
	SMSServicesStatusNotAvailableInPLMN
	SMSServicesStatusNetworkFailure
	SMSServicesStatusCongestion
)

// SMSServicesStatusFields is the structured value of the SMS services status IE
// (TS 24.301 clause 9.9.3.4B).
type SMSServicesStatusFields struct {
	Value SMSServicesStatusValue
}

// NewSMSServicesStatusWithFields creates a tagged SMS services status IE from
// fields.
func NewSMSServicesStatusWithFields(fields *SMSServicesStatusFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewSMSServicesStatus(payload[0]), nil
}

// SMSServicesStatus returns the IE as structured fields.
func (i *IE) SMSServicesStatus() (*SMSServicesStatusFields, error) {
	if err := i.requireDefinition(SMSServicesStatusDefinition); err != nil {
		return nil, err
	}
	return ParseSMSServicesStatusFields(i.Payload)
}

// ParseSMSServicesStatusFields decodes an SMS services status value without its IEI.
func ParseSMSServicesStatusFields(payload []byte) (*SMSServicesStatusFields, error) {
	fields := &SMSServicesStatusFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Value is a defined status.
func (f *SMSServicesStatusFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: SMS services status fields are nil", ErrInvalidValue)
	}
	if f.Value > SMSServicesStatusCongestion {
		return fmt.Errorf("%w: SMS services status %d is reserved", ErrInvalidValue, f.Value)
	}
	return nil
}

// Marshal encodes fields into an SMS services status value, without its IEI.
func (f *SMSServicesStatusFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.Value)}, nil
}

// MarshalTo encodes fields into b.
func (f *SMSServicesStatusFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the SMS services status value payload.
func (f *SMSServicesStatusFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes an SMS services status value without its IEI.
func (f *SMSServicesStatusFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: SMS services status fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: SMS services status has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xf0 != 0 {
		return fmt.Errorf("%w: SMS services status spare bits are nonzero", ErrInvalidValue)
	}
	v := SMSServicesStatusValue(payload[0] & 0x07)
	if v > SMSServicesStatusCongestion {
		return fmt.Errorf("%w: SMS services status %d is reserved", ErrInvalidValue, v)
	}
	*f = SMSServicesStatusFields{
		Value: v,
	}
	return nil
}
