package ie

import (
	"encoding/binary"
	"fmt"
	"io"
)

// NonceFields is the structured value of the Nonce IE (TS 24.301 clause
// 9.9.3.25).
type NonceFields struct {
	Nonce uint32
}

// NewNonceUEWithFields creates a tagged UE nonce IE from fields.
func NewNonceUEWithFields(fields *NonceFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewNonceUE(payload), nil
}

// NewNonceMMEWithFields creates a tagged MME nonce IE from fields.
func NewNonceMMEWithFields(fields *NonceFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewNonceMME(payload), nil
}

// NewPositionalNonceWithFields creates a positional Nonce IE from fields.
func NewPositionalNonceWithFields(fields *NonceFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalNonce(payload), nil
}

// Nonce returns the IE as structured fields.
func (i *IE) Nonce() (*NonceFields, error) {
	if err := i.requireDefinition(NonceDefinition); err != nil {
		return nil, err
	}
	return ParseNonceFields(i.Payload)
}

// ParseNonceFields decodes a Nonce value without its IEI.
func ParseNonceFields(payload []byte) (*NonceFields, error) {
	fields := &NonceFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *NonceFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: nonce fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a four-octet Nonce value.
func (f *NonceFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, f.Nonce)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *NonceFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 4 {
		return io.ErrShortBuffer
	}
	binary.BigEndian.PutUint32(b, f.Nonce)
	return nil
}

// MarshalLen returns the size of the Nonce value payload.
func (f *NonceFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 4
}

// UnmarshalBinary decodes a four-octet Nonce value.
func (f *NonceFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: nonce fields are nil", ErrInvalidValue)
	}
	if len(payload) != 4 {
		return fmt.Errorf("%w: nonce has %d value octets, want 4", ErrInvalidLength, len(payload))
	}
	*f = NonceFields{Nonce: binary.BigEndian.Uint32(payload)}
	return nil
}
