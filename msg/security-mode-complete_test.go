package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestSecurityModeCompleteRoundTrip(t *testing.T) {
	t.Parallel()

	m := msg.NewSecurityModeComplete(
		ie.NewMobileIdentity([]byte{0x01, 0x12, 0x34, 0xAB}),  // IMEISV
		ie.NewReplayedNASMessageContainer([]byte{0xAA, 0xBB}), // replayed container
		ie.NewUERadioCapabilityID([]byte{0x12, 0x34}),         // radio cap ID
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMSecurityModeComplete,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.SecurityModeComplete
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

	if got := decoded.MessageType(); got != msg.EMMSecurityModeComplete {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMSecurityModeComplete)
	}
	if got := decoded.MessageTypeName(); got != "Security mode complete" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Security mode complete")
	}
	if !decoded.IsRequest() {
		t.Error("IsRequest = false, want true")
	}

	if decoded.IMEISV == nil {
		t.Error("IMEISV is nil")
	}
	if decoded.ReplayedNASMessageContainer == nil {
		t.Error("ReplayedNASMessageContainer is nil")
	}
	if decoded.UERadioCapabilityID == nil {
		t.Error("UERadioCapabilityID is nil")
	}
}

func TestSecurityModeCompleteEmptyOptional(t *testing.T) {
	t.Parallel()

	m := msg.NewSecurityModeComplete()

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMSecurityModeComplete,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.SecurityModeComplete
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

	if decoded.IMEISV != nil {
		t.Error("unexpected IMEISV")
	}
	if decoded.ReplayedNASMessageContainer != nil {
		t.Error("unexpected ReplayedNASMessageContainer")
	}
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs")
	}
}

func TestSecurityModeCompleteErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.SecurityModeComplete

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

	local := &msg.SecurityModeComplete{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.SecurityModeComplete

	// Valid header only (2 bytes, SHT=0x00 = plain EMM) — no content, which is fine
	err = decoded.UnmarshalBinary([]byte{0x07, 0x5E})
	if err != nil {
		t.Fatalf("valid header-only = %v, want nil", err)
	}

	err = decoded.UnmarshalBinary([]byte{0x07}) // truncated header
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0xAA, 0x5E}) // invalid SHT
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid SHT = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}
}

func TestSecurityModeCompleteWithUECoarseLocation(t *testing.T) {
	t.Parallel()

	m := msg.NewSecurityModeComplete(
		ie.NewUECoarseLocationInformation([]byte{0x12, 0x34, 0x56, 0xAA, 0xBB, 0xCC}), // 6 bytes fixed
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMSecurityModeComplete,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.SecurityModeComplete
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.UECoarseLocationInformation == nil {
		t.Error("missing UECoarseLocationInformation after decode")
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

func TestSecurityModeCompleteAllOptional(t *testing.T) {
	t.Parallel()

	m := msg.NewSecurityModeComplete(
		ie.NewMobileIdentity([]byte{0x01, 0x12, 0x34}),
		ie.NewReplayedNASMessageContainer([]byte{0xAA, 0xBB, 0xCC}),
		ie.NewUERadioCapabilityID([]byte{0x12, 0x34}),
		ie.NewUECoarseLocationInformation([]byte{0x12, 0x34, 0x56, 0xAA, 0xBB, 0xCC}), // 6 bytes
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMSecurityModeComplete,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.SecurityModeComplete
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.IMEISV == nil {
		t.Error("missing IMEISV")
	}
	if decoded.ReplayedNASMessageContainer == nil {
		t.Error("missing ReplayedNASMessageContainer")
	}
	if decoded.UERadioCapabilityID == nil {
		t.Error("missing UERadioCapabilityID")
	}
	if decoded.UECoarseLocationInformation == nil {
		t.Error("missing UECoarseLocationInformation")
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

func TestSecurityModeCompleteMarshalLenExtraIEs(t *testing.T) {
	t.Parallel()

	collisionIE := &ie.IE{Kind: 0x32, Format: ie.IEFormatFour, Payload: []byte{0xAA}}

	m := &msg.SecurityModeComplete{
		EmmHeader:   msg.EmmHeader{SecurityHeaderType: msg.SecurityHeaderTypePlain, ProtocolDiscriminator: byte(msg.PdEmm), MsgType: msg.EMMSecurityModeComplete},
		IMEISV:      ie.NewMobileIdentity([]byte{0x23, 0x54}),
		IEs:         []*ie.IE{collisionIE},
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo with collision IE failed: %v", err)
	}

	var decoded msg.SecurityModeComplete
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

	collisionMsg := &msg.SecurityModeComplete{
		EmmHeader:   msg.EmmHeader{SecurityHeaderType: msg.SecurityHeaderTypePlain, ProtocolDiscriminator: byte(msg.PdEmm), MsgType: msg.EMMSecurityModeComplete},
		IMEISV:      ie.NewMobileIdentity([]byte{0x23, 0x54}),
		IEs:         []*ie.IE{collisionIE},
	}
	shortBuf := make([]byte, 1)
	if err := collisionMsg.MarshalTo(shortBuf); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("expected ErrShortBuffer for collision IE error path, got %v", err)
	}
}

func TestSecurityModeCompleteNewDefaultBranch(t *testing.T) {
	t.Parallel()

	nonceIE := ie.NewNonceMME([]byte{0xAA})
	m := msg.NewSecurityModeComplete(nonceIE)
	if len(m.IEs) == 0 {
		t.Fatal("NonceMME should go to IEs field (default branch)")
	}
	if m.IEs[0].Kind != 0x56 {
		t.Fatalf("IEs[0].Kind = %#02x, want 0x56", m.IEs[0].Kind)
	}
}
