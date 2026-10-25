package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestSecurityModeCommandRoundTrip(t *testing.T) {
	t.Parallel()

	// Security Mode Command with full IE set, SHT 0x03.
	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtectedWithNewEPSContext,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0xAA, 0xBB, 0xCC, 0xDD},
		SequenceNumber:        2,
		MsgType:               msg.EMMSecurityModeCommand,
	}

	nasAlgos := ie.NewPositionalNASSecurityAlgorithms(0x24) // EEA1 + EIA1
	spare := ie.NewPositionalSpareHalfOctet(0x00)
	ueCaps := ie.NewPositionalUESecurityCapability([]byte{
		0x02, 0xAA, 0xBB, 0xCC, 0xDD,
	}) // 5 bytes within [2..5]

	imeisvReq := ie.NewIMEISVRequest(0x01)
	nonceUE := ie.NewNonceUE([]byte{0x11, 0x22, 0x33, 0x44})
	nonceMME := ie.NewNonceMME([]byte{0x55, 0x66, 0x77, 0x88})
	hashMME := ie.NewHashMME([]byte{0xAB, 0xCD, 0xEF, 0x01, 0x23, 0x45, 0x67, 0x89}) // 8 bytes
	addSecCap := ie.NewUEAdditionalSecurityCapability([]byte{0x01, 0x02, 0x03, 0x04})
	radioCapReq := ie.NewUERadioCapabilityIDRequest(0x01)

	m := msg.NewSecurityModeCommand(nasAlgos, spare, ueCaps,
		imeisvReq, nonceUE, nonceMME, hashMME, addSecCap, radioCapReq,
	)

	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.SecurityModeCommand
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

	if got := decoded.MessageType(); got != msg.EMMSecurityModeCommand {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMSecurityModeCommand)
	}
	if got := decoded.MessageTypeName(); got != "Security mode command" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Security mode command")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}

	if decoded.SelectedNASSecurityAlgorithms == nil {
		t.Fatal("SelectedNASSecurityAlgorithms is nil")
	}
	if decoded.SpareHalfOctet == nil {
		t.Fatal("SpareHalfOctet is nil")
	}
	if decoded.ReplayedUESecurityCapability == nil {
		t.Fatal("ReplayedUESecurityCapability is nil")
	}
	if decoded.IMEISVRequest == nil {
		t.Error("IMEISVRequest is nil")
	}
	if decoded.NonceUE == nil {
		t.Error("NonceUE is nil")
	}
	if decoded.NonceMME == nil {
		t.Error("NonceMME is nil")
	}
	if decoded.HashMME == nil {
		t.Error("HashMME is nil")
	}
	if decoded.ReplayedUEAdditionalSecurityCapability == nil {
		t.Error("ReplayedUEAdditionalSecurityCapability is nil")
	}
	if decoded.UERadioCapabilityIDRequest == nil {
		t.Error("UERadioCapabilityIDRequest is nil")
	}
}

func TestSecurityModeCommandSHT04(t *testing.T) {
	t.Parallel()

	// Security Mode Command with SHT 0x04 (ciphered with new context).
	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtectedAndCipheredWithNewEPSContext,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
		SequenceNumber:        3,
		MsgType:               msg.EMMSecurityModeCommand,
	}

	nasAlgos := ie.NewPositionalNASSecurityAlgorithms(0x36)
	spare := ie.NewPositionalSpareHalfOctet(0x00)
	ueCaps := ie.NewPositionalUESecurityCapability([]byte{0x01, 0x02}) // 2 bytes min

	m := msg.NewSecurityModeCommand(nasAlgos, spare, ueCaps)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.SecurityModeCommand
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.EmmHeader.SecurityHeaderType != msg.SecurityHeaderTypeIntegrityProtectedAndCipheredWithNewEPSContext {
		t.Errorf("SHT = %#02x, want %#02x", decoded.EmmHeader.SecurityHeaderType,
			msg.SecurityHeaderTypeIntegrityProtectedAndCipheredWithNewEPSContext)
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

func TestSecurityModeCommandEmptyOptional(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtectedWithNewEPSContext,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
		SequenceNumber:        1,
		MsgType:               msg.EMMSecurityModeCommand,
	}

	nasAlgos := ie.NewPositionalNASSecurityAlgorithms(0x24)
	spare := ie.NewPositionalSpareHalfOctet(0x00)
	ueCaps := ie.NewPositionalUESecurityCapability([]byte{0x01, 0x02})

	m := msg.NewSecurityModeCommand(nasAlgos, spare, ueCaps)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.SecurityModeCommand
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

	if decoded.IMEISVRequest != nil {
		t.Error("unexpected IMEISVRequest")
	}
	if decoded.NonceUE != nil {
		t.Error("unexpected NonceUE")
	}
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs")
	}
}

func TestSecurityModeCommandErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.SecurityModeCommand

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

	local := &msg.SecurityModeCommand{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.SecurityModeCommand

	// Valid protected header (7 bytes), SHT=0x03, PD=0x7 → 0x37
	// Truncated protected header (only 5 of 7 header bytes)
	err = decoded.UnmarshalBinary([]byte{0x37, 0xAA, 0xBB, 0xCC, 0xDD})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated protected header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Valid header but no content — need at least NAS security algorithms (1 byte)
	err = decoded.UnmarshalBinary([]byte{0x37, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0x5D})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("no content after header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Valid header + 1 byte (NAS security algorithms OK), then missing spare half octet
	err = decoded.UnmarshalBinary([]byte{0x37, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0x5D, 0x24})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("missing spare half octet = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Valid header + NAS algo + spare half octet, then missing replayed UE sec cap
	err = decoded.UnmarshalBinary([]byte{0x37, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0x5D, 0x24, 0x00})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("missing replayed UE sec cap = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0xAA, 0xAA}) // invalid SHT (SHT=1, PD=0xE not EMM)
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid SHT = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}

	// Plain header for SMC (2 bytes, no payload → unexpected EOF)
	err = decoded.UnmarshalBinary([]byte{0x07, 0x5D})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("plain header for SMC = %v, want %T", err, io.ErrUnexpectedEOF)
	}
}

func TestSecurityModeCommandWithDuplicates(t *testing.T) {
	t.Parallel()

	nasAlgos := ie.NewPositionalNASSecurityAlgorithms(0x24)
	spare := ie.NewPositionalSpareHalfOctet(0x00)
	ueCaps := ie.NewPositionalUESecurityCapability([]byte{0x01, 0x02})

	// Both HashMME have exactly 8 bytes as required.
	hash1 := ie.NewHashMME([]byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88})
	dupHash := ie.NewHashMME([]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22})

	m := msg.NewSecurityModeCommand(nasAlgos, spare, ueCaps, hash1, dupHash)

	if m.HashMME == nil {
		t.Fatal("first HashMME was not stored")
	}

	var foundDup bool
	for _, i := range m.IEs {
		if i.Kind == ie.HashMME && len(i.Payload) == 8 {
			foundDup = true
		}
	}

	if !foundDup {
		t.Error("expected duplicate HashMME in IEs")
	}
}

func TestSecurityModeCommandWithAllTags(t *testing.T) {
	t.Parallel()

	nasAlgos := ie.NewPositionalNASSecurityAlgorithms(0x24)
	spare := ie.NewPositionalSpareHalfOctet(0x00)
	ueCaps := ie.NewPositionalUESecurityCapability([]byte{0x01, 0x02})

	// All tagged optional IEs with correct lengths.
	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtectedWithNewEPSContext,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0xFF, 0xEE, 0xDD, 0xCC},
		SequenceNumber:        5,
		MsgType:               msg.EMMSecurityModeCommand,
	}

	m := msg.NewSecurityModeCommand(nasAlgos, spare, ueCaps,
		ie.NewIMEISVRequest(0x01),
		ie.NewNonceUE([]byte{0x01, 0x02, 0x03, 0x04}),
		ie.NewNonceMME([]byte{0x05, 0x06, 0x07, 0x08}),
		ie.NewHashMME([]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22}), // 8 bytes
		ie.NewUEAdditionalSecurityCapability([]byte{0x01, 0x02, 0x03, 0x04}),
		ie.NewUERadioCapabilityIDRequest(0x01),
		ie.NewUEInformationRequest(0x02),
	)

	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.SecurityModeCommand
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.IMEISVRequest == nil {
		t.Error("missing IMEISVRequest after decode")
	}
	if decoded.NonceUE == nil {
		t.Error("missing NonceUE after decode")
	}
	if decoded.NonceMME == nil {
		t.Error("missing NonceMME after decode")
	}
	if decoded.HashMME == nil {
		t.Error("missing HashMME after decode")
	}
	if decoded.ReplayedUEAdditionalSecurityCapability == nil {
		t.Error("missing ReplayedUEAdditionalSecurityCapability after decode")
	}
	if decoded.UERadioCapabilityIDRequest == nil {
		t.Error("missing UERadioCapabilityIDRequest after decode")
	}
	if decoded.UEInformationRequest == nil {
		t.Error("missing UEInformationRequest after decode")
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
