package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestAdditionalUpdateResultFields(t *testing.T) {
	t.Parallel()
	fields := &ie.AdditionalUpdateResultFields{Value: ie.AdditionalUpdateResultSMSOnly}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x02}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewAdditionalUpdateResultWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.AdditionalUpdateResult()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("AdditionalUpdateResult = %+v, want %+v", decoded, fields)
	}
}

func TestAdditionalUpdateResultFields_Zero(t *testing.T) {
	t.Parallel()
	fields := &ie.AdditionalUpdateResultFields{Value: ie.AdditionalUpdateResultNoAdditionalInfo}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewAdditionalUpdateResultWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.AdditionalUpdateResult()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("AdditionalUpdateResult = %+v, want %+v", decoded, fields)
	}
}

func TestAdditionalUpdateResultFields_Invalid(t *testing.T) {
	t.Parallel()
	fields := &ie.AdditionalUpdateResultFields{Value: 3}
	if _, err := fields.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Marshal error = %v, want ErrInvalidValue", err)
	}
	if err := fields.MarshalTo(nil); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("MarshalTo error = %v, want ErrInvalidValue", err)
	}
	if _, err := ie.ParseAdditionalUpdateResultFields([]byte{0x04}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseAdditionalUpdateResultFields error = %v, want ErrInvalidValue", err)
	}
	fields2 := &ie.AdditionalUpdateResultFields{}
	if err := fields2.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
