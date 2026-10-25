package msg_test

import (
	"testing"

	"github.com/kepinsu/lte-nas/msg"
)

func TestAllConstantsUnique(t *testing.T) {
	t.Parallel()

	all := []uint8{
		msg.EMMAttachRequest, msg.EMMAttachAccept, msg.EMMAttachComplete, msg.EMMAttachReject,
		msg.EMMDetachRequest, msg.EMMDetachAccept,
		msg.EMMTrackingAreaUpdateRequest, msg.EMMTrackingAreaUpdateAccept, msg.EMMTrackingAreaUpdateComplete, msg.EMMTrackingAreaUpdateReject,
		msg.EMMExtendedServiceRequest, msg.EMMControlPlaneServiceRequest, msg.EMMServiceReject, msg.EMMServiceAccept,
		msg.EMMGUTIReallocationCommand, msg.EMMGUTIReallocationComplete,
		msg.EMMAuthenticationRequest, msg.EMMAuthenticationResponse, msg.EMMAuthenticationReject, msg.EMMAuthenticationFailure,
		msg.EMMIdentityRequest, msg.EMMIdentityResponse,
		msg.EMMSecurityModeCommand, msg.EMMSecurityModeComplete, msg.EMMSecurityModeReject,
		msg.EMMStatus, msg.EMMInformation,
		msg.EMMDownlinkNASTransport, msg.EMMUplinkNASTransport,
		msg.EMMCSServiceNotification,
		msg.EMMDownlinkGenericNASTransport, msg.EMMUplinkGenericNASTransport,
		msg.ESMActivateDefaultEPSBearerContextRequest, msg.ESMActivateDefaultEPSBearerContextAccept, msg.ESMActivateDefaultEPSBearerContextReject,
		msg.ESMActivateDedicatedEPSBearerContextRequest, msg.ESMActivateDedicatedEPSBearerContextAccept, msg.ESMActivateDedicatedEPSBearerContextReject,
		msg.ESMModifyEPSBearerContextRequest, msg.ESMModifyEPSBearerContextAccept, msg.ESMModifyEPSBearerContextReject,
		msg.ESMDeactivateEPSBearerContextRequest, msg.ESMDeactivateEPSBearerContextAccept,
		msg.ESMPDNConnectivityRequest, msg.ESMPDNConnectivityReject,
		msg.ESMPDNDisconnectRequest, msg.ESMPDNDisconnectReject,
		msg.ESMBearerResourceAllocationRequest, msg.ESMBearerResourceAllocationReject,
		msg.ESMBearerResourceModificationRequest, msg.ESMBearerResourceModificationReject,
		msg.ESMInformationRequest, msg.ESMInformationResponse, msg.ESMNotification,
		msg.ESMDummy,
		msg.ESMStatus, msg.ESMRemoteUEReport, msg.ESMRemoteUEReportResponse, msg.ESMDataTransport,
	}

	if len(all) != 59 {
		t.Fatalf("expected 59 constants, got %d", len(all))
	}

	seen := make(map[uint8]int, 59)
	for i, v := range all {
		if dup, ok := seen[v]; ok {
			t.Errorf("duplicate value %#02x at indices %d and %d", v, dup, i)
		}
		seen[v] = i
	}
}

func TestEMMTopBits(t *testing.T) {
	t.Parallel()

	emmValues := []uint8{
		msg.EMMAttachRequest, msg.EMMAttachAccept, msg.EMMAttachComplete, msg.EMMAttachReject,
		msg.EMMDetachRequest, msg.EMMDetachAccept,
		msg.EMMTrackingAreaUpdateRequest, msg.EMMTrackingAreaUpdateAccept, msg.EMMTrackingAreaUpdateComplete, msg.EMMTrackingAreaUpdateReject,
		msg.EMMExtendedServiceRequest, msg.EMMControlPlaneServiceRequest, msg.EMMServiceReject, msg.EMMServiceAccept,
		msg.EMMGUTIReallocationCommand, msg.EMMGUTIReallocationComplete,
		msg.EMMAuthenticationRequest, msg.EMMAuthenticationResponse, msg.EMMAuthenticationReject, msg.EMMAuthenticationFailure,
		msg.EMMIdentityRequest, msg.EMMIdentityResponse,
		msg.EMMSecurityModeCommand, msg.EMMSecurityModeComplete, msg.EMMSecurityModeReject,
		msg.EMMStatus, msg.EMMInformation,
		msg.EMMDownlinkNASTransport, msg.EMMUplinkNASTransport,
		msg.EMMCSServiceNotification,
		msg.EMMDownlinkGenericNASTransport, msg.EMMUplinkGenericNASTransport,
	}

	for _, v := range emmValues {
		if (v>>6)&0x03 != 0x01 {
			t.Errorf("EMM value %#02x top two bits = %b, want 01", v, v>>6)
		}
	}
	if len(emmValues) != 32 {
		t.Errorf("EMM count = %d, want 32", len(emmValues))
	}
}

func TestESMTopBits(t *testing.T) {
	t.Parallel()

	esmValues := []uint8{
		msg.ESMActivateDefaultEPSBearerContextRequest, msg.ESMActivateDefaultEPSBearerContextAccept, msg.ESMActivateDefaultEPSBearerContextReject,
		msg.ESMActivateDedicatedEPSBearerContextRequest, msg.ESMActivateDedicatedEPSBearerContextAccept, msg.ESMActivateDedicatedEPSBearerContextReject,
		msg.ESMModifyEPSBearerContextRequest, msg.ESMModifyEPSBearerContextAccept, msg.ESMModifyEPSBearerContextReject,
		msg.ESMDeactivateEPSBearerContextRequest, msg.ESMDeactivateEPSBearerContextAccept,
		msg.ESMPDNConnectivityRequest, msg.ESMPDNConnectivityReject,
		msg.ESMPDNDisconnectRequest, msg.ESMPDNDisconnectReject,
		msg.ESMBearerResourceAllocationRequest, msg.ESMBearerResourceAllocationReject,
		msg.ESMBearerResourceModificationRequest, msg.ESMBearerResourceModificationReject,
		msg.ESMInformationRequest, msg.ESMInformationResponse, msg.ESMNotification,
		msg.ESMDummy,
		msg.ESMStatus, msg.ESMRemoteUEReport, msg.ESMRemoteUEReportResponse, msg.ESMDataTransport,
	}

	for _, v := range esmValues {
		if (v>>6)&0x03 != 0x03 {
			t.Errorf("ESM value %#02x top two bits = %b, want 11", v, v>>6)
		}
	}
	if len(esmValues) != 27 {
		t.Errorf("ESM count = %d, want 27", len(esmValues))
	}
}

func TestConstantSpotChecks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		got  uint8
		want uint8
	}{
		{"EMM AttachRequest", msg.EMMAttachRequest, 0x41},
		{"EMM AuthenticationRequest", msg.EMMAuthenticationRequest, 0x52},
		{"EMM SecurityModeCommand", msg.EMMSecurityModeCommand, 0x5D},
		{"ESM ActivateDefault Request", msg.ESMActivateDefaultEPSBearerContextRequest, 0xC1},
		{"ESM PDNConnectivityRequest", msg.ESMPDNConnectivityRequest, 0xD0},
		{"ESM DataTransport", msg.ESMDataTransport, 0xEB},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s = %#02x, want %#02x", tt.name, tt.got, tt.want)
			}
		})
	}
}
