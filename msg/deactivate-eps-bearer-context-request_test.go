package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestDeactivateEPSBearerContextRequestRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x01,
		MsgType:                      msg.ESMDeactivateEPSBearerContextRequest,
	}

	esmCause := ie.NewPositionalESMCause(0x34) // Requested bearer context modified

	m := msg.NewDeactivateEPSBearerContextRequest(esmCause)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	if len(buf) != m.MarshalLen() {
		t.Fatalf("marshal length = %d, want %d", len(buf), m.MarshalLen())
	}

	var decoded msg.DeactivateEPSBearerContextRequest
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

	if got := decoded.MessageType(); got != msg.ESMDeactivateEPSBearerContextRequest {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.ESMDeactivateEPSBearerContextRequest)
	}
	if got := decoded.MessageTypeName(); got != "Deactivate EPS bearer context request" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Deactivate EPS bearer context request")
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

func TestDeactivateEPSBearerContextRequestMultiplePTI(t *testing.T) {
	t.Parallel()

	for _, bearerID := range []uint8{0x01, 0x02, 0x05, 0x07} {
		for _, pti := range []uint8{0x01, 0x02, 0x0F} {
			t.Run("", func(t *testing.T) {
				t.Parallel()

				h := msg.EsmHeader{
					BearerIdentity:               bearerID,
					ProtocolDiscriminator:        byte(msg.PdEsm),
					ProcedureTransactionIdentity: pti,
					MsgType:                      msg.ESMDeactivateEPSBearerContextRequest,
				}

				esmCause := ie.NewPositionalESMCause(0x34)

				m := msg.NewDeactivateEPSBearerContextRequest(esmCause)
				m.EsmHeader = h

				buf := make([]byte, m.MarshalLen())
				if err := m.MarshalTo(buf); err != nil {
					t.Fatalf("MarshalTo failed: %v", err)
				}

				var decoded msg.DeactivateEPSBearerContextRequest
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

func TestDeactivateEPSBearerContextRequestEmptyOptional(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x03,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x03,
		MsgType:                      msg.ESMDeactivateEPSBearerContextRequest,
	}

	esmCause := ie.NewPositionalESMCause(0x34)

	m := msg.NewDeactivateEPSBearerContextRequest(esmCause)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.DeactivateEPSBearerContextRequest
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
	if decoded.ExtendedProtocolConfigOpts != nil {
		t.Error("unexpected ExtendedProtocolConfigOpts in empty optional")
	}
	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs in empty optional")
	}
}

func TestDeactivateEPSBearerContextRequestErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.DeactivateEPSBearerContextRequest

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

	local := &msg.DeactivateEPSBearerContextRequest{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.DeactivateEPSBearerContextRequest

	// Truncated ESM header (2 bytes, need 3).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Invalid PD in 3-octet header.
	err = decoded.UnmarshalBinary([]byte{0xAA, 0x01, 0xCD})
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid PD = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}

	// Valid header but no payload — need at least ESM cause (1 byte).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, 0xCD})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated no cause = %v, want %T", err, io.ErrUnexpectedEOF)
	}
}

func TestDeactivateEPSBearerContextRequestParseHelper(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x05,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x06,
		MsgType:                      msg.ESMDeactivateEPSBearerContextRequest,
	}

	esmCause := ie.NewPositionalESMCause(0x34)

	m := msg.NewDeactivateEPSBearerContextRequest(esmCause)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	parsed, err := msg.ParseDeactivateEPSBearerContextRequest(buf)
	if err != nil {
		t.Fatalf("ParseDeactivateEPSBearerContextRequest failed: %v", err)
	}

	if parsed == nil {
		t.Fatal("ParseDeactivateEPSBearerContextRequest returned nil")
	}

	if parsed.ESMCause == nil {
		t.Fatal("ESMCause is nil")
	}

	if parsed.EsmHeader.BearerIdentity != 0x05 {
		t.Errorf("BearerIdentity = %#02x, want 0x05", parsed.EsmHeader.BearerIdentity)
	}
	if parsed.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if parsed.MessageType() != msg.ESMDeactivateEPSBearerContextRequest {
		t.Errorf("MessageType = %#02x, want %#02x", parsed.MessageType(), msg.ESMDeactivateEPSBearerContextRequest)
	}
}

func TestDeactivateEPSBearerContextRequestWithLLCSAPIAndEPCO(t *testing.T) {
	t.Parallel()

	llcSAPI := ie.NewRaw(0x32, ie.IEFormatFour, []byte{5})
	extPCO := ie.NewExtendedProtocolConfigurationOptions([]byte{0x50, 0x02, 0x21, 0x80})
	esmCause := ie.NewPositionalESMCause(0x34)

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x01,
		MsgType:                      msg.ESMDeactivateEPSBearerContextRequest,
	}

	m := msg.NewDeactivateEPSBearerContextRequest(esmCause, llcSAPI, extPCO)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.DeactivateEPSBearerContextRequest
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
	if decoded.ExtendedProtocolConfigOpts == nil {
		t.Fatal("ExtendedProtocolConfigOpts is nil after round-trip")
	}
	if len(decoded.ExtendedProtocolConfigOpts.Payload) != 4 {
		t.Errorf("ExtendedProtocolConfigOpts payload len = %d, want 4", len(decoded.ExtendedProtocolConfigOpts.Payload))
	}
}

func TestDeactivateEPSBearerContextRequestDuplicateIEs(t *testing.T) {
	t.Parallel()

	llc1 := ie.NewRaw(0x32, ie.IEFormatFour, []byte{5})
	llc2 := ie.NewRaw(0x32, ie.IEFormatFour, []byte{7})
	esmCause := ie.NewPositionalESMCause(0x34)

	m := msg.NewDeactivateEPSBearerContextRequest(esmCause, llc1, llc2)

	if m.LLCSAPI == nil {
		t.Fatal("first LLCSAPI was not stored")
	}

	var foundDup bool
	for _, i := range m.IEs {
		if i.Kind == ie.LLCServiceAccessPointIdentifier && len(i.Payload) == 1 && i.Payload[0] == 7 {
			foundDup = true
		}
	}

	if !foundDup {
		t.Error("expected duplicate LLCSAPI in IEs")
	}
}
