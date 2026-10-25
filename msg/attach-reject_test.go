package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestAttachRejectRoundTrip(t *testing.T) {
	t.Parallel()

	m := msg.NewAttachReject(
		ie.NewPositionalEMMCause(0x12),
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
		SequenceNumber:        5,
		MsgType:               msg.EMMAttachReject,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.AttachReject
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

	if got := decoded.MessageType(); got != msg.EMMAttachReject {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMAttachReject)
	}
	if got := decoded.MessageTypeName(); got != "Attach reject" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Attach reject")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if decoded.EmmHeader.SecurityHeaderType != msg.SecurityHeaderTypeIntegrityProtected {
		t.Errorf("SHT = %#02x, want %#02x", decoded.EmmHeader.SecurityHeaderType, msg.SecurityHeaderTypeIntegrityProtected)
	}
	if decoded.EMMCause == nil {
		t.Error("EMMCause is nil")
	}
}

func TestAttachRejectWithOptionalTags(t *testing.T) {
	t.Parallel()

	esmc := ie.NewESMMessageContainer([]byte{0xD1, 0x01})
	t3402 := ie.NewT3402ValueForAttachReject(0x05)
	extCause := ie.NewExtendedEMMCause(0x01)

	m := msg.NewAttachReject(
		ie.NewPositionalEMMCause(0x12),
		esmc, t3402, extCause,
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0xAA, 0xBB, 0xCC, 0xDD},
		SequenceNumber:        3,
		MsgType:               msg.EMMAttachReject,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.AttachReject
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	rebuf := make([]byte, decoded.MarshalLen())
	if err := decoded.MarshalTo(rebuf); err != nil {
		t.Fatalf("re-MarshalTo failed: %v", err)
	}
	for i := range buf {
		if buf[i] != rebuf[i] {
			t.Fatalf("byte[%d] = %#02x, want %#02x", i, rebuf[i], buf[i])
		}
	}

	if decoded.EMMCause == nil {
		t.Error("EMMCause is nil")
	}
	if decoded.ESMMessageContainer == nil {
		t.Error("ESMMessageContainer is nil")
	}
	if decoded.T3402Value == nil {
		t.Error("T3402Value is nil")
	}
	if decoded.ExtendedEMMCause == nil {
		t.Error("ExtendedEMMCause is nil")
	}
}

func TestAttachRejectPlainHeader(t *testing.T) {
	t.Parallel()

	m := msg.NewAttachReject(
		ie.NewPositionalEMMCause(0x12),
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMAttachReject,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.AttachReject
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.EmmHeader.SecurityHeaderType != msg.SecurityHeaderTypePlain {
		t.Errorf("SHT = %#02x, want %#02x", decoded.EmmHeader.SecurityHeaderType, msg.SecurityHeaderTypePlain)
	}

	rebuf := make([]byte, decoded.MarshalLen())
	if err := decoded.MarshalTo(rebuf); err != nil {
		t.Fatalf("re-MarshalTo failed: %v", err)
	}
	for i := range buf {
		if buf[i] != rebuf[i] {
			t.Fatalf("byte[%d] = %#02x, want %#02x", i, rebuf[i], buf[i])
		}
	}
}

func TestAttachRejectAllOptionalKindCases(t *testing.T) {
	t.Parallel()

	// Exercise every switch case in NewAttachReject.
	esmc := ie.NewESMMessageContainer([]byte{0xD1, 0x01})
	t3346 := ie.NewT3346Value(0x05)
	t3402 := ie.NewT3402ValueForAttachReject(0x05)
	extCause := ie.NewExtendedEMMCause(0x01)
	lbt := ie.NewLowerBoundTimerValue(0x03)

	m := msg.NewAttachReject(
		ie.NewPositionalEMMCause(0x12),
		esmc, t3346, t3402, extCause, lbt,
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMAttachReject,
	}

	if m.EMMCause == nil {
		t.Fatal("EMMCause is nil")
	}
	if m.ESMMessageContainer == nil {
		t.Fatal("ESMMessageContainer is nil")
	}
	if m.T3346Value == nil {
		t.Fatal("T3346Value is nil")
	}
	if m.T3402Value == nil {
		t.Fatal("T3402Value is nil")
	}
	if m.ExtendedEMMCause == nil {
		t.Fatal("ExtendedEMMCause is nil")
	}
	if m.LowerBoundTimer == nil {
		t.Fatal("LowerBoundTimer is nil")
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.AttachReject
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.ESMMessageContainer == nil {
		t.Error("decoded ESMMessageContainer is nil")
	}
	if decoded.T3346Value == nil {
		t.Error("decoded T3346Value is nil")
	}
	if decoded.T3402Value == nil {
		t.Error("decoded T3402Value is nil")
	}
	if decoded.ExtendedEMMCause == nil {
		t.Error("decoded ExtendedEMMCause is nil")
	}
	if decoded.LowerBoundTimer == nil {
		t.Error("decoded LowerBoundTimer is nil")
	}
}

func TestAttachRejectDuplicateAndDefaultCase(t *testing.T) {
	t.Parallel()

	// First ESMMessageContainer goes to primary, second is a duplicate → m.IEs.
	first := ie.NewESMMessageContainer([]byte{0xD1, 0x01})
	dup := ie.NewESMMessageContainer([]byte{0xD1, 0x02})

	m := msg.NewAttachReject(
		ie.NewPositionalEMMCause(0x12),
		first, dup,
	)

	if m.EMMCause == nil {
		t.Fatal("EMMCause is nil")
	}
	if m.ESMMessageContainer == nil {
		t.Fatal("ESMMessageContainer is nil")
	}
	if len(m.IEs) != 1 {
		t.Fatalf("m.IEs length = %d, want 1 (dup)", len(m.IEs))
	}

	if m.IEs[0].Kind != ie.ESMMessageContainer {
		t.Errorf("m.IEs[0] Kind = 0x%02x, want 0x%02x", m.IEs[0].Kind, ie.ESMMessageContainer)
	}

	// Test nil extra IE (loop should skip it).
	m2 := msg.NewAttachReject(
		ie.NewPositionalEMMCause(0x12),
		nil, // nil IE in variadic — should be skipped by loop
	)
	if m2.EMMCause == nil {
		t.Fatal("m2 EMMCause is nil")
	}
	if len(m2.IEs) != 0 {
		t.Fatalf("m2 IEs length = %d, want 0", len(m2.IEs))
	}

	// Test unknown IE Kind → default case, appended to m.IEs.
	unknownIE := ie.NewRaw(0x7e, ie.IEFormatFour, []byte{0x99})
	m3 := msg.NewAttachReject(
		ie.NewPositionalEMMCause(0x12),
		unknownIE,
	)
	if len(m3.IEs) != 1 {
		t.Fatalf("m3 IEs length = %d, want 1", len(m3.IEs))
	}

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMAttachReject,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.AttachReject
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if len(decoded.IEs) != 1 {
		t.Fatalf("decoded IEs length = %d, want 1", len(decoded.IEs))
	}
}

func TestAttachRejectErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.AttachReject

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

	local := &msg.AttachReject{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.AttachReject
	err = decoded.UnmarshalBinary([]byte{0x17, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0x44})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated protected = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0x07, 0x44}) // plain header - no content
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("plain header only = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0xAA, 0x44}) // invalid SHT
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid SHT = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}
}
