package ie

import (
	"fmt"
	"io"
)

// ExtendedEmergencyNumberListFields is the structured value of the extended
// emergency number list IE (TS 24.301 clause 9.9.3.37A). It carries an opaque
// byte string between 4 and 65535 octets.
type ExtendedEmergencyNumberListFields struct {
	Value []byte
}

// NewPositionalExtendedEmergencyNumberListWithFields creates a positional
// Extended emergency number list IE from fields.
func NewPositionalExtendedEmergencyNumberListWithFields(fields *ExtendedEmergencyNumberListFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalExtendedEmergencyNumberList(payload), nil
}

// ExtendedEmergencyNumberList returns the IE as structured fields.
func (i *IE) ExtendedEmergencyNumberList() (*ExtendedEmergencyNumberListFields, error) {
	if err := i.requireDefinition(ExtendedEmergencyNumberListDefinition); err != nil {
		return nil, err
	}
	return ParseExtendedEmergencyNumberListFields(i.Payload)
}

// ParseExtendedEmergencyNumberListFields decodes an extended emergency number
// list value without its IEI and length octet.
func ParseExtendedEmergencyNumberListFields(payload []byte) (*ExtendedEmergencyNumberListFields, error) {
	fields := &ExtendedEmergencyNumberListFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the payload is within [4, 65535] octets.
func (f *ExtendedEmergencyNumberListFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: extended emergency number list fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 4 {
		return fmt.Errorf("%w: extended emergency number list has %d value octets, want >= 4", ErrInvalidLength, n)
	}
	if n > 65535 {
		return fmt.Errorf("%w: extended emergency number list has %d value octets, want <= 65535", ErrInvalidLength, n)
	}
	return nil
}

// Marshal encodes fields into an extended emergency number list value without
// its IEI and length octet.
func (f *ExtendedEmergencyNumberListFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *ExtendedEmergencyNumberListFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the extended emergency number list value
// payload, or zero for a nil receiver or nil Value.
func (f *ExtendedEmergencyNumberListFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes an extended emergency number list value without its
// IEI and length octet. It copies the payload into a temporary struct first,
// leaving f unchanged on failure.
func (f *ExtendedEmergencyNumberListFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: extended emergency number list fields are nil", ErrInvalidValue)
	}
	if len(payload) < 4 {
		return fmt.Errorf("%w: extended emergency number list has %d value octets, want >= 4", ErrInvalidLength, len(payload))
	}
	if len(payload) > 65535 {
		return fmt.Errorf("%w: extended emergency number list has %d value octets, want <= 65535", ErrInvalidLength, len(payload))
	}
	var decoded ExtendedEmergencyNumberListFields
	decoded.Value = make([]byte, len(payload))
	copy(decoded.Value, payload)
	*f = decoded
	return nil
}
