package ie_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestPDNTypeFields(t *testing.T) {
	t.Parallel()
	for _, pdnType := range []ie.PDNType{
		ie.PDNTypeIPv4,
		ie.PDNTypeIPv6,
		ie.PDNTypeIPv4v6,
		ie.PDNTypeNonIP,
		ie.PDNTypeEthernet,
	} {
		pdnType := pdnType
		t.Run("type", func(t *testing.T) {
			t.Parallel()
			fields := &ie.PDNTypeFields{Type: pdnType}
			value, err := ie.NewPositionalPDNTypeWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := value.PDNType()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("PDNType = %+v, want %+v", decoded, fields)
			}
		})
	}

	if _, err := ie.ParsePDNTypeFields([]byte{0x04}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParsePDNTypeFields = %v, want ErrInvalidValue", err)
	}
}
