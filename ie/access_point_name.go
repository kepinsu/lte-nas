package ie

import (
	"fmt"
	"io"
	"strings"
)

const maximumAccessPointNameLength = 100

// AccessPointNameFields is the structured value of the Access point name IE
// (TS 24.301 clause 9.9.4.1). Name uses its familiar dotted representation,
// for example "internet.example"; the on-wire label lengths are handled by
// Marshal and UnmarshalBinary.
type AccessPointNameFields struct {
	Name string
}

// NewAccessPointNameWithFields creates a tagged Access point name IE from
// fields.
func NewAccessPointNameWithFields(fields *AccessPointNameFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewAccessPointName(payload), nil
}

// NewPositionalAccessPointNameWithFields creates a positional Access point
// name IE from fields.
func NewPositionalAccessPointNameWithFields(fields *AccessPointNameFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalAccessPointName(payload), nil
}

// AccessPointName returns the IE as structured fields.
func (i *IE) AccessPointName() (*AccessPointNameFields, error) {
	if err := i.requireDefinition(AccessPointNameDefinition); err != nil {
		return nil, err
	}
	return ParseAccessPointNameFields(i.Payload)
}

// ParseAccessPointNameFields decodes an Access point name value without its
// IEI and length octet.
func ParseAccessPointNameFields(payload []byte) (*AccessPointNameFields, error) {
	fields := &AccessPointNameFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Name can be encoded as a TS 24.008 APN.
func (f *AccessPointNameFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: access point name fields are nil", ErrInvalidValue)
	}
	if f.Name == "" {
		return fmt.Errorf("%w: access point name is empty", ErrInvalidValue)
	}
	labels := strings.Split(f.Name, ".")
	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 {
			return fmt.Errorf("%w: access point name label has %d octets, want 1..63", ErrInvalidValue, len(label))
		}
	}
	if encodedLength := len(f.Name) + 1; encodedLength > maximumAccessPointNameLength {
		return fmt.Errorf("%w: access point name has %d encoded octets, want at most %d", ErrInvalidLength, encodedLength, maximumAccessPointNameLength)
	}
	return nil
}

// Marshal encodes fields into an Access point name value without its IEI and
// length octet.
func (f *AccessPointNameFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	labels := strings.Split(f.Name, ".")
	payload := make([]byte, 0, len(f.Name)+1)
	for _, label := range labels {
		payload = append(payload, byte(len(label)))
		payload = append(payload, label...)
	}
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *AccessPointNameFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the Access point name value payload.
func (f *AccessPointNameFields) MarshalLen() int {
	if f == nil || f.Name == "" {
		return 0
	}
	return len(f.Name) + 1
}

// UnmarshalBinary decodes an Access point name value without its IEI and
// length octet. It leaves f unchanged when the payload is invalid.
func (f *AccessPointNameFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: access point name fields are nil", ErrInvalidValue)
	}
	if len(payload) == 0 || len(payload) > maximumAccessPointNameLength {
		return fmt.Errorf("%w: access point name has %d value octets, want 1..%d", ErrInvalidLength, len(payload), maximumAccessPointNameLength)
	}
	labels := make([]string, 0, 4)
	for offset := 0; offset < len(payload); {
		labelLength := int(payload[offset])
		offset++
		if labelLength == 0 || labelLength > 63 || labelLength > len(payload)-offset {
			return fmt.Errorf("%w: invalid access point name label length %d", ErrInvalidValue, labelLength)
		}
		labels = append(labels, string(payload[offset:offset+labelLength]))
		offset += labelLength
	}
	decoded := AccessPointNameFields{Name: strings.Join(labels, ".")}
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}
