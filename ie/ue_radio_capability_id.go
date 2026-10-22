package ie

import (
	"fmt"
	"io"
)

// UERadioCapabilityIDFields is the structured value of the tagged UE radio
// capability ID IE (TS 24.301 clause 9.9.3.60). It carries an opaque byte string
// between 1 and 255 octets.
type UERadioCapabilityIDFields struct {
	Value []byte
}

// NewUERadioCapabilityIDWithFields creates a tagged UE radio capability ID IE
// from fields.
func NewUERadioCapabilityIDWithFields(fields *UERadioCapabilityIDFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewUERadioCapabilityID(payload), nil
}

// NewPositionalUERadioCapabilityIDWithFields creates a positional UE radio
// capability ID IE from fields.
func NewPositionalUERadioCapabilityIDWithFields(fields *UERadioCapabilityIDFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalUERadioCapabilityID(payload), nil
}

// UERadioCapabilityID returns the IE as structured fields.
func (i *IE) UERadioCapabilityID() (*UERadioCapabilityIDFields, error) {
	if err := i.requireDefinition(UERadioCapabilityIDDefinition); err != nil {
		return nil, err
	}
	return ParseUERadioCapabilityIDFields(i.Payload)
}

// ParseUERadioCapabilityIDFields decodes a UE radio capability ID value without
// its IEI and length octet.
func ParseUERadioCapabilityIDFields(payload []byte) (*UERadioCapabilityIDFields, error) {
	fields := &UERadioCapabilityIDFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the payload is within [1, 255] octets.
func (f *UERadioCapabilityIDFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: UE radio capability ID fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 1 {
		return fmt.Errorf("%w: UE radio capability ID has %d value octets, want >= 1", ErrInvalidLength, n)
	}
	if n > 255 {
		return fmt.Errorf("%w: UE radio capability ID has %d value octets, want <= 255", ErrInvalidLength, n)
	}
	return nil
}

// Marshal encodes fields into a UE radio capability ID value without its IEI and
// length octet.
func (f *UERadioCapabilityIDFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *UERadioCapabilityIDFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the UE radio capability ID value payload, or zero
// for a nil receiver or nil Value.
func (f *UERadioCapabilityIDFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes a UE radio capability ID value without its IEI and length
// octet. It copies the payload into a temporary struct first, leaving f unchanged
// on failure.
func (f *UERadioCapabilityIDFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: UE radio capability ID fields are nil", ErrInvalidValue)
	}
	if len(payload) < 1 {
		return fmt.Errorf("%w: UE radio capability ID has %d value octets, want >= 1", ErrInvalidLength, len(payload))
	}
	if len(payload) > 255 {
		return fmt.Errorf("%w: UE radio capability ID has %d value octets, want <= 255", ErrInvalidLength, len(payload))
	}
	var decoded UERadioCapabilityIDFields
	decoded.Value = make([]byte, len(payload))
	copy(decoded.Value, payload)
	*f = decoded
	return nil
}
