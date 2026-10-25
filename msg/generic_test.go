package msg_test

import (
	"errors"
	"io"
	"strconv"
	"testing"

	"github.com/kepinsu/lte-nas/msg"
)

func TestGenericRoundTrip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		raw    []byte
		wantMt uint8
	}{
		{
			name:   "arbitrary bytes",
			raw:    []byte{0xDE, 0xAD, 0xBE, 0xEF},
			wantMt: 0x00, // not a valid header, so best-effort returns 0
		},
		{
			name:   "embedded EMM header",
			raw:    []byte{0x07, msg.EMMAttachRequest},
			wantMt: msg.EMMAttachRequest,
		},
		{
			name:   "embedded ESM header",
			raw:    []byte{0x02, 0x01, msg.ESMDataTransport},
			wantMt: msg.ESMDataTransport,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := msg.NewGeneric(tt.raw)

			if g == nil {
				t.Fatal("NewGeneric returned nil")
			}

			if len(g.Raw) != len(tt.raw) {
				t.Fatalf("Raw length = %d, want %d", len(g.Raw), len(tt.raw))
			}

			for i := range tt.raw {
				if g.Raw[i] != tt.raw[i] {
					t.Errorf("Raw[%d] = %#02x, want %#02x", i, g.Raw[i], tt.raw[i])
				}
			}

			gotMt := g.MessageType()
			if gotMt != tt.wantMt {
				t.Errorf("MessageType() = %#02x, want %#02x", gotMt, tt.wantMt)
			}

			name := g.MessageTypeName()
			expectedName := "Generic (0x00)" // fallback when MessageType returns 0
			if tt.wantMt != 0 {
				expectedName = "Generic (0x" + hexByte(tt.wantMt) + ")"
			}
			if name != expectedName {
				t.Errorf("MessageTypeName() = %q, want %q", name, expectedName)
			}

			if g.IsRequest() {
				t.Error("IsRequest() = true, want false")
			}

			marshaled := make([]byte, 16)
			err := g.MarshalTo(marshaled)
			if err != nil {
				t.Fatalf("MarshalTo failed: %v", err)
			}

			for i := range tt.raw {
				if marshaled[i] != tt.raw[i] {
					t.Errorf("marshaled[%d] = %#02x, want %#02x", i, marshaled[i], tt.raw[i])
				}
			}

			var g2 msg.Generic
			err = g2.UnmarshalBinary(tt.raw)
			if err != nil {
				t.Fatalf("UnmarshalBinary failed: %v", err)
			}
			for i := range tt.raw {
				if g2.Raw[i] != tt.raw[i] {
					t.Errorf("unmarshaled[%d] = %#02x, want %#02x", i, g2.Raw[i], tt.raw[i])
				}
			}
		})
	}
}

// hexByte formats a uint8 as a two-digit hex string (no "0x" prefix).
func hexByte(v uint8) string {
	if v < 0x10 {
		return "0" + strconv.FormatUint(uint64(v), 16)
	}
	return strconv.FormatUint(uint64(v), 16)
}

func TestGenericShortBuffer(t *testing.T) {
	t.Parallel()

	g := msg.NewGeneric([]byte{0x01, 0x02})
	buf := make([]byte, 1)
	err := g.MarshalTo(buf)
	if !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("MarshalTo on short buffer = %T, want error type", err)
	}
}

func TestGenericNilReceiver(t *testing.T) {
	t.Parallel()

	var g *msg.Generic

	if got := g.MarshalLen(); got != 0 {
		t.Errorf("nil MarshalLen = %d, want 0", got)
	}

	err := g.MarshalTo(nil)
	if err == nil {
		t.Fatal("nil MarshalTo returned nil error")
	}

	if g.MessageType() != 0 {
		t.Errorf("nil MessageType = %#02x, want 0x00", g.MessageType())
	}

	if name := g.MessageTypeName(); name != "Generic (0x00)" {
		t.Errorf("nil MessageTypeName = %q, want %q", name, "Generic (0x00)")
	}

	if g.IsRequest() {
		t.Error("nil IsRequest = true, want false")
	}
}

func TestGenericEmptyRaw(t *testing.T) {
	t.Parallel()

	g := msg.NewGeneric(nil)
	buf := make([]byte, 4)
	err := g.MarshalTo(buf)
	if err == nil {
		t.Error("MarshalTo on empty Raw = nil, want error")
	}

	g2 := msg.NewGeneric([]byte{})
	if got := g2.MarshalLen(); got != 0 {
		t.Errorf("empty MarshalLen = %d, want 0", got)
	}

	var g3 msg.Generic
	err = g3.UnmarshalBinary(nil)
	if err != nil {
		t.Fatalf("UnmarshalBinary(nil) = %v, want nil", err)
	}
}
