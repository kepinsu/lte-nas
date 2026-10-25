package ie

import (
	"fmt"
	"io"
)

// UnavailabilityInformationFields is the structured value of the unavailability
// information IE (TS 24.301 clause 9.9.3.69). It carries an opaque byte string
// between 1 and 7 octets.
type UnavailabilityInformationFields struct {
	Value []byte
}

// NewUnavailabilityInformationWithFields creates a tagged Unavailability
// information IE from fields.
func NewUnavailabilityInformationWithFields(fields *UnavailabilityInformationFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewUnavailabilityInformation(payload), nil
}

// NewPositionalUnavailabilityInformationWithFields creates a positional
// Unavailability information IE from fields.
func NewPositionalUnavailabilityInformationWithFields(fields *UnavailabilityInformationFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalUnavailabilityInformation(payload), nil
}

// UnavailabilityInformation returns the IE as structured fields.
func (i *IE) UnavailabilityInformation() (*UnavailabilityInformationFields, error) {
	if err := i.requireDefinition(UnavailabilityInformationDefinition); err != nil {
		return nil, err
	}
	return ParseUnavailabilityInformationFields(i.Payload)
}

// ParseUnavailabilityInformationFields decodes an unavailability information value
// without its IEI and length octet.
func ParseUnavailabilityInformationFields(payload []byte) (*UnavailabilityInformationFields, error) {
	fields := &UnavailabilityInformationFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the payload is within [1, 7] octets.
func (f *UnavailabilityInformationFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: unavailability information fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 1 {
		return fmt.Errorf("%w: unavailability information has %d value octets, want >= 1", ErrInvalidLength, n)
	}
	if n > 7 {
		return fmt.Errorf("%w: unavailability information has %d value octets, want <= 7", ErrInvalidLength, n)
	}
	return nil
}

// Marshal encodes fields into an unavailability information value without its IEI
// and length octet.
func (f *UnavailabilityInformationFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *UnavailabilityInformationFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the unavailability information value payload, or
// zero for a nil receiver or nil Value.
func (f *UnavailabilityInformationFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes an unavailability information value without its IEI and
// length octet. It copies the payload into a temporary struct first, leaving f
// unchanged on failure.
func (f *UnavailabilityInformationFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: unavailability information fields are nil", ErrInvalidValue)
	}
	if len(payload) < 1 {
		return fmt.Errorf("%w: unavailability information has %d value octets, want >= 1", ErrInvalidLength, len(payload))
	}
	if len(payload) > 7 {
		return fmt.Errorf("%w: unavailability information has %d value octets, want <= 7", ErrInvalidLength, len(payload))
	}
	var decoded UnavailabilityInformationFields
	decoded.Value = make([]byte, len(payload))
	copy(decoded.Value, payload)
	*f = decoded
	return nil
}
