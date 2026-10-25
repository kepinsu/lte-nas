package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestGUTIReallocationCompleteRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMGUTIReallocationComplete,
	}

	m := msg.NewGUTIReallocationComplete()
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.GUTIReallocationComplete
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	rebuf := make([]byte, decoded.MarshalLen())
	if err := decoded.MarshalTo(rebuf); err != nil {
		t.Fatalf("re-MarshalTo failed: %v", err)
	}
	if len(buf) != len(rebuf) {
		t.Fatalf("re-marshal length = %d, want %d", len(rebuf), len(buf))
	}
	for i := range buf {
		if buf[i] != rebuf[i] {
			t.Fatalf("byte[%d] = %#02x, want %#02x", i, rebuf[i], buf[i])
		}
	}

	if got := decoded.MessageType(); got != msg.EMMGUTIReallocationComplete {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMGUTIReallocationComplete)
	}
	if got := decoded.MessageTypeName(); got != "Guti reallocation complete" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Guti reallocation complete")
	}
	if !decoded.IsRequest() {
		t.Error("IsRequest = false, want true")
	}
}

func TestGUTIReallocationCompleteEmpty(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMGUTIReallocationComplete,
	}

	m := msg.NewGUTIReallocationComplete()
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	if len(buf) != 2 {
		t.Fatalf("buffer length = %d, want 2 (header only)", len(buf))
	}

	var decoded msg.GUTIReallocationComplete
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs")
	}
}

func TestGUTIReallocationCompleteWithIEs(t *testing.T) {
	t.Parallel()

	// Use a catalogue-type IE (NetworkResourceIdentifierContainer, Kind 0x10)
	// that parses correctly through ie.ParseMultiIEs.
	container := ie.NewNetworkResourceIdentifierContainer([]byte{0x12, 0x34})

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMGUTIReallocationComplete,
	}

	m := msg.NewGUTIReallocationComplete(container)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo: %v", err)
	}

	var decoded msg.GUTIReallocationComplete
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary: %v", err)
	}

	if len(decoded.IEs) != 1 {
		t.Fatalf("expected 1 IE, got %d", len(decoded.IEs))
	}

	if decoded.IsRequest() != true {
		t.Error("IsRequest = false, want true")
	}
}

func TestGUTIReallocationCompleteErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.GUTIReallocationComplete
	if got := nilMsg.MarshalLen(); got != 0 {
		t.Errorf("nil MarshalLen = %d, want 0", got)
	}

	err := nilMsg.MarshalTo(nil)
	if err == nil {
		t.Fatal("nil MarshalTo = nil, want error")
	}
	if !errors.Is(err, ie.ErrNilIE) {
		t.Errorf("nil MarshalTo = %v, want wrapped %T", err, ie.ErrNilIE)
	}

	local := &msg.GUTIReallocationComplete{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.GUTIReallocationComplete
	err = decoded.UnmarshalBinary([]byte{0x17, 0x51}) // header only - no payload
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header-only = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0xAA, 0x51}) // invalid SHT
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid SHT = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}
}
