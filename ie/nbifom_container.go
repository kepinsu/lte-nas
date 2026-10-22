package ie

import (
	"fmt"
	"io"
)

// NBIFOMContainerFields is the structured value of the NBIFOM container IE (TS
// 24.301 clause 9.9.4.19). It carries an opaque byte string between 1 and 255
// octets.
type NBIFOMContainerFields struct {
	Value []byte
}

// NewNBIFOMContainerWithFields creates a tagged NBIFOM container IE from
// fields.
func NewNBIFOMContainerWithFields(fields *NBIFOMContainerFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewNBIFOMContainer(payload), nil
}

// NewPositionalNBIFOMContainerWithFields creates a positional NBIFOM container
// IE from fields.
func NewPositionalNBIFOMContainerWithFields(fields *NBIFOMContainerFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalNBIFOMContainer(payload), nil
}

// NBIFOMContainer returns the IE as structured fields.
func (i *IE) NBIFOMContainer() (*NBIFOMContainerFields, error) {
	if err := i.requireDefinition(NBIFOMContainerDefinition); err != nil {
		return nil, err
	}
	return ParseNBIFOMContainerFields(i.Payload)
}

// ParseNBIFOMContainerFields decodes an NBIFOM container value without its IEI
// and length octet.
func ParseNBIFOMContainerFields(payload []byte) (*NBIFOMContainerFields, error) {
	fields := &NBIFOMContainerFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the payload is within [1, 255] octets.
func (f *NBIFOMContainerFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: NBIFOM container fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 1 {
		return fmt.Errorf("%w: NBIFOM container has %d value octets, want >= 1", ErrInvalidLength, n)
	}
	if n > 255 {
		return fmt.Errorf("%w: NBIFOM container has %d value octets, want <= 255", ErrInvalidLength, n)
	}
	return nil
}

// Marshal encodes fields into an NBIFOM container value without its IEI and
// length octet.
func (f *NBIFOMContainerFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *NBIFOMContainerFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the NBIFOM container value payload, or zero
// for a nil receiver or nil Value.
func (f *NBIFOMContainerFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes an NBIFOM container value without its IEI and length
// octet. It copies the payload into a temporary struct first, leaving f
// unchanged on failure.
func (f *NBIFOMContainerFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: NBIFOM container fields are nil", ErrInvalidValue)
	}
	if len(payload) < 1 {
		return fmt.Errorf("%w: NBIFOM container has %d value octets, want >= 1", ErrInvalidLength, len(payload))
	}
	if len(payload) > 255 {
		return fmt.Errorf("%w: NBIFOM container has %d value octets, want <= 255", ErrInvalidLength, len(payload))
	}
	var decoded NBIFOMContainerFields
	decoded.Value = make([]byte, len(payload))
	copy(decoded.Value, payload)
	*f = decoded
	return nil
}
