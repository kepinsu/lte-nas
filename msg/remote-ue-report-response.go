package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// RemoteUEReportResponse represents a REMOTE UE REPORT RESPONSE message as
// defined in TS 24.301 clause 8.3.24 (Table 8.3.24.1, direction network to UE).
//
// The message carries an ESM header followed by tagged optional IEs: UE remote ID
// (LV-E, tag 0x7F) and relay selection info (LV, tag 0x65). Both are resolved via
// prototypes in ie.ParseMultiIEs; duplicates and unrecognized IEs are appended to
// the IEs field.
type RemoteUEReportResponse struct {
	EsmHeader
	UERemoteID         *ie.IE // tagged optional LV-E (tag 0x7F, clause 9.9.x), pass-through
	RelaySelectionInfo *ie.IE // tagged optional LV (tag 0x65, clause 9.9.4.13A), pass-through
	IEs                []*ie.IE
}

// NewRemoteUEReportResponse creates a REMOTE UE REPORT RESPONSE message with
// optional tagged IEs. First occurrence of UERemoteID (tag 0x7F) and
// RelaySelectionInfo (tag 0x65) wins; duplicates and unrecognized IEs are
// appended to the IEs field.
func NewRemoteUEReportResponse(ies ...*ie.IE) *RemoteUEReportResponse {
	m := &RemoteUEReportResponse{}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
		case ie.Kind(0x7F): // UE remote ID
			if m.UERemoteID == nil {
				m.UERemoteID = i
			} else {
				m.IEs = append(m.IEs, i)
			}
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

// MarshalLen returns the wire size of the REMOTE UE REPORT RESPONSE message, or 0
// for a nil receiver. The result includes the 3-octet ESM header plus all named
// tagged IEs in spec-table order.
func (m *RemoteUEReportResponse) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EsmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	if m.UERemoteID != nil {
		n += m.UERemoteID.MarshalLen()
	}

	if m.RelaySelectionInfo != nil {
		n += m.RelaySelectionInfo.MarshalLen()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the REMOTE UE REPORT RESPONSE message into b. The nil
// receiver returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b
// is too short. Writes header, then named tagged IEs in spec-table order.
func (m *RemoteUEReportResponse) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: RemoteUEReportResponse is nil", ie.ErrNilIE)
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
	if m.UERemoteID != nil {
		w, err := writeIE(b, n, m.UERemoteID)
		if err != nil {
			return fmt.Errorf("UE remote ID: %w", err)
		}
		n += w
	}

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

// UnmarshalBinary decodes a REMOTE UE REPORT RESPONSE message from b. The
// receiver is left unchanged on error (uses local tmp struct). Header via
// EsmHeader UnmarshalBinary; tagged optional IEs resolved via ie.ParseMultiIEs.
func (m *RemoteUEReportResponse) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: RemoteUEReportResponse is nil", ie.ErrNilIE)
	}

	var tmp RemoteUEReportResponse

	// Decode header (3 octets for plain ESM).
	pay, err := tmp.EsmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	// Tagged optional IEs — prototypes resolve ambiguous tag 0x7F to UE remote ID
	// (LV-E, clause 9.9.x) and tag 0x65 to relay selection info (LV, clause
	// 9.9.4.13A).
	tagged, err := ie.ParseMultiIEs(pay)
	if err != nil {
		return fmt.Errorf("optional: %w", err)
	}

	for _, i := range tagged {
		switch i.Kind {
		case ie.Kind(0x7F): // UE remote ID
			if tmp.UERemoteID == nil {
				tmp.UERemoteID = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
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

// MessageType returns the REMOTE UE REPORT RESPONSE message type constant (0xEA).
func (m *RemoteUEReportResponse) MessageType() uint8 {
	return ESMRemoteUEReportResponse
}

// MessageTypeName returns the descriptive name from Table 9.8.2.
func (m *RemoteUEReportResponse) MessageTypeName() string {
	return "Remote UE report response"
}

// IsRequest returns false — REMOTE UE REPORT RESPONSE is sent network to UE.
func (m *RemoteUEReportResponse) IsRequest() bool {
	return false
}

// RemoteUEReportResponse implements the Message interface.
var _ Message = (*RemoteUEReportResponse)(nil)

// ParseRemoteUEReportResponse parses a byte slice as a REMOTE UE REPORT RESPONSE
// message. Convenience wrapper around UnmarshalBinary on a fresh receiver.
func ParseRemoteUEReportResponse(b []byte) (*RemoteUEReportResponse, error) {
	m := &RemoteUEReportResponse{}
	if err := m.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return m, nil
}
