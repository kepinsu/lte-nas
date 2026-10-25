package ie

import (
	"fmt"
	"io"
)

// NASSecurityParametersFromEUTRAFields is the structured value of the NAS
// security parameters from E-UTRA IE (TS 24.301 clause 9.9.2.6).
type NASSecurityParametersFromEUTRAFields struct {
	// KSIASME is the Key Set Identifier for NAS security from E-UTRA.
	KSIASME uint8
}

// NewNASSecurityParametersFromEUTRAWithFields creates a tagged NAS security
// parameters from E-UTRA IE from fields.
func NewNASSecurityParametersFromEUTRAWithFields(fields *NASSecurityParametersFromEUTRAFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalNASSecurityParametersFromEUTRA(payload[0]), nil
}

// NewPositionalNASSecurityParametersFromEUTRAWithFields creates a positional NAS
// security parameters from E-UTRA IE from fields.
func NewPositionalNASSecurityParametersFromEUTRAWithFields(fields *NASSecurityParametersFromEUTRAFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalNASSecurityParametersFromEUTRA(payload[0]), nil
}

// NASSecurityParametersFromEUTRA returns the IE as structured fields.
func (i *IE) NASSecurityParametersFromEUTRA() (*NASSecurityParametersFromEUTRAFields, error) {
	if err := i.requireDefinition(NASSecurityParametersFromEUTRADefinition); err != nil {
		return nil, err
	}
	return ParseNASSecurityParametersFromEUTRAFields(i.Payload)
}

// ParseNASSecurityParametersFromEUTRAFields decodes a NAS security parameters
// from E-UTRA value without its IEI.
func ParseNASSecurityParametersFromEUTRAFields(payload []byte) (*NASSecurityParametersFromEUTRAFields, error) {
	fields := &NASSecurityParametersFromEUTRAFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that KSIASME fits three bits.
func (f *NASSecurityParametersFromEUTRAFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: NAS security parameters from E-UTRA fields are nil", ErrInvalidValue)
	}
	if f.KSIASME > 0x07 {
		return fmt.Errorf("%w: KSI-ASME exceeds its encoded width", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a NAS security parameters from E-UTRA value,
// without its IEI.
func (f *NASSecurityParametersFromEUTRAFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{f.KSIASME}, nil
}

// MarshalTo encodes fields into b.
func (f *NASSecurityParametersFromEUTRAFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the NAS security parameters from E-UTRA value
// payload.
func (f *NASSecurityParametersFromEUTRAFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a NAS security parameters from E-UTRA value without its
// IEI.
func (f *NASSecurityParametersFromEUTRAFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: NAS security parameters from E-UTRA fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: NAS security parameters from E-UTRA has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xf8 != 0 {
		return fmt.Errorf("%w: NAS security parameters from E-UTRA spare bits are nonzero", ErrInvalidValue)
	}
	*f = NASSecurityParametersFromEUTRAFields{
		KSIASME: payload[0] & 0x07,
	}
	return nil
}
