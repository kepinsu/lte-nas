package ie

import (
	"fmt"
	"io"
)

// UECoarseLocationInformationFields is the structured value of the UE coarse
// location information IE (TS 24.301 clause 9.9.3.72). The six-octet payload
// is an opaque location information element.
type UECoarseLocationInformationFields struct {
	// Payload contains the 6-octet UE coarse location information value.
	Payload [6]byte
}

// NewUECoarseLocationInformationWithFields creates a tagged UE coarse location
// information IE from fields.
func NewUECoarseLocationInformationWithFields(fields *UECoarseLocationInformationFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewUECoarseLocationInformation(payload), nil
}

// NewPositionalUECoarseLocationInformationWithFields creates a positional UE coarse
// location information IE from fields.
func NewPositionalUECoarseLocationInformationWithFields(fields *UECoarseLocationInformationFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalUECoarseLocationInformation(payload), nil
}

// UECoarseLocationInformation returns the IE as structured fields.
func (i *IE) UECoarseLocationInformation() (*UECoarseLocationInformationFields, error) {
	if err := i.requireDefinition(UECoarseLocationInformationDefinition); err != nil {
		return nil, err
	}
	return ParseUECoarseLocationInformationFields(i.Payload)
}

// ParseUECoarseLocationInformationFields decodes a UE coarse location information
// value without its IEI and length octet.
func ParseUECoarseLocationInformationFields(payload []byte) (*UECoarseLocationInformationFields, error) {
	fields := &UECoarseLocationInformationFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *UECoarseLocationInformationFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: UE coarse location information fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a 6-octet UE coarse location information value.
func (f *UECoarseLocationInformationFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Payload))
	copy(payload, f.Payload[:])
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *UECoarseLocationInformationFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 6 {
		return io.ErrShortBuffer
	}
	copy(b, f.Payload[:])
	return nil
}

// MarshalLen returns the size of the UE coarse location information value payload.
func (f *UECoarseLocationInformationFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 6
}

// UnmarshalBinary decodes a 6-octet UE coarse location information value.
func (f *UECoarseLocationInformationFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: UE coarse location information fields are nil", ErrInvalidValue)
	}
	if len(payload) != 6 {
		return fmt.Errorf("%w: UE coarse location information has %d value octets, want 6", ErrInvalidLength, len(payload))
	}
	var decoded UECoarseLocationInformationFields
	copy(decoded.Payload[:], payload)
	*f = decoded
	return nil
}
