package ie

import (
	"fmt"
	"io"
)

// ControlPlaneOnlyIndicationFields is the structured value of the Control plane
// only indication IE (TS 24.301 clause 9.9.4.23).
type ControlPlaneOnlyIndicationFields struct {
	CPOI bool
}

// NewControlPlaneOnlyIndicationWithFields creates a tagged Control plane only
// indication IE from fields.
func NewControlPlaneOnlyIndicationWithFields(fields *ControlPlaneOnlyIndicationFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewControlPlaneOnlyIndication(payload[0]), nil
}

// ControlPlaneOnlyIndication returns the IE as structured fields.
func (i *IE) ControlPlaneOnlyIndication() (*ControlPlaneOnlyIndicationFields, error) {
	if err := i.requireDefinition(ControlPlaneOnlyIndicationDefinition); err != nil {
		return nil, err
	}
	return ParseControlPlaneOnlyIndicationFields(i.Payload)
}

// ParseControlPlaneOnlyIndicationFields decodes a Control plane only indication
// value without its IEI.
func ParseControlPlaneOnlyIndicationFields(payload []byte) (*ControlPlaneOnlyIndicationFields, error) {
	fields := &ControlPlaneOnlyIndicationFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *ControlPlaneOnlyIndicationFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: control plane only indication fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a Control plane only indication value, without its
// IEI.
func (f *ControlPlaneOnlyIndicationFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	var encoded byte
	if f.CPOI {
		encoded |= 0x01
	}
	return []byte{encoded}, nil
}

// MarshalTo encodes fields into b.
func (f *ControlPlaneOnlyIndicationFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the Control plane only indication value payload.
func (f *ControlPlaneOnlyIndicationFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a Control plane only indication value without its IEI.
func (f *ControlPlaneOnlyIndicationFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: control plane only indication fields are nil", ErrInvalidValue)
	}
	if len(payload) == 0 {
		return io.ErrShortBuffer
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: control plane only indication has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xfe != 0 {
		return fmt.Errorf("%w: control plane only indication spare bits are nonzero", ErrInvalidValue)
	}
	*f = ControlPlaneOnlyIndicationFields{
		CPOI: payload[0]&0x01 != 0,
	}
	return nil
}
