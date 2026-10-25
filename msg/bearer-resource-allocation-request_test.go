package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestBearerResourceAllocationRequestRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x01,
		MsgType:                      msg.ESMBearerResourceAllocationRequest,
	}

	bi := ie.NewPositionalLinkedEPSBearerIdentity(0x01)
	sp := ie.NewPositionalSpareHalfOctet(0x0F)
	qosFields := &ie.EPSQualityOfServiceFields{QCI: 9}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}
	ePSQoS := ie.NewPositionalEPSQualityOfService(qosPayload)
	tftIE := ie.NewPositionalTrafficFlowTemplate([]byte{0x01, 0xFF})

	m := msg.NewBearerResourceAllocationRequest(bi, sp, ePSQoS, tftIE)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.BearerResourceAllocationRequest
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

	if got := decoded.MessageType(); got != msg.ESMBearerResourceAllocationRequest {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.ESMBearerResourceAllocationRequest)
	}
	if got := decoded.MessageTypeName(); got != "Bearer resource allocation request" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Bearer resource allocation request")
	}
	if !decoded.IsRequest() {
		t.Error("IsRequest = false, want true")
	}
	if decoded.BearerIdentity == nil {
		t.Fatal("BearerIdentity is nil")
	}
	if decoded.SpareHalfOctet == nil {
		t.Fatal("SpareHalfOctet is nil")
	}
	if decoded.EPSQoS == nil {
		t.Fatal("EPSQoS is nil")
	}

	if decoded.EsmHeader.BearerIdentity != 0x01 {
		t.Errorf("BearerIdentity = %#02x, want 0x01", decoded.EsmHeader.BearerIdentity)
	}
}

func TestBearerResourceAllocationRequestWithLLCSAPIAndExtendedQoS(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x02,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x02,
		MsgType:                      msg.ESMBearerResourceAllocationRequest,
	}

	bi := ie.NewPositionalLinkedEPSBearerIdentity(0x02)
	sp := ie.NewPositionalSpareHalfOctet(0x0F)
	qosFields := &ie.EPSQualityOfServiceFields{QCI: 5}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}
	ePSQoS := ie.NewPositionalEPSQualityOfService(qosPayload)
	tftIE := ie.NewPositionalTrafficFlowTemplate([]byte{0x02, 0xFE})

	llcSAPI := ie.NewLLCServiceAccessPointIdentifier(3)
	extQoSFields := &ie.ExtendedQualityOfServiceFields{QCI: 9}
	extQoSIE, err := ie.NewExtendedQualityOfServiceWithFields(extQoSFields)
	if err != nil {
		t.Fatalf("Extended QoS marshal failed: %v", err)
	}

	m := msg.NewBearerResourceAllocationRequest(bi, sp, ePSQoS, tftIE, llcSAPI, extQoSIE)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.BearerResourceAllocationRequest
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
		t.Fatal("LLCSAPI is nil after round-trip")
	}
	if decoded.ExtendedQoS == nil {
		t.Fatal("ExtendedQoS is nil after round-trip")
	}
	if decoded.EsmHeader.BearerIdentity != 0x02 {
		t.Errorf("BearerIdentity = %#02x, want 0x02", decoded.EsmHeader.BearerIdentity)
	}
}

func TestBearerResourceAllocationRequestMultiplePTI(t *testing.T) {
	t.Parallel()

	for _, bearerID := range []uint8{0x01, 0x02, 0x05, 0x07} {
		for _, pti := range []uint8{0x01, 0x02, 0x0F} {
			t.Run("", func(t *testing.T) {
				t.Parallel()

				h := msg.EsmHeader{
					BearerIdentity:               bearerID,
					ProtocolDiscriminator:        byte(msg.PdEsm),
					ProcedureTransactionIdentity: pti,
					MsgType:                      msg.ESMBearerResourceAllocationRequest,
				}

				bi := ie.NewPositionalLinkedEPSBearerIdentity(bearerID)
				sp := ie.NewPositionalSpareHalfOctet(0x0F)
				qosFields := &ie.EPSQualityOfServiceFields{QCI: 9}
				qosPayload, err := qosFields.Marshal()
				if err != nil {
					t.Fatalf("QoS marshal failed: %v", err)
				}
				ePSQoS := ie.NewPositionalEPSQualityOfService(qosPayload)
				tftIE := ie.NewPositionalTrafficFlowTemplate([]byte{0x01, 0xFF})

				m := msg.NewBearerResourceAllocationRequest(bi, sp, ePSQoS, tftIE)
				m.EsmHeader = h

				buf := make([]byte, m.MarshalLen())
				if err := m.MarshalTo(buf); err != nil {
					t.Fatalf("MarshalTo failed: %v", err)
				}

				var decoded msg.BearerResourceAllocationRequest
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

func TestBearerResourceAllocationRequestDuplicateIEs(t *testing.T) {
	t.Parallel()

	c1 := ie.NewLLCServiceAccessPointIdentifier(3)
	c2 := ie.NewLLCServiceAccessPointIdentifier(5)
	bi := ie.NewPositionalLinkedEPSBearerIdentity(0x01)
	sp := ie.NewPositionalSpareHalfOctet(0x0F)
	qosFields := &ie.EPSQualityOfServiceFields{QCI: 9}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}
	ePSQoS := ie.NewPositionalEPSQualityOfService(qosPayload)
	tftIE := ie.NewPositionalTrafficFlowTemplate([]byte{0x01, 0xFF})

	m := msg.NewBearerResourceAllocationRequest(bi, sp, ePSQoS, tftIE, c1, c2)

	if m.LLCSAPI == nil {
		t.Fatal("first LLCSAPI was not stored")
	}

	var foundDup bool
	for _, i := range m.IEs {
		if i.Kind == ie.LLCServiceAccessPointIdentifier && len(i.Payload) == 1 && i.Payload[0] == 5 {
			foundDup = true
		}
	}

	if !foundDup {
		t.Error("expected duplicate LLCSAPI in IEs")
	}
}

func TestBearerResourceAllocationRequestErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.BearerResourceAllocationRequest

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

	local := &msg.BearerResourceAllocationRequest{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.BearerResourceAllocationRequest

	// Truncated ESM header (2 bytes, need 3).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Invalid PD in 3-octet header.
	err = decoded.UnmarshalBinary([]byte{0xAA, 0x01, 0xD4})
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid PD = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}

	// Valid header but no payload — need packed half-octet (1 byte).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, 0xD4})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated no half-octet = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Truncated after header (need EPS QoS).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, 0xD4, 0})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated no qos = %v, want %T", err, io.ErrUnexpectedEOF)
	}
}

func TestBearerResourceAllocationRequestParseHelper(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x05,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x06,
		MsgType:                      msg.ESMBearerResourceAllocationRequest,
	}

	bi := ie.NewPositionalLinkedEPSBearerIdentity(0x05)
	sp := ie.NewPositionalSpareHalfOctet(0x0F)
	qosFields := &ie.EPSQualityOfServiceFields{QCI: 9}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}
	ePSQoS := ie.NewPositionalEPSQualityOfService(qosPayload)
	tftIE := ie.NewPositionalTrafficFlowTemplate([]byte{0x01, 0xFF})

	m := msg.NewBearerResourceAllocationRequest(bi, sp, ePSQoS, tftIE)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	parsed, err := msg.ParseBearerResourceAllocationRequest(buf)
	if err != nil {
		t.Fatalf("ParseBearerResourceAllocationRequest failed: %v", err)
	}

	if parsed == nil {
		t.Fatal("ParseBearerResourceAllocationRequest returned nil")
	}

	if parsed.BearerIdentity == nil {
		t.Fatal("BearerIdentity is nil")
	}
	if parsed.IsRequest() != true {
		t.Error("IsRequest = false, want true")
	}
	if parsed.MessageType() != msg.ESMBearerResourceAllocationRequest {
		t.Errorf("MessageType = %#02x, want %#02x", parsed.MessageType(), msg.ESMBearerResourceAllocationRequest)
	}
	if parsed.EsmHeader.BearerIdentity != 0x05 {
		t.Errorf("BearerIdentity = %#02x, want 0x05", parsed.EsmHeader.BearerIdentity)
	}
}

func TestBearerResourceAllocationRequestMarshalLenExtraIEs(t *testing.T) {
	t.Parallel()

	collisionIE := &ie.IE{Kind: 0x32, Format: ie.IEFormatFour, Payload: []byte{0xAA}}

	bi := ie.NewPositionalLinkedEPSBearerIdentity(0x01)
	sho := ie.NewPositionalSpareHalfOctet(0x0F)
	qosFields := &ie.EPSQualityOfServiceFields{QCI: 9}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal: %v", err)
	}
	qos := ie.NewPositionalEPSQualityOfService(qosPayload)
	tft := ie.NewTrafficFlowTemplate([]byte{0xAA, 0xBB})

	m := &msg.BearerResourceAllocationRequest{
		EsmHeader:           msg.EsmHeader{BearerIdentity: 0x01, ProtocolDiscriminator: byte(msg.PdEsm), ProcedureTransactionIdentity: 0x01, MsgType: msg.ESMBearerResourceAllocationRequest},
		BearerIdentity:      bi,
		SpareHalfOctet:      sho,
		EPSQoS:              qos,
		TrafficFlowTemplate: tft,
		LLCSAPI:             collisionIE, // tag 0x32 resolves to LLCSAPI
		IEs:                 []*ie.IE{collisionIE}, // duplicate goes to IEs
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo with extra IEs failed: %v", err)
	}

	var decoded msg.BearerResourceAllocationRequest
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

	// Collision TLV short-buffer error path for extra IEs loop.
	collisionMsg := &msg.BearerResourceAllocationRequest{
		EsmHeader:           msg.EsmHeader{BearerIdentity: 0x01, ProtocolDiscriminator: byte(msg.PdEsm), ProcedureTransactionIdentity: 0x01, MsgType: msg.ESMBearerResourceAllocationRequest},
		BearerIdentity:      bi,
		SpareHalfOctet:      sho,
		EPSQoS:              qos,
		TrafficFlowTemplate: tft,
		IEs:                 []*ie.IE{collisionIE},
	}
	shortBuf := make([]byte, 1)
	if err := collisionMsg.MarshalTo(shortBuf); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("expected ErrShortBuffer for collision IE error path, got %v", err)
	}
}
