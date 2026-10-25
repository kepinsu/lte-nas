package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestNASSecurityParametersFromEUTRAFields(t *testing.T) {
	t.Parallel()
	fields := &ie.NASSecurityParametersFromEUTRAFields{KSIASME: 6}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x06}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewNASSecurityParametersFromEUTRAWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.NASSecurityParametersFromEUTRA()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("NASSecurityParametersFromEUTRA = %+v, want %+v", decoded, fields)
	}
}

func TestNASSecurityParametersFromEUTRAFields_Zero(t *testing.T) {
	t.Parallel()
	fields := &ie.NASSecurityParametersFromEUTRAFields{KSIASME: 0}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewNASSecurityParametersFromEUTRAWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.NASSecurityParametersFromEUTRA()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("NASSecurityParametersFromEUTRA = %+v, want %+v", decoded, fields)
	}
}

func TestNASSecurityParametersFromEUTRAFields_Invalid(t *testing.T) {
	t.Parallel()
	fields := &ie.NASSecurityParametersFromEUTRAFields{KSIASME: 8}
	if _, err := fields.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Marshal error = %v, want ErrInvalidValue", err)
	}
	if _, err := ie.ParseNASSecurityParametersFromEUTRAFields([]byte{0x10}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseNASSecurityParametersFromEUTRAFields error = %v, want ErrInvalidValue", err)
	}
	fields2 := &ie.NASSecurityParametersFromEUTRAFields{}
	if err := fields2.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
