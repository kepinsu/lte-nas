package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestAttachAcceptRoundTrip(t *testing.T) {
	t.Parallel()

	esmc := ie.NewPositionalESMMessageContainer([]byte{0xC1, 0x02, 0xAA, 0xBB})

	m := msg.NewAttachAccept(
		ie.NewPositionalEPSAttachResult(0x01),
		ie.NewPositionalSpareHalfOctet(0x0F),
		ie.NewPositionalGPRSTimer(0x05),
		ie.NewPositionalTrackingAreaIdentityList([]byte{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC}),
		esmc,
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMAttachAccept,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.AttachAccept
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

	if got := decoded.MessageType(); got != msg.EMMAttachAccept {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMAttachAccept)
	}
	if got := decoded.MessageTypeName(); got != "Attach accept" {
		t.Errorf("MessageTypeName = %q, want %q", got, "Attach accept")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if decoded.EmmHeader.SecurityHeaderType != msg.SecurityHeaderTypePlain {
		t.Errorf("SHT = %#02x, want %#02x", decoded.EmmHeader.SecurityHeaderType, msg.SecurityHeaderTypePlain)
	}
}

func TestAttachAcceptWithOptionalTags(t *testing.T) {
	t.Parallel()

	guti := ie.NewMobileIdentity([]byte{0x23, 0x05, 0x12, 0x34, 0xAB})
	emc := ie.NewEMMCause(ie.PlmnNotAllowedMmCause)

	m := msg.NewAttachAccept(
		ie.NewPositionalEPSAttachResult(0x01),
		ie.NewPositionalSpareHalfOctet(0x0F),
		ie.NewPositionalGPRSTimer(0x05),
		ie.NewPositionalTrackingAreaIdentityList([]byte{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC}),
		ie.NewPositionalESMMessageContainer([]byte{0xC1, 0x01}),
		guti, emc,
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMAttachAccept,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.AttachAccept
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

	if decoded.GUTI == nil {
		t.Error("GUTI is nil")
	}
	if decoded.EMMCause == nil {
		t.Error("EMMCause is nil")
	}
}

func TestAttachAcceptProtectedHeader(t *testing.T) {
	t.Parallel()

	m := msg.NewAttachAccept(
		ie.NewPositionalEPSAttachResult(0x01),
		ie.NewPositionalSpareHalfOctet(0x0F),
		ie.NewPositionalGPRSTimer(0x05),
		ie.NewPositionalTrackingAreaIdentityList([]byte{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC}),
		ie.NewPositionalESMMessageContainer([]byte{0xC1, 0x01}),
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
		SequenceNumber:        5,
		MsgType:               msg.EMMAttachAccept,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.AttachAccept
	if err := decoded.UnmarshalBinary(buf); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	if decoded.EmmHeader.SecurityHeaderType != msg.SecurityHeaderTypeIntegrityProtected {
		t.Errorf("SHT = %#02x, want %#02x", decoded.EmmHeader.SecurityHeaderType, msg.SecurityHeaderTypeIntegrityProtected)
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

func TestAttachAcceptErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.AttachAccept

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

	local := &msg.AttachAccept{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.AttachAccept
	err = decoded.UnmarshalBinary([]byte{0x17, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0x42})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated protected content = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	err = decoded.UnmarshalBinary([]byte{0x07, 0x42}) // plain header - no content
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("plain header only = %v, want %T", err, io.ErrUnexpectedEOF)
	}
}

func TestAttachAcceptNilFieldMarshal(t *testing.T) {
	t.Parallel()

	m := msg.NewAttachAccept(
		nil,
		ie.NewPositionalSpareHalfOctet(0x0F),
		ie.NewPositionalGPRSTimer(0x05),
		ie.NewPositionalTrackingAreaIdentityList([]byte{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC}),
		ie.NewPositionalESMMessageContainer([]byte{0xC1, 0x01}),
	)

	m.EmmHeader = msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMAttachAccept,
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.AttachAccept
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
}

// New* constructor edge cases: nil IE skip, duplicate handling, unknown IE default.
func TestAttachAcceptNewConstructorEdgeCases(t *testing.T) {
	t.Parallel()

	result := ie.NewPositionalEPSAttachResult(0x01)
	spare := ie.NewPositionalSpareHalfOctet(0x0F)
	t3412 := ie.NewPositionalGPRSTimer(0x05)
	taiList := ie.NewPositionalTrackingAreaIdentityList([]byte{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC})
	esmc := ie.NewPositionalESMMessageContainer([]byte{0xC1, 0x02, 0xAA})

	unknownIE := ie.NewRaw(0x99, ie.IEFormatFour, []byte{0xAA, 0xBB})

	// Test nil extra IE — the if i == nil { continue } path
	t.Run("nil extra IE skip", func(t *testing.T) {
		t.Parallel()

		m := msg.NewAttachAccept(result, spare, t3412, taiList, esmc, nil, nil)
		if len(m.IEs) != 0 {
			t.Errorf("nil IEs should not be appended, got %d IEs", len(m.IEs))
		}
	})

	// Test duplicate known kind — the else { m.IEs = append(m.IEs, i) } path
	t.Run("duplicate known kind", func(t *testing.T) {
		t.Parallel()

		guti1 := ie.NewMobileIdentity([]byte{0x23, 0x05, 0x12, 0x34, 0xAB})
		guti2 := ie.NewMobileIdentity([]byte{0x23, 0x05, 0xFF, 0xEE, 0xDD})

		m := msg.NewAttachAccept(result, spare, t3412, taiList, esmc, guti1, guti2)

		// First GUTI should be assigned to m.GUTI
		if m.GUTI == nil {
			t.Fatal("first GUTI should be assigned to m.GUTI")
		}

		// Second (duplicate) GUTI should be in m.IEs
		dupFound := false
		for _, item := range m.IEs {
			if item.Kind == ie.MobileIdentity {
				dupFound = true
			}
		}
		if !dupFound {
			t.Errorf("duplicate GUTI should be in m.IEs, got %d IEs: %+v", len(m.IEs), m.IEs)
		}
	})

	// Test unknown IE kind — the default: m.IEs = append(m.IEs, i) path
	t.Run("unknown IE kind", func(t *testing.T) {
		t.Parallel()

		m := msg.NewAttachAccept(result, spare, t3412, taiList, esmc, unknownIE)

		if len(m.IEs) != 1 {
			t.Fatalf("unknown IE should be in m.IEs, got %d IEs", len(m.IEs))
		}

		if m.IEs[0].Kind != 0x99 {
			t.Errorf("IE Kind = %#02x, want %#02x", m.IEs[0].Kind, 0x99)
		}

		if m.IEs[0].Payload[0] != 0xAA {
			t.Errorf("IE Payload = %v, want [0xaa]", m.IEs[0].Payload)
		}
	})

	// Test all optional fields assigned correctly (no duplicates)
	t.Run("all known kinds assigned", func(t *testing.T) {
		t.Parallel()

		guti := ie.NewMobileIdentity([]byte{0x23, 0x05, 0x12, 0x34, 0xAB})
		emc := ie.NewEMMCause(ie.PlmnNotAllowedMmCause)
		t3412Ext := ie.NewT3412ExtendedValue(0x06) // T3412ExtendedValue (tag 0x5E)

		m := msg.NewAttachAccept(result, spare, t3412, taiList, esmc, guti, emc, t3412Ext)

		if m.GUTI == nil {
			t.Error("GUTI should be assigned")
		}
		if m.EMMCause == nil {
			t.Error("EMMCause should be assigned")
		}
		if m.T3412Extended == nil {
			t.Error("T3412Extended should be assigned")
		}

		if len(m.IEs) != 0 {
			t.Errorf("no extras expected, got %d IEs", len(m.IEs))
		}
	})

	// Test mixed: one known + one duplicate known + one unknown
	t.Run("mixed duplicates and unknown", func(t *testing.T) {
		t.Parallel()

		guti1 := ie.NewMobileIdentity([]byte{0x23, 0x05, 0x12, 0x34, 0xAB})
		guti2 := ie.NewMobileIdentity([]byte{0x23, 0x05, 0xFF, 0xEE, 0xDD})
		emc1 := ie.NewEMMCause(ie.PlmnNotAllowedMmCause)
		unknown1 := ie.NewRaw(0x99, ie.IEFormatFour, []byte{0x11})
		unknown2 := ie.NewRaw(0x88, ie.IEFormatOne, []byte{0x22})

		m := msg.NewAttachAccept(result, spare, t3412, taiList, esmc,
			guti1, guti2, emc1, unknown1, unknown2)

		if m.GUTI == nil {
			t.Error("first GUTI should be assigned")
		}
		if m.EMMCause == nil {
			t.Error("first EMM Cause should be assigned")
		}

		// Should have: guti2 (dup), unknown1, unknown2 in IEs
		if len(m.IEs) != 3 {
			t.Errorf("expected 3 extra IEs, got %d", len(m.IEs))
		}

		found := map[byte]int{}
		for _, item := range m.IEs {
			found[byte(item.Kind)]++
		}
		if found[0x23] != 1 {
			t.Errorf("duplicate GUTI (0x23) count = %d, want 1", found[0x23])
		}
		if found[0x99] != 1 {
			t.Errorf("unknown IE 0x99 count = %d, want 1", found[0x99])
		}
		if found[0x88] != 1 {
			t.Errorf("unknown IE 0x88 count = %d, want 1", found[0x88])
		}
	})

	// Test all nil positional fields — exercises MarshalLen short-circuit paths
	t.Run("all nil positional", func(t *testing.T) {
		t.Parallel()

		m := msg.NewAttachAccept(nil, nil, nil, nil, nil)
		m.EmmHeader = msg.EmmHeader{
			SecurityHeaderType:    msg.SecurityHeaderTypePlain,
			ProtocolDiscriminator: byte(msg.PdEmm),
			MsgType:               msg.EMMAttachAccept,
		}

		n := m.MarshalLen()
		if n != 2 {
			t.Errorf("MarshalLen with nil IEs and valid header = %d, want 2 (header only)", n)
		}
	})

	// Test single nil in varargs — exercises continue path specifically
	t.Run("single nil vararg", func(t *testing.T) {
		t.Parallel()

		guti := ie.NewMobileIdentity([]byte{0x23, 0x05, 0x12, 0x34, 0xAB})
		m := msg.NewAttachAccept(result, spare, t3412, taiList, esmc, guti, nil)

		if len(m.IEs) != 0 {
			t.Errorf("nil should be skipped, got %d IEs", len(m.IEs))
		}
		if m.GUTI == nil {
			t.Error("GUTI should be assigned from guti")
		}

		// Add an extra IE after the nil
		m2 := msg.NewAttachAccept(result, spare, t3412, taiList, esmc, nil, unknownIE)
		if len(m2.IEs) != 1 {
			t.Errorf("nil skipped, unknownIE appended: expected 1 IE, got %d", len(m2.IEs))
		}
	})
}
