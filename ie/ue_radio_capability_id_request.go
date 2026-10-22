package ie

import (
	"fmt"
	"io"
)

// UERadioCapabilityIDRequestFields is the structured value of the UE radio
// capability ID request IE (TS 24.301 clause 9.9.3.59).
type UERadioCapabilityIDRequestFields struct {
	// URCIDR indicates whether the UE radio capability ID is requested.
	URCIDR bool
}

// NewUERadioCapabilityIDRequestWithFields creates a tagged UE radio capability
// ID request IE from fields.
func NewUERadioCapabilityIDRequestWithFields(fields *UERadioCapabilityIDRequestFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewUERadioCapabilityIDRequest(payload[0]), nil
}

// UERadioCapabilityIDRequest returns the IE as structured fields.
func (i *IE) UERadioCapabilityIDRequest() (*UERadioCapabilityIDRequestFields, error) {
	if err := i.requireDefinition(UERadioCapabilityIDRequestDefinition); err != nil {
		return nil, err
	}
	return ParseUERadioCapabilityIDRequestFields(i.Payload)
}

// ParseUERadioCapabilityIDRequestFields decodes a UE radio capability ID request
// value without its IEI.
func ParseUERadioCapabilityIDRequestFields(payload []byte) (*UERadioCapabilityIDRequestFields, error) {
	fields := &UERadioCapabilityIDRequestFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *UERadioCapabilityIDRequestFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: UE radio capability ID request fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a UE radio capability ID request value, without its IEI.
func (f *UERadioCapabilityIDRequestFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	var encoded byte
	if f.URCIDR {
		encoded |= 0x01
	}
	return []byte{encoded}, nil
}

// MarshalTo encodes fields into b.
func (f *UERadioCapabilityIDRequestFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the UE radio capability ID request value payload.
func (f *UERadioCapabilityIDRequestFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a UE radio capability ID request value without its IEI.
func (f *UERadioCapabilityIDRequestFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: UE radio capability ID request fields are nil", ErrInvalidValue)
	}
	if len(payload) == 0 {
		return io.ErrShortBuffer
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: UE radio capability ID request has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xfe != 0 {
		return fmt.Errorf("%w: UE radio capability ID request spare bits are nonzero", ErrInvalidValue)
	}
	*f = UERadioCapabilityIDRequestFields{
		URCIDR: payload[0]&0x01 != 0,
	}
	return nil
}
