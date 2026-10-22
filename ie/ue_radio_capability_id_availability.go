package ie

import (
	"fmt"
	"io"
)

// UERadioCapabilityIDAvailabilityFields is the structured value of the UE radio
// capability ID availability IE (TS 24.301 clause 9.9.3.58).
type UERadioCapabilityIDAvailabilityFields struct {
	// URCA indicates whether a UE radio capability ID is available.
	URCA bool
}

// NewUERadioCapabilityIDAvailabilityWithFields creates a tagged UE radio
// capability ID availability IE from fields.
func NewUERadioCapabilityIDAvailabilityWithFields(fields *UERadioCapabilityIDAvailabilityFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewUERadioCapabilityIDAvailability(payload[0]), nil
}

// UERadioCapabilityIDAvailability returns the IE as structured fields.
func (i *IE) UERadioCapabilityIDAvailability() (*UERadioCapabilityIDAvailabilityFields, error) {
	if err := i.requireDefinition(UERadioCapabilityIDAvailabilityDefinition); err != nil {
		return nil, err
	}
	return ParseUERadioCapabilityIDAvailabilityFields(i.Payload)
}

// ParseUERadioCapabilityIDAvailabilityFields decodes a UE radio capability ID
// availability value without its IEI.
func ParseUERadioCapabilityIDAvailabilityFields(payload []byte) (*UERadioCapabilityIDAvailabilityFields, error) {
	fields := &UERadioCapabilityIDAvailabilityFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *UERadioCapabilityIDAvailabilityFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: UE radio capability ID availability fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a UE radio capability ID availability value, without its
// IEI.
func (f *UERadioCapabilityIDAvailabilityFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	var encoded byte
	if f.URCA {
		encoded |= 0x02
	}
	return []byte{encoded}, nil
}

// MarshalTo encodes fields into b.
func (f *UERadioCapabilityIDAvailabilityFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the UE radio capability ID availability value payload.
func (f *UERadioCapabilityIDAvailabilityFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a UE radio capability ID availability value without its IEI.
func (f *UERadioCapabilityIDAvailabilityFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: UE radio capability ID availability fields are nil", ErrInvalidValue)
	}
	if len(payload) == 0 {
		return io.ErrShortBuffer
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: UE radio capability ID availability has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xf8 != 0 {
		return fmt.Errorf("%w: UE radio capability ID availability spare bits are nonzero", ErrInvalidValue)
	}
	*f = UERadioCapabilityIDAvailabilityFields{
		URCA: payload[0]&0x02 != 0,
	}
	return nil
}
