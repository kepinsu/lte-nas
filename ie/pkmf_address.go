package ie

import (
	"fmt"
	"io"
)

// PKMFAddressFields is the structured value of the PKMF address IE (TS 24.301
// clause 9.9.4.21). It carries an IPv4 (4 bytes) or IPv6 (16 bytes) address
// as an opaque byte string between 1 and 17 octets.
type PKMFAddressFields struct {
	Value []byte
}

// NewPKMFAddressWithFields creates a tagged PKMF address IE from fields.
func NewPKMFAddressWithFields(fields *PKMFAddressFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPKMFAddress(payload), nil
}

// NewPositionalPKMFAddressWithFields creates a positional PKMF address IE from
// fields.
func NewPositionalPKMFAddressWithFields(fields *PKMFAddressFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalPKMFAddress(payload), nil
}

// PKMFAddress returns the IE as structured fields.
func (i *IE) PKMFAddress() (*PKMFAddressFields, error) {
	if err := i.requireDefinition(PKMFAddressDefinition); err != nil {
		return nil, err
	}
	return ParsePKMFAddressFields(i.Payload)
}

// ParsePKMFAddressFields decodes a PKMF address value without its IEI and
// length octet.
func ParsePKMFAddressFields(payload []byte) (*PKMFAddressFields, error) {
	fields := &PKMFAddressFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the payload is within [1, 17] octets.
func (f *PKMFAddressFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: PKMF address fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 1 {
		return fmt.Errorf("%w: PKMF address has %d value octets, want >= 1", ErrInvalidLength, n)
	}
	if n > 17 {
		return fmt.Errorf("%w: PKMF address has %d value octets, want <= 17", ErrInvalidLength, n)
	}
	return nil
}

// Marshal encodes fields into a PKMF address value without its IEI and length
// octet.
func (f *PKMFAddressFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *PKMFAddressFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the PKMF address value payload, or zero for a
// nil receiver or nil Value.
func (f *PKMFAddressFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes a PKMF address value without its IEI and length
// octet. It copies the payload into a temporary struct first, leaving f
// unchanged on failure.
func (f *PKMFAddressFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: PKMF address fields are nil", ErrInvalidValue)
	}
	if len(payload) < 1 {
		return fmt.Errorf("%w: PKMF address has %d value octets, want >= 1", ErrInvalidLength, len(payload))
	}
	if len(payload) > 17 {
		return fmt.Errorf("%w: PKMF address has %d value octets, want <= 17", ErrInvalidLength, len(payload))
	}
	var decoded PKMFAddressFields
	decoded.Value = make([]byte, len(payload))
	copy(decoded.Value, payload)
	*f = decoded
	return nil
}
