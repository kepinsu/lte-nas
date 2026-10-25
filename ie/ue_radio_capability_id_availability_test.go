package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestUERadioCapabilityIDAvailabilityFields(t *testing.T) {
	t.Parallel()
	fields := &ie.UERadioCapabilityIDAvailabilityFields{URCA: true}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x02}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewUERadioCapabilityIDAvailabilityWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.UERadioCapabilityIDAvailability()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("UERadioCapabilityIDAvailability = %+v, want %+v", decoded, fields)
	}
}

func TestUERadioCapabilityIDAvailabilityFields_Zero(t *testing.T) {
	t.Parallel()
	fields := &ie.UERadioCapabilityIDAvailabilityFields{URCA: false}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewUERadioCapabilityIDAvailabilityWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.UERadioCapabilityIDAvailability()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("UERadioCapabilityIDAvailability = %+v, want %+v", decoded, fields)
	}
}

func TestUERadioCapabilityIDAvailabilityFields_Invalid(t *testing.T) {
	t.Parallel()
	if _, err := ie.ParseUERadioCapabilityIDAvailabilityFields([]byte{0x10}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseUERadioCapabilityIDAvailabilityFields error = %v, want ErrInvalidValue", err)
	}
	if _, err := ie.ParseUERadioCapabilityIDAvailabilityFields([]byte{}); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("ParseUERadioCapabilityIDAvailabilityFields error = %v, want io.ErrShortBuffer", err)
	}
	fields := &ie.UERadioCapabilityIDAvailabilityFields{}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
