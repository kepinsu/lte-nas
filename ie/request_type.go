package ie

import (
	"fmt"
	"io"
)

// RequestType identifies the request type, as defined in TS 24.301 clause
// 9.9.4.14.
type RequestType uint8

// Request type values.
const (
	RequestTypeReserved  RequestType = 0x00
	RequestTypeInitial   RequestType = 0x01
	RequestTypeEmergency RequestType = 0x02
	RequestTypePeriodic  RequestType = 0x03
	RequestTypeMobilized RequestType = 0x04
	RequestTypeReserved2 RequestType = 0x05
)

// RequestTypeFields is the structured value of the positional Request type IE.
type RequestTypeFields struct {
	Type RequestType
}

// NewPositionalRequestTypeWithFields creates a positional Request type IE from
// fields.
func NewPositionalRequestTypeWithFields(fields *RequestTypeFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalRequestType(payload[0]), nil
}

// RequestType returns the IE as structured fields.
func (i *IE) RequestType() (*RequestTypeFields, error) {
	if err := i.requireDefinition(RequestTypeDefinition); err != nil {
		return nil, err
	}
	return ParseRequestTypeFields(i.Payload)
}

// ParseRequestTypeFields decodes a Request type value without its IEI.
func ParseRequestTypeFields(payload []byte) (*RequestTypeFields, error) {
	fields := &RequestTypeFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Type is a defined Request type value.
func (f *RequestTypeFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: Request type fields are nil", ErrInvalidValue)
	}
	if !f.Type.valid() {
		return fmt.Errorf("%w: Request type %#02x is reserved", ErrInvalidValue, f.Type)
	}
	return nil
}

func (t RequestType) valid() bool {
	switch t {
	case RequestTypeInitial, RequestTypeEmergency, RequestTypePeriodic,
		RequestTypeMobilized:
		return true
	default:
		return false
	}
}

// Marshal encodes fields into a Request type value without its IEI.
func (f *RequestTypeFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.Type)}, nil
}

// MarshalTo encodes fields into b.
func (f *RequestTypeFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	b[0] = byte(f.Type)
	return nil
}

// MarshalLen returns the size of the Request type value payload.
func (f *RequestTypeFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a Request type value without its IEI.
func (f *RequestTypeFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: Request type fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: Request type has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0] == 0 || payload[0] >= 0x05 {
		return fmt.Errorf("%w: Request type is reserved", ErrInvalidValue)
	}
	decoded := RequestTypeFields{Type: RequestType(payload[0])}
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}
