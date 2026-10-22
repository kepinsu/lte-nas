package ie

import (
	"fmt"
	"io"
)

// ConnType identifies the connectivity type, as defined in TS 24.301
// clause 9.9.4.2A.
type ConnType uint8

// Connectivity type values.
const (
	ConnTypeReserved  ConnType = 0x00
	ConnTypeNonIP     ConnType = 0x01
	ConnTypeEthernet  ConnType = 0x02
)

// ConnectivityTypeFields is the structured value of the tagged Connectivity type
// IE.
type ConnectivityTypeFields struct {
	Type ConnType
}

// NewConnectivityTypeWithFields creates a tagged Connectivity type IE from
// fields.
func NewConnectivityTypeWithFields(fields *ConnectivityTypeFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewConnectivityType(payload[0]), nil
}

// ConnectivityType returns the IE as structured fields.
func (i *IE) ConnectivityType() (*ConnectivityTypeFields, error) {
	if err := i.requireDefinition(ConnectivityTypeDefinition); err != nil {
		return nil, err
	}
	return ParseConnectivityTypeFields(i.Payload)
}

// ParseConnectivityTypeFields decodes a Connectivity type value without its IEI.
func ParseConnectivityTypeFields(payload []byte) (*ConnectivityTypeFields, error) {
	fields := &ConnectivityTypeFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Type is a defined Connectivity type value.
func (f *ConnectivityTypeFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: Connectivity type fields are nil", ErrInvalidValue)
	}
	if !f.Type.valid() {
		return fmt.Errorf("%w: Connectivity type %#02x is reserved", ErrInvalidValue, f.Type)
	}
	return nil
}

func (t ConnType) valid() bool {
	switch t {
	case ConnTypeNonIP, ConnTypeEthernet:
		return true
	default:
		return false
	}
}

// Marshal encodes fields into a Connectivity type value without its IEI.
func (f *ConnectivityTypeFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.Type)}, nil
}

// MarshalTo encodes fields into b.
func (f *ConnectivityTypeFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	b[0] = byte(f.Type)
	return nil
}

// MarshalLen returns the size of the Connectivity type value payload.
func (f *ConnectivityTypeFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a Connectivity type value without its IEI.
func (f *ConnectivityTypeFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: Connectivity type fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: Connectivity type has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0] != 0x01 && payload[0] != 0x02 {
		return fmt.Errorf("%w: Connectivity type is reserved", ErrInvalidValue)
	}
	decoded := ConnectivityTypeFields{Type: ConnType(payload[0])}
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}
