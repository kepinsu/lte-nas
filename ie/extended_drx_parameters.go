package ie

import (
	"fmt"
	"io"
)

// ExtendedDRXParametersFields is the structured value of the Extended DRX
// parameters IE (TS 24.301 clause 9.9.3.46, see TS 24.008 clause 10.5.5.32).
type ExtendedDRXParametersFields struct {
	// PFC indicates the preferred CIoT EPS optimization.
	PFC bool
}

// NewExtendedDRXParametersWithFields creates a tagged Extended DRX parameters IE
// from fields.
func NewExtendedDRXParametersWithFields(fields *ExtendedDRXParametersFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewExtendedDRXParameters(payload[0]), nil
}

// ExtendedDRXParameters returns the IE as structured fields.
func (i *IE) ExtendedDRXParameters() (*ExtendedDRXParametersFields, error) {
	if err := i.requireDefinition(ExtendedDRXParametersDefinition); err != nil {
		return nil, err
	}
	return ParseExtendedDRXParametersFields(i.Payload)
}

// ParseExtendedDRXParametersFields decodes an Extended DRX parameters value
// without its IEI.
func ParseExtendedDRXParametersFields(payload []byte) (*ExtendedDRXParametersFields, error) {
	fields := &ExtendedDRXParametersFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *ExtendedDRXParametersFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: extended DRX parameters fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into an Extended DRX parameters value, without its IEI.
func (f *ExtendedDRXParametersFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	var encoded byte
	if f.PFC {
		encoded |= 0x80
	}
	return []byte{encoded}, nil
}

// MarshalTo encodes fields into b.
func (f *ExtendedDRXParametersFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the Extended DRX parameters value payload.
func (f *ExtendedDRXParametersFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes an Extended DRX parameters value without its IEI.
func (f *ExtendedDRXParametersFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: extended DRX parameters fields are nil", ErrInvalidValue)
	}
	if len(payload) == 0 {
		return io.ErrShortBuffer
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: extended DRX parameters has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	*f = ExtendedDRXParametersFields{
		PFC: payload[0]&0x80 != 0,
	}
	return nil
}
