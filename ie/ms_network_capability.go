package ie

import (
	"fmt"
	"io"
)

// MSNetworkCapabilityFields is the structured value of the MS network
// capability IE (TS 24.301 clause 9.9.3.20). The first byte encodes the core
// network type indicator (bit 7) and capability flags such as SRVCC support.
// Remaining bytes carry encryption/integrity algorithm bitmaps and other
// capability flags.
type MSNetworkCapabilityFields struct {
	Value []byte
}

// NewMSNetworkCapabilityWithFields creates a tagged MS network capability IE
// from fields.
func NewMSNetworkCapabilityWithFields(fields *MSNetworkCapabilityFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewMSNetworkCapability(payload), nil
}

// NewPositionalMSNetworkCapabilityWithFields creates a positional MS network
// capability IE from fields.
func NewPositionalMSNetworkCapabilityWithFields(fields *MSNetworkCapabilityFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalMSNetworkCapability(payload), nil
}

// MSNetworkCapability returns the IE as structured fields.
func (i *IE) MSNetworkCapability() (*MSNetworkCapabilityFields, error) {
	if err := i.requireDefinition(MSNetworkCapabilityDefinition); err != nil {
		return nil, err
	}
	return ParseMSNetworkCapabilityFields(i.Payload)
}

// ParseMSNetworkCapabilityFields decodes an MS network capability value without
// its IEI and length octet.
func ParseMSNetworkCapabilityFields(payload []byte) (*MSNetworkCapabilityFields, error) {
	fields := &MSNetworkCapabilityFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the value slice is non-nil and within the permitted
// length range (2–8 octets). Reserved bits in the first byte are checked.
func (f *MSNetworkCapabilityFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: MS network capability fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 2 || n > 8 {
		return fmt.Errorf("%w: MS network capability has %d value octets, want 2..8", ErrInvalidLength, n)
	}
	// First byte: bits 2-6 (0x7c) are spare.
	if f.Value[0]&0x7c != 0 {
		return fmt.Errorf("%w: MS network capability first byte spare bits (2-6) are nonzero", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into an MS network capability value without its IEI and
// length octet.
func (f *MSNetworkCapabilityFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *MSNetworkCapabilityFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the MS network capability value payload. It
// returns zero for a nil receiver or nil Value.
func (f *MSNetworkCapabilityFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes an MS network capability value without its IEI and
// length octet. It leaves f unchanged when payload is invalid.
func (f *MSNetworkCapabilityFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: MS network capability fields are nil", ErrInvalidValue)
	}
	if len(payload) < 2 || len(payload) > 8 {
		return fmt.Errorf("%w: MS network capability has %d value octets, want 2..8", ErrInvalidLength, len(payload))
	}
	if payload[0]&0x7c != 0 {
		return fmt.Errorf("%w: MS network capability first byte spare bits (2-6) are nonzero", ErrInvalidValue)
	}
	decoded := MSNetworkCapabilityFields{
		Value: make([]byte, len(payload)),
	}
	copy(decoded.Value, payload)
	*f = decoded
	return nil
}
