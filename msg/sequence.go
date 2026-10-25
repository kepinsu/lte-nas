package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// writePassThroughTagged writes a non-catalogue tagged IE into b at offset off,
// returning the bytes consumed. It handles tags that are not in the ie package
// catalogue by encoding them directly: for type-1, the Kind byte is written as
// the tag nibble with Payload[0] OR'd in; for type-4, the full tag byte +
// length byte + payload are written. Returns io.ErrShortBuffer if b is too short.
func writePassThroughTagged(b []byte, off int, i *ie.IE) (int, error) {
	if i == nil {
		return 0, nil
	}

	n := i.MarshalLen()
	if len(b)-off < n {
		return 0, io.ErrShortBuffer
	}

	switch i.Format {
	case ie.IEFormatOne: // Type-1 TV: single wire byte = payload | tag_nibble
		if len(i.Payload) == 0 {
			return 0, fmt.Errorf("%w: pass-through type-1 IE has nil payload", ie.ErrNilIE)
		}
		b[off] = i.Payload[0] | byte(i.Kind)
		return 1, nil

	case ie.IEFormatFour: // Type-4 TLV: tag + length + payload
		b[off] = byte(i.Kind)
		if i.Length > uint16(len(b)-off-1) {
			return 0, io.ErrShortBuffer
		}
		b[off+1] = byte(len(i.Payload))
		copy(b[off+2:], i.Payload)
		return 2 + len(i.Payload), nil

	default:
		return 0, fmt.Errorf("%w: unsupported format %d for pass-through tagged IE", ie.ErrInvalidFormat, i.Format)
	}
}
