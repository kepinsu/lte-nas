package ie

import (
	"fmt"
	"io"
)

// EPSQualityOfServiceFields is the structured value of the EPS quality of
// service IE (TS 24.301 clause 9.9.4.3). QCI is the QoS class identifier. The
// bit-rate fields are optional as one group: either all four are nil or all
// four are non-nil. Bit rates use the same bit/s representation as
// QualityOfServiceFields and select their extended encodings automatically.
type EPSQualityOfServiceFields struct {
	QCI                          uint8
	MaximumBitRateForUplink      *QualityOfServiceBitRate
	MaximumBitRateForDownlink    *QualityOfServiceBitRate
	GuaranteedBitRateForUplink   *QualityOfServiceBitRate
	GuaranteedBitRateForDownlink *QualityOfServiceBitRate
}

// NewEPSQualityOfServiceWithFields creates a tagged EPS quality of service IE
// from fields.
func NewEPSQualityOfServiceWithFields(fields *EPSQualityOfServiceFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewEPSQualityOfService(payload), nil
}

// NewPositionalEPSQualityOfServiceWithFields creates a positional EPS quality
// of service IE from fields.
func NewPositionalEPSQualityOfServiceWithFields(fields *EPSQualityOfServiceFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalEPSQualityOfService(payload), nil
}

// EPSQualityOfService returns the IE as structured fields.
func (i *IE) EPSQualityOfService() (*EPSQualityOfServiceFields, error) {
	if err := i.requireDefinition(EPSQualityOfServiceDefinition); err != nil {
		return nil, err
	}
	return ParseEPSQualityOfServiceFields(i.Payload)
}

// ParseEPSQualityOfServiceFields decodes an EPS quality of service value
// without its IEI and length octet.
func ParseEPSQualityOfServiceFields(payload []byte) (*EPSQualityOfServiceFields, error) {
	fields := &EPSQualityOfServiceFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies the QCI and all-or-none bit-rate group.
func (f *EPSQualityOfServiceFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: EPS quality of service fields are nil", ErrInvalidValue)
	}
	if f.QCI == 0xff {
		return fmt.Errorf("%w: QCI 255 is reserved", ErrInvalidValue)
	}
	if !f.hasBitRates() {
		if f.MaximumBitRateForUplink != nil || f.MaximumBitRateForDownlink != nil ||
			f.GuaranteedBitRateForUplink != nil || f.GuaranteedBitRateForDownlink != nil {
			return fmt.Errorf("%w: EPS quality of service bit rates must be present as a group", ErrInvalidValue)
		}
		return nil
	}
	_, err := f.encodedBitRates()
	return err
}

// Marshal encodes fields into an EPS quality of service value without its IEI
// and length octet.
func (f *EPSQualityOfServiceFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	if !f.hasBitRates() {
		return []byte{f.QCI}, nil
	}
	rates, err := f.encodedBitRates()
	if err != nil {
		return nil, err
	}
	payload := make([]byte, epsQualityOfServiceBitRateLength(rates)+1)
	payload[0] = f.QCI
	payload[1] = rates.uplinkMaximum.base
	payload[2] = rates.downlinkMaximum.base
	payload[3] = rates.uplinkGuaranteed.base
	payload[4] = rates.downlinkGuaranteed.base
	if len(payload) >= 9 {
		payload[5] = rates.uplinkMaximum.extended
		payload[6] = rates.downlinkMaximum.extended
		payload[7] = rates.uplinkGuaranteed.extended
		payload[8] = rates.downlinkGuaranteed.extended
	}
	if len(payload) == 13 {
		payload[9] = rates.uplinkMaximum.extended2
		payload[10] = rates.downlinkMaximum.extended2
		payload[11] = rates.uplinkGuaranteed.extended2
		payload[12] = rates.downlinkGuaranteed.extended2
	}
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *EPSQualityOfServiceFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the EPS quality of service value payload.
func (f *EPSQualityOfServiceFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	if !f.hasBitRates() {
		return 1
	}
	return f.estimatedBitRateLength() + 1
}

// UnmarshalBinary decodes an EPS quality of service value without its IEI and
// length octet. It leaves f unchanged when the payload is invalid.
func (f *EPSQualityOfServiceFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: EPS quality of service fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 && len(payload) != 5 && len(payload) != 9 && len(payload) != 13 {
		return fmt.Errorf("%w: EPS quality of service has %d value octets, want 1, 5, 9, or 13", ErrInvalidLength, len(payload))
	}
	decoded := EPSQualityOfServiceFields{QCI: payload[0]}
	if len(payload) == 1 {
		if err := decoded.Validate(); err != nil {
			return err
		}
		*f = decoded
		return nil
	}

	var uplinkMaximumExtended, downlinkMaximumExtended, uplinkGuaranteedExtended, downlinkGuaranteedExtended byte
	var uplinkMaximumExtended2, downlinkMaximumExtended2, uplinkGuaranteedExtended2, downlinkGuaranteedExtended2 byte
	if len(payload) >= 9 {
		uplinkMaximumExtended = payload[5]
		downlinkMaximumExtended = payload[6]
		uplinkGuaranteedExtended = payload[7]
		downlinkGuaranteedExtended = payload[8]
	}
	if len(payload) == 13 {
		uplinkMaximumExtended2 = payload[9]
		downlinkMaximumExtended2 = payload[10]
		uplinkGuaranteedExtended2 = payload[11]
		downlinkGuaranteedExtended2 = payload[12]
		if uplinkMaximumExtended2 != 0 && uplinkMaximumExtended != 0xfa ||
			downlinkMaximumExtended2 != 0 && downlinkMaximumExtended != 0xfa ||
			uplinkGuaranteedExtended2 != 0 && uplinkGuaranteedExtended != 0xfa ||
			downlinkGuaranteedExtended2 != 0 && downlinkGuaranteedExtended != 0xfa {
			return fmt.Errorf("%w: EPS quality of service extended-2 values require extended code 0xfa", ErrInvalidValue)
		}
	}

	uplinkMaximum, err := decodeQualityOfServiceBitRate(payload[1], uplinkMaximumExtended, uplinkMaximumExtended2)
	if err != nil {
		return err
	}
	downlinkMaximum, err := decodeQualityOfServiceBitRate(payload[2], downlinkMaximumExtended, downlinkMaximumExtended2)
	if err != nil {
		return err
	}
	uplinkGuaranteed, err := decodeQualityOfServiceBitRate(payload[3], uplinkGuaranteedExtended, uplinkGuaranteedExtended2)
	if err != nil {
		return err
	}
	downlinkGuaranteed, err := decodeQualityOfServiceBitRate(payload[4], downlinkGuaranteedExtended, downlinkGuaranteedExtended2)
	if err != nil {
		return err
	}
	decoded.MaximumBitRateForUplink = &uplinkMaximum
	decoded.MaximumBitRateForDownlink = &downlinkMaximum
	decoded.GuaranteedBitRateForUplink = &uplinkGuaranteed
	decoded.GuaranteedBitRateForDownlink = &downlinkGuaranteed
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}

func (f *EPSQualityOfServiceFields) hasBitRates() bool {
	return f.MaximumBitRateForUplink != nil && f.MaximumBitRateForDownlink != nil &&
		f.GuaranteedBitRateForUplink != nil && f.GuaranteedBitRateForDownlink != nil
}

func (f *EPSQualityOfServiceFields) encodedBitRates() (qualityOfServiceEncodedBitRates, error) {
	uplinkMaximum, err := encodeQualityOfServiceBitRate(*f.MaximumBitRateForUplink)
	if err != nil {
		return qualityOfServiceEncodedBitRates{}, fmt.Errorf("maximum bit rate for uplink: %w", err)
	}
	downlinkMaximum, err := encodeQualityOfServiceBitRate(*f.MaximumBitRateForDownlink)
	if err != nil {
		return qualityOfServiceEncodedBitRates{}, fmt.Errorf("maximum bit rate for downlink: %w", err)
	}
	uplinkGuaranteed, err := encodeQualityOfServiceBitRate(*f.GuaranteedBitRateForUplink)
	if err != nil {
		return qualityOfServiceEncodedBitRates{}, fmt.Errorf("guaranteed bit rate for uplink: %w", err)
	}
	downlinkGuaranteed, err := encodeQualityOfServiceBitRate(*f.GuaranteedBitRateForDownlink)
	if err != nil {
		return qualityOfServiceEncodedBitRates{}, fmt.Errorf("guaranteed bit rate for downlink: %w", err)
	}
	return qualityOfServiceEncodedBitRates{
		uplinkMaximum:      uplinkMaximum,
		downlinkMaximum:    downlinkMaximum,
		uplinkGuaranteed:   uplinkGuaranteed,
		downlinkGuaranteed: downlinkGuaranteed,
	}, nil
}

func (f *EPSQualityOfServiceFields) estimatedBitRateLength() int {
	rates := qualityOfServiceEncodedBitRates{
		uplinkMaximum:      estimateQualityOfServiceBitRate(*f.MaximumBitRateForUplink),
		downlinkMaximum:    estimateQualityOfServiceBitRate(*f.MaximumBitRateForDownlink),
		uplinkGuaranteed:   estimateQualityOfServiceBitRate(*f.GuaranteedBitRateForUplink),
		downlinkGuaranteed: estimateQualityOfServiceBitRate(*f.GuaranteedBitRateForDownlink),
	}
	return epsQualityOfServiceBitRateLength(rates)
}

func epsQualityOfServiceBitRateLength(rates qualityOfServiceEncodedBitRates) int {
	if rates.uplinkMaximum.extended2 != 0 || rates.downlinkMaximum.extended2 != 0 ||
		rates.uplinkGuaranteed.extended2 != 0 || rates.downlinkGuaranteed.extended2 != 0 {
		return 12
	}
	if rates.uplinkMaximum.extended != 0 || rates.downlinkMaximum.extended != 0 ||
		rates.uplinkGuaranteed.extended != 0 || rates.downlinkGuaranteed.extended != 0 {
		return 8
	}
	return 4
}
