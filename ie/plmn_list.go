package ie

import (
	"fmt"
	"io"
)

const maximumPLMNListEntries = 15

// PLMNListFields is the structured value of the PLMN list IE (TS 24.301
// clause 9.9.2.8 and TS 24.008 clause 10.5.1.13).
type PLMNListFields struct {
	PLMNs []PLMN
}

// NewPLMNListWithFields creates a tagged PLMN list IE from fields.
func NewPLMNListWithFields(fields *PLMNListFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPLMNList(payload), nil
}

// NewPositionalPLMNListWithFields creates a positional PLMN list IE from
// fields.
func NewPositionalPLMNListWithFields(fields *PLMNListFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalPLMNList(payload), nil
}

// PLMNList returns the IE as structured fields.
func (i *IE) PLMNList() (*PLMNListFields, error) {
	if err := i.requireDefinition(PLMNListDefinition); err != nil {
		return nil, err
	}
	return ParsePLMNListFields(i.Payload)
}

// ParsePLMNListFields decodes a PLMN list value without its IEI and length
// octet.
func ParsePLMNListFields(payload []byte) (*PLMNListFields, error) {
	fields := &PLMNListFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the list has a valid 3GPP size and valid PLMN values.
func (f *PLMNListFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: PLMN list fields are nil", ErrInvalidValue)
	}
	if len(f.PLMNs) == 0 || len(f.PLMNs) > maximumPLMNListEntries {
		return fmt.Errorf("%w: PLMN list has %d entries, want 1..%d", ErrInvalidLength, len(f.PLMNs), maximumPLMNListEntries)
	}
	for index, plmn := range f.PLMNs {
		if err := plmn.validate(); err != nil {
			return fmt.Errorf("PLMN %d: %w", index, err)
		}
	}
	return nil
}

// Marshal encodes fields into a PLMN list value without its IEI and length
// octet.
func (f *PLMNListFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, f.MarshalLen())
	for index, plmn := range f.PLMNs {
		encoded, _ := plmn.marshal()
		copy(payload[index*3:], encoded[:])
	}
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *PLMNListFields) MarshalTo(b []byte) error {
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

// MarshalLen returns the size of the PLMN list value payload.
func (f *PLMNListFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return len(f.PLMNs) * 3
}

// UnmarshalBinary decodes a PLMN list value without its IEI and length octet.
// It leaves f unchanged when payload is invalid.
func (f *PLMNListFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: PLMN list fields are nil", ErrInvalidValue)
	}
	if len(payload) == 0 || len(payload)%3 != 0 || len(payload)/3 > maximumPLMNListEntries {
		return fmt.Errorf("%w: PLMN list has %d value octets, want 3..45 in three-octet entries", ErrInvalidLength, len(payload))
	}
	plmns := make([]PLMN, 0, len(payload)/3)
	for index := 0; index < len(payload); index += 3 {
		plmn, err := parsePLMN(payload[index : index+3])
		if err != nil {
			return fmt.Errorf("PLMN %d: %w", index/3, err)
		}
		plmns = append(plmns, plmn)
	}
	*f = PLMNListFields{PLMNs: plmns}
	return nil
}
