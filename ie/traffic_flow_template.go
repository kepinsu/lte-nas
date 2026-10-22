package ie

import (
	"fmt"
	"io"
)

// TrafficFlowTemplateFields is the structured value of the traffic flow
// template IE (TS 24.301 clause 9.9.4.16). It carries an opaque byte string
// between 1 and 255 octets encoding filter information for a traffic flow.
type TrafficFlowTemplateFields struct {
	Value []byte
}

// NewTrafficFlowTemplateWithFields creates a tagged Traffic flow template IE
// from fields.
func NewTrafficFlowTemplateWithFields(fields *TrafficFlowTemplateFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewTrafficFlowTemplate(payload), nil
}

// NewPositionalTrafficFlowTemplateWithFields creates a positional Traffic flow
// template IE from fields.
func NewPositionalTrafficFlowTemplateWithFields(fields *TrafficFlowTemplateFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalTrafficFlowTemplate(payload), nil
}

// TrafficFlowTemplate returns the IE as structured fields.
func (i *IE) TrafficFlowTemplate() (*TrafficFlowTemplateFields, error) {
	if err := i.requireDefinition(TrafficFlowTemplateDefinition); err != nil {
		return nil, err
	}
	return ParseTrafficFlowTemplateFields(i.Payload)
}

// ParseTrafficFlowTemplateFields decodes a traffic flow template value without
// its IEI and length octet.
func ParseTrafficFlowTemplateFields(payload []byte) (*TrafficFlowTemplateFields, error) {
	fields := &TrafficFlowTemplateFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the payload is within [1, 255] octets.
func (f *TrafficFlowTemplateFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: traffic flow template fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 1 {
		return fmt.Errorf("%w: traffic flow template has %d value octets, want >= 1", ErrInvalidLength, n)
	}
	if n > 255 {
		return fmt.Errorf("%w: traffic flow template has %d value octets, want <= 255", ErrInvalidLength, n)
	}
	return nil
}

// Marshal encodes fields into a traffic flow template value without its IEI and
// length octet.
func (f *TrafficFlowTemplateFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *TrafficFlowTemplateFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the traffic flow template value payload, or
// zero for a nil receiver or nil Value.
func (f *TrafficFlowTemplateFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes a traffic flow template value without its IEI and
// length octet. It copies the payload into a temporary struct first, leaving f
// unchanged on failure.
func (f *TrafficFlowTemplateFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: traffic flow template fields are nil", ErrInvalidValue)
	}
	if len(payload) < 1 {
		return fmt.Errorf("%w: traffic flow template has %d value octets, want >= 1", ErrInvalidLength, len(payload))
	}
	if len(payload) > 255 {
		return fmt.Errorf("%w: traffic flow template has %d value octets, want <= 255", ErrInvalidLength, len(payload))
	}
	var decoded TrafficFlowTemplateFields
	decoded.Value = make([]byte, len(payload))
	copy(decoded.Value, payload)
	*f = decoded
	return nil
}
