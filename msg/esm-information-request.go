package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// EsmInfoRequest represents an ESM INFORMATION REQUEST message as defined in
// TS 24.301 clause 8.3.13 (Table 8.3.13.1, dual-direction; direction network to UE
// for response purposes in this implementation).
//
// The positional mandatory part encodes UE Network Capability (LV, type-4) followed by
// the ESM message container. Tagged optional IEs (Extended Protocol Configuration Options
// tag 0x7B, Network Feature Downlink tag 0xC3) are resolved via prototypes.
type EsmInfoRequest struct {
	EsmHeader
	UENetworkCapability  *ie.IE // positional LV type-4 (9.9.3.34)
	ESMMessageContainer  *ie.IE // positional LV-E type-6 (9.9.3.15), pass-through
	ExtendedProtocolOpts *ie.IE // tagged optional LV-E (0x7B, 9.9.4.26)
	NetworkFeatureDown   *ie.IE // tagged optional LV-E (0xC3, 9.9.4.50A)
	IEs                  []*ie.IE
}

// NewEsmInfoRequest creates an ESM INFORMATION REQUEST with mandatory
// positional IEs (UE Network Capability, ESM Message Container) from TS 24.301
// clause 8.3.13, followed by optional tagged IEs classified by Kind.
func NewEsmInfoRequest(ueNetworkCapability, esmMessageContainer *ie.IE,
	ies ...*ie.IE,
) *EsmInfoRequest {
	if ueNetworkCapability == nil || esmMessageContainer == nil {
		panic("EsmInfoRequest positional IEs must not be nil")
	}

	m := &EsmInfoRequest{
		UENetworkCapability: ueNetworkCapability,
		ESMMessageContainer: esmMessageContainer,
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
		case ie.ExtendedProtocolConfigurationOptions:
			if m.ExtendedProtocolOpts == nil {
				m.ExtendedProtocolOpts = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.Kind(0xC3):
			if m.NetworkFeatureDown == nil {
				m.NetworkFeatureDown = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the ESM INFORMATION REQUEST message, or 0 for
// a nil receiver. The result includes header (3 octets), positional LV IEs, and all
// named tagged IEs in spec-table order.
func (m *EsmInfoRequest) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EsmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	if m.UENetworkCapability != nil {
		n += m.UENetworkCapability.MarshalLenWithoutTag()
	}
	if m.ESMMessageContainer != nil {
		n += m.ESMMessageContainer.MarshalLenWithoutTag()
	}
	if m.ExtendedProtocolOpts != nil {
		n += m.ExtendedProtocolOpts.MarshalLen()
	}
	if m.NetworkFeatureDown != nil {
		n += m.NetworkFeatureDown.MarshalLen()
	}
	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the ESM INFORMATION REQUEST message into b. The nil receiver
// returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too short.
func (m *EsmInfoRequest) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: EsmInfoRequest is nil", ie.ErrNilIE)
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

	// Positional UE Network Capability (LV, no tag).
	if m.UENetworkCapability == nil {
		return io.ErrUnexpectedEOF
	}
	if err := m.UENetworkCapability.MarshalToWithoutTag(b[n:]); err != nil {
		return fmt.Errorf("UENetworkCapability: %w", err)
	}
	n += m.UENetworkCapability.MarshalLenWithoutTag()

	// Positional ESM Message Container (LV-E, no tag).
	if m.ESMMessageContainer == nil {
		return io.ErrUnexpectedEOF
	}
	if err := m.ESMMessageContainer.MarshalToWithoutTag(b[n:]); err != nil {
		return fmt.Errorf("ESMMessageContainer: %w", err)
	}
	n += m.ESMMessageContainer.MarshalLenWithoutTag()

	// Named tagged IEs in spec-table order.
	if m.ExtendedProtocolOpts != nil {
		w, err := writeIE(b, n, m.ExtendedProtocolOpts)
		if err != nil {
			return fmt.Errorf("ExtendedProtocolOpts: %w", err)
		}
		n += w
	}

	if m.NetworkFeatureDown != nil {
		w, err := writeIE(b, n, m.NetworkFeatureDown)
		if err != nil {
			return fmt.Errorf("NetworkFeatureDown: %w", err)
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

// UnmarshalBinary decodes an ESM INFORMATION REQUEST message from b. The receiver is
// left unchanged on error (uses local tmp struct). Header, positional IEs via
// tagged optional IEs resolved via ie.ParseMultiIEs.
func (m *EsmInfoRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: EsmInfoRequest is nil", ie.ErrNilIE)
	}

	var tmp EsmInfoRequest

	// Decode header (3 octets for plain ESM).
	pay, err := tmp.EsmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	pos := 0

	// Positional UE Network Capability (LV, no tag).
	tmp.UENetworkCapability, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalUENetworkCapability([]byte{0, 0}))
	if err != nil {
		return fmt.Errorf("UENetworkCapability: %w", err)
	}
	pos += tmp.UENetworkCapability.MarshalLenWithoutTag()

	// Positional ESM Message Container (LV-E, no tag).
	tmp.ESMMessageContainer, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalESMMessageContainer([]byte{0, 0, 0}))
	if err != nil {
		return fmt.Errorf("ESMMessageContainer: %w", err)
	}
	pos += tmp.ESMMessageContainer.MarshalLenWithoutTag()

	// Tagged optional IEs — prototypes resolve ambiguous tags.
	if len(pay[pos:]) > 0 {
		tagged, err := ie.ParseMultiIEs(pay[pos:])
		if err != nil {
			return fmt.Errorf("optional: %w", err)
		}

		for _, i := range tagged {
			switch i.Kind {
			case ie.ExtendedProtocolConfigurationOptions:
				if tmp.ExtendedProtocolOpts == nil {
					tmp.ExtendedProtocolOpts = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case ie.Kind(0xC3):
				if tmp.NetworkFeatureDown == nil {
					tmp.NetworkFeatureDown = i
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

// MessageType returns the ESM INFORMATION REQUEST message type constant (0xD9).
func (m *EsmInfoRequest) MessageType() uint8 {
	return ESMInformationRequest
}

// MessageTypeName returns the descriptive name from Table 9.8.2.
func (m *EsmInfoRequest) MessageTypeName() string {
	return "ESM information request"
}

// IsRequest returns false — ESM INFORMATION REQUEST is network to UE.
func (m *EsmInfoRequest) IsRequest() bool {
	return false
}

// EsmInfoRequest implements the Message interface.
var _ Message = (*EsmInfoRequest)(nil)

// ParseEsmInfoRequest parses a byte slice as an ESM INFORMATION REQUEST
// message. Convenience wrapper around UnmarshalBinary on a fresh receiver.
func ParseEsmInfoRequest(b []byte) (*EsmInfoRequest, error) {
	m := &EsmInfoRequest{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}
