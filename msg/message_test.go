package msg

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestGenericImplementsMessage(t *testing.T) {
	t.Parallel()

	var _ Message = (*Generic)(nil)
}

func TestMessageInterfaceMethodSignatures(t *testing.T) {
	t.Parallel()

	g := NewGeneric([]byte{0x07, 0x41})

	// Verify all Message interface methods compile and work correctly.
	gotLen := g.MarshalLen()
	if gotLen != 2 {
		t.Errorf("MarshalLen = %d, want 2", gotLen)
	}

	mt := g.MessageType()
	if mt != 0x41 {
		t.Errorf("MessageType = %#02x, want %#02x", mt, 0x41)
	}

	name := g.MessageTypeName()
	if name != "Generic (0x41)" {
		t.Errorf("MessageTypeName = %q, want %q", name, "Generic (0x41)")
	}

	if g.IsRequest() {
		t.Error("IsRequest = true, want false")
	}

	buf := make([]byte, 4)
	err := g.MarshalTo(buf)
	if err != nil {
		t.Fatalf("MarshalTo = %v", err)
	}

	var g2 Generic
	err = g2.UnmarshalBinary([]byte{0x07, 0x41})
	if err != nil {
		t.Fatalf("UnmarshalBinary = %v", err)
	}
}

func TestGenericMarshalLenNil(t *testing.T) {
	t.Parallel()

	var g *Generic
	if got := g.MarshalLen(); got != 0 {
		t.Errorf("nil MarshalLen = %d, want 0", got)
	}
}

func TestWriteIE(t *testing.T) {
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
			name:    "short buffer type-4",
			ie:      ie.NewRaw(0x52, ie.IEFormatFour, []byte{1, 2}),
			buf:     make([]byte, 2), // needs tag+length+payload = 4 bytes
			wantOff: 0,
			wantErr: true,
			errType: io.ErrShortBuffer,
		},
		{
			name:    "success type-4",
			ie:      ie.NewRaw(0x52, ie.IEFormatFour, []byte{1, 2}),
			buf:     make([]byte, 4), // tag(1)+len(1)+payload(2)
			wantOff: 4,
			wantErr: false,
		},
		{
			name:    "success type-1",
			ie:      ie.NewRaw(0x80, ie.IEFormatOne, []byte{0x01}), // EMM Cause placeholder (type-1)
			buf:     make([]byte, 1),
			wantOff: 1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			off, err := writeIE(tt.buf, 0, tt.ie)
			if (err != nil) != tt.wantErr {
				t.Errorf("writeIE error = %v, wantErr %v", err, tt.wantErr)
			}
			if !errors.Is(err, tt.errType) {
				t.Errorf("writeIE error type = %v, want %T", err, tt.errType)
			}
			if off != tt.wantOff {
				t.Errorf("writeIE offset = %d, want %d", off, tt.wantOff)
			}
		})
	}

	// MarshalTo error case: short buffer after length check — FormatSix with huge
	// payload where the buffer has exactly MarshalLen bytes but off > 0 leaves too
	// little space. This forces the MarshalTo path to return a non-short-buffer
	// error. Actually, short buffer at offset will be caught by MarshalTo as well.
	// Let's test an actual MarshalTo error: FormatThree payload length doesn't match.
	t.Run("marshalTo error", func(t *testing.T) {
		t.Parallel()

		// IE with payload too short for FormatThree minimum length.
		// Many catalogue IEs expect min 3 bytes in payload for type-3 format.
		badIE := &ie.IE{Kind: 0x54, Format: ie.IEFormatThree, Payload: []byte{1}}
		buf := make([]byte, 3) // would be enough for MarshalLen but Validate rejects it
		off, err := writeIE(buf, 0, badIE)
		if err == nil {
			t.Fatal("writeIE expected error for mismatched payload length, got nil")
		}
		if off != 0 {
			t.Errorf("writeIE offset = %d, want 0 on error", off)
		}
	})

	// Verify short buffer at offset (off > 0) triggers MarshalTo's length check
	t.Run("short buffer at offset", func(t *testing.T) {
		t.Parallel()

		ie := ie.NewRaw(0x52, ie.IEFormatFour, []byte{1, 2})
		buf := make([]byte, 5)
		off, err := writeIE(buf, 3, ie) // only 2 bytes left at offset 3
		if err == nil {
			t.Fatal("writeIE expected error, got nil")
		}
		if !errors.Is(err, io.ErrShortBuffer) {
			t.Errorf("writeIE error = %v, want %T", err, io.ErrShortBuffer)
		}
		if off != 0 {
			t.Errorf("writeIE offset = %d, want 0", off)
		}
	})
}
