package ie

import (
	"encoding/binary"
	"fmt"
	"io"
)

// ShortMACFields is the structured value of the Short MAC IE (TS 24.301
// clause 9.9.3.28).
type ShortMACFields struct {
	MAC uint16
}

// NewPositionalShortMACWithFields creates a positional Short MAC IE from
// fields.
func NewPositionalShortMACWithFields(fields *ShortMACFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalShortMAC(payload), nil
}

// ShortMAC returns the IE as structured fields.
func (i *IE) ShortMAC() (*ShortMACFields, error) {
	if err := i.requireDefinition(ShortMACDefinition); err != nil {
		return nil, err
	}
	return ParseShortMACFields(i.Payload)
}

// ParseShortMACFields decodes a Short MAC value without its IEI.
func ParseShortMACFields(payload []byte) (*ShortMACFields, error) {
	fields := &ShortMACFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *ShortMACFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: short MAC fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a two-octet Short MAC value.
func (f *ShortMACFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, 2)
	binary.BigEndian.PutUint16(payload, f.MAC)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *ShortMACFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 2 {
		return io.ErrShortBuffer
	}
	binary.BigEndian.PutUint16(b, f.MAC)
	return nil
}

// MarshalLen returns the size of the Short MAC value payload.
func (f *ShortMACFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 2
}

// UnmarshalBinary decodes a two-octet Short MAC value.
func (f *ShortMACFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: short MAC fields are nil", ErrInvalidValue)
	}
	if len(payload) != 2 {
		return fmt.Errorf("%w: short MAC has %d value octets, want 2", ErrInvalidLength, len(payload))
	}
	*f = ShortMACFields{MAC: binary.BigEndian.Uint16(payload)}
	return nil
}
