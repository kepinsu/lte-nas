package ie

import (
	"fmt"
	"io"
)

// SupportedCodecListFields is the structured value of the tagged Supported codec
// list IE (TS 24.301 clause 9.9.2.10). It carries an opaque byte string between
// 3 and 255 octets.
type SupportedCodecListFields struct {
	Value []byte
}

// NewSupportedCodecListWithFields creates a tagged Supported codec list IE from
// fields.
func NewSupportedCodecListWithFields(fields *SupportedCodecListFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewSupportedCodecList(payload), nil
}

// NewPositionalSupportedCodecListWithFields creates a positional Supported codec
// list IE from fields.
func NewPositionalSupportedCodecListWithFields(fields *SupportedCodecListFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalSupportedCodecList(payload), nil
}

// SupportedCodecList returns the IE as structured fields.
func (i *IE) SupportedCodecList() (*SupportedCodecListFields, error) {
	if err := i.requireDefinition(SupportedCodecListDefinition); err != nil {
		return nil, err
	}
	return ParseSupportedCodecListFields(i.Payload)
}

// ParseSupportedCodecListFields decodes a supported codec list value without its
// IEI and length octet.
func ParseSupportedCodecListFields(payload []byte) (*SupportedCodecListFields, error) {
	fields := &SupportedCodecListFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the payload is within [3, 255] octets.
func (f *SupportedCodecListFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: supported codec list fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 3 {
		return fmt.Errorf("%w: supported codec list has %d value octets, want >= 3", ErrInvalidLength, n)
	}
	if n > 255 {
		return fmt.Errorf("%w: supported codec list has %d value octets, want <= 255", ErrInvalidLength, n)
	}
	return nil
}

// Marshal encodes fields into a supported codec list value without its IEI and
// length octet.
func (f *SupportedCodecListFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *SupportedCodecListFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the supported codec list value payload, or zero
// for a nil receiver or nil Value.
func (f *SupportedCodecListFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes a supported codec list value without its IEI and length
// octet. It copies the payload into a temporary struct first, leaving f unchanged
// on failure.
func (f *SupportedCodecListFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: supported codec list fields are nil", ErrInvalidValue)
	}
	if len(payload) < 3 {
		return fmt.Errorf("%w: supported codec list has %d value octets, want >= 3", ErrInvalidLength, len(payload))
	}
	if len(payload) > 255 {
		return fmt.Errorf("%w: supported codec list has %d value octets, want <= 255", ErrInvalidLength, len(payload))
	}
	var decoded SupportedCodecListFields
	decoded.Value = make([]byte, len(payload))
	copy(decoded.Value, payload)
	*f = decoded
	return nil
}
