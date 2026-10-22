package ie_test

import (
	"bytes"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestNASSecurityAlgorithmsFields(t *testing.T) {
	t.Parallel()
	fields := &ie.NASSecurityAlgorithmsFields{
		CipheringAlgorithm: ie.NASEncryptionAlgorithm128EEA2,
		IntegrityAlgorithm: ie.NASIntegrityAlgorithm128EIA1,
	}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x21}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	decoded, err := ie.ParseNASSecurityAlgorithmsFields(payload)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("ParseNASSecurityAlgorithmsFields = %+v, want %+v", decoded, fields)
	}
}
