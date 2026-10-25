package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// Security header types from TS 24.301 Table 9.3.1.
const (
	SecurityHeaderTypePlain                                          uint8 = 0x00
	SecurityHeaderTypeIntegrityProtected                             uint8 = 0x01
	SecurityHeaderTypeIntegrityProtectedAndCiphered                  uint8 = 0x02
	SecurityHeaderTypeIntegrityProtectedWithNewEPSContext            uint8 = 0x03
	SecurityHeaderTypeIntegrityProtectedAndCipheredWithNewEPSContext uint8 = 0x04
	SecurityHeaderTypeIntegrityProtectedAndPartiallyCiphered         uint8 = 0x05
	SecurityHeaderTypeServiceRequest                                 uint8 = 0xC0 // legacy non-standard L3
)

// Protocol discriminators from TS 24.301 clause 9.1 (stored in the low nibble of octet 1).
const (
	ProtocolDiscriminatorEPSMobilityManagement uint8 = 0x7E
	ProtocolDiscriminatorEPSSessionManagement  uint8 = 0x2E

	// PD half-octet wire values (low nibble).
	PdEmm uint8 = 0x7
	PdEsm uint8 = 0x2
)

// EmmHeader is the EMM NAS message header defined by TS 24.301 clause 9.3.1.
//
// The first octet packs SecurityHeaderType in the high nibble and ProtocolDiscriminator
// in the low nibble. Plain messages (SHT=0x00) have a 2-octet layout: the second octet
// carries the message type. Integrity-protected messages (SHT 0x01..0x05 or 0xC0) have a
// 7-octet layout including MAC (octets 1-4), sequence number (octet 5), and message type
// (octet 6). See TS 24.301 clauses 9.5 and 9.6 for integrity protection details.
type EmmHeader struct {
	SecurityHeaderType    uint8
	ProtocolDiscriminator uint8
	Mac                   [4]byte
	SequenceNumber        uint8
	MsgType               uint8
}

// MarshalLen returns the size of EmmHeader on the wire: 2 octets for plain,
// 7 for integrity-protected. The nil receiver returns 0.
func (h *EmmHeader) MarshalLen() int {
	if h == nil {
		return 0
	}
	if h.SecurityHeaderType == SecurityHeaderTypePlain || h.SecurityHeaderType == 0xC0 {
		return 2
	}
	sht := h.SecurityHeaderType
	if sht >= SecurityHeaderTypeIntegrityProtected &&
		sht <= SecurityHeaderTypeIntegrityProtectedAndPartiallyCiphered {
		return 7
	}
	return 0
}

// MarshalTo encodes the header into b. Validates that ProtocolDiscriminator == 0x7E
// and SecurityHeaderType is in {0..5, 0xC0}; returns a wrapped ie.ErrInvalidValue
// on bad values. Returns io.ErrShortBuffer without writing if b is too short.
func (h *EmmHeader) MarshalTo(b []byte) error {
	if h == nil {
		return fmt.Errorf("%w: EmmHeader is nil", ie.ErrNilIE)
	}

	expectedPD := PdEmm
	if h.ProtocolDiscriminator != expectedPD {
		return fmt.Errorf("%w: EMM protocol discriminator %#02x, want %#02x",
			ie.ErrInvalidValue, h.ProtocolDiscriminator, expectedPD)
	}

	sht := h.SecurityHeaderType
	if sht != SecurityHeaderTypePlain && sht != 0xC0 &&
		(sht < SecurityHeaderTypeIntegrityProtected || sht > SecurityHeaderTypeIntegrityProtectedAndPartiallyCiphered) {
		return fmt.Errorf("%w: security header type %#02x", ie.ErrInvalidValue, sht)
	}

	required := h.MarshalLen()
	if required == 0 {
		return fmt.Errorf("%w: security header type %#02x yields zero length", ie.ErrInvalidValue, sht)
	}

	if len(b) < required {
		return io.ErrShortBuffer
	}

	if sht == SecurityHeaderTypePlain || sht == 0xC0 {
		b[0] = sht<<4 | h.ProtocolDiscriminator
		b[1] = h.MsgType
		return nil
	}

	b[0] = sht<<4 | h.ProtocolDiscriminator
	copy(b[1:5], h.Mac[:])
	b[5] = h.SequenceNumber
	b[6] = h.MsgType
	return nil
}

// UnmarshalBinary decodes the EMM header from b, returning the remaining bytes
// (payload after the header) and any error. The input buffer b is unchanged on
// error. Truncated input returns io.ErrUnexpectedEOF.
func (h *EmmHeader) UnmarshalBinary(b []byte) ([]byte, error) {
	if len(b) < 2 {
		return nil, io.ErrUnexpectedEOF
	}

	sht := b[0] >> 4
	pd := b[0] & 0x0f

	if pd != PdEmm {
		return nil, fmt.Errorf("%w: EMM protocol discriminator %#02x, want %#02x",
			ie.ErrInvalidValue, pd, PdEmm)
	}

	if sht == SecurityHeaderTypePlain || sht == 0xC0 {
		h.SecurityHeaderType = sht
		h.ProtocolDiscriminator = pd
		h.MsgType = b[1]
		return b[2:], nil
	}

	if sht >= SecurityHeaderTypeIntegrityProtected && sht <= SecurityHeaderTypeIntegrityProtectedAndPartiallyCiphered {
		if len(b) < 7 {
			return nil, io.ErrUnexpectedEOF
		}
		h.SecurityHeaderType = sht
		h.ProtocolDiscriminator = pd
		copy(h.Mac[:], b[1:5])
		h.SequenceNumber = b[5]
		h.MsgType = b[6]
		return b[7:], nil
	}

	return nil, fmt.Errorf("%w: security header type %#02x", ie.ErrInvalidValue, sht)
}

// EsmHeader is the ESM NAS message header defined by TS 24.301 clauses 9.3.2 and 9.4.
//
// The header is always exactly 3 octets: octet 1 packs BearerIdentity (high nibble)
// and ProtocolDiscriminator (low nibble), octet 2 is the Procedure Transaction
// Identity (PTI), and octet 3 is the message type.
type EsmHeader struct {
	BearerIdentity               uint8
	ProtocolDiscriminator        uint8
	ProcedureTransactionIdentity uint8
	MsgType                      uint8
}

// EPSBearerIdentity returns the bearer identity extracted from BearerIdentity's
// low nibble. The high nibble bits are reserved and ignored.
func (h EsmHeader) EPSBearerIdentity() uint8 {
	return h.BearerIdentity & 0x0f
}

// PTI returns the Procedure Transaction Identity.
func (h EsmHeader) PTI() uint8 {
	return h.ProcedureTransactionIdentity
}

// MarshalLen returns 3 for the fixed ESM header length. The nil receiver returns 0.
func (h *EsmHeader) MarshalLen() int {
	if h == nil {
		return 0
	}
	return 3
}

// MarshalTo encodes the ESM header into b (always 3 octets). Validates that
// ProtocolDiscriminator == 0x2E and BearerIdentity <= 0x0F; returns a wrapped
// ie.ErrInvalidValue on bad values. Returns io.ErrShortBuffer if b is too short.
func (h *EsmHeader) MarshalTo(b []byte) error {
	if h == nil {
		return fmt.Errorf("%w: EsmHeader is nil", ie.ErrNilIE)
	}

	expectedPD := PdEsm
	if h.ProtocolDiscriminator != expectedPD {
		return fmt.Errorf("%w: ESM protocol discriminator %#02x, want %#02x",
			ie.ErrInvalidValue, h.ProtocolDiscriminator, expectedPD)
	}

	if h.BearerIdentity > 0x0f {
		return fmt.Errorf("%w: bearer identity %#02x exceeds 4-bit mask", ie.ErrInvalidValue, h.BearerIdentity)
	}

	if len(b) < 3 {
		return io.ErrShortBuffer
	}

	b[0] = h.BearerIdentity<<4 | h.ProtocolDiscriminator
	b[1] = h.ProcedureTransactionIdentity
	b[2] = h.MsgType
	return nil
}

// UnmarshalBinary decodes the ESM header from b (always 3 octets), returning the
// remaining bytes (payload after the header) and any error. The input buffer b is
// unchanged on error. Truncated input returns io.ErrUnexpectedEOF.
func (h *EsmHeader) UnmarshalBinary(b []byte) ([]byte, error) {
	if len(b) < 3 {
		return nil, io.ErrUnexpectedEOF
	}

	pd := b[0] & 0x0f
	if pd != PdEsm {
		return nil, fmt.Errorf("%w: ESM protocol discriminator %#02x, want %#02x",
			ie.ErrInvalidValue, pd, PdEsm)
	}

	bi := b[0] >> 4
	if bi > 0x0f {
		return nil, fmt.Errorf("%w: bearer identity %#02x exceeds 4-bit mask", ie.ErrInvalidValue, bi)
	}

	h.BearerIdentity = b[0] >> 4
	h.ProtocolDiscriminator = pd
	h.ProcedureTransactionIdentity = b[1]
	h.MsgType = b[2]
	return b[3:], nil
}
