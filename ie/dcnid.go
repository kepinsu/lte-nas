package ie

import (
	"fmt"
	"io"
)

// DCNIDFields is the structured value of the DCN-ID IE (TS 24.301 clause 9.9.3.48).
// DCNID is a 2-octet field identifying a core network domain.
type DCNIDFields struct {
	DCNID [2]byte
}

// NewPositionalDCNIDWithFields creates a positional DCN-ID IE from fields.
func NewPositionalDCNIDWithFields(fields *DCNIDFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalDCNID(payload), nil
}

// DCNID returns the IE as structured fields.
func (i *IE) DCNID() (*DCNIDFields, error) {
	if err := i.requireDefinition(DCNIDDefinition); err != nil {
		return nil, err
	}
	return ParseDCNIDFields(i.Payload)
}

// ParseDCNIDFields decodes a DCN-ID value without its IEI and length octet.
func ParseDCNIDFields(payload []byte) (*DCNIDFields, error) {
	fields := &DCNIDFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *DCNIDFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: DCN-ID fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a 2-octet DCN-ID value.
func (f *DCNIDFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.DCNID))
	copy(payload, f.DCNID[:])
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *DCNIDFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < f.MarshalLen() {
		return io.ErrShortBuffer
	}
	copy(b, f.DCNID[:])
	return nil
}

// MarshalLen returns the size of the DCN-ID value payload.
func (f *DCNIDFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 2
}

// UnmarshalBinary decodes a 2-octet DCN-ID value.
func (f *DCNIDFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: DCN-ID fields are nil", ErrInvalidValue)
	}
	if len(payload) != f.MarshalLen() {
		return fmt.Errorf("%w: DCN-ID has %d value octets, want 2", ErrInvalidLength, len(payload))
	}
	var decoded DCNIDFields
	copy(decoded.DCNID[:], payload)
	*f = decoded
	return nil
}
