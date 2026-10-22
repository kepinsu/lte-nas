package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestTransactionIdentifierFields_RoundTrip(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		payload []byte
	}{
		{name: "one byte", payload: []byte{0x01}},
		{name: "two bytes", payload: []byte{0x01, 0x02}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fields := &ie.TransactionIdentifierFields{Value: tc.payload}

			encoded, err := fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(encoded, tc.payload) {
				t.Fatalf("Marshal = %x, want %x", encoded, tc.payload)
			}

			value, err := ie.NewTransactionIdentifierWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}

			decoded, err := value.TransactionIdentifier()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("TransactionIdentifier = %+v, want %+v", decoded, fields)
			}
		})
	}
}

func TestTransactionIdentifierFields_PositionalRoundTrip(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		payload []byte
	}{
		{name: "one byte", payload: []byte{0xaa}},
		{name: "two bytes", payload: []byte{0xaa, 0xbb}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fields := &ie.TransactionIdentifierFields{Value: tc.payload}

			value, err := ie.NewPositionalTransactionIdentifierWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}

			decoded, err := value.TransactionIdentifier()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("TransactionIdentifier = %+v, want %+v", decoded, fields)
			}
		})
	}
}

func TestTransactionIdentifierFields_BelowMinError(t *testing.T) {
	t.Parallel()
	fields := &ie.TransactionIdentifierFields{Value: nil}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseTransactionIdentifierFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse nil: error = %v, want ErrInvalidLength", err)
	}
}

func TestTransactionIdentifierFields_AboveMaxError(t *testing.T) {
	t.Parallel()
	payload := []byte{0x01, 0x02, 0x03}
	fields := &ie.TransactionIdentifierFields{Value: payload}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseTransactionIdentifierFields(payload); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse: error = %v, want ErrInvalidLength", err)
	}
}

func TestTransactionIdentifierFields_NilFields(t *testing.T) {
	t.Parallel()
	var nilFields *ie.TransactionIdentifierFields
	if err := nilFields.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Validate nil: error = %v, want ErrInvalidValue", err)
	}
	_, err := nilFields.Marshal()
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Marshal nil: error = %v, want ErrInvalidValue", err)
	}
	if nilFields.MarshalLen() != 0 {
		t.Fatal("MarshalLen nil: expected 0")
	}
}

func TestTransactionIdentifierFields_MarshalTo(t *testing.T) {
	t.Parallel()
	payload := []byte{0x10, 0x20}
	fields := &ie.TransactionIdentifierFields{Value: payload}

	buffer := make([]byte, 10)
	if err := fields.MarshalTo(buffer); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buffer[:len(payload)], payload) {
		t.Fatalf("MarshalTo = %x, want %x", buffer[:len(payload)], payload)
	}

	short := make([]byte, 1)
	if err := fields.MarshalTo(short); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo short buffer: error = %v, want io.ErrShortBuffer", err)
	}
}
