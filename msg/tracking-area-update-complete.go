package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
)

// TrackingAreaUpdateComplete represents a TRACKING AREA UPDATE COMPLETE message as
// defined in TS 24.301 clause 8.2.27 (Table 8.2.27.1, direction UE to network).
//
// This message has no mandatory positional payload body — only the EMM header
// (plain SHT=0x00, 2 octets) and message identity. All content IEs (if any) are
// tagged optional and not listed in Table 8.2.27.1.
type TrackingAreaUpdateComplete struct {
	EmmHeader
	IEs []*ie.IE
}

// NewTrackingAreaUpdateComplete creates a TRACKING AREA UPDATE COMPLETE message
// with no mandatory positional IEs. Tagged IEs (if present) are classified by
// their Kind; duplicates and unrecognized IEs are appended to the IEs field.
func NewTrackingAreaUpdateComplete(ies ...*ie.IE) *TrackingAreaUpdateComplete {
	m := &TrackingAreaUpdateComplete{}

	for _, i := range ies {
		if i == nil {
			continue
		}
		m.IEs = append(m.IEs, i)
	}

	return m
}

// MarshalLen returns the wire size of the TRACKING AREA UPDATE COMPLETE message,
// or 0 for a nil receiver. The result includes header only — no positional body.
func (m *TrackingAreaUpdateComplete) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the TRACKING AREA UPDATE COMPLETE message into b. The nil
// receiver returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b
// is too short. Header only (no positional body); optional tagged IEs follow.
func (m *TrackingAreaUpdateComplete) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: TrackingAreaUpdateComplete is nil", ie.ErrNilIE)
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

	n := hlen

	// Optional tagged IEs.
	for _, i := range m.IEs {
		w, err := writeIE(b, n, i)
		if err != nil {
			return fmt.Errorf("extra IE: %w", err)
		}
		n += w
	}

	return nil
}

// UnmarshalBinary decodes a TRACKING AREA UPDATE COMPLETE message from b. The
// receiver is left unchanged on error. Header via UnmarshalBinary; no mandatory
// positional body (see TS 24.301 Table 8.2.27.1).
func (m *TrackingAreaUpdateComplete) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: TrackingAreaUpdateComplete is nil", ie.ErrNilIE)
	}

	var tmp TrackingAreaUpdateComplete

	// Decode header (2 octets for plain).
	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	// No mandatory positional body per Table 8.2.27.1. Any remaining bytes are
	// optional tagged IEs (though none are defined in the spec table).
	if len(pay) > 0 {
		tagged, err := ie.ParseMultiIEs(pay)
		if err != nil {
			return fmt.Errorf("optional: %w", err)
		}

		tmp.IEs = append(tmp.IEs, tagged...)
	}

	*m = tmp
	return nil
}

// MessageType returns the TRACKING AREA UPDATE COMPLETE message type constant (0x4A).
func (m *TrackingAreaUpdateComplete) MessageType() uint8 {
	return EMMTrackingAreaUpdateComplete
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *TrackingAreaUpdateComplete) MessageTypeName() string {
	return "Tracking area update complete"
}

// IsRequest returns true — TRACKING AREA UPDATE COMPLETE is sent UE to network.
func (m *TrackingAreaUpdateComplete) IsRequest() bool {
	return true
}

// TrackingAreaUpdateComplete implements the Message interface.
var _ Message = (*TrackingAreaUpdateComplete)(nil)

// Header returns the EMM security header embedded in TrackingAreaUpdateComplete.
func (m *TrackingAreaUpdateComplete) Header() *EmmHeader {
	return &m.EmmHeader
}
