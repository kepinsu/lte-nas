package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestServiceRejectRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
		SequenceNumber:        3,
		MsgType:               msg.EMMServiceReject,
	}

	m := msg.NewServiceReject(
		ie.NewEMMCause(ie.CsDomainNotAvailableMmCause), // CS domain not available
		ie.NewT3402Value(0x55),
		ie.NewBackOffTimerValue(0x78),
		ie.NewLowerBoundTimerValue(0x9A),
	)

	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ServiceReject
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

	if got := decoded.MessageType(); got != msg.EMMServiceReject {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMServiceReject)
	}
	if got := decoded.MessageTypeName(); got != "Service reject" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Service reject")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if decoded.EMMCause == nil {
		t.Fatal("EMMCause is nil")
	}
}

func TestServiceRejectEmptyOptional(t *testing.T) {
	t.Parallel()

	m := msg.NewServiceReject(ie.NewEMMCause(ie.CongestionMmCause)) // Congestion

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
		SequenceNumber:        3,
		MsgType:               msg.EMMServiceReject,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ServiceReject
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

	if decoded.T3402Value != nil {
		t.Error("unexpected T3402Value")
	}
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs")
	}
}

func TestServiceRejectProtectedHeader(t *testing.T) {
	t.Parallel()

	m := msg.NewServiceReject(ie.NewEMMCause(ie.UeSecurityCapabilitiesMismatchMmCause)) // Congestion

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtectedAndCiphered,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0xBB, 0xCC, 0xDD, 0xEE},
		SequenceNumber:        8,
		MsgType:               msg.EMMServiceReject,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ServiceReject
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.EmmHeader.SecurityHeaderType != msg.SecurityHeaderTypeIntegrityProtectedAndCiphered {
		t.Errorf("SHT = %#02x, want %#02x", decoded.EmmHeader.SecurityHeaderType, msg.SecurityHeaderTypeIntegrityProtectedAndCiphered)
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

func TestServiceRejectErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.ServiceReject

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

	local := &msg.ServiceReject{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.ServiceReject
	err = decoded.UnmarshalBinary([]byte{0x17, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0x4E}) // header only
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header-only = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0x17, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0x4E}) // header only - no content
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated protected = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0xAA, 0x4E}) // invalid SHT
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid SHT = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}
}

func TestServiceRejectAllOptionalIEs(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
		SequenceNumber:        5,
		MsgType:               msg.EMMServiceReject,
	}

	m := msg.NewServiceReject(
		ie.NewEMMCause(0x17), // UE security capabilities mismatch
		ie.NewT3402Value(0x55),
		ie.NewBackOffTimerValue(0x78),
		ie.NewLowerBoundTimerValue(0x9A),
		ie.NewForbiddenTAIsForRoaming([]byte{0x23, 0x1A, 0xAB, 0xCD, 0xEF, 0x12, 0x34}),
		ie.NewForbiddenTAIsForRegionalProvisionOfService([]byte{0x24, 0x1B, 0xAC, 0xDE, 0x12, 0x34}),
	)

	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ServiceReject
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.T3402Value == nil {
		t.Error("T3402Value is nil")
	}
	if decoded.BackOffTimerValue == nil {
		t.Error("BackOffTimerValue is nil")
	}
	if decoded.LowerBoundTimerValue == nil {
		t.Error("LowerBoundTimerValue is nil")
	}
	if decoded.ForbiddenTAIsForRoaming == nil {
		t.Error("ForbiddenTAIsForRoaming is nil")
	}
	if decoded.ForbiddenTAIsForRegionalProvisionOfService == nil {
		t.Error("ForbiddenTAIsForRegionalProvisionOfService is nil")
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

func TestServiceRejectWithExtraIEs(t *testing.T) {
	t.Parallel()

	// Send an unrecognized IE tag (0x9F is not defined as type-4; use valid unknown-like 0x6f).
	m := msg.NewServiceReject(
		ie.NewEMMCause(0x02),
		ie.NewT3402Value(0x55),
		ie.NewRaw(0x6f, ie.IEFormatFour, []byte{0xDE, 0xAD}),
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
		SequenceNumber:        2,
		MsgType:               msg.EMMServiceReject,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ServiceReject
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.ForbiddenTAIsForRoaming != nil {
		t.Error("unexpected ForbiddenTAIsForRoaming")
	}

	// Extra IEs should be in decoded.IEs.
	if len(decoded.IEs) != 1 {
		t.Fatalf("IEs length = %d, want 1", len(decoded.IEs))
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
