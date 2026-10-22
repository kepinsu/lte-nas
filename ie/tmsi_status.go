package ie

import (
	"fmt"
	"io"
)

// TMSIStatusFields is the structured value of the tagged TMSI status IE (TS
// 24.301 clause 9.9.3.31). Only bit 0 is used (0=clear, 1=set); bits 2-7 are
// reserved and must be zero.
type TMSIStatusFields struct {
	Status uint8 // bit 0: 0=clear, 1=set; bits 2-7 reserved (zero)
}

// NewTMSIStatusWithFields creates a tagged TMSI status IE from fields.
func NewTMSIStatusWithFields(fields *TMSIStatusFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewTMSIStatus(payload[0]), nil
}

// NewPositionalTMSIStatusWithFields creates a positional TMSI status IE from
// fields.
func NewPositionalTMSIStatusWithFields(fields *TMSIStatusFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalTMSIStatus(payload[0]), nil
}

// TMSIStatus returns the IE as structured fields.
func (i *IE) TMSIStatus() (*TMSIStatusFields, error) {
	if err := i.requireDefinition(TMSIStatusDefinition); err != nil {
		return nil, err
	}
	return ParseTMSIStatusFields(i.Payload)
}

// ParseTMSIStatusFields decodes a TMSI status value without its IEI.
func ParseTMSIStatusFields(payload []byte) (*TMSIStatusFields, error) {
	fields := &TMSIStatusFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the TMSI status value uses only bit 1 (set=0x02, clear=0x00).
func (f *TMSIStatusFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: TMSI status fields are nil", ErrInvalidValue)
	}
	if f.Status != 0x00 && f.Status != 0x02 {
		return fmt.Errorf("%w: TMSI status %#02x is reserved", ErrInvalidValue, f.Status)
	}
	return nil
}

// Marshal encodes fields into a TMSI status value without its IEI.
func (f *TMSIStatusFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{f.Status}, nil
}

// MarshalTo encodes fields into b.
func (f *TMSIStatusFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	b[0] = f.Status
	return nil
}

// MarshalLen returns the size of the TMSI status value payload.
func (f *TMSIStatusFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a TMSI status value without its IEI. It leaves f
// unchanged when the payload is invalid.
func (f *TMSIStatusFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: TMSI status fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: TMSI status has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xfc != 0 || payload[0]&0x01 != 0 {
		return fmt.Errorf("%w: TMSI status reserved bits are nonzero", ErrInvalidValue)
	}
	var decoded TMSIStatusFields
	decoded.Status = payload[0]
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}
