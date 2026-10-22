package ie

import (
	"fmt"
	"io"
)

// ExtendedEMMCauseFields is the structured value of the Extended EMM cause IE
// (TS 24.301 clause 9.9.3.26A).
type ExtendedEMMCauseFields struct {
	EUTRANNotAllowed                bool
	EPSOptimizationInfoNotSupported bool
	NBIoTNotAllowed                 bool
}

// NewExtendedEMMCauseWithFields creates a tagged Extended EMM cause IE from
// fields.
func NewExtendedEMMCauseWithFields(fields *ExtendedEMMCauseFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewExtendedEMMCause(payload[0]), nil
}

// ExtendedEMMCause returns the IE as structured fields.
func (i *IE) ExtendedEMMCause() (*ExtendedEMMCauseFields, error) {
	if err := i.requireDefinition(ExtendedEMMCauseDefinition); err != nil {
		return nil, err
	}
	return ParseExtendedEMMCauseFields(i.Payload)
}

// ParseExtendedEMMCauseFields decodes an Extended EMM cause value without its
// IEI.
func ParseExtendedEMMCauseFields(payload []byte) (*ExtendedEMMCauseFields, error) {
	fields := &ExtendedEMMCauseFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *ExtendedEMMCauseFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: extended EMM cause fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into the low nibble of an Extended EMM cause value,
// without its IEI.
func (f *ExtendedEMMCauseFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	var encoded byte
	if f.EUTRANNotAllowed {
		encoded |= 0x01
	}
	if f.EPSOptimizationInfoNotSupported {
		encoded |= 0x02
	}
	if f.NBIoTNotAllowed {
		encoded |= 0x04
	}
	return []byte{encoded}, nil
}

// MarshalTo encodes fields into b.
func (f *ExtendedEMMCauseFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the Extended EMM cause value payload.
func (f *ExtendedEMMCauseFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes an Extended EMM cause value without its IEI.
func (f *ExtendedEMMCauseFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: extended EMM cause fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: extended EMM cause has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xf8 != 0 {
		return fmt.Errorf("%w: extended EMM cause spare bit is nonzero", ErrInvalidValue)
	}
	*f = ExtendedEMMCauseFields{
		EUTRANNotAllowed:                payload[0]&0x01 != 0,
		EPSOptimizationInfoNotSupported: payload[0]&0x02 != 0,
		NBIoTNotAllowed:                 payload[0]&0x04 != 0,
	}
	return nil
}
