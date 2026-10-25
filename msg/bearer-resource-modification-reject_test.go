package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestBearerResourceModificationRejectRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x01,
		MsgType:                      msg.ESMBearerResourceModificationReject,
	}

	esmCause := ie.NewPositionalESMCause(ie.NetworkFailureSmCause)
	m := msg.NewBearerResourceModificationReject(esmCause)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.BearerResourceModificationReject
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

	if got := decoded.MessageType(); got != msg.ESMBearerResourceModificationReject {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.ESMBearerResourceModificationReject)
	}
	if got := decoded.MessageTypeName(); got != "Bearer resource modification reject" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Bearer resource modification reject")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if decoded.ESMCause == nil {
		t.Fatal("ESMCause is nil")
	}

	if decoded.EsmHeader.BearerIdentity != 0x01 {
		t.Errorf("BearerIdentity = %#02x, want 0x01", decoded.EsmHeader.BearerIdentity)
	}
}

func TestBearerResourceModificationRejectWithLLCSAPI(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x02,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x02,
		MsgType:                      msg.ESMBearerResourceModificationReject,
	}

	esmCause := ie.NewPositionalESMCause(ie.ServiceOptionNotSupportedSmCause)
	llcSAPI := ie.NewRaw(0x32, ie.IEFormatFour, []byte{5})

	m := msg.NewBearerResourceModificationReject(esmCause, llcSAPI)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.BearerResourceModificationReject
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
	if decoded.ESMCause == nil {
		t.Fatal("ESMCause is nil after round-trip")
	}
}

func TestBearerResourceModificationRejectWithExtendedQoS(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x03,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x03,
		MsgType:                      msg.ESMBearerResourceModificationReject,
	}

	esmCause := ie.NewPositionalESMCause(ie.EpsQosNotAcceptedSmCause)
	extQoSFields := &ie.ExtendedQualityOfServiceFields{QCI: 7}
	extQoSIE, err := ie.NewExtendedQualityOfServiceWithFields(extQoSFields)
	if err != nil {
		t.Fatalf("Extended QoS marshal failed: %v", err)
	}

	m := msg.NewBearerResourceModificationReject(esmCause, extQoSIE)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.BearerResourceModificationReject
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

	if decoded.ExtendedQoS == nil {
		t.Fatal("ExtendedQoS is nil after round-trip")
	}
}

func TestBearerResourceModificationRejectMultiplePTI(t *testing.T) {
	t.Parallel()

	for _, bearerID := range []uint8{0x01, 0x02, 0x05, 0x07} {
		for _, pti := range []uint8{0x01, 0x02, 0x0F} {
			t.Run("", func(t *testing.T) {
				t.Parallel()

				h := msg.EsmHeader{
					BearerIdentity:               bearerID,
					ProtocolDiscriminator:        byte(msg.PdEsm),
					ProcedureTransactionIdentity: pti,
					MsgType:                      msg.ESMBearerResourceModificationReject,
				}

				esmCause := ie.NewPositionalESMCause(ie.NetworkFailureSmCause)
				m := msg.NewBearerResourceModificationReject(esmCause)
				m.EsmHeader = h

				buf := make([]byte, m.MarshalLen())
				if err := m.MarshalTo(buf); err != nil {
					t.Fatalf("MarshalTo failed: %v", err)
				}

				var decoded msg.BearerResourceModificationReject
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

func TestBearerResourceModificationRejectDuplicateIEs(t *testing.T) {
	t.Parallel()

	esmCause := ie.NewPositionalESMCause(ie.NetworkFailureSmCause)
	c1 := ie.NewLLCServiceAccessPointIdentifier(3)
	c2 := ie.NewLLCServiceAccessPointIdentifier(5)

	m := msg.NewBearerResourceModificationReject(esmCause, c1, c2)

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

func TestBearerResourceModificationRejectErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.BearerResourceModificationReject

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

	local := &msg.BearerResourceModificationReject{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.BearerResourceModificationReject

	// Truncated ESM header (2 bytes, need 3).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Invalid PD in 3-octet header.
	err = decoded.UnmarshalBinary([]byte{0xAA, 0x01, 0xD7})
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid PD = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}

	// Valid header but no payload.
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, 0xD7})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated no cause = %v, want %T", err, io.ErrUnexpectedEOF)
	}
}

func TestBearerResourceModificationRejectParseHelper(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x05,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x06,
		MsgType:                      msg.ESMBearerResourceModificationReject,
	}

	esmCause := ie.NewPositionalESMCause(ie.NetworkFailureSmCause)
	m := msg.NewBearerResourceModificationReject(esmCause)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	parsed, err := msg.ParseBearerResourceModificationReject(buf)
	if err != nil {
		t.Fatalf("ParseBearerResourceModificationReject failed: %v", err)
	}

	if parsed == nil {
		t.Fatal("ParseBearerResourceModificationReject returned nil")
	}

	if parsed.ESMCause == nil {
		t.Fatal("ESMCause is nil")
	}

	if parsed.IsRequest() {
		t.Error("IsRequest = true, want false")
	}

	if parsed.EsmHeader.BearerIdentity != 0x05 {
		t.Errorf("BearerIdentity = %#02x, want 0x05", parsed.EsmHeader.BearerIdentity)
	}
}
