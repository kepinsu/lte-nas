package ie

import (
	"fmt"
	"io"
)

// TrafficFlowAggregateDescriptionFields is the structured value of the positional
// Traffic flow aggregate description IE (TS 24.301 clause 9.9.4.15). It carries an
// opaque byte string between 1 and 255 octets.
type TrafficFlowAggregateDescriptionFields struct {
	Value []byte
}

// NewPositionalTrafficFlowAggregateDescriptionWithFields creates a positional Traffic
// flow aggregate description IE from fields.
func NewPositionalTrafficFlowAggregateDescriptionWithFields(fields *TrafficFlowAggregateDescriptionFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalTrafficFlowAggregateDescription(payload), nil
}

// TrafficFlowAggregateDescription returns the IE as structured fields.
func (i *IE) TrafficFlowAggregateDescription() (*TrafficFlowAggregateDescriptionFields, error) {
	if err := i.requireDefinition(TrafficFlowAggregateDescriptionDefinition); err != nil {
		return nil, err
	}
	return ParseTrafficFlowAggregateDescriptionFields(i.Payload)
}

// ParseTrafficFlowAggregateDescriptionFields decodes a traffic flow aggregate
// description value without its IEI and length octet.
func ParseTrafficFlowAggregateDescriptionFields(payload []byte) (*TrafficFlowAggregateDescriptionFields, error) {
	fields := &TrafficFlowAggregateDescriptionFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the payload is within [1, 255] octets.
func (f *TrafficFlowAggregateDescriptionFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: traffic flow aggregate description fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 1 {
		return fmt.Errorf("%w: traffic flow aggregate description has %d value octets, want >= 1", ErrInvalidLength, n)
	}
	if n > 255 {
		return fmt.Errorf("%w: traffic flow aggregate description has %d value octets, want <= 255", ErrInvalidLength, n)
	}
	return nil
}

// Marshal encodes fields into a traffic flow aggregate description value without its
// IEI and length octet.
func (f *TrafficFlowAggregateDescriptionFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *TrafficFlowAggregateDescriptionFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the traffic flow aggregate description value payload,
// or zero for a nil receiver or nil Value.
func (f *TrafficFlowAggregateDescriptionFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes a traffic flow aggregate description value without its IEI and
// length octet. It copies the payload into a temporary struct first, leaving f unchanged
// on failure.
func (f *TrafficFlowAggregateDescriptionFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: traffic flow aggregate description fields are nil", ErrInvalidValue)
	}
	if len(payload) < 1 {
		return fmt.Errorf("%w: traffic flow aggregate description has %d value octets, want >= 1", ErrInvalidLength, len(payload))
	}
	if len(payload) > 255 {
		return fmt.Errorf("%w: traffic flow aggregate description has %d value octets, want <= 255", ErrInvalidLength, len(payload))
	}
	var decoded TrafficFlowAggregateDescriptionFields
	decoded.Value = make([]byte, len(payload))
	copy(decoded.Value, payload)
	*f = decoded
	return nil
}
