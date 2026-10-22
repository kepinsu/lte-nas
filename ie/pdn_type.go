package ie

import (
	"fmt"
	"io"
)

// PDNType identifies the packet data network address family requested or
// selected for an EPS bearer, as defined in TS 24.301 clause 9.9.4.10.
type PDNType uint8

// PDN type values.
const (
	PDNTypeIPv4     PDNType = 0x01
	PDNTypeIPv6     PDNType = 0x02
	PDNTypeIPv4v6   PDNType = 0x03
	PDNTypeNonIP    PDNType = 0x05
	PDNTypeEthernet PDNType = 0x06
)

// PDNTypeFields is the structured value of the positional PDN type-1 IE.
type PDNTypeFields struct {
	Type PDNType
}

// NewPositionalPDNTypeWithFields creates a positional PDN type IE from
// fields.
func NewPositionalPDNTypeWithFields(fields *PDNTypeFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalPDNType(payload[0]), nil
}

// PDNType returns the PDN type IE as structured fields.
func (i *IE) PDNType() (*PDNTypeFields, error) {
	if err := i.requireDefinition(PDNTypeDefinition); err != nil {
		return nil, err
	}
	return ParsePDNTypeFields(i.Payload)
}

// ParsePDNTypeFields decodes a PDN type value without its IEI.
func ParsePDNTypeFields(payload []byte) (*PDNTypeFields, error) {
	fields := &PDNTypeFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Type is a defined PDN type value.
func (f *PDNTypeFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: PDN type fields are nil", ErrInvalidValue)
	}
	if !f.Type.valid() {
		return fmt.Errorf("%w: PDN type %#02x is reserved", ErrInvalidValue, f.Type)
	}
	return nil
}

func (t PDNType) valid() bool {
	switch t {
	case PDNTypeIPv4, PDNTypeIPv6, PDNTypeIPv4v6, PDNTypeNonIP, PDNTypeEthernet:
		return true
	default:
		return false
	}
}

// Marshal encodes fields into a PDN type value without its IEI.
func (f *PDNTypeFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.Type)}, nil
}

// MarshalTo encodes fields into b.
func (f *PDNTypeFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	b[0] = byte(f.Type)
	return nil
}

// MarshalLen returns the size of the PDN type value payload.
func (f *PDNTypeFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a PDN type value without its IEI.
func (f *PDNTypeFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: PDN type fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: PDN type has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xf0 != 0 {
		return fmt.Errorf("%w: PDN type spare bits are nonzero", ErrInvalidValue)
	}
	decoded := PDNTypeFields{Type: PDNType(payload[0])}
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}
