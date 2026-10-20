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
// A tag with several known wire formats returns ErrAmbiguousIE: even optional
// NAS IEIs can collide between messages. Such tags must be decoded by msg with
// Parse and the expected Kind. On failure, no partial list is returned.
// Payloads are copied from b; skipped IEs are not included in the result.
func ParseMultiIEs(b []byte) ([]*IE, error) {
	var ies []*IE
	for offset := 0; offset < len(b); {
		tag := b[offset]
		if tag >= 0x80 {
			tag &= 0xf0
		}
		kind := Kind(tag)
		format, known := _taggedFormats[kind]
		if !known {
			if _, ambiguous := _ambiguousTaggedFormats[kind]; ambiguous {
				return nil, fmt.Errorf("%w: tag %#02x at offset %d requires message context", ErrAmbiguousIE, tag, offset)
			}
			n, err := skipUnknownIE(b[offset:])
			if err != nil {
				return nil, fmt.Errorf("IE at offset %d: %w", offset, err)
			}
			offset += n
			continue
		}
		i := &IE{Kind: kind, Format: format.format, Length: format.length}
		if err := i.UnmarshalBinary(b[offset:]); err != nil {
			return nil, fmt.Errorf("IE at offset %d: %w", offset, err)
		}
		ies = append(ies, i)
		offset += i.MarshalLen()
	}
	return ies, nil
}

// skipUnknownIE implements TS 24.007 V18.1.0 clauses 11.2.4 and 11.2.5.
func skipUnknownIE(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, io.ErrUnexpectedEOF
	}
	tag := b[0]
	// Type-4 IEIs 00..0F and type-6 IEIs 7E..7F require comprehension.
	// 7C/7D were also comprehension required in older specification versions,
	// but are explicitly excluded from that range in Release 18.
	if tag&0xf0 == 0 || tag&0xfe == 0x7e {
		return 0, fmt.Errorf("%w: comprehension-required tag %#02x", ErrUnknownIE, tag)
	}
	if tag&0x80 != 0 {
		return 1, nil
	}
	header, length := 2, 0
	if tag&0x78 == 0x78 {
		header = 3
		if len(b) < header {
			return 0, io.ErrUnexpectedEOF
		}
		length = int(binary.BigEndian.Uint16(b[1:3]))
	} else {
		if len(b) < header {
			return 0, io.ErrUnexpectedEOF
		}
		length = int(b[1])
	}
	if length > len(b)-header {
		return 0, io.ErrUnexpectedEOF
	}
	return header + length, nil
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
