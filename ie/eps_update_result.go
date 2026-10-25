package ie

import (
	"fmt"
	"io"
)

// EPSUpdateResult identifies the result of an EPS updating procedure, as
// defined in TS 24.301 clause 9.9.3.13.
type EPSUpdateResult uint8

// EPS update result values.
const (
	EPSUpdateResultTrackingAreaUpdated                               EPSUpdateResult = 0x00
	EPSUpdateResultCombinedTrackingAreaAndLocationAreaUpdated        EPSUpdateResult = 0x01
	EPSUpdateResultTrackingAreaUpdatedWithISR                        EPSUpdateResult = 0x04
	EPSUpdateResultCombinedTrackingAreaAndLocationAreaUpdatedWithISR EPSUpdateResult = 0x05
)

// EPSUpdateResultFields is the structured value of the positional EPS update
// result type-1 IE.
type EPSUpdateResultFields struct {
	Result EPSUpdateResult
}

// NewPositionalEPSUpdateResultWithFields creates a positional EPS update
// result IE from fields.
func NewPositionalEPSUpdateResultWithFields(fields *EPSUpdateResultFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalEPSUpdateResult(payload[0]), nil
}

// EPSUpdateResult returns the EPS update result IE as structured fields.
func (i *IE) EPSUpdateResult() (*EPSUpdateResultFields, error) {
	if err := i.requireDefinition(EPSUpdateResultDefinition); err != nil {
		return nil, err
	}
	return ParseEPSUpdateResultFields(i.Payload)
}

// ParseEPSUpdateResultFields decodes an EPS update result value without its
// IEI.
func ParseEPSUpdateResultFields(payload []byte) (*EPSUpdateResultFields, error) {
	fields := &EPSUpdateResultFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Result is defined by TS 24.301.
func (f *EPSUpdateResultFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: EPS update result fields are nil", ErrInvalidValue)
	}
	switch f.Result {
	case EPSUpdateResultTrackingAreaUpdated,
		EPSUpdateResultCombinedTrackingAreaAndLocationAreaUpdated,
		EPSUpdateResultTrackingAreaUpdatedWithISR,
		EPSUpdateResultCombinedTrackingAreaAndLocationAreaUpdatedWithISR:
		return nil
	default:
		return fmt.Errorf("%w: EPS update result %#02x is reserved", ErrInvalidValue, f.Result)
	}
}

// Marshal encodes fields into an EPS update result value without its IEI.
func (f *EPSUpdateResultFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.Result)}, nil
}

// MarshalTo encodes fields into b.
func (f *EPSUpdateResultFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	b[0] = byte(f.Result)
	return nil
}

// MarshalLen returns the size of the EPS update result value payload.
func (f *EPSUpdateResultFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes an EPS update result value without its IEI.
func (f *EPSUpdateResultFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: EPS update result fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: EPS update result has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xf0 != 0 {
		return fmt.Errorf("%w: EPS update result spare bits are nonzero", ErrInvalidValue)
	}
	decoded := EPSUpdateResultFields{Result: EPSUpdateResult(payload[0])}
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}
