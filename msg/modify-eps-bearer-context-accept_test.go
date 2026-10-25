package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestModifyEPSBearerContextAcceptRoundTrip(t *testing.T) {
	t.Parallel()

	qosFields := &ie.EPSQualityOfServiceFields{QCI: 9}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x01,
		MsgType:                      msg.ESMModifyEPSBearerContextAccept,
	}

	esmQoS := ie.NewPositionalEPSQualityOfService(qosPayload)
	extQoS := ie.NewExtendedQualityOfService([]byte{0, 0, 0, 0, 9, 0, 0, 0, 0, 0})
	pco := ie.NewProtocolConfigurationOptions([]byte{0x50, 0x01, 0x21})

	m := msg.NewModifyEPSBearerContextAccept(
		esmQoS, extQoS, pco,
	)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	if len(buf) != m.MarshalLen() {
		t.Fatalf("marshal length = %d, want %d", len(buf), m.MarshalLen())
	}

	var decoded msg.ModifyEPSBearerContextAccept
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

	if got := decoded.MessageType(); got != msg.ESMModifyEPSBearerContextAccept {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.ESMModifyEPSBearerContextAccept)
	}
	if got := decoded.MessageTypeName(); got != "Modify EPS bearer context accept" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Modify EPS bearer context accept")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if decoded.EPSQoS == nil {
		t.Fatal("EPSQoS is nil")
	}

	if decoded.EsmHeader.BearerIdentity != 0x01 {
		t.Errorf("BearerIdentity = %#02x, want 0x01", decoded.EsmHeader.BearerIdentity)
	}
	if decoded.EsmHeader.ProcedureTransactionIdentity != 0x01 {
		t.Errorf("PTI = %#02x, want 0x01", decoded.EsmHeader.ProcedureTransactionIdentity)
	}
}

func TestModifyEPSBearerContextAcceptMultiplePTI(t *testing.T) {
	t.Parallel()

	qosFields := &ie.EPSQualityOfServiceFields{QCI: 5}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}

	for _, bearerID := range []uint8{0x01, 0x02, 0x05, 0x07} {
		for _, pti := range []uint8{0x01, 0x02, 0x0F} {
			t.Run("", func(t *testing.T) {
				t.Parallel()

				h := msg.EsmHeader{
					BearerIdentity:               bearerID,
					ProtocolDiscriminator:        byte(msg.PdEsm),
					ProcedureTransactionIdentity: pti,
					MsgType:                      msg.ESMModifyEPSBearerContextAccept,
				}

				esmQoS := ie.NewPositionalEPSQualityOfService(qosPayload)

				m := msg.NewModifyEPSBearerContextAccept(esmQoS)
				m.EsmHeader = h

				buf := make([]byte, m.MarshalLen())
				if err := m.MarshalTo(buf); err != nil {
					t.Fatalf("MarshalTo failed: %v", err)
				}

				var decoded msg.ModifyEPSBearerContextAccept
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

func TestModifyEPSBearerContextAcceptEmptyOptional(t *testing.T) {
	t.Parallel()

	qosFields := &ie.EPSQualityOfServiceFields{QCI: 9}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}

	h := msg.EsmHeader{
		BearerIdentity:               0x03,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x03,
		MsgType:                      msg.ESMModifyEPSBearerContextAccept,
	}

	esmQoS := ie.NewPositionalEPSQualityOfService(qosPayload)

	m := msg.NewModifyEPSBearerContextAccept(esmQoS)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ModifyEPSBearerContextAccept
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
	if decoded.ProtocolConfigOpts != nil {
		t.Error("unexpected ProtocolConfigOpts in empty optional")
	}
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs in empty optional")
	}
}

func TestModifyEPSBearerContextAcceptErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.ModifyEPSBearerContextAccept

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

	local := &msg.ModifyEPSBearerContextAccept{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.ModifyEPSBearerContextAccept

	// Truncated ESM header (2 bytes, need 3).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Invalid PD in 3-octet header.
	err = decoded.UnmarshalBinary([]byte{0xAA, 0x01, 0xCA})
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid PD = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}

	// Valid header but no payload — need at least QoS (2 bytes min).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, 0xCA})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated no QoS = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Valid header + 1 byte QoS length but no payload.
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, 0xCA, 0x01})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated no QoS payload = %v, want %T", err, io.ErrUnexpectedEOF)
	}
}

func TestModifyEPSBearerContextAcceptParseHelper(t *testing.T) {
	t.Parallel()

	qosFields := &ie.EPSQualityOfServiceFields{QCI: 9}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}

	h := msg.EsmHeader{
		BearerIdentity:               0x05,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x06,
		MsgType:                      msg.ESMModifyEPSBearerContextAccept,
	}

	esmQoS := ie.NewPositionalEPSQualityOfService(qosPayload)

	m := msg.NewModifyEPSBearerContextAccept(esmQoS)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	parsed, err := msg.ParseModifyEPSBearerContextAccept(buf)
	if err != nil {
		t.Fatalf("ParseModifyEPSBearerContextAccept failed: %v", err)
	}

	if parsed == nil {
		t.Fatal("ParseModifyEPSBearerContextAccept returned nil")
	}

	if parsed.EPSQoS == nil {
		t.Fatal("EPSQoS is nil")
	}

	if parsed.EsmHeader.BearerIdentity != 0x05 {
		t.Errorf("BearerIdentity = %#02x, want 0x05", parsed.EsmHeader.BearerIdentity)
	}
	if parsed.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if parsed.MessageType() != msg.ESMModifyEPSBearerContextAccept {
		t.Errorf("MessageType = %#02x, want %#02x", parsed.MessageType(), msg.ESMModifyEPSBearerContextAccept)
	}
}

func TestModifyEPSBearerContextAcceptWithExtendedQoSAndPCO(t *testing.T) {
	t.Parallel()

	qosFields := &ie.EPSQualityOfServiceFields{QCI: 1}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}

	extQOSPayload := make([]byte, 10)
	extQOSPayload[4] = 1

	extQoS := ie.NewExtendedQualityOfService(extQOSPayload)
	pco := ie.NewProtocolConfigurationOptions([]byte{0x50, 0x01, 0x21})
	esmQoS := ie.NewPositionalEPSQualityOfService(qosPayload)

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x01,
		MsgType:                      msg.ESMModifyEPSBearerContextAccept,
	}

	m := msg.NewModifyEPSBearerContextAccept(esmQoS, extQoS, pco)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ModifyEPSBearerContextAccept
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
	if decoded.ProtocolConfigOpts == nil {
		t.Fatal("ProtocolConfigOpts is nil after round-trip")
	}
	if len(decoded.ProtocolConfigOpts.Payload) != 3 {
		t.Errorf("ProtocolConfigOpts payload len = %d, want 3", len(decoded.ProtocolConfigOpts.Payload))
	}
}
