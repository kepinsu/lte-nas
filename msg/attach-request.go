package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// AttachRequest represents an ATTACH REQUEST message as defined in TS 24.301
// clause 8.2.4 (Table 8.2.4.1, direction UE to network).
//
// The positional mandatory part encodes EPS attach type and NAS key set
// identifier as half-octets (packed in a single wire byte), followed by
// EPS mobile identity, UE network capability and ESM message container
// without IEI tags. Tagged optional IEs (e.g. DRX parameter, tag 0x5C) are
// resolved via message context.
type AttachRequest struct {
	EmmHeader
	EPSAttachType       *ie.IE // positional type-1 (9.9.3.11), half-octet with NASKeySetIdentifier
	NASKeySetIdentifier *ie.IE // positional type-1 (9.9.3.21), half-octet with EPSAttachType
	EPSMobileIdentity   *ie.IE // positional LV type-4 (9.9.3.12)
	UENetworkCapability *ie.IE // positional LV type-4 (9.9.3.34)
	ESMMessageContainer *ie.IE // positional LV-E type-6 (9.9.3.15), pass-through
	DRXParameter        *ie.IE // tagged optional TV type-1, ambiguous tag 0x5C (9.9.3.8)
	IEs                 []*ie.IE
}

// NewAttachRequest creates an ATTACH REQUEST message with the mandatory
// positional IEs from TS 24.301 clause 8.2.4, followed by optional tagged
// IEs classified by their Kind. DRXParameter (Kind 0x5C) is resolved first;
// duplicates and unrecognized IEs are appended to the IEs field.
func NewAttachRequest(epsAttachType, nskid, epsMobileIdentity, ueNetworkCapability,
	esmMsgContainer *ie.IE, ies ...*ie.IE,
) *AttachRequest {
	if epsAttachType == nil || nskid == nil || epsMobileIdentity == nil ||
		ueNetworkCapability == nil || esmMsgContainer == nil {
		panic("AttachRequest positional IEs must not be nil")
	}

	m := &AttachRequest{
		EPSAttachType:       epsAttachType,
		NASKeySetIdentifier: nskid,
		EPSMobileIdentity:   epsMobileIdentity,
		UENetworkCapability: ueNetworkCapability,
		ESMMessageContainer: esmMsgContainer,
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		// Normalize collision tags (0x32, 0x58, 0x5b, 0x5c) to FormatFour TLV.
		if ie.IsAmbiguousTagged(i.Kind) {
			i.Format = ie.IEFormatFour
			i.Length = uint16(len(i.Payload))
		}
		switch i.Kind {
		case ie.DRXParameter:
			if m.DRXParameter == nil {
				m.DRXParameter = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the ATTACH REQUEST message, or 0 for
// a nil receiver. The result includes header, packed half-octets (1 byte),
// positional IEs and all named tagged IEs in spec-table order.
func (m *AttachRequest) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	// Half-octets: combined (1 byte) or individual.
	if m.EPSAttachType != nil && m.NASKeySetIdentifier != nil {
		n++ // packed half-octet
	} else {
		if m.EPSAttachType != nil {
			n++
		}
		if m.NASKeySetIdentifier != nil {
			n++
		}
	}

	if m.EPSMobileIdentity != nil {
		n += m.EPSMobileIdentity.MarshalLen()
	}
	if m.UENetworkCapability != nil {
		n += m.UENetworkCapability.MarshalLen()
	}
	if m.ESMMessageContainer != nil {
		n += m.ESMMessageContainer.MarshalLen()
	}

	if m.DRXParameter != nil {
		n += m.DRXParameter.MarshalLen()
	}
	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the ATTACH REQUEST message into b. The nil receiver
// returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too
// short. Header, packed half-octets (or individual), positional IEs, named
// tagged IEs and extras follow TS 24.301 clause 8.2.4 table order.
func (m *AttachRequest) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: AttachRequest is nil", ie.ErrNilIE)
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

	// Packed half-octets (EPS attach type low nibble, NAS key set ID high nibble).
	if m.EPSAttachType != nil && m.NASKeySetIdentifier != nil {
		hb, err := ie.MarshalHalfOctets(m.EPSAttachType, m.NASKeySetIdentifier)
		if err != nil {
			return fmt.Errorf("half-octets: %w", err)
		}
		if len(b) < n+len(hb) {
			return io.ErrShortBuffer
		}
		copy(b[n:n+len(hb)], hb)
		n += len(hb)
	} else {
		if m.EPSAttachType != nil {
			if err := m.EPSAttachType.MarshalToWithoutTag(b[n:]); err != nil {
				return fmt.Errorf("EPSAttachType: %w", err)
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
		n0 := n
		if err := m.EPSMobileIdentity.MarshalToWithoutTag(b[n:]); err != nil {
			return fmt.Errorf("EPSMobileIdentity: %w", err)
		}
		n += m.EPSMobileIdentity.MarshalLenWithoutTag()
		_ = n0 // unused; prevents unused-variable lint for gap.
	}
	if m.UENetworkCapability != nil {
		n0 := n
		if err := m.UENetworkCapability.MarshalToWithoutTag(b[n:]); err != nil {
			return fmt.Errorf("UENetworkCapability: %w", err)
		}
		n += m.UENetworkCapability.MarshalLenWithoutTag()
		_ = n0
	}
	if m.ESMMessageContainer != nil {
		n0 := n
		if err := m.ESMMessageContainer.MarshalToWithoutTag(b[n:]); err != nil {
			return fmt.Errorf("ESMMessageContainer: %w", err)
		}
		n += m.ESMMessageContainer.MarshalLenWithoutTag()
		_ = n0
	}

	// Named tagged IEs in spec-table order.
	if m.DRXParameter != nil {
		n0 := n
		w, err := writeIE(b, n, m.DRXParameter)
		if err != nil {
			return fmt.Errorf("DRXParameter: %w", err)
		}
		n += w
		_ = n0
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

// UnmarshalBinary decodes an ATTACH REQUEST message from b. The receiver is
// left unchanged on error. Header, half-octets via ParseHalfOctets, positional
// IEs via ie.ParseWithoutTag with prototypes, tagged optional IEs resolved
// via prototypes.
func (m *AttachRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: AttachRequest is nil", ie.ErrNilIE)
	}

	var tmp AttachRequest

	// Decode header (2 octets for plain).
	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	// Positional mandatory IEs. Half-octets (EPSAttachType + NASKeySetIdentifier)
	// are packed into a single wire byte. ParseHalfOctets extracts both values
	// from pay[pos].
	pos := 0

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	lProto := ie.NewPositionalEPSAttachType(0)
	hProto := ie.NewPositionalNASKeySetIdentifier(0)
	lIE, hIE, err := ie.ParseHalfOctets(pay[pos:], lProto, hProto)
	if err != nil {
		return fmt.Errorf("half-octets: %w", err)
	}
	tmp.EPSAttachType = lIE
	tmp.NASKeySetIdentifier = hIE

	pos++ // consumed 1 byte for packed half-octets

	// Remaining positional IEs using ParseWithoutTag with prototypes.
	tmp.EPSMobileIdentity, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalEPSMobileIdentity([]byte{0}))
	if err != nil {
		return fmt.Errorf("EPSMobileIdentity: %w", err)
	}
	pos += tmp.EPSMobileIdentity.MarshalLenWithoutTag()

	tmp.UENetworkCapability, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalUENetworkCapability([]byte{0, 0}))
	if err != nil {
		return fmt.Errorf("UENetworkCapability: %w", err)
	}
	pos += tmp.UENetworkCapability.MarshalLenWithoutTag()

	tmp.ESMMessageContainer, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalESMMessageContainer([]byte{0, 0, 0}))
	if err != nil {
		return fmt.Errorf("ESMMessageContainer: %w", err)
	}
	pos += tmp.ESMMessageContainer.MarshalLenWithoutTag()

	// Tagged optional IEs — prototypes resolve ambiguous tag 0x5C to DRX parameter.
	tagged, err := ie.ParseMultiIEs(pay[pos:])
	if err != nil {
		return fmt.Errorf("optional: %w", err)
	}

	for _, i := range tagged {
		switch i.Kind {
		case ie.DRXParameter:
			if tmp.DRXParameter == nil {
				tmp.DRXParameter = i
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

// MessageType returns the ATTACH REQUEST message type constant (0x41).
func (m *AttachRequest) MessageType() uint8 {
	return EMMAttachRequest
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *AttachRequest) MessageTypeName() string {
	return "Attach request"
}

// IsRequest returns true — ATTACH REQUEST is sent UE to network.
func (m *AttachRequest) IsRequest() bool {
	return true
}

// AttachRequest implements the Message interface.
var _ Message = (*AttachRequest)(nil)

// Header returns the EMM security header embedded in AttachRequest.
func (m *AttachRequest) Header() *EmmHeader {
	return &m.EmmHeader
}
