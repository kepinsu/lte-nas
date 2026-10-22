package ie

import (
	"fmt"
	"io"
)

// WUSAssistanceInformationFields is the structured value of the tagged WUS
// assistance information IE (TS 24.301 clause 9.9.3.62). It carries an opaque
// byte string between 1 and 255 octets.
type WUSAssistanceInformationFields struct {
	Value []byte
}

// NewWUSAssistanceInformationWithFields creates a tagged WUS assistance information
// IE from fields.
func NewWUSAssistanceInformationWithFields(fields *WUSAssistanceInformationFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewWUSAssistanceInformation(payload), nil
}

// NewPositionalWUSAssistanceInformationWithFields creates a positional WUS
// assistance information IE from fields.
func NewPositionalWUSAssistanceInformationWithFields(fields *WUSAssistanceInformationFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalWUSAssistanceInformation(payload), nil
}

// WUSAssistanceInformation returns the IE as structured fields.
func (i *IE) WUSAssistanceInformation() (*WUSAssistanceInformationFields, error) {
	if err := i.requireDefinition(WUSAssistanceInformationDefinition); err != nil {
		return nil, err
	}
	return ParseWUSAssistanceInformationFields(i.Payload)
}

// ParseWUSAssistanceInformationFields decodes a WUS assistance information value
// without its IEI and length octet.
func ParseWUSAssistanceInformationFields(payload []byte) (*WUSAssistanceInformationFields, error) {
	fields := &WUSAssistanceInformationFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the payload is within [1, 255] octets.
func (f *WUSAssistanceInformationFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: WUS assistance information fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 1 {
		return fmt.Errorf("%w: WUS assistance information has %d value octets, want >= 1", ErrInvalidLength, n)
	}
	if n > 255 {
		return fmt.Errorf("%w: WUS assistance information has %d value octets, want <= 255", ErrInvalidLength, n)
	}
	return nil
}

// Marshal encodes fields into a WUS assistance information value without its IEI
// and length octet.
func (f *WUSAssistanceInformationFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *WUSAssistanceInformationFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the WUS assistance information value payload, or
// zero for a nil receiver or nil Value.
func (f *WUSAssistanceInformationFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes a WUS assistance information value without its IEI and
// length octet. It copies the payload into a temporary struct first, leaving f
// unchanged on failure.
func (f *WUSAssistanceInformationFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: WUS assistance information fields are nil", ErrInvalidValue)
	}
	if len(payload) < 1 {
		return fmt.Errorf("%w: WUS assistance information has %d value octets, want >= 1", ErrInvalidLength, len(payload))
	}
	if len(payload) > 255 {
		return fmt.Errorf("%w: WUS assistance information has %d value octets, want <= 255", ErrInvalidLength, len(payload))
	}
	var decoded WUSAssistanceInformationFields
	decoded.Value = make([]byte, len(payload))
	copy(decoded.Value, payload)
	*f = decoded
	return nil
}
