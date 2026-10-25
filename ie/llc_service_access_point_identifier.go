package ie

import (
	"fmt"
	"io"
)

// LlcSapi identifies the LLC service access point, as defined in TS 24.301
// clause 9.9.4.7.
type LlcSapi uint8

// LLC service access point identifier values.
const (
	LlcSapiReserved LlcSapi = 0x00
	LlcSapiLLC      LlcSapi = 0x01
)

// LLCServiceAccessPointIdentifierFields is the structured value of the tagged
// LLC service access point identifier IE.
type LLCServiceAccessPointIdentifierFields struct {
	Type LlcSapi
}

// NewLLCServiceAccessPointIdentifierWithFields creates a tagged LLC service
// access point identifier IE from fields.
func NewLLCServiceAccessPointIdentifierWithFields(fields *LLCServiceAccessPointIdentifierFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewLLCServiceAccessPointIdentifier(payload[0]), nil
}

// LLCServiceAccessPointIdentifier returns the IE as structured fields.
func (i *IE) LLCServiceAccessPointIdentifier() (*LLCServiceAccessPointIdentifierFields, error) {
	if err := i.requireDefinition(LLCServiceAccessPointIdentifierDefinition); err != nil {
		return nil, err
	}
	return ParseLLCServiceAccessPointIdentifierFields(i.Payload)
}

// ParseLLCServiceAccessPointIdentifierFields decodes an LLC service access
// point identifier value without its IEI.
func ParseLLCServiceAccessPointIdentifierFields(payload []byte) (*LLCServiceAccessPointIdentifierFields, error) {
	fields := &LLCServiceAccessPointIdentifierFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Type is a defined LLC service access point
// identifier value.
func (f *LLCServiceAccessPointIdentifierFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: LLC service access point identifier fields are nil", ErrInvalidValue)
	}
	if !f.Type.valid() {
		return fmt.Errorf("%w: LLC service access point identifier %#02x is reserved", ErrInvalidValue, f.Type)
	}
	return nil
}

func (t LlcSapi) valid() bool {
	switch t {
	case LlcSapiReserved, LlcSapiLLC:
		return true
	default:
		return false
	}
}

// Marshal encodes fields into an LLC service access point identifier value
// without its IEI.
func (f *LLCServiceAccessPointIdentifierFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.Type)}, nil
}

// MarshalTo encodes fields into b.
func (f *LLCServiceAccessPointIdentifierFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	b[0] = byte(f.Type)
	return nil
}

// MarshalLen returns the size of the LLC service access point identifier value
// payload.
func (f *LLCServiceAccessPointIdentifierFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes an LLC service access point identifier value without
// its IEI.
func (f *LLCServiceAccessPointIdentifierFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: LLC service access point identifier fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: LLC service access point identifier has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	decoded := LLCServiceAccessPointIdentifierFields{Type: LlcSapi(payload[0])}
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}
