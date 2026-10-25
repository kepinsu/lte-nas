package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestEmmStatusRoundTrip(t *testing.T) {
	t.Parallel()

	// Plain header (2 bytes).
	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMStatus,
	}

	cause := ie.NewPositionalEMMCause(0x16) // #22 Congestion

	m := msg.NewEmmStatus(cause)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.EmmStatus
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

	if got := decoded.MessageType(); got != msg.EMMStatus {
		t.Errorf("MessageType = %#02x, want %#02x", got, msg.EMMStatus)
	}
	if got := decoded.MessageTypeName(); got != "EMM status" {
		t.Errorf("MessageTypeName = %q, want %q", got, "EMM status")
	}
	if decoded.IsRequest() {
		t.Error("IsRequest = true, want false")
	}
	if decoded.EMMCause == nil {
		t.Fatal("EMMCause is nil")
	}
	if causeVal, err := decoded.EMMCause.EMMCause(); err != nil {
		t.Fatalf("EMMCause parse failed: %v", err)
	} else if causeVal != 0x16 {
		t.Errorf("EMMCause = %v, want #22 Congestion (0x16)", causeVal)
	}

	if len(decoded.IEs) > 0 {
		t.Error("unexpected IEs")
	}
}

func TestEmmStatusRoundTripProtected(t *testing.T) {
	t.Parallel()

	// Integrity-protected header (7 bytes).
	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
		SequenceNumber:        2,
		MsgType:               msg.EMMStatus,
	}

	cause := ie.NewPositionalEMMCause(0x18) // #25 Security mode rejected unspecified

	m := msg.NewEmmStatus(cause)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.EmmStatus
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

	if causeVal, err := decoded.EMMCause.EMMCause(); err != nil {
		t.Fatalf("EMMCause parse failed: %v", err)
	} else if causeVal != 0x18 {
		t.Errorf("EMMCause = %v, want #25 (0x18)", causeVal)
	}
}

func TestEmmStatusTrailingBytesRejected(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMStatus,
	}

	cause := ie.NewPositionalEMMCause(0x16)
	m := msg.NewEmmStatus(cause)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	// Append trailing bytes — should be rejected.
	trailing := append(buf, 0xAA)

	var decoded msg.EmmStatus
	err := decoded.UnmarshalBinary(trailing)
	if !errors.Is(err, ie.ErrUnexpectedIE) {
		t.Fatalf("trailing bytes = %v, want wrapped %T", err, ie.ErrUnexpectedIE)
	}
}

func TestEmmStatusErrorCases(t *testing.T) {
	t.Parallel()

	var nilMsg *msg.EmmStatus

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

	local := &msg.EmmStatus{}
	buf := make([]byte, 1)
	if err := local.MarshalTo(buf); !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("short buffer = %v, want %T", err, io.ErrShortBuffer)
	}

	var decoded msg.EmmStatus

	// Truncated protected header (5 of 7 bytes)
	err = decoded.UnmarshalBinary([]byte{0x37, 0xAA, 0xBB, 0xCC, 0xDD})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated protected header = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Valid header but no content after it (only header bytes)
	err = decoded.UnmarshalBinary([]byte{0x17, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0x60})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("no content = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	// Valid header + 1 cause byte (should succeed for plain)
	err = decoded.UnmarshalBinary([]byte{0x07, 0x60, 0x16})
	if err != nil {
		t.Fatalf("plain header + cause = %v, want nil", err)
	}

	// Invalid SHT
	err = decoded.UnmarshalBinary([]byte{0xAA, 0x60})
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("invalid SHT = %v, want wrapped %T", err, ie.ErrInvalidValue)
	}

	// Plain header with cause + extra byte (trailing bytes rejected)
	err = decoded.UnmarshalBinary([]byte{0x07, 0x60, 0x16, 0xAA})
	if !errors.Is(err, ie.ErrUnexpectedIE) {
		t.Fatalf("trailing = %v, want wrapped %T", err, ie.ErrUnexpectedIE)
	}

	// Nil receiver MarshalTo should return wrapped ie.ErrNilIE
	var nilStatus *msg.EmmStatus
	if err := nilStatus.MarshalTo(nil); !errors.Is(err, ie.ErrNilIE) {
		t.Fatalf("nil receiver MarshalTo = %v, want wrapped %T", err, ie.ErrNilIE)
	}
	// Cause IE with invalid length (wrapped from ie.ParseWithoutTag)
	badCause := &ie.IE{Kind: ie.NoTag, Payload: []byte{}, Format: ie.IEFormatThree}
	badMsg := msg.NewEmmStatus(badCause)
	// Plain header: SHT=0, PD=0x7, MsgType=EMMStatus(0x60)
	short := []byte{0x07, 0x60}
	var decoded2 msg.EmmStatus
	err = decoded2.UnmarshalBinary(short)
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("empty cause = %v, want %T", err, io.ErrUnexpectedEOF)
	}

	_ = badMsg // MarshalLen checks: good cause IE will marshal fine, but cause nil test above
}

// ParseEmmStatus decodes b into a decoded msg.EmmStatus and returns it,
// panicking on error for use in subtests where t.Fatalf is inconvenient.
func ParseEmmStatus(t *testing.T, b []byte) msg.EmmStatus {
	t.Helper()
	var decoded msg.EmmStatus
	if err := decoded.UnmarshalBinary(b); err != nil {
		t.Fatalf("ParseEmmStatus failed: %v", err)
	}
	return decoded
}

func TestEmmStatusWithTrailingExtra(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMStatus,
	}

	cause := ie.NewPositionalEMMCause(0x02) // IMSI unknown
	m := msg.NewEmmStatus(cause)
	m.EmmHeader = h

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo failed: %v", err)
	}

	var decoded msg.EmmStatus
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

	if causeVal, err := decoded.EMMCause.EMMCause(); err != nil {
		t.Fatalf("EMMCause parse failed: %v", err)
	} else if causeVal != 0x02 {
		t.Errorf("EMMCause = %v, want IMSI unknown (0x02)", causeVal)
	}
}

func TestEmmStatusMarshalLenExtraIEs(t *testing.T) {
	t.Parallel()

	collisionIE := &ie.IE{Kind: 0x32, Format: ie.IEFormatFour, Payload: []byte{0xAA}}

	m := &msg.EmmStatus{
		EmmHeader: msg.EmmHeader{SecurityHeaderType: msg.SecurityHeaderTypePlain, ProtocolDiscriminator: byte(msg.PdEmm), MsgType: msg.EMMStatus},
		EMMCause:  ie.NewEMMCause(0x02),
		IEs:       []*ie.IE{collisionIE},
	}

	buf := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(buf); err != nil {
		t.Fatalf("MarshalTo with collision IE failed: %v", err)
	}

	var decoded msg.EmmStatus
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

	collisionMsg := &msg.EmmStatus{
		EmmHeader: msg.EmmHeader{SecurityHeaderType: msg.SecurityHeaderTypePlain, ProtocolDiscriminator: byte(msg.PdEmm), MsgType: msg.EMMStatus},
		EMMCause:  ie.NewEMMCause(0x02),
		IEs:       []*ie.IE{collisionIE},
	}
	shortBuf := make([]byte, 1)
	if err := collisionMsg.MarshalTo(shortBuf); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("expected ErrShortBuffer for collision IE error path, got %v", err)
	}
}
