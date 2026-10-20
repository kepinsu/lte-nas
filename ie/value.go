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
