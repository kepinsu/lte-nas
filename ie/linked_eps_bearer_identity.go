package ie

import (
	"fmt"
	"io"
)

// LinkedEPSBearerIdentity identifies the linked EPS bearer identity, as defined
// in TS 24.301 clause 9.9.4.6.
type LinkedEPSBearerIdentity uint8

// Linked EPS bearer identity values.
const (
	LinkedEPSBearerIdentityReserved LinkedEPSBearerIdentity = 0x00
	LinkedEPSBearerIdentity1        LinkedEPSBearerIdentity = 0x01
	LinkedEPSBearerIdentity2        LinkedEPSBearerIdentity = 0x02
	LinkedEPSBearerIdentity3        LinkedEPSBearerIdentity = 0x03
	LinkedEPSBearerIdentity4        LinkedEPSBearerIdentity = 0x04
	LinkedEPSBearerIdentity5        LinkedEPSBearerIdentity = 0x05
	LinkedEPSBearerIdentity6        LinkedEPSBearerIdentity = 0x06
	LinkedEPSBearerIdentity7        LinkedEPSBearerIdentity = 0x07
)

// LinkedEPSBearerIdentityFields is the structured value of the positional Linked
// EPS bearer identity IE.
type LinkedEPSBearerIdentityFields struct {
	Type LinkedEPSBearerIdentity
}

// NewPositionalLinkedEPSBearerIdentityWithFields creates a positional Linked EPS
// bearer identity IE from fields.
func NewPositionalLinkedEPSBearerIdentityWithFields(fields *LinkedEPSBearerIdentityFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalLinkedEPSBearerIdentity(payload[0]), nil
}

// LinkedEPSBearerIdentity returns the IE as structured fields.
func (i *IE) LinkedEPSBearerIdentity() (*LinkedEPSBearerIdentityFields, error) {
	if err := i.requireDefinition(LinkedEPSBearerIdentityDefinition); err != nil {
		return nil, err
	}
	return ParseLinkedEPSBearerIdentityFields(i.Payload)
}

// ParseLinkedEPSBearerIdentityFields decodes a Linked EPS bearer identity value
// without its IEI.
func ParseLinkedEPSBearerIdentityFields(payload []byte) (*LinkedEPSBearerIdentityFields, error) {
	fields := &LinkedEPSBearerIdentityFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Type is a defined Linked EPS bearer identity value.
func (f *LinkedEPSBearerIdentityFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: Linked EPS bearer identity fields are nil", ErrInvalidValue)
	}
	if !f.Type.valid() {
		return fmt.Errorf("%w: Linked EPS bearer identity %#02x is reserved", ErrInvalidValue, f.Type)
	}
	return nil
}

func (t LinkedEPSBearerIdentity) valid() bool {
	switch t {
	case LinkedEPSBearerIdentity1, LinkedEPSBearerIdentity2,
		LinkedEPSBearerIdentity3, LinkedEPSBearerIdentity4,
		LinkedEPSBearerIdentity5, LinkedEPSBearerIdentity6,
		LinkedEPSBearerIdentity7:
		return true
	default:
		return false
	}
}

// Marshal encodes fields into a Linked EPS bearer identity value without its IEI.
func (f *LinkedEPSBearerIdentityFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.Type)}, nil
}

// MarshalTo encodes fields into b.
func (f *LinkedEPSBearerIdentityFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	b[0] = byte(f.Type)
	return nil
}

// MarshalLen returns the size of the Linked EPS bearer identity value payload.
func (f *LinkedEPSBearerIdentityFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a Linked EPS bearer identity value without its IEI.
func (f *LinkedEPSBearerIdentityFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: Linked EPS bearer identity fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: Linked EPS bearer identity has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0] == 0 {
		return fmt.Errorf("%w: Linked EPS bearer identity is reserved", ErrInvalidValue)
	}
	decoded := LinkedEPSBearerIdentityFields{Type: LinkedEPSBearerIdentity(payload[0])}
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}
