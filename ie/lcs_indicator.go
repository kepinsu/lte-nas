package ie

import (
	"fmt"
	"io"
)

// LcsInd identifies the LCS indicator, as defined in TS 24.301 clause
// 9.9.3.40.
type LcsInd uint8

// LCS indicator values.
const (
	LcsIndNormal    LcsInd = 0x00
	LcsIndEmergency LcsInd = 0x01
	LcsIndReserved  LcsInd = 0x02
)

// LCSIndicatorFields is the structured value of the tagged LCS indicator IE.
type LCSIndicatorFields struct {
	Type LcsInd
}

// NewLCSIndicatorWithFields creates a tagged LCS indicator IE from fields.
func NewLCSIndicatorWithFields(fields *LCSIndicatorFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewLCSIndicator(payload[0]), nil
}

// LCSIndicator returns the IE as structured fields.
func (i *IE) LCSIndicator() (*LCSIndicatorFields, error) {
	if err := i.requireDefinition(LCSIndicatorDefinition); err != nil {
		return nil, err
	}
	return ParseLCSIndicatorFields(i.Payload)
}

// ParseLCSIndicatorFields decodes an LCS indicator value without its IEI.
func ParseLCSIndicatorFields(payload []byte) (*LCSIndicatorFields, error) {
	fields := &LCSIndicatorFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Type is a defined LCS indicator value.
func (f *LCSIndicatorFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: LCS indicator fields are nil", ErrInvalidValue)
	}
	if !f.Type.valid() {
		return fmt.Errorf("%w: LCS indicator %#02x is reserved", ErrInvalidValue, f.Type)
	}
	return nil
}

func (t LcsInd) valid() bool {
	switch t {
	case LcsIndNormal, LcsIndEmergency:
		return true
	default:
		return false
	}
}

// Marshal encodes fields into an LCS indicator value without its IEI.
func (f *LCSIndicatorFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.Type)}, nil
}

// MarshalTo encodes fields into b.
func (f *LCSIndicatorFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	b[0] = byte(f.Type)
	return nil
}

// MarshalLen returns the size of the LCS indicator value payload.
func (f *LCSIndicatorFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes an LCS indicator value without its IEI.
func (f *LCSIndicatorFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: LCS indicator fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: LCS indicator has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	decoded := LCSIndicatorFields{Type: LcsInd(payload[0])}
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}
