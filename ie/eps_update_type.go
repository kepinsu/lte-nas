package ie

import (
	"fmt"
	"io"
)

// EPSUpdateType identifies the area covered by an EPS updating procedure, as
// defined in TS 24.301 clause 9.9.3.14.
type EPSUpdateType uint8

// EPS update type values.
const (
	EPSUpdateTypeTrackingArea                                EPSUpdateType = 0x00
	EPSUpdateTypeCombinedTrackingAreaAndLocationArea         EPSUpdateType = 0x01
	EPSUpdateTypeCombinedTrackingAreaAndLocationAreaWithIMSI EPSUpdateType = 0x02
	EPSUpdateTypePeriodic                                    EPSUpdateType = 0x03
)

// EPSUpdateTypeFields is the structured value of the positional EPS update
// type-1 IE.
type EPSUpdateTypeFields struct {
	Type   EPSUpdateType
	Active bool
}

// NewPositionalEPSUpdateTypeWithFields creates a positional EPS update type
// IE from fields.
func NewPositionalEPSUpdateTypeWithFields(fields *EPSUpdateTypeFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalEPSUpdateType(payload[0]), nil
}

// EPSUpdateType returns the EPS update type IE as structured fields.
func (i *IE) EPSUpdateType() (*EPSUpdateTypeFields, error) {
	if err := i.requireDefinition(EPSUpdateTypeDefinition); err != nil {
		return nil, err
	}
	return ParseEPSUpdateTypeFields(i.Payload)
}

// ParseEPSUpdateTypeFields decodes an EPS update type value without its IEI.
func ParseEPSUpdateTypeFields(payload []byte) (*EPSUpdateTypeFields, error) {
	fields := &EPSUpdateTypeFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Type is a defined EPS update type value.
func (f *EPSUpdateTypeFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: EPS update type fields are nil", ErrInvalidValue)
	}
	switch f.Type {
	case EPSUpdateTypeTrackingArea,
		EPSUpdateTypeCombinedTrackingAreaAndLocationArea,
		EPSUpdateTypeCombinedTrackingAreaAndLocationAreaWithIMSI,
		EPSUpdateTypePeriodic:
		return nil
	default:
		return fmt.Errorf("%w: EPS update type %#02x is not defined", ErrInvalidValue, f.Type)
	}
}

// Marshal encodes fields into an EPS update type value without its IEI.
func (f *EPSUpdateTypeFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	value := byte(f.Type)
	if f.Active {
		value |= 0x08
	}
	return []byte{value}, nil
}

// MarshalTo encodes fields into b.
func (f *EPSUpdateTypeFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the EPS update type value payload.
func (f *EPSUpdateTypeFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes an EPS update type value without its IEI.
func (f *EPSUpdateTypeFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: EPS update type fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: EPS update type has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xf0 != 0 {
		return fmt.Errorf("%w: EPS update type spare bits are nonzero", ErrInvalidValue)
	}
	decoded := EPSUpdateTypeFields{
		Type:   EPSUpdateType(payload[0] & 0x07),
		Active: payload[0]&0x08 != 0,
	}
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}
