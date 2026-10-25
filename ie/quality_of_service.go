package ie

import (
	"fmt"
	"io"
	"time"
)

const (
	qualityOfServiceBaseBitRateMaximum      = 8_640_000
	qualityOfServiceExtendedBitRateMaximum  = 256_000_000
	qualityOfServiceExtended2BitRateMaximum = 10_000_000_000
)

// QualityOfServiceDelayClass is the delay class coded in octet 3 of the
// quality of service IE (TS 24.008, clause 10.5.6.5).
type QualityOfServiceDelayClass uint8

// Quality of service delay class values.
const (
	QualityOfServiceDelayClassSubscribed QualityOfServiceDelayClass = iota
	QualityOfServiceDelayClassOne
	QualityOfServiceDelayClassTwo
	QualityOfServiceDelayClassThree
	QualityOfServiceDelayClassFour
	_
	_
	QualityOfServiceDelayClassReserved
)

// QualityOfServiceReliabilityClass is the reliability class coded in octet 3
// of the quality of service IE.
type QualityOfServiceReliabilityClass uint8

// Quality of service reliability class values.
const (
	QualityOfServiceReliabilityClassSubscribed QualityOfServiceReliabilityClass = iota
	QualityOfServiceReliabilityClassUnused
	QualityOfServiceReliabilityClassAcknowledgedLLCAndRLC
	QualityOfServiceReliabilityClassAcknowledgedRLC
	QualityOfServiceReliabilityClassUnacknowledgedRLCProtected
	QualityOfServiceReliabilityClassUnacknowledgedRLCUnprotected
	_
	QualityOfServiceReliabilityClassReserved
)

// QualityOfServicePeakThroughputClass is the peak throughput class coded in
// octet 4 of the quality of service IE.
type QualityOfServicePeakThroughputClass uint8

// Quality of service peak throughput class values, expressed in octets per
// second.
const (
	QualityOfServicePeakThroughputClassSubscribed QualityOfServicePeakThroughputClass = iota
	QualityOfServicePeakThroughputClass1000OctetsPerSecond
	QualityOfServicePeakThroughputClass2000OctetsPerSecond
	QualityOfServicePeakThroughputClass4000OctetsPerSecond
	QualityOfServicePeakThroughputClass8000OctetsPerSecond
	QualityOfServicePeakThroughputClass16000OctetsPerSecond
	QualityOfServicePeakThroughputClass32000OctetsPerSecond
	QualityOfServicePeakThroughputClass64000OctetsPerSecond
	QualityOfServicePeakThroughputClass128000OctetsPerSecond
	QualityOfServicePeakThroughputClass256000OctetsPerSecond
)

// QualityOfServicePrecedenceClass is the precedence class coded in octet 4 of
// the quality of service IE.
type QualityOfServicePrecedenceClass uint8

// Quality of service precedence class values.
const (
	QualityOfServicePrecedenceClassSubscribed QualityOfServicePrecedenceClass = iota
	QualityOfServicePrecedenceClassHigh
	QualityOfServicePrecedenceClassNormal
	QualityOfServicePrecedenceClassLow
	_
	_
	_
	QualityOfServicePrecedenceClassReserved
)

// QualityOfServiceMeanThroughputClass is the mean throughput class coded in
// octet 5 of the quality of service IE.
type QualityOfServiceMeanThroughputClass uint8

// Quality of service mean throughput class values, expressed in octets per
// hour. Best effort is encoded separately because it is not a fixed rate.
const (
	QualityOfServiceMeanThroughputClassSubscribed QualityOfServiceMeanThroughputClass = iota
	QualityOfServiceMeanThroughputClass100OctetsPerHour
	QualityOfServiceMeanThroughputClass200OctetsPerHour
	QualityOfServiceMeanThroughputClass500OctetsPerHour
	QualityOfServiceMeanThroughputClass1000OctetsPerHour
	QualityOfServiceMeanThroughputClass2000OctetsPerHour
	QualityOfServiceMeanThroughputClass5000OctetsPerHour
	QualityOfServiceMeanThroughputClass10000OctetsPerHour
	QualityOfServiceMeanThroughputClass20000OctetsPerHour
	QualityOfServiceMeanThroughputClass50000OctetsPerHour
	QualityOfServiceMeanThroughputClass100000OctetsPerHour
	QualityOfServiceMeanThroughputClass200000OctetsPerHour
	QualityOfServiceMeanThroughputClass500000OctetsPerHour
	QualityOfServiceMeanThroughputClass1000000OctetsPerHour
	QualityOfServiceMeanThroughputClass2000000OctetsPerHour
	QualityOfServiceMeanThroughputClass5000000OctetsPerHour
	QualityOfServiceMeanThroughputClass10000000OctetsPerHour
	QualityOfServiceMeanThroughputClass20000000OctetsPerHour
	QualityOfServiceMeanThroughputClass50000000OctetsPerHour
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	QualityOfServiceMeanThroughputClassReserved
	QualityOfServiceMeanThroughputClassBestEffort
)

// QualityOfServiceTrafficClass is the traffic class coded in octet 6 of the
// quality of service IE.
type QualityOfServiceTrafficClass uint8

// Quality of service traffic class values.
const (
	QualityOfServiceTrafficClassSubscribed QualityOfServiceTrafficClass = iota
	QualityOfServiceTrafficClassConversational
	QualityOfServiceTrafficClassStreaming
	QualityOfServiceTrafficClassInteractive
	QualityOfServiceTrafficClassBackground
	_
	_
	QualityOfServiceTrafficClassReserved
)

// QualityOfServiceDeliveryOrder is the delivery order coded in octet 6 of the
// quality of service IE.
type QualityOfServiceDeliveryOrder uint8

// Quality of service delivery order values.
const (
	QualityOfServiceDeliveryOrderSubscribed QualityOfServiceDeliveryOrder = iota
	QualityOfServiceDeliveryOrderRequired
	QualityOfServiceDeliveryOrderNotRequired
	QualityOfServiceDeliveryOrderReserved
)

// QualityOfServiceDeliveryOfErroneousSDUs is the delivery-of-erroneous-SDUs
// value coded in octet 6 of the quality of service IE.
type QualityOfServiceDeliveryOfErroneousSDUs uint8

// Quality of service delivery-of-erroneous-SDUs values.
const (
	QualityOfServiceDeliveryOfErroneousSDUsSubscribed QualityOfServiceDeliveryOfErroneousSDUs = iota
	QualityOfServiceDeliveryOfErroneousSDUsNoDetection
	QualityOfServiceDeliveryOfErroneousSDUsAllowed
	QualityOfServiceDeliveryOfErroneousSDUsNotAllowed
	_
	_
	_
	QualityOfServiceDeliveryOfErroneousSDUsReserved
)

// QualityOfServiceResidualBitErrorRate is the residual BER class coded in
// octet 10 of the quality of service IE.
type QualityOfServiceResidualBitErrorRate uint8

// Quality of service residual BER class values.
const (
	QualityOfServiceResidualBitErrorRateSubscribed QualityOfServiceResidualBitErrorRate = iota
	QualityOfServiceResidualBitErrorRateFiveIn100
	QualityOfServiceResidualBitErrorRateOneIn100
	QualityOfServiceResidualBitErrorRateFiveIn1000
	QualityOfServiceResidualBitErrorRateFourIn1000
	QualityOfServiceResidualBitErrorRateOneIn1000
	QualityOfServiceResidualBitErrorRateOneIn10000
	QualityOfServiceResidualBitErrorRateOneIn100000
	QualityOfServiceResidualBitErrorRateOneIn1000000
	QualityOfServiceResidualBitErrorRateSixIn100000000
)

// QualityOfServiceSDUErrorRatio is the SDU error ratio class coded in octet 10
// of the quality of service IE.
type QualityOfServiceSDUErrorRatio uint8

// Quality of service SDU error ratio class values.
const (
	QualityOfServiceSDUErrorRatioSubscribed QualityOfServiceSDUErrorRatio = iota
	QualityOfServiceSDUErrorRatioOneIn100
	QualityOfServiceSDUErrorRatioSevenIn1000
	QualityOfServiceSDUErrorRatioOneIn1000
	QualityOfServiceSDUErrorRatioOneIn10000
	QualityOfServiceSDUErrorRatioOneIn100000
	QualityOfServiceSDUErrorRatioOneIn1000000
	QualityOfServiceSDUErrorRatioOneIn10
)

// QualityOfServiceTrafficHandlingPriority is the traffic handling priority
// coded in octet 11 of the quality of service IE.
type QualityOfServiceTrafficHandlingPriority uint8

// Quality of service traffic handling priority values.
const (
	QualityOfServiceTrafficHandlingPrioritySubscribed QualityOfServiceTrafficHandlingPriority = iota
	QualityOfServiceTrafficHandlingPriorityOne
	QualityOfServiceTrafficHandlingPriorityTwo
	QualityOfServiceTrafficHandlingPriorityThree
)

// QualityOfServiceSourceStatisticsDescriptor is the source statistics
// descriptor coded in octet 14 of the quality of service IE.
type QualityOfServiceSourceStatisticsDescriptor uint8

// Quality of service source statistics descriptor values.
const (
	QualityOfServiceSourceStatisticsDescriptorUnknown QualityOfServiceSourceStatisticsDescriptor = iota
	QualityOfServiceSourceStatisticsDescriptorSpeech
)

// QualityOfServiceBitRate is a QoS bit rate. BitsPerSecond is expressed in
// bits per second; it must be exactly representable by the NAS QoS coding.
// Set IsSubscribed to encode the subscribed value instead of a concrete rate.
type QualityOfServiceBitRate struct {
	BitsPerSecond uint64
	IsSubscribed  bool
}

// QualityOfServiceFields is the structured value of the Quality of service IE
// defined by TS 24.301 clause 9.9.4.12 and TS 24.008 clause 10.5.6.5.
//
// MaximumSDUSize is expressed in octets; zero encodes the subscribed maximum
// SDU size. TransferDelay is expressed as a time.Duration; zero encodes the
// subscribed transfer delay. The four bit rates are expressed in bits per
// second. Their optional extended octets are selected automatically when
// marshaling, up to 10 Gbit/s.
type QualityOfServiceFields struct {
	DelayClass                   QualityOfServiceDelayClass
	ReliabilityClass             QualityOfServiceReliabilityClass
	PeakThroughputClass          QualityOfServicePeakThroughputClass
	PrecedenceClass              QualityOfServicePrecedenceClass
	MeanThroughputClass          QualityOfServiceMeanThroughputClass
	TrafficClass                 QualityOfServiceTrafficClass
	DeliveryOrder                QualityOfServiceDeliveryOrder
	DeliveryOfErroneousSDUs      QualityOfServiceDeliveryOfErroneousSDUs
	MaximumSDUSize               uint16
	MaximumBitRateForUplink      QualityOfServiceBitRate
	MaximumBitRateForDownlink    QualityOfServiceBitRate
	ResidualBitErrorRate         QualityOfServiceResidualBitErrorRate
	SDUErrorRatio                QualityOfServiceSDUErrorRatio
	TransferDelay                time.Duration
	TrafficHandlingPriority      QualityOfServiceTrafficHandlingPriority
	GuaranteedBitRateForUplink   QualityOfServiceBitRate
	GuaranteedBitRateForDownlink QualityOfServiceBitRate
	SignallingIndication         bool
	SourceStatisticsDescriptor   QualityOfServiceSourceStatisticsDescriptor
}

// NewQualityOfServiceWithFields creates a tagged Quality of service IE from a
// structured value. NewQualityOfService remains available for callers that
// need to construct the IE from a raw payload.
func NewQualityOfServiceWithFields(fields *QualityOfServiceFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewQualityOfService(payload), nil
}

// NewPositionalQualityOfServiceWithFields creates a positional Quality of
// service IE from a structured value.
func NewPositionalQualityOfServiceWithFields(fields *QualityOfServiceFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalQualityOfService(payload), nil
}

// QualityOfService returns the Quality of service IE as a structured value.
// It accepts both tagged and positional instances created by this package.
func (i *IE) QualityOfService() (*QualityOfServiceFields, error) {
	if err := i.requireDefinition(QualityOfServiceDefinition); err != nil {
		return nil, err
	}
	return ParseQualityOfServiceFields(i.Payload)
}

// ParseQualityOfServiceFields decodes a Quality of service value payload. The
// payload excludes the IEI and length octet.
func ParseQualityOfServiceFields(payload []byte) (*QualityOfServiceFields, error) {
	fields := &QualityOfServiceFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that fields can be represented by the Quality of service
// IE. It does not validate direction-specific requirements because the IE does
// not identify whether it was sent by the MS or by the network.
func (f *QualityOfServiceFields) Validate() error {
	_, err := f.encodedBitRates()
	return err
}

// Marshal encodes fields into a Quality of service value payload, without the
// IEI or length octet.
func (f *QualityOfServiceFields) Marshal() ([]byte, error) {
	encoded, err := f.encodedBitRates()
	if err != nil {
		return nil, err
	}
	b := make([]byte, encoded.length())
	f.marshalTo(b, encoded)
	return b, nil
}

// MarshalTo encodes fields into b. It returns io.ErrShortBuffer without
// modifying b when b is too short.
func (f *QualityOfServiceFields) MarshalTo(b []byte) error {
	encoded, err := f.encodedBitRates()
	if err != nil {
		return err
	}
	if len(b) < encoded.length() {
		return io.ErrShortBuffer
	}
	f.marshalTo(b[:encoded.length()], encoded)
	return nil
}

// MarshalLen returns the size of the Quality of service value payload produced
// by Marshal. It returns zero for a nil receiver.
func (f *QualityOfServiceFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return qualityOfServiceEncodedBitRates{
		uplinkMaximum:      estimateQualityOfServiceBitRate(f.MaximumBitRateForUplink),
		downlinkMaximum:    estimateQualityOfServiceBitRate(f.MaximumBitRateForDownlink),
		uplinkGuaranteed:   estimateQualityOfServiceBitRate(f.GuaranteedBitRateForUplink),
		downlinkGuaranteed: estimateQualityOfServiceBitRate(f.GuaranteedBitRateForDownlink),
	}.length()
}

// UnmarshalBinary decodes a Quality of service value payload, excluding the
// IEI and length octet. It leaves f unchanged if decoding fails.
func (f *QualityOfServiceFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: quality of service fields are nil", ErrInvalidValue)
	}
	if !validQualityOfServicePayloadLength(len(payload)) {
		return fmt.Errorf("%w: quality of service value has %d octets, want 12, 14, 16, 18, or 20", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xc0 != 0 || payload[1]&0x08 != 0 || payload[2]&0xe0 != 0 || payload[11]&0xe0 != 0 {
		return fmt.Errorf("%w: quality of service spare bits are nonzero", ErrInvalidValue)
	}

	maximumSDUSize, err := decodeQualityOfServiceMaximumSDUSize(payload[4])
	if err != nil {
		return err
	}

	extensions := qualityOfServiceBitRateExtensions{}
	if len(payload) >= 14 {
		extensions.downlinkMaximum = payload[12]
		extensions.downlinkGuaranteed = payload[13]
	}
	if len(payload) >= 16 {
		extensions.uplinkMaximum = payload[14]
		extensions.uplinkGuaranteed = payload[15]
	}
	if len(payload) >= 18 {
		extensions.downlinkMaximum2 = payload[16]
		extensions.downlinkGuaranteed2 = payload[17]
	}
	if len(payload) == 20 {
		extensions.uplinkMaximum2 = payload[18]
		extensions.uplinkGuaranteed2 = payload[19]
	}

	uplinkMaximum, err := decodeQualityOfServiceBitRate(payload[5], extensions.uplinkMaximum, extensions.uplinkMaximum2)
	if err != nil {
		return err
	}
	downlinkMaximum, err := decodeQualityOfServiceBitRate(payload[6], extensions.downlinkMaximum, extensions.downlinkMaximum2)
	if err != nil {
		return err
	}
	uplinkGuaranteed, err := decodeQualityOfServiceBitRate(payload[9], extensions.uplinkGuaranteed, extensions.uplinkGuaranteed2)
	if err != nil {
		return err
	}
	downlinkGuaranteed, err := decodeQualityOfServiceBitRate(payload[10], extensions.downlinkGuaranteed, extensions.downlinkGuaranteed2)
	if err != nil {
		return err
	}
	transferDelay, err := decodeQualityOfServiceTransferDelay(payload[8] >> 2)
	if err != nil {
		return err
	}

	decoded := QualityOfServiceFields{
		DelayClass:                   QualityOfServiceDelayClass((payload[0] >> 3) & 0x07),
		ReliabilityClass:             QualityOfServiceReliabilityClass(payload[0] & 0x07),
		PeakThroughputClass:          QualityOfServicePeakThroughputClass(payload[1] >> 4),
		PrecedenceClass:              QualityOfServicePrecedenceClass(payload[1] & 0x07),
		MeanThroughputClass:          QualityOfServiceMeanThroughputClass(payload[2] & 0x1f),
		TrafficClass:                 QualityOfServiceTrafficClass(payload[3] >> 5),
		DeliveryOrder:                QualityOfServiceDeliveryOrder((payload[3] >> 3) & 0x03),
		DeliveryOfErroneousSDUs:      QualityOfServiceDeliveryOfErroneousSDUs(payload[3] & 0x07),
		MaximumSDUSize:               maximumSDUSize,
		MaximumBitRateForUplink:      uplinkMaximum,
		MaximumBitRateForDownlink:    downlinkMaximum,
		ResidualBitErrorRate:         QualityOfServiceResidualBitErrorRate(payload[7] >> 4),
		SDUErrorRatio:                QualityOfServiceSDUErrorRatio(payload[7] & 0x0f),
		TransferDelay:                transferDelay,
		TrafficHandlingPriority:      QualityOfServiceTrafficHandlingPriority(payload[8] & 0x03),
		GuaranteedBitRateForUplink:   uplinkGuaranteed,
		GuaranteedBitRateForDownlink: downlinkGuaranteed,
		SignallingIndication:         payload[11]&0x10 != 0,
		SourceStatisticsDescriptor:   QualityOfServiceSourceStatisticsDescriptor(payload[11] & 0x0f),
	}
	*f = decoded
	return nil
}

func (f *QualityOfServiceFields) marshalTo(b []byte, encoded qualityOfServiceEncodedBitRates) {
	maximumSDUSize, _ := encodeQualityOfServiceMaximumSDUSize(f.MaximumSDUSize)
	transferDelay, _ := encodeQualityOfServiceTransferDelay(f.TransferDelay)

	b[0] = byte(f.DelayClass<<3) | byte(f.ReliabilityClass)
	b[1] = byte(f.PeakThroughputClass<<4) | byte(f.PrecedenceClass)
	b[2] = byte(f.MeanThroughputClass)
	b[3] = byte(f.TrafficClass<<5) | byte(f.DeliveryOrder<<3) | byte(f.DeliveryOfErroneousSDUs)
	b[4] = maximumSDUSize
	b[5] = encoded.uplinkMaximum.base
	b[6] = encoded.downlinkMaximum.base
	b[7] = byte(f.ResidualBitErrorRate<<4) | byte(f.SDUErrorRatio)
	b[8] = transferDelay<<2 | byte(f.TrafficHandlingPriority)
	b[9] = encoded.uplinkGuaranteed.base
	b[10] = encoded.downlinkGuaranteed.base
	b[11] = byte(f.SourceStatisticsDescriptor)
	if f.SignallingIndication {
		b[11] |= 0x10
	}
	if len(b) >= 14 {
		b[12] = encoded.downlinkMaximum.extended
		b[13] = encoded.downlinkGuaranteed.extended
	}
	if len(b) >= 16 {
		b[14] = encoded.uplinkMaximum.extended
		b[15] = encoded.uplinkGuaranteed.extended
	}
	if len(b) >= 18 {
		b[16] = encoded.downlinkMaximum.extended2
		b[17] = encoded.downlinkGuaranteed.extended2
	}
	if len(b) == 20 {
		b[18] = encoded.uplinkMaximum.extended2
		b[19] = encoded.uplinkGuaranteed.extended2
	}
}

func (f *QualityOfServiceFields) encodedBitRates() (qualityOfServiceEncodedBitRates, error) {
	if f == nil {
		return qualityOfServiceEncodedBitRates{}, fmt.Errorf("%w: quality of service fields are nil", ErrInvalidValue)
	}
	if err := validateQualityOfServiceFieldWidths(f); err != nil {
		return qualityOfServiceEncodedBitRates{}, err
	}
	if _, err := encodeQualityOfServiceMaximumSDUSize(f.MaximumSDUSize); err != nil {
		return qualityOfServiceEncodedBitRates{}, err
	}
	if _, err := encodeQualityOfServiceTransferDelay(f.TransferDelay); err != nil {
		return qualityOfServiceEncodedBitRates{}, err
	}

	uplinkMaximum, err := encodeQualityOfServiceBitRate(f.MaximumBitRateForUplink)
	if err != nil {
		return qualityOfServiceEncodedBitRates{}, fmt.Errorf("maximum bit rate for uplink: %w", err)
	}
	downlinkMaximum, err := encodeQualityOfServiceBitRate(f.MaximumBitRateForDownlink)
	if err != nil {
		return qualityOfServiceEncodedBitRates{}, fmt.Errorf("maximum bit rate for downlink: %w", err)
	}
	uplinkGuaranteed, err := encodeQualityOfServiceBitRate(f.GuaranteedBitRateForUplink)
	if err != nil {
		return qualityOfServiceEncodedBitRates{}, fmt.Errorf("guaranteed bit rate for uplink: %w", err)
	}
	downlinkGuaranteed, err := encodeQualityOfServiceBitRate(f.GuaranteedBitRateForDownlink)
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

func validQualityOfServicePayloadLength(length int) bool {
	return length == 12 || length == 14 || length == 16 || length == 18 || length == 20
}

func validateQualityOfServiceFieldWidths(fields *QualityOfServiceFields) error {
	if fields.DelayClass > 0x07 || fields.ReliabilityClass > 0x07 || fields.PeakThroughputClass > 0x0f ||
		fields.PrecedenceClass > 0x07 || fields.MeanThroughputClass > 0x1f || fields.TrafficClass > 0x07 ||
		fields.DeliveryOrder > 0x03 || fields.DeliveryOfErroneousSDUs > 0x07 || fields.ResidualBitErrorRate > 0x0f ||
		fields.SDUErrorRatio > 0x0f || fields.TrafficHandlingPriority > 0x03 || fields.SourceStatisticsDescriptor > 0x0f {
		return fmt.Errorf("%w: quality of service field exceeds its encoded width", ErrInvalidValue)
	}
	return nil
}

func encodeQualityOfServiceMaximumSDUSize(size uint16) (byte, error) {
	switch size {
	case 0:
		return 0, nil
	case 1502:
		return 0x97, nil
	case 1510:
		return 0x98, nil
	case 1520:
		return 0x99, nil
	}
	if size >= 10 && size <= 1500 && size%10 == 0 {
		return byte(size / 10), nil
	}
	return 0, fmt.Errorf("%w: maximum SDU size %d is not representable", ErrInvalidValue, size)
}

func decodeQualityOfServiceMaximumSDUSize(encoded byte) (uint16, error) {
	switch encoded {
	case 0:
		return 0, nil
	case 0x97:
		return 1502, nil
	case 0x98:
		return 1510, nil
	case 0x99:
		return 1520, nil
	}
	if encoded <= 0x96 {
		return uint16(encoded) * 10, nil
	}
	return 0, fmt.Errorf("%w: maximum SDU size code %#02x is reserved", ErrInvalidValue, encoded)
}

func encodeQualityOfServiceTransferDelay(delay time.Duration) (byte, error) {
	if delay == 0 {
		return 0, nil
	}
	if delay < 0 {
		return 0, fmt.Errorf("%w: transfer delay cannot be negative", ErrInvalidValue)
	}
	if delay >= 10*time.Millisecond && delay <= 150*time.Millisecond && delay%(10*time.Millisecond) == 0 {
		return byte(delay / (10 * time.Millisecond)), nil
	}
	if delay >= 200*time.Millisecond && delay <= 950*time.Millisecond && (delay-200*time.Millisecond)%(50*time.Millisecond) == 0 {
		return byte(16 + (delay-200*time.Millisecond)/(50*time.Millisecond)), nil
	}
	if delay >= time.Second && delay <= 4*time.Second && (delay-time.Second)%(100*time.Millisecond) == 0 {
		return byte(32 + (delay-time.Second)/(100*time.Millisecond)), nil
	}
	return 0, fmt.Errorf("%w: transfer delay %s is not representable", ErrInvalidValue, delay)
}

func decodeQualityOfServiceTransferDelay(encoded byte) (time.Duration, error) {
	switch {
	case encoded == 0:
		return 0, nil
	case encoded <= 15:
		return time.Duration(encoded) * 10 * time.Millisecond, nil
	case encoded <= 31:
		return 200*time.Millisecond + time.Duration(encoded-16)*50*time.Millisecond, nil
	case encoded <= 62:
		return time.Second + time.Duration(encoded-32)*100*time.Millisecond, nil
	default:
		return 0, fmt.Errorf("%w: transfer delay code %#02x is reserved", ErrInvalidValue, encoded)
	}
}

type qualityOfServiceEncodedBitRate struct {
	base      byte
	extended  byte
	extended2 byte
}

type qualityOfServiceEncodedBitRates struct {
	uplinkMaximum      qualityOfServiceEncodedBitRate
	downlinkMaximum    qualityOfServiceEncodedBitRate
	uplinkGuaranteed   qualityOfServiceEncodedBitRate
	downlinkGuaranteed qualityOfServiceEncodedBitRate
}

func (r qualityOfServiceEncodedBitRates) length() int {
	if r.uplinkMaximum.extended2 != 0 || r.uplinkGuaranteed.extended2 != 0 {
		return 20
	}
	if r.downlinkMaximum.extended2 != 0 || r.downlinkGuaranteed.extended2 != 0 {
		return 18
	}
	if r.uplinkMaximum.extended != 0 || r.uplinkGuaranteed.extended != 0 {
		return 16
	}
	if r.downlinkMaximum.extended != 0 || r.downlinkGuaranteed.extended != 0 {
		return 14
	}
	return 12
}

type qualityOfServiceBitRateExtensions struct {
	uplinkMaximum       byte
	downlinkMaximum     byte
	uplinkGuaranteed    byte
	downlinkGuaranteed  byte
	uplinkMaximum2      byte
	downlinkMaximum2    byte
	uplinkGuaranteed2   byte
	downlinkGuaranteed2 byte
}

func estimateQualityOfServiceBitRate(rate QualityOfServiceBitRate) qualityOfServiceEncodedBitRate {
	if rate.IsSubscribed || rate.BitsPerSecond <= qualityOfServiceBaseBitRateMaximum {
		return qualityOfServiceEncodedBitRate{}
	}
	if rate.BitsPerSecond <= qualityOfServiceExtendedBitRateMaximum {
		return qualityOfServiceEncodedBitRate{extended: 1}
	}
	return qualityOfServiceEncodedBitRate{extended: 1, extended2: 1}
}

func encodeQualityOfServiceBitRate(rate QualityOfServiceBitRate) (qualityOfServiceEncodedBitRate, error) {
	if rate.IsSubscribed {
		if rate.BitsPerSecond != 0 {
			return qualityOfServiceEncodedBitRate{}, fmt.Errorf("%w: subscribed bit rate has a concrete value", ErrInvalidValue)
		}
		return qualityOfServiceEncodedBitRate{}, nil
	}
	if rate.BitsPerSecond == 0 {
		return qualityOfServiceEncodedBitRate{base: 0xff}, nil
	}
	if rate.BitsPerSecond <= qualityOfServiceBaseBitRateMaximum {
		base, err := encodeQualityOfServiceBaseBitRate(rate.BitsPerSecond)
		return qualityOfServiceEncodedBitRate{base: base}, err
	}
	if rate.BitsPerSecond <= qualityOfServiceExtendedBitRateMaximum {
		extended, err := encodeQualityOfServiceExtendedBitRate(rate.BitsPerSecond)
		return qualityOfServiceEncodedBitRate{base: 0xfe, extended: extended}, err
	}
	if rate.BitsPerSecond <= qualityOfServiceExtended2BitRateMaximum {
		extended2, err := encodeQualityOfServiceExtended2BitRate(rate.BitsPerSecond)
		return qualityOfServiceEncodedBitRate{base: 0xfe, extended: 0xfa, extended2: extended2}, err
	}
	return qualityOfServiceEncodedBitRate{}, fmt.Errorf("%w: bit rate %d bit/s exceeds 10 Gbit/s", ErrInvalidValue, rate.BitsPerSecond)
}

func decodeQualityOfServiceBitRate(base, extended, extended2 byte) (QualityOfServiceBitRate, error) {
	if extended2 != 0 {
		bitsPerSecond, err := decodeQualityOfServiceExtended2BitRate(extended2)
		return QualityOfServiceBitRate{BitsPerSecond: bitsPerSecond}, err
	}
	if extended != 0 {
		bitsPerSecond, err := decodeQualityOfServiceExtendedBitRate(extended)
		return QualityOfServiceBitRate{BitsPerSecond: bitsPerSecond}, err
	}
	bitsPerSecond, subscribed, err := decodeQualityOfServiceBaseBitRate(base)
	return QualityOfServiceBitRate{BitsPerSecond: bitsPerSecond, IsSubscribed: subscribed}, err
}

func encodeQualityOfServiceBaseBitRate(bitsPerSecond uint64) (byte, error) {
	if bitsPerSecond%1_000 != 0 {
		return 0, fmt.Errorf("%w: bit rate %d bit/s is not a whole number of kbit/s", ErrInvalidValue, bitsPerSecond)
	}
	kilobitsPerSecond := bitsPerSecond / 1_000
	switch {
	case kilobitsPerSecond >= 1 && kilobitsPerSecond <= 63:
		return byte(kilobitsPerSecond), nil
	case kilobitsPerSecond >= 64 && kilobitsPerSecond <= 568 && (kilobitsPerSecond-64)%8 == 0:
		return byte(64 + (kilobitsPerSecond-64)/8), nil
	case kilobitsPerSecond >= 576 && kilobitsPerSecond <= 8_640 && (kilobitsPerSecond-576)%64 == 0:
		return byte(128 + (kilobitsPerSecond-576)/64), nil
	default:
		return 0, fmt.Errorf("%w: bit rate %d bit/s is not representable", ErrInvalidValue, bitsPerSecond)
	}
}

func decodeQualityOfServiceBaseBitRate(encoded byte) (uint64, bool, error) {
	switch {
	case encoded == 0:
		return 0, true, nil
	case encoded <= 63:
		return uint64(encoded) * 1_000, false, nil
	case encoded <= 127:
		return (64 + uint64(encoded-64)*8) * 1_000, false, nil
	case encoded <= 254:
		return (576 + uint64(encoded-128)*64) * 1_000, false, nil
	default:
		return 0, false, nil
	}
}

func encodeQualityOfServiceExtendedBitRate(bitsPerSecond uint64) (byte, error) {
	switch {
	case bitsPerSecond >= 8_700_000 && bitsPerSecond <= 16_000_000 && (bitsPerSecond-8_600_000)%100_000 == 0:
		return byte((bitsPerSecond - 8_600_000) / 100_000), nil
	case bitsPerSecond >= 17_000_000 && bitsPerSecond <= 128_000_000 && (bitsPerSecond-16_000_000)%1_000_000 == 0:
		return byte(74 + (bitsPerSecond-16_000_000)/1_000_000), nil
	case bitsPerSecond >= 130_000_000 && bitsPerSecond <= 256_000_000 && (bitsPerSecond-128_000_000)%2_000_000 == 0:
		return byte(186 + (bitsPerSecond-128_000_000)/2_000_000), nil
	default:
		return 0, fmt.Errorf("%w: extended bit rate %d bit/s is not representable", ErrInvalidValue, bitsPerSecond)
	}
}

func decodeQualityOfServiceExtendedBitRate(encoded byte) (uint64, error) {
	switch {
	case encoded >= 1 && encoded <= 74:
		return (8_600_000 + uint64(encoded)*100_000), nil
	case encoded <= 186:
		return 16_000_000 + uint64(encoded-74)*1_000_000, nil
	case encoded <= 250:
		return 128_000_000 + uint64(encoded-186)*2_000_000, nil
	default:
		return 0, fmt.Errorf("%w: extended bit rate code %#02x is reserved", ErrInvalidValue, encoded)
	}
}

func encodeQualityOfServiceExtended2BitRate(bitsPerSecond uint64) (byte, error) {
	switch {
	case bitsPerSecond >= 260_000_000 && bitsPerSecond <= 500_000_000 && (bitsPerSecond-256_000_000)%4_000_000 == 0:
		return byte((bitsPerSecond - 256_000_000) / 4_000_000), nil
	case bitsPerSecond >= 510_000_000 && bitsPerSecond <= 1_500_000_000 && (bitsPerSecond-500_000_000)%10_000_000 == 0:
		return byte(61 + (bitsPerSecond-500_000_000)/10_000_000), nil
	case bitsPerSecond >= 1_600_000_000 && bitsPerSecond <= 10_000_000_000 && (bitsPerSecond-1_500_000_000)%100_000_000 == 0:
		return byte(161 + (bitsPerSecond-1_500_000_000)/100_000_000), nil
	default:
		return 0, fmt.Errorf("%w: extended-2 bit rate %d bit/s is not representable", ErrInvalidValue, bitsPerSecond)
	}
}

func decodeQualityOfServiceExtended2BitRate(encoded byte) (uint64, error) {
	switch {
	case encoded >= 1 && encoded <= 61:
		return 256_000_000 + uint64(encoded)*4_000_000, nil
	case encoded <= 161:
		return 500_000_000 + uint64(encoded-61)*10_000_000, nil
	case encoded <= 246:
		return 1_500_000_000 + uint64(encoded-161)*100_000_000, nil
	default:
		return 0, fmt.Errorf("%w: extended-2 bit rate code %#02x is reserved", ErrInvalidValue, encoded)
	}
}
