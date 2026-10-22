package ie

import (
	"fmt"
	"io"
)

const (
	apnAMBRBaseMaximum      uint64 = 8_640_000
	apnAMBRExtendedMaximum  uint64 = 256_000_000
	apnAMBRExtended2Maximum uint64 = 65_280_000_000
)

// APNAggregateMaximumBitRateFields is the structured value of the APN
// aggregate maximum bit rate IE (TS 24.301 clause 9.9.4.2). The bit rates are
// expressed in bit/s and must be exactly representable by the NAS coding.
// A zero bit rate encodes the explicit 0 kbps value.
type APNAggregateMaximumBitRateFields struct {
	DownlinkBitsPerSecond uint64
	UplinkBitsPerSecond   uint64
}

// NewAPNAggregateMaximumBitRateWithFields creates a tagged APN aggregate
// maximum bit rate IE from fields.
func NewAPNAggregateMaximumBitRateWithFields(fields *APNAggregateMaximumBitRateFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewAPNAggregateMaximumBitRate(payload), nil
}

// NewPositionalAPNAggregateMaximumBitRateWithFields creates a positional APN
// aggregate maximum bit rate IE from fields.
func NewPositionalAPNAggregateMaximumBitRateWithFields(fields *APNAggregateMaximumBitRateFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalAPNAggregateMaximumBitRate(payload), nil
}

// APNAggregateMaximumBitRate returns the IE as structured fields.
func (i *IE) APNAggregateMaximumBitRate() (*APNAggregateMaximumBitRateFields, error) {
	if err := i.requireDefinition(APNAggregateMaximumBitRateDefinition); err != nil {
		return nil, err
	}
	return ParseAPNAggregateMaximumBitRateFields(i.Payload)
}

// ParseAPNAggregateMaximumBitRateFields decodes an APN aggregate maximum bit
// rate value without its IEI and length octet.
func ParseAPNAggregateMaximumBitRateFields(payload []byte) (*APNAggregateMaximumBitRateFields, error) {
	fields := &APNAggregateMaximumBitRateFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that both bit rates are representable by the APN-AMBR
// coding.
func (f *APNAggregateMaximumBitRateFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: APN aggregate maximum bit rate fields are nil", ErrInvalidValue)
	}
	if _, err := encodeAPNAMBR(f.DownlinkBitsPerSecond); err != nil {
		return fmt.Errorf("downlink bit rate: %w", err)
	}
	if _, err := encodeAPNAMBR(f.UplinkBitsPerSecond); err != nil {
		return fmt.Errorf("uplink bit rate: %w", err)
	}
	return nil
}

// Marshal encodes fields into an APN aggregate maximum bit rate value without
// its IEI and length octet.
func (f *APNAggregateMaximumBitRateFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, fmt.Errorf("%w: APN aggregate maximum bit rate fields are nil", ErrInvalidValue)
	}
	downlink, err := encodeAPNAMBR(f.DownlinkBitsPerSecond)
	if err != nil {
		return nil, fmt.Errorf("downlink bit rate: %w", err)
	}
	uplink, err := encodeAPNAMBR(f.UplinkBitsPerSecond)
	if err != nil {
		return nil, fmt.Errorf("uplink bit rate: %w", err)
	}
	payload := make([]byte, apnAMBRPayloadLength(downlink, uplink))
	payload[0] = downlink.base
	payload[1] = uplink.base
	if len(payload) >= 4 {
		payload[2] = downlink.extended
		payload[3] = uplink.extended
	}
	if len(payload) == 6 {
		payload[4] = downlink.extended2
		payload[5] = uplink.extended2
	}
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *APNAggregateMaximumBitRateFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the APN aggregate maximum bit rate value
// payload produced by Marshal. It returns zero for a nil receiver.
func (f *APNAggregateMaximumBitRateFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	if f.DownlinkBitsPerSecond > apnAMBRExtendedMaximum || f.UplinkBitsPerSecond > apnAMBRExtendedMaximum {
		return 6
	}
	if f.DownlinkBitsPerSecond > apnAMBRBaseMaximum || f.UplinkBitsPerSecond > apnAMBRBaseMaximum {
		return 4
	}
	return 2
}

// UnmarshalBinary decodes an APN aggregate maximum bit rate value without its
// IEI and length octet. It leaves f unchanged when payload is invalid.
func (f *APNAggregateMaximumBitRateFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: APN aggregate maximum bit rate fields are nil", ErrInvalidValue)
	}
	if len(payload) != 2 && len(payload) != 4 && len(payload) != 6 {
		return fmt.Errorf("%w: APN aggregate maximum bit rate has %d value octets, want 2, 4, or 6", ErrInvalidLength, len(payload))
	}
	var downlinkExtended, uplinkExtended, downlinkExtended2, uplinkExtended2 byte
	if len(payload) >= 4 {
		downlinkExtended = payload[2]
		uplinkExtended = payload[3]
	}
	if len(payload) == 6 {
		downlinkExtended2 = payload[4]
		uplinkExtended2 = payload[5]
	}
	downlink, err := decodeAPNAMBR(payload[0], downlinkExtended, downlinkExtended2)
	if err != nil {
		return fmt.Errorf("downlink bit rate: %w", err)
	}
	uplink, err := decodeAPNAMBR(payload[1], uplinkExtended, uplinkExtended2)
	if err != nil {
		return fmt.Errorf("uplink bit rate: %w", err)
	}
	*f = APNAggregateMaximumBitRateFields{
		DownlinkBitsPerSecond: downlink,
		UplinkBitsPerSecond:   uplink,
	}
	return nil
}

type encodedAPNAMBR struct {
	base      byte
	extended  byte
	extended2 byte
}

func apnAMBRPayloadLength(downlink, uplink encodedAPNAMBR) int {
	if downlink.extended2 != 0 || uplink.extended2 != 0 {
		return 6
	}
	if downlink.extended != 0 || uplink.extended != 0 {
		return 4
	}
	return 2
}

func encodeAPNAMBR(bitsPerSecond uint64) (encodedAPNAMBR, error) {
	if bitsPerSecond <= apnAMBRBaseMaximum {
		base, err := encodeAPNAMBRBase(bitsPerSecond)
		return encodedAPNAMBR{base: base}, err
	}
	if bitsPerSecond <= apnAMBRExtendedMaximum {
		extended, err := encodeAPNAMBRExtended(bitsPerSecond)
		return encodedAPNAMBR{base: 0xfe, extended: extended}, err
	}
	if bitsPerSecond > apnAMBRExtended2Maximum {
		return encodedAPNAMBR{}, fmt.Errorf("%w: bit rate %d bit/s exceeds 65.28 Gbit/s", ErrInvalidValue, bitsPerSecond)
	}

	multiplier := (bitsPerSecond - 1) / apnAMBRExtendedMaximum
	remainder := bitsPerSecond - multiplier*apnAMBRExtendedMaximum
	base, extended, err := encodeAPNAMBRBaseAndExtended(remainder)
	if err != nil {
		return encodedAPNAMBR{}, err
	}
	return encodedAPNAMBR{base: base, extended: extended, extended2: byte(multiplier)}, nil
}

func encodeAPNAMBRBaseAndExtended(bitsPerSecond uint64) (byte, byte, error) {
	if bitsPerSecond <= apnAMBRBaseMaximum {
		base, err := encodeAPNAMBRBase(bitsPerSecond)
		return base, 0, err
	}
	extended, err := encodeAPNAMBRExtended(bitsPerSecond)
	return 0xfe, extended, err
}

func encodeAPNAMBRBase(bitsPerSecond uint64) (byte, error) {
	if bitsPerSecond == 0 {
		return 0xff, nil
	}
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

func encodeAPNAMBRExtended(bitsPerSecond uint64) (byte, error) {
	switch {
	case bitsPerSecond >= 8_700_000 && bitsPerSecond <= 16_000_000 && (bitsPerSecond-8_600_000)%100_000 == 0:
		return byte((bitsPerSecond - 8_600_000) / 100_000), nil
	case bitsPerSecond >= 17_000_000 && bitsPerSecond <= 128_000_000 && (bitsPerSecond-16_000_000)%1_000_000 == 0:
		return byte(74 + (bitsPerSecond-16_000_000)/1_000_000), nil
	case bitsPerSecond >= 130_000_000 && bitsPerSecond <= apnAMBRExtendedMaximum && (bitsPerSecond-128_000_000)%2_000_000 == 0:
		return byte(186 + (bitsPerSecond-128_000_000)/2_000_000), nil
	default:
		return 0, fmt.Errorf("%w: extended bit rate %d bit/s is not representable", ErrInvalidValue, bitsPerSecond)
	}
}

func decodeAPNAMBR(base, extended, extended2 byte) (uint64, error) {
	if extended2 == 0xff {
		extended2 = 0
	}
	if extended > 0xfa {
		extended = 0xfa
	}
	lower, err := decodeAPNAMBRBaseAndExtended(base, extended)
	if err != nil {
		return 0, err
	}
	if extended2 == 0 {
		return lower, nil
	}
	return uint64(extended2)*apnAMBRExtendedMaximum + lower, nil
}

func decodeAPNAMBRBaseAndExtended(base, extended byte) (uint64, error) {
	if extended == 0 {
		return decodeAPNAMBRBase(base)
	}
	switch {
	case extended <= 74:
		return 8_600_000 + uint64(extended)*100_000, nil
	case extended <= 186:
		return 16_000_000 + uint64(extended-74)*1_000_000, nil
	default:
		return 128_000_000 + uint64(extended-186)*2_000_000, nil
	}
}

func decodeAPNAMBRBase(encoded byte) (uint64, error) {
	switch {
	case encoded == 0:
		return 0, fmt.Errorf("%w: APN-AMBR base bit rate code 0 is reserved", ErrInvalidValue)
	case encoded <= 63:
		return uint64(encoded) * 1_000, nil
	case encoded <= 127:
		return (64 + uint64(encoded-64)*8) * 1_000, nil
	case encoded <= 254:
		return (576 + uint64(encoded-128)*64) * 1_000, nil
	default:
		return 0, nil
	}
}
