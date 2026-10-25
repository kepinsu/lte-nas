package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestModifyEPSBearerContextRejectRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x05,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x01,
		MsgType:                      msg.ESMModifyEPSBearerContextReject,
	}

	esmCause := ie.NewPositionalESMCause(0x36) // IMSI unknown in HSS

	m := msg.NewModifyEPSBearerContextReject(esmCause)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	if len(buf) != m.MarshalLen() {
		t.Fatalf("marshal length = %d, want %d", len(buf), m.MarshalLen())
	}

	var decoded msg.ModifyEPSBearerContextReject
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

	if got := decoded.MessageType(); got != msg.ESMModifyEPSBearerContextReject {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.ESMModifyEPSBearerContextReject)
	}
	if got := decoded.MessageTypeName(); got != "Modify EPS bearer context reject" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Modify EPS bearer context reject")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if decoded.ESMCause == nil {
		t.Fatal("ESMCause is nil")
	}

	if decoded.EsmHeader.BearerIdentity != 0x05 {
		t.Errorf("BearerIdentity = %#02x, want 0x05", decoded.EsmHeader.BearerIdentity)
	}
}

func TestModifyEPSBearerContextRejectMultiplePTI(t *testing.T) {
	t.Parallel()

	for _, bearerID := range []uint8{0x01, 0x02, 0x05, 0x07} {
		for _, pti := range []uint8{0x01, 0x02, 0x0F} {
			t.Run("", func(t *testing.T) {
				t.Parallel()

				h := msg.EsmHeader{
					BearerIdentity:               bearerID,
					ProtocolDiscriminator:        byte(msg.PdEsm),
					ProcedureTransactionIdentity: pti,
					MsgType:                      msg.ESMModifyEPSBearerContextReject,
				}

				esmCause := ie.NewPositionalESMCause(0x36)

				m := msg.NewModifyEPSBearerContextReject(esmCause)
				m.EsmHeader = h

				buf := make([]byte, m.MarshalLen())
				if err := m.MarshalTo(buf); err != nil {
					t.Fatalf("MarshalTo failed: %v", err)
				}

				var decoded msg.ModifyEPSBearerContextReject
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

				if decoded.EsmHeader.BearerIdentity != bearerID {
					t.Errorf("BearerIdentity = %#02x, want %#02x", decoded.EsmHeader.BearerIdentity, bearerID)
				}
				if decoded.EsmHeader.ProcedureTransactionIdentity != pti {
					t.Errorf("PTI = %#02x, want %#02x", decoded.EsmHeader.ProcedureTransactionIdentity, pti)
				}
			})
		}
	}
}

func TestModifyEPSBearerContextRejectEmptyOptional(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x01,
		MsgType:                      msg.ESMModifyEPSBearerContextReject,
	}

	esmCause := ie.NewPositionalESMCause(0x36)

	m := msg.NewModifyEPSBearerContextReject(esmCause)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ModifyEPSBearerContextReject
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

	if decoded.ExtendedQoS != nil {
		t.Error("unexpected ExtendedQoS in empty optional")
	}
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs in empty optional")
	}
}

func TestModifyEPSBearerContextRejectErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.ModifyEPSBearerContextReject

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

	local := &msg.ModifyEPSBearerContextReject{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.ModifyEPSBearerContextReject

	// Truncated ESM header (2 bytes, need 3).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Invalid PD in 3-octet header.
	err = decoded.UnmarshalBinary([]byte{0xAA, 0x01, 0xCB})
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid PD = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}

	// Valid header but no payload — need at least ESM cause (1 byte).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, 0xCB})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated no cause = %v, want %T", err, io.ErrUnexpectedEOF)
	}
}

func TestModifyEPSBearerContextRejectParseHelper(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x03,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x04,
		MsgType:                      msg.ESMModifyEPSBearerContextReject,
	}

	esmCause := ie.NewPositionalESMCause(0x36)

	m := msg.NewModifyEPSBearerContextReject(esmCause)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	parsed, err := msg.ParseModifyEPSBearerContextReject(buf)
	if err != nil {
		t.Fatalf("ParseModifyEPSBearerContextReject failed: %v", err)
	}

	if parsed == nil {
		t.Fatal("ParseModifyEPSBearerContextReject returned nil")
	}

	if parsed.ESMCause == nil {
		t.Fatal("ESMCause is nil")
	}

	if parsed.EsmHeader.BearerIdentity != 0x03 {
		t.Errorf("BearerIdentity = %#02x, want 0x03", parsed.EsmHeader.BearerIdentity)
	}
	if parsed.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if parsed.MessageType() != msg.ESMModifyEPSBearerContextReject {
		t.Errorf("MessageType = %#02x, want %#02x", parsed.MessageType(), msg.ESMModifyEPSBearerContextReject)
	}
}

func TestModifyEPSBearerContextRejectWithExtendedQoS(t *testing.T) {
	t.Parallel()

	extQOSPayload := make([]byte, 10)
	extQOSPayload[4] = 9

	extQoS := ie.NewExtendedQualityOfService(extQOSPayload)
	esmCause := ie.NewPositionalESMCause(0x36)

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x01,
		MsgType:                      msg.ESMModifyEPSBearerContextReject,
	}

	m := msg.NewModifyEPSBearerContextReject(esmCause, extQoS)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ModifyEPSBearerContextReject
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

	if decoded.ExtendedQoS == nil {
		t.Fatal("ExtendedQoS is nil after round-trip")
	}
	if len(decoded.ExtendedQoS.Payload) != 10 {
		t.Errorf("ExtendedQoS payload len = %d, want 10", len(decoded.ExtendedQoS.Payload))
	}
}

func TestModifyEPSBearerContextRejectDuplicateIEs(t *testing.T) {
	t.Parallel()

	extQoS1 := ie.NewExtendedQualityOfService([]byte{0, 0, 0, 0, 9, 0, 0, 0, 0, 0})
	extQoS2 := ie.NewExtendedQualityOfService([]byte{0, 0, 0, 0, 5, 0, 0, 0, 0, 0})
	esmCause := ie.NewPositionalESMCause(0x36)

	m := msg.NewModifyEPSBearerContextReject(esmCause, extQoS1, extQoS2)

	if m.ExtendedQoS == nil {
		t.Fatal("first ExtendedQoS was not stored")
	}

	var foundDup bool
	for _, i := range m.IEs {
		if i.Kind == ie.ExtendedQualityOfService && len(i.Payload) >= 5 && i.Payload[4] == 5 {
			foundDup = true
		}
	}

	if !foundDup {
		t.Error("expected duplicate ExtendedQoS in IEs")
	}
}
