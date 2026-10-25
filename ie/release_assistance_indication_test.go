package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestReleaseAssistanceIndicationFields(t *testing.T) {
	t.Parallel()
	fields := &ie.ReleaseAssistanceIndicationFields{DDX: 2}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x02}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewReleaseAssistanceIndicationWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.ReleaseAssistanceIndication()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("ReleaseAssistanceIndication = %+v, want %+v", decoded, fields)
	}
}

func TestReleaseAssistanceIndicationFields_Zero(t *testing.T) {
	t.Parallel()
	fields := &ie.ReleaseAssistanceIndicationFields{DDX: 0}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewReleaseAssistanceIndicationWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.ReleaseAssistanceIndication()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("ReleaseAssistanceIndication = %+v, want %+v", decoded, fields)
	}
}

func TestReleaseAssistanceIndicationFields_Invalid(t *testing.T) {
	t.Parallel()
	fields := &ie.ReleaseAssistanceIndicationFields{DDX: 4}
	if _, err := fields.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Marshal error = %v, want ErrInvalidValue", err)
	}
	if _, err := ie.ParseReleaseAssistanceIndicationFields([]byte{0x04}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseReleaseAssistanceIndicationFields error = %v, want ErrInvalidValue", err)
	}
	fields2 := &ie.ReleaseAssistanceIndicationFields{}
	if err := fields2.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
