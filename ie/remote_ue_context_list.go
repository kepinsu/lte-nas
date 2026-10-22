package ie

import (
	"fmt"
	"io"
)

// RemoteUEContextListFields is the structured value of the remote UE context
// list IE (TS 24.301 clause 9.9.4.20). It carries an opaque byte string
// between 2 and 65535 octets.
type RemoteUEContextListFields struct {
	Value []byte
}

// NewPositionalRemoteUEContextListWithFields creates a positional Remote UE
// context list IE from fields.
func NewPositionalRemoteUEContextListWithFields(fields *RemoteUEContextListFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalRemoteUEContextList(payload), nil
}

// RemoteUEContextList returns the IE as structured fields.
func (i *IE) RemoteUEContextList() (*RemoteUEContextListFields, error) {
	if err := i.requireDefinition(RemoteUEContextListDefinition); err != nil {
		return nil, err
	}
	return ParseRemoteUEContextListFields(i.Payload)
}

// ParseRemoteUEContextListFields decodes a remote UE context list value without
// its IEI and length octet.
func ParseRemoteUEContextListFields(payload []byte) (*RemoteUEContextListFields, error) {
	fields := &RemoteUEContextListFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the payload is within [2, 65535] octets.
func (f *RemoteUEContextListFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: remote UE context list fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 2 {
		return fmt.Errorf("%w: remote UE context list has %d value octets, want >= 2", ErrInvalidLength, n)
	}
	if n > 65535 {
		return fmt.Errorf("%w: remote UE context list has %d value octets, want <= 65535", ErrInvalidLength, n)
	}
	return nil
}

// Marshal encodes fields into a remote UE context list value without its IEI and
// length octet.
func (f *RemoteUEContextListFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *RemoteUEContextListFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the remote UE context list value payload, or
// zero for a nil receiver or nil Value.
func (f *RemoteUEContextListFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes a remote UE context list value without its IEI and
// length octet. It copies the payload into a temporary struct first, leaving f
// unchanged on failure.
func (f *RemoteUEContextListFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: remote UE context list fields are nil", ErrInvalidValue)
	}
	if len(payload) < 2 {
		return fmt.Errorf("%w: remote UE context list has %d value octets, want >= 2", ErrInvalidLength, len(payload))
	}
	if len(payload) > 65535 {
		return fmt.Errorf("%w: remote UE context list has %d value octets, want <= 65535", ErrInvalidLength, len(payload))
	}
	var decoded RemoteUEContextListFields
	decoded.Value = make([]byte, len(payload))
	copy(decoded.Value, payload)
	*f = decoded
	return nil
}
