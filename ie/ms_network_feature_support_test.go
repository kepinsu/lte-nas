package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestMSNetworkFeatureSupportFields(t *testing.T) {
	t.Parallel()
	fields := &ie.MSNetworkFeatureSupportFields{T3412ExtendedValue: true}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x80}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewMSNetworkFeatureSupportWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.MSNetworkFeatureSupport()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("MSNetworkFeatureSupport = %+v, want %+v", decoded, fields)
	}
}

func TestMSNetworkFeatureSupportFields_Zero(t *testing.T) {
	t.Parallel()
	fields := &ie.MSNetworkFeatureSupportFields{T3412ExtendedValue: false}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewMSNetworkFeatureSupportWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.MSNetworkFeatureSupport()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("MSNetworkFeatureSupport = %+v, want %+v", decoded, fields)
	}
}

func TestMSNetworkFeatureSupportFields_Invalid(t *testing.T) {
	t.Parallel()
	if _, err := ie.ParseMSNetworkFeatureSupportFields([]byte{0x81}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseMSNetworkFeatureSupportFields error = %v, want ErrInvalidValue", err)
	}
	fields := &ie.MSNetworkFeatureSupportFields{}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
