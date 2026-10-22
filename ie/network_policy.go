package ie

import (
	"fmt"
	"io"
)

// NetworkPolicyFields is the structured value of the Network policy IE (TS
// 24.301 clause 9.9.3.52).
type NetworkPolicyFields struct {
	// RedirPolicy indicates the redirection security policy.
	RedirPolicy bool
}

// NewNetworkPolicyWithFields creates a tagged Network policy IE from fields.
func NewNetworkPolicyWithFields(fields *NetworkPolicyFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewNetworkPolicy(payload[0]), nil
}

// NetworkPolicy returns the IE as structured fields.
func (i *IE) NetworkPolicy() (*NetworkPolicyFields, error) {
	if err := i.requireDefinition(NetworkPolicyDefinition); err != nil {
		return nil, err
	}
	return ParseNetworkPolicyFields(i.Payload)
}

// ParseNetworkPolicyFields decodes a Network policy value without its IEI.
func ParseNetworkPolicyFields(payload []byte) (*NetworkPolicyFields, error) {
	fields := &NetworkPolicyFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *NetworkPolicyFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: network policy fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a Network policy value, without its IEI.
func (f *NetworkPolicyFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	var encoded byte
	if f.RedirPolicy {
		encoded |= 0x01
	}
	return []byte{encoded}, nil
}

// MarshalTo encodes fields into b.
func (f *NetworkPolicyFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the Network policy value payload.
func (f *NetworkPolicyFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a Network policy value without its IEI.
func (f *NetworkPolicyFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: network policy fields are nil", ErrInvalidValue)
	}
	if len(payload) == 0 {
		return io.ErrShortBuffer
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: network policy has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xfe != 0 {
		return fmt.Errorf("%w: network policy spare bits are nonzero", ErrInvalidValue)
	}
	*f = NetworkPolicyFields{
		RedirPolicy: payload[0]&0x01 != 0,
	}
	return nil
}
