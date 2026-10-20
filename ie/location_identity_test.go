package ie_test

import (
	"bytes"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestLocationAreaIdentificationFields(t *testing.T) {
	t.Parallel()
	fields := &ie.LocationAreaIdentificationFields{
		PLMN:             ie.PLMN{MCC: "208", MNC: "01"},
		LocationAreaCode: 0x1234,
	}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x02, 0xf8, 0x10, 0x12, 0x34}
	if !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewLocationAreaIdentificationWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.LocationAreaIdentification()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("LocationAreaIdentification = %+v, want %+v", decoded, fields)
	}
}

func TestTrackingAreaIdentityFields(t *testing.T) {
	t.Parallel()
	fields := &ie.TrackingAreaIdentityFields{
		PLMN:             ie.PLMN{MCC: "310", MNC: "260"},
		TrackingAreaCode: 0xabcd,
	}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x13, 0x00, 0x62, 0xab, 0xcd}
	if !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	decoded, err := ie.ParseTrackingAreaIdentityFields(payload)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("ParseTrackingAreaIdentityFields = %+v, want %+v", decoded, fields)
	}
}
