package ie

import (
	"fmt"
	"io"
)

// TransactionIdentifierFields is the structured value of the tagged Transaction
// identifier IE (TS 24.301 clause 9.9.4.17). It carries a PTI (Personal
// Transaction Identity) between 1 and 2 octets.
type TransactionIdentifierFields struct {
	Value []byte
}

// NewTransactionIdentifierWithFields creates a tagged Transaction identifier IE
// from fields.
func NewTransactionIdentifierWithFields(fields *TransactionIdentifierFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewTransactionIdentifier(payload), nil
}

// NewPositionalTransactionIdentifierWithFields creates a positional Transaction
// identifier IE from fields.
func NewPositionalTransactionIdentifierWithFields(fields *TransactionIdentifierFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalTransactionIdentifier(payload), nil
}

// TransactionIdentifier returns the IE as structured fields.
func (i *IE) TransactionIdentifier() (*TransactionIdentifierFields, error) {
	if err := i.requireDefinition(TransactionIdentifierDefinition); err != nil {
		return nil, err
	}
	return ParseTransactionIdentifierFields(i.Payload)
}

// ParseTransactionIdentifierFields decodes a transaction identifier value without
// its IEI and length octet.
func ParseTransactionIdentifierFields(payload []byte) (*TransactionIdentifierFields, error) {
	fields := &TransactionIdentifierFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the payload is within [1, 2] octets.
func (f *TransactionIdentifierFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: transaction identifier fields are nil", ErrInvalidValue)
	}
	n := len(f.Value)
	if n < 1 {
		return fmt.Errorf("%w: transaction identifier has %d value octets, want >= 1", ErrInvalidLength, n)
	}
	if n > 2 {
		return fmt.Errorf("%w: transaction identifier has %d value octets, want <= 2", ErrInvalidLength, n)
	}
	return nil
}

// Marshal encodes fields into a transaction identifier value without its IEI and
// length octet.
func (f *TransactionIdentifierFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, len(f.Value))
	copy(payload, f.Value)
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *TransactionIdentifierFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < len(f.Value) {
		return io.ErrShortBuffer
	}
	copy(b, f.Value)
	return nil
}

// MarshalLen returns the size of the transaction identifier value payload, or zero
// for a nil receiver or nil Value.
func (f *TransactionIdentifierFields) MarshalLen() int {
	if f == nil || f.Value == nil {
		return 0
	}
	return len(f.Value)
}

// UnmarshalBinary decodes a transaction identifier value without its IEI and length
// octet. It copies the payload into a temporary struct first, leaving f unchanged
// on failure.
func (f *TransactionIdentifierFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: transaction identifier fields are nil", ErrInvalidValue)
	}
	if len(payload) < 1 {
		return fmt.Errorf("%w: transaction identifier has %d value octets, want >= 1", ErrInvalidLength, len(payload))
	}
	if len(payload) > 2 {
		return fmt.Errorf("%w: transaction identifier has %d value octets, want <= 2", ErrInvalidLength, len(payload))
	}
	var decoded TransactionIdentifierFields
	decoded.Value = make([]byte, len(payload))
	copy(decoded.Value, payload)
	*f = decoded
	return nil
}
