package ie

import (
	"fmt"
	"io"
)

// EPSAdditionalRequestResultFields is the structured value of the EPS additional
// request result IE (TS 24.301 clause 9.9.3.67).
type EPSAdditionalRequestResultFields struct {
	// PRD is the paging restriction decision.
	PRD uint8
}

// NewEPSAdditionalRequestResultWithFields creates a tagged EPS additional request
// result IE from fields.
func NewEPSAdditionalRequestResultWithFields(fields *EPSAdditionalRequestResultFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewEPSAdditionalRequestResult(payload[0]), nil
}

// EPSAdditionalRequestResult returns the IE as structured fields.
func (i *IE) EPSAdditionalRequestResult() (*EPSAdditionalRequestResultFields, error) {
	if err := i.requireDefinition(EPSAdditionalRequestResultDefinition); err != nil {
		return nil, err
	}
	return ParseEPSAdditionalRequestResultFields(i.Payload)
}

// ParseEPSAdditionalRequestResultFields decodes an EPS additional request result
// value without its IEI.
func ParseEPSAdditionalRequestResultFields(payload []byte) (*EPSAdditionalRequestResultFields, error) {
	fields := &EPSAdditionalRequestResultFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that PRD is a defined value.
func (f *EPSAdditionalRequestResultFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: EPS additional request result fields are nil", ErrInvalidValue)
	}
	if f.PRD > 0x03 {
		return fmt.Errorf("%w: PRD value %d is reserved", ErrInvalidValue, f.PRD)
	}
	return nil
}

// Marshal encodes fields into an EPS additional request result value, without its IEI.
func (f *EPSAdditionalRequestResultFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.PRD)}, nil
}

// MarshalTo encodes fields into b.
func (f *EPSAdditionalRequestResultFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the EPS additional request result value payload.
func (f *EPSAdditionalRequestResultFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes an EPS additional request result value without its IEI.
func (f *EPSAdditionalRequestResultFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: EPS additional request result fields are nil", ErrInvalidValue)
	}
	if len(payload) == 0 {
		return io.ErrShortBuffer
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: EPS additional request result has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xfc != 0 {
		return fmt.Errorf("%w: EPS additional request result spare bits are nonzero", ErrInvalidValue)
	}
	v := payload[0] & 0x03
	if v > 0x03 {
		return fmt.Errorf("%w: PRD value %d is reserved", ErrInvalidValue, v)
	}
	*f = EPSAdditionalRequestResultFields{
		PRD: v,
	}
	return nil
}
