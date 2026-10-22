package ie

import (
	"fmt"
	"io"
)

// PagingRestrictionFields is the structured value of the paging restriction IE
// (TS 24.301 clause 9.9.3.66). It carries an opaque byte string between 1 and
// 3 octets encoding bearer capability, extended RB type indicator, and/or max
// MRU value.
type PagingRestrictionFields struct {
	Value []byte
}

// NewPagingRestrictionWithFields creates a tagged Paging restriction IE from
// fields.
func NewPagingRestrictionWithFields(fields *PagingRestrictionFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPagingRestriction(payload), nil
}

// NewPositionalPagingRestrictionWithFields creates a positional Paging
// restriction IE from fields.
func NewPositionalPagingRestrictionWithFields(fields *PagingRestrictionFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalPagingRestriction(payload), nil
}

// PagingRestriction returns the IE as structured fields.
func (i *IE) PagingRestriction() (*PagingRestrictionFields, error) {
	if err := i.requireDefinition(PagingRestrictionDefinition); err != nil {
		return nil, err
	}
	return ParsePagingRestrictionFields(i.Payload)
}

// ParsePagingRestrictionFields decodes a paging restriction value without its
// IEI and length octet.
func ParsePagingRestrictionFields(payload []byte) (*PagingRestrictionFields, error) {
	fields := &PagingRestrictionFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the payload is within [1, 3] octets.
func (f *PagingRestrictionFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: paging restriction fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 1 {
		return fmt.Errorf("%w: paging restriction has %d value octets, want >= 1", ErrInvalidLength, n)
	}
	if n > 3 {
		return fmt.Errorf("%w: paging restriction has %d value octets, want <= 3", ErrInvalidLength, n)
	}
	return nil
}

// Marshal encodes fields into a paging restriction value without its IEI and
// length octet.
func (f *PagingRestrictionFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *PagingRestrictionFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the paging restriction value payload, or zero
// for a nil receiver or nil Value.
func (f *PagingRestrictionFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes a paging restriction value without its IEI and length
// octet. It copies the payload into a temporary struct first, leaving f
// unchanged on failure.
func (f *PagingRestrictionFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: paging restriction fields are nil", ErrInvalidValue)
	}
	if len(payload) < 1 {
		return fmt.Errorf("%w: paging restriction has %d value octets, want >= 1", ErrInvalidLength, len(payload))
	}
	if len(payload) > 3 {
		return fmt.Errorf("%w: paging restriction has %d value octets, want <= 3", ErrInvalidLength, len(payload))
	}
	var decoded PagingRestrictionFields
	decoded.Value = make([]byte, len(payload))
	copy(decoded.Value, payload)
	*f = decoded
	return nil
}
