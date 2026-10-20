package ie

// EmmCause values from TS 24.301 Table 9.9.3.9.1.
type MmCause uint8

const (
	/// #2 IMSI unknown in HSS.
	ImsiUnknownInHssMmCause MmCause = 0x02
	/// #3 Illegal UE.
	IllegalUeMmCause MmCause = 0x03
	/// #5 IMEI not accepted.
	ImeiNotAcceptedMmCause MmCause = 0x05
	/// #6 Illegal ME.
	IllegalMeMmCause MmCause = 0x06
	/// #7 EPS services not allowed.
	EpsServicesNotAllowedMmCause MmCause = 0x07
	/// #8 EPS services and non-EPS services not allowed.
	EpsServicesAndNonEpsServicesNotAllowedMmCause MmCause = 0x08
	/// #9 UE identity cannot be derived by the network.
	UeIdentityCannotBeDerivedByTheNetworkMmCause MmCause = 0x09
	/// #10 Implicitly detached.
	ImplicitlyDetachedMmCause MmCause = 0x0A
	/// #11 PLMN not allowed.
	PlmnNotAllowedMmCause MmCause = 0x0B
	/// #12 Tracking Area not allowed.
	TrackingAreaNotAllowedMmCause MmCause = 0x0C
	/// #13 Roaming not allowed in this tracking area.
	RoamingNotAllowedInThisTrackingAreaMmCause MmCause = 0x0D
	/// #14 EPS services not allowed in this PLMN.
	EpsServicesNotAllowedInThisPlmnMmCause MmCause = 0x0E
	/// #15 No Suitable Cells In tracking area.
	NoSuitableCellsInTrackingAreaMmCause MmCause = 0x0F
	/// #16 MSC temporarily not reachable.
	MscTemporarilyNotReachableMmCause MmCause = 0x10
	/// #17 Network failure.
	NetworkFailureMmCause MmCause = 0x11
	/// #18 CS domain not available.
	CsDomainNotAvailableMmCause MmCause = 0x12
	/// #19 ESM failure.
	EsmFailureMmCause MmCause = 0x13
	/// #20 MAC failure.
	MacFailureMmCause MmCause = 0x14
	/// #21 Synch failure.
	SynchFailureMmCause MmCause = 0x15
	/// #22 Congestion.
	CongestionMmCause MmCause = 0x16
	/// #23 UE security capabilities mismatch.
	UeSecurityCapabilitiesMismatchMmCause MmCause = 0x17
	/// #24 Security mode rejected unspecified.
	SecurityModeRejectedUnspecifiedMmCause MmCause = 0x18
	/// #25 Not authorized for this CSG.
	NotAuthorizedForThisCsgMmCause MmCause = 0x19
	/// #26 Non-EPS authentication unacceptable.
	NonEpsAuthenticationUnacceptableMmCause MmCause = 0x1A
	/// #31 Redirection to 5GCN required.
	RedirectionToFiveGcnRequiredMmCause MmCause = 0x1F
	/// #35 Requested service option not authorized in this PLMN.
	RequestedServiceOptionNotAuthorizedInThisPlmnMmCause MmCause = 0x23
	/// #36 IAB-node operation not authorized.
	IabNodeOperationNotAuthorizedMmCause MmCause = 0x24
	/// #39 CS service temporarily not available.
	CsServiceTemporarilyNotAvailableMmCause MmCause = 0x27
	/// #40 No EPS bearer context activated.
	NoEpsBearerContextActivatedMmCause MmCause = 0x28
	/// #42 Severe network failure.
	SevereNetworkFailureMmCause MmCause = 0x2A
	/// #78 PLMN not allowed to operate at the present UE location.
	PlmnNotAllowedToOperateAtThePresentUeLocationMmCause MmCause = 0x4E
	/// #80 Disaster roaming for the determined PLMN with disaster condition not allowed.
	DisasterRoamingForTheDeterminedPlmnWithDisasterConditionNotAllowedMmCause MmCause = 0x50
	/// #83 Procedure cannot be completed due to unavailable feeder link while MME is operating in S&F mode.
	ProcedureCannotBeCompletedDueToUnavailableFeederLinkWhileMmeIsOperatingInSAndFModeMmCause MmCause = 0x53
	/// #95 Semantically incorrect message.
	SemanticallyIncorrectMessageMmCause MmCause = 0x5F
	/// #96 Invalid mandatory information.
	InvalidMandatoryInformationMmCause MmCause = 0x60
	/// #97 Message type non-existent or not implemented.
	MessageTypeNonExistentOrNotImplementedMmCause MmCause = 0x61
	/// #98 Message type not compatible with the protocol state.
	MessageTypeNotCompatibleWithTheProtocolStateMmCause MmCause = 0x62
	/// #99 Information element non-existent or not implemented.
	InformationElementNonExistentOrNotImplementedMmCause MmCause = 0x63
	/// #100 Conditional IE error.
	ConditionalIeErrorMmCause MmCause = 0x64
	/// #101 Message not compatible with the protocol state.
	MessageNotCompatibleWithTheProtocolStateMmCause MmCause = 0x65
	/// #111 Protocol error unspecified.
	ProtocolErrorUnspecifiedMmCause MmCause = 0x6F
)
