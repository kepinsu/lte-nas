package ie

import (
	"encoding/binary"
	"fmt"
	"io"
)

// ParseMultiIEs decodes only the tagged optional part of a message, using a
// private wire-format registry. Order and duplicates are preserved. Each returned
// IE has Kind equal to its actual on-wire IEI. msg owns message-specific meaning
// and validation, as well as parsing all positional IEs before calling this
// function.
//
// Unknown non-comprehension-required tags are skipped using TS 24.007 EPS wire
// rules. Unknown comprehension-required tags return ErrUnknownIE (TS 24.301
// clauses 7.5.1 and 7.6.1). Truncated IEs, even ignored ones, return
// io.ErrUnexpectedEOF rather than losing the next IE's boundary.
//
// Collision tagged IEs (0x32, 0x58, 0x5b, 0x5c) can appear in multiple formats.
// ParseMultiIEs resolves them by trying FormatFour (TLV) first, then falling back
// to FormatThree (TV). On failure for non-collision tags, ErrAmbiguousIE is
// returned; collision tags return the last parsing error encountered. Unknown
// non-comprehension-required tags are skipped; comprehension-required unknown
// tags return ErrUnknownIE. Truncated IEs return io.ErrUnexpectedEOF.
func ParseMultiIEs(b []byte) ([]*IE, error) {
	var ies []*IE
	for offset := 0; offset < len(b); {
		tag := b[offset]
		if tag >= 0x80 {
			tag &= 0xf0
		}
		kind := Kind(tag)

		// Collision tags: must use FormatFour (TLV). If TLV parsing fails,
		// the tag is in an unexpected format (Type-3 TV). Return ErrAmbiguousIE.
		if _, isAmbiguous := _ambiguousTaggedFormats[kind]; isAmbiguous {
			i4 := &IE{Kind: kind, Format: IEFormatFour, Length: 0}
			if err := i4.UnmarshalBinary(b[offset:]); err == nil && i4.Length <= uint16(len(b[offset:])-2) {
				ies = append(ies, i4)
				offset += i4.MarshalLen()
				continue
			}
			return nil, fmt.Errorf("%w: collision tag %#02x", ErrAmbiguousIE, b[offset])
		}

		format, known := _taggedFormats[kind]

		if !known {
			// Extension IEI (0x78..0x7F): 2-byte BE length. If truncated,
			// return io.ErrUnexpectedEOF (format is known; only data is
			// incomplete).  Well-formed data is decoded as opaque type-6 Raw
			// payload — message-level code validates comprehension.
			if tag >= 0x78 && tag <= 0x7f {
				if len(b[offset:]) < 3 {
					return nil, io.ErrUnexpectedEOF
				}
				payloadLen := int(binary.BigEndian.Uint16(b[offset+1 : offset+3]))
				if 3+payloadLen <= len(b[offset:]) {
					i := &IE{Kind: Kind(b[offset]), Format: IEFormatSix}
					i.Payload = make([]byte, payloadLen)
					copy(i.Payload, b[offset+3:offset+3+payloadLen])
					ies = append(ies, i)
					offset += 3 + payloadLen
					continue
				}
			}
			// For everything else, skip via skipUnknownIE.
			n, err := skipUnknownIE(b[offset:])
			if err != nil {
				return nil, fmt.Errorf("IE at offset %d: %w", offset, err)
			}
			offset += n
			continue
		}

		i := &IE{Kind: kind, Format: format.format, Length: format.length}

		// High-bit pass-through tags (kind 0x80–0xf0, IEFormatOne) can be
		// either type-1 TV (single byte: tag | value) or type-4 TLV
		// (tag + length + data). When more than one byte is available, try
		// TLV parsing first — if the data is actually type-1 TV followed by
		// another IE, the length byte (next IE tag) is usually ≥ 0x10,
		// making the TLV total implausibly large and causing a truncated-
		// buffer error, at which point we fall back to type-1 TV.
		if tag >= 0x80 && format.format == IEFormatOne && len(b[offset:]) >= 2 {
			tlvLen := int(b[offset+1])
			tlvTotal := 2 + tlvLen
			if tlvTotal <= len(b[offset:]) {
				// Read the full TLV directly (bypassing validation for unknown tags)
				i = &IE{Kind: Kind(b[offset]), Format: IEFormatFour, Length: uint16(tlvLen)}
				i.Payload = make([]byte, tlvLen)
				copy(i.Payload, b[offset+2:offset+tlvTotal])
				ies = append(ies, i)
				offset += tlvTotal
				continue
			}
		}

		if err := i.UnmarshalBinary(b[offset:]); err != nil {
			return nil, fmt.Errorf("IE at offset %d: %w", offset, err)
		}
		ies = append(ies, i)
		offset += i.MarshalLen()
	}
	return ies, nil
}

// IsAmbiguousTagged reports whether kind is a collision tag that can appear in
// multiple message types with different encodings (TS 24.301 collision tags 0x32,
// 0x58, 0x5b, 0x5c). These tags must use FormatFour (TLV) wire encoding in
// tagged-optional parts.
func IsAmbiguousTagged(kind Kind) bool {
	_, ok := _ambiguousTaggedFormats[kind]
	return ok
}

// skipUnknownIE implements TS 24.007 V18.1.0 clauses 11.2.4 and 11.2.5.
// Comprehension-required tags (0x00..0x0F) return ErrUnknownIE.
// Type-6 IEIs (0x78..0x7F) are skipped by reading a 2-byte BE length (3+length).
// High-bit tags (0x80+) consume just 1 byte. Everything else consumes 2 bytes
// (tag + 1 payload octet for type-3 TV format), except TLV with nibble != 0
// which reads a length byte (2+length). Truncated type-6 data returns io.ErrUnexpectedEOF.
func skipUnknownIE(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, io.ErrUnexpectedEOF
	}
	tag := b[0]

	// Extension IEI (0x78..0x7F): 2-byte BE length. Truncated data returns
	// io.ErrUnexpectedEOF; well-formed data is skipped (3+length bytes).
	if tag >= 0x78 && tag <= 0x7f {
		if len(b) < 3 {
			return 0, io.ErrUnexpectedEOF
		}
		length := int(binary.BigEndian.Uint16(b[1:3]))
		total := 3 + length
		if total > len(b) {
			return 0, io.ErrUnexpectedEOF
		}
		return total, nil
	}

	if tag&0xf0 == 0 {
		// Type-4 unknown tags: nibble 0..F requires comprehension.
		return 0, ErrUnknownIE
	}

	if tag&0x80 != 0 {
		// High-bit tags: consume just the tag byte (1 octet).
		return 1, nil
	}

	// Non-high-bit type-4 (nibble != 0): TLV with length byte.
	// Tag 0x12 has nibble 1 (type-4 IEI), so read length byte.
	if len(b) < 2 {
		return 0, io.ErrUnexpectedEOF
	}
	length := int(b[1])
	total := 2 + length
	if total > len(b) {
		return 0, io.ErrUnexpectedEOF
	}
	return total, nil
}

// MarshalHalfOctets packs two positional type-1 values into one octet. The first
// argument occupies bits 1..4, the second bits 5..8. Tags on the inputs are ignored.
func MarshalHalfOctets(low, high *IE) ([]byte, error) {
	for _, i := range []*IE{low, high} {
		if err := i.validate(false); err != nil {
			return nil, err
		}
		if i.Format != IEFormatOne {
			return nil, fmt.Errorf("%w: expected type-1 IE", ErrInvalidFormat)
		}
	}
	return []byte{high.Payload[0]<<4 | low.Payload[0]}, nil
}

// ParseHalfOctets decodes the low and high positional type-1 IEs of b's first
// octet. low and high are prototypes normally obtained from NewPositional...
// constructors. Exactly one octet is consumed, even though two IE values are
// returned.
func ParseHalfOctets(b []byte, low, high *IE) (*IE, *IE, error) {
	for _, prototype := range []*IE{low, high} {
		if prototype == nil {
			return nil, nil, ErrNilIE
		}
		d, err := prototype.definitionFor()
		if err != nil {
			return nil, nil, err
		}
		if d.format != IEFormatOne {
			return nil, nil, fmt.Errorf("%w: %s is not type 1", ErrInvalidFormat, d.name)
		}
	}
	if len(b) == 0 {
		return nil, nil, io.ErrUnexpectedEOF
	}
	lowValue, err := ParseWithoutTag(b, low)
	if err != nil {
		return nil, nil, err
	}
	highValue, err := ParseWithoutTag([]byte{b[0] >> 4}, high)
	if err != nil {
		return nil, nil, err
	}
	return lowValue, highValue, nil
}
