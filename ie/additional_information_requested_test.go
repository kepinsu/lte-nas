package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestAdditionalInformationRequestedFields(t *testing.T) {
	t.Parallel()
	fields := &ie.AdditionalInformationRequestedFields{CipherKey: true}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x01}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewAdditionalInformationRequestedWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.AdditionalInformationRequested()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("AdditionalInformationRequested = %+v, want %+v", decoded, fields)
	}
}

func TestAdditionalInformationRequestedFields_Zero(t *testing.T) {
	t.Parallel()
	fields := &ie.AdditionalInformationRequestedFields{CipherKey: false}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewAdditionalInformationRequestedWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.AdditionalInformationRequested()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("AdditionalInformationRequested = %+v, want %+v", decoded, fields)
	}
}

func TestAdditionalInformationRequestedFields_Invalid(t *testing.T) {
	t.Parallel()
	if _, err := ie.ParseAdditionalInformationRequestedFields([]byte{0x02}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseAdditionalInformationRequestedFields error = %v, want ErrInvalidValue", err)
	}
	if _, err := ie.ParseAdditionalInformationRequestedFields([]byte{0x80}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseAdditionalInformationRequestedFields error = %v, want ErrInvalidValue", err)
	}
	if _, err := ie.ParseAdditionalInformationRequestedFields([]byte{}); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("ParseAdditionalInformationRequestedFields error = %v, want io.ErrShortBuffer", err)
	}
}
