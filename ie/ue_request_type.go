package ie

import (
	"fmt"
	"io"
)

// UeReqType identifies the UE request type, as defined in TS 24.301 clause
// 9.9.3.65.
type UeReqType uint8

// UE request type values.
const (
	UeReqTypeReserved          UeReqType = 0x00
	UeReqTypeMobileOriginating UeReqType = 0x01
	UeReqTypeMobileTerminating UeReqType = 0x02
	UeReqTypeReserved2         UeReqType = 0x03
)

// UERequestTypeFields is the structured value of the tagged UE request type IE.
type UERequestTypeFields struct {
	Type UeReqType
}

// NewUERequestTypeWithFields creates a tagged UE request type IE from fields.
func NewUERequestTypeWithFields(fields *UERequestTypeFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewUERequestType(payload[0]), nil
}

// UERequestType returns the IE as structured fields.
func (i *IE) UERequestType() (*UERequestTypeFields, error) {
	if err := i.requireDefinition(UERequestTypeDefinition); err != nil {
		return nil, err
	}
	return ParseUERequestTypeFields(i.Payload)
}

// ParseUERequestTypeFields decodes a UE request type value without its IEI.
func ParseUERequestTypeFields(payload []byte) (*UERequestTypeFields, error) {
	fields := &UERequestTypeFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Type is a defined UE request type value.
func (f *UERequestTypeFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: UE request type fields are nil", ErrInvalidValue)
	}
	if !f.Type.valid() {
		return fmt.Errorf("%w: UE request type %#02x is reserved", ErrInvalidValue, f.Type)
	}
	return nil
}

func (t UeReqType) valid() bool {
	switch t {
	case UeReqTypeMobileOriginating, UeReqTypeMobileTerminating:
		return true
	default:
		return false
	}
}

// Marshal encodes fields into a UE request type value without its IEI.
func (f *UERequestTypeFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.Type)}, nil
}

// MarshalTo encodes fields into b.
func (f *UERequestTypeFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	b[0] = byte(f.Type)
	return nil
}

// MarshalLen returns the size of the UE request type value payload.
func (f *UERequestTypeFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a UE request type value without its IEI.
func (f *UERequestTypeFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: UE request type fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: UE request type has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0] != 0x01 && payload[0] != 0x02 {
		return fmt.Errorf("%w: UE request type is reserved", ErrInvalidValue)
	}
	decoded := UERequestTypeFields{Type: UeReqType(payload[0])}
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}
