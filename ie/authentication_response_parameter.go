package ie

import (
	"fmt"
	"io"
	"slices"
)

// AuthenticationResponseParameterFields is the structured value of the
// Authentication response parameter IE (TS 24.301 clause 9.9.3.4). Response is
// the RES value returned by EPS AKA and contains from 4 to 16 octets.
type AuthenticationResponseParameterFields struct {
	Response []byte
}

// NewPositionalAuthenticationResponseParameterWithFields creates a positional
// Authentication response parameter IE from fields.
func NewPositionalAuthenticationResponseParameterWithFields(fields *AuthenticationResponseParameterFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalAuthenticationResponseParameter(payload), nil
}

// AuthenticationResponseParameter returns the IE as structured fields.
func (i *IE) AuthenticationResponseParameter() (*AuthenticationResponseParameterFields, error) {
	if err := i.requireDefinition(AuthenticationResponseParameterDefinition); err != nil {
		return nil, err
	}
	return ParseAuthenticationResponseParameterFields(i.Payload)
}

// ParseAuthenticationResponseParameterFields decodes an Authentication
// response parameter value without its IEI and length octet.
func ParseAuthenticationResponseParameterFields(payload []byte) (*AuthenticationResponseParameterFields, error) {
	fields := &AuthenticationResponseParameterFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies the RES length.
func (f *AuthenticationResponseParameterFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: authentication response parameter fields are nil", ErrInvalidValue)
	}
	if len(f.Response) < 4 || len(f.Response) > 16 {
		return fmt.Errorf("%w: authentication response has %d octets, want 4..16", ErrInvalidLength, len(f.Response))
	}
	return nil
}

// Marshal encodes fields into an Authentication response parameter value.
func (f *AuthenticationResponseParameterFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return slices.Clone(f.Response), nil
}

// MarshalTo encodes fields into b.
func (f *AuthenticationResponseParameterFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Response) {
		return io.ErrShortBuffer
	}
	copy(b, f.Response)
	return nil
}

// MarshalLen returns the size of the Authentication response parameter value
// payload.
func (f *AuthenticationResponseParameterFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return len(f.Response)
}

// UnmarshalBinary decodes an Authentication response parameter value.
func (f *AuthenticationResponseParameterFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: authentication response parameter fields are nil", ErrInvalidValue)
	}
	decoded := AuthenticationResponseParameterFields{Response: slices.Clone(payload)}
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}
