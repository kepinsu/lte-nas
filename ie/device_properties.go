package ie

import (
	"fmt"
	"io"
)

// DevicePropertiesFields is the structured value of the Device properties IE (TS
// 24.301 clause 9.9.2.0A, see TS 24.008 clause 10.5.7.8).
type DevicePropertiesFields struct {
	// EmergencyCallSupported indicates whether emergency calls are supported.
	EmergencyCallSupported bool
}

// NewDevicePropertiesWithFields creates a tagged Device properties IE from
// fields.
func NewDevicePropertiesWithFields(fields *DevicePropertiesFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewDeviceProperties(payload[0]), nil
}

// DeviceProperties returns the IE as structured fields.
func (i *IE) DeviceProperties() (*DevicePropertiesFields, error) {
	if err := i.requireDefinition(DevicePropertiesDefinition); err != nil {
		return nil, err
	}
	return ParseDevicePropertiesFields(i.Payload)
}

// ParseDevicePropertiesFields decodes a Device properties value without its IEI.
func ParseDevicePropertiesFields(payload []byte) (*DevicePropertiesFields, error) {
	fields := &DevicePropertiesFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *DevicePropertiesFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: device properties fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a Device properties value, without its IEI.
func (f *DevicePropertiesFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	var encoded byte
	if f.EmergencyCallSupported {
		encoded |= 0x80
	}
	return []byte{encoded}, nil
}

// MarshalTo encodes fields into b.
func (f *DevicePropertiesFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the Device properties value payload.
func (f *DevicePropertiesFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a Device properties value without its IEI.
func (f *DevicePropertiesFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: device properties fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: device properties has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0x7f != 0 {
		return fmt.Errorf("%w: device properties spare bits are nonzero", ErrInvalidValue)
	}
	*f = DevicePropertiesFields{
		EmergencyCallSupported: payload[0]&0x80 != 0,
	}
	return nil
}
