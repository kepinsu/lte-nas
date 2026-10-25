package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestUERadioCapabilityIDDeletionIndicationFields(t *testing.T) {
	t.Parallel()
	fields := &ie.UERadioCapabilityIDDeletionIndicationFields{URCDI: true}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x80}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewUERadioCapabilityIDDeletionIndicationWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.UERadioCapabilityIDDeletionIndication()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("UERadioCapabilityIDDeletionIndication = %+v, want %+v", decoded, fields)
	}
}

func TestUERadioCapabilityIDDeletionIndicationFields_Zero(t *testing.T) {
	t.Parallel()
	fields := &ie.UERadioCapabilityIDDeletionIndicationFields{URCDI: false}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewUERadioCapabilityIDDeletionIndicationWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.UERadioCapabilityIDDeletionIndication()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("UERadioCapabilityIDDeletionIndication = %+v, want %+v", decoded, fields)
	}
}

func TestUERadioCapabilityIDDeletionIndicationFields_Invalid(t *testing.T) {
	t.Parallel()
	if _, err := ie.ParseUERadioCapabilityIDDeletionIndicationFields([]byte{0x81}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseUERadioCapabilityIDDeletionIndicationFields error = %v, want ErrInvalidValue", err)
	}
	fields := &ie.UERadioCapabilityIDDeletionIndicationFields{}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
