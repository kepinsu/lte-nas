package ie

import (
	"fmt"
	"io"
)

// NASSecurityParametersToEUTRADefinition is the tagged IE (TS 24.301 clause
// 9.9.2.7). NAS security parameters to E-UTRA has format 3 (fixed-length TV).
// This package provides semantic field access for the opaque 6-octet value.

// NASSecurityParametersToEUTRAFields is the structured value of the NAS security
// parameters to E-UTRA IE (TS 24.301 clause 9.9.2.7). The six-octet payload
// contains encryption algorithm and integrity protection algorithm selection,
// along with additional security parameters.
type NASSecurityParametersToEUTRAFields struct {
	// Payload contains the 6-octet NAS security parameters to E-UTRA value.
	Payload [6]byte
}

// NewNASSecurityParametersToEUTRAWithFields creates a NAS security parameters
// to E-UTRA IE from fields. This uses the positional constructor since there is
// no tagged NewNASSecurityParametersToEUTRA in constructors_generated.go.
func NewNASSecurityParametersToEUTRAWithFields(fields *NASSecurityParametersToEUTRAFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalNASSecurityParametersToEUTRA(payload), nil
}

// NewPositionalNASSecurityParametersToEUTRAWithFields creates a positional NAS
// security parameters to E-UTRA IE from fields.
func NewPositionalNASSecurityParametersToEUTRAWithFields(fields *NASSecurityParametersToEUTRAFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalNASSecurityParametersToEUTRA(payload), nil
}

// NASSecurityParametersToEUTRA returns the IE as structured fields.
func (i *IE) NASSecurityParametersToEUTRA() (*NASSecurityParametersToEUTRAFields, error) {
	if err := i.requireDefinition(NASSecurityParametersToEUTRADefinition); err != nil {
		return nil, err
	}
	return ParseNASSecurityParametersToEUTRAFields(i.Payload)
}

// ParseNASSecurityParametersToEUTRAFields decodes a NAS security parameters to
// E-UTRA value without its IEI.
func ParseNASSecurityParametersToEUTRAFields(payload []byte) (*NASSecurityParametersToEUTRAFields, error) {
	fields := &NASSecurityParametersToEUTRAFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *NASSecurityParametersToEUTRAFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: NAS security parameters to E-UTRA fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a 6-octet NAS security parameters to E-UTRA value.
func (f *NASSecurityParametersToEUTRAFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Payload))
	copy(payload, f.Payload[:])
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *NASSecurityParametersToEUTRAFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 6 {
		return io.ErrShortBuffer
	}
	copy(b, f.Payload[:])
	return nil
}

// MarshalLen returns the size of the NAS security parameters to E-UTRA value payload.
func (f *NASSecurityParametersToEUTRAFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 6
}

// UnmarshalBinary decodes a 6-octet NAS security parameters to E-UTRA value.
func (f *NASSecurityParametersToEUTRAFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: NAS security parameters to E-UTRA fields are nil", ErrInvalidValue)
	}
	if len(payload) != 6 {
		return fmt.Errorf("%w: NAS security parameters to E-UTRA has %d value octets, want 6", ErrInvalidLength, len(payload))
	}
	var decoded NASSecurityParametersToEUTRAFields
	copy(decoded.Payload[:], payload)
	*f = decoded
	return nil
}
