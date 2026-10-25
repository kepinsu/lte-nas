package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// ServiceReject represents a SERVICE REJECT message as defined in TS 24.301
// clause 8.2.24 (Table 8.2.24.1, direction network to UE).
//
// The positional mandatory part encodes EMM cause (type-1 V). Optional tagged IEs
// include T3402 value (TV 0x17), back off timer value (TLV 0x37), lower bound
// timer value (TV 0x1C), and tracking area identity lists for roaming / regional
// provision (TLV 0x1D / 0x1E). The message is integrity-protected (SHT 0x01..0x05,
// 7-octet header).
type ServiceReject struct {
	EmmHeader
	EMMCause                                   *ie.IE // positional type-1 (9.9.3.9), EMM cause
	T3402Value                                 *ie.IE // tagged optional TV type-1 (9.9.3.16), tag 0x17
	BackOffTimerValue                          *ie.IE // tagged optional TLV type-4 (9.9.3.16A), tag 0x37
	LowerBoundTimerValue                       *ie.IE // tagged optional TV type-1 (9.9.3.16B), tag 0x1C
	ForbiddenTAIsForRoaming                    *ie.IE // tagged optional TLV type-4 (9.9.3.33), tag 0x1D
	ForbiddenTAIsForRegionalProvisionOfService *ie.IE // tagged optional TLV type-4 (9.9.3.33), tag 0x1E
	IEs                                        []*ie.IE
}

// NewServiceReject creates a SERVICE REJECT message with the mandatory positional
// EMM cause IE, followed by optional tagged IEs classified by their Kind. First
// occurrence wins; duplicates and unrecognized IEs are appended to the IEs field.
func NewServiceReject(emmCause *ie.IE, ies ...*ie.IE) *ServiceReject {
	if emmCause == nil {
		panic("ServiceReject EMM cause must not be nil")
	}

	m := &ServiceReject{EMMCause: emmCause}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
		case ie.T3402Value:
			if m.T3402Value == nil {
				m.T3402Value = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.BackOffTimerValue:
			if m.BackOffTimerValue == nil {
				m.BackOffTimerValue = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.LowerBoundTimerValue:
			if m.LowerBoundTimerValue == nil {
				m.LowerBoundTimerValue = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.ForbiddenTAIsForRoaming:
			if m.ForbiddenTAIsForRoaming == nil {
				m.ForbiddenTAIsForRoaming = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.ForbiddenTAIsForRegionalProvisionOfService:
			if m.ForbiddenTAIsForRegionalProvisionOfService == nil {
				m.ForbiddenTAIsForRegionalProvisionOfService = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the SERVICE REJECT message, or 0 for a nil
// receiver. The result includes header (7 octets), EMM cause (1 byte) and all
// named tagged IEs in spec-table order.
func (m *ServiceReject) MarshalLen() int {
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

	if m.T3402Value != nil {
		n += m.T3402Value.MarshalLen()
	}
	if m.BackOffTimerValue != nil {
		n += m.BackOffTimerValue.MarshalLen()
	}
	if m.LowerBoundTimerValue != nil {
		n += m.LowerBoundTimerValue.MarshalLen()
	}
	if m.ForbiddenTAIsForRoaming != nil {
		n += m.ForbiddenTAIsForRoaming.MarshalLen()
	}
	if m.ForbiddenTAIsForRegionalProvisionOfService != nil {
		n += m.ForbiddenTAIsForRegionalProvisionOfService.MarshalLen()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the SERVICE REJECT message into b. The nil receiver returns
// a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too short.
// Header, EMM cause (positional) and named tagged IEs follow TS 24.301 clause
// 8.2.24 table order.
func (m *ServiceReject) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: ServiceReject is nil", ie.ErrNilIE)
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

	if m.EMMCause != nil {
		if err := m.EMMCause.MarshalToWithoutTag(b); err != nil {
			return fmt.Errorf("EMMCause: %w", err)
		}
		b = b[1:]
	}

	if m.T3402Value != nil {
		w, err := writeIE(b, 0, m.T3402Value)
		if err != nil {
			return fmt.Errorf("T3402Value: %w", err)
		}
		b = b[w:]
	}
	if m.BackOffTimerValue != nil {
		w, err := writeIE(b, 0, m.BackOffTimerValue)
		if err != nil {
			return fmt.Errorf("BackOffTimerValue: %w", err)
		}
		b = b[w:]
	}
	if m.LowerBoundTimerValue != nil {
		w, err := writeIE(b, 0, m.LowerBoundTimerValue)
		if err != nil {
			return fmt.Errorf("LowerBoundTimerValue: %w", err)
		}
		b = b[w:]
	}
	if m.ForbiddenTAIsForRoaming != nil {
		w, err := writeIE(b, 0, m.ForbiddenTAIsForRoaming)
		if err != nil {
			return fmt.Errorf("ForbiddenTAIsForRoaming: %w", err)
		}
		b = b[w:]
	}
	if m.ForbiddenTAIsForRegionalProvisionOfService != nil {
		w, err := writeIE(b, 0, m.ForbiddenTAIsForRegionalProvisionOfService)
		if err != nil {
			return fmt.Errorf("ForbiddenTAIsForRegionalProvisionOfService: %w", err)
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

// UnmarshalBinary decodes a SERVICE REJECT message from b. The receiver is left
// unchanged on error. Header, EMM cause (positional) via ParseWithoutTag, tagged
// optional IEs resolved via expected prototypes.
func (m *ServiceReject) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: ServiceReject is nil", ie.ErrNilIE)
	}

	var tmp ServiceReject

	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	pos := 0

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	// Positional mandatory: EMM cause.
	tmp.EMMCause, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalEMMCause(0))
	if err != nil {
		return fmt.Errorf("EMMCause: %w", err)
	}
	pos += tmp.EMMCause.MarshalLenWithoutTag()

	// Tagged optional IEs — prototypes resolve ambiguous tags.
	tagged, err := ie.ParseMultiIEs(pay[pos:])
	if err != nil {
		return fmt.Errorf("optional: %w", err)
	}

	for _, i := range tagged {
		switch i.Kind {
		case ie.T3402Value:
			if tmp.T3402Value == nil {
				tmp.T3402Value = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.BackOffTimerValue:
			if tmp.BackOffTimerValue == nil {
				tmp.BackOffTimerValue = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.LowerBoundTimerValue:
			if tmp.LowerBoundTimerValue == nil {
				tmp.LowerBoundTimerValue = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.ForbiddenTAIsForRoaming:
			if tmp.ForbiddenTAIsForRoaming == nil {
				tmp.ForbiddenTAIsForRoaming = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.ForbiddenTAIsForRegionalProvisionOfService:
			if tmp.ForbiddenTAIsForRegionalProvisionOfService == nil {
				tmp.ForbiddenTAIsForRegionalProvisionOfService = i
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

// MessageType returns the SERVICE REJECT message type constant (0x4E).
func (m *ServiceReject) MessageType() uint8 {
	return EMMServiceReject
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *ServiceReject) MessageTypeName() string {
	return "Service reject"
}

// IsRequest returns false — SERVICE REJECT is sent network to UE.
func (m *ServiceReject) IsRequest() bool {
	return false
}

// ServiceReject implements the Message interface.
var _ Message = (*ServiceReject)(nil)

// Header returns the EMM security header embedded in ServiceReject.
func (m *ServiceReject) Header() *EmmHeader {
	return &m.EmmHeader
}
