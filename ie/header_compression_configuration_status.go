package ie

import (
	"fmt"
	"io"
)

// HeaderCompressionConfigurationStatusFields is the structured value of the
// header compression configuration status IE (TS 24.301 clause 9.9.4.27).
// The two-octet payload encodes a bitmap indicating which header compression
// protocols have their configuration updated.
type HeaderCompressionConfigurationStatusFields struct {
	// Payload contains the 2-octet bitmap. Bit 0 of the first octet indicates
	// ROHC status; bit 0 of the second octet indicates Robust header compression.
	Payload [2]byte
}

// NewHeaderCompressionConfigurationStatusWithFields creates a tagged header
// compression configuration status IE from fields.
func NewHeaderCompressionConfigurationStatusWithFields(fields *HeaderCompressionConfigurationStatusFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewHeaderCompressionConfigurationStatus(payload), nil
}

// NewPositionalHeaderCompressionConfigurationStatusWithFields creates a positional
// header compression configuration status IE from fields.
func NewPositionalHeaderCompressionConfigurationStatusWithFields(fields *HeaderCompressionConfigurationStatusFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalHeaderCompressionConfigurationStatus(payload), nil
}

// HeaderCompressionConfigurationStatus returns the IE as structured fields.
func (i *IE) HeaderCompressionConfigurationStatus() (*HeaderCompressionConfigurationStatusFields, error) {
	if err := i.requireDefinition(HeaderCompressionConfigurationStatusDefinition); err != nil {
		return nil, err
	}
	return ParseHeaderCompressionConfigurationStatusFields(i.Payload)
}

// ParseHeaderCompressionConfigurationStatusFields decodes a header compression
// configuration status value without its IEI and length octet.
func ParseHeaderCompressionConfigurationStatusFields(payload []byte) (*HeaderCompressionConfigurationStatusFields, error) {
	fields := &HeaderCompressionConfigurationStatusFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *HeaderCompressionConfigurationStatusFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: header compression configuration status fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a 2-octet header compression configuration status value.
func (f *HeaderCompressionConfigurationStatusFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Payload))
	copy(payload, f.Payload[:])
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *HeaderCompressionConfigurationStatusFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 2 {
		return io.ErrShortBuffer
	}
	copy(b, f.Payload[:])
	return nil
}

// MarshalLen returns the size of the header compression configuration status value payload.
func (f *HeaderCompressionConfigurationStatusFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 2
}

// UnmarshalBinary decodes a 2-octet header compression configuration status value.
func (f *HeaderCompressionConfigurationStatusFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: header compression configuration status fields are nil", ErrInvalidValue)
	}
	if len(payload) != 2 {
		return fmt.Errorf("%w: header compression configuration status has %d value octets, want 2", ErrInvalidLength, len(payload))
	}
	var decoded HeaderCompressionConfigurationStatusFields
	copy(decoded.Payload[:], payload)
	*f = decoded
	return nil
}
