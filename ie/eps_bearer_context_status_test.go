package ie_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestEPSBearerContextStatusFields(t *testing.T) {
	t.Parallel()
	fields := &ie.EPSBearerContextStatusFields{ActiveBearerIDs: []uint8{1, 5, 7, 8, 15}}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0xa2, 0x81}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}

	value, err := ie.NewEPSBearerContextStatusWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.EPSBearerContextStatus()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("EPSBearerContextStatus = %+v, want %+v", decoded, fields)
	}

	decoded.ActiveBearerIDs[0] = 2
	decodedAgain, err := value.EPSBearerContextStatus()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decodedAgain, fields) {
		t.Fatalf("EPSBearerContextStatus exposed mutable payload: got %+v, want %+v", decodedAgain, fields)
	}
}

func TestEPSBearerContextStatusFields_RejectInvalid(t *testing.T) {
	t.Parallel()
	for _, fields := range []*ie.EPSBearerContextStatusFields{
		{ActiveBearerIDs: []uint8{0}},
		{ActiveBearerIDs: []uint8{16}},
		{ActiveBearerIDs: []uint8{5, 5}},
	} {
		if _, err := fields.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
			t.Fatalf("Marshal = %v, want ErrInvalidValue", err)
		}
	}
	if _, err := ie.ParseEPSBearerContextStatusFields([]byte{0x01, 0x00}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseEPSBearerContextStatusFields = %v, want ErrInvalidValue", err)
	}
}
