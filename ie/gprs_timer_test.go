package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/kepinsu/lte-nas/ie"
)

func TestGPRSTimerFields(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name     string
		duration time.Duration
		encoded  byte
	}{
		{name: "two seconds", duration: 2 * time.Second, encoded: 0x01},
		{name: "one minute", duration: time.Minute, encoded: 0x1e},
		{name: "long duration", duration: 186 * time.Minute, encoded: 0x5f},
		{name: "deactivated", duration: 0, encoded: 0xe0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fields := &ie.GPRSTimerFields{Duration: tc.duration}
			payload, err := fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if want := []byte{tc.encoded}; !bytes.Equal(payload, want) {
				t.Fatalf("Marshal = %x, want %x", payload, want)
			}
			value, err := ie.NewT3412ValueWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := value.GPRSTimer()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("GPRSTimer = %+v, want %+v", decoded, fields)
			}
		})
	}
}

func TestGPRSTimer2Fields(t *testing.T) {
	t.Parallel()
	fields := &ie.GPRSTimer2Fields{Duration: 31 * time.Minute}
	value, err := ie.NewT3346ValueWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.GPRSTimer2()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("GPRSTimer2 = %+v, want %+v", decoded, fields)
	}
}

func TestGPRSTimer3Fields(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name     string
		duration time.Duration
		encoded  byte
	}{
		{name: "two seconds", duration: 2 * time.Second, encoded: 0x61},
		{name: "one minute", duration: time.Minute, encoded: 0x7e},
		{name: "ten hours", duration: 10 * time.Hour, encoded: 0x2a},
		{name: "three hundred twenty hours", duration: 320 * time.Hour, encoded: 0xc1},
		{name: "deactivated", duration: 0, encoded: 0xe0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fields := &ie.GPRSTimer3Fields{Duration: tc.duration}
			payload, err := fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if want := []byte{tc.encoded}; !bytes.Equal(payload, want) {
				t.Fatalf("Marshal = %x, want %x", payload, want)
			}
			value, err := ie.NewT3412ExtendedValueWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := value.GPRSTimer3()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("GPRSTimer3 = %+v, want %+v", decoded, fields)
			}
		})
	}
}

func TestGPRSTimerFields_InvalidDurationAndShortBuffer(t *testing.T) {
	t.Parallel()
	for _, fields := range []interface {
		Marshal() ([]byte, error)
		MarshalTo([]byte) error
	}{
		&ie.GPRSTimerFields{Duration: 3 * time.Second},
		&ie.GPRSTimer2Fields{Duration: -time.Minute},
		&ie.GPRSTimer3Fields{Duration: 17 * time.Second},
	} {
		if _, err := fields.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
			t.Fatalf("Marshal error = %v, want ErrInvalidValue", err)
		}
		if err := fields.MarshalTo(nil); !errors.Is(err, ie.ErrInvalidValue) {
			t.Fatalf("MarshalTo error = %v, want ErrInvalidValue", err)
		}
	}

	fields := &ie.GPRSTimer3Fields{Duration: time.Minute}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
	if _, err := ie.ParseGPRSTimerFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseGPRSTimerFields error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseGPRSTimer2Fields([]byte{0, 0}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseGPRSTimer2Fields error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseGPRSTimer3Fields([]byte{0, 0}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseGPRSTimer3Fields error = %v, want ErrInvalidLength", err)
	}
}
