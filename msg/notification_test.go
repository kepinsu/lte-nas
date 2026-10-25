package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestNotificationRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x01,
		MsgType:                      msg.ESMNotification,
	}

	nind := ie.NewPositionalNotificationIndicator(0x01)

	m := msg.NewNotification(nind)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.Notification
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

	if decoded.MessageType() != msg.ESMNotification {
		t.Errorf("MessageType = %#02x, want %#02x", decoded.MessageType(), msg.ESMNotification)
	}
	if decoded.MessageTypeName() != "Notification" {
		t.Errorf("MessageTypeName = %q, want %q", decoded.MessageTypeName(), "Notification")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if decoded.NotificationIndicator == nil {
		t.Fatal("NotificationIndicator is nil")
	}

	// Verify wire format: header(3) + LV indicator length byte(1) + value(1) = 5
	if len(buf) != 5 {
		t.Errorf("buf length = %d, want 5", len(buf))
	}

	// Check header: BearerIdentity<<4 | PdEsm should give first byte.
	wantFirst := byte(0x12) // BearerIdentity=1 (0x01<<4) | PdEsm=0x2
	if buf[0] != wantFirst {
		t.Errorf("header[0] = %#02x, want %#02x", buf[0], wantFirst)
	}
	if buf[1] != 0x01 || buf[2] != msg.ESMNotification {
		t.Errorf("header = %#02x %#02x %#02x, want %#02x %#02x %#02x",
			buf[0], buf[1], buf[2], wantFirst, 0x01, msg.ESMNotification)
	}
	if buf[3] != 0x01 || buf[4] != 0x01 {
		t.Errorf("notification indicator = %#02x %#02x, want 0x01 0x01", buf[3], buf[4])
	}
}

func TestNotificationReservedValues(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x02,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x02,
		MsgType:                      msg.ESMNotification,
	}

	for _, value := range []uint8{0x60, 0x7F} {
		t.Run("", func(t *testing.T) {
			t.Parallel()

			nind := ie.NewPositionalNotificationIndicator(value)

			m := msg.NewNotification(nind)
			m.EsmHeader = h

			buf := make([]byte, m.MarshalLen())
			if err := m.MarshalTo(buf); err != nil {
				t.Fatalf("MarshalTo failed: %v", err)
			}

			var decoded msg.Notification
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
		})
	}
}

func TestNotificationMultiplePTI(t *testing.T) {
	t.Parallel()

	for _, bearerID := range []uint8{0x01, 0x02, 0x05, 0x07} {
		for _, pti := range []uint8{0x01, 0x02, 0x0F} {
			t.Run("", func(t *testing.T) {
				t.Parallel()

				h := msg.EsmHeader{
					BearerIdentity:               bearerID,
					ProtocolDiscriminator:        byte(msg.PdEsm),
					ProcedureTransactionIdentity: pti,
					MsgType:                      msg.ESMNotification,
				}

				nind := ie.NewPositionalNotificationIndicator(0x01)

				m := msg.NewNotification(nind)
				m.EsmHeader = h

				buf := make([]byte, m.MarshalLen())
				if err := m.MarshalTo(buf); err != nil {
					t.Fatalf("MarshalTo failed: %v", err)
				}

				var decoded msg.Notification
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

func TestNotificationWithExtraIEs(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x01,
		MsgType:                      msg.ESMNotification,
	}

	nind := ie.NewPositionalNotificationIndicator(0x01)
	extraIE := ie.NewRaw(0x2D, ie.IEFormatThree, []byte{0xFF})

	m := msg.NewNotification(nind, extraIE)
	m.EsmHeader = h

	// New* populates m.IEs even though UnmarshalBinary rejects trailing bytes.
	if len(m.IEs) != 1 {
		t.Fatalf("m.IEs length = %d, want 1", len(m.IEs))
	}
	if m.IEs[0].Kind != 0x2D {
		t.Errorf("IE[0] Kind = 0x%02x, want 0x2D", m.IEs[0].Kind)
	}

	// Marshal to verify the extra-IE loops in MarshalTo and MarshalLen execute.
	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	// UnmarshalBinary rejects trailing bytes, so skip round-trip.
	// We verified New* and MarshalTo work correctly above.
}

func TestNotificationErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.Notification

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

	local := &msg.Notification{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.Notification

	// Truncated ESM header (2 bytes, need 3).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Valid header but no payload.
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, msg.ESMNotification})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("no payload = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Only 1 byte of payload — need LV length(1) + value(1) = 2 bytes.
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, msg.ESMNotification, 0x01})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("partial LV = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Valid LV but trailing byte after indicator.
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, msg.ESMNotification, 0x01, 0x01, 0xFF})
	if !errors.Is(err, ie.ErrUnexpectedIE) {
		t.Fatalf("trailing byte = %v, want wrapped %T", err, ie.ErrUnexpectedIE)
	}

	// Invalid PD.
	err = decoded.UnmarshalBinary([]byte{0xAA, 0x01, msg.ESMNotification})
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid PD = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}
}

func TestNotificationParseHelper(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x03,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x04,
		MsgType:                      msg.ESMNotification,
	}

	nind := ie.NewPositionalNotificationIndicator(0x01)

	m := msg.NewNotification(nind)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	parsed, err := msg.ParseNotification(buf)
	if err != nil {
		t.Fatalf("ParseNotification failed: %v", err)
	}

	if parsed == nil {
		t.Fatal("ParseNotification returned nil")
	}
	if parsed.NotificationIndicator == nil {
		t.Fatal("NotificationIndicator is nil")
	}
	if parsed.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if parsed.MessageType() != msg.ESMNotification {
		t.Errorf("MessageType = %#02x, want %#02x", parsed.MessageType(), msg.ESMNotification)
	}
	if parsed.EsmHeader.BearerIdentity != 0x03 {
		t.Errorf("BearerIdentity = %#02x, want 0x03", parsed.EsmHeader.BearerIdentity)
	}
}
