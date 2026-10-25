package ie

import (
	"fmt"
	"io"
)

const maximumMobileIdentityDigits = 15

// MobileIdentityType identifies the form carried by a mobile identity IE
// (TS 24.301 clause 9.9.2.3).
type MobileIdentityType uint8

// Mobile identity types (TS 24.008 clause 10.5.1.4).
const (
	MobileIdentityTypeIMSI MobileIdentityType = 0x01
	MobileIdentityTypeIMEI MobileIdentityType = 0x03
	MobileIdentityTypeTMSI MobileIdentityType = 0x05
	MobileIdentityTypeGUTI MobileIdentityType = 0x06
)

// MobileIdentityFields is the structured value of the tagged Mobile identity IE
// (TS 24.301 clause 9.9.2.3). Digits encodes the BCD-encoded identity (IMSI,
// IMEI, or TMSI digits). The max length of the value octet is 9 (32+ digits for
// IMEI/TMSI, 16+ digits for IMSI), while GUTI (11 bytes) exceeds this IE's limit
// and should use EPS mobile identity instead.
type MobileIdentityFields struct {
	Type   MobileIdentityType
	Digits string
}

// NewMobileIdentityWithFields creates a tagged Mobile identity IE from fields.
func NewMobileIdentityWithFields(fields *MobileIdentityFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewMobileIdentity(payload), nil
}

// NewPositionalMobileIdentityWithFields creates a positional Mobile identity IE
// from fields.
func NewPositionalMobileIdentityWithFields(fields *MobileIdentityFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalMobileIdentity(payload), nil
}

// MobileIdentity returns the IE as structured fields.
func (i *IE) MobileIdentity() (*MobileIdentityFields, error) {
	if err := i.requireDefinition(MobileIdentityDefinition); err != nil {
		return nil, err
	}
	return ParseMobileIdentityFields(i.Payload)
}

// ParseMobileIdentityFields decodes a mobile identity value without its IEI and
// length octet.
func ParseMobileIdentityFields(payload []byte) (*MobileIdentityFields, error) {
	fields := &MobileIdentityFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the identity alternatives are mutually exclusive and
// valid for the selected Type.
func (f *MobileIdentityFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: mobile identity fields are nil", ErrInvalidValue)
	}
	if f.Type != MobileIdentityTypeIMSI && f.Type != MobileIdentityTypeIMEI && f.Type != MobileIdentityTypeTMSI {
		return fmt.Errorf("%w: mobile identity type %#02x is reserved for this IE", ErrInvalidValue, f.Type)
	}
	if len(f.Digits) == 0 || len(f.Digits) > maximumMobileIdentityDigits {
		return fmt.Errorf("%w: mobile identity has %d digits, want 1..%d", ErrInvalidValue, len(f.Digits), maximumMobileIdentityDigits)
	}
	for _, digit := range f.Digits {
		if digit < '0' || digit > '9' {
			return fmt.Errorf("%w: mobile identity contains a non-decimal digit", ErrInvalidValue)
		}
	}
	return nil
}

// Marshal encodes fields into a mobile identity value without its IEI and length
// octet.
func (f *MobileIdentityFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return marshalMobileIdentityDigits(f.Type, f.Digits), nil
}

// MarshalTo encodes fields into b.
func (f *MobileIdentityFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the mobile identity value payload.
func (f *MobileIdentityFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	if len(f.Digits) == 0 {
		return 0
	}
	return 1 + len(f.Digits)/2
}

// UnmarshalBinary decodes a mobile identity value without its IEI and length octet.
// It leaves f unchanged when the payload is invalid.
func (f *MobileIdentityFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: mobile identity fields are nil", ErrInvalidValue)
	}
	if len(payload) == 0 || len(payload) > 9 {
		return fmt.Errorf("%w: mobile identity has %d value octets, want 1..9", ErrInvalidLength, len(payload))
	}
	typeValue := MobileIdentityType(payload[0] & 0x07)
	switch typeValue {
	case MobileIdentityTypeIMSI, MobileIdentityTypeIMEI, MobileIdentityTypeTMSI:
		digits, err := parseMobileIdentityDigits(payload)
		if err != nil {
			return err
		}
		var decoded MobileIdentityFields
		decoded.Type = typeValue
		decoded.Digits = digits
		if err := decoded.Validate(); err != nil {
			return err
		}
		*f = decoded
	default:
		return fmt.Errorf("%w: mobile identity type %#02x is reserved", ErrInvalidValue, typeValue)
	}
	return nil
}

func marshalMobileIdentityDigits(identityType MobileIdentityType, digits string) []byte {
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

func parseMobileIdentityDigits(payload []byte) (string, error) {
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
	if len(digits) > maximumMobileIdentityDigits {
		return "", fmt.Errorf("%w: mobile identity has %d digits, want at most %d", ErrInvalidLength, len(digits), maximumMobileIdentityDigits)
	}
	return string(digits), nil
}
