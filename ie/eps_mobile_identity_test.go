package ie_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestEPSMobileIdentityFields_Digits(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name     string
		identity *ie.EPSMobileIdentityFields
		want     []byte
	}{
		{
			name:     "IMSI with odd digits",
			identity: &ie.EPSMobileIdentityFields{Type: ie.EPSMobileIdentityTypeIMSI, Digits: "12345"},
			want:     []byte{0x19, 0x32, 0x54},
		},
		{
			name:     "IMEI with even digits",
			identity: &ie.EPSMobileIdentityFields{Type: ie.EPSMobileIdentityTypeIMEI, Digits: "1234"},
			want:     []byte{0x13, 0x32, 0xf4},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			payload, err := tc.identity.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(payload, tc.want) {
				t.Fatalf("Marshal = %x, want %x", payload, tc.want)
			}
			decoded, err := ie.ParseEPSMobileIdentityFields(payload)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, tc.identity) {
				t.Fatalf("ParseEPSMobileIdentityFields = %+v, want %+v", decoded, tc.identity)
			}
		})
	}
}

func TestEPSMobileIdentityFields_GUTI(t *testing.T) {
	t.Parallel()
	fields := &ie.EPSMobileIdentityFields{
		Type: ie.EPSMobileIdentityTypeGUTI,
		GUTI: &ie.GUTI{
			PLMN:       ie.PLMN{MCC: "208", MNC: "01"},
			MMEGroupID: 0x1234,
			MMECode:    0x56,
			MTMSI:      0x789abcde,
		},
	}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0xf6, 0x02, 0xf8, 0x10, 0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde}
	if !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewEPSMobileIdentityWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.EPSMobileIdentity()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("EPSMobileIdentity = %+v, want %+v", decoded, fields)
	}
}

func TestEPSMobileIdentityFields_RejectInvalid(t *testing.T) {
	t.Parallel()
	for _, fields := range []*ie.EPSMobileIdentityFields{
		{Type: ie.EPSMobileIdentityTypeIMSI},
		{Type: ie.EPSMobileIdentityTypeIMEI, Digits: "12a"},
		{Type: ie.EPSMobileIdentityTypeGUTI, Digits: "123"},
		{Type: ie.EPSMobileIdentityTypeGUTI, GUTI: &ie.GUTI{PLMN: ie.PLMN{MCC: "20", MNC: "01"}}},
	} {
		if _, err := fields.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
			t.Fatalf("Marshal = %v, want ErrInvalidValue", err)
		}
	}
	if _, err := ie.ParseEPSMobileIdentityFields([]byte{0x11, 0x32, 0x34}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseEPSMobileIdentityFields = %v, want ErrInvalidValue", err)
	}
}
