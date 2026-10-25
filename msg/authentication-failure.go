package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// AuthenticationFailure represents an AUTHENTICATION FAILURE message as defined in TS 24.301
// clause 8.2.5 (Table 8.2.5.1, direction UE to network).
//
// The positional mandatory part encodes the EMM cause as a type-1 IE (1 octet). The tagged
// optional part contains the authentication failure parameter (IEI 0x30, TLV 16 bytes) which
// is present only when the EMM cause is #21 "synch failure" (see TS 24.301 clause 8.2.5.2).
type AuthenticationFailure struct {
	EmmHeader
	EMMCause                       *ie.IE   // positional type-1 (9.9.3.9), 1 octet
	AuthenticationFailureParameter *ie.IE   // tagged optional IEI 0x30 (9.9.3.1), TLV 16 bytes
	IEs                            []*ie.IE // trailing unrecognized IEs
}

// NewAuthenticationFailure creates an AUTHENTICATION FAILURE message with the mandatory
// positional EMM cause from TS 24.301 clause 8.2.5, followed by optional tagged IEs
// classified by their Kind. AuthenticationFailureParameter (Kind 0x30) is resolved first;
// duplicates and unrecognized IEs are appended to the IEs field.
func NewAuthenticationFailure(cause *ie.IE, ies ...*ie.IE) *AuthenticationFailure {
	if cause == nil {
		panic("AuthenticationFailure positional IEs must not be nil")
	}

	m := &AuthenticationFailure{
		EMMCause: cause,
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
		case ie.AuthenticationFailureParameter:
			if m.AuthenticationFailureParameter == nil {
				m.AuthenticationFailureParameter = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the AUTHENTICATION FAILURE message, or 0 for a nil
// receiver. The result includes header, positional EMM cause (1 byte), optional tagged
// authentication failure parameter and extras.
func (m *AuthenticationFailure) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	if m.EMMCause != nil {
		n += m.EMMCause.MarshalLenWithoutTag()
	}

	if m.AuthenticationFailureParameter != nil {
		n += m.AuthenticationFailureParameter.MarshalLen()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the AUTHENTICATION FAILURE message into b. The nil receiver returns
// a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too short. Header,
// positional EMM cause, tagged authentication failure parameter and extras follow TS 24.301
// clause 8.2.5 table order.
func (m *AuthenticationFailure) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: AuthenticationFailure is nil", ie.ErrNilIE)
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

	if m.EMMCause != nil {
		if err := m.EMMCause.MarshalToWithoutTag(b[n:]); err != nil {
			return fmt.Errorf("EMMCause: %w", err)
		}
		n += m.EMMCause.MarshalLenWithoutTag()
	}

	if m.AuthenticationFailureParameter != nil {
		w, err := writeIE(b, n, m.AuthenticationFailureParameter)
		if err != nil {
			return fmt.Errorf("AuthenticationFailureParameter: %w", err)
		}
		n += w
	}

	for _, i := range m.IEs {
		w, err := writeIE(b, n, i)
		if err != nil {
			return fmt.Errorf("extra IE: %w", err)
		}
		n += w
	}

	return nil
}

// UnmarshalBinary decodes an AUTHENTICATION FAILURE message from b. The receiver is left
// unchanged on error. Header, positional EMM cause via ie.ParseWithoutTag with prototype,
// resolved via ie.ParseMultiIEs.
func (m *AuthenticationFailure) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: AuthenticationFailure is nil", ie.ErrNilIE)
	}

	var tmp AuthenticationFailure

	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	pos := 0

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	// Positional mandatory EMM cause (type-1, 1 octet).
	tmp.EMMCause, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalEMMCause(0))
	if err != nil {
		return fmt.Errorf("EMMCause: %w", err)
	}
	pos += tmp.EMMCause.MarshalLenWithoutTag()

	// Tagged optional authentication failure parameter (IEI 0x30, TLV).
	// Note: IEI 0x30 is ambiguous in the catalogue (QualityOfService, UnavailabilityInformation
	// share it). We resolve it per-message context to AuthenticationFailureParameter.
	tagged, err := ie.ParseMultiIEs(pay[pos:])
	if err != nil {
		return fmt.Errorf("optional: %w", err)
	}

	for _, i := range tagged {
		switch i.Kind {
		case ie.AuthenticationFailureParameter:
			if tmp.AuthenticationFailureParameter == nil {
				tmp.AuthenticationFailureParameter = i
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

// MessageType returns the AUTHENTICATION FAILURE message type constant (0x5C).
func (m *AuthenticationFailure) MessageType() uint8 {
	return EMMAuthenticationFailure
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *AuthenticationFailure) MessageTypeName() string {
	return "Authentication failure"
}

// IsRequest returns true — AUTHENTICATION FAILURE is sent UE to network.
func (m *AuthenticationFailure) IsRequest() bool {
	return true
}

// AuthenticationFailure implements the Message interface.
var _ Message = (*AuthenticationFailure)(nil)

// Header returns the EMM security header embedded in AuthenticationFailure.
func (m *AuthenticationFailure) Header() *EmmHeader {
	return &m.EmmHeader
}
