package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestWLANOffloadAcceptabilityFields(t *testing.T) {
	t.Parallel()
	fields := &ie.WLANOffloadAcceptabilityFields{OffloadIndication: true}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x80}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewWLANOffloadAcceptabilityWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.WLANOffloadAcceptability()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("WLANOffloadAcceptability = %+v, want %+v", decoded, fields)
	}
}

func TestWLANOffloadAcceptabilityFields_Zero(t *testing.T) {
	t.Parallel()
	fields := &ie.WLANOffloadAcceptabilityFields{OffloadIndication: false}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewWLANOffloadAcceptabilityWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.WLANOffloadAcceptability()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("WLANOffloadAcceptability = %+v, want %+v", decoded, fields)
	}
}

func TestWLANOffloadAcceptabilityFields_Invalid(t *testing.T) {
	t.Parallel()
	if _, err := ie.ParseWLANOffloadAcceptabilityFields([]byte{0x81}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseWLANOffloadAcceptabilityFields error = %v, want ErrInvalidValue", err)
	}
	fields := &ie.WLANOffloadAcceptabilityFields{}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
