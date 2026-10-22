package ie_test

import (
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestDRXParameterFields(t *testing.T) {
	t.Parallel()
	fields := &ie.DRXParameterFields{DRXValue: [1]byte{0x60}}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x60}; !reflect.DeepEqual(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewDRXParameterWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.DRXParameter()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("DRXParameter = %+v, want %+v", decoded, fields)
	}
}

func TestDRXParameterFields_Invalid(t *testing.T) {
	t.Parallel()
	if _, err := ie.ParseDRXParameterFields([]byte{0x00, 0x00}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseDRXParameterFields error = %v, want ErrInvalidLength", err)
	}

	fields := &ie.DRXParameterFields{}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
