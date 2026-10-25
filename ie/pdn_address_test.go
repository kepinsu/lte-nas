package ie_test

import (
	"bytes"
	"net/netip"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestPDNAddressFields(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		fields *ie.PDNAddressFields
		want   []byte
	}{
		{
			name:   "IPv4",
			fields: &ie.PDNAddressFields{Type: ie.PDNTypeIPv4, IPv4: netip.MustParseAddr("192.0.2.1")},
			want:   []byte{0x01, 192, 0, 2, 1},
		},
		{
			name: "IPv6 interface identifier",
			fields: &ie.PDNAddressFields{
				Type:                    ie.PDNTypeIPv6,
				IPv6InterfaceIdentifier: [8]byte{0, 1, 2, 3, 4, 5, 6, 7},
			},
			want: []byte{0x02, 0, 1, 2, 3, 4, 5, 6, 7},
		},
		{
			name: "IPv4v6",
			fields: &ie.PDNAddressFields{
				Type:                    ie.PDNTypeIPv4v6,
				IPv4:                    netip.MustParseAddr("198.51.100.2"),
				IPv6InterfaceIdentifier: [8]byte{0, 1, 2, 3, 4, 5, 6, 7},
			},
			want: []byte{0x03, 0, 1, 2, 3, 4, 5, 6, 7, 198, 51, 100, 2},
		},
		{
			name:   "non-IP",
			fields: &ie.PDNAddressFields{Type: ie.PDNTypeNonIP},
			want:   []byte{0x05, 0, 0, 0, 0},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			payload, err := tc.fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(payload, tc.want) {
				t.Fatalf("Marshal = %x, want %x", payload, tc.want)
			}
			decoded, err := ie.ParsePDNAddressFields(payload)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, tc.fields) {
				t.Fatalf("ParsePDNAddressFields = %+v, want %+v", decoded, tc.fields)
			}
		})
	}
}
