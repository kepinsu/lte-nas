package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// PDNConnectivityRequest represents a PDN CONNECTIVITY REQUEST message as defined
// in TS 24.301 clause 8.3.20 (Table 8.3.20.1, direction UE to network).
//
// The positional mandatory part encodes the PDN type and request type as a packed
// half-octet, followed by EPS QoS (LV, positional tag 0x5B) and access point name
// without IEI tags. Tagged optional IEs (ESM cause, tag 0x58) are resolved via
// prototypes in ie.ParseMultiIEs.
type PDNConnectivityRequest struct {
	EsmHeader
	PDNType         *ie.IE // positional half-octet, low nibble (9.9.4.10)
	RequestType     *ie.IE // positional half-octet, high nibble (9.9.4.14)
	EPSQoS          *ie.IE // positional LV-E (9.9.4.3), tag 0x5B
	AccessPointName *ie.IE // positional LV-E (9.9.4.1)
	ESMCause        *ie.IE // tagged optional type-3 TV (tag 0x58, 9.9.4.4)
	IEs             []*ie.IE
}

// NewPDNConnectivityRequest creates a PDN CONNECTIVITY REQUEST message with the
// mandatory positional IEs from TS 24.301 clause 8.3.20, followed by optional
// tagged IEs classified by Kind. ESMCause (Kind 0x58) is resolved first;
// duplicates and unrecognized IEs are appended to the IEs field.
func NewPDNConnectivityRequest(pdnType, requestType, epsQoS, accessPointName *ie.IE,
	ies ...*ie.IE,
) *PDNConnectivityRequest {
	if pdnType == nil || requestType == nil || epsQoS == nil || accessPointName == nil {
		panic("PDNConnectivityRequest positional IEs must not be nil")
	}

	m := &PDNConnectivityRequest{
		PDNType:         pdnType,
		RequestType:     requestType,
		EPSQoS:          epsQoS,
		AccessPointName: accessPointName,
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
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

// MarshalLen returns the wire size of the PDN CONNECTIVITY REQUEST message, or 0
// for a nil receiver. The result includes the 3-octet ESM header, packed half-
// octets (1 byte), positional IEs and all named tagged IEs in spec-table order.
func (m *PDNConnectivityRequest) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EsmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	// Packed half-octets: PDN type (low nibble) + request type (high nibble).
	if m.PDNType != nil && m.RequestType != nil {
		n++
	}

	if m.EPSQoS != nil {
		n += m.EPSQoS.MarshalLenWithoutTag()
	}
	if m.AccessPointName != nil {
		n += m.AccessPointName.MarshalLenWithoutTag()
	}

	if m.ESMCause != nil {
		n += m.ESMCause.MarshalLen()
	}
	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the PDN CONNECTIVITY REQUEST message into b. The nil receiver
// returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too short.
// Header, packed half-octets (1 byte), positional IEs in spec-table order, then
// named tagged IEs follow TS 24.301 clause 8.3.20 table order.
func (m *PDNConnectivityRequest) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: PDNConnectivityRequest is nil", ie.ErrNilIE)
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
	n := 0

	// Packed half-octets: PDN type (low nibble) + request type (high nibble).
	if m.PDNType == nil || m.RequestType == nil {
		return io.ErrUnexpectedEOF
	}

	hb, err := ie.MarshalHalfOctets(m.PDNType, m.RequestType)
	if err != nil {
		return fmt.Errorf("half-octets: %w", err)
	}

	if len(b)-n < len(hb) {
		return io.ErrShortBuffer
	}
	copy(b[n:], hb)
	n += len(hb)

	// Positional EPS QoS (LV, no tag).
	if m.EPSQoS == nil {
		return io.ErrUnexpectedEOF
	}
	if err := m.EPSQoS.MarshalToWithoutTag(b[n:]); err != nil {
		return fmt.Errorf("EPSQoS: %w", err)
	}
	n += m.EPSQoS.MarshalLenWithoutTag()

	// Positional Access point name (LV, no tag).
	if m.AccessPointName == nil {
		return io.ErrUnexpectedEOF
	}
	if err := m.AccessPointName.MarshalToWithoutTag(b[n:]); err != nil {
		return fmt.Errorf("AccessPointName: %w", err)
	}
	n += m.AccessPointName.MarshalLenWithoutTag()

	// Named tagged IEs in spec-table order.
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

// UnmarshalBinary decodes a PDN CONNECTIVITY REQUEST message from b. The receiver
// is left unchanged on error (uses local tmp struct). Header, packed half-octets
// via ParseHalfOctets, positional IEs via ie.ParseWithoutTag with prototypes, tagged
// optional IEs resolved via ie.ParseMultiIEs.
func (m *PDNConnectivityRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: PDNConnectivityRequest is nil", ie.ErrNilIE)
	}

	var tmp PDNConnectivityRequest

	// Decode header (3 octets for plain ESM).
	pay, err := tmp.EsmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	pos := 0

	// Packed half-octets: PDN type (low nibble) + request type (high nibble).
	if len(pay[pos:]) == 0 {
		return io.ErrUnexpectedEOF
	}

	pdnProto := ie.NewPositionalPDNType(0)
	reqProto := ie.NewPositionalRequestType(0)
	tmp.PDNType, tmp.RequestType, err = ie.ParseHalfOctets(pay[pos:], pdnProto, reqProto)
	if err != nil {
		return fmt.Errorf("half-octets: %w", err)
	}
	pos++ // consumed 1 byte for packed half-octets

	// Positional EPS QoS (LV, no tag).
	tmp.EPSQoS, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalEPSQualityOfService([]byte{0, 0}))
	if err != nil {
		return fmt.Errorf("EPSQoS: %w", err)
	}
	pos += tmp.EPSQoS.MarshalLenWithoutTag()

	// Positional Access point name (LV, no tag).
	tmp.AccessPointName, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalAccessPointName([]byte{0, 0}))
	if err != nil {
		return fmt.Errorf("AccessPointName: %w", err)
	}
	pos += tmp.AccessPointName.MarshalLenWithoutTag()

	// Tagged optional IEs — prototype resolves ambiguous tag 0x58 to ESM cause.
	if len(pay[pos:]) > 0 {
		tagged, err := ie.ParseMultiIEs(pay[pos:])
		if err != nil {
			return fmt.Errorf("optional: %w", err)
		}

		for _, i := range tagged {
			switch i.Kind {
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
	}

	*m = tmp
	return nil
}

// MessageType returns the PDN CONNECTIVITY REQUEST message type constant (0xD0).
func (m *PDNConnectivityRequest) MessageType() uint8 {
	return ESMPDNConnectivityRequest
}

// MessageTypeName returns the descriptive name from Table 9.8.2.
func (m *PDNConnectivityRequest) MessageTypeName() string {
	return "PDN connectivity request"
}

// IsRequest returns true — PDN CONNECTIVITY REQUEST is sent UE to network.
func (m *PDNConnectivityRequest) IsRequest() bool {
	return true
}

// PDNConnectivityRequest implements the Message interface.
var _ Message = (*PDNConnectivityRequest)(nil)

// ParsePDNConnectivityRequest parses a byte slice as a PDN CONNECTIVITY REQUEST
// message. Convenience wrapper around UnmarshalBinary on a fresh receiver.
func ParsePDNConnectivityRequest(b []byte) (*PDNConnectivityRequest, error) {
	m := &PDNConnectivityRequest{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}
