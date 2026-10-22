package ie

import (
	"fmt"
	"io"
)

// AdditionalInformationFields is the structured value of the tagged Additional
// information IE (TS 24.301 clause 9.9.2.0). It carries an opaque byte string
// between 1 and 255 octets.
type AdditionalInformationFields struct {
	Value []byte
}

// NewAdditionalInformationWithFields creates a tagged Additional information IE
// from fields.
func NewAdditionalInformationWithFields(fields *AdditionalInformationFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewAdditionalInformation(payload), nil
}

// NewPositionalAdditionalInformationWithFields creates a positional Additional
// information IE from fields.
func NewPositionalAdditionalInformationWithFields(fields *AdditionalInformationFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalAdditionalInformation(payload), nil
}

// AdditionalInformation returns the IE as structured fields.
func (i *IE) AdditionalInformation() (*AdditionalInformationFields, error) {
	if err := i.requireDefinition(AdditionalInformationDefinition); err != nil {
		return nil, err
	}
	return ParseAdditionalInformationFields(i.Payload)
}

// ParseAdditionalInformationFields decodes an additional information value without
// its IEI and length octet.
func ParseAdditionalInformationFields(payload []byte) (*AdditionalInformationFields, error) {
	fields := &AdditionalInformationFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the payload is within [1, 255] octets.
func (f *AdditionalInformationFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: additional information fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 1 {
		return fmt.Errorf("%w: additional information has %d value octets, want >= 1", ErrInvalidLength, n)
	}
	if n > 255 {
		return fmt.Errorf("%w: additional information has %d value octets, want <= 255", ErrInvalidLength, n)
	}
	return nil
}

// Marshal encodes fields into an additional information value without its IEI and
// length octet.
func (f *AdditionalInformationFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *AdditionalInformationFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the additional information value payload, or
// zero for a nil receiver or nil Value.
func (f *AdditionalInformationFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes an additional information value without its IEI and
// length octet. It copies the payload into a temporary struct first, leaving f
// unchanged on failure.
func (f *AdditionalInformationFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: additional information fields are nil", ErrInvalidValue)
	}
	if len(payload) < 1 {
		return fmt.Errorf("%w: additional information has %d value octets, want >= 1", ErrInvalidLength, len(payload))
	}
	if len(payload) > 255 {
		return fmt.Errorf("%w: additional information has %d value octets, want <= 255", ErrInvalidLength, len(payload))
	}
	var decoded AdditionalInformationFields
	decoded.Value = make([]byte, len(payload))
	copy(decoded.Value, payload)
	*f = decoded
	return nil
}
