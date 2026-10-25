package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestActivateDedicatedEPSBearerContextAcceptRoundTrip(t *testing.T) {
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
		MsgType:                      msg.ESMActivateDedicatedEPSBearerContextAccept,
	}

	esmQoS := ie.NewPositionalEPSQualityOfService(qosPayload)

	// Tagged optional: LLC SAPI, ESM cause and Traffic Flow Template.
	llcSAPI := ie.NewRaw(0x32, ie.IEFormatFour, []byte{5})
	esmCause := ie.NewRaw(0x58, ie.IEFormatFour, []byte{byte(ie.RequestRejectedUnspecifiedSmCause)})
	tft := ie.NewRaw(0x36, ie.IEFormatFour, []byte{0xAA, 0xBB, 0xCC})

	m := msg.NewActivateDedicatedEPSBearerContextAccept(esmQoS, llcSAPI, esmCause, tft)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	if len(buf) != m.MarshalLen() {
		t.Fatalf("marshal length = %d, want %d", len(buf), m.MarshalLen())
	}

	var decoded msg.ActivateDedicatedEPSBearerContextAccept
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

	if got := decoded.MessageType(); got != msg.ESMActivateDedicatedEPSBearerContextAccept {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.ESMActivateDedicatedEPSBearerContextAccept)
	}
	if got := decoded.MessageTypeName(); got != "Activate dedicated EPS bearer context accept" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Activate dedicated EPS bearer context accept")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}

	if decoded.EPSQoS == nil {
		t.Fatal("EPSQoS is nil")
	}
	if decoded.LLCSAPI == nil {
		t.Error("LLCSAPI is nil")
	} else if len(decoded.LLCSAPI.Payload) != 1 || decoded.LLCSAPI.Payload[0] != 5 {
		t.Errorf("LLCSAPI = %v, want [5]", decoded.LLCSAPI.Payload)
	}
	if decoded.ESMCause == nil {
		t.Error("ESMCause is nil")
	} else if len(decoded.ESMCause.Payload) != 1 {
		t.Errorf("ESMCause payload len = %d, want 1", len(decoded.ESMCause.Payload))
	}
	if decoded.TrafficFlowTemplate == nil {
		t.Error("TrafficFlowTemplate is nil")
	}

	if decoded.EsmHeader.BearerIdentity != 0x01 {
		t.Errorf("BearerIdentity = %#02x, want 0x01", decoded.EsmHeader.BearerIdentity)
	}
	if decoded.EsmHeader.ProcedureTransactionIdentity != 0x01 {
		t.Errorf("PTI = %#02x, want 0x01", decoded.EsmHeader.ProcedureTransactionIdentity)
	}
}

func TestActivateDedicatedEPSBearerContextAcceptMultiplePTI(t *testing.T) {
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
					MsgType:                      msg.ESMActivateDedicatedEPSBearerContextAccept,
				}

				esmQoS := ie.NewPositionalEPSQualityOfService(qosPayload)

				m := msg.NewActivateDedicatedEPSBearerContextAccept(esmQoS)
				m.EsmHeader = h

				buf := make([]byte, m.MarshalLen())
				if err := m.MarshalTo(buf); err != nil {
					t.Fatalf("MarshalTo failed: %v", err)
				}

				var decoded msg.ActivateDedicatedEPSBearerContextAccept
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

func TestActivateDedicatedEPSBearerContextAcceptEmptyOptional(t *testing.T) {
	t.Parallel()

	qosFields := &ie.EPSQualityOfServiceFields{QCI: 5}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x02,
		MsgType:                      msg.ESMActivateDedicatedEPSBearerContextAccept,
	}

	esmQoS := ie.NewPositionalEPSQualityOfService(qosPayload)

	m := msg.NewActivateDedicatedEPSBearerContextAccept(esmQoS)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ActivateDedicatedEPSBearerContextAccept
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

	if decoded.LLCSAPI != nil {
		t.Error("unexpected LLCSAPI in empty optional")
	}
	if decoded.ESMCause != nil {
		t.Error("unexpected ESMCause in empty optional")
	}
	if decoded.TrafficFlowTemplate != nil {
		t.Error("unexpected TrafficFlowTemplate in empty optional")
	}
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs in empty optional")
	}
}

func TestActivateDedicatedEPSBearerContextAcceptOnlyLLCSAPI(t *testing.T) {
	t.Parallel()

	qosFields := &ie.EPSQualityOfServiceFields{QCI: 9}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}

	h := msg.EsmHeader{
		BearerIdentity:               0x03,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x05,
		MsgType:                      msg.ESMActivateDedicatedEPSBearerContextAccept,
	}

	esmQoS := ie.NewPositionalEPSQualityOfService(qosPayload)
	llcSAPI := ie.NewRaw(0x32, ie.IEFormatFour, []byte{3})

	m := msg.NewActivateDedicatedEPSBearerContextAccept(esmQoS, llcSAPI)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ActivateDedicatedEPSBearerContextAccept
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

	if decoded.LLCSAPI == nil {
		t.Error("LLCSAPI is nil")
	} else if len(decoded.LLCSAPI.Payload) != 1 || decoded.LLCSAPI.Payload[0] != 3 {
		t.Errorf("LLCSAPI value = %v, want [3]", decoded.LLCSAPI.Payload)
	}

	if decoded.ESMCause != nil {
		t.Error("unexpected ESMCause")
	}
	if decoded.TrafficFlowTemplate != nil {
		t.Error("unexpected TrafficFlowTemplate")
	}
}

func TestActivateDedicatedEPSBearerContextAcceptErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.ActivateDedicatedEPSBearerContextAccept

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

	local := &msg.ActivateDedicatedEPSBearerContextAccept{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.ActivateDedicatedEPSBearerContextAccept

	// Truncated ESM header (2 bytes, need 3).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Invalid PD in 3-octet header.
	err = decoded.UnmarshalBinary([]byte{0xAA, 0x01, 0xC6})
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid PD = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}

	// Valid header but no payload — need QoS length byte.
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, 0xC6})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated no QoS = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Valid header + 1 byte QoS length, no QoS payload.
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, 0xC6, 0x02})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated no QoS payload = %v, want %T", err, io.ErrUnexpectedEOF)
	}
}

func TestActivateDedicatedEPSBearerContextAcceptParseHelper(t *testing.T) {
	t.Parallel()

	qosFields := &ie.EPSQualityOfServiceFields{QCI: 9}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}

	h := msg.EsmHeader{
		BearerIdentity:               0x04,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x06,
		MsgType:                      msg.ESMActivateDedicatedEPSBearerContextAccept,
	}

	esmQoS := ie.NewPositionalEPSQualityOfService(qosPayload)

	m := msg.NewActivateDedicatedEPSBearerContextAccept(esmQoS)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	parsed, err := msg.ParseActivateDedicatedEPSBearerContextAccept(buf)
	if err != nil {
		t.Fatalf("ParseActivateDedicatedEPSBearerContextAccept failed: %v", err)
	}

	if parsed == nil {
		t.Fatal("ParseActivateDedicatedEPSBearerContextAccept returned nil")
	}

	if parsed.EPSQoS == nil {
		t.Fatal("EPSQoS is nil")
	}
	if parsed.EsmHeader.BearerIdentity != 0x04 {
		t.Errorf("BearerIdentity = %#02x, want 0x04", parsed.EsmHeader.BearerIdentity)
	}
	if parsed.EsmHeader.ProcedureTransactionIdentity != 0x06 {
		t.Errorf("PTI = %#02x, want 0x06", parsed.EsmHeader.ProcedureTransactionIdentity)
	}
	if parsed.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if parsed.MessageType() != msg.ESMActivateDedicatedEPSBearerContextAccept {
		t.Errorf("MessageType = %#02x, want %#02x", parsed.MessageType(), msg.ESMActivateDedicatedEPSBearerContextAccept)
	}
}

func TestActivateDedicatedEPSBearerContextAcceptWithDuplicates(t *testing.T) {
	t.Parallel()

	qosFields := &ie.EPSQualityOfServiceFields{QCI: 9}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}

	esmQoS := ie.NewPositionalEPSQualityOfService(qosPayload)

	// Two ESM cause IEs.
	cause1 := ie.NewESMCause(ie.RequestRejectedUnspecifiedSmCause)
	dupCause := ie.NewESMCause(ie.InsufficientResourcesSmCause)

	m := msg.NewActivateDedicatedEPSBearerContextAccept(esmQoS, cause1, dupCause)

	if m.ESMCause == nil {
		t.Fatal("first ESMCause was not stored")
	}

	var foundDup bool
	for _, i := range m.IEs {
		if i.Kind == ie.ESMCause && len(i.Payload) == 1 && i.Payload[0] == byte(ie.InsufficientResourcesSmCause) {
			foundDup = true
		}
	}

	if !foundDup {
		t.Error("expected duplicate ESMCause in IEs")
	}
}
