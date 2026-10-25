package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestEsmDataTransportRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x01,
		MsgType:                      msg.ESMDataTransport,
	}

	m := msg.NewEsmDataTransport()
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.EsmDataTransport
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

	if decoded.MessageType() != msg.ESMDataTransport {
		t.Errorf("MessageType = %#02x, want %#02x", decoded.MessageType(), msg.ESMDataTransport)
	}
	if decoded.MessageTypeName() != "ESM data transport" {
		t.Errorf("MessageTypeName = %q, want %q", decoded.MessageTypeName(), "ESM data transport")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
}

func TestEsmDataTransportWithUserDataContainer(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x02,
		MsgType:                      msg.ESMDataTransport,
	}

	// User data container: tag 0x78, LV-E (tag + length BE2 + payload)
	userData := ie.NewRaw(0x78, ie.IEFormatSix, []byte{0x0A, 0x0B, 0x0C})

	m := msg.NewEsmDataTransport(userData)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.EsmDataTransport
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.UserDataContainer == nil {
		t.Fatal("UserDataContainer is nil")
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

	if decoded.UserDataContainer.MarshalLen() != userData.MarshalLen() {
		t.Errorf("user data length = %d, want %d", decoded.UserDataContainer.MarshalLen(), userData.MarshalLen())
	}
}

func TestEsmDataTransportWithAllOptionalIEs(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x03,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x04,
		MsgType:                      msg.ESMDataTransport,
	}

	userData := ie.NewRaw(0x78, ie.IEFormatSix, []byte{0x11, 0x22})
	cpOnly := ie.NewControlPlaneOnlyIndication(0x01)

	m := msg.NewEsmDataTransport(userData, cpOnly)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.EsmDataTransport
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.UserDataContainer == nil {
		t.Fatal("UserDataContainer is nil")
	}
	if decoded.ControlPlaneOnlyIndic == nil {
		t.Fatal("ControlPlaneOnlyIndic is nil")
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

func TestEsmDataTransportWithUnknownExtraIE(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x01,
		MsgType:                      msg.ESMDataTransport,
	}

	m := msg.NewEsmDataTransport()
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.EsmDataTransport
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

func TestEsmDataTransportDuplicateIEs(t *testing.T) {
	t.Parallel()

	ud1 := ie.NewRaw(0x78, ie.IEFormatSix, []byte{0x01})
	ud2 := ie.NewRaw(0x78, ie.IEFormatSix, []byte{0x02})
	cp1 := ie.NewControlPlaneOnlyIndication(0x01)
	cp2 := ie.NewControlPlaneOnlyIndication(0x02)

	m := msg.NewEsmDataTransport(ud1, ud2, cp1, cp2)

	if m.UserDataContainer == nil {
		t.Fatal("first UserDataContainer was not stored")
	}
	if m.ControlPlaneOnlyIndic == nil {
		t.Fatal("first ControlPlaneOnlyIndic was not stored")
	}

	if len(m.IEs) != 2 {
		t.Fatalf("IEs length = %d, want 2", len(m.IEs))
	}
}

func TestEsmDataTransportMultiplePTI(t *testing.T) {
	t.Parallel()

	for _, bearerID := range []uint8{0x01, 0x02, 0x05, 0x07} {
		for _, pti := range []uint8{0x01, 0x02, 0x0F} {
			t.Run("", func(t *testing.T) {
				t.Parallel()

				h := msg.EsmHeader{
					BearerIdentity:               bearerID,
					ProtocolDiscriminator:        byte(msg.PdEsm),
					ProcedureTransactionIdentity: pti,
					MsgType:                      msg.ESMDataTransport,
				}

				m := msg.NewEsmDataTransport()
				m.EsmHeader = h

				buf := make([]byte, m.MarshalLen())
				if err := m.MarshalTo(buf); err != nil {
					t.Fatalf("MarshalTo failed: %v", err)
				}

				var decoded msg.EsmDataTransport
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

func TestEsmDataTransportErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.EsmDataTransport

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

	local := &msg.EsmDataTransport{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.EsmDataTransport

	// Truncated ESM header (2 bytes, need 3).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Valid header but empty payload — must have at least one tagged IE.
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, msg.ESMDataTransport})
	if err != nil {
		t.Fatalf("empty payload = %v, want no error", err)
	}

	// Truncated LV-E: tag 0x78 + only length high byte.
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, msg.ESMDataTransport, 0x78, 0x00})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated LV-E = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Complete user data container only (LV-E: tag + length BE2=2 + payload).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, msg.ESMDataTransport, 0x78, 0x00, 0x02, 0x11, 0x22})
	if err != nil {
		t.Fatalf("complete with user data only = %v, want no error", err)
	}

	// Complete message with both user data container and control plane only indication.
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, msg.ESMDataTransport, 0x78, 0x00, 0x02, 0x11, 0x22, 0x90})
	if err != nil {
		t.Fatalf("complete with both IEs = %v, want no error", err)
	}

	// Complete with user data container only — valid message, CPIO is optional.
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, msg.ESMDataTransport, 0x78, 0x00, 0x02, 0x11, 0x22})
	if err != nil {
		t.Fatalf("user data only = %v, want no error", err)
	}

	// Invalid ESM header: non-ESM protocol discriminator.
	err = decoded.UnmarshalBinary([]byte{0x07, 0x41, 0xEB})
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid PD = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}
}

func TestEsmDataTransportParseHelper(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x05,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x06,
		MsgType:                      msg.ESMDataTransport,
	}

	userData := ie.NewRaw(0x78, ie.IEFormatSix, []byte{0xDE, 0xAD, 0xBE, 0xEF})
	cpOnly := ie.NewControlPlaneOnlyIndication(0x01)

	m := msg.NewEsmDataTransport(userData, cpOnly)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	parsed, err := msg.ParseEsmDataTransport(buf)
	if err != nil {
		t.Fatalf("ParseEsmDataTransport failed: %v", err)
	}

	if parsed == nil {
		t.Fatal("ParseEsmDataTransport returned nil")
	}
	if parsed.UserDataContainer == nil {
		t.Fatal("UserDataContainer is nil")
	}
	if parsed.ControlPlaneOnlyIndic == nil {
		t.Fatal("ControlPlaneOnlyIndic is nil")
	}
	if parsed.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if parsed.MessageType() != msg.ESMDataTransport {
		t.Errorf("MessageType = %#02x, want %#02x", parsed.MessageType(), msg.ESMDataTransport)
	}
	if parsed.EsmHeader.BearerIdentity != 0x05 {
		t.Errorf("BearerIdentity = %#02x, want 0x05", parsed.EsmHeader.BearerIdentity)
	}
}
