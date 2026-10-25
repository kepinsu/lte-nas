package ie

import (
	"fmt"
	"io"
)

// ESMInformationTransferFlagFields is the structured value of the ESM
// information transfer flag IE (TS 24.301 clause 9.9.4.5).
type ESMInformationTransferFlagFields struct {
	EIT bool
}

// NewESMInformationTransferFlagWithFields creates a tagged ESM information
// transfer flag IE from fields.
func NewESMInformationTransferFlagWithFields(fields *ESMInformationTransferFlagFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewESMInformationTransferFlag(payload[0]), nil
}

// ESMInformationTransferFlag returns the IE as structured fields.
func (i *IE) ESMInformationTransferFlag() (*ESMInformationTransferFlagFields, error) {
	if err := i.requireDefinition(ESMInformationTransferFlagDefinition); err != nil {
		return nil, err
	}
	return ParseESMInformationTransferFlagFields(i.Payload)
}

// ParseESMInformationTransferFlagFields decodes an ESM information transfer flag
// value without its IEI.
func ParseESMInformationTransferFlagFields(payload []byte) (*ESMInformationTransferFlagFields, error) {
	fields := &ESMInformationTransferFlagFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *ESMInformationTransferFlagFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: ESM information transfer flag fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into an ESM information transfer flag value, without its
// IEI.
func (f *ESMInformationTransferFlagFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	var encoded byte
	if f.EIT {
		encoded |= 0x01
	}
	return []byte{encoded}, nil
}

// MarshalTo encodes fields into b.
func (f *ESMInformationTransferFlagFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the ESM information transfer flag value payload.
func (f *ESMInformationTransferFlagFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes an ESM information transfer flag value without its IEI.
func (f *ESMInformationTransferFlagFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: ESM information transfer flag fields are nil", ErrInvalidValue)
	}
	if len(payload) == 0 {
		return io.ErrShortBuffer
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: ESM information transfer flag has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xfe != 0 {
		return fmt.Errorf("%w: ESM information transfer flag spare bits are nonzero", ErrInvalidValue)
	}
	*f = ESMInformationTransferFlagFields{
		EIT: payload[0]&0x01 != 0,
	}
	return nil
}
