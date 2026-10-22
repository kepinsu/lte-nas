package ie_test

import (
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestDCNIDFields(t *testing.T) {
	t.Parallel()
	fields := &ie.DCNIDFields{DCNID: [2]byte{0xde, 0xad}}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0xde, 0xad}; !reflect.DeepEqual(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewPositionalDCNIDWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.DCNID()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("DCNID = %+v, want %+v", decoded, fields)
	}
}

func TestDCNIDFields_Invalid(t *testing.T) {
	t.Parallel()
	if _, err := ie.ParseDCNIDFields([]byte{0x00}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseDCNIDFields error = %v, want ErrInvalidLength", err)
	}

	fields := &ie.DCNIDFields{}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
