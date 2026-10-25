package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestRemoteUEReportResponseRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x01,
		MsgType:                      msg.ESMRemoteUEReportResponse,
	}

	m := msg.NewRemoteUEReportResponse()
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.RemoteUEReportResponse
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

	if decoded.MessageType() != msg.ESMRemoteUEReportResponse {
		t.Errorf("MessageType = %#02x, want %#02x", decoded.MessageType(), msg.ESMRemoteUEReportResponse)
	}
	if decoded.MessageTypeName() != "Remote UE report response" {
		t.Errorf("MessageTypeName = %q, want %q", decoded.MessageTypeName(), "Remote UE report response")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
}

func TestRemoteUEReportResponseWithUEremoteID(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x02,
		MsgType:                      msg.ESMRemoteUEReportResponse,
	}

	// UE remote ID: tag 0x7F, LV-E (tag + length BE2 + payload)
	ueID := ie.NewRaw(0x7f, ie.IEFormatSix, []byte{0xDE, 0xAD, 0xBE, 0xEF})

	m := msg.NewRemoteUEReportResponse(ueID)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.RemoteUEReportResponse
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.UERemoteID == nil {
		t.Fatal("UERemoteID is nil")
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

	if decoded.UERemoteID.MarshalLen() != ueID.MarshalLen() {
		t.Errorf("UE ID length = %d, want %d", decoded.UERemoteID.MarshalLen(), ueID.MarshalLen())
	}
}

func TestRemoteUEReportResponseWithAllOptionalIEs(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x03,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x04,
		MsgType:                      msg.ESMRemoteUEReportResponse,
	}

	ueID := ie.NewRaw(0x7f, ie.IEFormatSix, []byte{0x11, 0x22})
	relayInfo := ie.NewRaw(0x65, ie.IEFormatFour, []byte{0xAB})

	m := msg.NewRemoteUEReportResponse(ueID, relayInfo)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.RemoteUEReportResponse
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.UERemoteID == nil {
		t.Fatal("UERemoteID is nil")
	}
	if decoded.RelaySelectionInfo == nil {
		t.Fatal("RelaySelectionInfo is nil")
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

func TestRemoteUEReportResponseWithUnknownExtraIE(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x01,
		MsgType:                      msg.ESMRemoteUEReportResponse,
	}

	m := msg.NewRemoteUEReportResponse()
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.RemoteUEReportResponse
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if len(decoded.IEs) != 0 {
		t.Fatalf("IEs length = %d, want 0", len(decoded.IEs))
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

func TestRemoteUEReportResponseDuplicateIEs(t *testing.T) {
	t.Parallel()

	ue1 := ie.NewRaw(0x7f, ie.IEFormatSix, []byte{0x01})
	ue2 := ie.NewRaw(0x7f, ie.IEFormatSix, []byte{0x02})
	relay1 := ie.NewRaw(0x65, ie.IEFormatFour, []byte{0x10})
	relay2 := ie.NewRaw(0x65, ie.IEFormatFour, []byte{0x20})

	m := msg.NewRemoteUEReportResponse(ue1, ue2, relay1, relay2)

	if m.UERemoteID == nil {
		t.Fatal("first UERemoteID was not stored")
	}
	if m.RelaySelectionInfo == nil {
		t.Fatal("first RelaySelectionInfo was not stored")
	}

	if len(m.IEs) != 2 {
		t.Fatalf("IEs length = %d, want 2", len(m.IEs))
	}
}

func TestRemoteUEReportResponseMultiplePTI(t *testing.T) {
	t.Parallel()

	for _, bearerID := range []uint8{0x01, 0x02, 0x05, 0x07} {
		for _, pti := range []uint8{0x01, 0x02, 0x0F} {
			t.Run("", func(t *testing.T) {
				t.Parallel()

				h := msg.EsmHeader{
					BearerIdentity:               bearerID,
					ProtocolDiscriminator:        byte(msg.PdEsm),
					ProcedureTransactionIdentity: pti,
					MsgType:                      msg.ESMRemoteUEReportResponse,
				}

				m := msg.NewRemoteUEReportResponse()
				m.EsmHeader = h

				buf := make([]byte, m.MarshalLen())
				if err := m.MarshalTo(buf); err != nil {
					t.Fatalf("MarshalTo failed: %v", err)
				}

				var decoded msg.RemoteUEReportResponse
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

func TestRemoteUEReportResponseErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.RemoteUEReportResponse

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

	local := &msg.RemoteUEReportResponse{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.RemoteUEReportResponse

	// Truncated ESM header (2 bytes, need 3).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Valid header but empty payload — no optional IEs.
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, msg.ESMRemoteUEReportResponse})
	if err != nil {
		t.Fatalf("empty payload = %v, want no error", err)
	}

	// Truncated LV-E: tag 0x7F + only length high byte.
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, msg.ESMRemoteUEReportResponse, 0x7f, 0x00})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated LV-E = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Truncated LV: tag 0x65 + length byte indicates more data.
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, msg.ESMRemoteUEReportResponse, 0x65, 0x03})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated LV = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Valid header + UE remote ID (LV-E: tag + length BE2=2 + payload) but missing relay selection info length.
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, msg.ESMRemoteUEReportResponse, 0x7f, 0x00, 0x02, 0x11, 0x22, 0x65})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated after UE ID = %v, want %T", err, io.ErrUnexpectedEOF)
	}
}

func TestRemoteUEReportResponseParseHelper(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x06,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x07,
		MsgType:                      msg.ESMRemoteUEReportResponse,
	}

	ueID := ie.NewRaw(0x7f, ie.IEFormatSix, []byte{0xAA, 0xBB, 0xCC})
	relayInfo := ie.NewRaw(0x65, ie.IEFormatFour, []byte{0x33})

	m := msg.NewRemoteUEReportResponse(ueID, relayInfo)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	parsed, err := msg.ParseRemoteUEReportResponse(buf)
	if err != nil {
		t.Fatalf("ParseRemoteUEReportResponse failed: %v", err)
	}

	if parsed == nil {
		t.Fatal("ParseRemoteUEReportResponse returned nil")
	}
	if parsed.UERemoteID == nil {
		t.Fatal("UERemoteID is nil")
	}
	if parsed.RelaySelectionInfo == nil {
		t.Fatal("RelaySelectionInfo is nil")
	}
	if parsed.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if parsed.MessageType() != msg.ESMRemoteUEReportResponse {
		t.Errorf("MessageType = %#02x, want %#02x", parsed.MessageType(), msg.ESMRemoteUEReportResponse)
	}
	if parsed.EsmHeader.BearerIdentity != 0x06 {
		t.Errorf("BearerIdentity = %#02x, want 0x06", parsed.EsmHeader.BearerIdentity)
	}
}
