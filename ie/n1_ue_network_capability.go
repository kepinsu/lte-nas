package ie

import (
	"fmt"
	"io"
)

// N1UENetworkCapabilityFields is the structured value of the N1 UE network
// capability IE (TS 24.301 clause 9.9.3.57). The first byte encodes the N1
// signalling negotiation information (bit 7 flag and security algorithms).
// Remaining bytes carry capability bitmaps.
type N1UENetworkCapabilityFields struct {
	Value []byte
}

// NewN1UENetworkCapabilityWithFields creates a tagged N1 UE network capability
// IE from fields.
func NewN1UENetworkCapabilityWithFields(fields *N1UENetworkCapabilityFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewN1UENetworkCapability(payload), nil
}

// NewPositionalN1UENetworkCapabilityWithFields creates a positional N1 UE
// network capability IE from fields.
func NewPositionalN1UENetworkCapabilityWithFields(fields *N1UENetworkCapabilityFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalN1UENetworkCapability(payload), nil
}

// N1UENetworkCapability returns the IE as structured fields.
func (i *IE) N1UENetworkCapability() (*N1UENetworkCapabilityFields, error) {
	if err := i.requireDefinition(N1UENetworkCapabilityDefinition); err != nil {
		return nil, err
	}
	return ParseN1UENetworkCapabilityFields(i.Payload)
}

// ParseN1UENetworkCapabilityFields decodes an N1 UE network capability value
// without its IEI and length octet.
func ParseN1UENetworkCapabilityFields(payload []byte) (*N1UENetworkCapabilityFields, error) {
	fields := &N1UENetworkCapabilityFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the value slice is non-nil and within the permitted
// length range (1–13 octets). Reserved bits in the first byte are checked.
func (f *N1UENetworkCapabilityFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: N1 UE network capability fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 1 || n > 13 {
		return fmt.Errorf("%w: N1 UE network capability has %d value octets, want 1..13", ErrInvalidLength, n)
	}
	// First byte: bits 2-4 (0x1c) are spare.
	if f.Value[0]&0x1c != 0 {
		return fmt.Errorf("%w: N1 UE network capability first byte spare bits (2-4) are nonzero", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into an N1 UE network capability value without its IEI
// and length octet.
func (f *N1UENetworkCapabilityFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *N1UENetworkCapabilityFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the N1 UE network capability value payload. It
// returns zero for a nil receiver or nil Value.
func (f *N1UENetworkCapabilityFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes an N1 UE network capability value without its IEI and
// length octet. It leaves f unchanged when payload is invalid.
func (f *N1UENetworkCapabilityFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: N1 UE network capability fields are nil", ErrInvalidValue)
	}
	if len(payload) < 1 || len(payload) > 13 {
		return fmt.Errorf("%w: N1 UE network capability has %d value octets, want 1..13", ErrInvalidLength, len(payload))
	}
	if payload[0]&0x1c != 0 {
		return fmt.Errorf("%w: N1 UE network capability first byte spare bits (2-4) are nonzero", ErrInvalidValue)
	}
	decoded := N1UENetworkCapabilityFields{
		Value: make([]byte, len(payload)),
	}
	copy(decoded.Value, payload)
	*f = decoded
	return nil
}
