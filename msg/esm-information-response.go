package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// EsmInfoResponse represents an ESM INFORMATION RESPONSE message as defined in
// TS 24.301 clause 8.3.14 (Table 8.3.14.1, direction UE to network).
//
// The body carries only tagged optional IEs (Extended Protocol Configuration Options
// tag 0x7B, Network Feature Uplink tag 0xC4) — no positional mandatory part.
type EsmInfoResponse struct {
	EsmHeader
	ExtendedProtocolOpts *ie.IE // tagged optional LV-E (0x7B, 9.9.4.26)
	NetworkFeatureUp     *ie.IE // tagged optional LV-E (0xC4, 9.9.4.51)
	IEs                  []*ie.IE
}

// NewEsmInfoResponse creates an ESM INFORMATION RESPONSE message. No
// positional mandatory IEs exist for this message type; all body IEs are optional
// tagged and classified by Kind.
func NewEsmInfoResponse(ies ...*ie.IE) *EsmInfoResponse {
	m := &EsmInfoResponse{}

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
		case ie.Kind(0xC4):
			if m.NetworkFeatureUp == nil {
				m.NetworkFeatureUp = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the ESM INFORMATION RESPONSE message, or 0 for
// a nil receiver. The result includes header (3 octets) and all named tagged IEs
// in spec-table order.
func (m *EsmInfoResponse) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EsmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	if m.ExtendedProtocolOpts != nil {
		n += m.ExtendedProtocolOpts.MarshalLen()
	}
	if m.NetworkFeatureUp != nil {
		n += m.NetworkFeatureUp.MarshalLen()
	}
	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the ESM INFORMATION RESPONSE message into b. The nil receiver
// returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too short.
func (m *EsmInfoResponse) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: EsmInfoResponse is nil", ie.ErrNilIE)
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

	n := hlen

	// Named tagged IEs in spec-table order.
	if m.ExtendedProtocolOpts != nil {
		w, err := writeIE(b, n, m.ExtendedProtocolOpts)
		if err != nil {
			return fmt.Errorf("ExtendedProtocolOpts: %w", err)
		}
		n += w
	}

	if m.NetworkFeatureUp != nil {
		w, err := writeIE(b, n, m.NetworkFeatureUp)
		if err != nil {
			return fmt.Errorf("NetworkFeatureUp: %w", err)
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

// UnmarshalBinary decodes an ESM INFORMATION RESPONSE message from b. The receiver is
// left unchanged on error (uses local tmp struct). Header only, then tagged optional
// IEs resolved via ie.ParseMultiIEs.
func (m *EsmInfoResponse) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: EsmInfoResponse is nil", ie.ErrNilIE)
	}

	var tmp EsmInfoResponse

	// Decode header (3 octets for plain ESM).
	pay, err := tmp.EsmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	// Tagged optional IEs — prototypes resolve ambiguous tags.
	if len(pay) > 0 {
		tagged, err := ie.ParseMultiIEs(pay)
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
			case ie.Kind(0xC4):
				if tmp.NetworkFeatureUp == nil {
					tmp.NetworkFeatureUp = i
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

// MessageType returns the ESM INFORMATION RESPONSE message type constant (0xDA).
func (m *EsmInfoResponse) MessageType() uint8 {
	return ESMInformationResponse
}

// MessageTypeName returns the descriptive name from Table 9.8.2.
func (m *EsmInfoResponse) MessageTypeName() string {
	return "ESM information response"
}

// IsRequest returns true — ESM INFORMATION RESPONSE is sent UE to network.
func (m *EsmInfoResponse) IsRequest() bool {
	return true
}

// EsmInfoResponse implements the Message interface.
var _ Message = (*EsmInfoResponse)(nil)

// ParseEsmInfoResponse parses a byte slice as an ESM INFORMATION RESPONSE
// message. Convenience wrapper around UnmarshalBinary on a fresh receiver.
func ParseEsmInfoResponse(b []byte) (*EsmInfoResponse, error) {
	m := &EsmInfoResponse{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}
