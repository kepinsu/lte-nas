package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestSMSServicesStatusFields(t *testing.T) {
	t.Parallel()
	fields := &ie.SMSServicesStatusFields{Value: ie.SMSServicesStatusCongestion}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x03}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewSMSServicesStatusWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.SMSServicesStatus()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("SMSServicesStatus = %+v, want %+v", decoded, fields)
	}
}

func TestSMSServicesStatusFields_Zero(t *testing.T) {
	t.Parallel()
	fields := &ie.SMSServicesStatusFields{Value: ie.SMSServicesStatusNotAvailable}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewSMSServicesStatusWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.SMSServicesStatus()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("SMSServicesStatus = %+v, want %+v", decoded, fields)
	}
}

func TestSMSServicesStatusFields_Invalid(t *testing.T) {
	t.Parallel()
	fields := &ie.SMSServicesStatusFields{Value: 7}
	if _, err := fields.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Marshal error = %v, want ErrInvalidValue", err)
	}
	if _, err := ie.ParseSMSServicesStatusFields([]byte{0x10}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseSMSServicesStatusFields error = %v, want ErrInvalidValue", err)
	}
	fields2 := &ie.SMSServicesStatusFields{}
	if err := fields2.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
