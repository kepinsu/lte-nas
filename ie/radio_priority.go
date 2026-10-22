package ie

import (
	"fmt"
	"io"
)

// RadioPri identifies the radio priority level, as defined in TS 24.301
// clause 9.9.4.13.
type RadioPri uint8

// Radio priority values.
const (
	RadioPri15      RadioPri = 0x0f
	RadioPri14      RadioPri = 0x0e
	RadioPri13      RadioPri = 0x0d
	RadioPri12      RadioPri = 0x0c
	RadioPri11      RadioPri = 0x0b
	RadioPri10      RadioPri = 0x0a
	RadioPri9       RadioPri = 0x09
	RadioPri8       RadioPri = 0x08
	RadioPri7       RadioPri = 0x07
	RadioPri6       RadioPri = 0x06
	RadioPri5       RadioPri = 0x05
	RadioPri4       RadioPri = 0x04
	RadioPri3       RadioPri = 0x03
	RadioPri2       RadioPri = 0x02
	RadioPri1       RadioPri = 0x01
	RadioPriDefault RadioPri = 0x00
)

// RadioPriorityFields is the structured value of the positional Radio priority
// IE.
type RadioPriorityFields struct {
	Type RadioPri
}

// NewPositionalRadioPriorityWithFields creates a positional Radio priority IE
// from fields.
func NewPositionalRadioPriorityWithFields(fields *RadioPriorityFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalRadioPriority(payload[0]), nil
}

// RadioPriority returns the IE as structured fields.
func (i *IE) RadioPriority() (*RadioPriorityFields, error) {
	if err := i.requireDefinition(RadioPriorityDefinition); err != nil {
		return nil, err
	}
	return ParseRadioPriorityFields(i.Payload)
}

// ParseRadioPriorityFields decodes a Radio priority value without its IEI.
func ParseRadioPriorityFields(payload []byte) (*RadioPriorityFields, error) {
	fields := &RadioPriorityFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Type is a defined Radio priority value.
func (f *RadioPriorityFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: Radio priority fields are nil", ErrInvalidValue)
	}
	if !f.Type.valid() {
		return fmt.Errorf("%w: Radio priority %#02x is reserved", ErrInvalidValue, f.Type)
	}
	return nil
}

func (t RadioPri) valid() bool {
	switch t {
	case RadioPri1, RadioPri2, RadioPri3, RadioPri4,
		RadioPri5, RadioPri6, RadioPri7, RadioPri8,
		RadioPri9, RadioPri10, RadioPri11, RadioPri12,
		RadioPri13, RadioPri14, RadioPri15, RadioPriDefault:
		return true
	default:
		return false
	}
}

// Marshal encodes fields into a Radio priority value without its IEI.
func (f *RadioPriorityFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.Type)}, nil
}

// MarshalTo encodes fields into b.
func (f *RadioPriorityFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	b[0] = byte(f.Type)
	return nil
}

// MarshalLen returns the size of the Radio priority value payload.
func (f *RadioPriorityFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a Radio priority value without its IEI.
func (f *RadioPriorityFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: Radio priority fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: Radio priority has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	decoded := RadioPriorityFields{Type: RadioPri(payload[0])}
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}
