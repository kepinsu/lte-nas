package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestSpareHalfOctetFields(t *testing.T) {
	t.Parallel()
	fields := &ie.SpareHalfOctetFields{}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewSpareHalfOctetWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.SpareHalfOctet()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("SpareHalfOctet = %+v, want %+v", decoded, fields)
	}
}

func TestSpareHalfOctetFields_Invalid(t *testing.T) {
	t.Parallel()
	if _, err := ie.ParseSpareHalfOctetFields([]byte{0x0f}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseSpareHalfOctetFields error = %v, want ErrInvalidValue", err)
	}
	if _, err := ie.ParseSpareHalfOctetFields([]byte{0x01}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseSpareHalfOctetFields error = %v, want ErrInvalidValue", err)
	}
	if _, err := ie.ParseSpareHalfOctetFields([]byte{}); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("ParseSpareHalfOctetFields error = %v, want io.ErrShortBuffer", err)
	}
}
