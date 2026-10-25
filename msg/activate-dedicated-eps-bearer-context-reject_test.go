package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestActivateDedicatedEPSBearerContextRejectRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x01,
		MsgType:                      msg.ESMActivateDedicatedEPSBearerContextReject,
	}

	esmCause := ie.NewESMCause(ie.RequestRejectedUnspecifiedSmCause)

	// Tagged optional: LLC SAPI.
	llcSAPI := ie.NewRaw(0x32, ie.IEFormatFour, []byte{5})

	m := msg.NewActivateDedicatedEPSBearerContextReject(esmCause, llcSAPI)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	if len(buf) != m.MarshalLen() {
		t.Fatalf("marshal length = %d, want %d", len(buf), m.MarshalLen())
	}

	var decoded msg.ActivateDedicatedEPSBearerContextReject
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

	if got := decoded.MessageType(); got != msg.ESMActivateDedicatedEPSBearerContextReject {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.ESMActivateDedicatedEPSBearerContextReject)
	}
	if got := decoded.MessageTypeName(); got != "Activate dedicated EPS bearer context reject" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Activate dedicated EPS bearer context reject")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}

	if decoded.ESMCause == nil {
		t.Fatal("ESMCause is nil")
	}
	if len(decoded.ESMCause.Payload) != 1 || decoded.ESMCause.Payload[0] != byte(ie.RequestRejectedUnspecifiedSmCause) {
		t.Errorf("ESMCause = %v, want [%d]", decoded.ESMCause.Payload, ie.RequestRejectedUnspecifiedSmCause)
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

func TestActivateDedicatedEPSBearerContextRejectMultiplePTI(t *testing.T) {
	t.Parallel()

	for _, bearerID := range []uint8{0x01, 0x02, 0x05, 0x07} {
		for _, pti := range []uint8{0x01, 0x02, 0x0F} {
			t.Run("", func(t *testing.T) {
				t.Parallel()

				h := msg.EsmHeader{
					BearerIdentity:               bearerID,
					ProtocolDiscriminator:        byte(msg.PdEsm),
					ProcedureTransactionIdentity: pti,
					MsgType:                      msg.ESMActivateDedicatedEPSBearerContextReject,
				}

				esmCause := ie.NewESMCause(ie.RequestRejectedUnspecifiedSmCause)

				m := msg.NewActivateDedicatedEPSBearerContextReject(esmCause)
				m.EsmHeader = h

				buf := make([]byte, m.MarshalLen())
				if err := m.MarshalTo(buf); err != nil {
					t.Fatalf("MarshalTo failed: %v", err)
				}

				var decoded msg.ActivateDedicatedEPSBearerContextReject
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

func TestActivateDedicatedEPSBearerContextRejectEmptyOptional(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x03,
		MsgType:                      msg.ESMActivateDedicatedEPSBearerContextReject,
	}

	esmCause := ie.NewESMCause(ie.RequestRejectedUnspecifiedSmCause)

	m := msg.NewActivateDedicatedEPSBearerContextReject(esmCause)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ActivateDedicatedEPSBearerContextReject
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

func TestActivateDedicatedEPSBearerContextRejectESMCauses(t *testing.T) {
	t.Parallel()

	for _, cause := range []ie.SmCause{
		ie.RequestRejectedUnspecifiedSmCause,
		ie.EpsQosNotAcceptedSmCause,
		ie.NetworkFailureSmCause,
		ie.InvalidEpsBearerIdentitySmCause,
	} {
		t.Run("", func(t *testing.T) {
			t.Parallel()

			h := msg.EsmHeader{
				BearerIdentity:               0x01,
				ProtocolDiscriminator:        byte(msg.PdEsm),
				ProcedureTransactionIdentity: 0x02,
				MsgType:                      msg.ESMActivateDedicatedEPSBearerContextReject,
			}

			esmCause := ie.NewESMCause(cause)

			m := msg.NewActivateDedicatedEPSBearerContextReject(esmCause)
			m.EsmHeader = h

			buf := make([]byte, m.MarshalLen())
			if err := m.MarshalTo(buf); err != nil {
				t.Fatalf("MarshalTo failed: %v", err)
			}

			var decoded msg.ActivateDedicatedEPSBearerContextReject
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
				t.Fatal("ESMCause is nil")
			}
			if decoded.ESMCause.Payload[0] != byte(cause) {
				t.Errorf("ESMCause = %d, want %d", decoded.ESMCause.Payload[0], cause)
			}
		})
	}
}

func TestActivateDedicatedEPSBearerContextRejectErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.ActivateDedicatedEPSBearerContextReject

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

	local := &msg.ActivateDedicatedEPSBearerContextReject{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.ActivateDedicatedEPSBearerContextReject

	// Truncated ESM header (2 bytes, need 3).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Invalid PD in 3-octet header.
	err = decoded.UnmarshalBinary([]byte{0xAA, 0x01, 0xC7})
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid PD = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}

	// Valid header but no payload — need ESM cause (1 byte).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, 0xC7})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated no ESM cause = %v, want %T", err, io.ErrUnexpectedEOF)
	}
}

func TestActivateDedicatedEPSBearerContextRejectParseHelper(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x02,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x04,
		MsgType:                      msg.ESMActivateDedicatedEPSBearerContextReject,
	}

	esmCause := ie.NewESMCause(ie.NetworkFailureSmCause)

	m := msg.NewActivateDedicatedEPSBearerContextReject(esmCause)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	parsed, err := msg.ParseActivateDedicatedEPSBearerContextReject(buf)
	if err != nil {
		t.Fatalf("ParseActivateDedicatedEPSBearerContextReject failed: %v", err)
	}

	if parsed == nil {
		t.Fatal("ParseActivateDedicatedEPSBearerContextReject returned nil")
	}

	if parsed.ESMCause == nil {
		t.Fatal("ESMCause is nil")
	}
	if parsed.ESMCause.Payload[0] != byte(ie.NetworkFailureSmCause) {
		t.Errorf("ESMCause = %d, want %d", parsed.ESMCause.Payload[0], ie.NetworkFailureSmCause)
	}

	if parsed.EsmHeader.BearerIdentity != 0x02 {
		t.Errorf("BearerIdentity = %#02x, want 0x02", parsed.EsmHeader.BearerIdentity)
	}
	if parsed.EsmHeader.ProcedureTransactionIdentity != 0x04 {
		t.Errorf("PTI = %#02x, want 0x04", parsed.EsmHeader.ProcedureTransactionIdentity)
	}
	if parsed.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if parsed.MessageType() != msg.ESMActivateDedicatedEPSBearerContextReject {
		t.Errorf("MessageType = %#02x, want %#02x", parsed.MessageType(), msg.ESMActivateDedicatedEPSBearerContextReject)
	}
}

func TestActivateDedicatedEPSBearerContextRejectWithLLCSAPI(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x03,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x05,
		MsgType:                      msg.ESMActivateDedicatedEPSBearerContextReject,
	}

	esmCause := ie.NewESMCause(ie.ServiceOptionNotSupportedSmCause)
	llcSAPI := ie.NewRaw(0x32, ie.IEFormatFour, []byte{4})

	m := msg.NewActivateDedicatedEPSBearerContextReject(esmCause, llcSAPI)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.ActivateDedicatedEPSBearerContextReject
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
		t.Fatal("ESMCause is nil")
	}
	if decoded.ESMCause.Payload[0] != byte(ie.ServiceOptionNotSupportedSmCause) {
		t.Errorf("ESMCause = %d, want %d", decoded.ESMCause.Payload[0], ie.ServiceOptionNotSupportedSmCause)
	}

	if decoded.LLCSAPI == nil {
		t.Error("LLCSAPI is nil")
	} else if len(decoded.LLCSAPI.Payload) != 1 || decoded.LLCSAPI.Payload[0] != 4 {
		t.Errorf("LLCSAPI = %v, want [4]", decoded.LLCSAPI.Payload)
	}
}
