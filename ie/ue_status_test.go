package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestUEStatusFields(t *testing.T) {
	t.Parallel()
	fields := &ie.UEStatusFields{EMMRegistrationStatus: true, FiveGMMRegistrationStatus: true}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x03}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewUEStatusWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.UEStatus()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("UEStatus = %+v, want %+v", decoded, fields)
	}
}

func TestUEStatusFields_Zero(t *testing.T) {
	t.Parallel()
	fields := &ie.UEStatusFields{EMMRegistrationStatus: false, FiveGMMRegistrationStatus: false}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewUEStatusWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.UEStatus()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("UEStatus = %+v, want %+v", decoded, fields)
	}
}

func TestUEStatusFields_Invalid(t *testing.T) {
	t.Parallel()
	if _, err := ie.ParseUEStatusFields([]byte{0x04}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseUEStatusFields error = %v, want ErrInvalidValue", err)
	}
	fields := &ie.UEStatusFields{}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
