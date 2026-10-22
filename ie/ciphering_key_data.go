package ie

import (
	"fmt"
	"io"
)

// CipheringKeyDataFields is the structured value of the ciphering key data IE
// (TS 24.301 clause 9.9.3.56). It carries an opaque byte string between 32 and
// 2288 octets.
type CipheringKeyDataFields struct {
	Value []byte
}

// NewPositionalCipheringKeyDataWithFields creates a positional Ciphering key data
// IE from fields.
func NewPositionalCipheringKeyDataWithFields(fields *CipheringKeyDataFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalCipheringKeyData(payload), nil
}

// CipheringKeyData returns the IE as structured fields.
func (i *IE) CipheringKeyData() (*CipheringKeyDataFields, error) {
	if err := i.requireDefinition(CipheringKeyDataDefinition); err != nil {
		return nil, err
	}
	return ParseCipheringKeyDataFields(i.Payload)
}

// ParseCipheringKeyDataFields decodes a ciphering key data value without its IEI
// and length octet.
func ParseCipheringKeyDataFields(payload []byte) (*CipheringKeyDataFields, error) {
	fields := &CipheringKeyDataFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the payload is within [32, 2288] octets.
func (f *CipheringKeyDataFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: ciphering key data fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 32 {
		return fmt.Errorf("%w: ciphering key data has %d value octets, want >= 32", ErrInvalidLength, n)
	}
	if n > 2288 {
		return fmt.Errorf("%w: ciphering key data has %d value octets, want <= 2288", ErrInvalidLength, n)
	}
	return nil
}

// Marshal encodes fields into a ciphering key data value without its IEI and length
// octet.
func (f *CipheringKeyDataFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *CipheringKeyDataFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the ciphering key data value payload, or zero for
// a nil receiver or nil Value.
func (f *CipheringKeyDataFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes a ciphering key data value without its IEI and length octet.
// It copies the payload into a temporary struct first, leaving f unchanged on failure.
func (f *CipheringKeyDataFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: ciphering key data fields are nil", ErrInvalidValue)
	}
	if len(payload) < 32 {
		return fmt.Errorf("%w: ciphering key data has %d value octets, want >= 32", ErrInvalidLength, len(payload))
	}
	if len(payload) > 2288 {
		return fmt.Errorf("%w: ciphering key data has %d value octets, want <= 2288", ErrInvalidLength, len(payload))
	}
	var decoded CipheringKeyDataFields
	decoded.Value = make([]byte, len(payload))
	copy(decoded.Value, payload)
	*f = decoded
	return nil
}
