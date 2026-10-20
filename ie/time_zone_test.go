package ie_test

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"lte-nas/ie"
)

func TestTimeZoneFields(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		offset time.Duration
		want   byte
	}{
		{offset: 5*time.Hour + 30*time.Minute, want: 0x22},
		{offset: -6 * time.Hour, want: 0x4a},
	} {
		fields := &ie.TimeZoneFields{Offset: tc.offset}
		payload, err := fields.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(payload, []byte{tc.want}) {
			t.Fatalf("Marshal = %x, want %02x", payload, tc.want)
		}
		decoded, err := ie.ParseTimeZoneFields(payload)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(decoded, fields) {
			t.Fatalf("ParseTimeZoneFields = %+v, want %+v", decoded, fields)
		}
	}
}

func TestTimeZoneAndTimeFields(t *testing.T) {
	t.Parallel()
	fields := &ie.TimeZoneAndTimeFields{
		Time:   time.Date(2026, time.October, 20, 13, 45, 7, 0, time.UTC),
		Offset: -6 * time.Hour,
	}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x62, 0x01, 0x02, 0x31, 0x54, 0x70, 0x4a}
	if !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	decoded, err := ie.ParseTimeZoneAndTimeFields(payload)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("ParseTimeZoneAndTimeFields = %+v, want %+v", decoded, fields)
	}
}
