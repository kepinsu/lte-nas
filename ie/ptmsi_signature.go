package ie

import (
	"fmt"
	"io"
)

// PTMSISignatureFields is the structured value of the P-TMSI signature IE
// (TS 24.301 clause 9.9.3.26). Signature is a 24-bit unsigned value.
type PTMSISignatureFields struct {
	Signature uint32
}

// NewPTMSISignatureWithFields creates a tagged P-TMSI signature IE from
// fields.
func NewPTMSISignatureWithFields(fields *PTMSISignatureFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPTMSISignature(payload), nil
}

// NewPositionalPTMSISignatureWithFields creates a positional P-TMSI signature
// IE from fields.
func NewPositionalPTMSISignatureWithFields(fields *PTMSISignatureFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalPTMSISignature(payload), nil
}

// PTMSISignature returns the IE as structured fields.
func (i *IE) PTMSISignature() (*PTMSISignatureFields, error) {
	if err := i.requireDefinition(PTMSISignatureDefinition); err != nil {
		return nil, err
	}
	return ParsePTMSISignatureFields(i.Payload)
}

// ParsePTMSISignatureFields decodes a P-TMSI signature value without its IEI.
func ParsePTMSISignatureFields(payload []byte) (*PTMSISignatureFields, error) {
	fields := &PTMSISignatureFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Signature fits 24 bits.
func (f *PTMSISignatureFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: P-TMSI signature fields are nil", ErrInvalidValue)
	}
	if f.Signature > 0x00ff_ffff {
		return fmt.Errorf("%w: P-TMSI signature %#08x exceeds 24 bits", ErrInvalidValue, f.Signature)
	}
	return nil
}

// Marshal encodes fields into a three-octet P-TMSI signature value.
func (f *PTMSISignatureFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.Signature >> 16), byte(f.Signature >> 8), byte(f.Signature)}, nil
}

// MarshalTo encodes fields into b.
func (f *PTMSISignatureFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 3 {
		return io.ErrShortBuffer
	}
	b[0] = byte(f.Signature >> 16)
	b[1] = byte(f.Signature >> 8)
	b[2] = byte(f.Signature)
	return nil
}

// MarshalLen returns the size of the P-TMSI signature value payload.
func (f *PTMSISignatureFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 3
}

// UnmarshalBinary decodes a three-octet P-TMSI signature value.
func (f *PTMSISignatureFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: P-TMSI signature fields are nil", ErrInvalidValue)
	}
	if len(payload) != 3 {
		return fmt.Errorf("%w: P-TMSI signature has %d value octets, want 3", ErrInvalidLength, len(payload))
	}
	*f = PTMSISignatureFields{Signature: uint32(payload[0])<<16 | uint32(payload[1])<<8 | uint32(payload[2])}
	return nil
}
