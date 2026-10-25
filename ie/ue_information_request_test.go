package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestUEInformationRequestFields(t *testing.T) {
	t.Parallel()
	fields := &ie.UEInformationRequestFields{UCLIR: true}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x01}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewUEInformationRequestWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.UEInformationRequest()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("UEInformationRequest = %+v, want %+v", decoded, fields)
	}
}

func TestUEInformationRequestFields_Zero(t *testing.T) {
	t.Parallel()
	fields := &ie.UEInformationRequestFields{UCLIR: false}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewUEInformationRequestWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.UEInformationRequest()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("UEInformationRequest = %+v, want %+v", decoded, fields)
	}
}

func TestUEInformationRequestFields_Invalid(t *testing.T) {
	t.Parallel()
	if _, err := ie.ParseUEInformationRequestFields([]byte{0x02}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseUEInformationRequestFields error = %v, want ErrInvalidValue", err)
	}
	if _, err := ie.ParseUEInformationRequestFields([]byte{}); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("ParseUEInformationRequestFields error = %v, want io.ErrShortBuffer", err)
	}
	fields := &ie.UEInformationRequestFields{}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
