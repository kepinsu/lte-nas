package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestESMInformationTransferFlagFields(t *testing.T) {
	t.Parallel()
	fields := &ie.ESMInformationTransferFlagFields{EIT: true}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x01}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewESMInformationTransferFlagWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.ESMInformationTransferFlag()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("ESMInformationTransferFlag = %+v, want %+v", decoded, fields)
	}
}

func TestESMInformationTransferFlagFields_Zero(t *testing.T) {
	t.Parallel()
	fields := &ie.ESMInformationTransferFlagFields{EIT: false}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewESMInformationTransferFlagWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.ESMInformationTransferFlag()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("ESMInformationTransferFlag = %+v, want %+v", decoded, fields)
	}
}

func TestESMInformationTransferFlagFields_Invalid(t *testing.T) {
	t.Parallel()
	if _, err := ie.ParseESMInformationTransferFlagFields([]byte{0x02}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseESMInformationTransferFlagFields error = %v, want ErrInvalidValue", err)
	}
	if _, err := ie.ParseESMInformationTransferFlagFields([]byte{}); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("ParseESMInformationTransferFlagFields error = %v, want io.ErrShortBuffer", err)
	}
	fields := &ie.ESMInformationTransferFlagFields{}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
