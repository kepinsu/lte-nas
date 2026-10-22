package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestDevicePropertiesFields(t *testing.T) {
	t.Parallel()
	fields := &ie.DevicePropertiesFields{EmergencyCallSupported: true}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x80}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewDevicePropertiesWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.DeviceProperties()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("DeviceProperties = %+v, want %+v", decoded, fields)
	}
}

func TestDevicePropertiesFields_Zero(t *testing.T) {
	t.Parallel()
	fields := &ie.DevicePropertiesFields{EmergencyCallSupported: false}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewDevicePropertiesWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.DeviceProperties()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("DeviceProperties = %+v, want %+v", decoded, fields)
	}
}

func TestDevicePropertiesFields_Invalid(t *testing.T) {
	t.Parallel()
	if _, err := ie.ParseDevicePropertiesFields([]byte{0x81}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseDevicePropertiesFields error = %v, want ErrInvalidValue", err)
	}
	fields := &ie.DevicePropertiesFields{}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
