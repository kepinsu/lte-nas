package ie_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestTrackingAreaIdentityListFields(t *testing.T) {
	t.Parallel()
	plmn := ie.PLMN{MCC: "208", MNC: "93"}
	for _, tc := range []struct {
		name    string
		fields  *ie.TrackingAreaIdentityListFields
		payload []byte
	}{
		{
			name: "non-consecutive TACs",
			fields: &ie.TrackingAreaIdentityListFields{TrackingAreas: []ie.TrackingAreaIdentityFields{
				{PLMN: plmn, TrackingAreaCode: 1},
				{PLMN: plmn, TrackingAreaCode: 3},
			}},
			payload: []byte{0x01, 0x02, 0xf8, 0x39, 0x00, 0x01, 0x00, 0x03},
		},
		{
			name: "consecutive TACs",
			fields: &ie.TrackingAreaIdentityListFields{TrackingAreas: []ie.TrackingAreaIdentityFields{
				{PLMN: plmn, TrackingAreaCode: 10},
				{PLMN: plmn, TrackingAreaCode: 11},
				{PLMN: plmn, TrackingAreaCode: 12},
			}},
			payload: []byte{0x22, 0x02, 0xf8, 0x39, 0x00, 0x0a},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			payload, err := tc.fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(payload, tc.payload) {
				t.Fatalf("Marshal = %x, want %x", payload, tc.payload)
			}
			value, err := ie.NewTAIListWithFields(tc.fields)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := value.TrackingAreaIdentityList()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, tc.fields) {
				t.Fatalf("TrackingAreaIdentityList = %+v, want %+v", decoded, tc.fields)
			}
		})
	}
}

func TestTrackingAreaIdentityListFields_DifferentPLMNs(t *testing.T) {
	t.Parallel()
	payload := []byte{0x41, 0x02, 0xf8, 0x39, 0x00, 0x01, 0x13, 0x00, 0x62, 0x00, 0x02}
	fields, err := ie.ParseTrackingAreaIdentityListFields(payload)
	if err != nil {
		t.Fatal(err)
	}
	want := &ie.TrackingAreaIdentityListFields{TrackingAreas: []ie.TrackingAreaIdentityFields{
		{PLMN: ie.PLMN{MCC: "208", MNC: "93"}, TrackingAreaCode: 1},
		{PLMN: ie.PLMN{MCC: "310", MNC: "260"}, TrackingAreaCode: 2},
	}}
	if !reflect.DeepEqual(fields, want) {
		t.Fatalf("ParseTrackingAreaIdentityListFields = %+v, want %+v", fields, want)
	}
	reencoded, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(reencoded, payload) {
		t.Fatalf("Marshal = %x, want %x", reencoded, payload)
	}
}

func TestTrackingAreaIdentityListFields_Invalid(t *testing.T) {
	t.Parallel()
	if _, err := (&ie.TrackingAreaIdentityListFields{}).Marshal(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Marshal error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParseTrackingAreaIdentityListFields([]byte{0x22, 0x02, 0xf8, 0x39, 0x00}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseTrackingAreaIdentityListFields error = %v, want ErrInvalidLength", err)
	}
}
