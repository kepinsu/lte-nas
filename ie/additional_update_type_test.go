package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestAdditionalUpdateTypeFields(t *testing.T) {
	t.Parallel()
	fields := &ie.AdditionalUpdateTypeFields{AUTV: ie.AdditionalUpdateTypeSMSOnly, SAF: true, PNBsCIoT: 2}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x45}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewAdditionalUpdateTypeWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.AdditionalUpdateType()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("AdditionalUpdateType = %+v, want %+v", decoded, fields)
	}
}

func TestAdditionalUpdateTypeFields_Zero(t *testing.T) {
	t.Parallel()
	fields := &ie.AdditionalUpdateTypeFields{AUTV: ie.AdditionalUpdateTypeNoAdditionalInfo}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewAdditionalUpdateTypeWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.AdditionalUpdateType()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("AdditionalUpdateType = %+v, want %+v", decoded, fields)
	}
}

func TestAdditionalUpdateTypeFields_Invalid(t *testing.T) {
	t.Parallel()
	fields := &ie.AdditionalUpdateTypeFields{AUTV: 1, PNBsCIoT: 4}
	if _, err := fields.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Marshal error = %v, want ErrInvalidValue", err)
	}
	if _, err := ie.ParseAdditionalUpdateTypeFields([]byte{0x10}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseAdditionalUpdateTypeFields error = %v, want ErrInvalidValue", err)
	}
	fields2 := &ie.AdditionalUpdateTypeFields{}
	if err := fields2.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
