package ie

import (
	"encoding/binary"
	"fmt"
	"io"
)

// TrackingAreaIdentityFields is the structured value of the Tracking area
// identity IE (TS 24.301 clause 9.9.3.32).
type TrackingAreaIdentityFields struct {
	PLMN             PLMN
	TrackingAreaCode uint16
}

// NewTrackingAreaIdentityWithFields creates a tagged Tracking area identity
// IE from fields.
func NewTrackingAreaIdentityWithFields(fields *TrackingAreaIdentityFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewTrackingAreaIdentity(payload), nil
}

// NewPositionalTrackingAreaIdentityWithFields creates a positional Tracking
// area identity IE from fields.
func NewPositionalTrackingAreaIdentityWithFields(fields *TrackingAreaIdentityFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalTrackingAreaIdentity(payload), nil
}

// TrackingAreaIdentity returns the IE as structured fields.
func (i *IE) TrackingAreaIdentity() (*TrackingAreaIdentityFields, error) {
	if err := i.requireDefinition(TrackingAreaIdentityDefinition); err != nil {
		return nil, err
	}
	return ParseTrackingAreaIdentityFields(i.Payload)
}

// ParseTrackingAreaIdentityFields decodes a Tracking area identity value
// without its IEI.
func ParseTrackingAreaIdentityFields(payload []byte) (*TrackingAreaIdentityFields, error) {
	fields := &TrackingAreaIdentityFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the PLMN is valid.
func (f *TrackingAreaIdentityFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: tracking area identity fields are nil", ErrInvalidValue)
	}
	return f.PLMN.validate()
}

// Marshal encodes fields into a five-octet Tracking area identity value.
func (f *TrackingAreaIdentityFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	plmn, err := f.PLMN.marshal()
	if err != nil {
		return nil, err
	}
	payload := make([]byte, 5)
	copy(payload, plmn[:])
	binary.BigEndian.PutUint16(payload[3:], f.TrackingAreaCode)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *TrackingAreaIdentityFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the Tracking area identity value.
func (f *TrackingAreaIdentityFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 5
}

// UnmarshalBinary decodes a five-octet Tracking area identity value.
func (f *TrackingAreaIdentityFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: tracking area identity fields are nil", ErrInvalidValue)
	}
	if len(payload) != 5 {
		return fmt.Errorf("%w: tracking area identity has %d value octets, want 5", ErrInvalidLength, len(payload))
	}
	plmn, err := parsePLMN(payload[:3])
	if err != nil {
		return err
	}
	*f = TrackingAreaIdentityFields{
		PLMN:             plmn,
		TrackingAreaCode: binary.BigEndian.Uint16(payload[3:]),
	}
	return nil
}
