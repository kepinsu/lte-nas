package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// ExtendedServiceRequest represents an EXTENDED SERVICE REQUEST message as
// defined in TS 24.301 clause 8.2.15 (Table 8.2.15.1, direction UE to network).
//
// The positional mandatory part encodes Service type and NAS key set identifier
// as half-octets (packed in a single wire byte), followed by M-TMSI without an
// IEI tag. Tagged optional IEs (e.g., CSFB response at 0x16, device properties
// at 0x0D) are resolved via message context.
type ExtendedServiceRequest struct {
	EmmHeader
	ServiceType            *ie.IE // positional type-1 (9.9.3.27), half-octet with NASKeySetIdentifier
	NASKeySetIdentifier    *ie.IE // positional type-1 (9.9.3.21), half-octet with ServiceType
	MobileIdentity         *ie.IE // positional LV-E type-4 (9.9.2.3), M-TMSI
	CSFBResponse           *ie.IE // tagged optional TV type-1 (9.9.3.5), tag 0x16
	EPSBearerContextStatus *ie.IE // tagged optional TLV type-4 (9.9.2.1), tag 0x57
	DeviceProperties       *ie.IE // tagged optional TV type-1 (9.9.2.0A), tag 0x0D
	UeRequestType          *ie.IE // tagged optional TLV type-4 (9.9.3.65), tag 0x29
	PagingRestriction      *ie.IE // tagged optional TLV type-4 (9.9.3.66), tag 0x28
	IEs                    []*ie.IE
}

// NewExtendedServiceRequest creates an EXTENDED SERVICE REQUEST message with the
// mandatory positional IEs from TS 24.301 clause 8.2.15, followed by optional
// tagged IEs classified by their Kind. First occurrence wins; duplicates and
// unrecognized IEs are appended to the IEs field.
func NewExtendedServiceRequest(serviceType, nskid, mobileIdentity *ie.IE,
	ies ...*ie.IE,
) *ExtendedServiceRequest {
	if serviceType == nil || nskid == nil || mobileIdentity == nil {
		panic("ExtendedServiceRequest positional IEs must not be nil")
	}

	m := &ExtendedServiceRequest{
		ServiceType:         serviceType,
		NASKeySetIdentifier: nskid,
		MobileIdentity:      mobileIdentity,
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
		case ie.CSFBResponse:
			if m.CSFBResponse == nil {
				m.CSFBResponse = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.EPSBearerContextStatus:
			if m.EPSBearerContextStatus == nil {
				m.EPSBearerContextStatus = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.DeviceProperties:
			if m.DeviceProperties == nil {
				m.DeviceProperties = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.UERequestType:
			if m.UeRequestType == nil {
				m.UeRequestType = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.PagingRestriction:
			if m.PagingRestriction == nil {
				m.PagingRestriction = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the EXTENDED SERVICE REQUEST message, or 0
// for a nil receiver. The result includes header, packed half-octets (1 byte),
// positional IEs and all named tagged IEs in spec-table order.
func (m *ExtendedServiceRequest) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	// Half-octets: combined (1 byte) or individual.
	if m.ServiceType != nil && m.NASKeySetIdentifier != nil {
		n++ // packed half-octet
	} else {
		if m.ServiceType != nil {
			n++
		}
		if m.NASKeySetIdentifier != nil {
			n++
		}
	}

	if m.MobileIdentity != nil {
		n += m.MobileIdentity.MarshalLenWithoutTag()
	}

	if m.CSFBResponse != nil {
		n += m.CSFBResponse.MarshalLen()
	}
	if m.EPSBearerContextStatus != nil {
		n += m.EPSBearerContextStatus.MarshalLen()
	}
	if m.DeviceProperties != nil {
		n += m.DeviceProperties.MarshalLen()
	}
	if m.UeRequestType != nil {
		n += m.UeRequestType.MarshalLen()
	}
	if m.PagingRestriction != nil {
		n += m.PagingRestriction.MarshalLen()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the EXTENDED SERVICE REQUEST message into b. The nil receiver
// returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too
// short. Header, packed half-octets, positional IEs and named tagged IEs follow
// TS 24.301 clause 8.2.15 table order.
func (m *ExtendedServiceRequest) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: ExtendedServiceRequest is nil", ie.ErrNilIE)
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

	// Packed half-octets (service type low nibble, NAS key set ID high nibble).
	if m.ServiceType != nil && m.NASKeySetIdentifier != nil {
		hb, err := ie.MarshalHalfOctets(m.ServiceType, m.NASKeySetIdentifier)
		if err != nil {
			return fmt.Errorf("half-octets: %w", err)
		}
		if len(b) < len(hb) {
			return io.ErrShortBuffer
		}
		copy(b, hb)
		b = b[len(hb):]
	} else {
		if m.ServiceType != nil {
			if err := m.ServiceType.MarshalToWithoutTag(b); err != nil {
				return fmt.Errorf("ServiceType: %w", err)
			}
			b = b[1:]
		}
		if m.NASKeySetIdentifier != nil {
			if err := m.NASKeySetIdentifier.MarshalToWithoutTag(b); err != nil {
				return fmt.Errorf("NASKeySetIdentifier: %w", err)
			}
			b = b[1:]
		}
	}

	if m.MobileIdentity != nil {
		if err := m.MobileIdentity.MarshalToWithoutTag(b); err != nil {
			return fmt.Errorf("MobileIdentity: %w", err)
		}
		n := m.MobileIdentity.MarshalLenWithoutTag()
		if len(b) < n {
			return io.ErrShortBuffer
		}
		b = b[n:]
	}

	// Named tagged IEs in spec-table order.
	if m.CSFBResponse != nil {
		w, err := writeIE(b, 0, m.CSFBResponse)
		if err != nil {
			return fmt.Errorf("CSFBResponse: %w", err)
		}
		b = b[w:]
	}
	if m.EPSBearerContextStatus != nil {
		w, err := writeIE(b, 0, m.EPSBearerContextStatus)
		if err != nil {
			return fmt.Errorf("EPSBearerContextStatus: %w", err)
		}
		b = b[w:]
	}
	if m.DeviceProperties != nil {
		w, err := writeIE(b, 0, m.DeviceProperties)
		if err != nil {
			return fmt.Errorf("DeviceProperties: %w", err)
		}
		b = b[w:]
	}
	if m.UeRequestType != nil {
		w, err := writeIE(b, 0, m.UeRequestType)
		if err != nil {
			return fmt.Errorf("UeRequestType: %w", err)
		}
		b = b[w:]
	}
	if m.PagingRestriction != nil {
		w, err := writeIE(b, 0, m.PagingRestriction)
		if err != nil {
			return fmt.Errorf("PagingRestriction: %w", err)
		}
		b = b[w:]
	}

	// Extra tagged IEs.
	for _, i := range m.IEs {
		w, err := writeIE(b, 0, i)
		if err != nil {
			return fmt.Errorf("extra IE: %w", err)
		}
		b = b[w:]
	}

	return nil
}

// UnmarshalBinary decodes an EXTENDED SERVICE REQUEST message from b. The receiver
// is left unchanged on error. Header, half-octets via ParseHalfOctets, positional
// IEs via ie.ParseWithoutTag with prototypes, tagged optional IEs resolved via
// expected prototypes.
func (m *ExtendedServiceRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: ExtendedServiceRequest is nil", ie.ErrNilIE)
	}

	var tmp ExtendedServiceRequest

	// Decode header (2 octets for plain).
	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	pos := 0

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	// Positional mandatory: ServiceType + NASKeySetIdentifier as half-octets.
	lProto := ie.NewPositionalServiceType(0)
	hProto := ie.NewPositionalNASKeySetIdentifier(0)
	lIE, hIE, err := ie.ParseHalfOctets(pay[pos:], lProto, hProto)
	if err != nil {
		return fmt.Errorf("half-octets: %w", err)
	}
	tmp.ServiceType = lIE
	tmp.NASKeySetIdentifier = hIE

	pos++ // consumed 1 byte for packed half-octets

	// Positional mandatory: M-TMSI (generic MobileIdentity, clause 9.9.2.3).
	tmp.MobileIdentity, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalMobileIdentity([]byte{0}))
	if err != nil {
		return fmt.Errorf("MobileIdentity: %w", err)
	}
	pos += tmp.MobileIdentity.MarshalLenWithoutTag()

	// Tagged optional IEs — prototypes resolve ambiguous tags to their correct types.
	tagged, err := ie.ParseMultiIEs(pay[pos:])
	if err != nil {
		return fmt.Errorf("optional: %w", err)
	}

	for _, i := range tagged {
		switch i.Kind {
		case ie.CSFBResponse:
			if tmp.CSFBResponse == nil {
				tmp.CSFBResponse = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.EPSBearerContextStatus:
			if tmp.EPSBearerContextStatus == nil {
				tmp.EPSBearerContextStatus = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.DeviceProperties:
			if tmp.DeviceProperties == nil {
				tmp.DeviceProperties = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.UERequestType:
			if tmp.UeRequestType == nil {
				tmp.UeRequestType = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.PagingRestriction:
			if tmp.PagingRestriction == nil {
				tmp.PagingRestriction = i
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

// MessageType returns the EXTENDED SERVICE REQUEST message type constant (0x4C).
func (m *ExtendedServiceRequest) MessageType() uint8 {
	return EMMExtendedServiceRequest
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *ExtendedServiceRequest) MessageTypeName() string {
	return "Extended service request"
}

// IsRequest returns true — EXTENDED SERVICE REQUEST is sent UE to network.
func (m *ExtendedServiceRequest) IsRequest() bool {
	return true
}

// ExtendedServiceRequest implements the Message interface.
var _ Message = (*ExtendedServiceRequest)(nil)

// Header returns the EMM security header embedded in ExtendedServiceRequest.
func (m *ExtendedServiceRequest) Header() *EmmHeader {
	return &m.EmmHeader
}
