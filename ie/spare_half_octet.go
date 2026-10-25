package ie

import (
	"fmt"
	"io"
)

// SpareHalfOctetFields is the structured value of the Spare half octet IE (TS
// 24.301 clause 9.9.2.9). The low nibble is used by the paired IE; this type
// represents only the high nibble.
type SpareHalfOctetFields struct{}

// NewSpareHalfOctetWithFields creates a tagged Spare half octet IE from fields.
func NewSpareHalfOctetWithFields(fields *SpareHalfOctetFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalSpareHalfOctet(payload[0]), nil
}

// NewPositionalSpareHalfOctetWithFields creates a positional Spare half octet IE
// from fields.
func NewPositionalSpareHalfOctetWithFields(fields *SpareHalfOctetFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalSpareHalfOctet(payload[0]), nil
}

// SpareHalfOctet returns the IE as structured fields.
func (i *IE) SpareHalfOctet() (*SpareHalfOctetFields, error) {
	if err := i.requireDefinition(SpareHalfOctetDefinition); err != nil {
		return nil, err
	}
	return ParseSpareHalfOctetFields(i.Payload)
}

// ParseSpareHalfOctetFields decodes a Spare half octet value without its IEI.
func ParseSpareHalfOctetFields(payload []byte) (*SpareHalfOctetFields, error) {
	fields := &SpareHalfOctetFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *SpareHalfOctetFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: spare half octet fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a Spare half octet value, without its IEI.
func (f *SpareHalfOctetFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{0x00}, nil
}

// MarshalTo encodes fields into b.
func (f *SpareHalfOctetFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the Spare half octet value payload.
func (f *SpareHalfOctetFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a Spare half octet value without its IEI.
func (f *SpareHalfOctetFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: spare half octet fields are nil", ErrInvalidValue)
	}
	if len(payload) == 0 {
		return io.ErrShortBuffer
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: spare half octet has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0x0f != 0 {
		return fmt.Errorf("%w: spare half octet low nibble is nonzero", ErrInvalidValue)
	}
	*f = SpareHalfOctetFields{}
	return nil
}
