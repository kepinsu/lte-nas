package ie

import (
	"fmt"
	"io"
)

// NASMessageContainerFields is the structured value of the tagged NAS message
// container IE (TS 24.301 clause 9.9.3.22). It carries an opaque byte string
// between 2 and 251 octets.
type NASMessageContainerFields struct {
	Value []byte
}

// NewNASMessageContainerWithFields creates a tagged NAS message container IE from
// fields.
func NewNASMessageContainerWithFields(fields *NASMessageContainerFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewNASMessageContainer(payload), nil
}

// NewPositionalNASMessageContainerWithFields creates a positional NAS message
// container IE from fields.
func NewPositionalNASMessageContainerWithFields(fields *NASMessageContainerFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalNASMessageContainer(payload), nil
}

// NASMessageContainer returns the IE as structured fields.
func (i *IE) NASMessageContainer() (*NASMessageContainerFields, error) {
	if err := i.requireDefinition(NASMessageContainerDefinition); err != nil {
		return nil, err
	}
	return ParseNASMessageContainerFields(i.Payload)
}

// ParseNASMessageContainerFields decodes a NAS message container value without its
// IEI and length octet.
func ParseNASMessageContainerFields(payload []byte) (*NASMessageContainerFields, error) {
	fields := &NASMessageContainerFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the payload is within [2, 251] octets.
func (f *NASMessageContainerFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: NAS message container fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 2 {
		return fmt.Errorf("%w: NAS message container has %d value octets, want >= 2", ErrInvalidLength, n)
	}
	if n > 251 {
		return fmt.Errorf("%w: NAS message container has %d value octets, want <= 251", ErrInvalidLength, n)
	}
	return nil
}

// Marshal encodes fields into a NAS message container value without its IEI and
// length octet.
func (f *NASMessageContainerFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *NASMessageContainerFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the NAS message container value payload, or zero
// for a nil receiver or nil Value.
func (f *NASMessageContainerFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes a NAS message container value without its IEI and length
// octet. It copies the payload into a temporary struct first, leaving f unchanged
// on failure.
func (f *NASMessageContainerFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: NAS message container fields are nil", ErrInvalidValue)
	}
	if len(payload) < 2 {
		return fmt.Errorf("%w: NAS message container has %d value octets, want >= 2", ErrInvalidLength, len(payload))
	}
	if len(payload) > 251 {
		return fmt.Errorf("%w: NAS message container has %d value octets, want <= 251", ErrInvalidLength, len(payload))
	}
	var decoded NASMessageContainerFields
	decoded.Value = make([]byte, len(payload))
	copy(decoded.Value, payload)
	*f = decoded
	return nil
}
