package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestActivateDefaultEPSBearerContextRequestRoundTrip(t *testing.T) {
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
		MsgType:                      msg.ESMActivateDefaultEPSBearerContextRequest,
	}

	esmQoS := ie.NewPositionalEPSQualityOfService(qosPayload)
	apn := ie.NewPositionalAccessPointName([]byte{7, 'e', 'x', 'a', 'm', 'p', 'l', 'e'})
	// PDN address: type 0x01 (IPv4) + 4 bytes = 5 octets
	pdnAddr := ie.NewPositionalPDNAddress([]byte{0x01, 0xAC, 0x10, 0x0A, 0x01})

	// Use raw constructors for ambiguous tags (0x32, 0x58) with IEFormatFour
	// to match the prototypes used in ie.ParseMultiIEs.
	llcSAPI := ie.NewRaw(0x32, ie.IEFormatFour, []byte{5})
	esmCause := ie.NewRaw(0x58, ie.IEFormatFour, []byte{byte(ie.RequestRejectedUnspecifiedSmCause)})

	m := msg.NewActivateDefaultEPSBearerContextRequest(esmQoS, apn, pdnAddr, llcSAPI, esmCause)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ActivateDefaultEPSBearerContextRequest
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

	if got := decoded.MessageType(); got != msg.ESMActivateDefaultEPSBearerContextRequest {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.ESMActivateDefaultEPSBearerContextRequest)
	}
	if got := decoded.MessageTypeName(); got != "Activate default EPS bearer context request" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Activate default EPS bearer context request")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}

	if decoded.EPSQoS == nil {
		t.Fatal("EPSQoS is nil")
	}
	if decoded.AccessPointName == nil {
		t.Fatal("AccessPointName is nil")
	}
	if decoded.PDNAddress == nil {
		t.Fatal("PDNAddress is nil")
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

	if decoded.EsmHeader.BearerIdentity != 0x01 {
		t.Errorf("BearerIdentity = %#02x, want 0x01", decoded.EsmHeader.BearerIdentity)
	}
	if decoded.EsmHeader.ProcedureTransactionIdentity != 0x01 {
		t.Errorf("PTI = %#02x, want 0x01", decoded.EsmHeader.ProcedureTransactionIdentity)
	}
}

func TestActivateDefaultEPSBearerContextRequestFullPTI(t *testing.T) {
	t.Parallel()

	qosFields := &ie.EPSQualityOfServiceFields{QCI: 9}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}

	h := msg.EsmHeader{
		BearerIdentity:               0x02,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x05,
		MsgType:                      msg.ESMActivateDefaultEPSBearerContextRequest,
	}

	esmQoS := ie.NewPositionalEPSQualityOfService(qosPayload)
	apn := ie.NewPositionalAccessPointName([]byte{5, 't', 'e', 's', 't'})
	pdnAddr := ie.NewPositionalPDNAddress([]byte{0x01, 0xAC, 0x10, 0x0A, 0x03})

	m := msg.NewActivateDefaultEPSBearerContextRequest(esmQoS, apn, pdnAddr)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ActivateDefaultEPSBearerContextRequest
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

	if decoded.EsmHeader.BearerIdentity != 0x02 {
		t.Errorf("BearerIdentity = %#02x, want 0x02", decoded.EsmHeader.BearerIdentity)
	}
	if decoded.EsmHeader.ProcedureTransactionIdentity != 0x05 {
		t.Errorf("PTI = %#02x, want 0x05", decoded.EsmHeader.ProcedureTransactionIdentity)
	}
	if decoded.LLCSAPI != nil {
		t.Error("unexpected LLC SAPI")
	}
	if decoded.ESMCause != nil {
		t.Error("unexpected ESM cause")
	}
}

func TestActivateDefaultEPSBearerContextRequestEmptyOptional(t *testing.T) {
	t.Parallel()

	qosFields := &ie.EPSQualityOfServiceFields{QCI: 5}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x03,
		MsgType:                      msg.ESMActivateDefaultEPSBearerContextRequest,
	}

	esmQoS := ie.NewPositionalEPSQualityOfService(qosPayload)
	apn := ie.NewPositionalAccessPointName([]byte{4, 'd', 'a', 't', 'a'})
	pdnAddr := ie.NewPositionalPDNAddress([]byte{0x01, 0xAC, 0x10, 0x0A, 0x05})

	m := msg.NewActivateDefaultEPSBearerContextRequest(esmQoS, apn, pdnAddr)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ActivateDefaultEPSBearerContextRequest
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
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs in empty optional")
	}
}

func TestActivateDefaultEPSBearerContextRequestErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.ActivateDefaultEPSBearerContextRequest

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

	local := &msg.ActivateDefaultEPSBearerContextRequest{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.ActivateDefaultEPSBearerContextRequest
	// 3-byte header with invalid PD (0x07) → error at UnmarshalBinary level
	err = decoded.UnmarshalBinary([]byte{0x7E, 0xC1}) // plain EMM header (2 octets), missing third
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Invalid PD in 3-octet header: PD must be 0x2E for ESM
	err = decoded.UnmarshalBinary([]byte{0xAA, 0x01, 0xC1})
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid PD = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}

	// Truncated after header (valid PD but no QoS payload)
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, 0xC1})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated no QoS = %v, want %T", err, io.ErrUnexpectedEOF)
	}
}

func TestActivateDefaultEPSBearerContextRequestHeaderWithQoS(t *testing.T) {
	t.Parallel()

	qosFields := &ie.EPSQualityOfServiceFields{QCI: 9}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x02,
		MsgType:                      msg.ESMActivateDefaultEPSBearerContextRequest,
	}

	esmQoS := ie.NewPositionalEPSQualityOfService(qosPayload)
	apn := ie.NewPositionalAccessPointName([]byte{7, 'e', 'x', 'a', 'm', 'p', 'l', 'e'})
	pdnAddr := ie.NewPositionalPDNAddress([]byte{0x01, 0xAC, 0x10, 0x0A, 0x07})

	// Only LLCSAPI (no ESM cause)
	llcSAPI := ie.NewRaw(0x32, ie.IEFormatFour, []byte{3})

	m := msg.NewActivateDefaultEPSBearerContextRequest(esmQoS, apn, pdnAddr, llcSAPI)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ActivateDefaultEPSBearerContextRequest
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
}
