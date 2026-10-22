package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestTrafficFlowAggregateDescriptionFields_RoundTrip(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		payload []byte
	}{
		{name: "one byte", payload: []byte{0x01}},
		{name: "three bytes", payload: []byte{0x01, 0x02, 0x03}},
		{name: "255 bytes", payload: make([]byte, 255)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fields := &ie.TrafficFlowAggregateDescriptionFields{Value: tc.payload}

			encoded, err := fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(encoded, tc.payload) {
				t.Fatalf("Marshal = %x, want %x", encoded, tc.payload)
			}

			value, err := ie.NewPositionalTrafficFlowAggregateDescriptionWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}

			decoded, err := value.TrafficFlowAggregateDescription()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("TrafficFlowAggregateDescription = %+v, want %+v", decoded, fields)
			}
		})
	}
}

func TestTrafficFlowAggregateDescriptionFields_BelowMinError(t *testing.T) {
	t.Parallel()
	fields := &ie.TrafficFlowAggregateDescriptionFields{Value: nil}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseTrafficFlowAggregateDescriptionFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse nil: error = %v, want ErrInvalidLength", err)
	}
}

func TestTrafficFlowAggregateDescriptionFields_AboveMaxError(t *testing.T) {
	t.Parallel()
	payload := make([]byte, 256)
	for i := range payload {
		payload[i] = byte(i)
	}
	fields := &ie.TrafficFlowAggregateDescriptionFields{Value: payload}
	if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Validate: error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseTrafficFlowAggregateDescriptionFields(payload); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Parse: error = %v, want ErrInvalidLength", err)
	}
}

func TestTrafficFlowAggregateDescriptionFields_NilFields(t *testing.T) {
	t.Parallel()
	var nilFields *ie.TrafficFlowAggregateDescriptionFields
	if err := nilFields.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Validate nil: error = %v, want ErrInvalidValue", err)
	}
	_, err := nilFields.Marshal()
	if !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Marshal nil: error = %v, want ErrInvalidValue", err)
	}
	if nilFields.MarshalLen() != 0 {
		t.Fatal("MarshalLen nil: expected 0")
	}
}

func TestTrafficFlowAggregateDescriptionFields_MarshalTo(t *testing.T) {
	t.Parallel()
	payload := []byte{0x10, 0x20, 0x30}
	fields := &ie.TrafficFlowAggregateDescriptionFields{Value: payload}

	buffer := make([]byte, 10)
	if err := fields.MarshalTo(buffer); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buffer[:len(payload)], payload) {
		t.Fatalf("MarshalTo = %x, want %x", buffer[:len(payload)], payload)
	}

	short := make([]byte, 1)
	if err := fields.MarshalTo(short); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo short buffer: error = %v, want io.ErrShortBuffer", err)
	}
}
