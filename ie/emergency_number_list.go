package ie

import (
	"fmt"
	"io"
)

// EmergencyNumberListFields is the structured value of the tagged Emergency number
// list IE (TS 24.301 clause 9.9.3.37). It carries an opaque byte string between
// 3 and 48 octets representing a list of dial digits.
type EmergencyNumberListFields struct {
	Value []byte
}

// NewEmergencyNumberListWithFields creates a tagged Emergency number list IE from
// fields.
func NewEmergencyNumberListWithFields(fields *EmergencyNumberListFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewEmergencyNumberList(payload), nil
}

// NewPositionalEmergencyNumberListWithFields creates a positional Emergency number
// list IE from fields.
func NewPositionalEmergencyNumberListWithFields(fields *EmergencyNumberListFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalEmergencyNumberList(payload), nil
}

// EmergencyNumberList returns the IE as structured fields.
func (i *IE) EmergencyNumberList() (*EmergencyNumberListFields, error) {
	if err := i.requireDefinition(EmergencyNumberListDefinition); err != nil {
		return nil, err
	}
	return ParseEmergencyNumberListFields(i.Payload)
}

// ParseEmergencyNumberListFields decodes an emergency number list value without its
// IEI and length octet.
func ParseEmergencyNumberListFields(payload []byte) (*EmergencyNumberListFields, error) {
	fields := &EmergencyNumberListFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the payload is within [3, 48] octets.
func (f *EmergencyNumberListFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: emergency number list fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 3 {
		return fmt.Errorf("%w: emergency number list has %d value octets, want >= 3", ErrInvalidLength, n)
	}
	if n > 48 {
		return fmt.Errorf("%w: emergency number list has %d value octets, want <= 48", ErrInvalidLength, n)
	}
	return nil
}

// Marshal encodes fields into an emergency number list value without its IEI and
// length octet.
func (f *EmergencyNumberListFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *EmergencyNumberListFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the emergency number list value payload, or zero
// for a nil receiver or nil Value.
func (f *EmergencyNumberListFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes an emergency number list value without its IEI and length
// octet. It copies the payload into a temporary struct first, leaving f unchanged
// on failure.
func (f *EmergencyNumberListFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: emergency number list fields are nil", ErrInvalidValue)
	}
	if len(payload) < 3 {
		return fmt.Errorf("%w: emergency number list has %d value octets, want >= 3", ErrInvalidLength, len(payload))
	}
	if len(payload) > 48 {
		return fmt.Errorf("%w: emergency number list has %d value octets, want <= 48", ErrInvalidLength, len(payload))
	}
	var decoded EmergencyNumberListFields
	decoded.Value = make([]byte, len(payload))
	copy(decoded.Value, payload)
	*f = decoded
	return nil
}
