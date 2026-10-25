package ie

import (
	"fmt"
	"io"
)

// HashMMEFields is the structured value of the Hash MME IE (TS 24.301 clause
// 9.9.3.50). The eight-octet payload contains the hash value used to identify
// an MME.
type HashMMEFields struct {
	// Payload contains the 8-octet Hash MME value.
	Payload [8]byte
}

// NewHashMMEWithFields creates a tagged Hash MME IE from fields.
func NewHashMMEWithFields(fields *HashMMEFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewHashMME(payload), nil
}

// NewPositionalHashMMEWithFields creates a positional Hash MME IE from fields.
func NewPositionalHashMMEWithFields(fields *HashMMEFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalHashMME(payload), nil
}

// HashMME returns the IE as structured fields.
func (i *IE) HashMME() (*HashMMEFields, error) {
	if err := i.requireDefinition(HashMMEDefinition); err != nil {
		return nil, err
	}
	return ParseHashMMEFields(i.Payload)
}

// ParseHashMMEFields decodes a Hash MME value without its IEI and length octet.
func ParseHashMMEFields(payload []byte) (*HashMMEFields, error) {
	fields := &HashMMEFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *HashMMEFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: Hash MME fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into an 8-octet Hash MME value.
func (f *HashMMEFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Payload))
	copy(payload, f.Payload[:])
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *HashMMEFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 8 {
		return io.ErrShortBuffer
	}
	copy(b, f.Payload[:])
	return nil
}

// MarshalLen returns the size of the Hash MME value payload.
func (f *HashMMEFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 8
}

// UnmarshalBinary decodes an 8-octet Hash MME value.
func (f *HashMMEFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: Hash MME fields are nil", ErrInvalidValue)
	}
	if len(payload) != 8 {
		return fmt.Errorf("%w: Hash MME has %d value octets, want 8", ErrInvalidLength, len(payload))
	}
	var decoded HashMMEFields
	copy(decoded.Payload[:], payload)
	*f = decoded
	return nil
}
