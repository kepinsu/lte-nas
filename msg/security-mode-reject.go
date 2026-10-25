package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// SecurityModeReject represents a SECURITY MODE REJECT message as defined in
// TS 24.301 clause 8.2.22 (Table 8.2.22.1, direction UE to network).
//
// The only content IE is the positional mandatory EMM cause (1 octet, type-3).
// There are no tagged optional IEs in TS 24.301 Table 8.2.22.1. This message is
// integrity-protected with a security header type of 0x01..0x05, yielding a
// 7-octet header per TS 24.301 clause 9.3.1.
type SecurityModeReject struct {
	EmmHeader
	EMMCause *ie.IE // positional type-3 (9.9.3.9)
	IEs      []*ie.IE
}

// NewSecurityModeReject creates a SECURITY MODE REJECT message with the mandatory
// positional EMM cause IE from TS 24.301 clause 8.2.22, followed by any optional
// tagged IEs (there are none per Table 8.2.22.1; extras are stored in the IEs
// field for future extension or passthrough).
func NewSecurityModeReject(emmCause *ie.IE, ies ...*ie.IE) *SecurityModeReject {
	if emmCause == nil {
		panic("SecurityModeReject positional IE (EMM cause) must not be nil")
	}

	m := &SecurityModeReject{
		EMMCause: emmCause,
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		m.IEs = append(m.IEs, i)
	}

	return m
}

// MarshalLen returns the wire size of the SECURITY MODE REJECT message, or 0 for
// a nil receiver. The result includes header and all IEs.
func (m *SecurityModeReject) MarshalLen() int {
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

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the SECURITY MODE REJECT message into b. The nil receiver
// returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too
// short. Header, positional mandatory IEs (no tags) and extras follow TS 24.301
// clause 8.2.22 table order.
func (m *SecurityModeReject) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: SecurityModeReject is nil", ie.ErrNilIE)
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

	// Positional mandatory IE without IEI tag.
	if m.EMMCause != nil {
		if err := m.EMMCause.MarshalToWithoutTag(b[n:]); err != nil {
			return fmt.Errorf("EMM cause: %w", err)
		}
		n += m.EMMCause.MarshalLenWithoutTag()
	}

	// Extra tagged IEs (there are none per Table 8.2.22.1; stored for future
	// extension or passthrough).
	for _, i := range m.IEs {
		w, err := writeIE(b, n, i)
		if err != nil {
			return fmt.Errorf("extra IE: %w", err)
		}
		n += w
	}

	return nil
}

// UnmarshalBinary decodes a SECURITY MODE REJECT message from b. The receiver
// is left unchanged on error. Header and positional mandatory EMM cause via
// ParseWithoutTag. No tagged optional part per TS 24.301 Table 8.2.22.1.
func (m *SecurityModeReject) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: SecurityModeReject is nil", ie.ErrNilIE)
	}

	var tmp SecurityModeReject

	// Decode header (7 octets for integrity-protected SHT 0x01..0x05).
	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	// Positional mandatory: EMM cause (1 byte, type-3).
	tmp.EMMCause, err = ie.ParseWithoutTag(pay, ie.NewPositionalEMMCause(0))
	if err != nil {
		return fmt.Errorf("EMM cause: %w", err)
	}

	*m = tmp
	return nil
}

// MessageType returns the SECURITY MODE REJECT message type constant (0x5F).
func (m *SecurityModeReject) MessageType() uint8 {
	return EMMSecurityModeReject
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *SecurityModeReject) MessageTypeName() string {
	return "Security mode reject"
}

// IsRequest returns true — SECURITY MODE REJECT is sent UE to network.
func (m *SecurityModeReject) IsRequest() bool {
	return true
}

// SecurityModeReject implements the Message interface.
var _ Message = (*SecurityModeReject)(nil)

// Header returns the EMM security header embedded in SecurityModeReject.
func (m *SecurityModeReject) Header() *EmmHeader {
	return &m.EmmHeader
}
