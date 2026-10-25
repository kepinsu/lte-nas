package ie

import (
	"fmt"
	"io"
)

// Non3GPPNWProvidedPoliciesFields is the structured value of the Non-3GPP NW
// provided policies IE (TS 24.301 clause 9.9.3.49, see TS 24.008 clause 10.5.5.37).
type Non3GPPNWProvidedPoliciesFields struct {
	// PCCRule indicates the presence of PCC rules.
	PCCRule bool
}

// NewNon3GPPNWProvidedPoliciesWithFields creates a tagged Non-3GPP NW provided
// policies IE from fields.
func NewNon3GPPNWProvidedPoliciesWithFields(fields *Non3GPPNWProvidedPoliciesFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewNon3GPPNWProvidedPolicies(payload[0]), nil
}

// Non3GPPNWProvidedPolicies returns the IE as structured fields.
func (i *IE) Non3GPPNWProvidedPolicies() (*Non3GPPNWProvidedPoliciesFields, error) {
	if err := i.requireDefinition(Non3GPPNWProvidedPoliciesDefinition); err != nil {
		return nil, err
	}
	return ParseNon3GPPNWProvidedPoliciesFields(i.Payload)
}

// ParseNon3GPPNWProvidedPoliciesFields decodes a Non-3GPP NW provided policies
// value without its IEI.
func ParseNon3GPPNWProvidedPoliciesFields(payload []byte) (*Non3GPPNWProvidedPoliciesFields, error) {
	fields := &Non3GPPNWProvidedPoliciesFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *Non3GPPNWProvidedPoliciesFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: non-3GPP NW provided policies fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a Non-3GPP NW provided policies value, without its
// IEI.
func (f *Non3GPPNWProvidedPoliciesFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	var encoded byte
	if f.PCCRule {
		encoded |= 0x80
	}
	return []byte{encoded}, nil
}

// MarshalTo encodes fields into b.
func (f *Non3GPPNWProvidedPoliciesFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the Non-3GPP NW provided policies value payload.
func (f *Non3GPPNWProvidedPoliciesFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a Non-3GPP NW provided policies value without its IEI.
func (f *Non3GPPNWProvidedPoliciesFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: non-3GPP NW provided policies fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: non-3GPP NW provided policies has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0x7f != 0 {
		return fmt.Errorf("%w: non-3GPP NW provided policies spare bits are nonzero", ErrInvalidValue)
	}
	*f = Non3GPPNWProvidedPoliciesFields{
		PCCRule: payload[0]&0x80 != 0,
	}
	return nil
}
