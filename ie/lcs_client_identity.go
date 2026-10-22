package ie

import (
	"fmt"
	"io"
)

// LCSClientIdentityFields is the structured value of the LCS client identity IE
// (TS 24.301 clause 9.9.3.41). Text is encoded using the selected character
// encoding. AddCountryInitials asks the receiver to prepend country initials and
// a separator when displaying the name.
type LCSClientIdentityFields struct {
	Text               string
	Encoding           NetworkNameEncoding
	AddCountryInitials bool
}

// NewLCSClientIdentityWithFields creates a tagged LCS client identity IE from
// fields.
func NewLCSClientIdentityWithFields(fields *LCSClientIdentityFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewLCSClientIdentity(payload), nil
}

// NewPositionalLCSClientIdentityWithFields creates a positional LCS client
// identity IE from fields.
func NewPositionalLCSClientIdentityWithFields(fields *LCSClientIdentityFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalLCSClientIdentity(payload), nil
}

// LCSClientIdentity returns the IE as structured fields.
func (i *IE) LCSClientIdentity() (*LCSClientIdentityFields, error) {
	if err := i.requireDefinition(LCSClientIdentityDefinition); err != nil {
		return nil, err
	}
	return ParseLCSClientIdentityFields(i.Payload)
}

// ParseLCSClientIdentityFields decodes an LCS client identity value without its
// IEI and length octet.
func ParseLCSClientIdentityFields(payload []byte) (*LCSClientIdentityFields, error) {
	fields := &LCSClientIdentityFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Text can be encoded with Encoding within the LCS client
// identity maximum length of 255 value octets.
func (f *LCSClientIdentityFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: LCS client identity fields are nil", ErrInvalidValue)
	}
	text, err := f.marshalText()
	if err != nil {
		return err
	}
	if len(text)+1 > 255 {
		return fmt.Errorf("%w: LCS client identity has %d value octets, want at most 255", ErrInvalidLength, len(text)+1)
	}
	return nil
}

// Marshal encodes fields into an LCS client identity value without its IEI and
// length octet.
func (f *LCSClientIdentityFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, fmt.Errorf("%w: LCS client identity fields are nil", ErrInvalidValue)
	}
	text, err := f.marshalText()
	if err != nil {
		return nil, err
	}
	if len(text)+1 > 255 {
		return nil, fmt.Errorf("%w: LCS client identity has %d value octets, want at most 255", ErrInvalidLength, len(text)+1)
	}
	header := byte(0x80 | byte(f.Encoding)<<4)
	if f.AddCountryInitials {
		header |= 0x08
	}
	if f.Encoding == NetworkNameEncodingGSM7 {
		spareBits, err := networkNameGSM7SpareBits(f.Text)
		if err != nil {
			return nil, err
		}
		header |= spareBits
	}
	return append([]byte{header}, text...), nil
}

// MarshalTo encodes fields into b.
func (f *LCSClientIdentityFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the LCS client identity value payload produced
// by Marshal. It returns zero for a nil receiver or an unrepresentable Text.
func (f *LCSClientIdentityFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	text, err := f.marshalText()
	if err != nil || len(text)+1 > 255 {
		return 0
	}
	return len(text) + 1
}

// UnmarshalBinary decodes an LCS client identity value without its IEI and length
// octet. It leaves f unchanged when payload is invalid.
func (f *LCSClientIdentityFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: LCS client identity fields are nil", ErrInvalidValue)
	}
	if len(payload) == 0 || len(payload) > 255 {
		return fmt.Errorf("%w: LCS client identity has %d value octets, want 1..255", ErrInvalidLength, len(payload))
	}
	header := payload[0]
	if header&0x80 == 0 {
		return fmt.Errorf("%w: LCS client identity extension bit is zero", ErrInvalidValue)
	}
	encoding := NetworkNameEncoding((header >> 4) & 0x07)
	var text string
	var err error
	switch encoding {
	case NetworkNameEncodingGSM7:
		text, err = unmarshalNetworkNameGSM7(payload[1:], header&0x07)
	case NetworkNameEncodingUCS2:
		text, err = unmarshalNetworkNameUCS2(payload[1:])
	default:
		return fmt.Errorf("%w: LCS client identity coding scheme %d is reserved", ErrInvalidValue, encoding)
	}
	if err != nil {
		return err
	}
	*f = LCSClientIdentityFields{
		Text:               text,
		Encoding:           encoding,
		AddCountryInitials: header&0x08 != 0,
	}
	return nil
}

func (f *LCSClientIdentityFields) marshalText() ([]byte, error) {
	if f.Text == "" {
		return nil, fmt.Errorf("%w: LCS client identity text is empty", ErrInvalidValue)
	}
	switch f.Encoding {
	case NetworkNameEncodingGSM7:
		return marshalNetworkNameGSM7(f.Text)
	case NetworkNameEncodingUCS2:
		return marshalNetworkNameUCS2(f.Text)
	default:
		return nil, fmt.Errorf("%w: LCS client identity coding scheme %d is reserved", ErrInvalidValue, f.Encoding)
	}
}
