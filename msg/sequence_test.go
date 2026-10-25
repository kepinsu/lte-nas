package msg

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

// writePassThroughTagged tests.
func TestWritePassThroughTagged(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		ie      *ie.IE
		buf     []byte
		wantOff int
		wantErr bool
		errType error
	}{
		{
			name:    "nil IE",
			ie:      nil,
			buf:     make([]byte, 10),
			wantOff: 0,
			wantErr: false,
		},
		{
			name:    "type-1 pass-through",
			ie:      ie.NewRaw(0x53, ie.IEFormatOne, []byte{0x01}),
			buf:     make([]byte, 1),
			wantOff: 1,
			wantErr: false,
		},
		{
			name:    "type-4 pass-through",
			ie:      ie.NewRaw(0x52, ie.IEFormatFour, []byte{1, 2}),
			buf:     make([]byte, 4), // tag + len + payload
			wantOff: 4,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			off, err := writePassThroughTagged(tt.buf, 0, tt.ie)
			if (err != nil) != tt.wantErr {
				t.Errorf("writePassThroughTagged error = %v, wantErr %v", err, tt.wantErr)
			}
			if !errors.Is(err, tt.errType) {
				t.Errorf("writePassThroughTagged error type = %v, want %T", err, tt.errType)
			}
			if off != tt.wantOff {
				t.Errorf("writePassThroughTagged offset = %d, want %d", off, tt.wantOff)
			}
		})
	}

	// Type-1 with nil payload should error
	t.Run("type-1 nil payload", func(t *testing.T) {
		t.Parallel()

		ie := &ie.IE{Kind: 0x53, Format: ie.IEFormatOne, Payload: nil}
		buf := make([]byte, 1)
		off, err := writePassThroughTagged(buf, 0, ie)
		if err == nil {
			t.Fatal("writePassThroughTagged type-1 nil payload expected error")
		}
		if off != 0 {
			t.Errorf("writePassThroughTagged offset = %d, want 0", off)
		}
	})

	// Short buffer type-1
	t.Run("type-1 short buffer", func(t *testing.T) {
		t.Parallel()

		ie := ie.NewRaw(0x53, ie.IEFormatOne, []byte{0x01})
		buf := make([]byte, 0)
		_, err := writePassThroughTagged(buf, 0, ie)
		if !errors.Is(err, io.ErrShortBuffer) {
			t.Errorf("writePassThroughTagged short buffer = %v, want %T", err, io.ErrShortBuffer)
		}
	})

	// Short buffer type-4 (length check in MarshalTo catches it)
	t.Run("type-4 short buffer", func(t *testing.T) {
		t.Parallel()

		ie := ie.NewRaw(0x52, ie.IEFormatFour, []byte{1, 2})
		buf := make([]byte, 2) // only tag + len byte, no room for payload
		_, err := writePassThroughTagged(buf, 0, ie)
		if !errors.Is(err, io.ErrShortBuffer) {
			t.Errorf("writePassThroughTagged short buffer = %v, want %T", err, io.ErrShortBuffer)
		}
	})

	// Unsupported format (type-2 tag-only)
	t.Run("unsupported format", func(t *testing.T) {
		t.Parallel()

		ie4 := ie.NewRaw(0xFF, ie.IEFormatTwo, []byte{})
		buf := make([]byte, 10)
		_, err := writePassThroughTagged(buf, 0, ie4)
		if err == nil {
			t.Fatal("writePassThroughTagged unsupported format expected error")
		}
		if !errors.Is(err, ie.ErrInvalidFormat) {
			t.Errorf("writePassThroughTagged error = %v, want wrapped %T", err, ie.ErrInvalidFormat)
		}
	})
}

func TestWritePassThroughTagged_PayloadIntegrity(t *testing.T) {
	t.Parallel()

	// Type-4: verify tag + length + payload are written correctly
	ie4 := ie.NewRaw(0x52, ie.IEFormatFour, []byte{0xDE, 0xAD})
	buf := make([]byte, 4)
	off, err := writePassThroughTagged(buf, 0, ie4)
	if err != nil {
		t.Fatalf("writePassThroughTagged = error: %v", err)
	}
	if off != 4 {
		t.Errorf("writePassThroughTagged offset = %d, want 4", off)
	}

	if buf[0] != 0x52 {
		t.Errorf("tag byte = %#02x, want %#02x", buf[0], 0x52)
	}
	if buf[1] != 2 {
		t.Errorf("length byte = %d, want 2", buf[1])
	}
	if len(buf) < 4 {
		t.Fatalf("buf too short: %d", len(buf))
	}
	if buf[2] != 0xDE || buf[3] != 0xAD {
		t.Errorf("payload = %v, want [0xde 0xad]", buf[2:4])
	}

	// Type-1: payload OR'd with tag nibble
	ie2 := ie.NewRaw(0x53, ie.IEFormatOne, []byte{0x07})
	buf2 := make([]byte, 1)
	off, err = writePassThroughTagged(buf2, 0, ie2)
	if err != nil {
		t.Fatalf("writePassThroughTagged type-1 = error: %v", err)
	}
	if off != 1 {
		t.Errorf("writePassThroughTagged offset = %d, want 1", off)
	}

	// payload[0] & 0x0f = 0x07, tag nibble = 0x50, so wire byte = 0x57
	wantWire := byte(0x07) | 0x50
	if buf2[0] != wantWire {
		t.Errorf("type-1 wire byte = %#02x, want %#02x", buf2[0], wantWire)
	}
}

func TestWritePassThroughTagged_CopyBehavior(t *testing.T) {
	t.Parallel()

	// Type-4: verify payload is copied (not shared reference)
	payload := []byte{0x11, 0x22}
	ie := &ie.IE{Kind: 0x52, Format: ie.IEFormatFour, Payload: payload}
	buf := make([]byte, 4)
	_, err := writePassThroughTagged(buf, 0, ie)
	if err != nil {
		t.Fatalf("writePassThroughTagged = error: %v", err)
	}

	// Mutate original payload after write — buffer should be unchanged
	payload[0] = 0xFF
	if buf[2] != 0x11 {
		t.Errorf("payload not copied: buf[2] = %#02x, want 0x11", buf[2])
	}
	if payload[0] != 0xFF {
		t.Errorf("original payload was mutated: payload[0] = %#02x", payload[0])
	}
}
