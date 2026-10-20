package ie

import (
	"fmt"
	"io"
)

// NASKeySetIdentifierFields is the structured value of the NAS key set
// identifier IE (TS 24.301 clause 9.9.3.21). KeySetIdentifier 7 means that no
// key is available in the UE-to-network direction. MappedSecurityContext is
// ignored by the protocol for that value, but retained to preserve the wire
// representation.
type NASKeySetIdentifierFields struct {
	KeySetIdentifier      uint8
	MappedSecurityContext bool
}

// NewNASKeySetIdentifierWithFields creates a tagged NAS key set identifier IE
// from fields.
func NewNASKeySetIdentifierWithFields(fields *NASKeySetIdentifierFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewNASKeySetIdentifier(payload[0]), nil
}

// NewPositionalNASKeySetIdentifierWithFields creates a positional NAS key set
// identifier IE from fields.
func NewPositionalNASKeySetIdentifierWithFields(fields *NASKeySetIdentifierFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalNASKeySetIdentifier(payload[0]), nil
}

// NASKeySetIdentifier returns the IE as structured fields.
func (i *IE) NASKeySetIdentifier() (*NASKeySetIdentifierFields, error) {
	if err := i.requireDefinition(NASKeySetIdentifierDefinition); err != nil {
		return nil, err
	}
	return ParseNASKeySetIdentifierFields(i.Payload)
}

// ParseNASKeySetIdentifierFields decodes a NAS key set identifier value
// without its IEI.
func ParseNASKeySetIdentifierFields(payload []byte) (*NASKeySetIdentifierFields, error) {
	fields := &NASKeySetIdentifierFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that KeySetIdentifier fits three bits.
func (f *NASKeySetIdentifierFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: NAS key set identifier fields are nil", ErrInvalidValue)
	}
	if f.KeySetIdentifier > 0x07 {
		return fmt.Errorf("%w: NAS key set identifier exceeds its encoded width", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into the low nibble of a NAS key set identifier
// value, without its IEI.
func (f *NASKeySetIdentifierFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	encoded := f.KeySetIdentifier
	if f.MappedSecurityContext {
		encoded |= 0x08
	}
	return []byte{encoded}, nil
}

// MarshalTo encodes fields into b.
func (f *NASKeySetIdentifierFields) MarshalTo(b []byte) error {
	payload, err := f.Marshal()
	if err != nil {
		return err
	}
	if len(b) < len(payload) {
		return io.ErrShortBuffer
	}
	b[0] = payload[0]
	return nil
}

// MarshalLen returns the size of the NAS key set identifier value payload.
func (f *NASKeySetIdentifierFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a NAS key set identifier value without its IEI.
func (f *NASKeySetIdentifierFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: NAS key set identifier fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: NAS key set identifier has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xf0 != 0 {
		return fmt.Errorf("%w: NAS key set identifier spare bits are nonzero", ErrInvalidValue)
	}
	*f = NASKeySetIdentifierFields{
		KeySetIdentifier:      payload[0] & 0x07,
		MappedSecurityContext: payload[0]&0x08 != 0,
	}
	return nil
}
