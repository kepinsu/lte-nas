package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestNon3GPPNWProvidedPoliciesFields(t *testing.T) {
	t.Parallel()
	fields := &ie.Non3GPPNWProvidedPoliciesFields{PCCRule: true}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x80}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewNon3GPPNWProvidedPoliciesWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.Non3GPPNWProvidedPolicies()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("Non3GPPNWProvidedPolicies = %+v, want %+v", decoded, fields)
	}
}

func TestNon3GPPNWProvidedPoliciesFields_Zero(t *testing.T) {
	t.Parallel()
	fields := &ie.Non3GPPNWProvidedPoliciesFields{PCCRule: false}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewNon3GPPNWProvidedPoliciesWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.Non3GPPNWProvidedPolicies()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("Non3GPPNWProvidedPolicies = %+v, want %+v", decoded, fields)
	}
}

func TestNon3GPPNWProvidedPoliciesFields_Invalid(t *testing.T) {
	t.Parallel()
	if _, err := ie.ParseNon3GPPNWProvidedPoliciesFields([]byte{0x81}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseNon3GPPNWProvidedPoliciesFields error = %v, want ErrInvalidValue", err)
	}
	fields := &ie.Non3GPPNWProvidedPoliciesFields{}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
