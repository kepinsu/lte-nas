package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestExtendedDRXParametersFields(t *testing.T) {
	t.Parallel()
	fields := &ie.ExtendedDRXParametersFields{PFC: true}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x80}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewExtendedDRXParametersWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.ExtendedDRXParameters()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("ExtendedDRXParameters = %+v, want %+v", decoded, fields)
	}
}

func TestExtendedDRXParametersFields_Zero(t *testing.T) {
	t.Parallel()
	fields := &ie.ExtendedDRXParametersFields{PFC: false}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewExtendedDRXParametersWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.ExtendedDRXParameters()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("ExtendedDRXParameters = %+v, want %+v", decoded, fields)
	}
}

func TestExtendedDRXParametersFields_Invalid(t *testing.T) {
	t.Parallel()
	if _, err := ie.ParseExtendedDRXParametersFields([]byte{}); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("ParseExtendedDRXParametersFields error = %v, want io.ErrShortBuffer", err)
	}
	fields := &ie.ExtendedDRXParametersFields{}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
