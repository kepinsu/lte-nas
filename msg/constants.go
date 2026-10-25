package msg

// Message type constants from TS 24.301 Tables 9.8.1 (EMM) and 9.8.2 (ESM).
// EMM values have top two bits 01 (xx01xxxx); ESM values have top two bits 11 (xx11xxxx).

// ---------------------------------------------------------------------------
// EMM message types — TS 24.301 Table 9.8.1
// ---------------------------------------------------------------------------

const (
	EMMAttachRequest               uint8 = 0x41 // TS 24.301 Table 9.8.1
	EMMAttachAccept                uint8 = 0x42
	EMMAttachComplete              uint8 = 0x43
	EMMAttachReject                uint8 = 0x44
	EMMDetachRequest               uint8 = 0x45
	EMMDetachAccept                uint8 = 0x46
	EMMTrackingAreaUpdateRequest   uint8 = 0x48
	EMMTrackingAreaUpdateAccept    uint8 = 0x49
	EMMTrackingAreaUpdateComplete  uint8 = 0x4A
	EMMTrackingAreaUpdateReject    uint8 = 0x4B
	EMMExtendedServiceRequest      uint8 = 0x4C
	EMMControlPlaneServiceRequest  uint8 = 0x4D
	EMMServiceReject               uint8 = 0x4E
	EMMServiceAccept               uint8 = 0x4F
	EMMGUTIReallocationCommand     uint8 = 0x50
	EMMGUTIReallocationComplete    uint8 = 0x51
	EMMAuthenticationRequest       uint8 = 0x52
	EMMAuthenticationResponse      uint8 = 0x53
	EMMAuthenticationReject        uint8 = 0x54
	EMMAuthenticationFailure       uint8 = 0x5C
	EMMIdentityRequest             uint8 = 0x55
	EMMIdentityResponse            uint8 = 0x56
	EMMSecurityModeCommand         uint8 = 0x5D
	EMMSecurityModeComplete        uint8 = 0x5E
	EMMSecurityModeReject          uint8 = 0x5F
	EMMStatus                      uint8 = 0x60
	EMMInformation                 uint8 = 0x61
	EMMDownlinkNASTransport        uint8 = 0x62
	EMMUplinkNASTransport          uint8 = 0x63
	EMMCSServiceNotification       uint8 = 0x64
	EMMDownlinkGenericNASTransport uint8 = 0x68
	EMMUplinkGenericNASTransport   uint8 = 0x69
)

// ---------------------------------------------------------------------------
// ESM message types — TS 24.301 Table 9.8.2
// ---------------------------------------------------------------------------

const (
	ESMActivateDefaultEPSBearerContextRequest   uint8 = 0xC1
	ESMActivateDefaultEPSBearerContextAccept    uint8 = 0xC2
	ESMActivateDefaultEPSBearerContextReject    uint8 = 0xC3
	ESMActivateDedicatedEPSBearerContextRequest uint8 = 0xC5
	ESMActivateDedicatedEPSBearerContextAccept  uint8 = 0xC6
	ESMActivateDedicatedEPSBearerContextReject  uint8 = 0xC7
	ESMModifyEPSBearerContextRequest            uint8 = 0xC9
	ESMModifyEPSBearerContextAccept             uint8 = 0xCA
	ESMModifyEPSBearerContextReject             uint8 = 0xCB
	ESMDeactivateEPSBearerContextRequest        uint8 = 0xCD
	ESMDeactivateEPSBearerContextAccept         uint8 = 0xCE
	ESMPDNConnectivityRequest                   uint8 = 0xD0
	ESMPDNConnectivityReject                    uint8 = 0xD1
	ESMPDNDisconnectRequest                     uint8 = 0xD2
	ESMPDNDisconnectReject                      uint8 = 0xD3
	ESMBearerResourceAllocationRequest          uint8 = 0xD4
	ESMBearerResourceAllocationReject           uint8 = 0xD5
	ESMBearerResourceModificationRequest        uint8 = 0xD6
	ESMBearerResourceModificationReject         uint8 = 0xD7
	ESMInformationRequest                       uint8 = 0xD9
	ESMInformationResponse                      uint8 = 0xDA
	ESMNotification                             uint8 = 0xDB
	ESMDummy                                    uint8 = 0xDC
	ESMStatus                                   uint8 = 0xE8
	ESMRemoteUEReport                           uint8 = 0xE9
	ESMRemoteUEReportResponse                   uint8 = 0xEA
	ESMDataTransport                            uint8 = 0xEB
)
