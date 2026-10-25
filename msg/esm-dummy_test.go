package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestESMDummyRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x01,
		MsgType:                      msg.ESMDummy,
	}

	m := msg.NewESMDummy()
	m.EsmHeader = h

	// ESM DUMMY always marshals to exactly 3 octets.
	if m.MarshalLen() != 3 {
		t.Fatalf("MarshalLen = %d, want 3", m.MarshalLen())
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	// Verify exact 3-octet output.
	if len(buf) != 3 {
		t.Fatalf("buf length = %d, want 3", len(buf))
	}

	var decoded msg.DummyMessage
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.MessageType() != msg.ESMDummy {
		t.Errorf("MessageType = %#02x, want %#02x", decoded.MessageType(), msg.ESMDummy)
	}
	if decoded.MessageTypeName() != "ESM dummy message" {
		t.Errorf("MessageTypeName = %q, want %q", decoded.MessageTypeName(), "ESM dummy message")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}

	if decoded.EsmHeader.BearerIdentity != 0x01 {
		t.Errorf("BearerIdentity = %#02x, want 0x01", decoded.EsmHeader.BearerIdentity)
	}
	if decoded.EsmHeader.ProcedureTransactionIdentity != 0x01 {
		t.Errorf("PTI = %#02x, want 0x01", decoded.EsmHeader.ProcedureTransactionIdentity)
	}
}

func TestESMDummyRejectsTrailingBytes(t *testing.T) {
	t.Parallel()

	var decoded msg.DummyMessage

	// 4 bytes: valid header + 1 trailing byte.
	errMsg := []byte{0xA2, 0x01, msg.ESMDummy, 0xFF}
	err := decoded.UnmarshalBinary(errMsg)
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("trailing byte = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// 5 bytes: valid header + 2 trailing bytes.
	errMsg = []byte{0xA2, 0x01, msg.ESMDummy, 0xFF, 0xAA}
	if err := decoded.UnmarshalBinary(errMsg); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("trailing bytes = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Long payload.
	errMsg = make([]byte, 20)
	copy(errMsg, []byte{0xA2, 0x01, msg.ESMDummy})
	if err := decoded.UnmarshalBinary(errMsg); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("long payload = %v, want %T", err, io.ErrUnexpectedEOF)
	}
}

func TestESMDummyNoIEsInConstructor(t *testing.T) {
	t.Parallel()

	m := msg.NewESMDummy()
	// NewESMDummy returns an unconfigured header; set the PD so MarshalTo succeeds.
	m.EsmHeader.ProtocolDiscriminator = byte(msg.PdEsm)

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	if len(buf) != 3 {
		t.Fatalf("buf length = %d, want 3", len(buf))
	}
}

func TestESMDummyDifferentHeaderValues(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		bearerID uint8
		pti      uint8
	}{
		{0x01, 0x01},
		{0x07, 0x0F},
		{0x05, 0x02},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			t.Parallel()

			m := msg.NewESMDummy()
			m.EsmHeader = msg.EsmHeader{
				BearerIdentity:               tc.bearerID,
				ProtocolDiscriminator:        byte(msg.PdEsm),
				ProcedureTransactionIdentity: tc.pti,
				MsgType:                      msg.ESMDummy,
			}

			buf := make([]byte, m.MarshalLen())
			if err := m.MarshalTo(buf); err != nil {
				t.Fatalf("MarshalTo failed: %v", err)
			}

			if len(buf) != 3 {
				t.Fatalf("buf length = %d, want 3", len(buf))
			}

			var decoded msg.DummyMessage
			if err := decoded.UnmarshalBinary(buf); err != nil {
				t.Fatalf("UnmarshalBinary failed: %v", err)
			}

			if decoded.EsmHeader.BearerIdentity != tc.bearerID {
				t.Errorf("BearerIdentity = %#02x, want %#02x", decoded.EsmHeader.BearerIdentity, tc.bearerID)
			}
			if decoded.EsmHeader.ProcedureTransactionIdentity != tc.pti {
				t.Errorf("PTI = %#02x, want %#02x", decoded.EsmHeader.ProcedureTransactionIdentity, tc.pti)
			}
		})
	}
}

func TestESMDummyErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.DummyMessage

	if got := nilMsg.MarshalLen(); got != 0 {
		t.Errorf("nil MarshalLen = %d, want 0", got)
	}

	err := nilMsg.MarshalTo(nil)
	if err == nil {
		t.Fatal("nil MarshalTo = nil, want error")
	}
	if !errors.Is(err, ie.ErrNilIE) {
		t.Errorf("nil MarshalTo = %v, want wrapped ie.ErrNilIE", err)
	}

	local := &msg.DummyMessage{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.DummyMessage

	// Truncated ESM header (2 bytes, need 3).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Invalid PD.
	err = decoded.UnmarshalBinary([]byte{0xAA, 0x01, msg.ESMDummy})
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid PD = %v, want wrapped ie.ErrInvalidValue", err)
	}

	// No trailing bytes should pass.
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, msg.ESMDummy})
	if err != nil {
		t.Fatalf("exact header = %v, want nil", err)
	}
}

func TestESMDummyParseHelper(t *testing.T) {
	t.Parallel()

	m := msg.NewESMDummy()
	m.EsmHeader = msg.EsmHeader{
		BearerIdentity:               0x04,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x05,
		MsgType:                      msg.ESMDummy,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	parsed, err := msg.ParseESMDummy(buf)
	if err != nil {
		t.Fatalf("ParseESMDummy failed: %v", err)
	}

	if parsed == nil {
		t.Fatal("ParseESMDummy returned nil")
	}
	if parsed.MessageType() != msg.ESMDummy {
		t.Errorf("MessageType = %#02x, want %#02x", parsed.MessageType(), msg.ESMDummy)
	}
	if parsed.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if parsed.EsmHeader.BearerIdentity != 0x04 {
		t.Errorf("BearerIdentity = %#02x, want 0x04", parsed.EsmHeader.BearerIdentity)
	}
}
