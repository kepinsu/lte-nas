package ie

import (
	"fmt"
	"io"
)

// ReleaseAssistanceIndicationFields is the structured value of the Release
// assistance indication IE (TS 24.301 clause 9.9.4.25).
type ReleaseAssistanceIndicationFields struct {
	// DDX indicates whether downlink data is expected.
	DDX uint8
}

// NewReleaseAssistanceIndicationWithFields creates a tagged Release assistance
// indication IE from fields.
func NewReleaseAssistanceIndicationWithFields(fields *ReleaseAssistanceIndicationFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewReleaseAssistanceIndication(payload[0]), nil
}

// ReleaseAssistanceIndication returns the IE as structured fields.
func (i *IE) ReleaseAssistanceIndication() (*ReleaseAssistanceIndicationFields, error) {
	if err := i.requireDefinition(ReleaseAssistanceIndicationDefinition); err != nil {
		return nil, err
	}
	return ParseReleaseAssistanceIndicationFields(i.Payload)
}

// ParseReleaseAssistanceIndicationFields decodes a Release assistance indication
// value without its IEI.
func ParseReleaseAssistanceIndicationFields(payload []byte) (*ReleaseAssistanceIndicationFields, error) {
	fields := &ReleaseAssistanceIndicationFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that DDX is a defined value.
func (f *ReleaseAssistanceIndicationFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: release assistance indication fields are nil", ErrInvalidValue)
	}
	if f.DDX > 0x03 {
		return fmt.Errorf("%w: DDX value %d is reserved", ErrInvalidValue, f.DDX)
	}
	return nil
}

// Marshal encodes fields into a Release assistance indication value, without its
// IEI.
func (f *ReleaseAssistanceIndicationFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.DDX)}, nil
}

// MarshalTo encodes fields into b.
func (f *ReleaseAssistanceIndicationFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the Release assistance indication value payload.
func (f *ReleaseAssistanceIndicationFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a Release assistance indication value without its IEI.
func (f *ReleaseAssistanceIndicationFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: release assistance indication fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: release assistance indication has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xfc != 0 {
		return fmt.Errorf("%w: release assistance indication spare bits are nonzero", ErrInvalidValue)
	}
	v := ReleaseAssistanceIndicationFields{DDX: payload[0] & 0x03}
	if v.DDX > 0x03 {
		return fmt.Errorf("%w: DDX value %d is reserved", ErrInvalidValue, v.DDX)
	}
	*f = v
	return nil
}
