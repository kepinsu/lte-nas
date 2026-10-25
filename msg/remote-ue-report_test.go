package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestRemoteUEReportRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x01,
		MsgType:                      msg.ESMRemoteUEReport,
	}

	m := msg.NewRemoteUEReport()
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.RemoteUEReport
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

	if decoded.MessageType() != msg.ESMRemoteUEReport {
		t.Errorf("MessageType = %#02x, want %#02x", decoded.MessageType(), msg.ESMRemoteUEReport)
	}
	if decoded.MessageTypeName() != "Remote UE report" {
		t.Errorf("MessageTypeName = %q, want %q", decoded.MessageTypeName(), "Remote UE report")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
}

func TestRemoteUEReportWithRelaySelectionInfo(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x02,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x03,
		MsgType:                      msg.ESMRemoteUEReport,
	}

	// Relay selection info: tag 0x65, LV (tag + length + payload)
	relayInfo := ie.NewRaw(0x65, ie.IEFormatFour, []byte{0x12, 0x34})

	m := msg.NewRemoteUEReport(relayInfo)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.RemoteUEReport
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
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

	if decoded.RelaySelectionInfo.MarshalLen() != relayInfo.MarshalLen() {
		t.Errorf("relay info length = %d, want %d", decoded.RelaySelectionInfo.MarshalLen(), relayInfo.MarshalLen())
	}
}

func TestRemoteUEReportWithAllOptionalIEs(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x04,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x05,
		MsgType:                      msg.ESMRemoteUEReport,
	}

	relayInfo := ie.NewRaw(0x65, ie.IEFormatFour, []byte{0xAB, 0xCD})
	dupRelay := ie.NewRaw(0x65, ie.IEFormatFour, []byte{0xEF})

	m := msg.NewRemoteUEReport(relayInfo, dupRelay)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.RemoteUEReport
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.RelaySelectionInfo == nil {
		t.Fatal("RelaySelectionInfo is nil")
	}

	if len(decoded.IEs) != 1 {
		t.Fatalf("IEs length = %d, want 1", len(decoded.IEs))
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

func TestRemoteUEReportDuplicateIEs(t *testing.T) {
	t.Parallel()

	relay1 := ie.NewRaw(0x65, ie.IEFormatFour, []byte{0x01})
	relay2 := ie.NewRaw(0x65, ie.IEFormatFour, []byte{0x02})

	m := msg.NewRemoteUEReport(relay1, relay2)

	if m.RelaySelectionInfo == nil {
		t.Fatal("first RelaySelectionInfo was not stored")
	}

	if len(m.IEs) != 1 {
		t.Fatalf("IEs length = %d, want 1", len(m.IEs))
	}

	var foundDup bool
	for _, i := range m.IEs {
		if i.Kind == ie.Kind(0x65) && len(i.Payload) >= 1 {
			foundDup = true
		}
	}

	if !foundDup {
		t.Error("expected duplicate RelaySelectionInfo in IEs")
	}
}

func TestRemoteUEReportMultiplePTI(t *testing.T) {
	t.Parallel()

	for _, bearerID := range []uint8{0x01, 0x02, 0x05, 0x07} {
		for _, pti := range []uint8{0x01, 0x02, 0x0F} {
			t.Run("", func(t *testing.T) {
				t.Parallel()

				h := msg.EsmHeader{
					BearerIdentity:               bearerID,
					ProtocolDiscriminator:        byte(msg.PdEsm),
					ProcedureTransactionIdentity: pti,
					MsgType:                      msg.ESMRemoteUEReport,
				}

				m := msg.NewRemoteUEReport()
				m.EsmHeader = h

				buf := make([]byte, m.MarshalLen())
				if err := m.MarshalTo(buf); err != nil {
					t.Fatalf("MarshalTo failed: %v", err)
				}

				var decoded msg.RemoteUEReport
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

func TestRemoteUEReportErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.RemoteUEReport

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

	local := &msg.RemoteUEReport{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.RemoteUEReport

	// Truncated ESM header (2 bytes, need 3).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Valid header but empty payload — no optional IEs.
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, msg.ESMRemoteUEReport})
	if err != nil {
		t.Fatalf("empty payload = %v, want no error", err)
	}

	// Short buffer for IE decoding — valid header + single byte (tag but no length).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, msg.ESMRemoteUEReport, 0x65})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated IE = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Valid IE but length indicates more data.
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, msg.ESMRemoteUEReport, 0x65, 0x03})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated LV = %v, want %T", err, io.ErrUnexpectedEOF)
	}
}

func TestRemoteUEReportParseHelper(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x05,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x06,
		MsgType:                      msg.ESMRemoteUEReport,
	}

	relayInfo := ie.NewRaw(0x65, ie.IEFormatFour, []byte{0xDE, 0xAD})

	m := msg.NewRemoteUEReport(relayInfo)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	parsed, err := msg.ParseRemoteUEReport(buf)
	if err != nil {
		t.Fatalf("ParseRemoteUEReport failed: %v", err)
	}

	if parsed == nil {
		t.Fatal("ParseRemoteUEReport returned nil")
	}
	if parsed.RelaySelectionInfo == nil {
		t.Fatal("RelaySelectionInfo is nil")
	}
	if parsed.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if parsed.MessageType() != msg.ESMRemoteUEReport {
		t.Errorf("MessageType = %#02x, want %#02x", parsed.MessageType(), msg.ESMRemoteUEReport)
	}
	if parsed.EsmHeader.BearerIdentity != 0x05 {
		t.Errorf("BearerIdentity = %#02x, want 0x05", parsed.EsmHeader.BearerIdentity)
	}
}
