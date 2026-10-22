package ie

import (
	"fmt"
	"io"
)

// SsCode identifies the SS code, as defined in TS 24.301 clause 9.9.3.39.
type SsCode uint8

// SS code values.
const (
	SsCodeReserved      SsCode = 0x00
	SsCodeValid         SsCode = 0x01
	SsCodeNotApplicable SsCode = 0x02
)

// SSCodeFields is the structured value of the tagged SS Code IE.
type SSCodeFields struct {
	Type SsCode
}

// NewSSCodeWithFields creates a tagged SS Code IE from fields.
func NewSSCodeWithFields(fields *SSCodeFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewSSCode(payload[0]), nil
}

// SSCode returns the IE as structured fields.
func (i *IE) SSCode() (*SSCodeFields, error) {
	if err := i.requireDefinition(SSCodeDefinition); err != nil {
		return nil, err
	}
	return ParseSSCodeFields(i.Payload)
}

// ParseSSCodeFields decodes an SS Code value without its IEI.
func ParseSSCodeFields(payload []byte) (*SSCodeFields, error) {
	fields := &SSCodeFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Type is a defined SS code value.
func (f *SSCodeFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: SS code fields are nil", ErrInvalidValue)
	}
	if !f.Type.valid() {
		return fmt.Errorf("%w: SS code %#02x is reserved", ErrInvalidValue, f.Type)
	}
	return nil
}

func (t SsCode) valid() bool {
	switch t {
	case SsCodeValid, SsCodeNotApplicable:
		return true
	default:
		return false
	}
}

// Marshal encodes fields into an SS Code value without its IEI.
func (f *SSCodeFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.Type)}, nil
}

// MarshalTo encodes fields into b.
func (f *SSCodeFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	b[0] = byte(f.Type)
	return nil
}

// MarshalLen returns the size of the SS Code value payload.
func (f *SSCodeFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes an SS Code value without its IEI.
func (f *SSCodeFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: SS code fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: SS code has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0] != 0x01 && payload[0] != 0x02 {
		return fmt.Errorf("%w: SS code is reserved", ErrInvalidValue)
	}
	decoded := SSCodeFields{Type: SsCode(payload[0])}
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}
