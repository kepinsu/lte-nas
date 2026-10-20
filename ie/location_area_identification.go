package ie

import (
	"encoding/binary"
	"fmt"
	"io"
)

// LocationAreaIdentificationFields is the structured value of the Location
// area identification IE (TS 24.301 clause 9.9.2.2).
type LocationAreaIdentificationFields struct {
	PLMN             PLMN
	LocationAreaCode uint16
}

// NewLocationAreaIdentificationWithFields creates a tagged Location area
// identification IE from fields.
func NewLocationAreaIdentificationWithFields(fields *LocationAreaIdentificationFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewLocationAreaIdentification(payload), nil
}

// NewPositionalLocationAreaIdentificationWithFields creates a positional
// Location area identification IE from fields.
func NewPositionalLocationAreaIdentificationWithFields(fields *LocationAreaIdentificationFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalLocationAreaIdentification(payload), nil
}

// LocationAreaIdentification returns the IE as structured fields.
func (i *IE) LocationAreaIdentification() (*LocationAreaIdentificationFields, error) {
	if err := i.requireDefinition(LocationAreaIdentificationDefinition); err != nil {
		return nil, err
	}
	return ParseLocationAreaIdentificationFields(i.Payload)
}

// ParseLocationAreaIdentificationFields decodes a Location area
// identification value without its IEI.
func ParseLocationAreaIdentificationFields(payload []byte) (*LocationAreaIdentificationFields, error) {
	fields := &LocationAreaIdentificationFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the PLMN is valid.
func (f *LocationAreaIdentificationFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: location area identification fields are nil", ErrInvalidValue)
	}
	return f.PLMN.validate()
}

// Marshal encodes fields into a five-octet Location area identification value.
func (f *LocationAreaIdentificationFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	plmn, err := f.PLMN.marshal()
	if err != nil {
		return nil, err
	}
	payload := make([]byte, 5)
	copy(payload, plmn[:])
	binary.BigEndian.PutUint16(payload[3:], f.LocationAreaCode)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *LocationAreaIdentificationFields) MarshalTo(b []byte) error {
	payload, err := f.Marshal()
	if err != nil {
		return err
	}
	if len(b) < len(payload) {
		return io.ErrShortBuffer
	}
	copy(b, payload)
	return nil
}

// MarshalLen returns the size of the Location area identification value.
func (f *LocationAreaIdentificationFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 5
}

// UnmarshalBinary decodes a five-octet Location area identification value.
func (f *LocationAreaIdentificationFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: location area identification fields are nil", ErrInvalidValue)
	}
	if len(payload) != 5 {
		return fmt.Errorf("%w: location area identification has %d value octets, want 5", ErrInvalidLength, len(payload))
	}
	plmn, err := parsePLMN(payload[:3])
	if err != nil {
		return err
	}
	*f = LocationAreaIdentificationFields{
		PLMN:             plmn,
		LocationAreaCode: binary.BigEndian.Uint16(payload[3:]),
	}
	return nil
}
