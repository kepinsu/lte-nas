package ie

import (
	"fmt"
	"io"
)

// PagingIdentity identifies the identity type used to page for non-EPS
// services.
type PagingIdentity uint8

// Paging identity values defined by TS 24.301 clause 9.9.3.25A.
const (
	PagingIdentityIMSI PagingIdentity = iota
	PagingIdentityTMSI
)

// PagingIdentityFields is the structured value of the Paging identity IE (TS
// 24.301 clause 9.9.3.25A).
type PagingIdentityFields struct {
	Identity PagingIdentity
}

// NewPositionalPagingIdentityWithFields creates a positional Paging identity
// IE from fields.
func NewPositionalPagingIdentityWithFields(fields *PagingIdentityFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalPagingIdentity(payload[0]), nil
}

// PagingIdentity returns the IE as structured fields.
func (i *IE) PagingIdentity() (*PagingIdentityFields, error) {
	if err := i.requireDefinition(PagingIdentityDefinition); err != nil {
		return nil, err
	}
	return ParsePagingIdentityFields(i.Payload)
}

// ParsePagingIdentityFields decodes a Paging identity value without its IEI.
func ParsePagingIdentityFields(payload []byte) (*PagingIdentityFields, error) {
	fields := &PagingIdentityFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Identity is a defined paging identity.
func (f *PagingIdentityFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: paging identity fields are nil", ErrInvalidValue)
	}
	if f.Identity != PagingIdentityIMSI && f.Identity != PagingIdentityTMSI {
		return fmt.Errorf("%w: paging identity %d is reserved", ErrInvalidValue, f.Identity)
	}
	return nil
}

// Marshal encodes fields into a Paging identity value without its IEI.
func (f *PagingIdentityFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.Identity)}, nil
}

// MarshalTo encodes fields into b.
func (f *PagingIdentityFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the Paging identity value payload.
func (f *PagingIdentityFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a Paging identity value without its IEI.
func (f *PagingIdentityFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: paging identity fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: paging identity has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xfe != 0 {
		return fmt.Errorf("%w: paging identity spare bits are nonzero", ErrInvalidValue)
	}
	*f = PagingIdentityFields{Identity: PagingIdentity(payload[0])}
	return nil
}
