package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestIMEISVRequestFields(t *testing.T) {
	t.Parallel()
	fields := &ie.IMEISVRequestFields{IMEISV: true}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x01}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewIMEISVRequestWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.IMEISVRequest()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("IMEISVRequest = %+v, want %+v", decoded, fields)
	}
}

func TestIMEISVRequestFields_Zero(t *testing.T) {
	t.Parallel()
	fields := &ie.IMEISVRequestFields{IMEISV: false}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewIMEISVRequestWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.IMEISVRequest()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("IMEISVRequest = %+v, want %+v", decoded, fields)
	}
}

func TestIMEISVRequestFields_Invalid(t *testing.T) {
	t.Parallel()
	if _, err := ie.ParseIMEISVRequestFields([]byte{0x02}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseIMEISVRequestFields error = %v, want ErrInvalidValue", err)
	}
	fields := &ie.IMEISVRequestFields{}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
