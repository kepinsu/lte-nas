package ie

import (
	"fmt"
	"io"
)

// EPSAttachResult identifies the result of the EPS attach procedure, as
// defined in TS 24.301 clause 9.9.3.10.
type EPSAttachResult uint8

// EPS attach result values.
const (
	EPSAttachResultEPSOnly            EPSAttachResult = 0x01
	EPSAttachResultCombinedEPSAndIMSI EPSAttachResult = 0x02
)

// EPSAttachResultFields is the structured value of the positional EPS attach
// result type-1 IE.
type EPSAttachResultFields struct {
	Result EPSAttachResult
}

// NewPositionalEPSAttachResultWithFields creates a positional EPS attach
// result IE from fields.
func NewPositionalEPSAttachResultWithFields(fields *EPSAttachResultFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalEPSAttachResult(payload[0]), nil
}

// EPSAttachResult returns the EPS attach result IE as structured fields.
func (i *IE) EPSAttachResult() (*EPSAttachResultFields, error) {
	if err := i.requireDefinition(EPSAttachResultDefinition); err != nil {
		return nil, err
	}
	return ParseEPSAttachResultFields(i.Payload)
}

// ParseEPSAttachResultFields decodes an EPS attach result value without its
// IEI.
func ParseEPSAttachResultFields(payload []byte) (*EPSAttachResultFields, error) {
	fields := &EPSAttachResultFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Result is defined by TS 24.301.
func (f *EPSAttachResultFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: EPS attach result fields are nil", ErrInvalidValue)
	}
	if f.Result != EPSAttachResultEPSOnly && f.Result != EPSAttachResultCombinedEPSAndIMSI {
		return fmt.Errorf("%w: EPS attach result %#02x is reserved", ErrInvalidValue, f.Result)
	}
	return nil
}

// Marshal encodes fields into an EPS attach result value without its IEI.
func (f *EPSAttachResultFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.Result)}, nil
}

// MarshalTo encodes fields into b.
func (f *EPSAttachResultFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	b[0] = byte(f.Result)
	return nil
}

// MarshalLen returns the size of the EPS attach result value payload.
func (f *EPSAttachResultFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes an EPS attach result value without its IEI.
func (f *EPSAttachResultFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: EPS attach result fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: EPS attach result has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xf0 != 0 {
		return fmt.Errorf("%w: EPS attach result spare bits are nonzero", ErrInvalidValue)
	}
	decoded := EPSAttachResultFields{Result: EPSAttachResult(payload[0])}
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}
