package ie

import (
	"fmt"
	"io"
	"slices"
)

// EPSBearerContextStatusFields is the structured value of the EPS bearer
// context status IE (TS 24.301 clause 9.9.2.1). ActiveBearerIDs contains the
// EPS bearer identities whose ESM state is not BEARER CONTEXT-INACTIVE.
type EPSBearerContextStatusFields struct {
	ActiveBearerIDs []uint8
}

// NewEPSBearerContextStatusWithFields creates a tagged EPS bearer context
// status IE from fields.
func NewEPSBearerContextStatusWithFields(fields *EPSBearerContextStatusFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewEPSBearerContextStatus(payload), nil
}

// NewPositionalEPSBearerContextStatusWithFields creates a positional EPS
// bearer context status IE from fields.
func NewPositionalEPSBearerContextStatusWithFields(fields *EPSBearerContextStatusFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalEPSBearerContextStatus(payload), nil
}

// EPSBearerContextStatus returns the IE as structured fields.
func (i *IE) EPSBearerContextStatus() (*EPSBearerContextStatusFields, error) {
	if err := i.requireDefinition(EPSBearerContextStatusDefinition); err != nil {
		return nil, err
	}
	return ParseEPSBearerContextStatusFields(i.Payload)
}

// ParseEPSBearerContextStatusFields decodes a value without its IEI and length
// octet.
func ParseEPSBearerContextStatusFields(payload []byte) (*EPSBearerContextStatusFields, error) {
	fields := &EPSBearerContextStatusFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that every active bearer identity is in the range 1..15
// and appears at most once.
func (f *EPSBearerContextStatusFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: EPS bearer context status fields are nil", ErrInvalidValue)
	}
	var seen uint16
	for _, bearerID := range f.ActiveBearerIDs {
		if bearerID == 0 || bearerID > 15 {
			return fmt.Errorf("%w: EPS bearer identity %d is outside 1..15", ErrInvalidValue, bearerID)
		}
		bit := uint16(1) << bearerID
		if seen&bit != 0 {
			return fmt.Errorf("%w: EPS bearer identity %d is repeated", ErrInvalidValue, bearerID)
		}
		seen |= bit
	}
	return nil
}

// Marshal encodes fields into the two-octet status value.
func (f *EPSBearerContextStatusFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, 2)
	for _, bearerID := range f.ActiveBearerIDs {
		if bearerID < 8 {
			payload[0] |= 1 << bearerID
			continue
		}
		payload[1] |= 1 << (bearerID - 8)
	}
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *EPSBearerContextStatusFields) MarshalTo(b []byte) error {
	payload, err := f.Marshal()
	if err != nil {
		return err
	}
	if len(b) < len(payload) {
		return io.ErrShortBuffer
	}
	copy(b, payload)
	return nil
}

// MarshalLen returns the size of the EPS bearer context status value payload.
func (f *EPSBearerContextStatusFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 2
}

// UnmarshalBinary decodes the two-octet status value. It leaves f unchanged on
// error.
func (f *EPSBearerContextStatusFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: EPS bearer context status fields are nil", ErrInvalidValue)
	}
	if len(payload) != 2 {
		return fmt.Errorf("%w: EPS bearer context status has %d value octets, want 2", ErrInvalidLength, len(payload))
	}
	if payload[0]&0x01 != 0 {
		return fmt.Errorf("%w: EPS bearer identity 0 is spare", ErrInvalidValue)
	}
	activeBearerIDs := make([]uint8, 0, 15)
	for bearerID := uint8(1); bearerID <= 15; bearerID++ {
		if bearerID < 8 && payload[0]&(1<<bearerID) != 0 {
			activeBearerIDs = append(activeBearerIDs, bearerID)
		}
		if bearerID >= 8 && payload[1]&(1<<(bearerID-8)) != 0 {
			activeBearerIDs = append(activeBearerIDs, bearerID)
		}
	}
	*f = EPSBearerContextStatusFields{ActiveBearerIDs: slices.Clone(activeBearerIDs)}
	return nil
}
