package ie

import (
	"fmt"
	"io"
)

// NetworkResourceIdentifier is the 10-bit part of a TMSI that identifies the
// network resource, as defined by TS 24.008 clause 10.5.5.31.
type NetworkResourceIdentifier uint16

const (
	// NetworkResourceIdentifierMinimum is the lowest valid NRI value.
	NetworkResourceIdentifierMinimum NetworkResourceIdentifier = 0
	// NetworkResourceIdentifierMaximum is the highest valid NRI value.
	NetworkResourceIdentifierMaximum NetworkResourceIdentifier = 0x03ff
)

// NetworkResourceIdentifierContainerFields is the structured value of the
// Network resource identifier container IE (TS 24.301 clause 9.9.3.24A and TS
// 24.008 clause 10.5.5.31). NetworkResourceIdentifier consists of bits 23 to
// 14 of a valid TMSI.
type NetworkResourceIdentifierContainerFields struct {
	NetworkResourceIdentifier NetworkResourceIdentifier
}

// NewNetworkResourceIdentifierContainerWithFields creates a tagged Network
// resource identifier container IE from fields.
func NewNetworkResourceIdentifierContainerWithFields(fields *NetworkResourceIdentifierContainerFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewNetworkResourceIdentifierContainer(payload), nil
}

// NewPositionalNetworkResourceIdentifierContainerWithFields creates a
// positional Network resource identifier container IE from fields.
func NewPositionalNetworkResourceIdentifierContainerWithFields(fields *NetworkResourceIdentifierContainerFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalNetworkResourceIdentifierContainer(payload), nil
}

// NetworkResourceIdentifierContainer returns the IE as structured fields.
func (i *IE) NetworkResourceIdentifierContainer() (*NetworkResourceIdentifierContainerFields, error) {
	if err := i.requireDefinition(NetworkResourceIdentifierContainerDefinition); err != nil {
		return nil, err
	}
	return ParseNetworkResourceIdentifierContainerFields(i.Payload)
}

// ParseNetworkResourceIdentifierContainerFields decodes a Network resource
// identifier container value without its IEI and length octet.
func ParseNetworkResourceIdentifierContainerFields(payload []byte) (*NetworkResourceIdentifierContainerFields, error) {
	fields := &NetworkResourceIdentifierContainerFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that NetworkResourceIdentifier fits ten bits.
func (f *NetworkResourceIdentifierContainerFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: network resource identifier container fields are nil", ErrInvalidValue)
	}
	if f.NetworkResourceIdentifier > NetworkResourceIdentifierMaximum {
		return fmt.Errorf("%w: network resource identifier %d exceeds ten bits", ErrInvalidValue, f.NetworkResourceIdentifier)
	}
	return nil
}

// Marshal encodes fields into a two-octet Network resource identifier
// container value without its IEI and length octet.
func (f *NetworkResourceIdentifierContainerFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.NetworkResourceIdentifier >> 2), byte(f.NetworkResourceIdentifier << 6)}, nil
}

// MarshalTo encodes fields into b.
func (f *NetworkResourceIdentifierContainerFields) MarshalTo(b []byte) error {
	payload, err := f.Marshal()
	if err != nil {
		return err
	}
	if len(b) < len(payload) {
		return io.ErrShortBuffer
	}
	copy(b, payload)
	return nil
}

// MarshalLen returns the size of the Network resource identifier container
// value payload.
func (f *NetworkResourceIdentifierContainerFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 2
}

// UnmarshalBinary decodes a two-octet Network resource identifier container
// value without its IEI and length octet.
func (f *NetworkResourceIdentifierContainerFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: network resource identifier container fields are nil", ErrInvalidValue)
	}
	if len(payload) != 2 {
		return fmt.Errorf("%w: network resource identifier container has %d value octets, want 2", ErrInvalidLength, len(payload))
	}
	if payload[1]&0x3f != 0 {
		return fmt.Errorf("%w: network resource identifier container spare bits are nonzero", ErrInvalidValue)
	}
	*f = NetworkResourceIdentifierContainerFields{
		NetworkResourceIdentifier: NetworkResourceIdentifier(uint16(payload[0])<<2 | uint16(payload[1]>>6)),
	}
	return nil
}
