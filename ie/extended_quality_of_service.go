package ie

import (
	"encoding/binary"
	"fmt"
	"io"
)

// ExtendedQualityOfServiceFields is the structured value of the Extended
// quality of service IE (TS 24.301 clause 9.9.4.30).
//
// Wire layout (10 bytes):
//
//	Octet 1: QCI (uint8, the IE field maps the full octet)
//	Octets 2-3: TrafficClass (uint16 big-endian)
//	Octet 4: PriorityLevel (uint8, bits 1-7 per spec)
//	Octet 5: ARPPreEmption (uint8, bits 7-8 per spec)
//	Octets 6-8: GBR (24-bit big-endian value in bits/sec)
//	Octets 9-10: MBR (24-bit big-endian value in bits/sec, upper 16 bits only
//		as only 2 bytes remain; lower byte is zeroed)
//	TOSField has no dedicated wire bytes in this version.
type ExtendedQualityOfServiceFields struct {
	QCI           uint8  // Octet 1 — QoS Class Identifier
	TrafficClass  uint16 // Octets 2-3 — Traffic Class Identifier (TCI)
	PriorityLevel uint8  // Octet 4 — bits 1-7 (Priority Level)
	ARPPreEmption uint8  // Octet 5 — bits 7-8 (Pre-emption Capability/Vulnerability)
	GBR           uint32 // Octets 6-8 — 24-bit big-endian (bits/sec)
	MBR           uint32 // Octets 9-10 — 24-bit big-endian (bits/sec, upper 16 bits only)
	TOSField      uint16 // Reserved — unused in this version
}

// NewExtendedQualityOfServiceWithFields creates a tagged Extended quality of
// service IE from structured fields. NewExtendedQualityOfService remains
// available for callers that need to construct the IE from a raw payload.
func NewExtendedQualityOfServiceWithFields(fields *ExtendedQualityOfServiceFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewExtendedQualityOfService(payload), nil
}

// NewPositionalExtendedQualityOfServiceWithFields creates a positional Extended
// quality of service IE from structured fields.
func NewPositionalExtendedQualityOfServiceWithFields(fields *ExtendedQualityOfServiceFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalExtendedQualityOfService(payload), nil
}

// ExtendedQualityOfService returns the IE as structured fields. It accepts
// both tagged and positional instances created by this package.
func (i *IE) ExtendedQualityOfService() (*ExtendedQualityOfServiceFields, error) {
	if err := i.requireDefinition(ExtendedQualityOfServiceDefinition); err != nil {
		return nil, err
	}
	return ParseExtendedQualityOfServiceFields(i.Payload)
}

// ParseExtendedQualityOfServiceFields decodes an Extended quality of service
// value payload. The payload excludes the IEI and length octet.
func ParseExtendedQualityOfServiceFields(payload []byte) (*ExtendedQualityOfServiceFields, error) {
	fields := &ExtendedQualityOfServiceFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil. It does not check field
// widths because this IE's wire format stores the raw values directly.
func (f *ExtendedQualityOfServiceFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: extended quality of service fields are nil", ErrInvalidValue)
	}
	return nil
}

const ExtendedQualityOfServiceLength = 10

// marshalExtendedQoS encodes f into the 10-octet wire format.
func marshalExtendedQoS(f *ExtendedQualityOfServiceFields, b []byte) {
	b[0] = f.QCI
	binary.BigEndian.PutUint16(b[1:3], f.TrafficClass)
	b[3] = f.PriorityLevel
	b[4] = f.ARPPreEmption
	// GBR: 24-bit big-endian — mask to 24 bits (upper 3 bytes of uint32)
	gbr := f.GBR & 0xffffff
	b[5] = byte(gbr >> 16)
	b[6] = byte(gbr >> 8)
	b[7] = byte(gbr)
	// MBR: 24-bit big-endian — mask to 24 bits
	mbr := f.MBR & 0xffffff
	b[8] = byte(mbr >> 16)
	b[9] = byte(mbr >> 8)
}

// Marshal encodes fields into an Extended quality of service value payload
// without the IEI and length octet. Returns exactly 10 bytes.
func (f *ExtendedQualityOfServiceFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	b := make([]byte, ExtendedQualityOfServiceLength)
	marshalExtendedQoS(f, b)
	return b, nil
}

// MarshalTo encodes fields into b. It returns io.ErrShortBuffer without
// modifying b when b is too short.
func (f *ExtendedQualityOfServiceFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < ExtendedQualityOfServiceLength {
		return io.ErrShortBuffer
	}
	marshalExtendedQoS(f, b)
	return nil
}

// MarshalLen returns the size of the Extended quality of service value payload
// produced by Marshal. It returns zero for a nil receiver.
func (f *ExtendedQualityOfServiceFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return ExtendedQualityOfServiceLength
}

// UnmarshalBinary decodes an Extended quality of service value payload,
// excluding the IEI and length octet. It leaves f unchanged if decoding fails.
func (f *ExtendedQualityOfServiceFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: extended quality of service fields are nil", ErrInvalidValue)
	}
	if len(payload) != ExtendedQualityOfServiceLength {
		return fmt.Errorf("%w: extended quality of service has %d value octets, want 10", ErrInvalidLength, len(payload))
	}

	var decoded ExtendedQualityOfServiceFields
	decoded.QCI = payload[0]
	decoded.TrafficClass = binary.BigEndian.Uint16(payload[1:3])
	decoded.PriorityLevel = payload[3]
	decoded.ARPPreEmption = payload[4]
	// GBR: 24-bit big-endian from octets 5-7 (indices)
	decoded.GBR = uint32(payload[5])<<16 | uint32(payload[6])<<8 | uint32(payload[7])
	// MBR: 24-bit big-endian from octets 8-9 — only 16 bits available in
	// a 10-octet payload (upper 2 bytes). The lower byte is unused/reserved.
	decoded.MBR = uint32(payload[8])<<16 | uint32(payload[9])<<8

	*f = decoded
	return nil
}
