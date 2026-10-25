package ie

import (
	"fmt"
	"io"
)

// MobileStationClassmark2Fields is the structured value of the mobile station
// classmark 2 IE (TS 24.301 clause 9.9.2.4). The three-octet payload is an
// opaque classmark 2 information element.
type MobileStationClassmark2Fields struct {
	// Payload contains the 3-octet mobile station classmark 2 value.
	Payload [3]byte
}

// NewMobileStationClassmark2WithFields creates a tagged mobile station classmark 2
// IE from fields.
func NewMobileStationClassmark2WithFields(fields *MobileStationClassmark2Fields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewMobileStationClassmark2(payload), nil
}

// NewPositionalMobileStationClassmark2WithFields creates a positional mobile station
// classmark 2 IE from fields.
func NewPositionalMobileStationClassmark2WithFields(fields *MobileStationClassmark2Fields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalMobileStationClassmark2(payload), nil
}

// MobileStationClassmark2 returns the IE as structured fields.
func (i *IE) MobileStationClassmark2() (*MobileStationClassmark2Fields, error) {
	if err := i.requireDefinition(MobileStationClassmark2Definition); err != nil {
		return nil, err
	}
	return ParseMobileStationClassmark2Fields(i.Payload)
}

// ParseMobileStationClassmark2Fields decodes a mobile station classmark 2 value
// without its IEI and length octet.
func ParseMobileStationClassmark2Fields(payload []byte) (*MobileStationClassmark2Fields, error) {
	fields := &MobileStationClassmark2Fields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *MobileStationClassmark2Fields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: mobile station classmark 2 fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a 3-octet mobile station classmark 2 value.
func (f *MobileStationClassmark2Fields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Payload))
	copy(payload, f.Payload[:])
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *MobileStationClassmark2Fields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 3 {
		return io.ErrShortBuffer
	}
	copy(b, f.Payload[:])
	return nil
}

// MarshalLen returns the size of the mobile station classmark 2 value payload.
func (f *MobileStationClassmark2Fields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 3
}

// UnmarshalBinary decodes a 3-octet mobile station classmark 2 value.
func (f *MobileStationClassmark2Fields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: mobile station classmark 2 fields are nil", ErrInvalidValue)
	}
	if len(payload) != 3 {
		return fmt.Errorf("%w: mobile station classmark 2 has %d value octets, want 3", ErrInvalidLength, len(payload))
	}
	var decoded MobileStationClassmark2Fields
	copy(decoded.Payload[:], payload)
	*f = decoded
	return nil
}
