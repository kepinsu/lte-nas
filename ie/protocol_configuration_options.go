package ie

import (
	"fmt"
	"io"
)

// ProtocolConfigurationOptionsFields is the structured value of the protocol
// configuration options IE (TS 24.301 clause 9.9.4.11). It carries an opaque
// byte string between 1 and 251 octets.
type ProtocolConfigurationOptionsFields struct {
	Value []byte
}

// NewProtocolConfigurationOptionsWithFields creates a tagged Protocol
// configuration options IE from fields.
func NewProtocolConfigurationOptionsWithFields(fields *ProtocolConfigurationOptionsFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewProtocolConfigurationOptions(payload), nil
}

// NewPositionalProtocolConfigurationOptionsWithFields creates a positional
// Protocol configuration options IE from fields.
func NewPositionalProtocolConfigurationOptionsWithFields(fields *ProtocolConfigurationOptionsFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalProtocolConfigurationOptions(payload), nil
}

// ProtocolConfigurationOptions returns the IE as structured fields.
func (i *IE) ProtocolConfigurationOptions() (*ProtocolConfigurationOptionsFields, error) {
	if err := i.requireDefinition(ProtocolConfigurationOptionsDefinition); err != nil {
		return nil, err
	}
	return ParseProtocolConfigurationOptionsFields(i.Payload)
}

// ParseProtocolConfigurationOptionsFields decodes a protocol configuration
// options value without its IEI and length octet.
func ParseProtocolConfigurationOptionsFields(payload []byte) (*ProtocolConfigurationOptionsFields, error) {
	fields := &ProtocolConfigurationOptionsFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the payload is within [1, 251] octets.
func (f *ProtocolConfigurationOptionsFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: protocol configuration options fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 1 {
		return fmt.Errorf("%w: protocol configuration options has %d value octets, want >= 1", ErrInvalidLength, n)
	}
	if n > 251 {
		return fmt.Errorf("%w: protocol configuration options has %d value octets, want <= 251", ErrInvalidLength, n)
	}
	return nil
}

// Marshal encodes fields into a protocol configuration options value without
// its IEI and length octet.
func (f *ProtocolConfigurationOptionsFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *ProtocolConfigurationOptionsFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the protocol configuration options value
// payload, or zero for a nil receiver or nil Value.
func (f *ProtocolConfigurationOptionsFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes a protocol configuration options value without its
// IEI and length octet. It copies the payload into a temporary struct first,
// leaving f unchanged on failure.
func (f *ProtocolConfigurationOptionsFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: protocol configuration options fields are nil", ErrInvalidValue)
	}
	if len(payload) < 1 {
		return fmt.Errorf("%w: protocol configuration options has %d value octets, want >= 1", ErrInvalidLength, len(payload))
	}
	if len(payload) > 251 {
		return fmt.Errorf("%w: protocol configuration options has %d value octets, want <= 251", ErrInvalidLength, len(payload))
	}
	var decoded ProtocolConfigurationOptionsFields
	decoded.Value = make([]byte, len(payload))
	copy(decoded.Value, payload)
	*f = decoded
	return nil
}
