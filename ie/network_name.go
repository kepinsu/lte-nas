package ie

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

// NetworkNameEncoding identifies the character encoding of a Network name.
type NetworkNameEncoding uint8

// Network name encodings defined by TS 24.008 clause 10.5.3.5a.
const (
	NetworkNameEncodingGSM7 NetworkNameEncoding = iota
	NetworkNameEncodingUCS2
)

// NetworkNameFields is the structured value of the Network name IE (TS 24.301
// clause 9.9.3.24 and TS 24.008 clause 10.5.3.5a). Text is encoded using the
// selected character encoding. AddCountryInitials asks the receiver to prepend
// country initials and a separator when displaying the name.
type NetworkNameFields struct {
	Text               string
	Encoding           NetworkNameEncoding
	AddCountryInitials bool
}

// NewFullNetworkNameWithFields creates a tagged full Network name IE from
// fields.
func NewFullNetworkNameWithFields(fields *NetworkNameFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewFullNetworkName(payload), nil
}

// NewShortNetworkNameWithFields creates a tagged short Network name IE from
// fields.
func NewShortNetworkNameWithFields(fields *NetworkNameFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewShortNetworkName(payload), nil
}

// NewPositionalNetworkNameWithFields creates a positional Network name IE
// from fields.
func NewPositionalNetworkNameWithFields(fields *NetworkNameFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalNetworkName(payload), nil
}

// NetworkName returns the IE as structured fields.
func (i *IE) NetworkName() (*NetworkNameFields, error) {
	if err := i.requireDefinition(NetworkNameDefinition); err != nil {
		return nil, err
	}
	return ParseNetworkNameFields(i.Payload)
}

// ParseNetworkNameFields decodes a Network name value without its IEI and
// length octet.
func ParseNetworkNameFields(payload []byte) (*NetworkNameFields, error) {
	fields := &NetworkNameFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Text can be encoded with Encoding within the Network
// name maximum length.
func (f *NetworkNameFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: network name fields are nil", ErrInvalidValue)
	}
	text, err := f.marshalText()
	if err != nil {
		return err
	}
	if len(text)+1 > 255 {
		return fmt.Errorf("%w: network name has %d value octets, want at most 255", ErrInvalidLength, len(text)+1)
	}
	return nil
}

// Marshal encodes fields into a Network name value without its IEI and length
// octet.
func (f *NetworkNameFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, fmt.Errorf("%w: network name fields are nil", ErrInvalidValue)
	}
	text, err := f.marshalText()
	if err != nil {
		return nil, err
	}
	if len(text)+1 > 255 {
		return nil, fmt.Errorf("%w: network name has %d value octets, want at most 255", ErrInvalidLength, len(text)+1)
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
func (f *NetworkNameFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the Network name value payload produced by
// Marshal. It returns zero for a nil receiver or an unrepresentable Text.
func (f *NetworkNameFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	text, err := f.marshalText()
	if err != nil || len(text)+1 > 255 {
		return 0
	}
	return len(text) + 1
}

// UnmarshalBinary decodes a Network name value without its IEI and length
// octet. It leaves f unchanged when payload is invalid.
func (f *NetworkNameFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: network name fields are nil", ErrInvalidValue)
	}
	if len(payload) == 0 || len(payload) > 255 {
		return fmt.Errorf("%w: network name has %d value octets, want 1..255", ErrInvalidLength, len(payload))
	}
	header := payload[0]
	if header&0x80 == 0 {
		return fmt.Errorf("%w: network name extension bit is zero", ErrInvalidValue)
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
		return fmt.Errorf("%w: network name coding scheme %d is reserved", ErrInvalidValue, encoding)
	}
	if err != nil {
		return err
	}
	*f = NetworkNameFields{
		Text:               text,
		Encoding:           encoding,
		AddCountryInitials: header&0x08 != 0,
	}
	return nil
}

func (f *NetworkNameFields) marshalText() ([]byte, error) {
	if !utf8.ValidString(f.Text) {
		return nil, fmt.Errorf("%w: network name text is not valid UTF-8", ErrInvalidValue)
	}
	switch f.Encoding {
	case NetworkNameEncodingGSM7:
		return marshalNetworkNameGSM7(f.Text)
	case NetworkNameEncodingUCS2:
		return marshalNetworkNameUCS2(f.Text)
	default:
		return nil, fmt.Errorf("%w: network name coding scheme %d is reserved", ErrInvalidValue, f.Encoding)
	}
}

func networkNameGSM7SpareBits(text string) (byte, error) {
	septetCount := 0
	for _, character := range text {
		encoded, err := encodeNetworkNameGSM7Character(character)
		if err != nil {
			return 0, err
		}
		septetCount += len(encoded)
	}
	if septetCount == 0 {
		return 0, nil
	}
	return byte((8 - (septetCount*7)%8) % 8), nil
}

func marshalNetworkNameGSM7(text string) ([]byte, error) {
	septets := make([]byte, 0, len(text))
	for _, character := range text {
		encoded, err := encodeNetworkNameGSM7Character(character)
		if err != nil {
			return nil, err
		}
		septets = append(septets, encoded...)
	}
	if len(septets) == 0 {
		return nil, nil
	}
	payload := make([]byte, (len(septets)*7+7)/8)
	for index, septet := range septets {
		bitOffset := index * 7
		octetOffset := bitOffset / 8
		shift := bitOffset % 8
		payload[octetOffset] |= septet << shift
		if shift > 1 {
			payload[octetOffset+1] |= septet >> (8 - shift)
		}
	}
	return payload, nil
}

func encodeNetworkNameGSM7Character(character rune) ([]byte, error) {
	for code, value := range networkNameGSM7DefaultAlphabet {
		if byte(code) != 0x1b && value == character {
			return []byte{byte(code)}, nil
		}
	}
	switch character {
	case '\f':
		return []byte{0x1b, 0x0a}, nil
	case '^':
		return []byte{0x1b, 0x14}, nil
	case '{':
		return []byte{0x1b, 0x28}, nil
	case '}':
		return []byte{0x1b, 0x29}, nil
	case '\\':
		return []byte{0x1b, 0x2f}, nil
	case '[':
		return []byte{0x1b, 0x3c}, nil
	case '~':
		return []byte{0x1b, 0x3d}, nil
	case ']':
		return []byte{0x1b, 0x3e}, nil
	case '|':
		return []byte{0x1b, 0x40}, nil
	case '€':
		return []byte{0x1b, 0x65}, nil
	default:
		return nil, fmt.Errorf("%w: network name character %q is not in the GSM 7-bit alphabet", ErrInvalidValue, character)
	}
}

func unmarshalNetworkNameGSM7(payload []byte, spareBits byte) (string, error) {
	if spareBits > 7 {
		return "", fmt.Errorf("%w: network name spare-bit count %d exceeds seven", ErrInvalidValue, spareBits)
	}
	bitLength := len(payload) * 8
	if spareBits != 0 {
		if len(payload) == 0 || bitLength < int(spareBits) {
			return "", fmt.Errorf("%w: network name GSM 7-bit text has invalid spare bits", ErrInvalidLength)
		}
		paddingMask := byte(0xff << (8 - spareBits))
		if payload[len(payload)-1]&paddingMask != 0 {
			return "", fmt.Errorf("%w: network name GSM 7-bit spare bits are nonzero", ErrInvalidValue)
		}
		bitLength -= int(spareBits)
		if bitLength%7 != 0 {
			return "", fmt.Errorf("%w: network name GSM 7-bit text has inconsistent spare bits", ErrInvalidLength)
		}
	}
	septetCount := bitLength / 7
	var text strings.Builder
	for index := 0; index < septetCount; index++ {
		septet := networkNameGSM7Septet(payload, index)
		if septet == 0x1b {
			index++
			if index >= septetCount {
				return "", fmt.Errorf("%w: network name GSM 7-bit extension escape is truncated", ErrInvalidLength)
			}
			septet = networkNameGSM7Septet(payload, index)
			character, ok := decodeNetworkNameGSM7Extension(septet)
			if !ok {
				return "", fmt.Errorf("%w: network name GSM 7-bit extension %#02x is reserved", ErrInvalidValue, septet)
			}
			text.WriteRune(character)
			continue
		}
		text.WriteRune(networkNameGSM7DefaultAlphabet[septet])
	}
	return text.String(), nil
}

func networkNameGSM7Septet(payload []byte, index int) byte {
	bitOffset := index * 7
	octetOffset := bitOffset / 8
	shift := bitOffset % 8
	septet := payload[octetOffset] >> shift
	if shift > 1 && octetOffset+1 < len(payload) {
		septet |= payload[octetOffset+1] << (8 - shift)
	}
	return septet & 0x7f
}

func decodeNetworkNameGSM7Extension(septet byte) (rune, bool) {
	switch septet {
	case 0x0a:
		return '\f', true
	case 0x14:
		return '^', true
	case 0x28:
		return '{', true
	case 0x29:
		return '}', true
	case 0x2f:
		return '\\', true
	case 0x3c:
		return '[', true
	case 0x3d:
		return '~', true
	case 0x3e:
		return ']', true
	case 0x40:
		return '|', true
	case 0x65:
		return '€', true
	default:
		return 0, false
	}
}

func marshalNetworkNameUCS2(text string) ([]byte, error) {
	payload := make([]byte, 0, len(text)*2)
	for _, character := range text {
		if character > 0xffff {
			return nil, fmt.Errorf("%w: network name character %q is outside UCS2", ErrInvalidValue, character)
		}
		payload = append(payload, byte(character>>8), byte(character))
	}
	return payload, nil
}

func unmarshalNetworkNameUCS2(payload []byte) (string, error) {
	if len(payload)%2 != 0 {
		return "", fmt.Errorf("%w: network name UCS2 text has an odd number of octets", ErrInvalidLength)
	}
	var text strings.Builder
	for offset := 0; offset < len(payload); offset += 2 {
		character := rune(payload[offset])<<8 | rune(payload[offset+1])
		if character >= 0xd800 && character <= 0xdfff {
			return "", fmt.Errorf("%w: network name UCS2 text contains a surrogate code point", ErrInvalidValue)
		}
		text.WriteRune(character)
	}
	return text.String(), nil
}

var networkNameGSM7DefaultAlphabet = [...]rune{
	'@', '£', '$', '¥', 'è', 'é', 'ù', 'ì', 'ò', 'Ç', '\n', 'Ø', 'ø', '\r', 'Å', 'å',
	'Δ', '_', 'Φ', 'Γ', 'Λ', 'Ω', 'Π', 'Ψ', 'Σ', 'Θ', 'Ξ', 0, 'Æ', 'æ', 'ß', 'É',
	' ', '!', '"', '#', '¤', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ':', ';', '<', '=', '>', '?',
	'¡', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O',
	'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'Ä', 'Ö', 'Ñ', 'Ü', '§',
	'¿', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
	'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'ä', 'ö', 'ñ', 'ü', 'à',
}
