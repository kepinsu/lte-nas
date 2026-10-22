package ie

import (
	"fmt"
	"io"
)

// DRXParameterFields is the structured value of the DRX parameter IE
// (TS 24.301 clause 9.9.3.8). The DRX parameter encodes DRX cycle length and
// related configuration in a single octet.
type DRXParameterFields struct {
	DRXValue [1]byte
}

// NewDRXParameterWithFields creates a tagged DRX parameter IE from fields.
func NewDRXParameterWithFields(fields *DRXParameterFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewDRXParameter(payload), nil
}

// NewPositionalDRXParameterWithFields creates a positional DRX parameter IE
// from fields.
func NewPositionalDRXParameterWithFields(fields *DRXParameterFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalDRXParameter(payload), nil
}

// DRXParameter returns the IE as structured fields.
func (i *IE) DRXParameter() (*DRXParameterFields, error) {
	if err := i.requireDefinition(DRXParameterDefinition); err != nil {
		return nil, err
	}
	return ParseDRXParameterFields(i.Payload)
}

// ParseDRXParameterFields decodes a DRX parameter value without its IEI.
func ParseDRXParameterFields(payload []byte) (*DRXParameterFields, error) {
	fields := &DRXParameterFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *DRXParameterFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: DRX parameter fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a 1-octet DRX parameter value.
func (f *DRXParameterFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.DRXValue))
	copy(payload, f.DRXValue[:])
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *DRXParameterFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < f.MarshalLen() {
		return io.ErrShortBuffer
	}
	copy(b, f.DRXValue[:])
	return nil
}

// MarshalLen returns the size of the DRX parameter value payload.
func (f *DRXParameterFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a 1-octet DRX parameter value.
func (f *DRXParameterFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: DRX parameter fields are nil", ErrInvalidValue)
	}
	if len(payload) != f.MarshalLen() {
		return fmt.Errorf("%w: DRX parameter has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	var decoded DRXParameterFields
	copy(decoded.DRXValue[:], payload)
	*f = decoded
	return nil
}
