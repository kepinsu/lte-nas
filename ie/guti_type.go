package ie

import (
	"fmt"
	"io"
)

// GUTIType identifies the GUTI type, as defined in TS 24.301 clause 9.9.3.45.
type GutIType uint8

// GUTI type values.
const (
	GutITypeEPS     GutIType = 0x00
	GutITypeNonEPS  GutIType = 0x01
	GutITypeReserved GutIType = 0x02
)

// GUTITypeFields is the structured value of the tagged GUTI type IE.
type GUTITypeFields struct {
	Type GutIType
}

// NewGUTITypeWithFields creates a tagged GUTI type IE from fields.
func NewGUTITypeWithFields(fields *GUTITypeFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewGUTIType(payload[0]), nil
}

// GUTIType returns the IE as structured fields.
func (i *IE) GUTIType() (*GUTITypeFields, error) {
	if err := i.requireDefinition(GUTITypeDefinition); err != nil {
		return nil, err
	}
	return ParseGUTITypeFields(i.Payload)
}

// ParseGUTITypeFields decodes a GUTI type value without its IEI.
func ParseGUTITypeFields(payload []byte) (*GUTITypeFields, error) {
	fields := &GUTITypeFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Type is a defined GUTI type value.
func (f *GUTITypeFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: GUTI type fields are nil", ErrInvalidValue)
	}
	if !f.Type.valid() {
		return fmt.Errorf("%w: GUTI type %#02x is reserved", ErrInvalidValue, f.Type)
	}
	return nil
}

func (t GutIType) valid() bool {
	switch t {
	case GutITypeEPS, GutITypeNonEPS:
		return true
	default:
		return false
	}
}

// Marshal encodes fields into a GUTI type value without its IEI.
func (f *GUTITypeFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.Type)}, nil
}

// MarshalTo encodes fields into b.
func (f *GUTITypeFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	b[0] = byte(f.Type)
	return nil
}

// MarshalLen returns the size of the GUTI type value payload.
func (f *GUTITypeFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a GUTI type value without its IEI.
func (f *GUTITypeFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: GUTI type fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: GUTI type has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	decoded := GUTITypeFields{Type: GutIType(payload[0])}
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}
