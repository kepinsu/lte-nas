package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestEsmInfoRequestRoundTrip(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x01,
		MsgType:                      msg.ESMInformationRequest,
	}

	// UE Network Capability: 6 bytes of capability flags.
	ueNcap := ie.NewPositionalUENetworkCapability([]byte{0x42, 0x88, 0x00, 0x00, 0x00, 0x00})

	// ESM Message Container wrapping an activation request.
	esmContainer := ie.NewPositionalESMMessageContainer([]byte{0xA2, 0x01, 0xC1, 0x42})

	m := msg.NewEsmInfoRequest(ueNcap, esmContainer)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.EsmInfoRequest
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

	if decoded.MessageType() != msg.ESMInformationRequest {
		t.Errorf("MessageType = %#02x, want %#02x", decoded.MessageType(), msg.ESMInformationRequest)
	}
	if decoded.MessageTypeName() != "ESM information request" {
		t.Errorf("MessageTypeName = %q, want %q", decoded.MessageTypeName(), "ESM information request")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if decoded.UENetworkCapability == nil {
		t.Fatal("UENetworkCapability is nil")
	}
	if decoded.ESMMessageContainer == nil {
		t.Fatal("ESMMessageContainer is nil")
	}
	if decoded.EsmHeader.BearerIdentity != 0x01 {
		t.Errorf("BearerIdentity = %#02x, want 0x01", decoded.EsmHeader.BearerIdentity)
	}
}

func TestEsmInfoRequestWithExtendedPCO(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x02,
		MsgType:                      msg.ESMInformationRequest,
	}

	ueNcap := ie.NewPositionalUENetworkCapability([]byte{0x42, 0x88})
	esmContainer := ie.NewPositionalESMMessageContainer([]byte{0xA2, 0x01, 0xC1})
	extendedPCO := ie.NewExtendedProtocolConfigurationOptions([]byte{0x03, 0x50, 0x16})

	m := msg.NewEsmInfoRequest(ueNcap, esmContainer, extendedPCO)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.EsmInfoRequest
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.ExtendedProtocolOpts == nil {
		t.Fatal("ExtendedProtocolOpts is nil")
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

func TestEsmInfoRequestWithNetworkFeatureDown(t *testing.T) {
	t.Parallel()

	// Network Feature Downlink (tag 0xC3) is an unknown extension IE not in the
	// ie catalogue. We verify it does NOT cause a parsing error when encountered
	// in the optional part. Full round-trip marshalling with unknown tags is skipped
	// because ie.MarshalTo rejects high-bit tags (>= 0x80) for type-4 IEs.

	h := msg.EsmHeader{
		BearerIdentity:               0x03,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x04,
		MsgType:                      msg.ESMInformationRequest,
	}

	ueNcap := ie.NewPositionalUENetworkCapability([]byte{0x42, 0x88})
	esmContainer := ie.NewPositionalESMMessageContainer([]byte{0xA2, 0x01})

	m := msg.NewEsmInfoRequest(ueNcap, esmContainer)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.EsmInfoRequest
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.EsmHeader.BearerIdentity != 0x03 {
		t.Errorf("BearerIdentity = %#02x, want 0x03", decoded.EsmHeader.BearerIdentity)
	}
	if decoded.UENetworkCapability == nil {
		t.Fatal("UENetworkCapability is nil")
	}
	if decoded.ESMMessageContainer == nil {
		t.Fatal("ESMMessageContainer is nil")
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

func TestEsmInfoRequestWithAllOptionalIEs(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x01,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x01,
		MsgType:                      msg.ESMInformationRequest,
	}

	ueNcap := ie.NewPositionalUENetworkCapability([]byte{0x42, 0x88})
	esmContainer := ie.NewPositionalESMMessageContainer([]byte{0xA2, 0x01, 0xC1})
	extendedPCO := ie.NewExtendedProtocolConfigurationOptions([]byte{0x03, 0x50})

	m := msg.NewEsmInfoRequest(ueNcap, esmContainer, extendedPCO)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.EsmInfoRequest
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.ExtendedProtocolOpts == nil {
		t.Fatal("ExtendedProtocolOpts is nil")
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

func TestEsmInfoRequestDuplicateIEs(t *testing.T) {
	t.Parallel()

	c1 := ie.NewExtendedProtocolConfigurationOptions([]byte{0x01})
	c2 := ie.NewExtendedProtocolConfigurationOptions([]byte{0x02})

	ueNcap := ie.NewPositionalUENetworkCapability([]byte{0x42, 0x88})
	esmContainer := ie.NewPositionalESMMessageContainer([]byte{0xA2, 0x01})

	m := msg.NewEsmInfoRequest(ueNcap, esmContainer, c1, c2)

	if m.ExtendedProtocolOpts == nil {
		t.Fatal("first ExtendedProtocolOpts was not stored")
	}

	var foundDup bool
	for _, i := range m.IEs {
		if i.Kind == ie.ExtendedProtocolConfigurationOptions && len(i.Payload) >= 1 {
			foundDup = true
		}
	}

	if !foundDup {
		t.Error("expected duplicate ExtendedProtocolOpts in IEs")
	}
}

func TestEsmInfoRequestMultiplePTI(t *testing.T) {
	t.Parallel()

	for _, bearerID := range []uint8{0x01, 0x02, 0x05, 0x07} {
		for _, pti := range []uint8{0x01, 0x02, 0x0F} {
			t.Run("", func(t *testing.T) {
				t.Parallel()

				h := msg.EsmHeader{
					BearerIdentity:               bearerID,
					ProtocolDiscriminator:        byte(msg.PdEsm),
					ProcedureTransactionIdentity: pti,
					MsgType:                      msg.ESMInformationRequest,
				}

				ueNcap := ie.NewPositionalUENetworkCapability([]byte{0x42, 0x88})
				esmContainer := ie.NewPositionalESMMessageContainer([]byte{0xA2, 0x01})

				m := msg.NewEsmInfoRequest(ueNcap, esmContainer)
				m.EsmHeader = h

				buf := make([]byte, m.MarshalLen())
				if err := m.MarshalTo(buf); err != nil {
					t.Fatalf("MarshalTo failed: %v", err)
				}

				var decoded msg.EsmInfoRequest
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

func TestEsmInfoRequestErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.EsmInfoRequest

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

	local := &msg.EsmInfoRequest{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.EsmInfoRequest

	// Truncated ESM header (2 bytes, need 3).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Valid header but empty payload — need UE Network Capability (at least 2 bytes).
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, msg.ESMInformationRequest})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("no payload = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Valid header + UE Ncap but no ESM container.
	err = decoded.UnmarshalBinary([]byte{0xA2, 0x01, msg.ESMInformationRequest, 0x06, 0x42, 0x88})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("no ESM container = %v, want %T", err, io.ErrUnexpectedEOF)
	}
}

func TestEsmInfoRequestParseHelper(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:               0x05,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x06,
		MsgType:                      msg.ESMInformationRequest,
	}

	ueNcap := ie.NewPositionalUENetworkCapability([]byte{0x42, 0x88})
	esmContainer := ie.NewPositionalESMMessageContainer([]byte{0xA2, 0x01, 0xC1})

	m := msg.NewEsmInfoRequest(ueNcap, esmContainer)
	m.EsmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	parsed, err := msg.ParseEsmInfoRequest(buf)
	if err != nil {
		t.Fatalf("ParseEsmInfoRequest failed: %v", err)
	}

	if parsed == nil {
		t.Fatal("ParseEsmInfoRequest returned nil")
	}
	if parsed.UENetworkCapability == nil {
		t.Fatal("UENetworkCapability is nil")
	}
	if parsed.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if parsed.MessageType() != msg.ESMInformationRequest {
		t.Errorf("MessageType = %#02x, want %#02x", parsed.MessageType(), msg.ESMInformationRequest)
	}
	if parsed.EsmHeader.BearerIdentity != 0x05 {
		t.Errorf("BearerIdentity = %#02x, want 0x05", parsed.EsmHeader.BearerIdentity)
	}
}

func TestEsmInfoRequestMarshalLenExtraIEs(t *testing.T) {
	t.Parallel()

	collisionIE := &ie.IE{Kind: 0x32, Format: ie.IEFormatFour, Payload: []byte{0xAA}}

	ueCap := ie.NewUENetworkCapability([]byte{0x01, 0x02})
	esmc := ie.NewESMMessageContainer([]byte{0xC1, 0x02})

	m := &msg.EsmInfoRequest{
		EsmHeader: msg.EsmHeader{
			BearerIdentity:               0x01,
			ProtocolDiscriminator:        byte(msg.PdEsm),
			ProcedureTransactionIdentity: 0x01,
			MsgType:                      msg.ESMInformationRequest,
		},
		UENetworkCapability:   ueCap,
		ESMMessageContainer:   esmc,
		ExtendedProtocolOpts:  ie.NewExtendedProtocolConfigurationOptions([]byte{0x01}),
		IEs:                   []*ie.IE{collisionIE},
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo with collision IE failed: %v", err)
	}

	var decoded msg.EsmInfoRequest
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

	// Collision TLV short-buffer error path.
	collisionMsg := &msg.EsmInfoRequest{
		EsmHeader:             msg.EsmHeader{BearerIdentity: 0x01, ProtocolDiscriminator: byte(msg.PdEsm), ProcedureTransactionIdentity: 0x01, MsgType: msg.ESMInformationRequest},
		UENetworkCapability:   ueCap,
		ESMMessageContainer:   esmc,
		ExtendedProtocolOpts:  ie.NewExtendedProtocolConfigurationOptions([]byte{0x01}),
		IEs:                   []*ie.IE{collisionIE},
	}
	shortBuf := make([]byte, 1)
	if err := collisionMsg.MarshalTo(shortBuf); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("expected ErrShortBuffer for collision IE error path, got %v", err)
	}
}

func TestEsmInfoRequestNewDefaultBranch(t *testing.T) {
	t.Parallel()

	// 0x56 (NonceMME) is not in the EsmInfoRequest switch.
	nonceIE := ie.NewNonceMME([]byte{0xAA})
	m := msg.NewEsmInfoRequest(
		ie.NewUENetworkCapability([]byte{0x01, 0x02}),
		ie.NewESMMessageContainer([]byte{0xC1, 0x02}),
		nonceIE,
	)
	if m.UENetworkCapability == nil || m.ESMMessageContainer == nil {
		t.Fatal("positional IEs should be set")
	}
	if len(m.IEs) == 0 {
		t.Fatal("NonceMME should go to IEs field (default branch)")
	}
	if m.IEs[0].Kind != 0x56 {
		t.Fatalf("IEs[0].Kind = %#02x, want 0x56", m.IEs[0].Kind)
	}
}
