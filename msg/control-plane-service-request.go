package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// ControlPlaneServiceRequest represents a CONTROL PLANE SERVICE REQUEST message
// as defined in TS 24.301 clause 8.2.33 (Table 8.2.33.1, direction UE to network).
//
// The positional mandatory part encodes Control plane service type and NAS key set
// identifier as half-octets (packed in a single wire byte). Tagged optional IEs
// include ESM message container (0x78), NAS message container (0x67) and others.
// May carry security header type 0x05 (partially ciphered) per Table 9.3.1 NOTE 4.
type ControlPlaneServiceRequest struct {
	EmmHeader
	ControlPlaneServiceType *ie.IE // positional type-1 (9.9.3.47), half-octet with NASKeySetIdentifier
	NASKeySetIdentifier     *ie.IE // positional type-1 (9.9.3.21), half-octet with ControlPlaneServiceType
	EsmMessageContainer     *ie.IE // tagged optional TLV-E type-6 (9.9.3.15), tag 0x78
	NasMessageContainer     *ie.IE // tagged optional TLV type-4 (9.9.3.22), tag 0x67
	EPSBearerContextStatus  *ie.IE // tagged optional TLV type-4 (9.9.2.1), tag 0x57
	DeviceProperties        *ie.IE // tagged optional TV type-1 (9.9.2.0A), tag 0x0D
	UeRequestType           *ie.IE // tagged optional TLV type-4 (9.9.3.65), tag 0x29
	PagingRestriction       *ie.IE // tagged optional TLV type-4 (9.9.3.66), tag 0x28
	IEs                     []*ie.IE
}

// NewControlPlaneServiceRequest creates a CONTROL PLANE SERVICE REQUEST message
// with the mandatory positional IEs from TS 24.301 clause 8.2.33, followed by
// optional tagged IEs classified by their Kind. First occurrence wins; duplicates
// and unrecognized IEs are appended to the IEs field.
func NewControlPlaneServiceRequest(controlPlaneServiceType, nskid *ie.IE,
	ies ...*ie.IE,
) *ControlPlaneServiceRequest {
	if controlPlaneServiceType == nil || nskid == nil {
		panic("ControlPlaneServiceRequest positional IEs must not be nil")
	}

	m := &ControlPlaneServiceRequest{
		ControlPlaneServiceType: controlPlaneServiceType,
		NASKeySetIdentifier:     nskid,
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
		case ie.ESMMessageContainer:
			if m.EsmMessageContainer == nil {
				m.EsmMessageContainer = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.NASMessageContainer:
			if m.NasMessageContainer == nil {
				m.NasMessageContainer = i
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

// MarshalLen returns the wire size of the CONTROL PLANE SERVICE REQUEST message,
// or 0 for a nil receiver. The result includes header, packed half-octets (1 byte),
// and all named tagged IEs in spec-table order.
func (m *ControlPlaneServiceRequest) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	// Half-octets: combined (1 byte) or individual.
	if m.ControlPlaneServiceType != nil && m.NASKeySetIdentifier != nil {
		n++ // packed half-octet
	} else {
		if m.ControlPlaneServiceType != nil {
			n++
		}
		if m.NASKeySetIdentifier != nil {
			n++
		}
	}

	if m.EsmMessageContainer != nil {
		n += m.EsmMessageContainer.MarshalLen()
	}
	if m.NasMessageContainer != nil {
		n += m.NasMessageContainer.MarshalLen()
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

// MarshalTo encodes the CONTROL PLANE SERVICE REQUEST message into b. The nil
// receiver returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b
// is too short. Header, packed half-octets and named tagged IEs follow TS 24.301
// clause 8.2.33 table order.
func (m *ControlPlaneServiceRequest) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: ControlPlaneServiceRequest is nil", ie.ErrNilIE)
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

	// Packed half-octets (control plane service type low nibble, NAS key set ID high nibble).
	if m.ControlPlaneServiceType != nil && m.NASKeySetIdentifier != nil {
		hb, err := ie.MarshalHalfOctets(m.ControlPlaneServiceType, m.NASKeySetIdentifier)
		if err != nil {
			return fmt.Errorf("half-octets: %w", err)
		}
		if len(b) < len(hb) {
			return io.ErrShortBuffer
		}
		copy(b, hb)
		b = b[len(hb):]
	} else {
		if m.ControlPlaneServiceType != nil {
			if err := m.ControlPlaneServiceType.MarshalToWithoutTag(b); err != nil {
				return fmt.Errorf("ControlPlaneServiceType: %w", err)
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

	// Named tagged IEs in spec-table order.
	if m.EsmMessageContainer != nil {
		w, err := writeIE(b, 0, m.EsmMessageContainer)
		if err != nil {
			return fmt.Errorf("EsmMessageContainer: %w", err)
		}
		b = b[w:]
	}
	if m.NasMessageContainer != nil {
		w, err := writeIE(b, 0, m.NasMessageContainer)
		if err != nil {
			return fmt.Errorf("NasMessageContainer: %w", err)
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

// UnmarshalBinary decodes a CONTROL PLANE SERVICE REQUEST message from b. The
// receiver is left unchanged on error. Header, half-octets via ParseHalfOctets,
// tagged optional IEs resolved via expected prototypes.
func (m *ControlPlaneServiceRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: ControlPlaneServiceRequest is nil", ie.ErrNilIE)
	}

	var tmp ControlPlaneServiceRequest

	// Decode header (2 octets for plain, 7 for integrity-protected).
	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	pos := 0

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	// Positional mandatory: ControlPlaneServiceType + NASKeySetIdentifier as half-octets.
	lProto := ie.NewPositionalControlPlaneServiceType(0)
	hProto := ie.NewPositionalNASKeySetIdentifier(0)
	lIE, hIE, err := ie.ParseHalfOctets(pay[pos:], lProto, hProto)
	if err != nil {
		return fmt.Errorf("half-octets: %w", err)
	}
	tmp.ControlPlaneServiceType = lIE
	tmp.NASKeySetIdentifier = hIE

	pos++ // consumed 1 byte for packed half-octets

	// Tagged optional IEs — prototypes resolve ambiguous tags to their correct types.
	tagged, err := ie.ParseMultiIEs(pay[pos:])
	if err != nil {
		return fmt.Errorf("optional: %w", err)
	}

	for _, i := range tagged {
		switch i.Kind {
		case ie.ESMMessageContainer:
			if tmp.EsmMessageContainer == nil {
				tmp.EsmMessageContainer = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.NASMessageContainer:
			if tmp.NasMessageContainer == nil {
				tmp.NasMessageContainer = i
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

// MessageType returns the CONTROL PLANE SERVICE REQUEST message type constant (0x4D).
func (m *ControlPlaneServiceRequest) MessageType() uint8 {
	return EMMControlPlaneServiceRequest
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *ControlPlaneServiceRequest) MessageTypeName() string {
	return "Control plane service request"
}

// IsRequest returns true — CONTROL PLANE SERVICE REQUEST is sent UE to network.
func (m *ControlPlaneServiceRequest) IsRequest() bool {
	return true
}

// ControlPlaneServiceRequest implements the Message interface.
var _ Message = (*ControlPlaneServiceRequest)(nil)

// Header returns the EMM security header embedded in ControlPlaneServiceRequest.
func (m *ControlPlaneServiceRequest) Header() *EmmHeader {
	return &m.EmmHeader
}
