package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

func TestEmmHeaderRoundTrip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		header  msg.EmmHeader
		wantLen int
		payload []byte
	}{
		{
			name: "plain header",
			header: msg.EmmHeader{
				SecurityHeaderType:    msg.SecurityHeaderTypePlain,
				ProtocolDiscriminator: byte(msg.PdEmm),
				MsgType:               msg.EMMAttachRequest,
			},
			wantLen: 2,
			payload: []byte{0xAA}, // post-header IE data
		},
		{
			name: "integrity protected header",
			header: msg.EmmHeader{
				SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
				ProtocolDiscriminator: byte(msg.PdEmm),
				Mac:                   [4]byte{0xAA, 0xBB, 0xCC, 0xDD},
				SequenceNumber:        5,
				MsgType:               msg.EMMAttachAccept,
			},
			wantLen: 7,
			payload: []byte{0xBB}, // post-header IE data
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := make([]byte, 16)

			err := tt.header.MarshalTo(buf)
			if err != nil {
				t.Fatalf("MarshalTo failed: %v", err)
			}

			// Place payload bytes after the header for unmarshal verification.
			copy(buf[tt.wantLen:], tt.payload)

			gotLen := tt.header.MarshalLen()
			if gotLen != tt.wantLen {
				t.Errorf("MarshalLen returned %d, want %d", gotLen, tt.wantLen)
			}

			var decoded msg.EmmHeader
			rem, err := decoded.UnmarshalBinary(buf[:tt.wantLen+len(tt.payload)])
			if err != nil {
				t.Fatalf("UnmarshalBinary failed: %v", err)
			}

			if decoded.SecurityHeaderType != tt.header.SecurityHeaderType {
				t.Errorf("SecurityHeaderType = %#02x, want %#02x", decoded.SecurityHeaderType, tt.header.SecurityHeaderType)
			}
			if decoded.ProtocolDiscriminator != tt.header.ProtocolDiscriminator {
				t.Errorf("ProtocolDiscriminator = %#02x, want %#02x", decoded.ProtocolDiscriminator, tt.header.ProtocolDiscriminator)
			}
			if decoded.MsgType != tt.header.MsgType {
				t.Errorf("MsgType = %#02x, want %#02x", decoded.MsgType, tt.header.MsgType)
			}

			if tt.wantLen == 7 {
				for i := range 4 {
					if decoded.Mac[i] != tt.header.Mac[i] {
						t.Errorf("Mac[%d] = %#02x, want %#02x", i, decoded.Mac[i], tt.header.Mac[i])
					}
				}
				if decoded.SequenceNumber != tt.header.SequenceNumber {
					t.Errorf("SequenceNumber = %d, want %d", decoded.SequenceNumber, tt.header.SequenceNumber)
				}
			}

			if len(rem) != len(tt.payload) {
				t.Fatalf("remaining length = %d, want %d", len(rem), len(tt.payload))
			}
			for i := range tt.payload {
				if rem[i] != tt.payload[i] {
					t.Errorf("payload[%d] = %#02x, want %#02x", i, rem[i], tt.payload[i])
				}
			}
		})
	}
}

func TestEmmHeaderNilMarshalLen(t *testing.T) {
	t.Parallel()

	var h *msg.EmmHeader
	if got := h.MarshalLen(); got != 0 {
		t.Errorf("nil MarshalLen = %d, want 0", got)
	}
}

func TestEmmHeaderMarshalShortBuffer(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypeIntegrityProtected,
		ProtocolDiscriminator: byte(msg.PdEmm),
		Mac:                   [4]byte{0x11, 0x22, 0x33, 0x44},
		SequenceNumber:        1,
		MsgType:               msg.EMMAttachRequest,
	}

	buf := make([]byte, 3)
	err := h.MarshalTo(buf)
	if !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("MarshalTo on short buffer = %T, want error type", err)
	}
}

func TestEmmHeaderInvalidPD(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: 0x11,
		MsgType:               msg.EMMAttachRequest,
	}

	buf := make([]byte, 2)
	err := h.MarshalTo(buf)
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Errorf("MarshalTo with invalid PD = %T, want error type", err)
	}

	var decoded msg.EmmHeader
	_, err = decoded.UnmarshalBinary([]byte{0x11, 0x41})
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Errorf("UnmarshalBinary with invalid PD = %T, want error type", err)
	}
}

func TestEmmHeaderInvalidSHT(t *testing.T) {
	t.Parallel()

	h := msg.EmmHeader{
		SecurityHeaderType:    0x7E,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msg.EMMAttachRequest,
	}

	buf := make([]byte, 2)
	err := h.MarshalTo(buf)
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Errorf("MarshalTo with invalid SHT = %T, want error type", err)
	}

	var decoded msg.EmmHeader
	_, err = decoded.UnmarshalBinary([]byte{0x77, 0x41})
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Errorf("UnmarshalBinary with invalid SHT = %T, want error type", err)
	}
}

func TestEmmHeaderTruncation(t *testing.T) {
	t.Parallel()

	var h msg.EmmHeader
	_, err := h.UnmarshalBinary([]byte{0x7E})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Errorf("UnmarshalBinary truncated = %T, want error type", err)
	}

	integ := []byte{0x17, 0xAA, 0xBB, 0xCC, 0xDD}
	_, err = h.UnmarshalBinary(integ)
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Errorf("UnmarshalBinary integrity truncated = %T, want error type", err)
	}
}

func TestEsmHeaderRoundTrip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		header  msg.EsmHeader
		payload []byte
	}{
		{
			name: "default bearer identity",
			header: msg.EsmHeader{
				BearerIdentity:               0x00,
				ProtocolDiscriminator:        byte(msg.PdEsm),
				ProcedureTransactionIdentity: 0x01,
				MsgType:                      msg.ESMActivateDefaultEPSBearerContextRequest,
			},
			payload: []byte{0x11},
		},
		{
			name: "bearer identity 7",
			header: msg.EsmHeader{
				BearerIdentity:               0x07,
				ProtocolDiscriminator:        byte(msg.PdEsm),
				ProcedureTransactionIdentity: 0xFF,
				MsgType:                      msg.ESMDataTransport,
			},
			payload: []byte{0x22},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := make([]byte, 16)

			err := tt.header.MarshalTo(buf)
			if err != nil {
				t.Fatalf("MarshalTo failed: %v", err)
			}

			copy(buf[3:], tt.payload)

			gotLen := tt.header.MarshalLen()
			if gotLen != 3 {
				t.Errorf("MarshalLen returned %d, want 3", gotLen)
			}

			var decoded msg.EsmHeader
			rem, err := decoded.UnmarshalBinary(buf[:3+len(tt.payload)])
			if err != nil {
				t.Fatalf("UnmarshalBinary failed: %v", err)
			}

			if decoded.BearerIdentity != tt.header.BearerIdentity {
				t.Errorf("BearerIdentity = %#02x, want %#02x", decoded.BearerIdentity, tt.header.BearerIdentity)
			}
			if decoded.MsgType != tt.header.MsgType {
				t.Errorf("MsgType = %#02x, want %#02x", decoded.MsgType, tt.header.MsgType)
			}
			if decoded.PTI() != tt.header.ProcedureTransactionIdentity {
				t.Errorf("PTI = %#02x, want %#02x", decoded.PTI(), tt.header.ProcedureTransactionIdentity)
			}

			if got := decoded.EPSBearerIdentity(); tt.header.BearerIdentity != got {
				t.Errorf("EPSBearerIdentity = %d, want %d", got, tt.header.BearerIdentity)
			}

			if len(rem) != len(tt.payload) {
				t.Fatalf("remaining length = %d, want %d", len(rem), len(tt.payload))
			}
		})
	}
}

func TestEsmHeaderNilMarshalLen(t *testing.T) {
	t.Parallel()

	var h *msg.EsmHeader
	if got := h.MarshalLen(); got != 0 {
		t.Errorf("nil MarshalLen = %d, want 0", got)
	}
}

func TestEsmHeaderShortBuffer(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:        0x00,
		ProtocolDiscriminator: msg.PdEsm,
		MsgType:               msg.ESMDataTransport,
	}

	buf := make([]byte, 2)
	err := h.MarshalTo(buf)
	if !errors.Is(err, io.ErrShortBuffer) {
		t.Errorf("MarshalTo on short buffer = %T, want error type", err)
	}
}

func TestEsmHeaderInvalidPD(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:        0x00,
		ProtocolDiscriminator: 0xFF,
		MsgType:               msg.ESMDataTransport,
	}

	buf := make([]byte, 3)
	err := h.MarshalTo(buf)
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Errorf("MarshalTo with invalid PD = %T, want error type", err)
	}

	var decoded msg.EsmHeader
	_, err = decoded.UnmarshalBinary([]byte{0xFF, 0x00, msg.ESMDataTransport})
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Errorf("UnmarshalBinary with invalid PD = %T, want error type", err)
	}
}

func TestEsmHeaderInvalidBearerID(t *testing.T) {
	t.Parallel()

	h := msg.EsmHeader{
		BearerIdentity:        0x10,
		ProtocolDiscriminator: byte(msg.PdEsm),
		MsgType:               msg.ESMDataTransport,
	}

	buf := make([]byte, 3)
	err := h.MarshalTo(buf)
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Errorf("MarshalTo with invalid bearer ID = %T, want error type", err)
	}
}

func TestEsmHeaderTruncation(t *testing.T) {
	t.Parallel()

	var h msg.EsmHeader
	_, err := h.UnmarshalBinary([]byte{0x2E, 0x00})
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Errorf("UnmarshalBinary truncated = %T, want error type", err)
	}
}
