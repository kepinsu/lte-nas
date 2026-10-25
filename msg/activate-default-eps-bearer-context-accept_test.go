package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestActivateDefaultEPSBearerContextAcceptRoundTrip(t *testing.T) {
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
		MsgType:                      msg.ESMActivateDefaultEPSBearerContextAccept,
	}

	esmQoS := ie.NewPositionalEPSQualityOfService(qosPayload)
	pco := ie.NewProtocolConfigurationOptions([]byte{0x00, 0x05, 0x21, 0x02, 0xC0, 0xA8})

	m := msg.NewActivateDefaultEPSBearerContextAccept(esmQoS, pco)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ActivateDefaultEPSBearerContextAccept
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

	if got := decoded.MessageType(); got != msg.ESMActivateDefaultEPSBearerContextAccept {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.ESMActivateDefaultEPSBearerContextAccept)
	}
	if got := decoded.MessageTypeName(); got != "Activate default EPS bearer context accept" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Activate default EPS bearer context accept")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}

	if decoded.EPSQoS == nil {
		t.Fatal("EPSQoS is nil")
	}
	if decoded.ProtocolConfigurationOpts == nil {
		t.Fatal("ProtocolConfigurationOpts is nil")
	}

	if decoded.EsmHeader.BearerIdentity != 0x01 {
		t.Errorf("BearerIdentity = %#02x, want 0x01", decoded.EsmHeader.BearerIdentity)
	}
	if decoded.EsmHeader.ProcedureTransactionIdentity != 0x01 {
		t.Errorf("PTI = %#02x, want 0x01", decoded.EsmHeader.ProcedureTransactionIdentity)
	}
}

func TestActivateDefaultEPSBearerContextAcceptWithOptional(t *testing.T) {
	t.Parallel()

	qosFields := &ie.EPSQualityOfServiceFields{QCI: 5}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}

	h := msg.EsmHeader{
		BearerIdentity:               0x03,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x07,
		MsgType:                      msg.ESMActivateDefaultEPSBearerContextAccept,
	}

	esmQoS := ie.NewPositionalEPSQualityOfService(qosPayload)
	pco := ie.NewProtocolConfigurationOptions([]byte{0x00, 0x05, 0x21, 0x02, 0xC0, 0xA8})
	extPCO := ie.NewExtendedProtocolConfigurationOptions([]byte{0x10, 0x20})
	tft := ie.NewTrafficFlowTemplate([]byte{0x0F, 0x1A})

	m := msg.NewActivateDefaultEPSBearerContextAccept(esmQoS, pco, extPCO, tft)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ActivateDefaultEPSBearerContextAccept
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

	if decoded.EsmHeader.BearerIdentity != 0x03 {
		t.Errorf("BearerIdentity = %#02x, want 0x03", decoded.EsmHeader.BearerIdentity)
	}
	if decoded.ExtendedPCO == nil {
		t.Error("ExtendedPCO is nil")
	}
	if decoded.TrafficFlowTemplate == nil {
		t.Error("TrafficFlowTemplate is nil")
	}
}

func TestActivateDefaultEPSBearerContextAcceptEmptyOptional(t *testing.T) {
	t.Parallel()

	qosFields := &ie.EPSQualityOfServiceFields{QCI: 9}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}

	esmQoS := ie.NewPositionalEPSQualityOfService(qosPayload)
	pco := ie.NewProtocolConfigurationOptions([]byte{0x00, 0x01, 0x21})

	m := msg.NewActivateDefaultEPSBearerContextAccept(esmQoS, pco)
	m.EsmHeader = msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x02,
		MsgType:                      msg.ESMActivateDefaultEPSBearerContextAccept,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ActivateDefaultEPSBearerContextAccept
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

	if decoded.ExtendedPCO != nil {
		t.Error("unexpected ExtendedPCO in empty optional")
	}
	if decoded.TrafficFlowTemplate != nil {
		t.Error("unexpected TrafficFlowTemplate in empty optional")
	}
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs in empty optional")
	}
}

func TestActivateDefaultEPSBearerContextAcceptErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.ActivateDefaultEPSBearerContextAccept

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

	local := &msg.ActivateDefaultEPSBearerContextAccept{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.ActivateDefaultEPSBearerContextAccept
	// Invalid PD (low nibble 0x0A ≠ 0x02): should get wrapped ie.ErrInvalidValue
	err = decoded.UnmarshalBinary([]byte{0xAA, 0x01, 0xC2})
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid PD = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}

	// Valid header (low nibble 0x2 matches PdEsm) but no payload → UnexpectedEOF
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, 0xC2})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header-only = %v, want %T", err, io.ErrUnexpectedEOF)
	}
}
