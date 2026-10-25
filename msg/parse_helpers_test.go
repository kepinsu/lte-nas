package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/msg"
)

// TestParseHelperErrorCases tests that every standalone ParseXXX function returns
// a proper error when given invalid/empty input, covering the error-return branch.

func TestParseHelperErrorCases(t *testing.T) {
	t.Parallel()

	// ---- EMM Parse helpers ----

	t.Run("ParseActivateDedicatedBearerContextAccept", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParseActivateDedicatedEPSBearerContextAccept(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("ParseActivateDedicatedBearerContextReject", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParseActivateDedicatedEPSBearerContextReject(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("ParseActivateDedicatedBearerContextRequest", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParseActivateDedicatedEPSBearerContextRequest(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("ParseActivateDefaultBearerContextAccept", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParseActivateDefaultEPSBearerContextAccept(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("ParseActivateDefaultBearerContextReject", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParseActivateDefaultEPSBearerContextReject(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("ParseActivateDefaultBearerContextRequest", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParseActivateDefaultEPSBearerContextRequest(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("ParseBearerResourceAllocationReject", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParseBearerResourceAllocationReject(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("ParseBearerResourceAllocationRequest", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParseBearerResourceAllocationRequest(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("ParseBearerResourceModificationReject", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParseBearerResourceModificationReject(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("ParseBearerResourceModificationRequest", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParseBearerResourceModificationRequest(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("ParseDeactivateBearerContextAccept", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParseDeactivateEPSBearerContextAccept(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("ParseDeactivateBearerContextRequest", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParseDeactivateEPSBearerContextRequest(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("ParseModifyBearerContextAccept", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParseModifyEPSBearerContextAccept(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("ParseModifyBearerContextReject", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParseModifyEPSBearerContextReject(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("ParseModifyBearerContextRequest", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParseModifyEPSBearerContextRequest(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("ParsePDNConnectivityReject", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParsePDNConnectivityReject(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("ParsePDNConnectivityRequest", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParsePDNConnectivityRequest(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("ParsePDNDisconnectReject", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParsePDNDisconnectReject(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("ParsePDNDisconnectRequest", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParsePDNDisconnectRequest(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("ParseRemoteUEReportResponse", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParseRemoteUEReportResponse(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("ParseNotification", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParseNotification(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("ParseESMStatus", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParseESMStatusMsg(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("ParseEsmInfoRequest", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParseEsmInfoRequest(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("ParseEsmInfoResponse", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParseEsmInfoResponse(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("ParseEsmDataTransport", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParseEsmDataTransport(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("ParseRemoteUEReport", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParseRemoteUEReport(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})

	t.Run("ParseESMDummy", func(t *testing.T) {
		t.Parallel()
		_, err := msg.ParseESMDummy(nil)
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("nil input: expected %T, got %v", io.ErrUnexpectedEOF, err)
		}
	})
}
