package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestActivateDefaultEPSBearerContextRejectRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x01,
		MsgType:                      msg.ESMActivateDefaultEPSBearerContextReject,
	}

	esmCause := ie.NewPositionalESMCause(ie.RequestRejectedUnspecifiedSmCause)

	m := msg.NewActivateDefaultEPSBearerContextReject(esmCause)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ActivateDefaultEPSBearerContextReject
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

	if got := decoded.MessageType(); got != msg.ESMActivateDefaultEPSBearerContextReject {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.ESMActivateDefaultEPSBearerContextReject)
	}
	if got := decoded.MessageTypeName(); got != "Activate default EPS bearer context reject" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Activate default EPS bearer context reject")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}

	if decoded.ESMCause == nil {
		t.Fatal("ESMCause is nil")
	}

	if decoded.EsmHeader.BearerIdentity != 0x01 {
		t.Errorf("BearerIdentity = %#02x, want 0x01", decoded.EsmHeader.BearerIdentity)
	}
	if decoded.EsmHeader.ProcedureTransactionIdentity != 0x01 {
		t.Errorf("PTI = %#02x, want 0x01", decoded.EsmHeader.ProcedureTransactionIdentity)
	}
}

func TestActivateDefaultEPSBearerContextRejectCauseValue(t *testing.T) {
	t.Parallel()

	cause := ie.NewPositionalESMCause(ie.NetworkFailureSmCause)

	m := msg.NewActivateDefaultEPSBearerContextReject(cause)
	m.EsmHeader = msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x03,
		MsgType:                      msg.ESMActivateDefaultEPSBearerContextReject,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ActivateDefaultEPSBearerContextReject
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

	if decoded.ESMCause.Payload[0] != byte(ie.NetworkFailureSmCause) {
		t.Errorf("ESMCause value = %#02x, want %#02x", decoded.ESMCause.Payload[0], byte(ie.NetworkFailureSmCause))
	}
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs")
	}
}

func TestActivateDefaultEPSBearerContextRejectMultipleCausValues(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		cause  ie.SmCause
		bearID uint8
		pti    uint8
	}{
		{name: "QoS not accepted", cause: ie.EpsQosNotAcceptedSmCause, bearID: 0x01, pti: 0x01},
		{name: "Network failure", cause: ie.NetworkFailureSmCause, bearID: 0x02, pti: 0x05},
		{name: "Reactivation requested", cause: ie.ReactivationRequestedSmCause, bearID: 0x03, pti: 0xFF},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cause := ie.NewPositionalESMCause(tc.cause)

			m := msg.NewActivateDefaultEPSBearerContextReject(cause)
			m.EsmHeader = msg.EsmHeader{
				BearerIdentity:               tc.bearID,
				ProtocolDiscriminator:        byte(msg.PdEsm),
				ProcedureTransactionIdentity: tc.pti,
				MsgType:                      msg.ESMActivateDefaultEPSBearerContextReject,
			}

			buf := make([]byte, m.MarshalLen())
			if err := m.MarshalTo(buf); err != nil {
				t.Fatalf("MarshalTo failed: %v", err)
			}

			var decoded msg.ActivateDefaultEPSBearerContextReject
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

			if decoded.EsmHeader.BearerIdentity != tc.bearID {
				t.Errorf("BearerIdentity = %#02x, want %#02x", decoded.EsmHeader.BearerIdentity, tc.bearID)
			}
			if decoded.EsmHeader.ProcedureTransactionIdentity != tc.pti {
				t.Errorf("PTI = %#02x, want %#02x", decoded.EsmHeader.ProcedureTransactionIdentity, tc.pti)
			}
		})
	}
}

func TestActivateDefaultEPSBearerContextRejectErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.ActivateDefaultEPSBearerContextReject

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

	local := &msg.ActivateDefaultEPSBearerContextReject{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.ActivateDefaultEPSBearerContextReject
	// 3-byte header with invalid PD: must be 0x2E for ESM
	err = decoded.UnmarshalBinary([]byte{0xAA, 0x01, 0xC3})
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid PD = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}

	// Valid header (low nibble 0x2 matches PdEsm) but no payload → UnexpectedEOF
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, 0xC3})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header-only = %v, want %T", err, io.ErrUnexpectedEOF)
	}
}

func TestActivateDefaultEPSBearerContextRejectMarshalLenExtraIEs(t *testing.T) {
	t.Parallel()

	collisionIE := &ie.IE{Kind: 0x32, Format: ie.IEFormatFour, Payload: []byte{0xAA}}

	m := &msg.ActivateDefaultEPSBearerContextReject{
		EsmHeader: msg.EsmHeader{BearerIdentity: 0x01, ProtocolDiscriminator: byte(msg.PdEsm), ProcedureTransactionIdentity: 0x01, MsgType: msg.ESMActivateDefaultEPSBearerContextReject},
		ESMCause:  ie.NewESMCause(0x02),
		IEs:       []*ie.IE{collisionIE},
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo with collision IE failed: %v", err)
	}

	// Collision TLV short-buffer error path for the IEs loop.
	collisionMsg := &msg.ActivateDefaultEPSBearerContextReject{
		EsmHeader: msg.EsmHeader{BearerIdentity: 0x01, ProtocolDiscriminator: byte(msg.PdEsm), ProcedureTransactionIdentity: 0x01, MsgType: msg.ESMActivateDefaultEPSBearerContextReject},
		ESMCause:  ie.NewESMCause(0x02),
		IEs:       []*ie.IE{collisionIE},
	}
	shortBuf := make([]byte, 1)
	if err := collisionMsg.MarshalTo(shortBuf); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("expected ErrShortBuffer for collision IE error path, got %v", err)
	}
}
