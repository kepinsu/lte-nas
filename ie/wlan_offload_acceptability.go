package ie

import (
	"fmt"
	"io"
)

// WLANOffloadAcceptabilityFields is the structured value of the WLAN offload
// acceptability IE (TS 24.301 clause 9.9.4.18, see TS 24.008 clause 10.5.6.20).
type WLANOffloadAcceptabilityFields struct {
	// OffloadIndication indicates whether offloading to WLAN is acceptable.
	OffloadIndication bool
}

// NewWLANOffloadAcceptabilityWithFields creates a tagged WLAN offload
// acceptability IE from fields.
func NewWLANOffloadAcceptabilityWithFields(fields *WLANOffloadAcceptabilityFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewWLANOffloadAcceptability(payload[0]), nil
}

// WLANOffloadAcceptability returns the IE as structured fields.
func (i *IE) WLANOffloadAcceptability() (*WLANOffloadAcceptabilityFields, error) {
	if err := i.requireDefinition(WLANOffloadAcceptabilityDefinition); err != nil {
		return nil, err
	}
	return ParseWLANOffloadAcceptabilityFields(i.Payload)
}

// ParseWLANOffloadAcceptabilityFields decodes a WLAN offload acceptability value
// without its IEI.
func ParseWLANOffloadAcceptabilityFields(payload []byte) (*WLANOffloadAcceptabilityFields, error) {
	fields := &WLANOffloadAcceptabilityFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the receiver is non-nil.
func (f *WLANOffloadAcceptabilityFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: WLAN offload acceptability fields are nil", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a WLAN offload acceptability value, without its IEI.
func (f *WLANOffloadAcceptabilityFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	var encoded byte
	if f.OffloadIndication {
		encoded |= 0x80
	}
	return []byte{encoded}, nil
}

// MarshalTo encodes fields into b.
func (f *WLANOffloadAcceptabilityFields) MarshalTo(b []byte) error {
	payload, err := f.Marshal()
	if err != nil {
		return err
	}
	if len(b) < len(payload) {
		return io.ErrShortBuffer
	}
	b[0] = payload[0]
	return nil
}

// MarshalLen returns the size of the WLAN offload acceptability value payload.
func (f *WLANOffloadAcceptabilityFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a WLAN offload acceptability value without its IEI.
func (f *WLANOffloadAcceptabilityFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: WLAN offload acceptability fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: WLAN offload acceptability has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0x7f != 0 {
		return fmt.Errorf("%w: WLAN offload acceptability spare bits are nonzero", ErrInvalidValue)
	}
	*f = WLANOffloadAcceptabilityFields{
		OffloadIndication: payload[0]&0x80 != 0,
	}
	return nil
}
