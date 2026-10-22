package ie_test

import (
	"bytes"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestNASKeySetIdentifierFields(t *testing.T) {
	t.Parallel()
	fields := &ie.NASKeySetIdentifierFields{KeySetIdentifier: 6, MappedSecurityContext: true}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x0e}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewNASKeySetIdentifierWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.NASKeySetIdentifier()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("NASKeySetIdentifier = %+v, want %+v", decoded, fields)
	}
}
