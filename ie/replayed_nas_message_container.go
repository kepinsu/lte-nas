package ie

import (
	"fmt"
	"io"
)

// ReplayedNASMessageContainerFields is the structured value of the replayed NAS
// message container IE (TS 24.301 clause 9.9.3.51). It carries an opaque byte
// string between 0 and 65535 octets.
type ReplayedNASMessageContainerFields struct {
	Value []byte
}

// NewPositionalReplayedNASMessageContainerWithFields creates a positional Replayed
// NAS message container IE from fields.
func NewPositionalReplayedNASMessageContainerWithFields(fields *ReplayedNASMessageContainerFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalReplayedNASMessageContainer(payload), nil
}

// ReplayedNASMessageContainer returns the IE as structured fields.
func (i *IE) ReplayedNASMessageContainer() (*ReplayedNASMessageContainerFields, error) {
	if err := i.requireDefinition(ReplayedNASMessageContainerDefinition); err != nil {
		return nil, err
	}
	return ParseReplayedNASMessageContainerFields(i.Payload)
}

// ParseReplayedNASMessageContainerFields decodes a replayed NAS message container
// value without its IEI and length octet.
func ParseReplayedNASMessageContainerFields(payload []byte) (*ReplayedNASMessageContainerFields, error) {
	fields := &ReplayedNASMessageContainerFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the payload is within [0, 65535] octets.
func (f *ReplayedNASMessageContainerFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: replayed NAS message container fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 0 {
		return fmt.Errorf("%w: replayed NAS message container has %d value octets, want >= 0", ErrInvalidLength, n)
	}
	if n > 65535 {
		return fmt.Errorf("%w: replayed NAS message container has %d value octets, want <= 65535", ErrInvalidLength, n)
	}
	return nil
}

// Marshal encodes fields into a replayed NAS message container value without its
// IEI and length octet.
func (f *ReplayedNASMessageContainerFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *ReplayedNASMessageContainerFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the replayed NAS message container value payload,
// or zero for a nil receiver or nil Value.
func (f *ReplayedNASMessageContainerFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes a replayed NAS message container value without its IEI
// and length octet. It copies the payload into a temporary struct first, leaving f
// unchanged on failure.
func (f *ReplayedNASMessageContainerFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: replayed NAS message container fields are nil", ErrInvalidValue)
	}
	if len(payload) > 65535 {
		return fmt.Errorf("%w: replayed NAS message container has %d value octets, want <= 65535", ErrInvalidLength, len(payload))
	}
	var decoded ReplayedNASMessageContainerFields
	if payload == nil {
		decoded.Value = nil
	} else {
		decoded.Value = make([]byte, len(payload))
		copy(decoded.Value, payload)
	}
	*f = decoded
	return nil
}
