package ie

import (
	"fmt"
	"io"
)

// UnavailabilityConfigurationFields is the structured value of the unavailability
// configuration IE (TS 24.301 clause 9.9.3.70). It carries an opaque byte string
// between 1 and 7 octets.
type UnavailabilityConfigurationFields struct {
	Value []byte
}

// NewUnavailabilityConfigurationWithFields creates a tagged Unavailability
// configuration IE from fields.
func NewUnavailabilityConfigurationWithFields(fields *UnavailabilityConfigurationFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewUnavailabilityConfiguration(payload), nil
}

// NewPositionalUnavailabilityConfigurationWithFields creates a positional
// Unavailability configuration IE from fields.
func NewPositionalUnavailabilityConfigurationWithFields(fields *UnavailabilityConfigurationFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalUnavailabilityConfiguration(payload), nil
}

// UnavailabilityConfiguration returns the IE as structured fields.
func (i *IE) UnavailabilityConfiguration() (*UnavailabilityConfigurationFields, error) {
	if err := i.requireDefinition(UnavailabilityConfigurationDefinition); err != nil {
		return nil, err
	}
	return ParseUnavailabilityConfigurationFields(i.Payload)
}

// ParseUnavailabilityConfigurationFields decodes an unavailability configuration
// value without its IEI and length octet.
func ParseUnavailabilityConfigurationFields(payload []byte) (*UnavailabilityConfigurationFields, error) {
	fields := &UnavailabilityConfigurationFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the payload is within [1, 7] octets.
func (f *UnavailabilityConfigurationFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: unavailability configuration fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 1 {
		return fmt.Errorf("%w: unavailability configuration has %d value octets, want >= 1", ErrInvalidLength, n)
	}
	if n > 7 {
		return fmt.Errorf("%w: unavailability configuration has %d value octets, want <= 7", ErrInvalidLength, n)
	}
	return nil
}

// Marshal encodes fields into an unavailability configuration value without its IEI
// and length octet.
func (f *UnavailabilityConfigurationFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *UnavailabilityConfigurationFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the unavailability configuration value payload, or
// zero for a nil receiver or nil Value.
func (f *UnavailabilityConfigurationFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes an unavailability configuration value without its IEI and
// length octet. It copies the payload into a temporary struct first, leaving f
// unchanged on failure.
func (f *UnavailabilityConfigurationFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: unavailability configuration fields are nil", ErrInvalidValue)
	}
	if len(payload) < 1 {
		return fmt.Errorf("%w: unavailability configuration has %d value octets, want >= 1", ErrInvalidLength, len(payload))
	}
	if len(payload) > 7 {
		return fmt.Errorf("%w: unavailability configuration has %d value octets, want <= 7", ErrInvalidLength, len(payload))
	}
	var decoded UnavailabilityConfigurationFields
	decoded.Value = make([]byte, len(payload))
	copy(decoded.Value, payload)
	*f = decoded
	return nil
}
