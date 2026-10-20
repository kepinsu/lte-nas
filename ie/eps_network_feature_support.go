package ie

import (
	"fmt"
	"io"
)

// CSLocationServicesSupport identifies the network support status of location
// services through the CS domain.
type CSLocationServicesSupport uint8

// CS location-services support values from TS 24.301 clause 9.9.3.12A.
const (
	CSLocationServicesSupportUnspecified CSLocationServicesSupport = iota
	CSLocationServicesSupportSupported
	CSLocationServicesSupportNotSupported
)

// EPSNetworkFeatureSupportFields is the structured value of the EPS network
// feature support IE (TS 24.301 clause 9.9.3.12A). Optional octets are emitted
// only when at least one feature in that octet is set.
type EPSNetworkFeatureSupportFields struct {
	IMSVoiceOverPSSession                   bool
	EmergencyBearerServices                 bool
	EPCLocationServices                     bool
	CSLocationServices                      CSLocationServicesSupport
	ExtendedServiceRequestForPacketServices bool
	EMMRegisteredWithoutPDN                 bool
	ControlPlaneCIoT                        bool
	UserPlaneCIoT                           bool
	S1UDataTransfer                         bool
	HeaderCompressionForControlPlaneCIoT    bool
	ExtendedProtocolConfigurationOptions    bool
	EnhancedCoverageRestricted              bool
	DualConnectivityWithNRRestricted        bool
	InterworkingWithoutN26                  bool
	FifteenEPSBearerContexts                bool
	NASSignallingConnectionRelease          bool
	PagingIndicationForVoiceServices        bool
	RejectPagingRequest                     bool
	PagingRestriction                       bool
	PagingTimingCollisionControl            bool
	EnhancedDiscontinuousCoverage           bool
}

// NewEPSNetworkFeatureSupportWithFields creates a tagged EPS network feature
// support IE from fields.
func NewEPSNetworkFeatureSupportWithFields(fields *EPSNetworkFeatureSupportFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewEPSNetworkFeatureSupport(payload), nil
}

// NewPositionalEPSNetworkFeatureSupportWithFields creates a positional EPS
// network feature support IE from fields.
func NewPositionalEPSNetworkFeatureSupportWithFields(fields *EPSNetworkFeatureSupportFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalEPSNetworkFeatureSupport(payload), nil
}

// EPSNetworkFeatureSupport returns the IE as structured fields.
func (i *IE) EPSNetworkFeatureSupport() (*EPSNetworkFeatureSupportFields, error) {
	if err := i.requireDefinition(EPSNetworkFeatureSupportDefinition); err != nil {
		return nil, err
	}
	return ParseEPSNetworkFeatureSupportFields(i.Payload)
}

// ParseEPSNetworkFeatureSupportFields decodes an EPS network feature support
// value without its IEI and length octet.
func ParseEPSNetworkFeatureSupportFields(payload []byte) (*EPSNetworkFeatureSupportFields, error) {
	fields := &EPSNetworkFeatureSupportFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that CSLocationServices is a defined value.
func (f *EPSNetworkFeatureSupportFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: EPS network feature support fields are nil", ErrInvalidValue)
	}
	if f.CSLocationServices > CSLocationServicesSupportNotSupported {
		return fmt.Errorf("%w: CS location-services support %d is reserved", ErrInvalidValue, f.CSLocationServices)
	}
	return nil
}

// Marshal encodes fields into an EPS network feature support value without its
// IEI and length octet.
func (f *EPSNetworkFeatureSupportFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, f.MarshalLen())
	payload[0] = f.firstOctet()
	if len(payload) >= 2 {
		payload[1] = f.secondOctet()
	}
	if len(payload) == 3 {
		payload[2] = f.thirdOctet()
	}
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *EPSNetworkFeatureSupportFields) MarshalTo(b []byte) error {
	payload, err := f.Marshal()
	if err != nil {
		return err
	}
	if len(b) < len(payload) {
		return io.ErrShortBuffer
	}
	copy(b, payload)
	return nil
}

// MarshalLen returns the size of the EPS network feature support value payload
// produced by Marshal. It returns zero for a nil receiver.
func (f *EPSNetworkFeatureSupportFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	if f.thirdOctet() != 0 {
		return 3
	}
	if f.secondOctet() != 0 {
		return 2
	}
	return 1
}

// UnmarshalBinary decodes an EPS network feature support value without its IEI
// and length octet. It leaves f unchanged when payload is invalid.
func (f *EPSNetworkFeatureSupportFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: EPS network feature support fields are nil", ErrInvalidValue)
	}
	if len(payload) < 1 || len(payload) > 3 {
		return fmt.Errorf("%w: EPS network feature support has %d value octets, want 1..3", ErrInvalidLength, len(payload))
	}
	decoded := EPSNetworkFeatureSupportFields{
		IMSVoiceOverPSSession:                   payload[0]&0x01 != 0,
		EmergencyBearerServices:                 payload[0]&0x02 != 0,
		EPCLocationServices:                     payload[0]&0x04 != 0,
		CSLocationServices:                      CSLocationServicesSupport((payload[0] >> 3) & 0x03),
		ExtendedServiceRequestForPacketServices: payload[0]&0x20 != 0,
		EMMRegisteredWithoutPDN:                 payload[0]&0x40 != 0,
		ControlPlaneCIoT:                        payload[0]&0x80 != 0,
	}
	if err := decoded.Validate(); err != nil {
		return err
	}
	if len(payload) >= 2 {
		decoded.UserPlaneCIoT = payload[1]&0x01 != 0
		decoded.S1UDataTransfer = payload[1]&0x02 != 0
		decoded.HeaderCompressionForControlPlaneCIoT = payload[1]&0x04 != 0
		decoded.ExtendedProtocolConfigurationOptions = payload[1]&0x08 != 0
		decoded.EnhancedCoverageRestricted = payload[1]&0x10 != 0
		decoded.DualConnectivityWithNRRestricted = payload[1]&0x20 != 0
		decoded.InterworkingWithoutN26 = payload[1]&0x40 != 0
		decoded.FifteenEPSBearerContexts = payload[1]&0x80 != 0
	}
	if len(payload) == 3 {
		if payload[2]&0xc0 != 0 {
			return fmt.Errorf("%w: EPS network feature support third-octet spare bits are nonzero", ErrInvalidValue)
		}
		decoded.NASSignallingConnectionRelease = payload[2]&0x01 != 0
		decoded.PagingIndicationForVoiceServices = payload[2]&0x02 != 0
		decoded.RejectPagingRequest = payload[2]&0x04 != 0
		decoded.PagingRestriction = payload[2]&0x08 != 0
		decoded.PagingTimingCollisionControl = payload[2]&0x10 != 0
		decoded.EnhancedDiscontinuousCoverage = payload[2]&0x20 != 0
	}
	*f = decoded
	return nil
}

func (f *EPSNetworkFeatureSupportFields) firstOctet() byte {
	var encoded byte
	if f.IMSVoiceOverPSSession {
		encoded |= 0x01
	}
	if f.EmergencyBearerServices {
		encoded |= 0x02
	}
	if f.EPCLocationServices {
		encoded |= 0x04
	}
	encoded |= byte(f.CSLocationServices) << 3
	if f.ExtendedServiceRequestForPacketServices {
		encoded |= 0x20
	}
	if f.EMMRegisteredWithoutPDN {
		encoded |= 0x40
	}
	if f.ControlPlaneCIoT {
		encoded |= 0x80
	}
	return encoded
}

func (f *EPSNetworkFeatureSupportFields) secondOctet() byte {
	var encoded byte
	if f.UserPlaneCIoT {
		encoded |= 0x01
	}
	if f.S1UDataTransfer {
		encoded |= 0x02
	}
	if f.HeaderCompressionForControlPlaneCIoT {
		encoded |= 0x04
	}
	if f.ExtendedProtocolConfigurationOptions {
		encoded |= 0x08
	}
	if f.EnhancedCoverageRestricted {
		encoded |= 0x10
	}
	if f.DualConnectivityWithNRRestricted {
		encoded |= 0x20
	}
	if f.InterworkingWithoutN26 {
		encoded |= 0x40
	}
	if f.FifteenEPSBearerContexts {
		encoded |= 0x80
	}
	return encoded
}

func (f *EPSNetworkFeatureSupportFields) thirdOctet() byte {
	var encoded byte
	if f.NASSignallingConnectionRelease {
		encoded |= 0x01
	}
	if f.PagingIndicationForVoiceServices {
		encoded |= 0x02
	}
	if f.RejectPagingRequest {
		encoded |= 0x04
	}
	if f.PagingRestriction {
		encoded |= 0x08
	}
	if f.PagingTimingCollisionControl {
		encoded |= 0x10
	}
	if f.EnhancedDiscontinuousCoverage {
		encoded |= 0x20
	}
	return encoded
}
