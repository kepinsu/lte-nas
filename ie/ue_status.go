package ie

import (
	"fmt"
	"io"
)

// UEStatusFields is the structured value of the UE status IE (TS 24.301 clause
// 9.9.3.54, see TS 24.501 clause 9.11.3.56).
type UEStatusFields struct {
	// EMMRegistrationStatus indicates the EMM registration status.
	EMMRegistrationStatus bool
	// FiveGMMRegistrationStatus indicates the 5GMM registration status.
	FiveGMMRegistrationStatus bool
}

// NewUEStatusWithFields creates a tagged UE status IE from fields.
func NewUEStatusWithFields(fields *UEStatusFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewUEStatus(payload[0]), nil
}

// UEStatus returns the IE as structured fields.
func (i *IE) UEStatus() (*UEStatusFields, error) {
	if err := i.requireDefinition(UEStatusDefinition); err != nil {
		return nil, err
	}
	return ParseUEStatusFields(i.Payload)
}

// ParseUEStatusFields decodes a UE status value without its IEI.
func ParseUEStatusFields(payload []byte) (*UEStatusFields, error) {
	fields := &UEStatusFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *UEStatusFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: UE status fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a UE status value, without its IEI.
func (f *UEStatusFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	var encoded byte
	if f.EMMRegistrationStatus {
		encoded |= 0x01
	}
	if f.FiveGMMRegistrationStatus {
		encoded |= 0x02
	}
	return []byte{encoded}, nil
}

// MarshalTo encodes fields into b.
func (f *UEStatusFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the UE status value payload.
func (f *UEStatusFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a UE status value without its IEI.
func (f *UEStatusFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: UE status fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: UE status has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xfc != 0 {
		return fmt.Errorf("%w: UE status spare bits are nonzero", ErrInvalidValue)
	}
	*f = UEStatusFields{
		EMMRegistrationStatus:     payload[0]&0x01 != 0,
		FiveGMMRegistrationStatus: payload[0]&0x02 != 0,
	}
	return nil
}
