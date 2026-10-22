package ie

import (
	"fmt"
	"io"
)

// ServiceType identifies the service type, as defined in TS 24.301 clause
// 9.9.3.27.
type ServiceType uint8

// Service type values.
const (
	ServiceTypeVoice  ServiceType = 0x01
	ServiceTypeSS     ServiceType = 0x02
	ServiceTypeSMS    ServiceType = 0x03
)

// ServiceTypeFields is the structured value of the positional Service type IE
// (TS 24.301 clause 9.9.3.27).
type ServiceTypeFields struct {
	Type ServiceType
}

// NewPositionalServiceTypeWithFields creates a positional Service type IE from
// fields.
func NewPositionalServiceTypeWithFields(fields *ServiceTypeFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalServiceType(payload[0]), nil
}

// ServiceType returns the IE as structured fields.
func (i *IE) ServiceType() (*ServiceTypeFields, error) {
	if err := i.requireDefinition(ServiceTypeDefinition); err != nil {
		return nil, err
	}
	return ParseServiceTypeFields(i.Payload)
}

// ParseServiceTypeFields decodes a service type value without its IEI.
func ParseServiceTypeFields(payload []byte) (*ServiceTypeFields, error) {
	fields := &ServiceTypeFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Type is a defined service type value.
func (f *ServiceTypeFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: service type fields are nil", ErrInvalidValue)
	}
	if !f.Type.valid() {
		return fmt.Errorf("%w: service type %#02x is reserved", ErrInvalidValue, f.Type)
	}
	return nil
}

func (t ServiceType) valid() bool {
	switch t {
	case ServiceTypeVoice, ServiceTypeSS, ServiceTypeSMS:
		return true
	default:
		return false
	}
}

// Marshal encodes fields into a service type value without its IEI.
func (f *ServiceTypeFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.Type)}, nil
}

// MarshalTo encodes fields into b.
func (f *ServiceTypeFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	b[0] = byte(f.Type)
	return nil
}

// MarshalLen returns the size of the service type value payload.
func (f *ServiceTypeFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a service type value without its IEI. It leaves f
// unchanged when the payload is invalid.
func (f *ServiceTypeFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: service type fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: service type has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	var decoded ServiceTypeFields
	decoded.Type = ServiceType(payload[0])
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}
