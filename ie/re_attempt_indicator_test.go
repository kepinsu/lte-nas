package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestReAttemptIndicatorFields(t *testing.T) {
	t.Parallel()
	fields := &ie.ReAttemptIndicatorFields{RATC: true}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x01}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewReAttemptIndicatorWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.ReAttemptIndicator()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("ReAttemptIndicator = %+v, want %+v", decoded, fields)
	}
}

func TestReAttemptIndicatorFields_Zero(t *testing.T) {
	t.Parallel()
	fields := &ie.ReAttemptIndicatorFields{RATC: false}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewReAttemptIndicatorWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.ReAttemptIndicator()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("ReAttemptIndicator = %+v, want %+v", decoded, fields)
	}
}

func TestReAttemptIndicatorFields_Invalid(t *testing.T) {
	t.Parallel()
	if _, err := ie.ParseReAttemptIndicatorFields([]byte{0x02}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseReAttemptIndicatorFields error = %v, want ErrInvalidValue", err)
	}
	if _, err := ie.ParseReAttemptIndicatorFields([]byte{}); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("ParseReAttemptIndicatorFields error = %v, want io.ErrShortBuffer", err)
	}
	fields := &ie.ReAttemptIndicatorFields{}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
