package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestESMStatusRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x01,
		MsgType:                      msg.ESMStatus,
	}

	cause := ie.NewESMCause(ie.NetworkFailureSmCause)

	m := msg.NewESMStatusMsg(cause)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ESMStatusMsg
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

	if decoded.MessageType() != msg.ESMStatus {
		t.Errorf("MessageType = %#02x, want %#02x", decoded.MessageType(), msg.ESMStatus)
	}
	if decoded.MessageTypeName() != "ESM status" {
		t.Errorf("MessageTypeName = %q, want %q", decoded.MessageTypeName(), "ESM status")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if decoded.ESMCause == nil {
		t.Fatal("ESMCause is nil")
	}

	if cause, err := decoded.ESMCause.ESMCause(); err != nil || cause != ie.NetworkFailureSmCause {
		t.Errorf("ESMCause = %v, want %v", cause, ie.NetworkFailureSmCause)
	}
	if decoded.EsmHeader.BearerIdentity != 0x01 {
		t.Errorf("BearerIdentity = %#02x, want 0x01", decoded.EsmHeader.BearerIdentity)
	}
}

func TestESMStatusAllCauseValues(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		cause ie.SmCause
	}{
		{"insufficient resources", ie.InsufficientResourcesSmCause},
		{"missing APN", ie.MissingOrUnknownApnSmCause},
		{"user auth failed", ie.UserAuthenticationOrAuthorizationFailedSmCause},
		{"PTI in use", ie.PtiAlreadyInUseSmCause},
		{"reactivation requested", ie.ReactivationRequestedSmCause},
		{"invalid bearer ID", ie.InvalidEpsBearerIdentitySmCause},
		{"protocol error", ie.ProtocolErrorUnspecifiedSmCause},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h := msg.EsmHeader{
				BearerIdentity:               0x01,
				ProtocolDiscriminator:        byte(msg.PdEsm),
				ProcedureTransactionIdentity: 0x01,
				MsgType:                      msg.ESMStatus,
			}

			cause := ie.NewESMCause(tc.cause)

			m := msg.NewESMStatusMsg(cause)
			m.EsmHeader = h

			buf := make([]byte, m.MarshalLen())
			if err := m.MarshalTo(buf); err != nil {
				t.Fatalf("MarshalTo failed: %v", err)
			}

			var decoded msg.ESMStatusMsg
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

			if decoded.EsmHeader.BearerIdentity != 0x01 {
				t.Errorf("BearerIdentity = %#02x, want 0x01", decoded.EsmHeader.BearerIdentity)
			}
		})
	}
}

func TestESMStatusMultiplePTI(t *testing.T) {
	t.Parallel()

	for _, bearerID := range []uint8{0x01, 0x02, 0x05, 0x07} {
		for _, pti := range []uint8{0x01, 0x02, 0x0F} {
			t.Run("", func(t *testing.T) {
				t.Parallel()

				h := msg.EsmHeader{
					BearerIdentity:               bearerID,
					ProtocolDiscriminator:        byte(msg.PdEsm),
					ProcedureTransactionIdentity: pti,
					MsgType:                      msg.ESMStatus,
				}

				cause := ie.NewESMCause(ie.NetworkFailureSmCause)

				m := msg.NewESMStatusMsg(cause)
				m.EsmHeader = h

				buf := make([]byte, m.MarshalLen())
				if err := m.MarshalTo(buf); err != nil {
					t.Fatalf("MarshalTo failed: %v", err)
				}

				var decoded msg.ESMStatusMsg
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

func TestESMStatusTrailingBytesRejected(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x01,
		MsgType:                      msg.ESMStatus,
	}

	cause := ie.NewESMCause(ie.NetworkFailureSmCause)

	m := msg.NewESMStatusMsg(cause)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	// Append trailing byte.
	bufWithTrailing := append(buf, 0xFF)

	var decoded msg.ESMStatusMsg
	err := decoded.UnmarshalBinary(bufWithTrailing)
	if !errors.Is(err, ie.ErrUnexpectedIE) {
		t.Fatalf("trailing byte = %v, want wrapped %T", err, ie.ErrUnexpectedIE)
	}

	// Multiple trailing bytes.
	bufWithTrailing = append(buf, 0xFF, 0xAA)
	if err := decoded.UnmarshalBinary(bufWithTrailing); !errors.Is(err, ie.ErrUnexpectedIE) {
		t.Fatalf("trailing bytes = %v, want wrapped %T", err, ie.ErrUnexpectedIE)
	}
}

func TestESMStatusErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.ESMStatusMsg

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

	local := &msg.ESMStatusMsg{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.ESMStatusMsg

	// Truncated ESM header (2 bytes, need 3).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Valid header but no payload.
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, msg.ESMStatus})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("no payload = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Invalid PD.
	err = decoded.UnmarshalBinary([]byte{0xAA, 0x01, msg.ESMStatus})
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid PD = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}

	// Trailing byte after cause.
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, msg.ESMStatus, 0x1A, 0xFF})
	if !errors.Is(err, ie.ErrUnexpectedIE) {
		t.Fatalf("trailing after cause = %v, want wrapped %T", err, ie.ErrUnexpectedIE)
	}
}

func TestESMStatusParseHelper(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x03,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x04,
		MsgType:                      msg.ESMStatus,
	}

	cause := ie.NewESMCause(ie.InvalidPtiValueSmCause)

	m := msg.NewESMStatusMsg(cause)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	parsed, err := msg.ParseESMStatusMsg(buf)
	if err != nil {
		t.Fatalf("ParseESMStatusMsg failed: %v", err)
	}

	if parsed == nil {
		t.Fatal("ParseESMStatusMsg returned nil")
	}
	if parsed.ESMCause == nil {
		t.Fatal("ESMCause is nil")
	}
	if parsed.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if parsed.MessageType() != msg.ESMStatus {
		t.Errorf("MessageType = %#02x, want %#02x", parsed.MessageType(), msg.ESMStatus)
	}
	if parsed.EsmHeader.BearerIdentity != 0x03 {
		t.Errorf("BearerIdentity = %#02x, want 0x03", parsed.EsmHeader.BearerIdentity)
	}
}

func TestESMStatusMsgMarshalLenExtraIEs(t *testing.T) {
	t.Parallel()

	collisionIE := &ie.IE{Kind: 0x32, Format: ie.IEFormatFour, Payload: []byte{0xAA}}

	m := &msg.ESMStatusMsg{
		EsmHeader: msg.EsmHeader{BearerIdentity: 0x01, ProtocolDiscriminator: byte(msg.PdEsm), ProcedureTransactionIdentity: 0x01, MsgType: msg.ESMStatus},
		ESMCause:  ie.NewESMCause(0x02),
		IEs:       []*ie.IE{collisionIE},
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo with collision IE failed: %v", err)
	}

	// Collision TLV short-buffer error path for the IEs loop.
	collisionMsg := &msg.ESMStatusMsg{
		EsmHeader: msg.EsmHeader{BearerIdentity: 0x01, ProtocolDiscriminator: byte(msg.PdEsm), ProcedureTransactionIdentity: 0x01, MsgType: msg.ESMStatus},
		ESMCause:  ie.NewESMCause(0x02),
		IEs:       []*ie.IE{collisionIE},
	}
	shortBuf := make([]byte, 1)
	if err := collisionMsg.MarshalTo(shortBuf); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("expected ErrShortBuffer for collision IE error path, got %v", err)
	}
}
