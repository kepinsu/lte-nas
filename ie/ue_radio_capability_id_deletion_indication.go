package ie

import (
	"fmt"
	"io"
)

// UERadioCapabilityIDDeletionIndicationFields is the structured value of the UE
// radio capability ID deletion indication IE (TS 24.301 clause 9.9.3.61, see TS
// 24.501 clause 9.11.3.69).
type UERadioCapabilityIDDeletionIndicationFields struct {
	// URCDI indicates whether network-assigned UE radio capability IDs should be
	// deleted.
	URCDI bool
}

// NewUERadioCapabilityIDDeletionIndicationWithFields creates a tagged UE radio
// capability ID deletion indication IE from fields.
func NewUERadioCapabilityIDDeletionIndicationWithFields(fields *UERadioCapabilityIDDeletionIndicationFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewUERadioCapabilityIDDeletionIndication(payload[0]), nil
}

// UERadioCapabilityIDDeletionIndication returns the IE as structured fields.
func (i *IE) UERadioCapabilityIDDeletionIndication() (*UERadioCapabilityIDDeletionIndicationFields, error) {
	if err := i.requireDefinition(UERadioCapabilityIDDeletionIndicationDefinition); err != nil {
		return nil, err
	}
	return ParseUERadioCapabilityIDDeletionIndicationFields(i.Payload)
}

// ParseUERadioCapabilityIDDeletionIndicationFields decodes a UE radio capability
// ID deletion indication value without its IEI.
func ParseUERadioCapabilityIDDeletionIndicationFields(payload []byte) (*UERadioCapabilityIDDeletionIndicationFields, error) {
	fields := &UERadioCapabilityIDDeletionIndicationFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *UERadioCapabilityIDDeletionIndicationFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: UE radio capability ID deletion indication fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a UE radio capability ID deletion indication value,
// without its IEI.
func (f *UERadioCapabilityIDDeletionIndicationFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	var encoded byte
	if f.URCDI {
		encoded |= 0x80
	}
	return []byte{encoded}, nil
}

// MarshalTo encodes fields into b.
func (f *UERadioCapabilityIDDeletionIndicationFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the UE radio capability ID deletion indication
// value payload.
func (f *UERadioCapabilityIDDeletionIndicationFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a UE radio capability ID deletion indication value
// without its IEI.
func (f *UERadioCapabilityIDDeletionIndicationFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: UE radio capability ID deletion indication fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: UE radio capability ID deletion indication has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0x7f != 0 {
		return fmt.Errorf("%w: UE radio capability ID deletion indication spare bits are nonzero", ErrInvalidValue)
	}
	*f = UERadioCapabilityIDDeletionIndicationFields{
		URCDI: payload[0]&0x80 != 0,
	}
	return nil
}
