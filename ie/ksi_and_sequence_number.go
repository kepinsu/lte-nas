package ie

import (
	"fmt"
	"io"
)

// KSIAndSequenceNumberFields is the structured value of the KSI and sequence
// number IE (TS 24.301 clause 9.9.3.19). SequenceNumber contains the five
// least significant bits of the NAS COUNT value.
type KSIAndSequenceNumberFields struct {
	KeySetIdentifier uint8
	SequenceNumber   uint8
}

// NewPositionalKSIAndSequenceNumberWithFields creates a positional KSI and
// sequence number IE from fields.
func NewPositionalKSIAndSequenceNumberWithFields(fields *KSIAndSequenceNumberFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalKSIAndSequenceNumber(payload[0]), nil
}

// KSIAndSequenceNumber returns the IE as structured fields.
func (i *IE) KSIAndSequenceNumber() (*KSIAndSequenceNumberFields, error) {
	if err := i.requireDefinition(KSIAndSequenceNumberDefinition); err != nil {
		return nil, err
	}
	return ParseKSIAndSequenceNumberFields(i.Payload)
}

// ParseKSIAndSequenceNumberFields decodes a KSI and sequence number value
// without its IEI.
func ParseKSIAndSequenceNumberFields(payload []byte) (*KSIAndSequenceNumberFields, error) {
	fields := &KSIAndSequenceNumberFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the KSI and short sequence number fit their encoded
// bit widths.
func (f *KSIAndSequenceNumberFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: KSI and sequence number fields are nil", ErrInvalidValue)
	}
	if f.KeySetIdentifier > 0x07 || f.SequenceNumber > 0x1f {
		return fmt.Errorf("%w: KSI and sequence number field exceeds its encoded width", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a KSI and sequence number value without its IEI.
func (f *KSIAndSequenceNumberFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{f.KeySetIdentifier<<5 | f.SequenceNumber}, nil
}

// MarshalTo encodes fields into b.
func (f *KSIAndSequenceNumberFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	b[0] = f.KeySetIdentifier<<5 | f.SequenceNumber
	return nil
}

// MarshalLen returns the size of the KSI and sequence number value payload.
func (f *KSIAndSequenceNumberFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a KSI and sequence number value without its IEI.
func (f *KSIAndSequenceNumberFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: KSI and sequence number fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: KSI and sequence number has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	*f = KSIAndSequenceNumberFields{
		KeySetIdentifier: payload[0] >> 5,
		SequenceNumber:   payload[0] & 0x1f,
	}
	return nil
}
