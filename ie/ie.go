package ie

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"slices"
)

// IEFormat is an information element category defined in TS 24.007, clause 11.2.
type IEFormat uint8

// Information element categories from TS 24.007, clause 11.2.
const (
	// IEFormatOne is a half-octet V, or a one-octet TV with a four-bit IEI.
	IEFormatOne IEFormat = 1
	// IEFormatTwo is a tag-only (T) IE, not a tag followed by a value octet.
	IEFormatTwo IEFormat = 2
	// IEFormatThree is a fixed-length V or TV, including one-octet values.
	IEFormatThree IEFormat = 3
	// IEFormatFour is LV or TLV with a one-octet length indicator.
	IEFormatFour IEFormat = 4
	// IEFormatSix is LV-E or TLV-E with a two-octet big-endian length.
	IEFormatSix IEFormat = 6
)

// Kind is an on-wire NAS IEI. Tagged IEs always have Kind equal to their encoded
// tag. NoTag identifies a positional IE, which has no IEI on the wire.
type Kind uint8

// Codec errors can be inspected with errors.Is.
var (
	ErrInvalidLength = errors.New("invalid IE length")
	ErrInvalidFormat = errors.New("invalid IE format")
	ErrInvalidTag    = errors.New("invalid IE tag")
	ErrInvalidValue  = errors.New("invalid IE value")
	ErrUnknownIE     = errors.New("unknown IE")
	ErrNilIE         = errors.New("nil IE")
	ErrUnexpectedIE  = errors.New("unexpected IE")
	ErrAmbiguousIE   = errors.New("ambiguous tagged IE format")
)

// IE is an LTE NAS information element. Payload contains only the value, without
// its tag or length indicator. A type-1 value occupies Payload[0]'s low nibble.
//
// Kind is the actual encoded IEI for tagged IEs, or NoTag for positional IEs.
// Length records the payload length at construction/decoding time. Encoding
// always derives the length from Payload, so editing Payload cannot emit a stale
// length. Format is required only when constructing a raw IE directly; use a
// generated constructor whenever the IE has a TS 24.301 definition.
type IE struct {
	Kind       Kind           `json:"kind"`
	Length     uint16         `json:"length"`
	Payload    []byte         `json:"payload"`
	Format     IEFormat       `json:"format"`
	definition definitionKind `json:"-"`
}

type definition struct {
	name      string
	clause    string
	format    IEFormat
	minLength int
	maxLength int
}

func newTagged(definition definitionKind, kind Kind, payload []byte) *IE {
	i := newIE(definition, payload)
	i.Kind = kind
	return i
}

func newPositional(definition definitionKind, payload []byte) *IE {
	return newIE(definition, payload)
}

func newIE(definition definitionKind, payload []byte) *IE {
	d := _definitions[definition]
	i := &IE{
		Length:     payloadLength(payload),
		Payload:    slices.Clone(payload),
		Format:     d.format,
		definition: definition,
	}
	return i
}

// NewRaw creates an IE without a TS 24.301 value definition. It is intended for
// extensions unknown to this package. Raw type-3 decoding needs Length to be set
// explicitly before calling UnmarshalBinary or UnmarshalWithoutTag.
func NewRaw(kind Kind, format IEFormat, payload []byte) *IE {
	return &IE{
		Kind:    kind,
		Length:  payloadLength(payload),
		Payload: slices.Clone(payload),
		Format:  format,
	}
}

func payloadLength(payload []byte) uint16 {
	if len(payload) > math.MaxUint16 {
		return 0
	}
	return uint16(len(payload))
}

// Marshal encodes i. A nonzero Kind is always encoded as the IEI; a positional
// IE (Kind == NoTag) is encoded without an IEI.
func (i *IE) Marshal() ([]byte, error) {
	return i.marshal(i != nil && i.Kind != NoTag)
}

// MarshalBinary implements encoding.BinaryMarshaler, with the same mode as Marshal.
func (i *IE) MarshalBinary() ([]byte, error) {
	return i.Marshal()
}

// MarshalWithoutTag encodes the positional V, LV or LV-E representation without
// modifying i. A type-1 value is returned in the low nibble, with a zero high nibble.
// Use MarshalHalfOctets to pack two positional type-1 IEs into one octet.
func (i *IE) MarshalWithoutTag() ([]byte, error) {
	return i.marshal(false)
}

func (i *IE) marshal(tagged bool) ([]byte, error) {
	if err := i.validate(tagged); err != nil {
		return nil, err
	}
	b := make([]byte, i.marshalLen(tagged))
	i.marshalTo(b, tagged)
	return b, nil
}

// MarshalTo encodes into b using the same mode as Marshal. It returns
// io.ErrShortBuffer without writing anything if b is too short.
func (i *IE) MarshalTo(b []byte) error {
	return i.marshalToChecked(b, i != nil && i.Kind != NoTag)
}

// MarshalToWithoutTag writes the positional representation into b.
func (i *IE) MarshalToWithoutTag(b []byte) error {
	return i.marshalToChecked(b, false)
}

func (i *IE) marshalToChecked(b []byte, tagged bool) error {
	if err := i.validate(tagged); err != nil {
		return err
	}
	if len(b) < i.marshalLen(tagged) {
		return io.ErrShortBuffer
	}
	i.marshalTo(b, tagged)
	return nil
}

func (i *IE) marshalTo(b []byte, tagged bool) {
	if i.Format == IEFormatOne {
		b[0] = i.Payload[0]
		if tagged {
			b[0] |= byte(i.Kind)
		}
		return
	}
	// Copy the value before writing the header, so MarshalTo also works when
	// the destination and Payload overlap.
	offset := i.marshalLen(tagged) - len(i.Payload)
	copy(b[offset:], i.Payload)
	offset = 0
	if tagged {
		b[0] = byte(i.Kind)
		offset++
	}
	switch i.Format {
	case IEFormatFour:
		b[offset] = byte(len(i.Payload))
	case IEFormatSix:
		binary.BigEndian.PutUint16(b[offset:offset+2], uint16(len(i.Payload)))
	}
}

// MarshalLen returns the size of Marshal's output, or zero for a nil receiver.
// For a positional type-1 IE it returns one storage byte, not half an octet.
// It does not validate the payload; Marshal reports invalid values.
func (i *IE) MarshalLen() int {
	return i.marshalLen(i != nil && i.Kind != NoTag)
}

// MarshalLenWithoutTag returns the size of the positional representation.
func (i *IE) MarshalLenWithoutTag() int {
	return i.marshalLen(false)
}

func (i *IE) marshalLen(tagged bool) int {
	if i == nil {
		return 0
	}
	if i.Format == IEFormatOne {
		return 1
	}
	n := len(i.Payload)
	if tagged {
		n++
	}
	switch i.Format {
	case IEFormatFour:
		n++
	case IEFormatSix:
		n += 2
	}
	return n
}

// Parse decodes the first tagged IE with a globally unambiguous wire envelope.
// For a tag whose format depends on its message context, use the corresponding
// generated constructor and call UnmarshalBinary on it.
func Parse(b []byte, kind Kind) (*IE, error) {
	format, ok := _taggedFormats[kind]
	if !ok {
		if _, ambiguous := _ambiguousTaggedFormats[kind]; ambiguous {
			return nil, fmt.Errorf("%w: tag %#02x requires message context", ErrAmbiguousIE, kind)
		}
		return nil, fmt.Errorf("%w: tag %#02x", ErrUnknownIE, kind)
	}
	i := &IE{Kind: kind, Format: format.format, Length: format.length}
	if err := i.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return i, nil
}

// ParseWithoutTag decodes a positional IE matching prototype. The prototype is
// normally obtained from a NewPositional... constructor. For type-1 pairs use
// ParseHalfOctets instead of advancing by MarshalLen.
func ParseWithoutTag(b []byte, prototype *IE) (*IE, error) {
	if prototype == nil {
		return nil, ErrNilIE
	}
	i := *prototype
	i.Kind = NoTag
	i.Payload = nil
	if err := i.UnmarshalWithoutTag(b); err != nil {
		return nil, err
	}
	return &i, nil
}

// UnmarshalBinary decodes using i's Kind: NoTag selects positional decoding;
// otherwise the configured on-wire IEI must match. It copies the value and
// leaves i unchanged on failure.
func (i *IE) UnmarshalBinary(b []byte) error {
	return i.unmarshal(b, i != nil && i.Kind != NoTag)
}

// UnmarshalWithTag decodes a tagged IE. A nonzero Kind is the expected IEI;
// a zero Kind accepts and stores the received IEI.
func (i *IE) UnmarshalWithTag(b []byte) error {
	return i.unmarshal(b, true)
}

// UnmarshalWithoutTag decodes a positional IE and sets Kind to NoTag on
// success. For a type-1 IE only the low nibble is read; the high nibble belongs
// to another IE.
func (i *IE) UnmarshalWithoutTag(b []byte) error {
	return i.unmarshal(b, false)
}

func (i *IE) unmarshal(b []byte, tagged bool) error {
	if i == nil {
		return ErrNilIE
	}
	d, err := i.definitionFor()
	if err != nil {
		return err
	}
	decoded := IE{Format: d.format, definition: i.definition}
	offset := 0
	if tagged {
		if len(b) == 0 {
			return io.ErrUnexpectedEOF
		}
		decoded.Kind = Kind(b[0])
		if d.format == IEFormatOne {
			decoded.Kind &= 0xf0
		}
		if err := validateTag(decoded.Kind, d.format); err != nil {
			return err
		}
		if i.Kind != NoTag && i.Kind != decoded.Kind {
			return fmt.Errorf("%w: got %#02x, want %#02x", ErrInvalidTag, decoded.Kind, i.Kind)
		}
		if d.format != IEFormatOne {
			offset++
		}
	}
	n := d.minLength
	switch d.format {
	case IEFormatOne:
		if len(b) == 0 {
			return io.ErrUnexpectedEOF
		}
		decoded.Payload = []byte{b[0] & 0x0f}
		decoded.Length = 1
		if err := decoded.validate(tagged); err != nil {
			return err
		}
		*i = decoded
		return nil
	case IEFormatTwo:
		if !tagged {
			return fmt.Errorf("%w: type-2 IE requires a tag", ErrInvalidFormat)
		}
		n = 0
	case IEFormatFour:
		if len(b)-offset < 1 {
			return io.ErrUnexpectedEOF
		}
		n = int(b[offset])
		offset++
	case IEFormatSix:
		if len(b)-offset < 2 {
			return io.ErrUnexpectedEOF
		}
		n = int(binary.BigEndian.Uint16(b[offset : offset+2]))
		offset += 2
	}
	if n > len(b)-offset {
		return io.ErrUnexpectedEOF
	}
	if err := validateLength(d, i.definition, n); err != nil {
		return err
	}
	decoded.Payload = slices.Clone(b[offset : offset+n])
	decoded.Length = uint16(n)
	*i = decoded
	return nil
}

// Validate checks the IE envelope, payload length and type-1 nibble. It does not
// interpret compound payloads or check message-specific presence requirements.
func (i *IE) Validate() error {
	return i.validate(i != nil && i.Kind != NoTag)
}

func (i *IE) validate(tagged bool) error {
	if i == nil {
		return ErrNilIE
	}
	d, err := i.definitionFor()
	if err != nil {
		return err
	}
	if tagged {
		if err := validateTag(i.Kind, d.format); err != nil {
			return err
		}
	} else if d.format == IEFormatTwo {
		return fmt.Errorf("%w: type-2 IE requires a tag", ErrInvalidFormat)
	}
	if err := validateLength(d, i.definition, len(i.Payload)); err != nil {
		return err
	}
	if d.format == IEFormatOne && i.Payload[0] > 0x0f {
		return fmt.Errorf("%w: type-1 value %#02x exceeds a nibble", ErrInvalidValue, i.Payload[0])
	}
	return nil
}

func (i *IE) definitionFor() (definition, error) {
	if i.definition != unknownDefinition {
		if int(i.definition) >= len(_definitions) {
			return definition{}, fmt.Errorf("%w: invalid definition", ErrUnknownIE)
		}
		d := _definitions[i.definition]
		if i.Format != d.format {
			return definition{}, fmt.Errorf("%w: %s requires type %d", ErrInvalidFormat, d.name, d.format)
		}
		return d, nil
	}
	d := definition{format: i.Format, name: "raw IE"}
	switch i.Format {
	case IEFormatOne:
		d.minLength, d.maxLength = 1, 1
	case IEFormatTwo:
		d.minLength, d.maxLength = 0, 0
	case IEFormatThree:
		if i.Length == 0 {
			return definition{}, fmt.Errorf("%w: raw type-3 IE requires Length", ErrInvalidLength)
		}
		d.minLength, d.maxLength = int(i.Length), int(i.Length)
	case IEFormatFour:
		d.maxLength = math.MaxUint8
	case IEFormatSix:
		d.maxLength = math.MaxUint16
	default:
		return definition{}, fmt.Errorf("%w: %d", ErrInvalidFormat, i.Format)
	}
	return d, nil
}

func validateTag(kind Kind, format IEFormat) error {
	if kind == NoTag || (format == IEFormatOne && (kind < 0x80 || kind&0x0f != 0)) ||
		(format != IEFormatOne && kind >= 0x80) {
		return fmt.Errorf("%w: %#02x for type %d", ErrInvalidTag, kind, format)
	}
	return nil
}

func validateLength(d definition, kind definitionKind, n int) error {
	if n < d.minLength || n > d.maxLength {
		return fmt.Errorf("%w: %s value has %d octets, want %d..%d", ErrInvalidLength, d.name, n, d.minLength, d.maxLength)
	}
	// These optional extensions are complete groups, not arbitrary lengths
	// within the minimum/maximum range (TS 24.301 clauses 9.9.4.2 and 9.9.4.3).
	if (kind == APNAggregateMaximumBitRateDefinition && n%2 != 0) ||
		(kind == EPSQualityOfServiceDefinition && (n-1)%4 != 0) ||
		(kind == PLMNListDefinition && n%3 != 0) ||
		(kind == PDNAddressDefinition && n != 5 && n != 9 && n != 13) {
		return fmt.Errorf("%w: incomplete %s value (%d octets)", ErrInvalidLength, d.name, n)
	}
	return nil
}

// Name returns the TS 24.301 definition name for a generated IE, or "Unknown"
// for a raw IE parsed without message context.
func (i *IE) Name() string {
	if i != nil && i.definition != unknownDefinition && int(i.definition) < len(_definitions) {
		return _definitions[i.definition].name
	}
	return "Unknown"
}
