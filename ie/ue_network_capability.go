package ie

import (
	"fmt"
	"io"
)

// UENetworkCapabilityFields is the structured value of the UE network
// capability IE (TS 24.301 clause 9.9.3.34). The first byte encodes the NAS
// security algorithm selection (bits 1-4) and core network type indicator
// (bit 7). Remaining bytes carry capability bitmaps.
type UENetworkCapabilityFields struct {
	Value []byte
}

// NewUENetworkCapabilityWithFields creates a tagged UE network capability IE
// from fields.
func NewUENetworkCapabilityWithFields(fields *UENetworkCapabilityFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewUENetworkCapability(payload), nil
}

// NewPositionalUENetworkCapabilityWithFields creates a positional UE network
// capability IE from fields.
func NewPositionalUENetworkCapabilityWithFields(fields *UENetworkCapabilityFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalUENetworkCapability(payload), nil
}

// UENetworkCapability returns the IE as structured fields.
func (i *IE) UENetworkCapability() (*UENetworkCapabilityFields, error) {
	if err := i.requireDefinition(UENetworkCapabilityDefinition); err != nil {
		return nil, err
	}
	return ParseUENetworkCapabilityFields(i.Payload)
}

// ParseUENetworkCapabilityFields decodes a UE network capability value without
// its IEI and length octet.
func ParseUENetworkCapabilityFields(payload []byte) (*UENetworkCapabilityFields, error) {
	fields := &UENetworkCapabilityFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the value slice is non-nil and within the permitted
// length range (2–13 octets). Reserved bits in the first byte are checked.
func (f *UENetworkCapabilityFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: UE network capability fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 2 || n > 13 {
		return fmt.Errorf("%w: UE network capability has %d value octets, want 2..13", ErrInvalidLength, n)
	}
	// First byte: bits 5-6 (value 0x20 and 0x40) are spare.
	if f.Value[0]&0x60 != 0 {
		return fmt.Errorf("%w: UE network capability first byte spare bits are nonzero", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a UE network capability value without its IEI and
// length octet.
func (f *UENetworkCapabilityFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *UENetworkCapabilityFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the UE network capability value payload. It
// returns zero for a nil receiver or nil Value.
func (f *UENetworkCapabilityFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes a UE network capability value without its IEI and
// length octet. It leaves f unchanged when payload is invalid.
func (f *UENetworkCapabilityFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: UE network capability fields are nil", ErrInvalidValue)
	}
	if len(payload) < 2 || len(payload) > 13 {
		return fmt.Errorf("%w: UE network capability has %d value octets, want 2..13", ErrInvalidLength, len(payload))
	}
	if payload[0]&0x60 != 0 {
		return fmt.Errorf("%w: UE network capability first byte spare bits (5-6) are nonzero", ErrInvalidValue)
	}
	decoded := UENetworkCapabilityFields{
		Value: make([]byte, len(payload)),
	}
	copy(decoded.Value, payload)
	*f = decoded
	return nil
}
