package ie

import (
	"fmt"
	"io"
)

// HeaderCompressionConfigurationFields is the structured value of the header
// compression configuration IE (TS 24.301 clause 9.9.4.22). It carries an
// opaque byte string with a minimum of 3 octets and a maximum of 255.
type HeaderCompressionConfigurationFields struct {
	Value []byte
}

// NewHeaderCompressionConfigurationWithFields creates a tagged Header
// compression configuration IE from fields.
func NewHeaderCompressionConfigurationWithFields(fields *HeaderCompressionConfigurationFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewHeaderCompressionConfiguration(payload), nil
}

// NewPositionalHeaderCompressionConfigurationWithFields creates a positional
// Header compression configuration IE from fields.
func NewPositionalHeaderCompressionConfigurationWithFields(fields *HeaderCompressionConfigurationFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalHeaderCompressionConfiguration(payload), nil
}

// HeaderCompressionConfiguration returns the IE as structured fields.
func (i *IE) HeaderCompressionConfiguration() (*HeaderCompressionConfigurationFields, error) {
	if err := i.requireDefinition(HeaderCompressionConfigurationDefinition); err != nil {
		return nil, err
	}
	return ParseHeaderCompressionConfigurationFields(i.Payload)
}

// ParseHeaderCompressionConfigurationFields decodes a header compression
// configuration value without its IEI and length octet.
func ParseHeaderCompressionConfigurationFields(payload []byte) (*HeaderCompressionConfigurationFields, error) {
	fields := &HeaderCompressionConfigurationFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the payload is within [3, 255] octets.
func (f *HeaderCompressionConfigurationFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: header compression configuration fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 3 {
		return fmt.Errorf("%w: header compression configuration has %d value octets, want >= 3", ErrInvalidLength, n)
	}
	if n > 255 {
		return fmt.Errorf("%w: header compression configuration has %d value octets, want <= 255", ErrInvalidLength, n)
	}
	return nil
}

// Marshal encodes fields into a header compression configuration value without
// its IEI and length octet.
func (f *HeaderCompressionConfigurationFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *HeaderCompressionConfigurationFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the header compression configuration value
// payload, or zero for a nil receiver or nil Value.
func (f *HeaderCompressionConfigurationFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes a header compression configuration value without its
// IEI and length octet. It copies the payload into a temporary struct first,
// leaving f unchanged on failure.
func (f *HeaderCompressionConfigurationFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: header compression configuration fields are nil", ErrInvalidValue)
	}
	if len(payload) < 3 {
		return fmt.Errorf("%w: header compression configuration has %d value octets, want >= 3", ErrInvalidLength, len(payload))
	}
	if len(payload) > 255 {
		return fmt.Errorf("%w: header compression configuration has %d value octets, want <= 255", ErrInvalidLength, len(payload))
	}
	var decoded HeaderCompressionConfigurationFields
	decoded.Value = make([]byte, len(payload))
	copy(decoded.Value, payload)
	*f = decoded
	return nil
}
