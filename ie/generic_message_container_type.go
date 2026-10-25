package ie

import (
	"fmt"
	"io"
)

// GenericMessageContainerType identifies the generic message container type, as
// defined in TS 24.301 clause 9.9.3.42.
type GenericMessageContainerType uint8

// Generic message container type values.
const (
	GenericMessageContainerTypeReserved GenericMessageContainerType = 0x00
	GenericMessageContainerTypeESM      GenericMessageContainerType = 0x01
)

// GenericMessageContainerTypeFields is the structured value of the positional
// Generic message container type IE.
type GenericMessageContainerTypeFields struct {
	Type GenericMessageContainerType
}

// NewPositionalGenericMessageContainerTypeWithFields creates a positional
// Generic message container type IE from fields.
func NewPositionalGenericMessageContainerTypeWithFields(fields *GenericMessageContainerTypeFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalGenericMessageContainerType(payload[0]), nil
}

// GenericMessageContainerType returns the IE as structured fields.
func (i *IE) GenericMessageContainerType() (*GenericMessageContainerTypeFields, error) {
	if err := i.requireDefinition(GenericMessageContainerTypeDefinition); err != nil {
		return nil, err
	}
	return ParseGenericMessageContainerTypeFields(i.Payload)
}

// ParseGenericMessageContainerTypeFields decodes a Generic message container type
// value without its IEI.
func ParseGenericMessageContainerTypeFields(payload []byte) (*GenericMessageContainerTypeFields, error) {
	fields := &GenericMessageContainerTypeFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Type is a defined Generic message container type value.
func (f *GenericMessageContainerTypeFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: Generic message container type fields are nil", ErrInvalidValue)
	}
	if !f.Type.valid() {
		return fmt.Errorf("%w: Generic message container type %#02x is reserved", ErrInvalidValue, f.Type)
	}
	return nil
}

func (t GenericMessageContainerType) valid() bool {
	switch t {
	case GenericMessageContainerTypeESM:
		return true
	default:
		return false
	}
}

// Marshal encodes fields into a Generic message container type value without its
// IEI.
func (f *GenericMessageContainerTypeFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.Type)}, nil
}

// MarshalTo encodes fields into b.
func (f *GenericMessageContainerTypeFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	b[0] = byte(f.Type)
	return nil
}

// MarshalLen returns the size of the Generic message container type value payload.
func (f *GenericMessageContainerTypeFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a Generic message container type value without its IEI.
func (f *GenericMessageContainerTypeFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: Generic message container type fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: Generic message container type has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0] != 0x01 {
		return fmt.Errorf("%w: Generic message container type is reserved", ErrInvalidValue)
	}
	decoded := GenericMessageContainerTypeFields{Type: GenericMessageContainerType(payload[0])}
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}
