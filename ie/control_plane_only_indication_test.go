package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestControlPlaneOnlyIndicationFields(t *testing.T) {
	t.Parallel()
	fields := &ie.ControlPlaneOnlyIndicationFields{CPOI: true}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x01}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewControlPlaneOnlyIndicationWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.ControlPlaneOnlyIndication()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("ControlPlaneOnlyIndication = %+v, want %+v", decoded, fields)
	}
}

func TestControlPlaneOnlyIndicationFields_Zero(t *testing.T) {
	t.Parallel()
	fields := &ie.ControlPlaneOnlyIndicationFields{CPOI: false}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewControlPlaneOnlyIndicationWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.ControlPlaneOnlyIndication()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("ControlPlaneOnlyIndication = %+v, want %+v", decoded, fields)
	}
}

func TestControlPlaneOnlyIndicationFields_Invalid(t *testing.T) {
	t.Parallel()
	if _, err := ie.ParseControlPlaneOnlyIndicationFields([]byte{0x02}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseControlPlaneOnlyIndicationFields error = %v, want ErrInvalidValue", err)
	}
	if _, err := ie.ParseControlPlaneOnlyIndicationFields([]byte{}); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("ParseControlPlaneOnlyIndicationFields error = %v, want io.ErrShortBuffer", err)
	}
	fields := &ie.ControlPlaneOnlyIndicationFields{}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
