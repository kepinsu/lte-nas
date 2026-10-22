package ie

import (
	"fmt"
	"io"
)

// AuthenticationParameterRANDFields is the structured value of the
// Authentication parameter RAND IE (TS 24.301 clause 9.9.3.3). RandomChallenge
// is the 128-bit RAND value used by EPS AKA.
type AuthenticationParameterRANDFields struct {
	RandomChallenge [16]byte
}

// NewPositionalAuthenticationParameterRANDWithFields creates a positional
// Authentication parameter RAND IE from fields.
func NewPositionalAuthenticationParameterRANDWithFields(fields *AuthenticationParameterRANDFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalAuthenticationParameterRAND(payload), nil
}

// AuthenticationParameterRAND returns the IE as structured fields.
func (i *IE) AuthenticationParameterRAND() (*AuthenticationParameterRANDFields, error) {
	if err := i.requireDefinition(AuthenticationParameterRANDDefinition); err != nil {
		return nil, err
	}
	return ParseAuthenticationParameterRANDFields(i.Payload)
}

// ParseAuthenticationParameterRANDFields decodes an Authentication parameter
// RAND value without its IEI.
func ParseAuthenticationParameterRANDFields(payload []byte) (*AuthenticationParameterRANDFields, error) {
	fields := &AuthenticationParameterRANDFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *AuthenticationParameterRANDFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: authentication parameter RAND fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a 16-octet Authentication parameter RAND value.
func (f *AuthenticationParameterRANDFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.RandomChallenge))
	copy(payload, f.RandomChallenge[:])
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *AuthenticationParameterRANDFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < f.MarshalLen() {
		return io.ErrShortBuffer
	}
	copy(b, f.RandomChallenge[:])
	return nil
}

// MarshalLen returns the size of the Authentication parameter RAND value
// payload.
func (f *AuthenticationParameterRANDFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 16
}

// UnmarshalBinary decodes a 16-octet Authentication parameter RAND value.
func (f *AuthenticationParameterRANDFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: authentication parameter RAND fields are nil", ErrInvalidValue)
	}
	if len(payload) != f.MarshalLen() {
		return fmt.Errorf("%w: authentication parameter RAND has %d value octets, want 16", ErrInvalidLength, len(payload))
	}
	var decoded AuthenticationParameterRANDFields
	copy(decoded.RandomChallenge[:], payload)
	*f = decoded
	return nil
}
