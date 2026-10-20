package ie

import (
	"encoding/binary"
	"fmt"
)

// ValueAsUint8 returns an exactly one-octet payload, including type-1 values.
func (i *IE) ValueAsUint8() (uint8, error) {
	if err := i.requireValueLength(1); err != nil {
		return 0, err
	}
	return i.Payload[0], nil
}

// ValueAsUint16 returns an exactly two-octet big-endian payload. IEs whose value
// is little-endian or has subfields must interpret Payload according to their clause.
func (i *IE) ValueAsUint16() (uint16, error) {
	if err := i.requireValueLength(2); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(i.Payload), nil
}

// ValueAsUint32 returns an exactly four-octet big-endian payload.
func (i *IE) ValueAsUint32() (uint32, error) {
	if err := i.requireValueLength(4); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(i.Payload), nil
}

// EMMCause returns the EMM cause using the existing MmCause enumeration.
func (i *IE) EMMCause() (MmCause, error) {
	if err := i.requireDefinition(EMMCauseDefinition); err != nil {
		return 0, err
	}
	v, err := i.ValueAsUint8()
	return MmCause(v), err
}

// ESMCause returns the ESM cause using the existing SmCause enumeration.
func (i *IE) ESMCause() (SmCause, error) {
	if err := i.requireDefinition(ESMCauseDefinition); err != nil {
		return 0, err
	}
	v, err := i.ValueAsUint8()
	return SmCause(v), err
}

func (i *IE) requireValueLength(n int) error {
	if i == nil {
		return ErrNilIE
	}
	if len(i.Payload) != n {
		return fmt.Errorf("%w: got %d value octets, want %d", ErrInvalidLength, len(i.Payload), n)
	}
	return nil
}

func (i *IE) requireDefinition(definition definitionKind) error {
	if i == nil {
		return ErrNilIE
	}
	if i.definition != definition {
		return fmt.Errorf("%w: got %s", ErrUnexpectedIE, i.Name())
	}
	return nil
}
