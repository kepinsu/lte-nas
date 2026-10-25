package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestDaylightSavingTimeFields(t *testing.T) {
	t.Parallel()
	fields := &ie.DaylightSavingTimeFields{Offset: 30}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x1e}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewDaylightSavingTimeWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.DaylightSavingTime()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("DaylightSavingTime = %+v, want %+v", decoded, fields)
	}
}

func TestDaylightSavingTimeFields_Negative(t *testing.T) {
	t.Parallel()
	fields := &ie.DaylightSavingTimeFields{Offset: -30}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x9e}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewDaylightSavingTimeWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.DaylightSavingTime()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("DaylightSavingTime = %+v, want %+v", decoded, fields)
	}
}

func TestDaylightSavingTimeFields_Zero(t *testing.T) {
	t.Parallel()
	fields := &ie.DaylightSavingTimeFields{Offset: 0}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewDaylightSavingTimeWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.DaylightSavingTime()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("DaylightSavingTime = %+v, want %+v", decoded, fields)
	}
}

func TestDaylightSavingTimeFields_Invalid(t *testing.T) {
	t.Parallel()
	if _, err := ie.ParseDaylightSavingTimeFields([]byte{}); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("ParseDaylightSavingTimeFields error = %v, want io.ErrShortBuffer", err)
	}
	fields := &ie.DaylightSavingTimeFields{}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
