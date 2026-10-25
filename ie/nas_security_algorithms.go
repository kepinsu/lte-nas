package ie

import (
	"fmt"
	"io"
)

// NASIntegrityAlgorithm identifies an EPS integrity protection algorithm.
type NASIntegrityAlgorithm uint8

// NAS integrity protection algorithm identifiers.
const (
	NASIntegrityAlgorithmEIA0 NASIntegrityAlgorithm = iota
	NASIntegrityAlgorithm128EIA1
	NASIntegrityAlgorithm128EIA2
	NASIntegrityAlgorithm128EIA3
	NASIntegrityAlgorithmEIA4
	NASIntegrityAlgorithmEIA5
	NASIntegrityAlgorithmEIA6
	NASIntegrityAlgorithmEIA7
)

// NASEncryptionAlgorithm identifies an EPS ciphering algorithm.
type NASEncryptionAlgorithm uint8

// NAS ciphering algorithm identifiers.
const (
	NASEncryptionAlgorithmEEA0 NASEncryptionAlgorithm = iota
	NASEncryptionAlgorithm128EEA1
	NASEncryptionAlgorithm128EEA2
	NASEncryptionAlgorithm128EEA3
	NASEncryptionAlgorithmEEA4
	NASEncryptionAlgorithmEEA5
	NASEncryptionAlgorithmEEA6
	NASEncryptionAlgorithmEEA7
)

// NASSecurityAlgorithmsFields is the structured value of the positional NAS
// security algorithms IE (TS 24.301 clause 9.9.3.23).
type NASSecurityAlgorithmsFields struct {
	CipheringAlgorithm NASEncryptionAlgorithm
	IntegrityAlgorithm NASIntegrityAlgorithm
}

// NewPositionalNASSecurityAlgorithmsWithFields creates a positional NAS
// security algorithms IE from fields.
func NewPositionalNASSecurityAlgorithmsWithFields(fields *NASSecurityAlgorithmsFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalNASSecurityAlgorithms(payload[0]), nil
}

// NASSecurityAlgorithms returns the IE as structured fields.
func (i *IE) NASSecurityAlgorithms() (*NASSecurityAlgorithmsFields, error) {
	if err := i.requireDefinition(NASSecurityAlgorithmsDefinition); err != nil {
		return nil, err
	}
	return ParseNASSecurityAlgorithmsFields(i.Payload)
}

// ParseNASSecurityAlgorithmsFields decodes a NAS security algorithms value
// without its IEI.
func ParseNASSecurityAlgorithmsFields(payload []byte) (*NASSecurityAlgorithmsFields, error) {
	fields := &NASSecurityAlgorithmsFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that both algorithm identifiers fit three bits.
func (f *NASSecurityAlgorithmsFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: NAS security algorithms fields are nil", ErrInvalidValue)
	}
	if f.CipheringAlgorithm > 0x07 || f.IntegrityAlgorithm > 0x07 {
		return fmt.Errorf("%w: NAS security algorithm exceeds its encoded width", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a NAS security algorithms value without its IEI.
func (f *NASSecurityAlgorithmsFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.CipheringAlgorithm<<4) | byte(f.IntegrityAlgorithm)}, nil
}

// MarshalTo encodes fields into b.
func (f *NASSecurityAlgorithmsFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	b[0] = byte(f.CipheringAlgorithm<<4) | byte(f.IntegrityAlgorithm)
	return nil
}

// MarshalLen returns the size of the NAS security algorithms value payload.
func (f *NASSecurityAlgorithmsFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a NAS security algorithms value without its IEI.
func (f *NASSecurityAlgorithmsFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: NAS security algorithms fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: NAS security algorithms has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0x88 != 0 {
		return fmt.Errorf("%w: NAS security algorithms spare bits are nonzero", ErrInvalidValue)
	}
	*f = NASSecurityAlgorithmsFields{
		CipheringAlgorithm: NASEncryptionAlgorithm((payload[0] >> 4) & 0x07),
		IntegrityAlgorithm: NASIntegrityAlgorithm(payload[0] & 0x07),
	}
	return nil
}
