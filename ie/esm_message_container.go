package ie

import (
	"fmt"
	"io"
)

// ESMMessageContainerFields is the structured value of the tagged ESM message
// container IE (TS 24.301 clause 9.9.3.15). It carries an opaque byte string
// with a minimum length of 0 octets.
type ESMMessageContainerFields struct {
	Value []byte
}

// NewPositionalESMMessageContainerWithFields creates a positional ESM message
// container IE from fields.
func NewPositionalESMMessageContainerWithFields(fields *ESMMessageContainerFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalESMMessageContainer(payload), nil
}

// ESMMessageContainer returns the IE as structured fields.
func (i *IE) ESMMessageContainer() (*ESMMessageContainerFields, error) {
	if err := i.requireDefinition(ESMMessageContainerDefinition); err != nil {
		return nil, err
	}
	return ParseESMMessageContainerFields(i.Payload)
}

// ParseESMMessageContainerFields decodes an ESM message container value without
// its IEI and length octet.
func ParseESMMessageContainerFields(payload []byte) (*ESMMessageContainerFields, error) {
	fields := &ESMMessageContainerFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the payload is within valid bounds.
func (f *ESMMessageContainerFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: ESM message container fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into an ESM message container value without its IEI and
// length octet.
func (f *ESMMessageContainerFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *ESMMessageContainerFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the ESM message container value payload, or zero
// for a nil receiver or nil Value.
func (f *ESMMessageContainerFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes an ESM message container value without its IEI and length
// octet. It copies the payload into a temporary struct first, leaving f unchanged
// on failure.
func (f *ESMMessageContainerFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: ESM message container fields are nil", ErrInvalidValue)
	}
	var decoded ESMMessageContainerFields
	decoded.Value = make([]byte, len(payload))
	copy(decoded.Value, payload)
	*f = decoded
	return nil
}
