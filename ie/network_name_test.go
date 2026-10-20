package ie_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestNetworkNameFields(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		fields  *ie.NetworkNameFields
		payload []byte
	}{
		{
			name: "GSM 7-bit",
			fields: &ie.NetworkNameFields{
				Text:               "Test",
				Encoding:           ie.NetworkNameEncodingGSM7,
				AddCountryInitials: false,
			},
			payload: []byte{0x84, 0xd4, 0xf2, 0x9c, 0x0e},
		},
		{
			name: "GSM 7-bit extension",
			fields: &ie.NetworkNameFields{
				Text:               "€",
				Encoding:           ie.NetworkNameEncodingGSM7,
				AddCountryInitials: true,
			},
			payload: []byte{0x8a, 0x9b, 0x32},
		},
		{
			name: "UCS2",
			fields: &ie.NetworkNameFields{
				Text:     "É",
				Encoding: ie.NetworkNameEncodingUCS2,
			},
			payload: []byte{0x90, 0x00, 0xc9},
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
			value, err := ie.NewFullNetworkNameWithFields(tc.fields)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := value.NetworkName()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, tc.fields) {
				t.Fatalf("NetworkName = %+v, want %+v", decoded, tc.fields)
			}
		})
	}
}

func TestNetworkNameFields_Invalid(t *testing.T) {
	t.Parallel()
	fields := &ie.NetworkNameFields{Text: "`", Encoding: ie.NetworkNameEncodingGSM7}
	if _, err := fields.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Marshal error = %v, want ErrInvalidValue", err)
	}
	if _, err := ie.ParseNetworkNameFields([]byte{0x84, 0xd4, 0xf2, 0x9c, 0x8e}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseNetworkNameFields error = %v, want ErrInvalidValue", err)
	}
}
