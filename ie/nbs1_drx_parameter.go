package ie

import (
	"fmt"
	"io"
)

// Nbs1DrxParam identifies the NB-S1 DRX parameter, as defined in TS 24.301
// clause 9.9.3.63.
type Nbs1DrxParam uint8

// NB-S1 DRX parameter values.
const (
	Nbs1DrxParamReserved  Nbs1DrxParam = 0x00
	Nbs1DrxParamSupported   Nbs1DrxParam = 0x01
	Nbs1DrxParamNotApplicable Nbs1DrxParam = 0x02
)

// NBS1DRXParameterFields is the structured value of the tagged NB-S1 DRX
// parameter IE.
type NBS1DRXParameterFields struct {
	Type Nbs1DrxParam
}

// NewNBS1DRXParameterWithFields creates a tagged NB-S1 DRX parameter IE from
// fields.
func NewNBS1DRXParameterWithFields(fields *NBS1DRXParameterFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewNBS1DRXParameter(payload[0]), nil
}

// NBS1DRXParameter returns the IE as structured fields.
func (i *IE) NBS1DRXParameter() (*NBS1DRXParameterFields, error) {
	if err := i.requireDefinition(NBS1DRXParameterDefinition); err != nil {
		return nil, err
	}
	return ParseNBS1DRXParameterFields(i.Payload)
}

// ParseNBS1DRXParameterFields decodes an NB-S1 DRX parameter value without its
// IEI.
func ParseNBS1DRXParameterFields(payload []byte) (*NBS1DRXParameterFields, error) {
	fields := &NBS1DRXParameterFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Type is a defined NB-S1 DRX parameter value.
func (f *NBS1DRXParameterFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: NB-S1 DRX parameter fields are nil", ErrInvalidValue)
	}
	if !f.Type.valid() {
		return fmt.Errorf("%w: NB-S1 DRX parameter %#02x is reserved", ErrInvalidValue, f.Type)
	}
	return nil
}

func (t Nbs1DrxParam) valid() bool {
	switch t {
	case Nbs1DrxParamSupported, Nbs1DrxParamNotApplicable:
		return true
	default:
		return false
	}
}

// Marshal encodes fields into an NB-S1 DRX parameter value without its IEI.
func (f *NBS1DRXParameterFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.Type)}, nil
}

// MarshalTo encodes fields into b.
func (f *NBS1DRXParameterFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	b[0] = byte(f.Type)
	return nil
}

// MarshalLen returns the size of the NB-S1 DRX parameter value payload.
func (f *NBS1DRXParameterFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes an NB-S1 DRX parameter value without its IEI.
func (f *NBS1DRXParameterFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: NB-S1 DRX parameter fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: NB-S1 DRX parameter has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0] != 0x01 && payload[0] != 0x02 {
		return fmt.Errorf("%w: NB-S1 DRX parameter is reserved", ErrInvalidValue)
	}
	decoded := NBS1DRXParameterFields{Type: Nbs1DrxParam(payload[0])}
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}
