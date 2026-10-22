package ie

import (
	"fmt"
	"io"
)

// AdditionalUpdateTypeValue identifies the type of request for an attach or tracking
// area updating procedure (TS 24.301 clause 9.9.3.0B).
type AdditionalUpdateTypeValue uint8

const (
	AdditionalUpdateTypeNoAdditionalInfo  AdditionalUpdateTypeValue = iota
	AdditionalUpdateTypeSMSOnly
)

// AdditionalUpdateTypeFields is the structured value of the Additional update
// type IE (TS 24.301 clause 9.9.3.0B).
type AdditionalUpdateTypeFields struct {
	// AUTV identifies the type of request.
	AUTV AdditionalUpdateTypeValue
	// SAF indicates whether keeping the NAS signalling connection is required.
	SAF bool
	// PNBsCIoT identifies the preferred CIoT network behaviour.
	PNBsCIoT uint8
}

// NewAdditionalUpdateTypeWithFields creates a tagged Additional update type IE
// from fields.
func NewAdditionalUpdateTypeWithFields(fields *AdditionalUpdateTypeFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewAdditionalUpdateType(payload[0]), nil
}

// NewPositionalAdditionalUpdateTypeWithFields creates a positional Additional
// update type IE from fields.
func NewPositionalAdditionalUpdateTypeWithFields(fields *AdditionalUpdateTypeFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalAdditionalUpdateType(payload[0]), nil
}

// AdditionalUpdateType returns the IE as structured fields.
func (i *IE) AdditionalUpdateType() (*AdditionalUpdateTypeFields, error) {
	if err := i.requireDefinition(AdditionalUpdateTypeDefinition); err != nil {
		return nil, err
	}
	return ParseAdditionalUpdateTypeFields(i.Payload)
}

// ParseAdditionalUpdateTypeFields decodes an Additional update type value
// without its IEI.
func ParseAdditionalUpdateTypeFields(payload []byte) (*AdditionalUpdateTypeFields, error) {
	fields := &AdditionalUpdateTypeFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that AUTV and PNBsCIoT are defined values.
func (f *AdditionalUpdateTypeFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: additional update type fields are nil", ErrInvalidValue)
	}
	if f.AUTV > AdditionalUpdateTypeSMSOnly {
		return fmt.Errorf("%w: additional update type AUTV %d is reserved", ErrInvalidValue, f.AUTV)
	}
	if f.PNBsCIoT > 0x03 {
		return fmt.Errorf("%w: additional update type PNB-CIoT %d is reserved", ErrInvalidValue, f.PNBsCIoT)
	}
	return nil
}

// Marshal encodes fields into an Additional update type value, without its IEI.
func (f *AdditionalUpdateTypeFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	var encoded byte
	encoded |= byte(f.AUTV) & 0x01
	if f.SAF {
		encoded |= 0x04
	}
	encoded |= byte(f.PNBsCIoT) << 5
	return []byte{encoded}, nil
}

// MarshalTo encodes fields into b.
func (f *AdditionalUpdateTypeFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the Additional update type value payload.
func (f *AdditionalUpdateTypeFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes an Additional update type value without its IEI.
func (f *AdditionalUpdateTypeFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: additional update type fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: additional update type has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0x9a != 0 {
		return fmt.Errorf("%w: additional update type spare bits are nonzero", ErrInvalidValue)
	}
	*f = AdditionalUpdateTypeFields{
		AUTV:      AdditionalUpdateTypeValue(payload[0] & 0x01),
		SAF:       payload[0]&0x04 != 0,
		PNBsCIoT:  payload[0] >> 5,
	}
	return nil
}
