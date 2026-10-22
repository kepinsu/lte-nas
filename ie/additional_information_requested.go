package ie

import (
	"fmt"
	"io"
)

// AdditionalInformationRequestedFields is the structured value of the
// Additional information requested IE (TS 24.301 clause 9.9.3.55).
type AdditionalInformationRequestedFields struct {
	CipherKey bool
}

// NewAdditionalInformationRequestedWithFields creates a tagged Additional
// information requested IE from fields.
func NewAdditionalInformationRequestedWithFields(fields *AdditionalInformationRequestedFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewAdditionalInformationRequested(payload[0]), nil
}

// AdditionalInformationRequested returns the IE as structured fields.
func (i *IE) AdditionalInformationRequested() (*AdditionalInformationRequestedFields, error) {
	if err := i.requireDefinition(AdditionalInformationRequestedDefinition); err != nil {
		return nil, err
	}
	return ParseAdditionalInformationRequestedFields(i.Payload)
}

// ParseAdditionalInformationRequestedFields decodes an Additional information
// requested value without its IEI.
func ParseAdditionalInformationRequestedFields(payload []byte) (*AdditionalInformationRequestedFields, error) {
	fields := &AdditionalInformationRequestedFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *AdditionalInformationRequestedFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: additional information requested fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into an Additional information requested value, without its
// IEI.
func (f *AdditionalInformationRequestedFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	var encoded byte
	if f.CipherKey {
		encoded |= 0x01
	}
	return []byte{encoded}, nil
}

// MarshalTo encodes fields into b.
func (f *AdditionalInformationRequestedFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the Additional information requested value payload.
func (f *AdditionalInformationRequestedFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes an Additional information requested value without its IEI.
func (f *AdditionalInformationRequestedFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: additional information requested fields are nil", ErrInvalidValue)
	}
	if len(payload) == 0 {
		return io.ErrShortBuffer
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: additional information requested has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xfe != 0 {
		return fmt.Errorf("%w: additional information requested spare bits are nonzero", ErrInvalidValue)
	}
	*f = AdditionalInformationRequestedFields{
		CipherKey: payload[0]&0x01 != 0,
	}
	return nil
}
