package ie

import (
	"fmt"
	"io"
)

// CipheringKeySequenceNumberFields is the structured value of the Ciphering
// key sequence number IE (TS 24.301 clause 9.9.3.4A and TS 24.008 clause
// 10.5.1.2). KeySequenceNumber 7 means that no key is available in the
// mobile-to-network direction.
type CipheringKeySequenceNumberFields struct {
	KeySequenceNumber uint8
}

// NewCipheringKeySequenceNumberWithFields creates a tagged Ciphering key
// sequence number IE from fields.
func NewCipheringKeySequenceNumberWithFields(fields *CipheringKeySequenceNumberFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewCipheringKeySequenceNumber(payload[0]), nil
}

// NewPositionalCipheringKeySequenceNumberWithFields creates a positional
// Ciphering key sequence number IE from fields.
func NewPositionalCipheringKeySequenceNumberWithFields(fields *CipheringKeySequenceNumberFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalCipheringKeySequenceNumber(payload[0]), nil
}

// CipheringKeySequenceNumber returns the IE as structured fields.
func (i *IE) CipheringKeySequenceNumber() (*CipheringKeySequenceNumberFields, error) {
	if err := i.requireDefinition(CipheringKeySequenceNumberDefinition); err != nil {
		return nil, err
	}
	return ParseCipheringKeySequenceNumberFields(i.Payload)
}

// ParseCipheringKeySequenceNumberFields decodes a Ciphering key sequence
// number value without its IEI.
func ParseCipheringKeySequenceNumberFields(payload []byte) (*CipheringKeySequenceNumberFields, error) {
	fields := &CipheringKeySequenceNumberFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that KeySequenceNumber fits three bits.
func (f *CipheringKeySequenceNumberFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: ciphering key sequence number fields are nil", ErrInvalidValue)
	}
	if f.KeySequenceNumber > 0x07 {
		return fmt.Errorf("%w: ciphering key sequence number exceeds its encoded width", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into the low nibble of a Ciphering key sequence
// number value, without its IEI.
func (f *CipheringKeySequenceNumberFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{f.KeySequenceNumber}, nil
}

// MarshalTo encodes fields into b.
func (f *CipheringKeySequenceNumberFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the Ciphering key sequence number value
// payload.
func (f *CipheringKeySequenceNumberFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a Ciphering key sequence number value without its
// IEI.
func (f *CipheringKeySequenceNumberFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: ciphering key sequence number fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: ciphering key sequence number has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xf8 != 0 {
		return fmt.Errorf("%w: ciphering key sequence number spare bit is nonzero", ErrInvalidValue)
	}
	*f = CipheringKeySequenceNumberFields{KeySequenceNumber: payload[0]}
	return nil
}
