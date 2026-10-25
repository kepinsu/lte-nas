package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestUERadioCapabilityIDRequestFields(t *testing.T) {
	t.Parallel()
	fields := &ie.UERadioCapabilityIDRequestFields{URCIDR: true}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x01}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewUERadioCapabilityIDRequestWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.UERadioCapabilityIDRequest()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("UERadioCapabilityIDRequest = %+v, want %+v", decoded, fields)
	}
}

func TestUERadioCapabilityIDRequestFields_Zero(t *testing.T) {
	t.Parallel()
	fields := &ie.UERadioCapabilityIDRequestFields{URCIDR: false}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewUERadioCapabilityIDRequestWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.UERadioCapabilityIDRequest()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("UERadioCapabilityIDRequest = %+v, want %+v", decoded, fields)
	}
}

func TestUERadioCapabilityIDRequestFields_Invalid(t *testing.T) {
	t.Parallel()
	if _, err := ie.ParseUERadioCapabilityIDRequestFields([]byte{0x02}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseUERadioCapabilityIDRequestFields error = %v, want ErrInvalidValue", err)
	}
	if _, err := ie.ParseUERadioCapabilityIDRequestFields([]byte{}); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("ParseUERadioCapabilityIDRequestFields error = %v, want io.ErrShortBuffer", err)
	}
	fields := &ie.UERadioCapabilityIDRequestFields{}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
