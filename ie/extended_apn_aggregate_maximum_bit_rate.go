package ie

import (
	"encoding/binary"
	"fmt"
	"io"
)

// ExtendedAPNAMBREntry holds an encoded bit-rate unit and value.
type ExtendedAPNAMBREntry struct {
	UnitCode  uint8
	ValueBits uint64
}

// ExtendedAPNAggregateMaximumBitRateUnit is the unit code for extended APN AMBR
// bit rates (TS 24.301 clause 9.9.4.29, table 9.9.4.29.1).
type ExtendedAPNAggregateMaximumBitRateUnit uint8

// Extended APN AMBR unit codes.
const (
	ExtendedAMBRUnitNotUsed ExtendedAPNAggregateMaximumBitRateUnit = 0x00
	ExtendedAMBRUnit1Mbps   ExtendedAPNAggregateMaximumBitRateUnit = 0x02
	ExtendedAMBRUnit4Mbps   ExtendedAPNAggregateMaximumBitRateUnit = 0x03
	ExtendedAMBRUnit16Mbps  ExtendedAPNAggregateMaximumBitRateUnit = 0x04
	ExtendedAMBRUnit64Mbps  ExtendedAPNAggregateMaximumBitRateUnit = 0x05
	ExtendedAMBRUnit256Mbps ExtendedAPNAggregateMaximumBitRateUnit = 0x06
	ExtendedAMBRUnit1Gbps   ExtendedAPNAggregateMaximumBitRateUnit = 0x07
	ExtendedAMBRUnit4Gbps   ExtendedAPNAggregateMaximumBitRateUnit = 0x08
	ExtendedAMBRUnit16Gbps  ExtendedAPNAggregateMaximumBitRateUnit = 0x09
	ExtendedAMBRUnit64Gbps  ExtendedAPNAggregateMaximumBitRateUnit = 0x0a
	ExtendedAMBRUnit256Gbps ExtendedAPNAggregateMaximumBitRateUnit = 0x0b
	ExtendedAMBRUnit1Tbps   ExtendedAPNAggregateMaximumBitRateUnit = 0x0c
	ExtendedAMBRUnit4Tbps   ExtendedAPNAggregateMaximumBitRateUnit = 0x0d
	ExtendedAMBRUnit16Tbps  ExtendedAPNAggregateMaximumBitRateUnit = 0x0e
	ExtendedAMBRUnit64Tbps  ExtendedAPNAggregateMaximumBitRateUnit = 0x0f
	ExtendedAMBRUnit256Tbps ExtendedAPNAggregateMaximumBitRateUnit = 0x10
	ExtendedAMBRUnit1Pbps   ExtendedAPNAggregateMaximumBitRateUnit = 0x11
	ExtendedAMBRUnit4Pbps   ExtendedAPNAggregateMaximumBitRateUnit = 0x12
	ExtendedAMBRUnit16Pbps  ExtendedAPNAggregateMaximumBitRateUnit = 0x13
	ExtendedAMBRUnit64Pbps  ExtendedAPNAggregateMaximumBitRateUnit = 0x14
	ExtendedAMBRUnit256Pbps ExtendedAPNAggregateMaximumBitRateUnit = 0x15
	_                                                              /* reserved in this version */
)

// extendedAPNAMBRBitsPerSecondToUnit converts a bit rate value to the best-fit
// unit code and returns the encoded 16-bit binary value.
func extendedAPNAMBRBitsPerSecondToUnit(bitsPerSecond uint64) (unit ExtendedAPNAggregateMaximumBitRateUnit, value uint16) {
	switch {
	case bitsPerSecond >= 1_000_000 && bitsPerSecond < 4_000_000:
		value = uint16((bitsPerSecond - 1_000_000) / 1_000_000)
		return ExtendedAMBRUnit1Mbps, value
	case bitsPerSecond >= 4_000_000 && bitsPerSecond < 16_000_000:
		value = uint16((bitsPerSecond - 4_000_000) / 4_000_000)
		return ExtendedAMBRUnit4Mbps, value
	case bitsPerSecond >= 16_000_000 && bitsPerSecond < 64_000_000:
		value = uint16((bitsPerSecond - 16_000_000) / 16_000_000)
		return ExtendedAMBRUnit16Mbps, value
	case bitsPerSecond >= 64_000_000 && bitsPerSecond < 256_000_000:
		value = uint16((bitsPerSecond - 64_000_000) / 64_000_000)
		return ExtendedAMBRUnit64Mbps, value
	case bitsPerSecond >= 256_000_000 && bitsPerSecond < 1_000_000_000:
		value = uint16((bitsPerSecond - 256_000_000) / 256_000_000)
		return ExtendedAMBRUnit256Mbps, value
	case bitsPerSecond >= 1_000_000_000 && bitsPerSecond < 4_000_000_000:
		value = uint16((bitsPerSecond - 1_000_000_000) / 1_000_000_000)
		return ExtendedAMBRUnit1Gbps, value
	case bitsPerSecond >= 4_000_000_000 && bitsPerSecond < 16_000_000_000:
		value = uint16((bitsPerSecond - 4_000_000_000) / 4_000_000_000)
		return ExtendedAMBRUnit4Gbps, value
	case bitsPerSecond >= 16_000_000_000 && bitsPerSecond < 64_000_000_000:
		value = uint16((bitsPerSecond - 16_000_000_000) / 16_000_000_000)
		return ExtendedAMBRUnit16Gbps, value
	case bitsPerSecond >= 64_000_000_000 && bitsPerSecond < 256_000_000_000:
		value = uint16((bitsPerSecond - 64_000_000_000) / 64_000_000_000)
		return ExtendedAMBRUnit64Gbps, value
	case bitsPerSecond >= 256_000_000_000 && bitsPerSecond < 1_000_000_000_000:
		value = uint16((bitsPerSecond - 256_000_000_000) / 256_000_000_000)
		return ExtendedAMBRUnit256Gbps, value
	case bitsPerSecond >= 1_000_000_000_000 && bitsPerSecond < 4_000_000_000_000:
		value = uint16((bitsPerSecond - 1_000_000_000_000) / 1_000_000_000_000)
		return ExtendedAMBRUnit1Tbps, value
	case bitsPerSecond >= 4_000_000_000_000 && bitsPerSecond < 16_000_000_000_000:
		value = uint16((bitsPerSecond - 4_000_000_000_000) / 4_000_000_000_000)
		return ExtendedAMBRUnit4Tbps, value
	case bitsPerSecond >= 16_000_000_000_000 && bitsPerSecond < 64_000_000_000_000:
		value = uint16((bitsPerSecond - 16_000_000_000_000) / 16_000_000_000_000)
		return ExtendedAMBRUnit16Tbps, value
	case bitsPerSecond >= 64_000_000_000_000 && bitsPerSecond < 256_000_000_000_000:
		value = uint16((bitsPerSecond - 64_000_000_000_000) / 64_000_000_000_000)
		return ExtendedAMBRUnit64Tbps, value
	case bitsPerSecond >= 256_000_000_000_000 && bitsPerSecond < 1_000_000_000_000_000:
		value = uint16((bitsPerSecond - 256_000_000_000_000) / 256_000_000_000_000)
		return ExtendedAMBRUnit256Tbps, value
	case bitsPerSecond >= 1_000_000_000_000_000 && bitsPerSecond < 4_000_000_000_000_000:
		value = uint16((bitsPerSecond - 1_000_000_000_000_000) / 1_000_000_000_000_000)
		return ExtendedAMBRUnit1Pbps, value
	case bitsPerSecond >= 4_000_000_000_000_000 && bitsPerSecond < 16_000_000_000_000_000:
		value = uint16((bitsPerSecond - 4_000_000_000_000_000) / 4_000_000_000_000_000)
		return ExtendedAMBRUnit4Pbps, value
	case bitsPerSecond >= 16_000_000_000_000_000 && bitsPerSecond < 64_000_000_000_000_000:
		value = uint16((bitsPerSecond - 16_000_000_000_000_000) / 16_000_000_000_000_000)
		return ExtendedAMBRUnit16Pbps, value
	case bitsPerSecond >= 64_000_000_000_000_000 && bitsPerSecond < 256_000_000_000_000_000:
		value = uint16((bitsPerSecond - 64_000_000_000_000_000) / 64_000_000_000_000_000)
		return ExtendedAMBRUnit64Pbps, value
	case bitsPerSecond >= 256_000_000_000_000_000 && bitsPerSecond < 1_000_000_000_000_000_000:
		value = uint16((bitsPerSecond - 256_000_000_000_000_000) / 256_000_000_000_000_000)
		return ExtendedAMBRUnit256Pbps, value
	default:
		if bitsPerSecond > 0 {
			value = uint16((bitsPerSecond - 256_000_000_000_000_000) / 256_000_000_000_000_000)
			return ExtendedAMBRUnit256Pbps, value
		}
		return ExtendedAMBRUnitNotUsed, 0
	}
}

// extendedAPNAMBRBitsPerSecondFromUnit converts a unit code and 16-bit value
// back to bits per second.
func extendedAPNAMBRBitsPerSecondFromUnit(unit ExtendedAPNAggregateMaximumBitRateUnit, value uint16) uint64 {
	switch unit {
	case ExtendedAMBRUnit1Mbps:
		return 1_000_000 + uint64(value)*1_000_000
	case ExtendedAMBRUnit4Mbps:
		return 4_000_000 + uint64(value)*4_000_000
	case ExtendedAMBRUnit16Mbps:
		return 16_000_000 + uint64(value)*16_000_000
	case ExtendedAMBRUnit64Mbps:
		return 64_000_000 + uint64(value)*64_000_000
	case ExtendedAMBRUnit256Mbps:
		return 256_000_000 + uint64(value)*256_000_000
	case ExtendedAMBRUnit1Gbps:
		return 1_000_000_000 + uint64(value)*1_000_000_000
	case ExtendedAMBRUnit4Gbps:
		return 4_000_000_000 + uint64(value)*4_000_000_000
	case ExtendedAMBRUnit16Gbps:
		return 16_000_000_000 + uint64(value)*16_000_000_000
	case ExtendedAMBRUnit64Gbps:
		return 64_000_000_000 + uint64(value)*64_000_000_000
	case ExtendedAMBRUnit256Gbps:
		return 256_000_000_000 + uint64(value)*256_000_000_000
	case ExtendedAMBRUnit1Tbps:
		return 1_000_000_000_000 + uint64(value)*1_000_000_000_000
	case ExtendedAMBRUnit4Tbps:
		return 4_000_000_000_000 + uint64(value)*4_000_000_000_000
	case ExtendedAMBRUnit16Tbps:
		return 16_000_000_000_000 + uint64(value)*16_000_000_000_000
	case ExtendedAMBRUnit64Tbps:
		return 64_000_000_000_000 + uint64(value)*64_000_000_000_000
	case ExtendedAMBRUnit256Tbps:
		return 256_000_000_000_000 + uint64(value)*256_000_000_000_000
	case ExtendedAMBRUnit1Pbps:
		return 1_000_000_000_000_000 + uint64(value)*1_000_000_000_000_000
	case ExtendedAMBRUnit4Pbps:
		return 4_000_000_000_000_000 + uint64(value)*4_000_000_000_000_000
	case ExtendedAMBRUnit16Pbps:
		return 16_000_000_000_000_000 + uint64(value)*16_000_000_000_000_000
	case ExtendedAMBRUnit64Pbps:
		return 64_000_000_000_000_000 + uint64(value)*64_000_000_000_000_000
	case ExtendedAMBRUnit256Pbps:
		return 256_000_000_000_000_000 + uint64(value)*256_000_000_000_000_000
	default:
		return 0
	}
}

// extendedAPNAMBRBitsPerSecondFromUnit32 converts a unit code and 32-bit value
// (used in the 8-octet extended format) back to bits per second. The uint32
// value is decoded as an offset added to the unit's base rate, then multiplied
// by the unit's multiplier.
func extendedAPNAMBRBitsPerSecondFromUnit32(unit ExtendedAPNAggregateMaximumBitRateUnit, value uint32) uint64 {
	switch unit {
	case ExtendedAMBRUnit64Gbps:
		return 64_000_000_000 + uint64(value)*64_000_000_000
	case ExtendedAMBRUnit1Tbps:
		return 1_000_000_000_000 + uint64(value)*1_000_000_000_000
	case ExtendedAMBRUnit4Tbps:
		return 4_000_000_000_000 + uint64(value)*4_000_000_000_000
	case ExtendedAMBRUnit16Tbps:
		return 16_000_000_000_000 + uint64(value)*16_000_000_000_000
	case ExtendedAMBRUnit64Tbps:
		return 64_000_000_000_000 + uint64(value)*64_000_000_000_000
	case ExtendedAMBRUnit256Tbps:
		return 256_000_000_000_000 + uint64(value)*256_000_000_000_000
	case ExtendedAMBRUnit1Pbps:
		return 1_000_000_000_000_000 + uint64(value)*1_000_000_000_000_000
	}
	return 0
}

// ExtendedAPNAggregateMaximumBitRateFields is the structured value of the
// Extended APN aggregate maximum bit rate IE (TS 24.301 clause 9.9.4.29).
//
// Wire layout (8 bytes):
//
//	Octet 3: Unit for extended APN-AMBR for downlink
//	Octets 4-5: Extended APN-AMBR for downlink (16-bit big-endian in units of octet 3)
//	Octet 6: Unit for extended APN-AMBR for uplink
//	Octets 7-8: Extended APN-AMBR for uplink (16-bit big-endian in units of octet 6)
//
// Each AMBR value is expressed in bits per second. The encoder selects the
// appropriate unit multiplier automatically when marshaling. Zero values are
// encoded as UnitNotUsed (0x00).
type ExtendedAPNAggregateMaximumBitRateFields struct {
	UplinkBitsPerSecond   uint64 // Downlink rate (unit code encodes the multiplier)
	DownlinkBitsPerSecond uint64 // Uplink rate (unit code encodes the multiplier)
}

// NewExtendedAPNAggregateMaximumBitRateWithFields creates a tagged Extended
// APN aggregate maximum bit rate IE from fields.
func NewExtendedAPNAggregateMaximumBitRateWithFields(fields *ExtendedAPNAggregateMaximumBitRateFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewExtendedAPNAggregateMaximumBitRate(payload), nil
}

// NewPositionalExtendedAPNAggregateMaximumBitRateWithFields creates a positional
// Extended APN aggregate maximum bit rate IE from fields.
func NewPositionalExtendedAPNAggregateMaximumBitRateWithFields(fields *ExtendedAPNAggregateMaximumBitRateFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalExtendedAPNAggregateMaximumBitRate(payload), nil
}

// ExtendedAPNAggregateMaximumBitRate returns the IE as structured fields. It
// accepts both tagged and positional instances created by this package.
func (i *IE) ExtendedAPNAggregateMaximumBitRate() (*ExtendedAPNAggregateMaximumBitRateFields, error) {
	if err := i.requireDefinition(ExtendedAPNAggregateMaximumBitRateDefinition); err != nil {
		return nil, err
	}
	return ParseExtendedAPNAggregateMaximumBitRateFields(i.Payload)
}

// ParseExtendedAPNAggregateMaximumBitRateFields decodes an Extended APN
// aggregate maximum bit rate value payload. The payload excludes the IEI and
// length octet.
func ParseExtendedAPNAggregateMaximumBitRateFields(payload []byte) (*ExtendedAPNAggregateMaximumBitRateFields, error) {
	fields := &ExtendedAPNAggregateMaximumBitRateFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that all non-zero bit rates are representable by the
// extended AMBR coding. A zero rate means "not used" and is valid.
func (f *ExtendedAPNAggregateMaximumBitRateFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: extended APN aggregate maximum bit rate fields are nil", ErrInvalidValue)
	}
	if f.DownlinkBitsPerSecond > 0 {
		_, _ = extendedAPNAMBRBitsPerSecondToUnit(f.DownlinkBitsPerSecond)
	}
	if f.UplinkBitsPerSecond > 0 {
		_, _ = extendedAPNAMBRBitsPerSecondToUnit(f.UplinkBitsPerSecond)
	}
	return nil
}

const ExtendedAPNAggregateMaximumBitRateLength = 8

// marshalExtendedAPNAMBR encodes f into b, using 6-octet format for rates <
// 64 Gbit/s and 8-octet format otherwise.
func marshalExtendedAPNAMBR(f *ExtendedAPNAggregateMaximumBitRateFields, b []byte) {
	if len(b) == 6 {
		// 6-octet format: unit (uint8) + value (uint16) per direction
		downUnit, downVal := extendedAPNAMBRBitsPerSecondToUnit(f.DownlinkBitsPerSecond)
		b[0] = byte(downUnit)
		binary.BigEndian.PutUint16(b[1:3], downVal)

		upUnit, upVal := extendedAPNAMBRBitsPerSecondToUnit(f.UplinkBitsPerSecond)
		b[3] = byte(upUnit)
		binary.BigEndian.PutUint16(b[4:6], upVal)
	} else {
		// 8-octet format: unit (uint8) + value (uint32) per direction
		downUnit, downVal := extendedAPNAMBRBitsPerSecondToUnit(f.DownlinkBitsPerSecond)
		b[0] = byte(downUnit)
		binary.BigEndian.PutUint32(b[1:5], uint32(downVal))

		upUnit, upVal := extendedAPNAMBRBitsPerSecondToUnit(f.UplinkBitsPerSecond)
		b[5] = byte(upUnit)
		binary.BigEndian.PutUint32(b[6:10], uint32(upVal))
	}
}

// Marshal encodes fields into an Extended APN aggregate maximum bit rate value
// without the IEI and length octet. Returns exactly 8 bytes.
func (f *ExtendedAPNAggregateMaximumBitRateFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	b := make([]byte, f.MarshalLen())
	marshalExtendedAPNAMBR(f, b)
	return b, nil
}

// MarshalTo encodes fields into b. It returns io.ErrShortBuffer without
// modifying b when b is too short.
func (f *ExtendedAPNAggregateMaximumBitRateFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < f.MarshalLen() {
		return io.ErrShortBuffer
	}
	marshalExtendedAPNAMBR(f, b)
	return nil
}

// MarshalLen returns the size of the Extended APN aggregate maximum bit rate
// value payload produced by Marshal. It returns zero for a nil receiver.
func (f *ExtendedAPNAggregateMaximumBitRateFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	if isSmallRate(f.UplinkBitsPerSecond) && isSmallRate(f.DownlinkBitsPerSecond) {
		return 6
	}
	return ExtendedAPNAggregateMaximumBitRateLength
}

func isSmallRate(rate uint64) bool {
	return rate < 64_000_000_000
}

// UnmarshalBinary decodes an Extended APN aggregate maximum bit rate value
// payload, excluding the IEI and length octet. It leaves f unchanged if
// decoding fails. Supports both 6-byte (unit+16bit value) and 8-byte
// (unit+32bit value) formats.
func (f *ExtendedAPNAggregateMaximumBitRateFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: extended APN aggregate maximum bit rate fields are nil", ErrInvalidValue)
	}
	if len(payload) != 6 && len(payload) != 8 {
		return fmt.Errorf("%w: extended APN aggregate maximum bit rate has %d value octets, want 6 or 8", ErrInvalidLength, len(payload))
	}

	var decoded ExtendedAPNAggregateMaximumBitRateFields

	if len(payload) == 6 {
		// 6-byte format: unit (uint8) + value (uint16) per direction
		downUnit := ExtendedAPNAggregateMaximumBitRateUnit(payload[0])
		if downUnit == ExtendedAMBRUnitNotUsed {
			decoded.DownlinkBitsPerSecond = 0
		} else {
			downVal := binary.BigEndian.Uint16(payload[1:3])
			decoded.DownlinkBitsPerSecond = extendedAPNAMBRBitsPerSecondFromUnit(downUnit, downVal)
		}
		upUnit := ExtendedAPNAggregateMaximumBitRateUnit(payload[3])
		if upUnit == ExtendedAMBRUnitNotUsed {
			decoded.UplinkBitsPerSecond = 0
		} else {
			upVal := binary.BigEndian.Uint16(payload[4:6])
			decoded.UplinkBitsPerSecond = extendedAPNAMBRBitsPerSecondFromUnit(upUnit, upVal)
		}
	} else {
		// 8-byte format: unit (uint8) + value (uint32) per direction
		downUnit := ExtendedAPNAggregateMaximumBitRateUnit(payload[0])
		if downUnit == ExtendedAMBRUnitNotUsed {
			decoded.DownlinkBitsPerSecond = 0
		} else {
			downVal := binary.BigEndian.Uint32(payload[1:5])
			decoded.DownlinkBitsPerSecond = extendedAPNAMBRBitsPerSecondFromUnit32(downUnit, downVal)
		}
		upUnit := ExtendedAPNAggregateMaximumBitRateUnit(payload[5])
		if upUnit == ExtendedAMBRUnitNotUsed {
			decoded.UplinkBitsPerSecond = 0
		} else {
			upVal := binary.BigEndian.Uint32(payload[6:10])
			decoded.UplinkBitsPerSecond = extendedAPNAMBRBitsPerSecondFromUnit32(upUnit, upVal)
		}
	}

	*f = decoded
	return nil
}
