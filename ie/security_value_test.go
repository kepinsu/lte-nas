package ie_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestNonceFields(t *testing.T) {
	t.Parallel()
	fields := &ie.NonceFields{Nonce: 0x12345678}
	value, err := ie.NewNonceMMEWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.Nonce()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("Nonce = %+v, want %+v", decoded, fields)
	}
}

func TestPTMSISignatureFields(t *testing.T) {
	t.Parallel()
	fields := &ie.PTMSISignatureFields{Signature: 0x123456}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x12, 0x34, 0x56}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
}

func TestShortMACFields(t *testing.T) {
	t.Parallel()
	fields := &ie.ShortMACFields{MAC: 0x1234}
	value, err := ie.NewPositionalShortMACWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.ShortMAC()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("ShortMAC = %+v, want %+v", decoded, fields)
	}
}
