package ie_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestNetworkResourceIdentifierContainerFields(t *testing.T) {
	t.Parallel()
	fields := &ie.NetworkResourceIdentifierContainerFields{NetworkResourceIdentifier: 0x02ab}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0xaa, 0xc0}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewNetworkResourceIdentifierContainerWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.NetworkResourceIdentifierContainer()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("NetworkResourceIdentifierContainer = %+v, want %+v", decoded, fields)
	}
}

func TestNetworkResourceIdentifierContainerFields_Invalid(t *testing.T) {
	t.Parallel()
	if _, err := (&ie.NetworkResourceIdentifierContainerFields{NetworkResourceIdentifier: ie.NetworkResourceIdentifierMaximum + 1}).Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Marshal error = %v, want ErrInvalidValue", err)
	}
	if _, err := ie.ParseNetworkResourceIdentifierContainerFields([]byte{0, 1}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseNetworkResourceIdentifierContainerFields error = %v, want ErrInvalidValue", err)
	}
}
