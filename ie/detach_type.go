package ie

import (
	"fmt"
	"io"
)

// DetachType identifies the detach type, as defined in TS 24.301 clause 9.9.3.7.
type DetachType uint8

// DetachTypeFields is the structured value of the tagged Detach type IE (TS
// 24.301 clause 9.9.3.7).
type DetachTypeFields struct {
	Type DetachType
}

// DetachTypeEPSOnly identifies a detach for the EPS attach procedure.
const DetachTypeEPSOnly DetachType = 0x01

// DetachTypeCombinedEPSIMSI identifies a detach for the combined EPS attach
// and IMSI detach procedure.
const DetachTypeCombinedEPSIMSI DetachType = 0x02

// NewDetachTypeWithFields creates a tagged Detach type IE from fields.
func NewDetachTypeWithFields(fields *DetachTypeFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalDetachType(payload[0]), nil
}

// NewPositionalDetachTypeWithFields creates a positional Detach type IE from
// fields.
func NewPositionalDetachTypeWithFields(fields *DetachTypeFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalDetachType(payload[0]), nil
}

// DetachType returns the IE as structured fields.
func (i *IE) DetachType() (*DetachTypeFields, error) {
	if err := i.requireDefinition(DetachTypeDefinition); err != nil {
		return nil, err
	}
	return ParseDetachTypeFields(i.Payload)
}

// ParseDetachTypeFields decodes a detach type value without its IEI.
func ParseDetachTypeFields(payload []byte) (*DetachTypeFields, error) {
	fields := &DetachTypeFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Type is a defined detach type value.
func (f *DetachTypeFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: detach type fields are nil", ErrInvalidValue)
	}
	if !f.Type.valid() {
		return fmt.Errorf("%w: detach type %#02x is reserved", ErrInvalidValue, f.Type)
	}
	return nil
}

func (t DetachType) valid() bool {
	switch t {
	case DetachTypeEPSOnly, DetachTypeCombinedEPSIMSI:
		return true
	default:
		return false
	}
}

// Marshal encodes fields into a detach type value without its IEI.
func (f *DetachTypeFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.Type)}, nil
}

// MarshalTo encodes fields into b.
func (f *DetachTypeFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	b[0] = byte(f.Type)
	return nil
}

// MarshalLen returns the size of the detach type value payload.
func (f *DetachTypeFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a detach type value without its IEI. It leaves f
// unchanged when the payload is invalid.
func (f *DetachTypeFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: detach type fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: detach type has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xfc != 0 {
		return fmt.Errorf("%w: detach type reserved bits are nonzero", ErrInvalidValue)
	}
	var decoded DetachTypeFields
	decoded.Type = DetachType(payload[0])
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}
