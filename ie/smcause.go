package ie

// EsmCause values from TS 24.301 Table 9.9.4.4.1.
type SmCause uint8

const (
	/// #8 Operator Determined Barring.
	OperatorDeterminedBarringSmCause SmCause = 0x08
	/// #26 Insufficient resources.
	InsufficientResourcesSmCause SmCause = 0x1A
	/// #27 Missing or unknown APN.
	MissingOrUnknownApnSmCause SmCause = 0x1B
	/// #28 Unknown PDN type.
	UnknownPdnTypeSmCause SmCause = 0x1C
	/// #29 User authentication or authorization failed.
	UserAuthenticationOrAuthorizationFailedSmCause SmCause = 0x1D
	/// #30 Request rejected by Serving GW or PDN GW.
	RequestRejectedByServingGwOrPdnGwSmCause = 0x1E
	/// #31 Request rejected unspecified.
	RequestRejectedUnspecifiedSmCause SmCause = 0x1F
	/// #32 Service option not supported.
	ServiceOptionNotSupportedSmCause SmCause = 0x20
	/// #33 Requested service option not subscribed.
	RequestedServiceOptionNotSubscribedSmCause SmCause = 0x21
	/// #34 Service option temporarily out of order.
	ServiceOptionTemporarilyOutOfOrderSmCause SmCause = 0x22
	/// #35 PTI already in use.
	PtiAlreadyInUseSmCause SmCause = 0x23
	/// #36 Regular deactivation.
	RegularDeactivationSmCause SmCause = 0x24
	/// #37 EPS QoS not accepted.
	EpsQosNotAcceptedSmCause SmCause = 0x25
	/// #38 Network failure.
	NetworkFailureSmCause SmCause = 0x26
	/// #39 Reactivation requested.
	ReactivationRequestedSmCause SmCause = 0x27
	/// #41 Semantic error in the TFT operation.
	SemanticErrorInTheTftOperationSmCause SmCause = 0x29
	/// #42 Syntactical error in the TFT operation.
	SyntacticalErrorInTheTftOperationSmCause SmCause = 0x2A
	/// #43 Invalid EPS bearer identity.
	InvalidEpsBearerIdentitySmCause SmCause = 0x2B
	/// #44 Semantic errors in packet filter(s).
	SemanticErrorsInPacketFiltersSmCause SmCause = 0x2C
	/// #45 Syntactical errors in packet filter(s).
	SyntacticalErrorsInPacketFiltersSmCause SmCause = 0x2D
	/// #47 PTI mismatch.
	PtiMismatchSmCause SmCause = 0x2F
	/// #49 Last PDN disconnection not allowed.
	LastPdnDisconnectionNotAllowedSmCause SmCause = 0x31
	/// #50 PDN type IPv4 only allowed.
	PdnTypeIpv4OnlyAllowedSmCause SmCause = 0x32
	/// #51 PDN type IPv6 only allowed.
	PdnTypeIpv6OnlyAllowedSmCause SmCause = 0x33
	/// #52 Single address bearers only allowed.
	SingleAddressBearersOnlyAllowedSmCause SmCause = 0x34
	/// #53 ESM information not received.
	EsmInformationNotReceivedSmCause SmCause = 0x35
	/// #54 PDN connection does not exist.
	PdnConnectionDoesNotExistSmCause SmCause = 0x36
	/// #55 Multiple PDN connections for a given APN not allowed.
	MultiplePdnConnectionsForAGivenApnNotAllowedSmCause SmCause = 0x37
	/// #56 Collision with network initiated request.
	CollisionWithNetworkInitiatedRequestSmCause SmCause = 0x38
	/// #57 PDN type IPv4v6 only allowed.
	PdnTypeIpv4v6OnlyAllowedSmCause SmCause = 0x39
	/// #58 PDN type non IP only allowed.
	PdnTypeNonIpOnlyAllowedSmCause SmCause = 0x3A
	/// #59 Unsupported QCI value.
	UnsupportedQciValueSmCause SmCause = 0x3B
	/// #60 Bearer handling not supported.
	BearerHandlingNotSupportedSmCause SmCause = 0x3C
	/// #61 PDN type Ethernet only allowed.
	PdnTypeEthernetOnlyAllowedSmCause SmCause = 0x3D
	/// #65 Maximum number of EPS bearers reached.
	MaximumNumberOfEpsBearersReachedSmCause SmCause = 0x41
	/// #66 Requested APN not supported in current RAT and PLMN combination.
	RequestedApnNotSupportedInCurrentRatAndPlmnCombinationSmCause SmCause = 0x42
	/// #81 Invalid PTI value.
	InvalidPtiValueSmCause SmCause = 0x51
	/// #95 Semantically incorrect message.
	SemanticallyIncorrectMessageSmCause SmCause = 0x5F
	/// #96 Invalid mandatory information.
	InvalidMandatoryInformationSmCause SmCause = 0x60
	/// #97 Message type non-existent or not implemented.
	MessageTypeNonExistentOrNotImplementedSmCause SmCause = 0x61
	/// #98 Message type not compatible with the protocol state.
	MessageTypeNotCompatibleWithTheProtocolStateSmCause SmCause = 0x62
	/// #99 Information element non-existent or not implemented.
	InformationElementNonExistentOrNotImplementedSmCause SmCause = 0x63
	/// #100 Conditional IE error.
	ConditionalIeErrorSmCause SmCause = 0x64
	/// #101 Message not compatible with the protocol state.
	MessageNotCompatibleWithTheProtocolStateSmCause SmCause = 0x65
	/// #111 Protocol error unspecified.
	ProtocolErrorUnspecifiedSmCause SmCause = 0x6F
	/// #112 APN restriction value incompatible with active EPS bearer context.
	ApnRestrictionValueIncompatibleWithActiveEpsBearerContextSmCause SmCause = 0x70
	/// #113 Multiple accesses to a PDN connection not allowed.
	MultipleAccessesToAPdnConnectionNotAllowedSmCause SmCause = 0x71
)
