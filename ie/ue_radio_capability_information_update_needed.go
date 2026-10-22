package ie

import (
	"fmt"
	"io"
)

// UERadioCapabilityInformationUpdateNeededFields is the structured value of the
// UE radio capability information update needed IE (TS 24.301 clause 9.9.3.35).
type UERadioCapabilityInformationUpdateNeededFields struct {
	// URCUpd indicates whether UE radio capability information update is needed.
	URCUpd bool
}

// NewUERadioCapabilityInformationUpdateNeededWithFields creates a tagged UE radio
// capability information update needed IE from fields.
func NewUERadioCapabilityInformationUpdateNeededWithFields(fields *UERadioCapabilityInformationUpdateNeededFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewUERadioCapabilityInformationUpdateNeeded(payload[0]), nil
}

// UERadioCapabilityInformationUpdateNeeded returns the IE as structured fields.
func (i *IE) UERadioCapabilityInformationUpdateNeeded() (*UERadioCapabilityInformationUpdateNeededFields, error) {
	if err := i.requireDefinition(UERadioCapabilityInformationUpdateNeededDefinition); err != nil {
		return nil, err
	}
	return ParseUERadioCapabilityInformationUpdateNeededFields(i.Payload)
}

// ParseUERadioCapabilityInformationUpdateNeededFields decodes a UE radio
// capability information update needed value without its IEI.
func ParseUERadioCapabilityInformationUpdateNeededFields(payload []byte) (*UERadioCapabilityInformationUpdateNeededFields, error) {
	fields := &UERadioCapabilityInformationUpdateNeededFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *UERadioCapabilityInformationUpdateNeededFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: UE radio capability information update needed fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a UE radio capability information update needed
// value, without its IEI.
func (f *UERadioCapabilityInformationUpdateNeededFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	var encoded byte
	if f.URCUpd {
		encoded |= 0x01
	}
	return []byte{encoded}, nil
}

// MarshalTo encodes fields into b.
func (f *UERadioCapabilityInformationUpdateNeededFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the UE radio capability information update needed
// value payload.
func (f *UERadioCapabilityInformationUpdateNeededFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a UE radio capability information update needed value
// without its IEI.
func (f *UERadioCapabilityInformationUpdateNeededFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: UE radio capability information update needed fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: UE radio capability information update needed has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xfe != 0 {
		return fmt.Errorf("%w: UE radio capability information update needed spare bits are nonzero", ErrInvalidValue)
	}
	*f = UERadioCapabilityInformationUpdateNeededFields{
		URCUpd: payload[0]&0x01 != 0,
	}
	return nil
}
