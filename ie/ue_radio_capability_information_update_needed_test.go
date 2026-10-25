package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestUERadioCapabilityInformationUpdateNeededFields(t *testing.T) {
	t.Parallel()
	fields := &ie.UERadioCapabilityInformationUpdateNeededFields{URCUpd: true}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x01}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewUERadioCapabilityInformationUpdateNeededWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.UERadioCapabilityInformationUpdateNeeded()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("UERadioCapabilityInformationUpdateNeeded = %+v, want %+v", decoded, fields)
	}
}

func TestUERadioCapabilityInformationUpdateNeededFields_Zero(t *testing.T) {
	t.Parallel()
	fields := &ie.UERadioCapabilityInformationUpdateNeededFields{URCUpd: false}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewUERadioCapabilityInformationUpdateNeededWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.UERadioCapabilityInformationUpdateNeeded()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("UERadioCapabilityInformationUpdateNeeded = %+v, want %+v", decoded, fields)
	}
}

func TestUERadioCapabilityInformationUpdateNeededFields_Invalid(t *testing.T) {
	t.Parallel()
	if _, err := ie.ParseUERadioCapabilityInformationUpdateNeededFields([]byte{0x02}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseUERadioCapabilityInformationUpdateNeededFields error = %v, want ErrInvalidValue", err)
	}
	fields := &ie.UERadioCapabilityInformationUpdateNeededFields{}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
