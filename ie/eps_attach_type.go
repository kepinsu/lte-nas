package ie

import (
	"fmt"
	"io"
)

// EPSAttachType identifies the requested EPS attach procedure, as defined in
// TS 24.301 clause 9.9.3.11.
type EPSAttachType uint8

// EPS attach type values.
const (
	EPSAttachTypeEPS                EPSAttachType = 0x01
	EPSAttachTypeCombinedEPSAndIMSI EPSAttachType = 0x02
	EPSAttachTypeEPSRLOS            EPSAttachType = 0x03
	EPSAttachTypeEPSEmergency       EPSAttachType = 0x06
)

// EPSAttachTypeFields is the structured value of the EPS attach type IE.
//
// The IE is a positional type-1 IE. Its IEI is not part of these fields; the
// encoder stores Type in the low nibble, and MarshalHalfOctets combines it with
// another positional type-1 IE when required by a message format.
type EPSAttachTypeFields struct {
	Type EPSAttachType
}

// NewPositionalEPSAttachTypeWithFields creates a positional EPS attach type IE
// from fields. NewPositionalEPSAttachType remains available for raw values.
func NewPositionalEPSAttachTypeWithFields(fields *EPSAttachTypeFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalEPSAttachType(payload[0]), nil
}

// EPSAttachType returns the EPS attach type IE as structured fields.
func (i *IE) EPSAttachType() (*EPSAttachTypeFields, error) {
	if err := i.requireDefinition(EPSAttachTypeDefinition); err != nil {
		return nil, err
	}
	return ParseEPSAttachTypeFields(i.Payload)
}

// ParseEPSAttachTypeFields decodes an EPS attach type value without its IEI.
func ParseEPSAttachTypeFields(payload []byte) (*EPSAttachTypeFields, error) {
	fields := &EPSAttachTypeFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Type is a defined EPS attach type value.
func (f *EPSAttachTypeFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: EPS attach type fields are nil", ErrInvalidValue)
	}
	if !f.Type.valid() {
		return fmt.Errorf("%w: EPS attach type %#02x is reserved", ErrInvalidValue, f.Type)
	}
	return nil
}

// Marshal encodes fields into an EPS attach type value without its IEI.
func (f *EPSAttachTypeFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.Type)}, nil
}

// MarshalTo encodes fields into b. It returns io.ErrShortBuffer without
// modifying b when b is too short.
func (f *EPSAttachTypeFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	b[0] = byte(f.Type)
	return nil
}

// MarshalLen returns the size of the EPS attach type value payload.
func (f *EPSAttachTypeFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes an EPS attach type value without its IEI. It leaves
// f unchanged when the payload is invalid.
func (f *EPSAttachTypeFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: EPS attach type fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: EPS attach type has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xf0 != 0 {
		return fmt.Errorf("%w: EPS attach type spare bits are nonzero", ErrInvalidValue)
	}
	decoded := EPSAttachTypeFields{Type: EPSAttachType(payload[0])}
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}

func (t EPSAttachType) valid() bool {
	switch t {
	case EPSAttachTypeEPS, EPSAttachTypeCombinedEPSAndIMSI, EPSAttachTypeEPSRLOS, EPSAttachTypeEPSEmergency:
		return true
	default:
		return false
	}
}
