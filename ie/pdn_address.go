package ie

import (
	"bytes"
	"fmt"
	"io"
	"net/netip"
)

// PDNAddressFields is the structured value of the PDN address IE (TS 24.301
// clause 9.9.4.9). IPv6InterfaceIdentifier is the 64-bit interface identifier,
// not a full IPv6 address. IPv4 is required for IPv4 and IPv4v6 PDN types;
// netip.AddrFrom4([4]byte{}) explicitly represents DHCPv4 allocation.
type PDNAddressFields struct {
	Type                    PDNType
	IPv4                    netip.Addr
	IPv6InterfaceIdentifier [8]byte
}

// NewPositionalPDNAddressWithFields creates a positional PDN address IE from
// fields.
func NewPositionalPDNAddressWithFields(fields *PDNAddressFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalPDNAddress(payload), nil
}

// PDNAddress returns the IE as structured fields.
func (i *IE) PDNAddress() (*PDNAddressFields, error) {
	if err := i.requireDefinition(PDNAddressDefinition); err != nil {
		return nil, err
	}
	return ParsePDNAddressFields(i.Payload)
}

// ParsePDNAddressFields decodes a PDN address value without its IEI and
// length octet.
func ParsePDNAddressFields(payload []byte) (*PDNAddressFields, error) {
	fields := &PDNAddressFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the address alternatives match Type.
func (f *PDNAddressFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: PDN address fields are nil", ErrInvalidValue)
	}
	if !f.Type.valid() {
		return fmt.Errorf("%w: PDN type %#02x is reserved", ErrInvalidValue, f.Type)
	}
	interfaceIDIsZero := bytes.Equal(f.IPv6InterfaceIdentifier[:], make([]byte, len(f.IPv6InterfaceIdentifier)))
	switch f.Type {
	case PDNTypeIPv4:
		if !f.IPv4.IsValid() || !f.IPv4.Is4() || !interfaceIDIsZero {
			return fmt.Errorf("%w: IPv4 PDN address requires IPv4 and no IPv6 interface identifier", ErrInvalidValue)
		}
	case PDNTypeIPv6:
		if f.IPv4.IsValid() {
			return fmt.Errorf("%w: IPv6 PDN address must not include IPv4", ErrInvalidValue)
		}
	case PDNTypeIPv4v6:
		if !f.IPv4.IsValid() || !f.IPv4.Is4() {
			return fmt.Errorf("%w: IPv4v6 PDN address requires IPv4", ErrInvalidValue)
		}
	case PDNTypeNonIP, PDNTypeEthernet:
		if f.IPv4.IsValid() || !interfaceIDIsZero {
			return fmt.Errorf("%w: non-IP or Ethernet PDN address must contain spare zero octets", ErrInvalidValue)
		}
	}
	return nil
}

// Marshal encodes fields into a PDN address value without its IEI and length
// octet.
func (f *PDNAddressFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, f.MarshalLen())
	payload[0] = byte(f.Type)
	switch f.Type {
	case PDNTypeIPv4:
		ipv4 := f.IPv4.As4()
		copy(payload[1:], ipv4[:])
	case PDNTypeIPv6:
		copy(payload[1:], f.IPv6InterfaceIdentifier[:])
	case PDNTypeIPv4v6:
		copy(payload[1:9], f.IPv6InterfaceIdentifier[:])
		ipv4 := f.IPv4.As4()
		copy(payload[9:], ipv4[:])
	}
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *PDNAddressFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the PDN address value payload.
func (f *PDNAddressFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	switch f.Type {
	case PDNTypeIPv4, PDNTypeNonIP, PDNTypeEthernet:
		return 5
	case PDNTypeIPv6:
		return 9
	case PDNTypeIPv4v6:
		return 13
	default:
		return 0
	}
}

// UnmarshalBinary decodes a PDN address value without its IEI and length
// octet. It leaves f unchanged when the payload is invalid.
func (f *PDNAddressFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: PDN address fields are nil", ErrInvalidValue)
	}
	if len(payload) < 1 || payload[0]&0xf8 != 0 {
		return fmt.Errorf("%w: invalid PDN address type octet", ErrInvalidValue)
	}
	decoded := PDNAddressFields{Type: PDNType(payload[0])}
	switch decoded.Type {
	case PDNTypeIPv4:
		if len(payload) != 5 {
			return invalidPDNAddressLength(len(payload), 5)
		}
		decoded.IPv4 = netip.AddrFrom4([4]byte(payload[1:5]))
	case PDNTypeIPv6:
		if len(payload) != 9 {
			return invalidPDNAddressLength(len(payload), 9)
		}
		copy(decoded.IPv6InterfaceIdentifier[:], payload[1:9])
	case PDNTypeIPv4v6:
		if len(payload) != 13 {
			return invalidPDNAddressLength(len(payload), 13)
		}
		copy(decoded.IPv6InterfaceIdentifier[:], payload[1:9])
		decoded.IPv4 = netip.AddrFrom4([4]byte(payload[9:13]))
	case PDNTypeNonIP, PDNTypeEthernet:
		if len(payload) != 5 {
			return invalidPDNAddressLength(len(payload), 5)
		}
		if !bytes.Equal(payload[1:], []byte{0, 0, 0, 0}) {
			return fmt.Errorf("%w: non-IP or Ethernet PDN address spare octets are nonzero", ErrInvalidValue)
		}
	default:
		return fmt.Errorf("%w: PDN type %#02x is reserved", ErrInvalidValue, decoded.Type)
	}
	*f = decoded
	return nil
}

func invalidPDNAddressLength(got, want int) error {
	return fmt.Errorf("%w: PDN address has %d value octets, want %d", ErrInvalidLength, got, want)
}
