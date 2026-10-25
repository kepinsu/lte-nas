package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestActivateDedicatedEPSBearerContextRequestRoundTrip(t *testing.T) {
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
		MsgType:                      msg.ESMActivateDedicatedEPSBearerContextRequest,
	}

	// Packed half-octets: linked bearer identity low nibble, spare high nibble.
	linkedBEID := ie.NewPositionalLinkedEPSBearerIdentity(0x01)
	spare := ie.NewPositionalSpareHalfOctet(0x00)
	esmQoS := ie.NewPositionalEPSQualityOfService(qosPayload)

	// Traffic flow template payload (TFT with one packet filter).
	tft := ie.NewTrafficFlowTemplate([]byte{0x12, 0x34, 0x56})

	// Tagged optional: LLC SAPI via raw prototype for ambiguous tag 0x32.
	llcSAPI := ie.NewRaw(0x32, ie.IEFormatFour, []byte{5})

	m := msg.NewActivateDedicatedEPSBearerContextRequest(
		linkedBEID, spare, esmQoS, tft, llcSAPI,
	)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	if len(buf) != m.MarshalLen() {
		t.Fatalf("marshal length = %d, want %d", len(buf), m.MarshalLen())
	}

	var decoded msg.ActivateDedicatedEPSBearerContextRequest
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

	if got := decoded.MessageType(); got != msg.ESMActivateDedicatedEPSBearerContextRequest {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.ESMActivateDedicatedEPSBearerContextRequest)
	}
	if got := decoded.MessageTypeName(); got != "Activate dedicated EPS bearer context request" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Activate dedicated EPS bearer context request")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}

	if decoded.LinkedEPSBearerIdentity == nil {
		t.Fatal("LinkedEPSBearerIdentity is nil")
	}
	if decoded.SpareHalfOctet == nil {
		t.Fatal("SpareHalfOctet is nil")
	}
	if decoded.EPSQoS == nil {
		t.Fatal("EPSQoS is nil")
	}
	if decoded.TrafficFlowTemplate == nil {
		t.Fatal("TrafficFlowTemplate is nil")
	}
	if decoded.LLCSAPI == nil {
		t.Error("LLCSAPI is nil")
	} else if len(decoded.LLCSAPI.Payload) != 1 || decoded.LLCSAPI.Payload[0] != 5 {
		t.Errorf("LLCSAPI = %v, want [5]", decoded.LLCSAPI.Payload)
	}

	if decoded.EsmHeader.BearerIdentity != 0x01 {
		t.Errorf("BearerIdentity = %#02x, want 0x01", decoded.EsmHeader.BearerIdentity)
	}
	if decoded.EsmHeader.ProcedureTransactionIdentity != 0x01 {
		t.Errorf("PTI = %#02x, want 0x01", decoded.EsmHeader.ProcedureTransactionIdentity)
	}
}

func TestActivateDedicatedEPSBearerContextRequestMultiplePTI(t *testing.T) {
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
					MsgType:                      msg.ESMActivateDedicatedEPSBearerContextRequest,
				}

				linkedBEID := ie.NewPositionalLinkedEPSBearerIdentity(bearerID)
				spare := ie.NewPositionalSpareHalfOctet(0x00)
				esmQoS := ie.NewPositionalEPSQualityOfService(qosPayload)
				tft := ie.NewTrafficFlowTemplate([]byte{0xAA, 0xBB})

				m := msg.NewActivateDedicatedEPSBearerContextRequest(
					linkedBEID, spare, esmQoS, tft,
				)
				m.EsmHeader = h

				buf := make([]byte, m.MarshalLen())
				if err := m.MarshalTo(buf); err != nil {
					t.Fatalf("MarshalTo failed: %v", err)
				}

				var decoded msg.ActivateDedicatedEPSBearerContextRequest
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

func TestActivateDedicatedEPSBearerContextRequestPackedHalfOctetsRoundTrip(t *testing.T) {
	t.Parallel()

	// Test multiple bearer-id / spare-half-octet combinations to verify the
	// packed type-1 pair round-trips correctly.
	testCases := []struct {
		name  string
		bid   uint8
		spare uint8
	}{
		{name: "bid1_spare0", bid: 0x01, spare: 0x00},
		{name: "bid3_spare0", bid: 0x03, spare: 0x00},
		{name: "bid5_spareF", bid: 0x05, spare: 0x0F},
		{name: "bid7_spare3", bid: 0x07, spare: 0x03},
	}

	qosFields := &ie.EPSQualityOfServiceFields{QCI: 9}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h := msg.EsmHeader{
				BearerIdentity:               0x00, // bearer ID is not in EsmHeader for dedicated
				ProtocolDiscriminator:        byte(msg.PdEsm),
				ProcedureTransactionIdentity: 0x01,
				MsgType:                      msg.ESMActivateDedicatedEPSBearerContextRequest,
			}

			linkedBEID := ie.NewPositionalLinkedEPSBearerIdentity(tc.bid)
			spare := ie.NewPositionalSpareHalfOctet(tc.spare)
			esmQoS := ie.NewPositionalEPSQualityOfService(qosPayload)
			tft := ie.NewTrafficFlowTemplate([]byte{0x12, 0x34})

			m := msg.NewActivateDedicatedEPSBearerContextRequest(
				linkedBEID, spare, esmQoS, tft,
			)
			m.EsmHeader = h

			buf := make([]byte, m.MarshalLen())
			if err := m.MarshalTo(buf); err != nil {
				t.Fatalf("MarshalTo failed: %v", err)
			}

			var decoded msg.ActivateDedicatedEPSBearerContextRequest
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

			if decoded.LinkedEPSBearerIdentity == nil {
				t.Fatal("LinkedEPSBearerIdentity is nil")
			}
			if len(decoded.LinkedEPSBearerIdentity.Payload) != 1 ||
				decoded.LinkedEPSBearerIdentity.Payload[0] != tc.bid {
				t.Errorf("LinkedEPSBearerIdentity = %v, want [%d]", decoded.LinkedEPSBearerIdentity.Payload, tc.bid)
			}

			if decoded.SpareHalfOctet == nil {
				t.Fatal("SpareHalfOctet is nil")
			}
			if len(decoded.SpareHalfOctet.Payload) != 1 ||
				decoded.SpareHalfOctet.Payload[0] != tc.spare {
				t.Errorf("SpareHalfOctet = %v, want [%d]", decoded.SpareHalfOctet.Payload, tc.spare)
			}
		})
	}
}

func TestActivateDedicatedEPSBearerContextRequestEmptyOptional(t *testing.T) {
	t.Parallel()

	qosFields := &ie.EPSQualityOfServiceFields{QCI: 5}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}

	h := msg.EsmHeader{
		BearerIdentity:               0x02,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x03,
		MsgType:                      msg.ESMActivateDedicatedEPSBearerContextRequest,
	}

	linkedBEID := ie.NewPositionalLinkedEPSBearerIdentity(0x02)
	spare := ie.NewPositionalSpareHalfOctet(0x00)
	esmQoS := ie.NewPositionalEPSQualityOfService(qosPayload)
	tft := ie.NewTrafficFlowTemplate([]byte{0x56, 0x78})

	m := msg.NewActivateDedicatedEPSBearerContextRequest(
		linkedBEID, spare, esmQoS, tft,
	)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ActivateDedicatedEPSBearerContextRequest
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
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs in empty optional")
	}
}

func TestActivateDedicatedEPSBearerContextRequestErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.ActivateDedicatedEPSBearerContextRequest

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

	local := &msg.ActivateDedicatedEPSBearerContextRequest{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.ActivateDedicatedEPSBearerContextRequest

	// Truncated ESM header (2 bytes, need 3).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Invalid PD in 3-octet header.
	err = decoded.UnmarshalBinary([]byte{0xAA, 0x01, 0xC5})
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid PD = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}

	// Valid header but no payload — need at least half-octet (1 byte).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, 0xC5})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated no payload = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Valid header + 1 byte half-octet, no QoS (need length byte).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, 0xC5, 0x00})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated no QoS = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Valid header + half-octet + QoS (LV, need at least length+1).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, 0xC5, 0x01})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated no QoS payload = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Valid header + half-octet + QoS LV (2 bytes), no TFT.
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, 0xC5, 0x00, 0x02, 0xAA})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated no TFT = %v, want %T", err, io.ErrUnexpectedEOF)
	}
}

func TestActivateDedicatedEPSBearerContextRequestParseHelper(t *testing.T) {
	t.Parallel()

	qosFields := &ie.EPSQualityOfServiceFields{QCI: 9}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}

	h := msg.EsmHeader{
		BearerIdentity:               0x03,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x04,
		MsgType:                      msg.ESMActivateDedicatedEPSBearerContextRequest,
	}

	linkedBEID := ie.NewPositionalLinkedEPSBearerIdentity(0x03)
	spare := ie.NewPositionalSpareHalfOctet(0x00)
	esmQoS := ie.NewPositionalEPSQualityOfService(qosPayload)
	tft := ie.NewTrafficFlowTemplate([]byte{0xCC, 0xDD})

	m := msg.NewActivateDedicatedEPSBearerContextRequest(
		linkedBEID, spare, esmQoS, tft,
	)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	parsed, err := msg.ParseActivateDedicatedEPSBearerContextRequest(buf)
	if err != nil {
		t.Fatalf("ParseActivateDedicatedEPSBearerContextRequest failed: %v", err)
	}

	if parsed == nil {
		t.Fatal("ParseActivateDedicatedEPSBearerContextRequest returned nil")
	}

	if parsed.LinkedEPSBearerIdentity == nil {
		t.Fatal("LinkedEPSBearerIdentity is nil")
	}
	if parsed.SpareHalfOctet == nil {
		t.Fatal("SpareHalfOctet is nil")
	}
	if parsed.EPSQoS == nil {
		t.Fatal("EPSQoS is nil")
	}
	if parsed.TrafficFlowTemplate == nil {
		t.Fatal("TrafficFlowTemplate is nil")
	}

	if parsed.EsmHeader.BearerIdentity != 0x03 {
		t.Errorf("BearerIdentity = %#02x, want 0x03", parsed.EsmHeader.BearerIdentity)
	}
	if parsed.EsmHeader.ProcedureTransactionIdentity != 0x04 {
		t.Errorf("PTI = %#02x, want 0x04", parsed.EsmHeader.ProcedureTransactionIdentity)
	}
	if parsed.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if parsed.MessageType() != msg.ESMActivateDedicatedEPSBearerContextRequest {
		t.Errorf("MessageType = %#02x, want %#02x", parsed.MessageType(), msg.ESMActivateDedicatedEPSBearerContextRequest)
	}
}

func TestActivateDedicatedEPSBearerContextRequestWithDuplicates(t *testing.T) {
	t.Parallel()

	qosFields := &ie.EPSQualityOfServiceFields{QCI: 9}
	qosPayload, err := qosFields.Marshal()
	if err != nil {
		t.Fatalf("QoS marshal failed: %v", err)
	}

	linkedBEID := ie.NewPositionalLinkedEPSBearerIdentity(0x01)
	spare := ie.NewPositionalSpareHalfOctet(0x00)
	esmQoS := ie.NewPositionalEPSQualityOfService(qosPayload)
	tft := ie.NewTrafficFlowTemplate([]byte{0x12, 0x34})

	// Two LLCSAPI IEs.
	llc1 := ie.NewLLCServiceAccessPointIdentifier(0x01)
	llc2 := ie.NewLLCServiceAccessPointIdentifier(0x02)

	m := msg.NewActivateDedicatedEPSBearerContextRequest(
		linkedBEID, spare, esmQoS, tft, llc1, llc2,
	)

	if m.LLCSAPI == nil {
		t.Fatal("first LLCSAPI was not stored")
	}

	var foundDup bool
	for _, i := range m.IEs {
		if i.Kind == ie.LLCServiceAccessPointIdentifier && len(i.Payload) == 1 && i.Payload[0] == 2 {
			foundDup = true
		}
	}

	if !foundDup {
		t.Error("expected duplicate LLCSAPI in IEs")
	}
}
