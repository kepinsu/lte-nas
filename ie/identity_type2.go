package ie

import (
	"fmt"
	"io"
)

// IdentityType2 identifies the type of identity in an Identity Type 2 IE, as
// defined in TS 24.301 clause 9.9.3.17.
type IdentityType2 uint8

// Identity Type 2 values.
const (
	IdentityType2IMEI     IdentityType2 = 0x01
	IdentityType2IMEISV   IdentityType2 = 0x02
	IdentityType2TMSI     IdentityType2 = 0x03
	IdentityType2GUTI     IdentityType2 = 0x04
	IdentityType2MSISDN   IdentityType2 = 0x05
	IdentityType2Reserved IdentityType2 = 0x00
)

// IdentityType2Fields is the structured value of the positional Identity type 2
// IE.
type IdentityType2Fields struct {
	Type IdentityType2
}

// NewPositionalIdentityType2WithFields creates a positional Identity type 2 IE
// from fields.
func NewPositionalIdentityType2WithFields(fields *IdentityType2Fields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalIdentityType2(payload[0]), nil
}

// IdentityType2 returns the IE as structured fields.
func (i *IE) IdentityType2() (*IdentityType2Fields, error) {
	if err := i.requireDefinition(IdentityType2Definition); err != nil {
		return nil, err
	}
	return ParseIdentityType2Fields(i.Payload)
}

// ParseIdentityType2Fields decodes an Identity type 2 value without its IEI.
func ParseIdentityType2Fields(payload []byte) (*IdentityType2Fields, error) {
	fields := &IdentityType2Fields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Type is a defined Identity Type 2 value.
func (f *IdentityType2Fields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: Identity type 2 fields are nil", ErrInvalidValue)
	}
	if !f.Type.valid() {
		return fmt.Errorf("%w: Identity type 2 %#02x is reserved", ErrInvalidValue, f.Type)
	}
	return nil
}

func (t IdentityType2) valid() bool {
	switch t {
	case IdentityType2IMEI, IdentityType2IMEISV, IdentityType2TMSI,
		IdentityType2GUTI, IdentityType2MSISDN:
		return true
	default:
		return false
	}
}

// Marshal encodes fields into an Identity type 2 value without its IEI.
func (f *IdentityType2Fields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.Type)}, nil
}

// MarshalTo encodes fields into b.
func (f *IdentityType2Fields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	b[0] = byte(f.Type)
	return nil
}

// MarshalLen returns the size of the Identity type 2 value payload.
func (f *IdentityType2Fields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes an Identity type 2 value without its IEI.
func (f *IdentityType2Fields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: Identity type 2 fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: Identity type 2 has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0] == 0 {
		return fmt.Errorf("%w: Identity type 2 is reserved", ErrInvalidValue)
	}
	decoded := IdentityType2Fields{Type: IdentityType2(payload[0])}
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}
