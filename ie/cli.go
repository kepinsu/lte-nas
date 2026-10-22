package ie

import (
	"fmt"
	"io"
)

// CLIFields is the structured value of the Calling Line Identity IE
// (TS 24.301 clause 9.9.3.38). Text is encoded using the selected
// character encoding. AddCountryInitials asks the receiver to prepend
// country initials and a separator when displaying the number.
type CLIFields struct {
	Text               string
	Encoding           NetworkNameEncoding
	AddCountryInitials bool
}

// NewCLIWithFields creates a tagged CLI IE from fields.
func NewCLIWithFields(fields *CLIFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewCLI(payload), nil
}

// NewPositionalCLIWithFields creates a positional CLI IE from fields.
func NewPositionalCLIWithFields(fields *CLIFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalCLI(payload), nil
}

// CLI returns the IE as structured fields.
func (i *IE) CLI() (*CLIFields, error) {
	if err := i.requireDefinition(CLIDefinition); err != nil {
		return nil, err
	}
	return ParseCLIFields(i.Payload)
}

// ParseCLIFields decodes a CLI value without its IEI and length octet.
func ParseCLIFields(payload []byte) (*CLIFields, error) {
	fields := &CLIFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Text can be encoded with Encoding within the CLI
// maximum length of 12 characters.
func (f *CLIFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: CLI fields are nil", ErrInvalidValue)
	}
	text, err := f.marshalText()
	if err != nil {
		return err
	}
	if len(text)+1 > 12 {
		return fmt.Errorf("%w: CLI has %d value octets, want at most 12", ErrInvalidLength, len(text)+1)
	}
	return nil
}

// Marshal encodes fields into a CLI value without its IEI and length octet.
func (f *CLIFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, fmt.Errorf("%w: CLI fields are nil", ErrInvalidValue)
	}
	text, err := f.marshalText()
	if err != nil {
		return nil, err
	}
	if len(text)+1 > 12 {
		return nil, fmt.Errorf("%w: CLI has %d value octets, want at most 12", ErrInvalidLength, len(text)+1)
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
func (f *CLIFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the CLI value payload produced by Marshal.
// It returns zero for a nil receiver or an unrepresentable Text.
func (f *CLIFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	text, err := f.marshalText()
	if err != nil || len(text)+1 > 12 {
		return 0
	}
	return len(text) + 1
}

// UnmarshalBinary decodes a CLI value without its IEI and length octet.
// It leaves f unchanged when payload is invalid.
func (f *CLIFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: CLI fields are nil", ErrInvalidValue)
	}
	if len(payload) == 0 || len(payload) > 12 {
		return fmt.Errorf("%w: CLI has %d value octets, want 1..12", ErrInvalidLength, len(payload))
	}
	header := payload[0]
	if header&0x80 == 0 {
		return fmt.Errorf("%w: CLI extension bit is zero", ErrInvalidValue)
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
		return fmt.Errorf("%w: CLI coding scheme %d is reserved", ErrInvalidValue, encoding)
	}
	if err != nil {
		return err
	}
	*f = CLIFields{
		Text:               text,
		Encoding:           encoding,
		AddCountryInitials: header&0x08 != 0,
	}
	return nil
}

func (f *CLIFields) marshalText() ([]byte, error) {
	if f.Text == "" {
		return nil, fmt.Errorf("%w: CLI text is empty", ErrInvalidValue)
	}
	switch f.Encoding {
	case NetworkNameEncodingGSM7:
		return marshalNetworkNameGSM7(f.Text)
	case NetworkNameEncodingUCS2:
		return marshalNetworkNameUCS2(f.Text)
	default:
		return nil, fmt.Errorf("%w: CLI coding scheme %d is reserved", ErrInvalidValue, f.Encoding)
	}
}
