package msg_test

import (
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/msg"
)

// ---------------------------------------------------------------------------
// Helper: plain EMM header for the given message type.
// ---------------------------------------------------------------------------

func makeEMMHeaderPlain(msgType uint8) *msg.EmmHeader {
	return &msg.EmmHeader{
		SecurityHeaderType:    msg.SecurityHeaderTypePlain,
		ProtocolDiscriminator: byte(msg.PdEmm),
		MsgType:               msgType,
	}
}

func makeESMHeader(msgType uint8) *msg.EsmHeader {
	return &msg.EsmHeader{
		BearerIdentity:               0x0,
		ProtocolDiscriminator:        byte(msg.PdEsm),
		ProcedureTransactionIdentity: 0x1,
		MsgType:                      msgType,
	}
}

// validQoS returns a 5-byte EPS Quality of Service payload (valid per IE validation).
var validQoS = []byte{0x72, 0x05, 0x0D, 0x0A, 0xBB}

// ---------------------------------------------------------------------------
// EMM message constructors
// ---------------------------------------------------------------------------

func newAttachRequest() msg.Message {
	m := msg.NewAttachRequest(
		ie.NewPositionalEPSAttachType(0x1),
		ie.NewPositionalNASKeySetIdentifier(0x1),
		ie.NewPositionalEPSMobileIdentity([]byte{0x23, 0x05}),
		ie.NewPositionalUENetworkCapability([]byte{0xAA, 0xBB}),
		ie.NewPositionalESMMessageContainer([]byte{0xD0, 0x01}),
	)
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMAttachRequest)
	return m
}

func newAttachAccept() msg.Message {
	m := msg.NewAttachAccept(
		ie.NewPositionalEPSAttachResult(0x1),
		ie.NewPositionalSpareHalfOctet(0x0F),
		ie.NewPositionalGPRSTimer(0x05),
		ie.NewPositionalTrackingAreaIdentityList([]byte{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC}),
		ie.NewPositionalESMMessageContainer([]byte{0xC1, 0x01}),
	)
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMAttachAccept)
	return m
}

func newAttachComplete() msg.Message {
	// AttachComplete has mandatory ESM Message Container even though New takes variadic only.
	// The existing tests construct it without positional IEs but UnmarshalBinary requires ESM container.
	// Use minimal header only — it round-trips via existing test pattern.
	m := msg.NewAttachComplete(ie.NewESMMessageContainer([]byte{0xC1, 0x02, 0xAA, 0xBB}))
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMAttachComplete)
	return m
}

func newAttachReject() msg.Message {
	m := msg.NewAttachReject(ie.NewEMMCause(0x1))
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMAttachReject)
	return m
}

func newDetachRequest() msg.Message {
	m := msg.NewDetachRequest(
		ie.NewPositionalDetachType(0x1),
		ie.NewPositionalNASKeySetIdentifier(0x1),
		ie.NewPositionalEPSMobileIdentity([]byte{0x23, 0x05}),
	)
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMDetachRequest)
	return m
}

func newDetachAccept() msg.Message {
	m := msg.NewDetachAccept(nil)
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMDetachAccept)
	return m
}

func newTAURequest() msg.Message {
	m := msg.NewTrackingAreaUpdateRequest(
		ie.NewPositionalEPSUpdateType(0x1),
		ie.NewPositionalNASKeySetIdentifier(0x1),
		ie.NewEPSMobileIdentity([]byte{0x40, 0x05, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE}),
	)
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMTrackingAreaUpdateRequest)
	return m
}

func newTAUAccept() msg.Message {
	m := msg.NewTrackingAreaUpdateAccept(
		ie.NewPositionalEPSUpdateResult(0x1),
		ie.NewPositionalSpareHalfOctet(0x0F),
	)
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMTrackingAreaUpdateAccept)
	return m
}

func newTAUComplete() msg.Message {
	m := msg.NewTrackingAreaUpdateComplete(nil)
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMTrackingAreaUpdateComplete)
	return m
}

func newTAUReject() msg.Message {
	m := msg.NewTrackingAreaUpdateReject(ie.NewEMMCause(0x1))
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMTrackingAreaUpdateReject)
	return m
}

func newExtServiceRequest() msg.Message {
	m := msg.NewExtendedServiceRequest(
		ie.NewPositionalServiceType(0x1),
		ie.NewPositionalNASKeySetIdentifier(0x1),
		ie.NewPositionalEPSMobileIdentity([]byte{0x23, 0x05}),
	)
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMExtendedServiceRequest)
	return m
}

func newCPServiceRequest() msg.Message {
	m := msg.NewControlPlaneServiceRequest(
		ie.NewPositionalControlPlaneServiceType(0x0),
		ie.NewPositionalNASKeySetIdentifier(0x1),
	)
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMControlPlaneServiceRequest)
	return m
}

func newServiceReject() msg.Message {
	m := msg.NewServiceReject(ie.NewEMMCause(0x1))
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMServiceReject)
	return m
}

func newServiceAccept() msg.Message {
	m := msg.NewServiceAccept(nil)
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMServiceAccept)
	return m
}

func newGUTICommand() msg.Message {
	m := msg.NewGUTIReallocationCommand(
		ie.NewPositionalEPSMobileIdentity([]byte{0x23, 0x05}),
	)
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMGUTIReallocationCommand)
	return m
}

func newGUTICComplete() msg.Message {
	m := msg.NewGUTIReallocationComplete(nil)
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMGUTIReallocationComplete)
	return m
}

func newAuthRequest() msg.Message {
	m := msg.NewAuthenticationRequest(
		ie.NewPositionalNASKeySetIdentifier(0x1),
		ie.NewPositionalSpareHalfOctet(0x0),
		ie.NewPositionalAuthenticationParameterRAND(make([]byte, 16)),
		ie.NewPositionalAuthenticationParameterAUTN(make([]byte, 16)),
	)
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMAuthenticationRequest)
	return m
}

func newAuthResponse() msg.Message {
	m := msg.NewAuthenticationResponse(
		ie.NewPositionalAuthenticationResponseParameter([]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE}),
	)
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMAuthenticationResponse)
	return m
}

func newAuthReject() msg.Message {
	m := msg.NewAuthenticationReject(nil)
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMAuthenticationReject)
	return m
}

func newAuthFailure() msg.Message {
	m := msg.NewAuthenticationFailure(ie.NewEMMCause(0x1))
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMAuthenticationFailure)
	return m
}

func newIdentityRequest() msg.Message {
	m := msg.NewIdentityRequest(
		ie.NewPositionalIdentityType2(0x1),
		ie.NewPositionalSpareHalfOctet(0x0),
	)
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMIdentityRequest)
	return m
}

func newIdentityResponse() msg.Message {
	m := msg.NewIdentityResponse(ie.NewMobileIdentity([]byte{0x23, 0x05}))
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMIdentityResponse)
	return m
}

func newSecModeCommand() msg.Message {
	m := msg.NewSecurityModeCommand(
		ie.NewPositionalNASSecurityAlgorithms(0x0),
		ie.NewPositionalSpareHalfOctet(0x0),
		ie.NewPositionalUESecurityCapability([]byte{0xAA, 0xBB}),
	)
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMSecurityModeCommand)
	return m
}

func newSecModeComplete() msg.Message {
	m := msg.NewSecurityModeComplete(nil)
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMSecurityModeComplete)
	return m
}

func newSecModeReject() msg.Message {
	m := msg.NewSecurityModeReject(ie.NewEMMCause(0x1))
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMSecurityModeReject)
	return m
}

func newEmmStatus() msg.Message {
	m := msg.NewEmmStatus(ie.NewEMMCause(0x1))
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMStatus)
	return m
}

func newEmmInfo() msg.Message {
	m := msg.NewEmmInformation(nil)
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMInformation)
	return m
}

func newDownlinkNASTransport() msg.Message {
	// GUTI is type-3 V, one octet.
	guti := ie.NewRaw(ie.NoTag, ie.IEFormatThree, []byte{0x05})
	tmsi := ie.NewRaw(ie.NoTag, ie.IEFormatFour, []byte{0x12, 0x34, 0xAB, 0xCD})
	m := msg.NewDownlinkNASTransport(guti, tmsi)
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMDownlinkNASTransport)
	return m
}

func newUplinkNASTransport() msg.Message {
	m := msg.NewUplinkNASTransport(nil)
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMUplinkNASTransport)
	return m
}

func newCSServiceNotification() msg.Message {
	m := msg.NewCSServiceNotification(ie.NewPositionalPagingIdentity(0x1))
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMCSServiceNotification)
	return m
}

func newDownlinkGenericNASTransport() msg.Message {
	m := msg.NewDownlinkGenericNASTransport(ie.NewRaw(0x7E, ie.IEFormatFour, []byte{0xD0, 0x01}))
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMDownlinkGenericNASTransport)
	return m
}

func newUplinkGenericNASTransport() msg.Message {
	m := msg.NewUplinkGenericNASTransport(ie.NewRaw(0x7E, ie.IEFormatFour, []byte{0xD0, 0x01}))
	m.EmmHeader = *makeEMMHeaderPlain(msg.EMMUplinkGenericNASTransport)
	return m
}

// ---------------------------------------------------------------------------
// ESM message constructors
// ---------------------------------------------------------------------------

func newActivateDefaultBearerRequest() msg.Message {
	m := msg.NewActivateDefaultEPSBearerContextRequest(
		ie.NewPositionalEPSQualityOfService(validQoS),
		ie.NewAccessPointName([]byte("pdn")),
		ie.NewPositionalPDNAddress([]byte{0x1, 0xC0, 0xA8, 0x01, 0x01}),
	)
	m.EsmHeader = *makeESMHeader(msg.ESMActivateDefaultEPSBearerContextRequest)
	return m
}

func newActivateDefaultBearerAccept() msg.Message {
	m := msg.NewActivateDefaultEPSBearerContextAccept(
		ie.NewPositionalEPSQualityOfService(validQoS),
		ie.NewProtocolConfigurationOptions([]byte{0x00, 0x05}),
	)
	m.EsmHeader = *makeESMHeader(msg.ESMActivateDefaultEPSBearerContextAccept)
	return m
}

func newActivateDefaultBearerReject() msg.Message {
	m := msg.NewActivateDefaultEPSBearerContextReject(ie.NewESMCause(0x1))
	m.EsmHeader = *makeESMHeader(msg.ESMActivateDefaultEPSBearerContextReject)
	return m
}

func newActivateDedicatedBearerRequest() msg.Message {
	m := msg.NewActivateDedicatedEPSBearerContextRequest(
		ie.NewPositionalLinkedEPSBearerIdentity(0x1),
		ie.NewPositionalSpareHalfOctet(0x0),
		ie.NewPositionalEPSQualityOfService(validQoS),
		ie.NewTrafficFlowTemplate([]byte{0x01, 0x02}),
	)
	m.EsmHeader = *makeESMHeader(msg.ESMActivateDedicatedEPSBearerContextRequest)
	return m
}

func newActivateDedicatedBearerAccept() msg.Message {
	m := msg.NewActivateDedicatedEPSBearerContextAccept(
		ie.NewPositionalEPSQualityOfService(validQoS),
	)
	m.EsmHeader = *makeESMHeader(msg.ESMActivateDedicatedEPSBearerContextAccept)
	return m
}

func newActivateDedicatedBearerReject() msg.Message {
	m := msg.NewActivateDedicatedEPSBearerContextReject(ie.NewESMCause(0x1))
	m.EsmHeader = *makeESMHeader(msg.ESMActivateDedicatedEPSBearerContextReject)
	return m
}

func newModifyBearerRequest() msg.Message {
	m := msg.NewModifyEPSBearerContextRequest(
		ie.NewPositionalEPSQualityOfService(validQoS),
	)
	m.EsmHeader = *makeESMHeader(msg.ESMModifyEPSBearerContextRequest)
	return m
}

func newModifyBearerAccept() msg.Message {
	m := msg.NewModifyEPSBearerContextAccept(
		ie.NewPositionalEPSQualityOfService(validQoS),
	)
	m.EsmHeader = *makeESMHeader(msg.ESMModifyEPSBearerContextAccept)
	return m
}

func newModifyBearerReject() msg.Message {
	m := msg.NewModifyEPSBearerContextReject(ie.NewESMCause(0x1))
	m.EsmHeader = *makeESMHeader(msg.ESMModifyEPSBearerContextReject)
	return m
}

func newDeactivateBearerRequest() msg.Message {
	m := msg.NewDeactivateEPSBearerContextRequest(ie.NewESMCause(0x1))
	m.EsmHeader = *makeESMHeader(msg.ESMDeactivateEPSBearerContextRequest)
	return m
}

func newDeactivateBearerAccept() msg.Message {
	m := msg.NewDeactivateEPSBearerContextAccept(
		ie.NewPositionalEPSQualityOfService(validQoS),
	)
	m.EsmHeader = *makeESMHeader(msg.ESMDeactivateEPSBearerContextAccept)
	return m
}

func newPDNConnectivityRequest() msg.Message {
	m := msg.NewPDNConnectivityRequest(
		ie.NewPositionalPDNType(0x1),
		ie.NewPositionalRequestType(0x1),
		ie.NewPositionalEPSQualityOfService(validQoS),
		ie.NewAccessPointName([]byte("pdn")),
	)
	m.EsmHeader = *makeESMHeader(msg.ESMPDNConnectivityRequest)
	return m
}

func newPDNConnectivityReject() msg.Message {
	m := msg.NewPDNConnectivityReject(ie.NewESMCause(0x1))
	m.EsmHeader = *makeESMHeader(msg.ESMPDNConnectivityReject)
	return m
}

func newPDNDisconnectRequest() msg.Message {
	m := msg.NewPDNDisconnectRequest(nil)
	m.EsmHeader = *makeESMHeader(msg.ESMPDNDisconnectRequest)
	return m
}

func newPDNDisconnectReject() msg.Message {
	m := msg.NewPDNDisconnectReject(ie.NewESMCause(0x1))
	m.EsmHeader = *makeESMHeader(msg.ESMPDNDisconnectReject)
	return m
}

func newBearerResourceAllocRequest() msg.Message {
	m := msg.NewBearerResourceAllocationRequest(
		ie.NewPositionalLinkedEPSBearerIdentity(0x1),
		ie.NewPositionalSpareHalfOctet(0x0),
		ie.NewPositionalEPSQualityOfService(validQoS),
		ie.NewTrafficFlowTemplate([]byte{0x01, 0x02}),
	)
	m.EsmHeader = *makeESMHeader(msg.ESMBearerResourceAllocationRequest)
	return m
}

func newBearerResourceAllocReject() msg.Message {
	m := msg.NewBearerResourceAllocationReject(ie.NewESMCause(0x1))
	m.EsmHeader = *makeESMHeader(msg.ESMBearerResourceAllocationReject)
	return m
}

func newBearerResourceModRequest() msg.Message {
	m := msg.NewBearerResourceModificationRequest(
		ie.NewPositionalLinkedEPSBearerIdentity(0x1),
		ie.NewPositionalSpareHalfOctet(0x0),
		ie.NewPositionalEPSQualityOfService(validQoS),
		ie.NewTrafficFlowTemplate([]byte{0x01, 0x02}),
	)
	m.EsmHeader = *makeESMHeader(msg.ESMBearerResourceModificationRequest)
	return m
}

func newBearerResourceModReject() msg.Message {
	m := msg.NewBearerResourceModificationReject(ie.NewESMCause(0x1))
	m.EsmHeader = *makeESMHeader(msg.ESMBearerResourceModificationReject)
	return m
}

func newEsmInfoRequest() msg.Message {
	m := msg.NewEsmInfoRequest(
		ie.NewPositionalUENetworkCapability([]byte{0xAA, 0xBB}),
		ie.NewPositionalESMMessageContainer([]byte{0xD0, 0x01}),
	)
	m.EsmHeader = *makeESMHeader(msg.ESMInformationRequest)
	return m
}

func newEsmInfoResponse() msg.Message {
	m := msg.NewEsmInfoResponse(nil)
	m.EsmHeader = *makeESMHeader(msg.ESMInformationResponse)
	return m
}

func newNotification() msg.Message {
	m := msg.NewNotification(ie.NewPositionalNotificationIndicator(0x1))
	m.EsmHeader = *makeESMHeader(msg.ESMNotification)
	return m
}

func newDummy() msg.Message {
	m := msg.NewESMDummy()
	m.EsmHeader = *makeESMHeader(msg.ESMDummy)
	return m
}

func newEsmStatus() msg.Message {
	m := msg.NewESMStatusMsg(ie.NewESMCause(0x1))
	m.EsmHeader = *makeESMHeader(msg.ESMStatus)
	return m
}

func newRemoteUEReport() msg.Message {
	m := msg.NewRemoteUEReport(nil)
	m.EsmHeader = *makeESMHeader(msg.ESMRemoteUEReport)
	return m
}

func newRemoteUEReportResponse() msg.Message {
	m := msg.NewRemoteUEReportResponse(nil)
	m.EsmHeader = *makeESMHeader(msg.ESMRemoteUEReportResponse)
	return m
}

func newEsmDataTransport() msg.Message {
	m := msg.NewEsmDataTransport(nil)
	m.EsmHeader = *makeESMHeader(msg.ESMDataTransport)
	return m
}

// ---------------------------------------------------------------------------
// messageEntry pairs a name, constructor function, expected type and type name.
// ---------------------------------------------------------------------------

type messageEntry struct {
	name     string
	newMsg   func() msg.Message
	msgType  uint8
	typeName string
}

func allEMMEntries() []messageEntry {
	return []messageEntry{
		{"Attach request", newAttachRequest, msg.EMMAttachRequest, "Attach request"},
		{"Attach accept", newAttachAccept, msg.EMMAttachAccept, "Attach accept"},
		{"Attach complete", newAttachComplete, msg.EMMAttachComplete, "Attach complete"},
		{"Attach reject", newAttachReject, msg.EMMAttachReject, "Attach reject"},
		{"Detach request", newDetachRequest, msg.EMMDetachRequest, "Detach request"},
		{"Detach accept", newDetachAccept, msg.EMMDetachAccept, "Detach accept"},
		{"Tracking area update request", newTAURequest, msg.EMMTrackingAreaUpdateRequest, "Tracking area update request"},
		{"Tracking area update accept", newTAUAccept, msg.EMMTrackingAreaUpdateAccept, "Tracking area update accept"},
		{"Tracking area update complete", newTAUComplete, msg.EMMTrackingAreaUpdateComplete, "Tracking area update complete"},
		{"Tracking area update reject", newTAUReject, msg.EMMTrackingAreaUpdateReject, "Tracking area update reject"},
		{"Extended service request", newExtServiceRequest, msg.EMMExtendedServiceRequest, "Extended service request"},
		{"Control plane service request", newCPServiceRequest, msg.EMMControlPlaneServiceRequest, "Control plane service request"},
		{"Service reject", newServiceReject, msg.EMMServiceReject, "Service reject"},
		{"Service accept", newServiceAccept, msg.EMMServiceAccept, "Service accept"},
		{"GUTI reallocation command", newGUTICommand, msg.EMMGUTIReallocationCommand, "Guti reallocation command"},
		{"GUTI reallocation complete", newGUTICComplete, msg.EMMGUTIReallocationComplete, "Guti reallocation complete"},
		{"Authentication request", newAuthRequest, msg.EMMAuthenticationRequest, "Authentication request"},
		{"Authentication response", newAuthResponse, msg.EMMAuthenticationResponse, "Authentication response"},
		{"Authentication reject", newAuthReject, msg.EMMAuthenticationReject, "Authentication reject"},
		{"Authentication failure", newAuthFailure, msg.EMMAuthenticationFailure, "Authentication failure"},
		{"Identity request", newIdentityRequest, msg.EMMIdentityRequest, "Identity request"},
		{"Identity response", newIdentityResponse, msg.EMMIdentityResponse, "Identity response"},
		{"Security mode command", newSecModeCommand, msg.EMMSecurityModeCommand, "Security mode command"},
		{"Security mode complete", newSecModeComplete, msg.EMMSecurityModeComplete, "Security mode complete"},
		{"Security mode reject", newSecModeReject, msg.EMMSecurityModeReject, "Security mode reject"},
		{"EMM status", newEmmStatus, msg.EMMStatus, "EMM status"},
		{"EMM information", newEmmInfo, msg.EMMInformation, "EMM information"},
		{"Downlink NAS transport", newDownlinkNASTransport, msg.EMMDownlinkNASTransport, "Downlink NAS transport"},
		{"Uplink NAS transport", newUplinkNASTransport, msg.EMMUplinkNASTransport, "Uplink NAS transport"},
		{"CSS service notification", newCSServiceNotification, msg.EMMCSServiceNotification, "CS service notification"},
		{"Downlink generic NAS transport", newDownlinkGenericNASTransport, msg.EMMDownlinkGenericNASTransport, "Downlink generic NAS transport"},
		{"Uplink generic NAS transport", newUplinkGenericNASTransport, msg.EMMUplinkGenericNASTransport, "Uplink generic NAS transport"},
	}
}

func allESMEntries() []messageEntry {
	return []messageEntry{
		{"Activate default EPS bearer context request", newActivateDefaultBearerRequest, msg.ESMActivateDefaultEPSBearerContextRequest, "Activate default EPS bearer context request"},
		{"Activate default EPS bearer context accept", newActivateDefaultBearerAccept, msg.ESMActivateDefaultEPSBearerContextAccept, "Activate default EPS bearer context accept"},
		{"Activate default EPS bearer context reject", newActivateDefaultBearerReject, msg.ESMActivateDefaultEPSBearerContextReject, "Activate default EPS bearer context reject"},
		{"Activate dedicated EPS bearer context request", newActivateDedicatedBearerRequest, msg.ESMActivateDedicatedEPSBearerContextRequest, "Activate dedicated EPS bearer context request"},
		{"Activate dedicated EPS bearer context accept", newActivateDedicatedBearerAccept, msg.ESMActivateDedicatedEPSBearerContextAccept, "Activate dedicated EPS bearer context accept"},
		{"Activate dedicated EPS bearer context reject", newActivateDedicatedBearerReject, msg.ESMActivateDedicatedEPSBearerContextReject, "Activate dedicated EPS bearer context reject"},
		{"Modify EPS bearer context request", newModifyBearerRequest, msg.ESMModifyEPSBearerContextRequest, "Modify EPS bearer context request"},
		{"Modify EPS bearer context accept", newModifyBearerAccept, msg.ESMModifyEPSBearerContextAccept, "Modify EPS bearer context accept"},
		{"Modify EPS bearer context reject", newModifyBearerReject, msg.ESMModifyEPSBearerContextReject, "Modify EPS bearer context reject"},
		{"Deactivate EPS bearer context request", newDeactivateBearerRequest, msg.ESMDeactivateEPSBearerContextRequest, "Deactivate EPS bearer context request"},
		{"Deactivate EPS bearer context accept", newDeactivateBearerAccept, msg.ESMDeactivateEPSBearerContextAccept, "Deactivate EPS bearer context accept"},
		{"PDN connectivity request", newPDNConnectivityRequest, msg.ESMPDNConnectivityRequest, "PDN connectivity request"},
		{"PDN connectivity reject", newPDNConnectivityReject, msg.ESMPDNConnectivityReject, "PDN connectivity reject"},
		{"PDN disconnect request", newPDNDisconnectRequest, msg.ESMPDNDisconnectRequest, "PDN disconnect request"},
		{"PDN disconnect reject", newPDNDisconnectReject, msg.ESMPDNDisconnectReject, "PDN disconnect reject"},
		{"Bearer resource allocation request", newBearerResourceAllocRequest, msg.ESMBearerResourceAllocationRequest, "Bearer resource allocation request"},
		{"Bearer resource allocation reject", newBearerResourceAllocReject, msg.ESMBearerResourceAllocationReject, "Bearer resource allocation reject"},
		{"Bearer resource modification request", newBearerResourceModRequest, msg.ESMBearerResourceModificationRequest, "Bearer resource modification request"},
		{"Bearer resource modification reject", newBearerResourceModReject, msg.ESMBearerResourceModificationReject, "Bearer resource modification reject"},
		{"ESM information request", newEsmInfoRequest, msg.ESMInformationRequest, "ESM information request"},
		{"ESM information response", newEsmInfoResponse, msg.ESMInformationResponse, "ESM information response"},
		{"ESM notification", newNotification, msg.ESMNotification, "Notification"},
		{"ESM dummy", newDummy, msg.ESMDummy, "ESM dummy message"},
		{"ESM status", newEsmStatus, msg.ESMStatus, "ESM status"},
		{"Remote UE report", newRemoteUEReport, msg.ESMRemoteUEReport, "Remote UE report"},
		{"Remote UE report response", newRemoteUEReportResponse, msg.ESMRemoteUEReportResponse, "Remote UE report response"},
		{"ESM data transport", newEsmDataTransport, msg.ESMDataTransport, "ESM data transport"},
	}
}

// ---------------------------------------------------------------------------
// TestParse — table-driven test exercising msg.Parse for all 59 types.
// ---------------------------------------------------------------------------

func TestParse(t *testing.T) {
	t.Parallel()

	all := append(allEMMEntries(), allESMEntries()...)

	for _, tc := range all {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			src := tc.newMsg()
			buf := make([]byte, src.MarshalLen())
			if err := src.MarshalTo(buf); err != nil {
				t.Fatalf("MarshalTo: %v", err)
			}

			got, err := msg.Parse(buf)
			if err != nil {
				t.Fatalf("Parse(%d bytes): %v", len(buf), err)
			}
			if got == nil {
				t.Fatal("Parse returned nil")
			}

			wantTypeName := tc.typeName
			gotTypeName := got.MessageTypeName()
			if gotTypeName != wantTypeName {
				t.Errorf("MessageTypeName = %q, want %q", gotTypeName, wantTypeName)
			}

			if got.MessageType() != tc.msgType {
				t.Errorf("MessageType = %#02x, want %#02x", got.MessageType(), tc.msgType)
			}
		})
	}

	// --- Error and fallback cases ---

	t.Run("errors_and_fallbacks", func(t *testing.T) {
		t.Parallel()

		testParseErrorCases(t)
	})
}

// testParseErrorCases tests Parse's error handling and fallback behaviour.
func testParseErrorCases(t *testing.T) {
	cases := []struct {
		name string
		data []byte
		want error
	}{
		{"empty input", nil, io.ErrUnexpectedEOF},
		{"single byte EMM header", []byte{0x07, 0x41}, io.ErrUnexpectedEOF},
		{"invalid SHT (PD=EMM, SHT=0xE)", []byte{0xE7, 0x41}, ie.ErrInvalidValue},
		{"plain EMM type byte present but payload empty", []byte{0x07, 0x41}, io.ErrUnexpectedEOF},
		{"protected EMM truncated MAC", []byte{0x17, 0x5E, 0x11, 0x22}, ie.ErrInvalidValue},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := msg.Parse(tc.data)
			if tc.want == nil {
				return
			}

			if err == nil {
				t.Fatalf("Parse(%v): expected error %T, got nil", tc.data, tc.want)
			}

			if !errors.Is(err, tc.want) {
				t.Errorf("Parse(%v): error = %v (%T), want %v (%T)",
					tc.data, err, err, tc.want, tc.want)
			}

			if got != nil {
				t.Errorf("Parse(%v): non-nil message on error: %T", tc.data, got)
			}
		})
	}

	t.Run("unknown EMM type -> Generic", func(t *testing.T) {
		t.Parallel()
		got, err := msg.Parse([]byte{0x07, 0xFF})
		if err != nil {
			t.Fatalf("Parse unknown EMM: %v", err)
		}
		if _, ok := got.(*msg.Generic); !ok {
			t.Errorf("expected *Generic, got %T", got)
		}
	})

	t.Run("unknown ESM type -> Generic", func(t *testing.T) {
		t.Parallel()
		got, err := msg.Parse([]byte{0x21, 0x01, 0xFF})
		if err != nil {
			t.Fatalf("Parse unknown ESM: %v", err)
		}
		if _, ok := got.(*msg.Generic); !ok {
			t.Errorf("expected *Generic, got %T", got)
		}
	})

	t.Run("non-EPS PD -> Generic", func(t *testing.T) {
		t.Parallel()
		got, err := msg.Parse([]byte{0x3F, 0x01})
		if err != nil {
			t.Fatalf("Parse non-EPS PD: %v", err)
		}
		if _, ok := got.(*msg.Generic); !ok {
			t.Errorf("expected *Generic, got %T", got)
		}
	})

}
func TestFuzzParseSeedHelper(t *testing.T) {
	t.Parallel()
	_ = &ie.IE{Kind: ie.Kind(0x7E), Payload: []byte{0xD0, 0x01}, Format: ie.IEFormatFour}
	_ = &ie.IE{Kind: ie.Kind(0x17), Payload: []byte{0x1}, Format: ie.IEFormatFour}
	_ = &ie.IE{Kind: ie.Kind(0x00), Payload: []byte{0x10}, Format: ie.IEFormatFour}
	_ = &ie.IE{Kind: ie.Kind(0x7E), Payload: []byte{}, Format: ie.IEFormatFour}
	_ = &ie.IE{Kind: ie.Kind(0x58), Payload: []byte{0x01, 0x01}, Format: ie.IEFormatFour}
	t.Log("fuzz seed IE constructors compile and don't panic")
}
