package ie

import (
	"encoding/binary"
	"fmt"
	"io"
)

// ImsioffsetFields is the structured value of the IMSI offset IE (TS 24.301 clause 9.9.3.64).
// The two-octet value encodes an IMSI offset as a big-endian uint16.
type ImsioffsetFields struct {
	// IMSIOffset is the offset value as a big-endian uint16.
	IMSIOffset uint16
}

// NewIMSIOffsetWithFields creates a tagged IMSI offset IE from fields.
func NewIMSIOffsetWithFields(fields *ImsioffsetFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewIMSIOffset(payload), nil
}

// NewPositionalIMSIOffsetWithFields creates a positional IMSI offset IE from fields.
func NewPositionalIMSIOffsetWithFields(fields *ImsioffsetFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalIMSIOffset(payload), nil
}

// IMSIOffset returns the IE as structured fields.
func (i *IE) IMSIOffset() (*ImsioffsetFields, error) {
	if err := i.requireDefinition(IMSIOffsetDefinition); err != nil {
		return nil, err
	}
	return ParseIMSIOffsetFields(i.Payload)
}

// ParseIMSIOffsetFields decodes an IMSI offset value without its IEI and length octet.
func ParseIMSIOffsetFields(payload []byte) (*ImsioffsetFields, error) {
	fields := &ImsioffsetFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *ImsioffsetFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: IMSI offset fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a 2-octet IMSI offset value.
func (f *ImsioffsetFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, 2)
	binary.BigEndian.PutUint16(payload, f.IMSIOffset)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *ImsioffsetFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 2 {
		return io.ErrShortBuffer
	}
	binary.BigEndian.PutUint16(b, f.IMSIOffset)
	return nil
}

// MarshalLen returns the size of the IMSI offset value payload.
func (f *ImsioffsetFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 2
}

// UnmarshalBinary decodes a 2-octet IMSI offset value.
func (f *ImsioffsetFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: IMSI offset fields are nil", ErrInvalidValue)
	}
	if len(payload) != 2 {
		return fmt.Errorf("%w: IMSI offset has %d value octets, want 2", ErrInvalidLength, len(payload))
	}
	var decoded ImsioffsetFields
	decoded.IMSIOffset = binary.BigEndian.Uint16(payload)
	*f = decoded
	return nil
}
