package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// DeactivateEPSBearerContextRequest represents a DEACTIVATE EPS BEARER CONTEXT
// REQUEST message as defined in TS 24.301 clause 8.3.12 (Table 8.3.12.1,
// direction network to UE).
//
// The positional mandatory part encodes ESM cause (V, 1 octet). Tagged optional
// IEs: LLC SAPI (tag 0x32), Extended protocol configuration options (tag 0x7b).
type DeactivateEPSBearerContextRequest struct {
	EsmHeader
	ESMCause                   *ie.IE // positional V type-3 (9.9.4.4), one octet
	LLCSAPI                    *ie.IE // tagged optional type-3 TV, ambiguous tag 0x32 (9.9.4.7)
	ExtendedProtocolConfigOpts *ie.IE // tagged optional TLV-E type-6, tag 0x7b (9.9.4.26)
	IEs                        []*ie.IE
}

// NewDeactivateEPSBearerContextRequest creates a DEACTIVATE EPS BEARER CONTEXT
// REQUEST message with the mandatory positional ESM cause from TS 24.301 clause
// 8.3.12, followed by optional tagged IEs classified by Kind. LLCSAPI (Kind 0x32)
// and ExtendedProtocolConfigOpts (Kind 0x7b) are resolved first; duplicates and
// unrecognized IEs go to the IEs field.
func NewDeactivateEPSBearerContextRequest(esmCause *ie.IE, ies ...*ie.IE) *DeactivateEPSBearerContextRequest {
	if esmCause == nil {
		panic("DeactivateEPSBearerContextRequest positional ESM cause must not be nil")
	}

	m := &DeactivateEPSBearerContextRequest{
		ESMCause: esmCause,
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
		case ie.ExtendedProtocolConfigurationOptions:
			if m.ExtendedProtocolConfigOpts == nil {
				m.ExtendedProtocolConfigOpts = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the DEACTIVATE EPS BEARER CONTEXT REQUEST
// message, or 0 for a nil receiver. The result includes the 3-octet ESM header,
// positional ESM cause and all named tagged IEs.
func (m *DeactivateEPSBearerContextRequest) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EsmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	if m.ESMCause != nil {
		n += m.ESMCause.MarshalLenWithoutTag()
	}
	if m.LLCSAPI != nil {
		n += m.LLCSAPI.MarshalLen()
	}
	if m.ExtendedProtocolConfigOpts != nil {
		n += m.ExtendedProtocolConfigOpts.MarshalLen()
	}
	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the DEACTIVATE EPS BEARER CONTEXT REQUEST message into b. The
// nil receiver returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if
// b is too short. Header, positional ESM cause, named tagged IEs and extras follow
// TS 24.301 clause 8.3.12 table order.
func (m *DeactivateEPSBearerContextRequest) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: DeactivateEPSBearerContextRequest is nil", ie.ErrNilIE)
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

	if m.ESMCause == nil {
		return io.ErrUnexpectedEOF
	}
	if err := m.ESMCause.MarshalToWithoutTag(b[n:]); err != nil {
		return fmt.Errorf("ESMCause: %w", err)
	}
	n += m.ESMCause.MarshalLenWithoutTag()

	// Named tagged IEs in spec-table order.
	if m.LLCSAPI != nil {
		w, err := writeIE(b, n, m.LLCSAPI)
		if err != nil {
			return fmt.Errorf("LLCSAPI: %w", err)
		}
		n += w
	}
	if m.ExtendedProtocolConfigOpts != nil {
		w, err := writeIE(b, n, m.ExtendedProtocolConfigOpts)
		if err != nil {
			return fmt.Errorf("ExtendedProtocolConfigOpts: %w", err)
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

// UnmarshalBinary decodes a DEACTIVATE EPS BEARER CONTEXT REQUEST message from b.
// The receiver is left unchanged on error (uses local tmp struct). Header via
// UnmarshalBinary, positional ESM cause via ParseWithoutTag with prototypes, tagged
// IEs resolved via ie.ParseMultiIEs (0x32 = LLC SAPI, 0x7b = extended PCO).
func (m *DeactivateEPSBearerContextRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: DeactivateEPSBearerContextRequest is nil", ie.ErrNilIE)
	}

	var tmp DeactivateEPSBearerContextRequest

	// Decode header (3 octets for plain ESM).
	pay, err := tmp.EsmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	pos := 0

	// Positional ESM cause (V, 1 octet).
	if len(pay[pos:]) < 1 {
		return io.ErrUnexpectedEOF
	}

	tmp.ESMCause, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalESMCause(0))
	if err != nil {
		return fmt.Errorf("ESMCause: %w", err)
	}
	pos += tmp.ESMCause.MarshalLenWithoutTag()

	// Tagged optional IEs — prototypes resolve ambiguous tag 0x32 to LLC SAPI and
	// tag 0x7b to extended protocol configuration options.
	if len(pay[pos:]) > 0 {
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
			case ie.ExtendedProtocolConfigurationOptions:
				if tmp.ExtendedProtocolConfigOpts == nil {
					tmp.ExtendedProtocolConfigOpts = i
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

// MessageType returns the DEACTIVATE EPS BEARER CONTEXT REQUEST message type constant (0xCD).
func (m *DeactivateEPSBearerContextRequest) MessageType() uint8 {
	return ESMDeactivateEPSBearerContextRequest
}

// MessageTypeName returns the descriptive name from Table 9.8.2.
func (m *DeactivateEPSBearerContextRequest) MessageTypeName() string {
	return "Deactivate EPS bearer context request"
}

// IsRequest returns false — DEACTIVATE EPS BEARER CONTEXT REQUEST is sent network to UE.
func (m *DeactivateEPSBearerContextRequest) IsRequest() bool {
	return false
}

// DeactivateEPSBearerContextRequest implements the Message interface.
var _ Message = (*DeactivateEPSBearerContextRequest)(nil)

// ParseDeactivateEPSBearerContextRequest parses a byte slice as a DEACTIVATE EPS
// BEARER CONTEXT REQUEST message. Convenience wrapper around UnmarshalBinary on a
// fresh receiver.
func ParseDeactivateEPSBearerContextRequest(b []byte) (*DeactivateEPSBearerContextRequest, error) {
	m := &DeactivateEPSBearerContextRequest{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}
