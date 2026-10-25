package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// CSServiceNotification represents a CS SERVICE NOTIFICATION message as defined in
// TS 24.301 clause 8.2.9 (Table 8.2.9.1, direction network to UE).
//
// The positional mandatory part encodes paging identity (type-3 V) per clause
// 9.9.3.25A. Tagged optional IEs include CLI (0x60), SS code (0x61), LCS
// indicator (0x62) and LCS client identity (0x63). The message is
// integrity-protected (SHT 0x01..0x05, 7-octet header).
type CSServiceNotification struct {
	EmmHeader
	PagingIdentity    *ie.IE // positional type-3 (9.9.3.25A), paging identity
	Cli               *ie.IE // tagged optional TLV type-4 (9.9.3.38), tag 0x60
	SsCode            *ie.IE // tagged optional TV type-3 (9.9.3.39), tag 0x61
	LcsIndicator      *ie.IE // tagged optional TV type-3 (9.9.3.40), tag 0x62
	LcsClientIdentity *ie.IE // tagged optional TLV type-4 (9.9.3.41), tag 0x63
	IEs               []*ie.IE
}

// NewCSServiceNotification creates a CS SERVICE NOTIFICATION message with the
// mandatory positional Paging identity IE, followed by optional tagged IEs
// classified by their Kind. First occurrence wins; duplicates and unrecognized
// IEs are appended to the IEs field.
func NewCSServiceNotification(pagingIdentity *ie.IE, ies ...*ie.IE) *CSServiceNotification {
	if pagingIdentity == nil {
		panic("CSServiceNotification paging identity must not be nil")
	}

	m := &CSServiceNotification{PagingIdentity: pagingIdentity}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
		case ie.CLI:
			if m.Cli == nil {
				m.Cli = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.SSCode:
			if m.SsCode == nil {
				m.SsCode = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.LCSIndicator:
			if m.LcsIndicator == nil {
				m.LcsIndicator = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.LCSClientIdentity:
			if m.LcsClientIdentity == nil {
				m.LcsClientIdentity = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the CS SERVICE NOTIFICATION message, or 0
// for a nil receiver. The result includes header (7 octets), paging identity (1
// byte) and all named tagged IEs in spec-table order.
func (m *CSServiceNotification) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	if m.PagingIdentity != nil {
		n += m.PagingIdentity.MarshalLenWithoutTag()
	}

	if m.Cli != nil {
		n += m.Cli.MarshalLen()
	}
	if m.SsCode != nil {
		n += m.SsCode.MarshalLen()
	}
	if m.LcsIndicator != nil {
		n += m.LcsIndicator.MarshalLen()
	}
	if m.LcsClientIdentity != nil {
		n += m.LcsClientIdentity.MarshalLen()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the CS SERVICE NOTIFICATION message into b. The nil receiver
// returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too
// short. Header, paging identity (positional) and named tagged IEs follow TS
// 24.301 clause 8.2.9 table order.
func (m *CSServiceNotification) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: CSServiceNotification is nil", ie.ErrNilIE)
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

	if m.PagingIdentity != nil {
		if err := m.PagingIdentity.MarshalToWithoutTag(b); err != nil {
			return fmt.Errorf("PagingIdentity: %w", err)
		}
		b = b[1:]
	}

	if m.Cli != nil {
		w, err := writeIE(b, 0, m.Cli)
		if err != nil {
			return fmt.Errorf("cli: %w", err)
		}
		b = b[w:]
	}
	if m.SsCode != nil {
		w, err := writeIE(b, 0, m.SsCode)
		if err != nil {
			return fmt.Errorf("SsCode: %w", err)
		}
		b = b[w:]
	}
	if m.LcsIndicator != nil {
		w, err := writeIE(b, 0, m.LcsIndicator)
		if err != nil {
			return fmt.Errorf("LcsIndicator: %w", err)
		}
		b = b[w:]
	}
	if m.LcsClientIdentity != nil {
		w, err := writeIE(b, 0, m.LcsClientIdentity)
		if err != nil {
			return fmt.Errorf("LcsClientIdentity: %w", err)
		}
		b = b[w:]
	}

	for _, i := range m.IEs {
		w, err := writeIE(b, 0, i)
		if err != nil {
			return fmt.Errorf("extra IE: %w", err)
		}
		b = b[w:]
	}

	return nil
}

// UnmarshalBinary decodes a CS SERVICE NOTIFICATION message from b. The receiver
// is left unchanged on error. Header, paging identity (positional) via
// ParseWithoutTag, tagged optional IEs resolved via expected prototypes.
func (m *CSServiceNotification) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: CSServiceNotification is nil", ie.ErrNilIE)
	}

	var tmp CSServiceNotification

	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	pos := 0

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	// Positional mandatory: Paging identity (9.9.3.25A).
	tmp.PagingIdentity, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalPagingIdentity(0))
	if err != nil {
		return fmt.Errorf("PagingIdentity: %w", err)
	}
	pos += tmp.PagingIdentity.MarshalLenWithoutTag()

	// Tagged optional IEs — prototypes resolve ambiguous tags.
	tagged, err := ie.ParseMultiIEs(pay[pos:])
	if err != nil {
		return fmt.Errorf("optional: %w", err)
	}

	for _, i := range tagged {
		switch i.Kind {
		case ie.CLI:
			if tmp.Cli == nil {
				tmp.Cli = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.SSCode:
			if tmp.SsCode == nil {
				tmp.SsCode = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.LCSIndicator:
			if tmp.LcsIndicator == nil {
				tmp.LcsIndicator = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.LCSClientIdentity:
			if tmp.LcsClientIdentity == nil {
				tmp.LcsClientIdentity = i
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

// MessageType returns the CS SERVICE NOTIFICATION message type constant (0x64).
func (m *CSServiceNotification) MessageType() uint8 {
	return EMMCSServiceNotification
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *CSServiceNotification) MessageTypeName() string {
	return "CS service notification"
}

// IsRequest returns false — CS SERVICE NOTIFICATION is sent network to UE.
func (m *CSServiceNotification) IsRequest() bool {
	return false
}

// CSServiceNotification implements the Message interface.
var _ Message = (*CSServiceNotification)(nil)

// Header returns the EMM security header embedded in CSServiceNotification.
func (m *CSServiceNotification) Header() *EmmHeader {
	return &m.EmmHeader
}
