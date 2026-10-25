package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// RemoteUEReport represents a REMOTE UE REPORT message as defined in TS 24.301
// clause 8.3.23 (Table 8.3.23.1, direction network to UE).
//
// The message carries only an ESM header followed by tagged optional IEs. Relay
// selection info (LV, tag 0x65) is resolved via prototype in ie.ParseMultiIEs;
// duplicates and unrecognized IEs are appended to the IEs field.
type RemoteUEReport struct {
	EsmHeader
	RelaySelectionInfo *ie.IE // tagged optional LV (tag 0x65, clause 9.9.4.13A), pass-through
	IEs                []*ie.IE
}

// NewRemoteUEReport creates a REMOTE UE REPORT message with optional tagged IEs.
// First RelaySelectionInfo (tag 0x65) occurrence wins; duplicates and unrecognized
// IEs are appended to the IEs field.
func NewRemoteUEReport(ies ...*ie.IE) *RemoteUEReport {
	m := &RemoteUEReport{}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
		case ie.Kind(0x65): // relay selection info
			if m.RelaySelectionInfo == nil {
				m.RelaySelectionInfo = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the REMOTE UE REPORT message, or 0 for a
// nil receiver. The result includes the 3-octet ESM header plus all named tagged
// IEs in spec-table order.
func (m *RemoteUEReport) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EsmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	if m.RelaySelectionInfo != nil {
		n += m.RelaySelectionInfo.MarshalLen()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the REMOTE UE REPORT message into b. The nil receiver
// returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too
// short. Writes header, then named tagged IEs in spec-table order.
func (m *RemoteUEReport) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: RemoteUEReport is nil", ie.ErrNilIE)
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

	// Named tagged IEs in spec-table order.
	if m.RelaySelectionInfo != nil {
		w, err := writeIE(b, n, m.RelaySelectionInfo)
		if err != nil {
			return fmt.Errorf("relay selection info: %w", err)
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

// UnmarshalBinary decodes a REMOTE UE REPORT message from b. The receiver is
// left unchanged on error (uses local tmp struct). Header via EsmHeader
// UnmarshalBinary; tagged optional IEs resolved via ie.ParseMultiIEs.
func (m *RemoteUEReport) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: RemoteUEReport is nil", ie.ErrNilIE)
	}

	var tmp RemoteUEReport

	// Decode header (3 octets for plain ESM).
	pay, err := tmp.EsmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	// Tagged optional IEs — prototypes resolve ambiguous tag 0x65 to relay
	// selection info. RelaySelectionInfo definition uses LV format (tag 0x65,
	// clause 9.9.4.13A).
	tagged, err := ie.ParseMultiIEs(pay)
	if err != nil {
		return fmt.Errorf("optional: %w", err)
	}

	for _, i := range tagged {
		switch i.Kind {
		case ie.Kind(0x65): // relay selection info
			if tmp.RelaySelectionInfo == nil {
				tmp.RelaySelectionInfo = i
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

// MessageType returns the REMOTE UE REPORT message type constant (0xE9).
func (m *RemoteUEReport) MessageType() uint8 {
	return ESMRemoteUEReport
}

// MessageTypeName returns the descriptive name from Table 9.8.2.
func (m *RemoteUEReport) MessageTypeName() string {
	return "Remote UE report"
}

// IsRequest returns false — REMOTE UE REPORT is sent network to UE.
func (m *RemoteUEReport) IsRequest() bool {
	return false
}

// RemoteUEReport implements the Message interface.
var _ Message = (*RemoteUEReport)(nil)

// ParseRemoteUEReport parses a byte slice as a REMOTE UE REPORT message.
// Convenience wrapper around UnmarshalBinary on a fresh receiver.
func ParseRemoteUEReport(b []byte) (*RemoteUEReport, error) {
	m := &RemoteUEReport{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}
