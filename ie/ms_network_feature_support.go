package ie

import (
	"fmt"
	"io"
)

// MSNetworkFeatureSupportFields is the structured value of the MS network
// feature support IE (TS 24.301 clause 9.9.3.20A, see TS 24.008 clause 10.5.1.15).
type MSNetworkFeatureSupportFields struct {
	// T3412ExtendedValue indicates support for extended periodic timer value.
	T3412ExtendedValue bool
}

// NewMSNetworkFeatureSupportWithFields creates a tagged MS network feature
// support IE from fields.
func NewMSNetworkFeatureSupportWithFields(fields *MSNetworkFeatureSupportFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewMSNetworkFeatureSupport(payload[0]), nil
}

// MSNetworkFeatureSupport returns the IE as structured fields.
func (i *IE) MSNetworkFeatureSupport() (*MSNetworkFeatureSupportFields, error) {
	if err := i.requireDefinition(MSNetworkFeatureSupportDefinition); err != nil {
		return nil, err
	}
	return ParseMSNetworkFeatureSupportFields(i.Payload)
}

// ParseMSNetworkFeatureSupportFields decodes an MS network feature support value
// without its IEI.
func ParseMSNetworkFeatureSupportFields(payload []byte) (*MSNetworkFeatureSupportFields, error) {
	fields := &MSNetworkFeatureSupportFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *MSNetworkFeatureSupportFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: MS network feature support fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into an MS network feature support value, without its
// IEI.
func (f *MSNetworkFeatureSupportFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	var encoded byte
	if f.T3412ExtendedValue {
		encoded |= 0x80
	}
	return []byte{encoded}, nil
}

// MarshalTo encodes fields into b.
func (f *MSNetworkFeatureSupportFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the MS network feature support value payload.
func (f *MSNetworkFeatureSupportFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes an MS network feature support value without its IEI.
func (f *MSNetworkFeatureSupportFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: MS network feature support fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: MS network feature support has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0x7f != 0 {
		return fmt.Errorf("%w: MS network feature support spare bits are nonzero", ErrInvalidValue)
	}
	*f = MSNetworkFeatureSupportFields{
		T3412ExtendedValue: payload[0]&0x80 != 0,
	}
	return nil
}
