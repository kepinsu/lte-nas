package ie

import (
	"fmt"
	"io"
)

// AuthenticationFailureParameterFields is the structured value of the
// Authentication failure parameter IE (TS 24.301 clause 9.9.3.1). It carries
// the AUTS synchronization-failure token defined by 3GPP TS 33.102.
type AuthenticationFailureParameterFields struct {
	SequenceNumberXORAnonymityKey [6]byte
	MessageAuthenticationCode     [8]byte
}

// NewAuthenticationFailureParameterWithFields creates a tagged Authentication
// failure parameter IE from fields.
func NewAuthenticationFailureParameterWithFields(fields *AuthenticationFailureParameterFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewAuthenticationFailureParameter(payload), nil
}

// NewPositionalAuthenticationFailureParameterWithFields creates a positional
// Authentication failure parameter IE from fields.
func NewPositionalAuthenticationFailureParameterWithFields(fields *AuthenticationFailureParameterFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalAuthenticationFailureParameter(payload), nil
}

// AuthenticationFailureParameter returns the IE as structured fields.
func (i *IE) AuthenticationFailureParameter() (*AuthenticationFailureParameterFields, error) {
	if err := i.requireDefinition(AuthenticationFailureParameterDefinition); err != nil {
		return nil, err
	}
	return ParseAuthenticationFailureParameterFields(i.Payload)
}

// ParseAuthenticationFailureParameterFields decodes an Authentication failure
// parameter value without its IEI and length octet.
func ParseAuthenticationFailureParameterFields(payload []byte) (*AuthenticationFailureParameterFields, error) {
	fields := &AuthenticationFailureParameterFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *AuthenticationFailureParameterFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: authentication failure parameter fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a 14-octet Authentication failure parameter
// value.
func (f *AuthenticationFailureParameterFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, f.MarshalLen())
	copy(payload[:6], f.SequenceNumberXORAnonymityKey[:])
	copy(payload[6:], f.MessageAuthenticationCode[:])
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *AuthenticationFailureParameterFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < f.MarshalLen() {
		return io.ErrShortBuffer
	}
	copy(b[:6], f.SequenceNumberXORAnonymityKey[:])
	copy(b[6:14], f.MessageAuthenticationCode[:])
	return nil
}

// MarshalLen returns the size of the Authentication failure parameter value
// payload.
func (f *AuthenticationFailureParameterFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 14
}

// UnmarshalBinary decodes a 14-octet Authentication failure parameter value.
func (f *AuthenticationFailureParameterFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: authentication failure parameter fields are nil", ErrInvalidValue)
	}
	if len(payload) != f.MarshalLen() {
		return fmt.Errorf("%w: authentication failure parameter has %d value octets, want 14", ErrInvalidLength, len(payload))
	}
	var decoded AuthenticationFailureParameterFields
	copy(decoded.SequenceNumberXORAnonymityKey[:], payload[:6])
	copy(decoded.MessageAuthenticationCode[:], payload[6:])
	*f = decoded
	return nil
}
