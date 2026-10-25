package ie

import (
	"fmt"
	"io"
)

// CsfbResponse identifies the CS fallback response, as defined in TS 24.301
// clause 9.9.3.5.
type CsfbResponse uint8

// CSFB response values.
const (
	CsfbResponseNotApplicable CsfbResponse = 0x00
	CsfbResponseSupported     CsfbResponse = 0x01
	CsfbResponseNotSupported  CsfbResponse = 0x02
	CsfbResponseReserved      CsfbResponse = 0x03
)

// CSFBResponseFields is the structured value of the tagged CSFB response IE.
type CSFBResponseFields struct {
	Type CsfbResponse
}

// NewCSFBResponseWithFields creates a tagged CSFB response IE from fields.
func NewCSFBResponseWithFields(fields *CSFBResponseFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewCSFBResponse(payload[0]), nil
}

// CSFBResponse returns the IE as structured fields.
func (i *IE) CSFBResponse() (*CSFBResponseFields, error) {
	if err := i.requireDefinition(CSFBResponseDefinition); err != nil {
		return nil, err
	}
	return ParseCSFBResponseFields(i.Payload)
}

// ParseCSFBResponseFields decodes a CSFB response value without its IEI.
func ParseCSFBResponseFields(payload []byte) (*CSFBResponseFields, error) {
	fields := &CSFBResponseFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Type is a defined CSFB response value.
func (f *CSFBResponseFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: CSFB response fields are nil", ErrInvalidValue)
	}
	if !f.Type.valid() {
		return fmt.Errorf("%w: CSFB response %#02x is reserved", ErrInvalidValue, f.Type)
	}
	return nil
}

func (t CsfbResponse) valid() bool {
	switch t {
	case CsfbResponseNotApplicable, CsfbResponseSupported, CsfbResponseNotSupported:
		return true
	default:
		return false
	}
}

// Marshal encodes fields into a CSFB response value without its IEI.
func (f *CSFBResponseFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.Type)}, nil
}

// MarshalTo encodes fields into b.
func (f *CSFBResponseFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	b[0] = byte(f.Type)
	return nil
}

// MarshalLen returns the size of the CSFB response value payload.
func (f *CSFBResponseFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a CSFB response value without its IEI.
func (f *CSFBResponseFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: CSFB response fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: CSFB response has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	decoded := CSFBResponseFields{Type: CsfbResponse(payload[0])}
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}
