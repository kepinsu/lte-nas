package ie

import (
	"fmt"
	"io"
)

// GenericMessageContainerFields is the structured value of the generic message
// container IE (TS 24.301 clause 9.9.3.43). It carries an opaque byte string
// between 1 and 65535 octets.
type GenericMessageContainerFields struct {
	Value []byte
}

// NewPositionalGenericMessageContainerWithFields creates a positional Generic
// message container IE from fields.
func NewPositionalGenericMessageContainerWithFields(fields *GenericMessageContainerFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalGenericMessageContainer(payload), nil
}

// GenericMessageContainer returns the IE as structured fields.
func (i *IE) GenericMessageContainer() (*GenericMessageContainerFields, error) {
	if err := i.requireDefinition(GenericMessageContainerDefinition); err != nil {
		return nil, err
	}
	return ParseGenericMessageContainerFields(i.Payload)
}

// ParseGenericMessageContainerFields decodes a generic message container value
// without its IEI and length octet.
func ParseGenericMessageContainerFields(payload []byte) (*GenericMessageContainerFields, error) {
	fields := &GenericMessageContainerFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the payload is within [1, 65535] octets.
func (f *GenericMessageContainerFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: generic message container fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 1 {
		return fmt.Errorf("%w: generic message container has %d value octets, want >= 1", ErrInvalidLength, n)
	}
	if n > 65535 {
		return fmt.Errorf("%w: generic message container has %d value octets, want <= 65535", ErrInvalidLength, n)
	}
	return nil
}

// Marshal encodes fields into a generic message container value without its IEI
// and length octet.
func (f *GenericMessageContainerFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *GenericMessageContainerFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the generic message container value payload, or
// zero for a nil receiver or nil Value.
func (f *GenericMessageContainerFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes a generic message container value without its IEI and
// length octet. It copies the payload into a temporary struct first, leaving f
// unchanged on failure.
func (f *GenericMessageContainerFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: generic message container fields are nil", ErrInvalidValue)
	}
	if len(payload) < 1 {
		return fmt.Errorf("%w: generic message container has %d value octets, want >= 1", ErrInvalidLength, len(payload))
	}
	if len(payload) > 65535 {
		return fmt.Errorf("%w: generic message container has %d value octets, want <= 65535", ErrInvalidLength, len(payload))
	}
	var decoded GenericMessageContainerFields
	decoded.Value = make([]byte, len(payload))
	copy(decoded.Value, payload)
	*f = decoded
	return nil
}
