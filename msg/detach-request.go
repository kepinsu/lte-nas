package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// DetachRequest represents a DETACH REQUEST message as defined in TS 24.301
// clause 8.2.11 (Table 8.2.11.1, direction UE to network).
//
// The positional mandatory part encodes detach type and NAS key set identifier
// as packed half-octets (wire byte), followed by EPS mobile identity without IEI
// tags. Tagged optional IEs (e.g. ESM message container tag 0x78, UE network
// capability tag 0x58) are resolved via message context prototypes.
type DetachRequest struct {
	EmmHeader
	DetachType          *ie.IE // positional type-1 (9.9.3.7), half-octet with NASKeySetIdentifier
	NASKeySetIdentifier *ie.IE // positional type-1 (9.9.3.21), half-octet with DetachType
	EPSMobileIdentity   *ie.IE // positional LV type-4 (9.9.3.12)
	// Tagged optional IEs resolved via prototypes in the spec table.
	ESMMessageContainer *ie.IE // tagged optional TLV-E, 0x78 (9.9.3.15), pass-through
	UENetworkCapability *ie.IE // tagged optional TLV, 0x58 ambiguous (9.9.3.34)
	IEs                 []*ie.IE
}

// NewDetachRequest creates a DETACH REQUEST message with the mandatory positional
// IEs from TS 24.301 clause 8.2.11, followed by optional tagged IEs classified
// by their Kind. Fields for recognized tags are resolved first; duplicates and
// unrecognized IEs are appended to the IEs field.
func NewDetachRequest(detachType, nskid, epsMobileIdentity *ie.IE,
	ies ...*ie.IE,
) *DetachRequest {
	if detachType == nil || nskid == nil || epsMobileIdentity == nil {
		panic("DetachRequest positional IEs must not be nil")
	}

	m := &DetachRequest{
		DetachType:          detachType,
		NASKeySetIdentifier: nskid,
		EPSMobileIdentity:   epsMobileIdentity,
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
		case ie.ESMMessageContainer:
			if m.ESMMessageContainer == nil {
				m.ESMMessageContainer = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.UENetworkCapability:
			if m.UENetworkCapability == nil {
				m.UENetworkCapability = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the DETACH REQUEST message, or 0 for a
// nil receiver. The result includes header, packed half-octets (1 byte),
// positional IEs and all named tagged IEs in spec-table order.
func (m *DetachRequest) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	// Packed half-octets (detach type low nibble, NAS key set ID high nibble).
	if m.DetachType != nil && m.NASKeySetIdentifier != nil {
		n++ // packed half-octet
	} else {
		if m.DetachType != nil {
			n++
		}
		if m.NASKeySetIdentifier != nil {
			n++
		}
	}

	if m.EPSMobileIdentity != nil {
		n += m.EPSMobileIdentity.MarshalLenWithoutTag()
	}

	if m.ESMMessageContainer != nil {
		n += m.ESMMessageContainer.MarshalLen()
	}
	if m.UENetworkCapability != nil {
		n += m.UENetworkCapability.MarshalLen()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the DETACH REQUEST message into b. The nil receiver
// returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too
// short. Header, packed half-octets, positional IEs, named tagged IEs and
// extras follow TS 24.301 clause 8.2.11 table order.
func (m *DetachRequest) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: DetachRequest is nil", ie.ErrNilIE)
	}

	hlen := m.EmmHeader.MarshalLen()
	if hlen == 0 {
		return fmt.Errorf("%w: EmmHeader MarshalLen returned 0", ie.ErrNilIE)
	}

	if len(b) < hlen {
		return io.ErrShortBuffer
	}

	if err := m.EmmHeader.MarshalTo(b); err != nil {
		return fmt.Errorf("header: %w", err)
	}

	b = b[hlen:]
	n := 0

	// Packed half-octets (detach type low nibble, NAS key set ID high nibble).
	if m.DetachType != nil && m.NASKeySetIdentifier != nil {
		hb, err := ie.MarshalHalfOctets(m.DetachType, m.NASKeySetIdentifier)
		if err != nil {
			return fmt.Errorf("half-octets: %w", err)
		}
		if len(b) < n+len(hb) {
			return io.ErrShortBuffer
		}
		copy(b[n:n+len(hb)], hb)
		n += len(hb)
	} else {
		if m.DetachType != nil {
			if err := m.DetachType.MarshalToWithoutTag(b[n:]); err != nil {
				return fmt.Errorf("DetachType: %w", err)
			}
			n++
		}
		if m.NASKeySetIdentifier != nil {
			if err := m.NASKeySetIdentifier.MarshalToWithoutTag(b[n:]); err != nil {
				return fmt.Errorf("NASKeySetIdentifier: %w", err)
			}
			n++
		}
	}

	// Positional IEs in spec-table order.
	if m.EPSMobileIdentity != nil {
		if err := m.EPSMobileIdentity.MarshalToWithoutTag(b[n:]); err != nil {
			return fmt.Errorf("EPSMobileIdentity: %w", err)
		}
		n += m.EPSMobileIdentity.MarshalLenWithoutTag()
	}

	// Named tagged IEs in spec-table order.
	if m.ESMMessageContainer != nil {
		w, err := writeIE(b, n, m.ESMMessageContainer)
		if err != nil {
			return fmt.Errorf("ESMMessageContainer: %w", err)
		}
		n += w
	}
	if m.UENetworkCapability != nil {
		w, err := writeIE(b, n, m.UENetworkCapability)
		if err != nil {
			return fmt.Errorf("UENetworkCapability: %w", err)
		}
		n += w
	}

	// Extra tagged IEs.
	for _, i := range m.IEs {
		w, err := writeIE(b, n, i)
		if err != nil {
			return fmt.Errorf("extra IE: %w", err)
		}
		n += w
	}

	return nil
}

// UnmarshalBinary decodes a DETACH REQUEST message from b. The receiver is left
// unchanged on error. Header, half-octets via ParseHalfOctets, positional IEs
// via ParseWithoutTag prototypes, tagged optional IEs resolved via prototypes.
func (m *DetachRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: DetachRequest is nil", ie.ErrNilIE)
	}

	var tmp DetachRequest

	// Decode header (2 octets for plain).
	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	pos := 0

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	// Positional mandatory: half-octets (DetachType low nibble, NASKeySetIdentifier
	// high nibble) packed in one wire byte. ParseHalfOctets extracts both values.
	lProto := ie.NewPositionalDetachType(0)
	hProto := ie.NewPositionalNASKeySetIdentifier(0)
	lIE, hIE, err := ie.ParseHalfOctets(pay[pos:], lProto, hProto)
	if err != nil {
		return fmt.Errorf("half-octets: %w", err)
	}
	tmp.DetachType = lIE
	tmp.NASKeySetIdentifier = hIE

	pos++ // consumed 1 byte for packed half-octets

	// Remaining positional IEs via ParseWithoutTag with prototypes.
	tmp.EPSMobileIdentity, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalEPSMobileIdentity([]byte{0, 0}))
	if err != nil {
		return fmt.Errorf("EPSMobileIdentity: %w", err)
	}
	pos += tmp.EPSMobileIdentity.MarshalLenWithoutTag()

	// Tagged optional IEs — prototypes resolve ambiguous tag 0x58 to UE NC
	// and pass-through ESM message container at 0x78.
	tagged, err := ie.ParseMultiIEs(pay[pos:])
	if err != nil {
		return fmt.Errorf("optional: %w", err)
	}

	for _, i := range tagged {
		switch i.Kind {
		case ie.ESMMessageContainer:
			if tmp.ESMMessageContainer == nil {
				tmp.ESMMessageContainer = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.UENetworkCapability:
			if tmp.UENetworkCapability == nil {
				tmp.UENetworkCapability = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		default:
			tmp.IEs = append(tmp.IEs, i)
		}
	}

	*m = tmp
	return nil
}

// MessageType returns the DETACH REQUEST message type constant (0x45).
func (m *DetachRequest) MessageType() uint8 {
	return EMMDetachRequest
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *DetachRequest) MessageTypeName() string {
	return "Detach request"
}

// IsRequest returns true — DETACH REQUEST is sent UE to network.
func (m *DetachRequest) IsRequest() bool {
	return true
}

// DetachRequest implements the Message interface.
var _ Message = (*DetachRequest)(nil)

// Header returns the EMM security header embedded in DetachRequest.
func (m *DetachRequest) Header() *EmmHeader {
	return &m.EmmHeader
}
