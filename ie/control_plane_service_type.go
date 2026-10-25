package ie

import (
	"fmt"
	"io"
)

// ControlPlaneServiceType identifies the control plane service type, as defined
// in TS 24.301 clause 9.9.3.47.
type ControlPlaneServiceType uint8

// Control plane service type values.
const (
	ControlPlaneServiceTypeReserved        ControlPlaneServiceType = 0x00
	ControlPlaneServiceTypeNASRegistration ControlPlaneServiceType = 0x01
	ControlPlaneServiceTypeNASServices     ControlPlaneServiceType = 0x02
)

// ControlPlaneServiceTypeFields is the structured value of the positional Control
// plane service type IE.
type ControlPlaneServiceTypeFields struct {
	Type ControlPlaneServiceType
}

// NewPositionalControlPlaneServiceTypeWithFields creates a positional Control
// plane service type IE from fields.
func NewPositionalControlPlaneServiceTypeWithFields(fields *ControlPlaneServiceTypeFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalControlPlaneServiceType(payload[0]), nil
}

// ControlPlaneServiceType returns the IE as structured fields.
func (i *IE) ControlPlaneServiceType() (*ControlPlaneServiceTypeFields, error) {
	if err := i.requireDefinition(ControlPlaneServiceTypeDefinition); err != nil {
		return nil, err
	}
	return ParseControlPlaneServiceTypeFields(i.Payload)
}

// ParseControlPlaneServiceTypeFields decodes a Control plane service type value
// without its IEI.
func ParseControlPlaneServiceTypeFields(payload []byte) (*ControlPlaneServiceTypeFields, error) {
	fields := &ControlPlaneServiceTypeFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Type is a defined Control plane service type value.
func (f *ControlPlaneServiceTypeFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: Control plane service type fields are nil", ErrInvalidValue)
	}
	if !f.Type.valid() {
		return fmt.Errorf("%w: Control plane service type %#02x is reserved", ErrInvalidValue, f.Type)
	}
	return nil
}

func (t ControlPlaneServiceType) valid() bool {
	switch t {
	case ControlPlaneServiceTypeNASRegistration, ControlPlaneServiceTypeNASServices:
		return true
	default:
		return false
	}
}

// Marshal encodes fields into a Control plane service type value without its IEI.
func (f *ControlPlaneServiceTypeFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.Type)}, nil
}

// MarshalTo encodes fields into b.
func (f *ControlPlaneServiceTypeFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	b[0] = byte(f.Type)
	return nil
}

// MarshalLen returns the size of the Control plane service type value payload.
func (f *ControlPlaneServiceTypeFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a Control plane service type value without its IEI.
func (f *ControlPlaneServiceTypeFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: Control plane service type fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: Control plane service type has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0] != 0x01 && payload[0] != 0x02 {
		return fmt.Errorf("%w: Control plane service type is reserved", ErrInvalidValue)
	}
	decoded := ControlPlaneServiceTypeFields{Type: ControlPlaneServiceType(payload[0])}
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}
