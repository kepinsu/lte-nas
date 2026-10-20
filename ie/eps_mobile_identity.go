package ie

import (
	"encoding/binary"
	"fmt"
	"io"
)

const maximumEPSMobileIdentityDigits = 15

// EPSMobileIdentityType identifies the form carried by an EPS mobile identity
// IE (TS 24.301 clause 9.9.3.12).
type EPSMobileIdentityType uint8

// EPS mobile identity types.
const (
	EPSMobileIdentityTypeIMSI EPSMobileIdentityType = 0x01
	EPSMobileIdentityTypeIMEI EPSMobileIdentityType = 0x03
	EPSMobileIdentityTypeGUTI EPSMobileIdentityType = 0x06
)

// GUTI is a globally unique temporary UE identity.
type GUTI struct {
	PLMN       PLMN
	MMEGroupID uint16
	MMECode    uint8
	MTMSI      uint32
}

// EPSMobileIdentityFields is the structured value of the EPS mobile identity
// IE. Set Digits for an IMSI or IMEI identity; set GUTI for a GUTI identity.
// The unused alternative must be empty or nil.
type EPSMobileIdentityFields struct {
	Type   EPSMobileIdentityType
	Digits string
	GUTI   *GUTI
}

// NewEPSMobileIdentityWithFields creates a tagged EPS mobile identity IE from
// fields.
func NewEPSMobileIdentityWithFields(fields *EPSMobileIdentityFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewEPSMobileIdentity(payload), nil
}

// NewPositionalEPSMobileIdentityWithFields creates a positional EPS mobile
// identity IE from fields.
func NewPositionalEPSMobileIdentityWithFields(fields *EPSMobileIdentityFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalEPSMobileIdentity(payload), nil
}

// EPSMobileIdentity returns the IE as structured fields.
func (i *IE) EPSMobileIdentity() (*EPSMobileIdentityFields, error) {
	if err := i.requireDefinition(EPSMobileIdentityDefinition); err != nil {
		return nil, err
	}
	return ParseEPSMobileIdentityFields(i.Payload)
}

// ParseEPSMobileIdentityFields decodes an EPS mobile identity value without
// its IEI and length octet.
func ParseEPSMobileIdentityFields(payload []byte) (*EPSMobileIdentityFields, error) {
	fields := &EPSMobileIdentityFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the identity alternatives are mutually exclusive and
// valid for the selected Type.
func (f *EPSMobileIdentityFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: EPS mobile identity fields are nil", ErrInvalidValue)
	}
	switch f.Type {
	case EPSMobileIdentityTypeIMSI, EPSMobileIdentityTypeIMEI:
		if f.GUTI != nil {
			return fmt.Errorf("%w: IMSI or IMEI identity cannot include a GUTI", ErrInvalidValue)
		}
		if len(f.Digits) == 0 || len(f.Digits) > maximumEPSMobileIdentityDigits {
			return fmt.Errorf("%w: mobile identity has %d digits, want 1..%d", ErrInvalidValue, len(f.Digits), maximumEPSMobileIdentityDigits)
		}
		for _, digit := range f.Digits {
			if digit < '0' || digit > '9' {
				return fmt.Errorf("%w: mobile identity contains a non-decimal digit", ErrInvalidValue)
			}
		}
		return nil
	case EPSMobileIdentityTypeGUTI:
		if f.Digits != "" || f.GUTI == nil {
			return fmt.Errorf("%w: GUTI identity requires GUTI and no digits", ErrInvalidValue)
		}
		return f.GUTI.PLMN.validate()
	default:
		return fmt.Errorf("%w: EPS mobile identity type %#02x is reserved", ErrInvalidValue, f.Type)
	}
}

// Marshal encodes fields into an EPS mobile identity value without its IEI and
// length octet.
func (f *EPSMobileIdentityFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	if f.Type == EPSMobileIdentityTypeGUTI {
		return marshalGUTI(f.GUTI)
	}
	return marshalEPSMobileIdentityDigits(f.Type, f.Digits), nil
}

// MarshalTo encodes fields into b.
func (f *EPSMobileIdentityFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the EPS mobile identity value payload.
func (f *EPSMobileIdentityFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	if f.Type == EPSMobileIdentityTypeGUTI {
		return 11
	}
	if len(f.Digits) == 0 {
		return 0
	}
	return 1 + len(f.Digits)/2
}

// UnmarshalBinary decodes an EPS mobile identity value without its IEI and
// length octet. It leaves f unchanged when the payload is invalid.
func (f *EPSMobileIdentityFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: EPS mobile identity fields are nil", ErrInvalidValue)
	}
	if len(payload) == 0 || len(payload) > 11 {
		return fmt.Errorf("%w: EPS mobile identity has %d value octets, want 1..11", ErrInvalidLength, len(payload))
	}
	typeValue := EPSMobileIdentityType(payload[0] & 0x07)
	var decoded EPSMobileIdentityFields
	switch typeValue {
	case EPSMobileIdentityTypeIMSI, EPSMobileIdentityTypeIMEI:
		digits, err := parseEPSMobileIdentityDigits(payload)
		if err != nil {
			return err
		}
		decoded = EPSMobileIdentityFields{Type: typeValue, Digits: digits}
	case EPSMobileIdentityTypeGUTI:
		guti, err := parseGUTI(payload)
		if err != nil {
			return err
		}
		decoded = EPSMobileIdentityFields{Type: typeValue, GUTI: &guti}
	default:
		return fmt.Errorf("%w: EPS mobile identity type %#02x is reserved", ErrInvalidValue, typeValue)
	}
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}

func marshalEPSMobileIdentityDigits(identityType EPSMobileIdentityType, digits string) []byte {
	payload := make([]byte, 1+len(digits)/2)
	payload[0] = (digits[0]-'0')<<4 | byte(identityType)
	if len(digits)%2 != 0 {
		payload[0] |= 0x08
	}
	for index := 1; index < len(digits); index += 2 {
		octet := (digits[index] - '0') & 0x0f
		if index+1 < len(digits) {
			octet |= (digits[index+1] - '0') << 4
		} else {
			octet |= 0xf0
		}
		payload[1+(index-1)/2] = octet
	}
	return payload
}

func parseEPSMobileIdentityDigits(payload []byte) (string, error) {
	firstDigit := payload[0] >> 4
	if firstDigit > 9 {
		return "", fmt.Errorf("%w: mobile identity contains an invalid BCD digit", ErrInvalidValue)
	}
	oddDigits := payload[0]&0x08 != 0
	digits := make([]byte, 0, len(payload)*2-1)
	digits = append(digits, '0'+firstDigit)
	for index, octet := range payload[1:] {
		low := octet & 0x0f
		high := octet >> 4
		if low > 9 {
			return "", fmt.Errorf("%w: mobile identity contains an invalid BCD digit", ErrInvalidValue)
		}
		digits = append(digits, '0'+low)
		last := index == len(payload)-2
		if last && !oddDigits {
			if high != 0x0f {
				return "", fmt.Errorf("%w: even mobile identity lacks its end mark", ErrInvalidValue)
			}
			continue
		}
		if high > 9 {
			return "", fmt.Errorf("%w: mobile identity contains an invalid BCD digit", ErrInvalidValue)
		}
		digits = append(digits, '0'+high)
	}
	if oddDigits != (len(digits)%2 != 0) {
		return "", fmt.Errorf("%w: mobile identity odd/even indication does not match its digits", ErrInvalidValue)
	}
	if len(digits) > maximumEPSMobileIdentityDigits {
		return "", fmt.Errorf("%w: mobile identity has %d digits, want at most %d", ErrInvalidLength, len(digits), maximumEPSMobileIdentityDigits)
	}
	return string(digits), nil
}

func marshalGUTI(guti *GUTI) ([]byte, error) {
	plmn, err := guti.PLMN.marshal()
	if err != nil {
		return nil, err
	}
	payload := make([]byte, 11)
	payload[0] = 0xf0 | byte(EPSMobileIdentityTypeGUTI)
	copy(payload[1:4], plmn[:])
	binary.BigEndian.PutUint16(payload[4:6], guti.MMEGroupID)
	payload[6] = guti.MMECode
	binary.BigEndian.PutUint32(payload[7:11], guti.MTMSI)
	return payload, nil
}

func parseGUTI(payload []byte) (GUTI, error) {
	if len(payload) != 11 {
		return GUTI{}, fmt.Errorf("%w: GUTI has %d value octets, want 11", ErrInvalidLength, len(payload))
	}
	if payload[0]&0xf8 != 0xf0 {
		return GUTI{}, fmt.Errorf("%w: GUTI spare and odd/even bits are invalid", ErrInvalidValue)
	}
	plmn, err := parsePLMN(payload[1:4])
	if err != nil {
		return GUTI{}, err
	}
	return GUTI{
		PLMN:       plmn,
		MMEGroupID: binary.BigEndian.Uint16(payload[4:6]),
		MMECode:    payload[6],
		MTMSI:      binary.BigEndian.Uint32(payload[7:11]),
	}, nil
}
