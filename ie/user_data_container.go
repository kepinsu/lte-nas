package ie

import (
	"fmt"
	"io"
)

// UserDataContainerFields is the structured value of the user data container IE
// (TS 24.301 clause 9.9.4.24). It carries an opaque byte string between 0 and
// 65535 octets.
type UserDataContainerFields struct {
	Value []byte
}

// NewPositionalUserDataContainerWithFields creates a positional User data container
// IE from fields.
func NewPositionalUserDataContainerWithFields(fields *UserDataContainerFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalUserDataContainer(payload), nil
}

// UserDataContainer returns the IE as structured fields.
func (i *IE) UserDataContainer() (*UserDataContainerFields, error) {
	if err := i.requireDefinition(UserDataContainerDefinition); err != nil {
		return nil, err
	}
	return ParseUserDataContainerFields(i.Payload)
}

// ParseUserDataContainerFields decodes a user data container value without its IEI
// and length octet.
func ParseUserDataContainerFields(payload []byte) (*UserDataContainerFields, error) {
	fields := &UserDataContainerFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the payload is within [0, 65535] octets.
func (f *UserDataContainerFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: user data container fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 0 {
		return fmt.Errorf("%w: user data container has %d value octets, want >= 0", ErrInvalidLength, n)
	}
	if n > 65535 {
		return fmt.Errorf("%w: user data container has %d value octets, want <= 65535", ErrInvalidLength, n)
	}
	return nil
}

// Marshal encodes fields into a user data container value without its IEI and
// length octet.
func (f *UserDataContainerFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *UserDataContainerFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the user data container value payload, or zero for
// a nil receiver or nil Value.
func (f *UserDataContainerFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes a user data container value without its IEI and length
// octet. It copies the payload into a temporary struct first, leaving f unchanged
// on failure.
func (f *UserDataContainerFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: user data container fields are nil", ErrInvalidValue)
	}
	if len(payload) > 65535 {
		return fmt.Errorf("%w: user data container has %d value octets, want <= 65535", ErrInvalidLength, len(payload))
	}
	var decoded UserDataContainerFields
	if payload == nil {
		decoded.Value = nil
	} else {
		decoded.Value = make([]byte, len(payload))
		copy(decoded.Value, payload)
	}
	*f = decoded
	return nil
}
