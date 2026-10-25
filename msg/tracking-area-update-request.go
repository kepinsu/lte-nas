package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// TrackingAreaUpdateRequest represents a TRACKING AREA UPDATE REQUEST message as
// defined in TS 24.301 clause 8.2.29 (Table 8.2.29.1, direction UE to network).
//
// The positional mandatory part encodes EPS update type and NAS key set identifier
// as packed half-octets, followed by Old GUTI (EPS mobile identity) without IEI
// tags. Tagged optional IEs are resolved via message context prototypes, with
// ambiguous tags (0x32, 0x58) handled per-table.
type TrackingAreaUpdateRequest struct {
	EmmHeader
	EPSUpdateType       *ie.IE // positional type-1 (9.9.3.14), half-octet with NASKeySetIdentifier
	NASKeySetIdentifier *ie.IE // positional type-1 (9.9.3.21), half-octet with EPSUpdateType
	OldGUTI             *ie.IE // positional LV type-4 (9.9.3.12)
	// Tagged optional IEs resolved via prototypes in the spec table.
	DRXParameter            *ie.IE // tagged optional TV type-3, 0x5c ambiguous (9.9.3.8)
	GUTI                    *ie.IE // tagged optional TLV, 0x50 (9.9.3.12)
	TAIList                 *ie.IE // tagged optional TLV, 0x54 (9.9.3.33)
	EMMCause                *ie.IE // tagged optional TV type-3, 0x53 (9.9.3.9)
	EPSBearerContextStatus  *ie.IE // tagged optional TLV, 0x57 (9.9.2.1)
	UENetworkCapability     *ie.IE // tagged optional TLV, 0x58 ambiguous (9.9.3.34)
	ESMMessageContainer     *ie.IE // tagged optional TLV-E, 0x78 (9.9.3.15)
	LastVisitedRegisteredTA *ie.IE // tagged optional TV type-3, 0x52 (9.9.3.32)
	IEs                     []*ie.IE
}

// NewTrackingAreaUpdateRequest creates a TRACKING AREA UPDATE REQUEST message with
// the mandatory positional IEs from TS 24.301 clause 8.2.29, followed by optional
// tagged IEs classified by their Kind. Fields for recognized tags are resolved
// first; duplicates and unrecognized IEs are appended to the IEs field.
func NewTrackingAreaUpdateRequest(epsUpdateType, nskid, oldGUTI *ie.IE,
	ies ...*ie.IE,
) *TrackingAreaUpdateRequest {
	if epsUpdateType == nil || nskid == nil || oldGUTI == nil {
		panic("TrackingAreaUpdateRequest positional IEs must not be nil")
	}

	m := &TrackingAreaUpdateRequest{
		EPSUpdateType:       epsUpdateType,
		NASKeySetIdentifier: nskid,
		OldGUTI:             oldGUTI,
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
		case ie.DRXParameter:
			if m.DRXParameter == nil {
				m.DRXParameter = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.EPSMobileIdentity:
			if m.GUTI == nil {
				m.GUTI = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.TAIList:
			if m.TAIList == nil {
				m.TAIList = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.EMMCause:
			if m.EMMCause == nil {
				m.EMMCause = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.EPSBearerContextStatus:
			if m.EPSBearerContextStatus == nil {
				m.EPSBearerContextStatus = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.UENetworkCapability:
			if m.UENetworkCapability == nil {
				m.UENetworkCapability = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.ESMMessageContainer:
			if m.ESMMessageContainer == nil {
				m.ESMMessageContainer = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.TrackingAreaIdentity:
			if m.LastVisitedRegisteredTA == nil {
				m.LastVisitedRegisteredTA = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the TRACKING AREA UPDATE REQUEST message, or
// 0 for a nil receiver. The result includes header, packed half-octets (1 byte),
// positional IEs and all named tagged IEs in spec-table order.
func (m *TrackingAreaUpdateRequest) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	// Packed half-octets (EPS update type low nibble, NAS key set ID high nibble).
	if m.EPSUpdateType != nil && m.NASKeySetIdentifier != nil {
		n++ // packed half-octet
	} else {
		if m.EPSUpdateType != nil {
			n++
		}
		if m.NASKeySetIdentifier != nil {
			n++
		}
	}

	if m.OldGUTI != nil {
		n += m.OldGUTI.MarshalLenWithoutTag()
	}

	// Named tagged IEs in spec-table order.
	if m.DRXParameter != nil {
		n += m.DRXParameter.MarshalLen()
	}
	if m.GUTI != nil {
		n += m.GUTI.MarshalLen()
	}
	if m.TAIList != nil {
		n += m.TAIList.MarshalLen()
	}
	if m.EMMCause != nil {
		n += m.EMMCause.MarshalLen()
	}
	if m.EPSBearerContextStatus != nil {
		n += m.EPSBearerContextStatus.MarshalLen()
	}
	if m.UENetworkCapability != nil {
		n += m.UENetworkCapability.MarshalLen()
	}
	if m.ESMMessageContainer != nil {
		n += m.ESMMessageContainer.MarshalLen()
	}
	if m.LastVisitedRegisteredTA != nil {
		n += m.LastVisitedRegisteredTA.MarshalLen()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the TRACKING AREA UPDATE REQUEST message into b. The nil
// receiver returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b
// is too short. Header, packed half-octets, positional IEs, named tagged IEs
// and extras follow TS 24.301 clause 8.2.29 table order.
func (m *TrackingAreaUpdateRequest) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: TrackingAreaUpdateRequest is nil", ie.ErrNilIE)
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

	// Packed half-octets (EPS update type low nibble, NAS key set ID high nibble).
	if m.EPSUpdateType != nil && m.NASKeySetIdentifier != nil {
		hb, err := ie.MarshalHalfOctets(m.EPSUpdateType, m.NASKeySetIdentifier)
		if err != nil {
			return fmt.Errorf("half-octets: %w", err)
		}
		if len(b) < n+len(hb) {
			return io.ErrShortBuffer
		}
		copy(b[n:n+len(hb)], hb)
		n += len(hb)
	} else {
		if m.EPSUpdateType != nil {
			if err := m.EPSUpdateType.MarshalToWithoutTag(b[n:]); err != nil {
				return fmt.Errorf("EPSUpdateType: %w", err)
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
	if m.OldGUTI != nil {
		if err := m.OldGUTI.MarshalToWithoutTag(b[n:]); err != nil {
			return fmt.Errorf("OldGUTI: %w", err)
		}
		n += m.OldGUTI.MarshalLenWithoutTag()
	}

	// Named tagged IEs in spec-table order.
	if m.DRXParameter != nil {
		w, err := writeIE(b, n, m.DRXParameter)
		if err != nil {
			return fmt.Errorf("DRXParameter: %w", err)
		}
		n += w
	}
	if m.GUTI != nil {
		w, err := writeIE(b, n, m.GUTI)
		if err != nil {
			return fmt.Errorf("GUTI: %w", err)
		}
		n += w
	}
	if m.TAIList != nil {
		w, err := writeIE(b, n, m.TAIList)
		if err != nil {
			return fmt.Errorf("TAIList: %w", err)
		}
		n += w
	}
	if m.EMMCause != nil {
		w, err := writeIE(b, n, m.EMMCause)
		if err != nil {
			return fmt.Errorf("EMMCause: %w", err)
		}
		n += w
	}
	if m.EPSBearerContextStatus != nil {
		w, err := writeIE(b, n, m.EPSBearerContextStatus)
		if err != nil {
			return fmt.Errorf("EPSBearerContextStatus: %w", err)
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
	if m.ESMMessageContainer != nil {
		w, err := writeIE(b, n, m.ESMMessageContainer)
		if err != nil {
			return fmt.Errorf("ESMMessageContainer: %w", err)
		}
		n += w
	}
	if m.LastVisitedRegisteredTA != nil {
		w, err := writeIE(b, n, m.LastVisitedRegisteredTA)
		if err != nil {
			return fmt.Errorf("LastVisitedRegisteredTA: %w", err)
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

// UnmarshalBinary decodes a TRACKING AREA UPDATE REQUEST message from b. The
// receiver is left unchanged on error. Header, half-octets via ParseHalfOctets,
// positional IEs via ParseWithoutTag prototypes, tagged optional IEs resolved
// via ambiguous tag prototypes.
func (m *TrackingAreaUpdateRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: TrackingAreaUpdateRequest is nil", ie.ErrNilIE)
	}

	var tmp TrackingAreaUpdateRequest

	// Decode header (2 octets for plain).
	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	pos := 0

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	// Positional mandatory: half-octets (EPSUpdateType low nibble,
	// NASKeySetIdentifier high nibble) packed in one wire byte.
	lProto := ie.NewPositionalEPSAttachType(0)
	hProto := ie.NewPositionalNASKeySetIdentifier(0)
	lIE, hIE, err := ie.ParseHalfOctets(pay[pos:], lProto, hProto)
	if err != nil {
		return fmt.Errorf("half-octets: %w", err)
	}
	tmp.EPSUpdateType = lIE
	tmp.NASKeySetIdentifier = hIE

	pos++ // consumed 1 byte for packed half-octets

	// Remaining positional IEs via ParseWithoutTag with prototypes.
	tmp.OldGUTI, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalEPSMobileIdentity([]byte{0, 0}))
	if err != nil {
		return fmt.Errorf("OldGUTI: %w", err)
	}
	pos += tmp.OldGUTI.MarshalLenWithoutTag()

	// Tagged optional IEs — prototypes resolve ambiguous tags per TS 24.301
	// Table 8.2.29.1: 0x32 LLC SAPI / N1 UE NC, 0x58 ESM cause / UE NC.
	// Tag 0x5c DRX parameter uses constructor prototype; format Three (TV).
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
		case ie.EPSQualityOfService:
			if tmp.DRXParameter == nil {
				tmp.DRXParameter = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.EPSMobileIdentity:
			if tmp.GUTI == nil {
				tmp.GUTI = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.TAIList:
			if tmp.TAIList == nil {
				tmp.TAIList = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.EMMCause:
			if tmp.EMMCause == nil {
				tmp.EMMCause = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.EPSBearerContextStatus:
			if tmp.EPSBearerContextStatus == nil {
				tmp.EPSBearerContextStatus = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.UENetworkCapability:
			if tmp.UENetworkCapability == nil {
				tmp.UENetworkCapability = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.ESMMessageContainer:
			if tmp.ESMMessageContainer == nil {
				tmp.ESMMessageContainer = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.TrackingAreaIdentity:
			if tmp.LastVisitedRegisteredTA == nil {
				tmp.LastVisitedRegisteredTA = i
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

// MessageType returns the TRACKING AREA UPDATE REQUEST message type constant (0x48).
func (m *TrackingAreaUpdateRequest) MessageType() uint8 {
	return EMMTrackingAreaUpdateRequest
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *TrackingAreaUpdateRequest) MessageTypeName() string {
	return "Tracking area update request"
}

// IsRequest returns true — TRACKING AREA UPDATE REQUEST is sent UE to network.
func (m *TrackingAreaUpdateRequest) IsRequest() bool {
	return true
}

// TrackingAreaUpdateRequest implements the Message interface.
var _ Message = (*TrackingAreaUpdateRequest)(nil)

// Header returns the EMM security header embedded in TrackingAreaUpdateRequest.
func (m *TrackingAreaUpdateRequest) Header() *EmmHeader {
	return &m.EmmHeader
}
