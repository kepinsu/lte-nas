package ie

import (
	"fmt"
	"io"
)

// UMTSCipheringAlgorithm identifies a UMTS encryption algorithm.
type UMTSCipheringAlgorithm uint8

// UMTS encryption algorithm identifiers.
const (
	UMTSCipheringAlgorithmUEA0 UMTSCipheringAlgorithm = iota
	UMTSCipheringAlgorithmUEA1
	UMTSCipheringAlgorithmUEA2
	UMTSCipheringAlgorithmUEA3
	UMTSCipheringAlgorithmUEA4
	UMTSCipheringAlgorithmUEA5
	UMTSCipheringAlgorithmUEA6
	UMTSCipheringAlgorithmUEA7
)

// UMTSIntegrityAlgorithm identifies a UMTS integrity algorithm.
type UMTSIntegrityAlgorithm uint8

// UMTS integrity algorithm identifiers.
const (
	UMTSIntegrityAlgorithmUIA1 UMTSIntegrityAlgorithm = 1 + iota
	UMTSIntegrityAlgorithmUIA2
	UMTSIntegrityAlgorithmUIA3
	UMTSIntegrityAlgorithmUIA4
	UMTSIntegrityAlgorithmUIA5
	UMTSIntegrityAlgorithmUIA6
	UMTSIntegrityAlgorithmUIA7
)

// GPRSCipheringAlgorithm identifies a GPRS encryption algorithm.
type GPRSCipheringAlgorithm uint8

// GPRS encryption algorithm identifiers.
const (
	GPRSCipheringAlgorithmGEA1 GPRSCipheringAlgorithm = 1 + iota
	GPRSCipheringAlgorithmGEA2
	GPRSCipheringAlgorithmGEA3
	GPRSCipheringAlgorithmGEA4
	GPRSCipheringAlgorithmGEA5
	GPRSCipheringAlgorithmGEA6
	GPRSCipheringAlgorithmGEA7
)

// UESecurityCapabilityFields is the structured value of the UE security
// capability IE (TS 24.301 clause 9.9.3.36). Each algorithm slice is a sorted
// set of supported algorithms. Empty UMTS and GPRS slices omit their optional
// octets as prescribed by the specification.
type UESecurityCapabilityFields struct {
	EPSCipheringAlgorithms          []NASEncryptionAlgorithm
	EPSIntegrityAlgorithms          []NASIntegrityAlgorithm
	EPSUserPlaneIntegrityProtection bool
	UMTSCipheringAlgorithms         []UMTSCipheringAlgorithm
	UMTSIntegrityAlgorithms         []UMTSIntegrityAlgorithm
	GPRSCipheringAlgorithms         []GPRSCipheringAlgorithm
}

// NewPositionalUESecurityCapabilityWithFields creates a positional UE security
// capability IE from fields.
func NewPositionalUESecurityCapabilityWithFields(fields *UESecurityCapabilityFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalUESecurityCapability(payload), nil
}

// UESecurityCapability returns the IE as structured fields.
func (i *IE) UESecurityCapability() (*UESecurityCapabilityFields, error) {
	if err := i.requireDefinition(UESecurityCapabilityDefinition); err != nil {
		return nil, err
	}
	return ParseUESecurityCapabilityFields(i.Payload)
}

// ParseUESecurityCapabilityFields decodes a UE security capability value
// without its IEI and length octet.
func ParseUESecurityCapabilityFields(payload []byte) (*UESecurityCapabilityFields, error) {
	fields := &UESecurityCapabilityFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that every algorithm slice is an ordered set of values
// defined by TS 24.301.
func (f *UESecurityCapabilityFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: UE security capability fields are nil", ErrInvalidValue)
	}
	if err := validateNASEncryptionAlgorithmSet(f.EPSCipheringAlgorithms); err != nil {
		return err
	}
	if err := validateNASIntegrityAlgorithmSet(f.EPSIntegrityAlgorithms); err != nil {
		return err
	}
	if err := validateUMTSCipheringAlgorithmSet(f.UMTSCipheringAlgorithms); err != nil {
		return err
	}
	if err := validateUMTSIntegrityAlgorithmSet(f.UMTSIntegrityAlgorithms); err != nil {
		return err
	}
	return validateGPRSCipheringAlgorithmSet(f.GPRSCipheringAlgorithms)
}

// Marshal encodes fields into a UE security capability value without its IEI
// and length octet.
func (f *UESecurityCapabilityFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, f.MarshalLen())
	payload[0] = encodeNASEncryptionAlgorithmSet(f.EPSCipheringAlgorithms)
	payload[1] = encodeNASIntegrityAlgorithmSet(f.EPSIntegrityAlgorithms)
	if f.EPSUserPlaneIntegrityProtection {
		payload[1] |= 0x01
	}
	if len(payload) >= 4 {
		payload[2] = encodeUMTSCipheringAlgorithmSet(f.UMTSCipheringAlgorithms)
		payload[3] = encodeUMTSIntegrityAlgorithmSet(f.UMTSIntegrityAlgorithms)
	}
	if len(payload) == 5 {
		payload[4] = encodeGPRSCipheringAlgorithmSet(f.GPRSCipheringAlgorithms)
	}
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *UESecurityCapabilityFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the UE security capability value payload
// produced by Marshal. It returns zero for a nil receiver.
func (f *UESecurityCapabilityFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	if len(f.GPRSCipheringAlgorithms) != 0 {
		return 5
	}
	if len(f.UMTSCipheringAlgorithms) != 0 || len(f.UMTSIntegrityAlgorithms) != 0 {
		return 4
	}
	return 2
}

// UnmarshalBinary decodes a UE security capability value without its IEI and
// length octet. It leaves f unchanged when payload is invalid.
func (f *UESecurityCapabilityFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: UE security capability fields are nil", ErrInvalidValue)
	}
	if len(payload) != 2 && len(payload) != 4 && len(payload) != 5 {
		return fmt.Errorf("%w: UE security capability has %d value octets, want 2, 4, or 5", ErrInvalidLength, len(payload))
	}
	if len(payload) >= 4 && payload[3]&0x80 != 0 {
		return fmt.Errorf("%w: UE security capability UMTS integrity spare bit is nonzero", ErrInvalidValue)
	}
	if len(payload) == 5 && payload[4]&0x80 != 0 {
		return fmt.Errorf("%w: UE security capability GPRS ciphering spare bit is nonzero", ErrInvalidValue)
	}
	decoded := UESecurityCapabilityFields{
		EPSCipheringAlgorithms:          decodeNASEncryptionAlgorithmSet(payload[0]),
		EPSIntegrityAlgorithms:          decodeNASIntegrityAlgorithmSet(payload[1]),
		EPSUserPlaneIntegrityProtection: payload[1]&0x01 != 0,
	}
	if len(payload) >= 4 {
		decoded.UMTSCipheringAlgorithms = decodeUMTSCipheringAlgorithmSet(payload[2])
		decoded.UMTSIntegrityAlgorithms = decodeUMTSIntegrityAlgorithmSet(payload[3])
	}
	if len(payload) == 5 {
		decoded.GPRSCipheringAlgorithms = decodeGPRSCipheringAlgorithmSet(payload[4])
	}
	*f = decoded
	return nil
}

func validateNASEncryptionAlgorithmSet(algorithms []NASEncryptionAlgorithm) error {
	for index, algorithm := range algorithms {
		if algorithm > NASEncryptionAlgorithmEEA7 {
			return fmt.Errorf("%w: EPS ciphering algorithm %d is reserved", ErrInvalidValue, algorithm)
		}
		if index > 0 && algorithms[index-1] >= algorithm {
			return fmt.Errorf("%w: EPS ciphering algorithms must be strictly ascending", ErrInvalidValue)
		}
	}
	return nil
}

func validateNASIntegrityAlgorithmSet(algorithms []NASIntegrityAlgorithm) error {
	for index, algorithm := range algorithms {
		if algorithm > NASIntegrityAlgorithmEIA6 {
			return fmt.Errorf("%w: EPS integrity algorithm %d is reserved", ErrInvalidValue, algorithm)
		}
		if index > 0 && algorithms[index-1] >= algorithm {
			return fmt.Errorf("%w: EPS integrity algorithms must be strictly ascending", ErrInvalidValue)
		}
	}
	return nil
}

func validateUMTSCipheringAlgorithmSet(algorithms []UMTSCipheringAlgorithm) error {
	for index, algorithm := range algorithms {
		if algorithm > UMTSCipheringAlgorithmUEA7 {
			return fmt.Errorf("%w: UMTS ciphering algorithm %d is reserved", ErrInvalidValue, algorithm)
		}
		if index > 0 && algorithms[index-1] >= algorithm {
			return fmt.Errorf("%w: UMTS ciphering algorithms must be strictly ascending", ErrInvalidValue)
		}
	}
	return nil
}

func validateUMTSIntegrityAlgorithmSet(algorithms []UMTSIntegrityAlgorithm) error {
	for index, algorithm := range algorithms {
		if algorithm < UMTSIntegrityAlgorithmUIA1 || algorithm > UMTSIntegrityAlgorithmUIA7 {
			return fmt.Errorf("%w: UMTS integrity algorithm %d is reserved", ErrInvalidValue, algorithm)
		}
		if index > 0 && algorithms[index-1] >= algorithm {
			return fmt.Errorf("%w: UMTS integrity algorithms must be strictly ascending", ErrInvalidValue)
		}
	}
	return nil
}

func validateGPRSCipheringAlgorithmSet(algorithms []GPRSCipheringAlgorithm) error {
	for index, algorithm := range algorithms {
		if algorithm < GPRSCipheringAlgorithmGEA1 || algorithm > GPRSCipheringAlgorithmGEA7 {
			return fmt.Errorf("%w: GPRS ciphering algorithm %d is reserved", ErrInvalidValue, algorithm)
		}
		if index > 0 && algorithms[index-1] >= algorithm {
			return fmt.Errorf("%w: GPRS ciphering algorithms must be strictly ascending", ErrInvalidValue)
		}
	}
	return nil
}

func encodeNASEncryptionAlgorithmSet(algorithms []NASEncryptionAlgorithm) byte {
	var encoded byte
	for _, algorithm := range algorithms {
		encoded |= 1 << (7 - algorithm)
	}
	return encoded
}

func encodeNASIntegrityAlgorithmSet(algorithms []NASIntegrityAlgorithm) byte {
	var encoded byte
	for _, algorithm := range algorithms {
		encoded |= 1 << (7 - algorithm)
	}
	return encoded
}

func encodeUMTSCipheringAlgorithmSet(algorithms []UMTSCipheringAlgorithm) byte {
	var encoded byte
	for _, algorithm := range algorithms {
		encoded |= 1 << (7 - algorithm)
	}
	return encoded
}

func encodeUMTSIntegrityAlgorithmSet(algorithms []UMTSIntegrityAlgorithm) byte {
	var encoded byte
	for _, algorithm := range algorithms {
		encoded |= 1 << (8 - algorithm)
	}
	return encoded
}

func encodeGPRSCipheringAlgorithmSet(algorithms []GPRSCipheringAlgorithm) byte {
	var encoded byte
	for _, algorithm := range algorithms {
		encoded |= 1 << (8 - algorithm)
	}
	return encoded
}

func decodeNASEncryptionAlgorithmSet(encoded byte) []NASEncryptionAlgorithm {
	algorithms := make([]NASEncryptionAlgorithm, 0, 8)
	for algorithm := NASEncryptionAlgorithmEEA0; algorithm <= NASEncryptionAlgorithmEEA7; algorithm++ {
		if encoded&(1<<(7-algorithm)) != 0 {
			algorithms = append(algorithms, algorithm)
		}
	}
	return algorithms
}

func decodeNASIntegrityAlgorithmSet(encoded byte) []NASIntegrityAlgorithm {
	algorithms := make([]NASIntegrityAlgorithm, 0, 7)
	for algorithm := NASIntegrityAlgorithmEIA0; algorithm <= NASIntegrityAlgorithmEIA6; algorithm++ {
		if encoded&(1<<(7-algorithm)) != 0 {
			algorithms = append(algorithms, algorithm)
		}
	}
	return algorithms
}

func decodeUMTSCipheringAlgorithmSet(encoded byte) []UMTSCipheringAlgorithm {
	algorithms := make([]UMTSCipheringAlgorithm, 0, 8)
	for algorithm := UMTSCipheringAlgorithmUEA0; algorithm <= UMTSCipheringAlgorithmUEA7; algorithm++ {
		if encoded&(1<<(7-algorithm)) != 0 {
			algorithms = append(algorithms, algorithm)
		}
	}
	return algorithms
}

func decodeUMTSIntegrityAlgorithmSet(encoded byte) []UMTSIntegrityAlgorithm {
	algorithms := make([]UMTSIntegrityAlgorithm, 0, 7)
	for algorithm := UMTSIntegrityAlgorithmUIA1; algorithm <= UMTSIntegrityAlgorithmUIA7; algorithm++ {
		if encoded&(1<<(8-algorithm)) != 0 {
			algorithms = append(algorithms, algorithm)
		}
	}
	return algorithms
}

func decodeGPRSCipheringAlgorithmSet(encoded byte) []GPRSCipheringAlgorithm {
	algorithms := make([]GPRSCipheringAlgorithm, 0, 7)
	for algorithm := GPRSCipheringAlgorithmGEA1; algorithm <= GPRSCipheringAlgorithmGEA7; algorithm++ {
		if encoded&(1<<(8-algorithm)) != 0 {
			algorithms = append(algorithms, algorithm)
		}
	}
	return algorithms
}
