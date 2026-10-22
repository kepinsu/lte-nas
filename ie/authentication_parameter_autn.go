package ie

import (
	"fmt"
	"io"
)

// AuthenticationParameterAUTNFields is the structured value of the
// Authentication parameter AUTN IE (TS 24.301 clause 9.9.3.2). The component
// fields use the EPS AKA AUTN layout defined by 3GPP TS 33.102.
type AuthenticationParameterAUTNFields struct {
	SequenceNumberXORAnonymityKey [6]byte
	AuthenticationManagementField [2]byte
	MessageAuthenticationCode     [8]byte
}

// NewPositionalAuthenticationParameterAUTNWithFields creates a positional
// Authentication parameter AUTN IE from fields.
func NewPositionalAuthenticationParameterAUTNWithFields(fields *AuthenticationParameterAUTNFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalAuthenticationParameterAUTN(payload), nil
}

// AuthenticationParameterAUTN returns the IE as structured fields.
func (i *IE) AuthenticationParameterAUTN() (*AuthenticationParameterAUTNFields, error) {
	if err := i.requireDefinition(AuthenticationParameterAUTNDefinition); err != nil {
		return nil, err
	}
	return ParseAuthenticationParameterAUTNFields(i.Payload)
}

// ParseAuthenticationParameterAUTNFields decodes an Authentication parameter
// AUTN value without its IEI and length octet.
func ParseAuthenticationParameterAUTNFields(payload []byte) (*AuthenticationParameterAUTNFields, error) {
	fields := &AuthenticationParameterAUTNFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *AuthenticationParameterAUTNFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: authentication parameter AUTN fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a 16-octet Authentication parameter AUTN value.
func (f *AuthenticationParameterAUTNFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, f.MarshalLen())
	copy(payload[:6], f.SequenceNumberXORAnonymityKey[:])
	copy(payload[6:8], f.AuthenticationManagementField[:])
	copy(payload[8:], f.MessageAuthenticationCode[:])
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *AuthenticationParameterAUTNFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < f.MarshalLen() {
		return io.ErrShortBuffer
	}
	copy(b[:6], f.SequenceNumberXORAnonymityKey[:])
	copy(b[6:8], f.AuthenticationManagementField[:])
	copy(b[8:16], f.MessageAuthenticationCode[:])
	return nil
}

// MarshalLen returns the size of the Authentication parameter AUTN value
// payload.
func (f *AuthenticationParameterAUTNFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 16
}

// UnmarshalBinary decodes a 16-octet Authentication parameter AUTN value.
func (f *AuthenticationParameterAUTNFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: authentication parameter AUTN fields are nil", ErrInvalidValue)
	}
	if len(payload) != f.MarshalLen() {
		return fmt.Errorf("%w: authentication parameter AUTN has %d value octets, want 16", ErrInvalidLength, len(payload))
	}
	var decoded AuthenticationParameterAUTNFields
	copy(decoded.SequenceNumberXORAnonymityKey[:], payload[:6])
	copy(decoded.AuthenticationManagementField[:], payload[6:8])
	copy(decoded.MessageAuthenticationCode[:], payload[8:])
	*f = decoded
	return nil
}
