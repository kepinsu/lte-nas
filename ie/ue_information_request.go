package ie

import (
	"fmt"
	"io"
)

// UEInformationRequestFields is the structured value of the UE information
// request IE (TS 24.301 clause 9.9.3.71).
type UEInformationRequestFields struct {
	// UCLIR indicates whether UE coarse location information is requested.
	UCLIR bool
}

// NewUEInformationRequestWithFields creates a tagged UE information request IE
// from fields.
func NewUEInformationRequestWithFields(fields *UEInformationRequestFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewUEInformationRequest(payload[0]), nil
}

// UEInformationRequest returns the IE as structured fields.
func (i *IE) UEInformationRequest() (*UEInformationRequestFields, error) {
	if err := i.requireDefinition(UEInformationRequestDefinition); err != nil {
		return nil, err
	}
	return ParseUEInformationRequestFields(i.Payload)
}

// ParseUEInformationRequestFields decodes a UE information request value without
// its IEI.
func ParseUEInformationRequestFields(payload []byte) (*UEInformationRequestFields, error) {
	fields := &UEInformationRequestFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *UEInformationRequestFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: UE information request fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a UE information request value, without its IEI.
func (f *UEInformationRequestFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	var encoded byte
	if f.UCLIR {
		encoded |= 0x01
	}
	return []byte{encoded}, nil
}

// MarshalTo encodes fields into b.
func (f *UEInformationRequestFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the UE information request value payload.
func (f *UEInformationRequestFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a UE information request value without its IEI.
func (f *UEInformationRequestFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: UE information request fields are nil", ErrInvalidValue)
	}
	if len(payload) == 0 {
		return io.ErrShortBuffer
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: UE information request has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xfe != 0 {
		return fmt.Errorf("%w: UE information request spare bits are nonzero", ErrInvalidValue)
	}
	*f = UEInformationRequestFields{
		UCLIR: payload[0]&0x01 != 0,
	}
	return nil
}
