package ie

import (
	"fmt"
	"io"
)

// MobileStationClassmark3Fields is the structured value of the mobile station
// classmark 3 IE (TS 24.301 clause 9.9.2.5). It carries an opaque byte string
// between 0 and 32 octets.
type MobileStationClassmark3Fields struct {
	Value []byte
}

// NewMobileStationClassmark3WithFields creates a tagged Mobile station classmark
// 3 IE from fields.
func NewMobileStationClassmark3WithFields(fields *MobileStationClassmark3Fields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewMobileStationClassmark3(payload), nil
}

// NewPositionalMobileStationClassmark3WithFields creates a positional Mobile station
// classmark 3 IE from fields.
func NewPositionalMobileStationClassmark3WithFields(fields *MobileStationClassmark3Fields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalMobileStationClassmark3(payload), nil
}

// MobileStationClassmark3 returns the IE as structured fields.
func (i *IE) MobileStationClassmark3() (*MobileStationClassmark3Fields, error) {
	if err := i.requireDefinition(MobileStationClassmark3Definition); err != nil {
		return nil, err
	}
	return ParseMobileStationClassmark3Fields(i.Payload)
}

// ParseMobileStationClassmark3Fields decodes a mobile station classmark 3 value
// without its IEI and length octet.
func ParseMobileStationClassmark3Fields(payload []byte) (*MobileStationClassmark3Fields, error) {
	fields := &MobileStationClassmark3Fields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the payload is within [0, 32] octets.
func (f *MobileStationClassmark3Fields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: mobile station classmark 3 fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 0 {
		return fmt.Errorf("%w: mobile station classmark 3 has %d value octets, want >= 0", ErrInvalidLength, n)
	}
	if n > 32 {
		return fmt.Errorf("%w: mobile station classmark 3 has %d value octets, want <= 32", ErrInvalidLength, n)
	}
	return nil
}

// Marshal encodes fields into a mobile station classmark 3 value without its IEI
// and length octet.
func (f *MobileStationClassmark3Fields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *MobileStationClassmark3Fields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the mobile station classmark 3 value payload, or
// zero for a nil receiver or nil Value.
func (f *MobileStationClassmark3Fields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes a mobile station classmark 3 value without its IEI and
// length octet. It copies the payload into a temporary struct first, leaving f
// unchanged on failure.
func (f *MobileStationClassmark3Fields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: mobile station classmark 3 fields are nil", ErrInvalidValue)
	}
	if len(payload) > 32 {
		return fmt.Errorf("%w: mobile station classmark 3 has %d value octets, want <= 32", ErrInvalidLength, len(payload))
	}
	var decoded MobileStationClassmark3Fields
	if len(payload) > 0 {
		decoded.Value = make([]byte, len(payload))
		copy(decoded.Value, payload)
	}
	*f = decoded
	return nil
}
