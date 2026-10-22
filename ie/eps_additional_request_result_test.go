package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestEPSAdditionalRequestResultFields(t *testing.T) {
	t.Parallel()
	fields := &ie.EPSAdditionalRequestResultFields{PRD: 2}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x02}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewEPSAdditionalRequestResultWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.EPSAdditionalRequestResult()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("EPSAdditionalRequestResult = %+v, want %+v", decoded, fields)
	}
}

func TestEPSAdditionalRequestResultFields_Zero(t *testing.T) {
	t.Parallel()
	fields := &ie.EPSAdditionalRequestResultFields{PRD: 0}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewEPSAdditionalRequestResultWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.EPSAdditionalRequestResult()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("EPSAdditionalRequestResult = %+v, want %+v", decoded, fields)
	}
}

func TestEPSAdditionalRequestResultFields_Invalid(t *testing.T) {
	t.Parallel()
	fields := &ie.EPSAdditionalRequestResultFields{PRD: 4}
	if _, err := fields.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Marshal error = %v, want ErrInvalidValue", err)
	}
	if _, err := ie.ParseEPSAdditionalRequestResultFields([]byte{0x04}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseEPSAdditionalRequestResultFields error = %v, want ErrInvalidValue", err)
	}
	if _, err := ie.ParseEPSAdditionalRequestResultFields([]byte{}); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("ParseEPSAdditionalRequestResultFields error = %v, want io.ErrShortBuffer", err)
	}
	fields2 := &ie.EPSAdditionalRequestResultFields{}
	if err := fields2.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
