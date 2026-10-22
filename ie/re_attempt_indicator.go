package ie

import (
	"fmt"
	"io"
)

// ReAttemptIndicatorFields is the structured value of the Re-attempt indicator IE
// (TS 24.301 clause 9.9.4.13A).
type ReAttemptIndicatorFields struct {
	// RATC indicates whether the UE is allowed to re-attempt.
	RATC bool
}

// NewReAttemptIndicatorWithFields creates a tagged Re-attempt indicator IE from
// fields.
func NewReAttemptIndicatorWithFields(fields *ReAttemptIndicatorFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewReAttemptIndicator(payload[0]), nil
}

// ReAttemptIndicator returns the IE as structured fields.
func (i *IE) ReAttemptIndicator() (*ReAttemptIndicatorFields, error) {
	if err := i.requireDefinition(ReAttemptIndicatorDefinition); err != nil {
		return nil, err
	}
	return ParseReAttemptIndicatorFields(i.Payload)
}

// ParseReAttemptIndicatorFields decodes a Re-attempt indicator value without its IEI.
func ParseReAttemptIndicatorFields(payload []byte) (*ReAttemptIndicatorFields, error) {
	fields := &ReAttemptIndicatorFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *ReAttemptIndicatorFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: re-attempt indicator fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a Re-attempt indicator value, without its IEI.
func (f *ReAttemptIndicatorFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	var encoded byte
	if f.RATC {
		encoded |= 0x01
	}
	return []byte{encoded}, nil
}

// MarshalTo encodes fields into b.
func (f *ReAttemptIndicatorFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the Re-attempt indicator value payload.
func (f *ReAttemptIndicatorFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a Re-attempt indicator value without its IEI.
func (f *ReAttemptIndicatorFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: re-attempt indicator fields are nil", ErrInvalidValue)
	}
	if len(payload) == 0 {
		return io.ErrShortBuffer
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: re-attempt indicator has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xfe != 0 {
		return fmt.Errorf("%w: re-attempt indicator spare bits are nonzero", ErrInvalidValue)
	}
	*f = ReAttemptIndicatorFields{
		RATC: payload[0]&0x01 != 0,
	}
	return nil
}
