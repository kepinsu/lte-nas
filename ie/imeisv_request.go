package ie

import (
	"fmt"
	"io"
)

// IMEISVRequestFields is the structured value of the IMEISV request IE (TS 24.301
// clause 9.9.3.18, see TS 24.008 clause 10.5.5.10).
type IMEISVRequestFields struct {
	// IMEISV indicates whether the IMEISV is requested.
	IMEISV bool
}

// NewIMEISVRequestWithFields creates a tagged IMEISV request IE from fields.
func NewIMEISVRequestWithFields(fields *IMEISVRequestFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewIMEISVRequest(payload[0]), nil
}

// IMEISVRequest returns the IE as structured fields.
func (i *IE) IMEISVRequest() (*IMEISVRequestFields, error) {
	if err := i.requireDefinition(IMEISVRequestDefinition); err != nil {
		return nil, err
	}
	return ParseIMEISVRequestFields(i.Payload)
}

// ParseIMEISVRequestFields decodes an IMEISV request value without its IEI.
func ParseIMEISVRequestFields(payload []byte) (*IMEISVRequestFields, error) {
	fields := &IMEISVRequestFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *IMEISVRequestFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: IMEISV request fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into an IMEISV request value, without its IEI.
func (f *IMEISVRequestFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	var encoded byte
	if f.IMEISV {
		encoded |= 0x01
	}
	return []byte{encoded}, nil
}

// MarshalTo encodes fields into b.
func (f *IMEISVRequestFields) MarshalTo(b []byte) error {
	payload, err := f.Marshal()
	if err != nil {
		return err
	}
	if len(b) < len(payload) {
		return io.ErrShortBuffer
	}
	b[0] = payload[0]
	return nil
}

// MarshalLen returns the size of the IMEISV request value payload.
func (f *IMEISVRequestFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes an IMEISV request value without its IEI.
func (f *IMEISVRequestFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: IMEISV request fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: IMEISV request has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xfe != 0 {
		return fmt.Errorf("%w: IMEISV request spare bits are nonzero", ErrInvalidValue)
	}
	*f = IMEISVRequestFields{
		IMEISV: payload[0]&0x01 != 0,
	}
	return nil
}
