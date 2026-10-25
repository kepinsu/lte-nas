package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// ActivateDefaultEPSBearerContextRequest represents an ACTIVATE DEFAULT EPS
// BEARER CONTEXT REQUEST message as defined in TS 24.301 clause 8.3.6
// (Table 8.3.6.1, direction network to UE).
//
// The positional mandatory part encodes EPS quality of service, access point
// name and PDN address without IEI tags. Tagged optional IEs (LLC SAPI tag
// 0x32, ESM cause tag 0x58) are resolved via prototypes in ie.ParseMultiIEs.
type ActivateDefaultEPSBearerContextRequest struct {
	EsmHeader
	EPSQoS          *ie.IE // positional LV (9.9.4.3)
	AccessPointName *ie.IE // positional LV (9.9.4.1)
	PDNAddress      *ie.IE // positional LV (9.9.4.9)
	LLCSAPI         *ie.IE // tagged optional type-1 TV (tag 0x32, 9.9.4.7)
	ESMCause        *ie.IE // tagged optional type-3 TV (tag 0x58, 9.9.4.4)
	IEs             []*ie.IE
}

// NewActivateDefaultEPSBearerContextRequest creates an ACTIVATE DEFAULT EPS
// BEARER CONTEXT REQUEST message with the mandatory positional IEs from TS
// 24.301 clause 8.3.6, followed by optional tagged IEs classified by Kind.
// LLCSAPI (tag 0x32) and ESMCause (tag 0x58) are resolved first; duplicates
// and unrecognized IEs are appended to the IEs field.
func NewActivateDefaultEPSBearerContextRequest(ePSQoS, accessPointName, pDNAddress *ie.IE,
	ies ...*ie.IE,
) *ActivateDefaultEPSBearerContextRequest {
	if ePSQoS == nil || accessPointName == nil || pDNAddress == nil {
		panic("ActivateDefaultEPSBearerContextRequest positional IEs must not be nil")
	}

	m := &ActivateDefaultEPSBearerContextRequest{
		EPSQoS:          ePSQoS,
		AccessPointName: accessPointName,
		PDNAddress:      pDNAddress,
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
		case ie.LLCServiceAccessPointIdentifier:
			if m.LLCSAPI == nil {
				m.LLCSAPI = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.ESMCause:
			if m.ESMCause == nil {
				m.ESMCause = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the ACTIVATE DEFAULT EPS BEARER CONTEXT
// REQUEST message, or 0 for a nil receiver. The result includes the 3-octet
// ESM header plus positional IEs and all named tagged IEs in spec-table order.
func (m *ActivateDefaultEPSBearerContextRequest) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EsmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	if m.EPSQoS != nil {
		n += m.EPSQoS.MarshalLenWithoutTag()
	}
	if m.AccessPointName != nil {
		n += m.AccessPointName.MarshalLenWithoutTag()
	}
	if m.PDNAddress != nil {
		n += m.PDNAddress.MarshalLenWithoutTag()
	}

	if m.LLCSAPI != nil {
		n += m.LLCSAPI.MarshalLen()
	}
	if m.ESMCause != nil {
		n += m.ESMCause.MarshalLen()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the ACTIVATE DEFAULT EPS BEARER CONTEXT REQUEST message
// into b. The nil receiver returns a wrapped ie.ErrNilIE; io.ErrShortBuffer
// without writing if b is too short. Header, positional IEs in spec-table
// order, then named tagged IEs follow.
func (m *ActivateDefaultEPSBearerContextRequest) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: ActivateDefaultEPSBearerContextRequest is nil", ie.ErrNilIE)
	}

	hlen := m.EsmHeader.MarshalLen()
	if hlen == 0 {
		return fmt.Errorf("%w: EsmHeader MarshalLen returned 0", ie.ErrNilIE)
	}

	if len(b) < hlen {
		return io.ErrShortBuffer
	}

	if err := m.EsmHeader.MarshalTo(b); err != nil {
		return fmt.Errorf("header: %w", err)
	}

	b = b[hlen:]

	if m.EPSQoS == nil {
		return io.ErrUnexpectedEOF
	}
	n := 0
	if err := m.EPSQoS.MarshalToWithoutTag(b[n:]); err != nil {
		return fmt.Errorf("EPSQoS: %w", err)
	}
	n += m.EPSQoS.MarshalLenWithoutTag()

	if m.AccessPointName == nil {
		return io.ErrUnexpectedEOF
	}
	if err := m.AccessPointName.MarshalToWithoutTag(b[n:]); err != nil {
		return fmt.Errorf("AccessPointName: %w", err)
	}
	n += m.AccessPointName.MarshalLenWithoutTag()

	if m.PDNAddress == nil {
		return io.ErrUnexpectedEOF
	}
	if err := m.PDNAddress.MarshalToWithoutTag(b[n:]); err != nil {
		return fmt.Errorf("PDNAddress: %w", err)
	}
	n += m.PDNAddress.MarshalLenWithoutTag()

	// Named tagged IEs in spec-table order.
	if m.LLCSAPI != nil {
		w, err := writeIE(b, n, m.LLCSAPI)
		if err != nil {
			return fmt.Errorf("LLCSAPI: %w", err)
		}
		n += w
	}

	if m.ESMCause != nil {
		w, err := writeIE(b, n, m.ESMCause)
		if err != nil {
			return fmt.Errorf("ESMCause: %w", err)
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

// UnmarshalBinary decodes an ACTIVATE DEFAULT EPS BEARER CONTEXT REQUEST
// message from b. The receiver is left unchanged on error (uses local tmp
// struct). Header, positional IEs via ie.ParseWithoutTag with prototypes,
// tagged optional IEs resolved via ie.ParseMultiIEs.
func (m *ActivateDefaultEPSBearerContextRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: ActivateDefaultEPSBearerContextRequest is nil", ie.ErrNilIE)
	}

	var tmp ActivateDefaultEPSBearerContextRequest

	// Decode header (3 octets for plain ESM).
	pay, err := tmp.EsmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	pos := 0

	tmp.EPSQoS, err = ie.ParseWithoutTag(pay[pos:],
		ie.NewPositionalEPSQualityOfService([]byte{0, 0}))
	if err != nil {
		return fmt.Errorf("EPSQoS: %w", err)
	}
	pos += tmp.EPSQoS.MarshalLenWithoutTag()

	tmp.AccessPointName, err = ie.ParseWithoutTag(pay[pos:],
		ie.NewPositionalAccessPointName([]byte{0, 0}))
	if err != nil {
		return fmt.Errorf("AccessPointName: %w", err)
	}
	pos += tmp.AccessPointName.MarshalLenWithoutTag()

	tmp.PDNAddress, err = ie.ParseWithoutTag(pay[pos:],
		ie.NewPositionalPDNAddress([]byte{0, 0}))
	if err != nil {
		return fmt.Errorf("PDNAddress: %w", err)
	}
	pos += tmp.PDNAddress.MarshalLenWithoutTag()

	// Tagged optional IEs — prototypes resolve ambiguous tags.
	tagged, err := ie.ParseMultiIEs(pay[pos:])
	if err != nil {
		return fmt.Errorf("optional: %w", err)
	}

	for _, i := range tagged {
		switch i.Kind {
		case ie.LLCServiceAccessPointIdentifier:
			if tmp.LLCSAPI == nil {
				tmp.LLCSAPI = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.ESMCause:
			if tmp.ESMCause == nil {
				tmp.ESMCause = i
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

// MessageType returns the ACTIVATE DEFAULT EPS BEARER CONTEXT REQUEST message
// type constant (0xC1).
func (m *ActivateDefaultEPSBearerContextRequest) MessageType() uint8 {
	return ESMActivateDefaultEPSBearerContextRequest
}

// MessageTypeName returns the descriptive name from Table 9.8.2.
func (m *ActivateDefaultEPSBearerContextRequest) MessageTypeName() string {
	return "Activate default EPS bearer context request"
}

// IsRequest returns false — ACTIVATE DEFAULT EPS BEARER CONTEXT REQUEST is
// sent network to UE.
func (m *ActivateDefaultEPSBearerContextRequest) IsRequest() bool {
	return false
}

// ActivateDefaultEPSBearerContextRequest implements the Message interface.
var _ Message = (*ActivateDefaultEPSBearerContextRequest)(nil)

// ParseActivateDefaultEPSBearerContextRequest parses a byte slice as an ACTIVATE
// DEFAULT EPS BEARER CONTEXT REQUEST message. Convenience wrapper around
// UnmarshalBinary on a fresh receiver.
func ParseActivateDefaultEPSBearerContextRequest(b []byte) (*ActivateDefaultEPSBearerContextRequest, error) {
	m := &ActivateDefaultEPSBearerContextRequest{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}
