package ie

import (
	"fmt"
	"io"
)

// UEAdditionalSecurityCapabilityFields is the structured value of the UE
// additional security capability IE (TS 24.301 clause 9.9.3.53). The four-octet
// payload contains two 32-bit bitmaps: encryption algorithm (EA) capabilities
// in the first octet and integrity protection algorithm (IA) capabilities in the
// next three octets. See clause 125.93 for bit definitions.
type UEAdditionalSecurityCapabilityFields struct {
	// Payload contains the 4-octet bitmap. Byte 0: EEA0–EEA3 + spare bits.
	// Bytes 1-3: IA bitmaps (IAA0–IAA2 etc.).
	Payload [4]byte
}

// NewUEAdditionalSecurityCapabilityWithFields creates a tagged UE additional
// security capability IE from fields.
func NewUEAdditionalSecurityCapabilityWithFields(fields *UEAdditionalSecurityCapabilityFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewUEAdditionalSecurityCapability(payload), nil
}

// NewPositionalUEAdditionalSecurityCapabilityWithFields creates a positional UE
// additional security capability IE from fields.
func NewPositionalUEAdditionalSecurityCapabilityWithFields(fields *UEAdditionalSecurityCapabilityFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalUEAdditionalSecurityCapability(payload), nil
}

// UEAdditionalSecurityCapability returns the IE as structured fields.
func (i *IE) UEAdditionalSecurityCapability() (*UEAdditionalSecurityCapabilityFields, error) {
	if err := i.requireDefinition(UEAdditionalSecurityCapabilityDefinition); err != nil {
		return nil, err
	}
	return ParseUEAdditionalSecurityCapabilityFields(i.Payload)
}

// ParseUEAdditionalSecurityCapabilityFields decodes a UE additional security
// capability value without its IEI and length octet.
func ParseUEAdditionalSecurityCapabilityFields(payload []byte) (*UEAdditionalSecurityCapabilityFields, error) {
	fields := &UEAdditionalSecurityCapabilityFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *UEAdditionalSecurityCapabilityFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: UE additional security capability fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a 4-octet UE additional security capability value.
func (f *UEAdditionalSecurityCapabilityFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Payload))
	copy(payload, f.Payload[:])
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *UEAdditionalSecurityCapabilityFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 4 {
		return io.ErrShortBuffer
	}
	copy(b, f.Payload[:])
	return nil
}

// MarshalLen returns the size of the UE additional security capability value payload.
func (f *UEAdditionalSecurityCapabilityFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 4
}

// UnmarshalBinary decodes a 4-octet UE additional security capability value.
func (f *UEAdditionalSecurityCapabilityFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: UE additional security capability fields are nil", ErrInvalidValue)
	}
	if len(payload) != 4 {
		return fmt.Errorf("%w: UE additional security capability has %d value octets, want 4", ErrInvalidLength, len(payload))
	}
	var decoded UEAdditionalSecurityCapabilityFields
	copy(decoded.Payload[:], payload)
	*f = decoded
	return nil
}
