package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestPDNConnectivityRequestRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x01,
		MsgType:                      msg.ESMPDNConnectivityRequest,
	}

	pdnTypeIE := ie.NewPositionalPDNType(byte(ie.PDNTypeIPv4))
	reqTypeIE := ie.NewPositionalRequestType(byte(ie.RequestTypeInitial))
	qosFields := &ie.EPSQualityOfServiceFields{QCI: 9}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}
	ePSQoS := ie.NewPositionalEPSQualityOfService(qosPayload)
	apnFields := &ie.AccessPointNameFields{Name: "internet"}
	apnIE, err := ie.NewPositionalAccessPointNameWithFields(apnFields)
	if err != nil {
		t.Fatalf("APN marshal failed: %v", err)
	}

	m := msg.NewPDNConnectivityRequest(pdnTypeIE, reqTypeIE, ePSQoS, apnIE)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.PDNConnectivityRequest
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

	if got := decoded.MessageType(); got != msg.ESMPDNConnectivityRequest {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.ESMPDNConnectivityRequest)
	}
	if got := decoded.MessageTypeName(); got != "PDN connectivity request" {
		t.Errorf("MessageTypeName = %q, want %q", got, "PDN connectivity request")
	}
	if !decoded.IsRequest() {
		t.Error("IsRequest = false, want true")
	}
	if decoded.PDNType == nil {
		t.Fatal("PDNType is nil")
	}
	if decoded.RequestType == nil {
		t.Fatal("RequestType is nil")
	}
	if decoded.EPSQoS == nil {
		t.Fatal("EPSQoS is nil")
	}
	if decoded.AccessPointName == nil {
		t.Fatal("AccessPointName is nil")
	}

	if decoded.EsmHeader.BearerIdentity != 0x01 {
		t.Errorf("BearerIdentity = %#02x, want 0x01", decoded.EsmHeader.BearerIdentity)
	}
}

func TestPDNConnectivityRequestMultiplePTI(t *testing.T) {
	t.Parallel()

	for _, bearerID := range []uint8{0x01, 0x02, 0x05, 0x07} {
		for _, pti := range []uint8{0x01, 0x02, 0x0F} {
			t.Run("", func(t *testing.T) {
				t.Parallel()

				h := msg.EsmHeader{
					BearerIdentity:               bearerID,
					ProtocolDiscriminator:        byte(msg.PdEsm),
					ProcedureTransactionIdentity: pti,
					MsgType:                      msg.ESMPDNConnectivityRequest,
				}

				pdnTypeIE := ie.NewPositionalPDNType(byte(ie.PDNTypeIPv6))
				reqTypeIE := ie.NewPositionalRequestType(byte(ie.RequestTypeEmergency))
				qosFields := &ie.EPSQualityOfServiceFields{QCI: 5}
				qosPayload, err := qosFields.Marshal()
				if err != nil {
					t.Fatalf("QoS marshal failed: %v", err)
				}
				ePSQoS := ie.NewPositionalEPSQualityOfService(qosPayload)
				apnFields := &ie.AccessPointNameFields{Name: "mms"}
				apnIE, err := ie.NewPositionalAccessPointNameWithFields(apnFields)
				if err != nil {
					t.Fatalf("APN marshal failed: %v", err)
				}

				m := msg.NewPDNConnectivityRequest(pdnTypeIE, reqTypeIE, ePSQoS, apnIE)
				m.EsmHeader = h

				buf := make([]byte, m.MarshalLen())
				if err := m.MarshalTo(buf); err != nil {
					t.Fatalf("MarshalTo failed: %v", err)
				}

				var decoded msg.PDNConnectivityRequest
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

func TestPDNConnectivityRequestEmptyOptional(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x03,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x03,
		MsgType:                      msg.ESMPDNConnectivityRequest,
	}

	pdnTypeIE := ie.NewPositionalPDNType(byte(ie.PDNTypeIPv4))
	reqTypeIE := ie.NewPositionalRequestType(byte(ie.RequestTypeInitial))
	qosFields := &ie.EPSQualityOfServiceFields{QCI: 9}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}
	ePSQoS := ie.NewPositionalEPSQualityOfService(qosPayload)
	apnFields := &ie.AccessPointNameFields{Name: "internet"}
	apnIE, err := ie.NewPositionalAccessPointNameWithFields(apnFields)
	if err != nil {
		t.Fatalf("APN marshal failed: %v", err)
	}

	m := msg.NewPDNConnectivityRequest(pdnTypeIE, reqTypeIE, ePSQoS, apnIE)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.PDNConnectivityRequest
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

	if decoded.ESMCause != nil {
		t.Error("unexpected ESMCause in empty optional")
	}
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs in empty optional")
	}
}

func TestPDNConnectivityRequestWithESMCause(t *testing.T) {
	t.Parallel()

	esmCause := ie.NewRaw(0x58, ie.IEFormatFour, []byte{1, byte(ie.NetworkFailureSmCause)})
	pdnTypeIE := ie.NewPositionalPDNType(byte(ie.PDNTypeIPv4))
	reqTypeIE := ie.NewPositionalRequestType(byte(ie.RequestTypeInitial))
	qosFields := &ie.EPSQualityOfServiceFields{QCI: 9}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}
	ePSQoS := ie.NewPositionalEPSQualityOfService(qosPayload)
	apnFields := &ie.AccessPointNameFields{Name: "internet"}
	apnIE, err := ie.NewPositionalAccessPointNameWithFields(apnFields)
	if err != nil {
		t.Fatalf("APN marshal failed: %v", err)
	}

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x01,
		MsgType:                      msg.ESMPDNConnectivityRequest,
	}

	m := msg.NewPDNConnectivityRequest(pdnTypeIE, reqTypeIE, ePSQoS, apnIE, esmCause)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.PDNConnectivityRequest
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

	if decoded.ESMCause == nil {
		t.Fatal("ESMCause is nil after round-trip")
	}
}

func TestPDNConnectivityRequestDuplicateIEs(t *testing.T) {
	t.Parallel()

	c1 := ie.NewESMCause(ie.NetworkFailureSmCause)
	c2 := ie.NewESMCause(ie.ServiceOptionNotSupportedSmCause)
	pdnTypeIE := ie.NewPositionalPDNType(byte(ie.PDNTypeIPv4))
	reqTypeIE := ie.NewPositionalRequestType(byte(ie.RequestTypeInitial))
	qosFields := &ie.EPSQualityOfServiceFields{QCI: 9}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}
	ePSQoS := ie.NewPositionalEPSQualityOfService(qosPayload)
	apnFields := &ie.AccessPointNameFields{Name: "internet"}
	apnIE, err := ie.NewPositionalAccessPointNameWithFields(apnFields)
	if err != nil {
		t.Fatalf("APN marshal failed: %v", err)
	}

	m := msg.NewPDNConnectivityRequest(pdnTypeIE, reqTypeIE, ePSQoS, apnIE, c1, c2)

	if m.ESMCause == nil {
		t.Fatal("first ESMCause was not stored")
	}

	var foundDup bool
	for _, i := range m.IEs {
		if i.Kind == ie.ESMCause && len(i.Payload) == 1 && i.Payload[0] == byte(ie.ServiceOptionNotSupportedSmCause) {
			foundDup = true
		}
	}

	if !foundDup {
		t.Error("expected duplicate ESMCause in IEs")
	}
}

func TestPDNConnectivityRequestErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.PDNConnectivityRequest

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

	local := &msg.PDNConnectivityRequest{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.PDNConnectivityRequest

	// Truncated ESM header (2 bytes, need 3).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Invalid PD in 3-octet header.
	err = decoded.UnmarshalBinary([]byte{0xAA, 0x01, 0xD0})
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid PD = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}

	// Valid header but no payload — need packed half-octet (1 byte).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, 0xD0})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated no half-octet = %v, want %T", err, io.ErrUnexpectedEOF)
	}
}

func TestPDNConnectivityRequestParseHelper(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x05,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x06,
		MsgType:                      msg.ESMPDNConnectivityRequest,
	}

	pdnTypeIE := ie.NewPositionalPDNType(byte(ie.PDNTypeIPv4))
	reqTypeIE := ie.NewPositionalRequestType(byte(ie.RequestTypeInitial))
	qosFields := &ie.EPSQualityOfServiceFields{QCI: 9}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}
	ePSQoS := ie.NewPositionalEPSQualityOfService(qosPayload)
	apnFields := &ie.AccessPointNameFields{Name: "internet"}
	apnIE, err := ie.NewPositionalAccessPointNameWithFields(apnFields)
	if err != nil {
		t.Fatalf("APN marshal failed: %v", err)
	}

	m := msg.NewPDNConnectivityRequest(pdnTypeIE, reqTypeIE, ePSQoS, apnIE)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	parsed, err := msg.ParsePDNConnectivityRequest(buf)
	if err != nil {
		t.Fatalf("ParsePDNConnectivityRequest failed: %v", err)
	}

	if parsed == nil {
		t.Fatal("ParsePDNConnectivityRequest returned nil")
	}

	if parsed.PDNType == nil {
		t.Fatal("PDNType is nil")
	}
	if parsed.RequestType == nil {
		t.Fatal("RequestType is nil")
	}
	if !parsed.IsRequest() {
		t.Error("IsRequest = false, want true")
	}
	if parsed.MessageType() != msg.ESMPDNConnectivityRequest {
		t.Errorf("MessageType = %#02x, want %#02x", parsed.MessageType(), msg.ESMPDNConnectivityRequest)
	}
	if parsed.EsmHeader.BearerIdentity != 0x05 {
		t.Errorf("BearerIdentity = %#02x, want 0x05", parsed.EsmHeader.BearerIdentity)
	}
}
