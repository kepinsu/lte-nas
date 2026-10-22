package ie

import (
	"encoding/binary"
	"fmt"
	"io"
)

// ExtendedProtocolConfigurationOptionsFields is the structured value of the
// extended protocol configuration options IE (TS 24.301 clause 9.9.4.26). It
// carries an opaque byte string with a maximum of 65535 octets encoded as an
// LV (two-octet big-endian length) for positional wire format, or TLV-E with
// IEI=0x7b for tagged format.
type ExtendedProtocolConfigurationOptionsFields struct {
	Value []byte
}

// NewExtendedProtocolConfigurationOptionsWithFields creates a tagged Extended
// protocol configuration options IE from fields.
func NewExtendedProtocolConfigurationOptionsWithFields(fields *ExtendedProtocolConfigurationOptionsFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewExtendedProtocolConfigurationOptions(payload), nil
}

// NewPositionalExtendedProtocolConfigurationOptionsWithFields creates a
// positional Extended protocol configuration options IE from fields.
func NewPositionalExtendedProtocolConfigurationOptionsWithFields(fields *ExtendedProtocolConfigurationOptionsFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalExtendedProtocolConfigurationOptions(payload), nil
}

// ExtendedProtocolConfigurationOptions returns the IE as structured fields.
func (i *IE) ExtendedProtocolConfigurationOptions() (*ExtendedProtocolConfigurationOptionsFields, error) {
	if err := i.requireDefinition(ExtendedProtocolConfigurationOptionsDefinition); err != nil {
		return nil, err
	}
	return ParseExtendedProtocolConfigurationOptionsFields(i.Payload)
}

// ParseExtendedProtocolConfigurationOptionsFields decodes an extended protocol
// configuration options value without its IEI and two-octet length indicator.
func ParseExtendedProtocolConfigurationOptionsFields(payload []byte) (*ExtendedProtocolConfigurationOptionsFields, error) {
	fields := &ExtendedProtocolConfigurationOptionsFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the payload is within [1, 65535] octets.
func (f *ExtendedProtocolConfigurationOptionsFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: extended protocol configuration options fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 1 {
		return fmt.Errorf("%w: extended protocol configuration options has %d value octets, want >= 1", ErrInvalidLength, n)
	}
	if n > 65535 {
		return fmt.Errorf("%w: extended protocol configuration options has %d value octets, want <= 65535", ErrInvalidLength, n)
	}
	return nil
}

// Marshal encodes fields into an extended protocol configuration options value
// without its IEI and two-octet length indicator.
func (f *ExtendedProtocolConfigurationOptionsFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *ExtendedProtocolConfigurationOptionsFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the extended protocol configuration options
// value payload, or zero for a nil receiver or nil Value.
func (f *ExtendedProtocolConfigurationOptionsFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes an extended protocol configuration options value
// without its IEI and two-octet length indicator. It copies the payload into a
// temporary struct first, leaving f unchanged on failure.
func (f *ExtendedProtocolConfigurationOptionsFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: extended protocol configuration options fields are nil", ErrInvalidValue)
	}
	if len(payload) < 1 {
		return fmt.Errorf("%w: extended protocol configuration options has %d value octets, want >= 1", ErrInvalidLength, len(payload))
	}
	if len(payload) > 65535 {
		return fmt.Errorf("%w: extended protocol configuration options has %d value octets, want <= 65535", ErrInvalidLength, len(payload))
	}
	var decoded ExtendedProtocolConfigurationOptionsFields
	decoded.Value = make([]byte, len(payload))
	copy(decoded.Value, payload)
	*f = decoded
	return nil
}

// MaxExtendedProtocolConfigurationOptionsLen is the maximum length for a tagged
// ExtendedProtocolConfigurationOptions wire representation: 1 (IEI) + 2 (length)
// + 65535 (payload).
const MaxExtendedProtocolConfigurationOptionsLen = 1 + 2 + 65535

// MarshalTagged encodes fields into the full tagged TLV-E wire format (IEI +
// 2-octet big-endian length + payload).
func (f *ExtendedProtocolConfigurationOptionsFields) MarshalTagged() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	n := len(f.Value)
	payload := make([]byte, 3+n)
	payload[0] = byte(ExtendedProtocolConfigurationOptions)
	binary.BigEndian.PutUint16(payload[1:3], uint16(n))
	copy(payload[3:], f.Value)
	return payload, nil
}

// UnmarshalBinaryTagged decodes a tagged TLV-E wire format. The first byte is
// the IEI, bytes 1-2 are the big-endian length, and the remaining bytes are
// the payload.
func (f *ExtendedProtocolConfigurationOptionsFields) UnmarshalBinaryTagged(b []byte) error {
	if f == nil {
		return fmt.Errorf("%w: extended protocol configuration options fields are nil", ErrInvalidValue)
	}
	if len(b) < 3 {
		return io.ErrUnexpectedEOF
	}
	length := int(binary.BigEndian.Uint16(b[1:3]))
	if length < 1 || length > 65535 {
		return fmt.Errorf("%w: extended protocol configuration options has %d value octets, want 1..65535", ErrInvalidLength, length)
	}
	if len(b) < 3+length {
		return io.ErrUnexpectedEOF
	}
	var decoded ExtendedProtocolConfigurationOptionsFields
	decoded.Value = make([]byte, length)
	copy(decoded.Value, b[3:3+length])
	*f = decoded
	return nil
}
