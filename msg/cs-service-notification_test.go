package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestCSServiceNotificationRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x99, 0x88, 0x77, 0x66},
		SequenceNumber:        4,
		MsgType:               msg.EMMCSServiceNotification,
	}

	m := msg.NewCSServiceNotification(
		ie.NewPositionalPagingIdentity(0x01),
		ie.NewCLI([]byte("04095551234")),
		ie.NewSSCode(0x07),
		ie.NewLCSIndicator(0x01),
	)

	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.CSServiceNotification
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

	if got := decoded.MessageType(); got != msg.EMMCSServiceNotification {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMCSServiceNotification)
	}
	if got := decoded.MessageTypeName(); got != "CS service notification" {
		t.Errorf("MessageTypeName = %q, want %q", got, "CS service notification")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if decoded.PagingIdentity == nil {
		t.Fatal("PagingIdentity is nil")
	}
	if decoded.Cli == nil {
		t.Error("Cli is nil")
	}
}

func TestCSServiceNotificationEmptyOptional(t *testing.T) {
	t.Parallel()

	m := msg.NewCSServiceNotification(ie.NewPositionalPagingIdentity(0x02))

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
		SequenceNumber:        1,
		MsgType:               msg.EMMCSServiceNotification,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.CSServiceNotification
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

	if decoded.Cli != nil {
		t.Error("unexpected Cli")
	}
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs")
	}
}

func TestCSServiceNotificationWithLcsClientIdentityAndExtras(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMCSServiceNotification,
	}

	m := msg.NewCSServiceNotification(
		ie.NewPositionalPagingIdentity(0x01),
		ie.NewCLI([]byte("04095551234")),
		ie.NewLCSClientIdentity([]byte{0x01}),
	)

	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo: %v", err)
	}

	var decoded msg.CSServiceNotification
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary: %v", err)
	}

	rebuf := make([]byte, decoded.MarshalLen())
	if err := decoded.MarshalTo(rebuf); err != nil {
		t.Fatalf("re-MarshalTo: %v", err)
	}
	for i := range buf {
		if buf[i] != rebuf[i] {
			t.Fatalf("byte[%d] = %#02x, want %#02x", i, rebuf[i], buf[i])
		}
	}

	if decoded.Cli == nil {
		t.Fatal("Cli is nil")
	}
	if decoded.LcsClientIdentity == nil {
		t.Fatal("LcsClientIdentity is nil")
	}
	if decoded.LcsIndicator != nil {
		t.Error("LcsIndicator should be nil")
	}
	if decoded.SsCode != nil {
		t.Error("SsCode should be nil")
	}
}

func TestCSServiceNotificationDuplicateAndExtraIEs(t *testing.T) {
	t.Parallel()

	// Create message with a duplicate CLI and an unrecognized IE.
	extraIE := ie.NewNetworkResourceIdentifierContainer([]byte{0x12, 0x34})
	duplicateCLI := ie.NewCLI([]byte("9999"))

	m := msg.NewCSServiceNotification(
		ie.NewPositionalPagingIdentity(0x01),
		ie.NewCLI([]byte("04095551234")),
		duplicateCLI, // duplicate CLI → should go to IEs
		extraIE,      // unrecognized IE → should go to IEs
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMCSServiceNotification,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo: %v", err)
	}

	var decoded msg.CSServiceNotification
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary: %v", err)
	}

	rebuf := make([]byte, decoded.MarshalLen())
	if err := decoded.MarshalTo(rebuf); err != nil {
		t.Fatalf("re-MarshalTo: %v", err)
	}
	for i := range buf {
		if buf[i] != rebuf[i] {
			t.Fatalf("byte[%d] = %#02x, want %#02x", i, rebuf[i], buf[i])
		}
	}

	if decoded.Cli == nil {
		t.Fatal("Cli is nil")
	}
	if len(decoded.IEs) != 2 {
		t.Fatalf("expected 2 extra IEs (duplicate CLI + extra IE), got %d", len(decoded.IEs))
	}
}

func TestCSServiceNotificationErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.CSServiceNotification

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

	local := &msg.CSServiceNotification{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.CSServiceNotification
	err = decoded.UnmarshalBinary([]byte{0x17, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0x64}) // header only
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header-only = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0xAA, 0x64}) // invalid SHT
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid SHT = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}
}
