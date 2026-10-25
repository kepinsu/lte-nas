package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestGUTIReallocationCommandRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
		SequenceNumber:        5,
		MsgType:               msg.EMMGUTIReallocationCommand,
	}

	guti := ie.NewPositionalEPSMobileIdentity([]byte{0x46, 0x05, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF})
	taiList := ie.NewTAIList([]byte{0x23, 0x1A, 0xAB, 0xCD, 0xEF, 0x12, 0x34})
	dcnID := ie.NewDCNID([]byte{0x01, 0x02})
	ucap := ie.NewUERadioCapabilityID([]byte{0xDE, 0xAD, 0xBE, 0xEF})
	ucapDel := ie.NewUERadioCapabilityIDDeletionIndication(1)

	m := msg.NewGUTIReallocationCommand(guti, taiList, dcnID, ucap, ucapDel)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.GUTIReallocationCommand
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

	if got := decoded.MessageType(); got != msg.EMMGUTIReallocationCommand {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMGUTIReallocationCommand)
	}
	if got := decoded.MessageTypeName(); got != "Guti reallocation command" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Guti reallocation command")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}

	if decoded.GUTI == nil {
		t.Fatal("GUTI is nil")
	}
	if decoded.TAIList == nil {
		t.Error("TAIList is nil")
	}
	if decoded.DCNID == nil {
		t.Error("DCNID is nil")
	}
	if decoded.UERadioCapabilityID == nil {
		t.Error("UERadioCapabilityID is nil")
	}
	if decoded.UERadioCapabilityIDDeletionInd == nil {
		t.Error("UERadioCapabilityIDDeletionInd is nil")
	}
}

func TestGUTIReallocationCommandPlainHeader(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMGUTIReallocationCommand,
	}

	guti := ie.NewPositionalEPSMobileIdentity([]byte{0x46, 0x05, 0xAA, 0xBB})
	m := msg.NewGUTIReallocationCommand(guti)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.GUTIReallocationCommand
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

	if decoded.GUTI == nil {
		t.Fatal("GUTI is nil")
	}
	if decoded.TAIList != nil {
		t.Error("unexpected TAIList")
	}
}

func TestGUTIReallocationCommandEmptyOptional(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0xAA, 0xBB, 0xCC, 0xDD},
		SequenceNumber:        1,
		MsgType:               msg.EMMGUTIReallocationCommand,
	}

	guti := ie.NewPositionalEPSMobileIdentity([]byte{0x46, 0x05, 0x12, 0x34, 0xAB, 0xCD})
	m := msg.NewGUTIReallocationCommand(guti)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.GUTIReallocationCommand
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

	if decoded.TAIList != nil {
		t.Error("unexpected TAIList")
	}
	if decoded.DCNID != nil {
		t.Error("unexpected DCNID")
	}
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs")
	}
}

func TestGUTIReallocationCommandErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.GUTIReallocationCommand
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

	local := &msg.GUTIReallocationCommand{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.GUTIReallocationCommand
	err = decoded.UnmarshalBinary([]byte{0x17, 0x50}) // header only - no payload
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header-only = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0x17, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0x50})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated content = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0xAA, 0x50}) // invalid SHT
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid SHT = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}
}
