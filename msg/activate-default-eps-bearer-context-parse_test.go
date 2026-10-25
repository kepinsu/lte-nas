package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

// TestActivateDefaultBearerContextParseHelpers covers the standalone ParseXXX
// functions for the three activate-default EPS bearer context message types,
// exercising both the success and error paths.

func TestActivateDefaultBearerContextParseHelpers(t *testing.T) {
	t.Parallel()

	// -----------------------------------------------------------------------
	// Activate Default EPS Bearer Context ACCEPT
	// Plain ESM header (3 bytes) + QoS LV + PCO LV + optional ExtraPCO.
	// -----------------------------------------------------------------------
	t.Run("ParseActivateDefaultBearerContextAccept", func(t *testing.T) {
		t.Parallel()

		header := []byte{
			0xA2, 0x01, msg.ESMActivateDefaultEPSBearerContextAccept,
		}

		// QoS (5 bytes positional): length=5 + 5 payload
		qos := []byte{0x05, 0x72, 0x05, 0x0D, 0x0A, 0xBB}

		// PCO (2 bytes positional): length=0 + padding
		pco := []byte{0x01, 0x00}

		// ExtraPCO: tag=0x7B, format six (LV-E), length 2, payload
		extraPCO := ie.NewRaw(0x7b, ie.IEFormatSix, []byte{0x02, 0xAA})

		// Build buffer
		buf := make([]byte, len(header)+len(qos)+len(pco)+extraPCO.MarshalLen())
		pos := copy(buf, header)
		pos += copy(buf[pos:], qos)
		pos += copy(buf[pos:], pco)
		extraPCO.MarshalTo(buf[pos:])

		got, err := msg.ParseActivateDefaultEPSBearerContextAccept(buf)
		if err != nil {
			t.Fatalf("ParseActivateDefaultEPSBearerContextAccept failed: %v", err)
		}
		if got == nil {
			t.Fatal("Parse returned nil")
		}
		if got.IsRequest() {
			t.Error("IsRequest = true, want false")
		}
		if got.EPSQoS == nil {
			t.Error("EPSQoS is nil")
		}
		if got.ProtocolConfigurationOpts == nil {
			t.Error("ProtocolConfigurationOpts is nil")
		}
		if got.ExtendedPCO == nil {
			t.Error("ExtendedPCO is nil")
		}

		// Round-trip
		rebuf := make([]byte, got.MarshalLen())
		if err := got.MarshalTo(rebuf); err != nil {
			t.Fatalf("re-MarshalTo failed: %v", err)
		}
		for i := range buf {
			if buf[i] != rebuf[i] {
				t.Fatalf("byte[%d] = %#02x, want %#02x", i, rebuf[i], buf[i])
			}
		}
	})

	// -----------------------------------------------------------------------
	// Activate Default EPS Bearer Context REJECT
	// Plain ESM header + ESM cause (1 byte).
	// -----------------------------------------------------------------------
	t.Run("ParseActivateDefaultBearerContextReject", func(t *testing.T) {
		t.Parallel()

		buf := make([]byte, 3+1)
		buf[0] = 0xA2 // PD=ESM, Bearer ID=0x0
		buf[1] = 0x01 // PTI
		buf[2] = msg.ESMActivateDefaultEPSBearerContextReject
		buf[3] = 0x68 // ESM cause

		got, err := msg.ParseActivateDefaultEPSBearerContextReject(buf)
		if err != nil {
			t.Fatalf("ParseActivateDefaultEPSBearerContextReject failed: %v", err)
		}
		if got == nil {
			t.Fatal("Parse returned nil")
		}
		if got.ESMCause == nil {
			t.Fatal("ESMCause is nil")
		}
		if got.IsRequest() {
			t.Error("IsRequest = true, want false")
		}

		rebuf := make([]byte, got.MarshalLen())
		if err := got.MarshalTo(rebuf); err != nil {
			t.Fatalf("re-MarshalTo failed: %v", err)
		}
		for i := range buf {
			if buf[i] != rebuf[i] {
				t.Fatalf("byte[%d] = %#02x, want %#02x", i, rebuf[i], buf[i])
			}
		}

		// Error case: empty buffer
		_, err = msg.ParseActivateDefaultEPSBearerContextReject(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil buffer: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}

		// Error case: header only (no payload)
		_, err = msg.ParseActivateDefaultEPSBearerContextReject(buf[:3])
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("header only: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	// -----------------------------------------------------------------------
	// Activate Default EPS Bearer Context REQUEST
	// Plain ESM header + QoS + APN + PDN Address (5+ bytes) + optional LLCSAPI + ESM Cause
	// -----------------------------------------------------------------------
	t.Run("ParseActivateDefaultBearerContextRequest", func(t *testing.T) {
		t.Parallel()

		// PDN address requires at least 5 bytes value per validation.
		buf := make([]byte, 3+6+2+6) // header + qos(1+5) + apn(1+len) + pdnaddr(1+5=len6)
		buf[0] = 0xA2
		buf[1] = 0x02
		buf[2] = msg.ESMActivateDefaultEPSBearerContextRequest

		pos := 3
		qos := []byte{0x05, 0x72, 0x05, 0x0D, 0x0A, 0xBB}
		copy(buf[pos:], qos)
		pos += len(qos)

		apn := []byte{0x01, 0x63} // len=1, "c"
		copy(buf[pos:], apn)
		pos += len(apn)

		pdn := []byte{0x05, 0x01, 0xC0, 0xA8, 0x01, 0x01} // len=5, IPv4 type+addr
		copy(buf[pos:], pdn)

		got, err := msg.ParseActivateDefaultEPSBearerContextRequest(buf)
		if err != nil {
			t.Fatalf("ParseActivateDefaultEPSBearerContextRequest failed: %v", err)
		}
		if got == nil {
			t.Fatal("Parse returned nil")
		}
		if got.EPSQoS == nil {
			t.Fatal("EPSQoS is nil")
		}
		if got.AccessPointName == nil {
			t.Fatal("AccessPointName is nil")
		}
		if got.PDNAddress == nil {
			t.Fatal("PDNAddress is nil")
		}

		rebuf := make([]byte, got.MarshalLen())
		if err := got.MarshalTo(rebuf); err != nil {
			t.Fatalf("re-MarshalTo failed: %v", err)
		}
		for i := range buf {
			if buf[i] != rebuf[i] {
				t.Fatalf("byte[%d] = %#02x, want %#02x", i, rebuf[i], buf[i])
			}
		}

		// Error: nil buffer
		_, err = msg.ParseActivateDefaultEPSBearerContextRequest(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}

		// Error: truncated payload
		_, err = msg.ParseActivateDefaultEPSBearerContextRequest(buf[:4])
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("truncated: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})
}
