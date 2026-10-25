package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestNotificationIndicatorFields(t *testing.T) {
	t.Parallel()
	fields := &ie.NotificationIndicatorFields{Value: 0x01}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x01}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewNotificationIndicatorWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.NotificationIndicator()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("NotificationIndicator = %+v, want %+v", decoded, fields)
	}
}

func TestNotificationIndicatorFields_ReservedValue(t *testing.T) {
	t.Parallel()
	fields := &ie.NotificationIndicatorFields{Value: 0x60}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x60}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewNotificationIndicatorWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.NotificationIndicator()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("NotificationIndicator = %+v, want %+v", decoded, fields)
	}
}

func TestNotificationIndicatorFields_Invalid(t *testing.T) {
	t.Parallel()
	if _, err := ie.ParseNotificationIndicatorFields([]byte{0x02}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseNotificationIndicatorFields error = %v, want ErrInvalidValue", err)
	}
	if _, err := ie.ParseNotificationIndicatorFields([]byte{}); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("ParseNotificationIndicatorFields error = %v, want io.ErrShortBuffer", err)
	}
	fields := &ie.NotificationIndicatorFields{}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
