package ie_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestConnectivityTypeFields(t *testing.T) {
	t.Parallel()

	for _, connType := range []ie.ConnType{
		ie.ConnTypeNonIP,
		ie.ConnTypeEthernet,
	} {
		connType := connType
		t.Run("type", func(t *testing.T) {
			t.Parallel()
			fields := &ie.ConnectivityTypeFields{Type: connType}
			value, err := ie.NewConnectivityTypeWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := value.ConnectivityType()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("ConnectivityType = %+v, want %+v", decoded, fields)
			}
		})
	}

	if _, err := ie.ParseConnectivityTypeFields([]byte{0x00}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseConnectivityTypeFields = %v, want ErrInvalidValue", err)
	}

	if _, err := ie.ParseConnectivityTypeFields([]byte{0x03}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseConnectivityTypeFields = %v, want ErrInvalidValue", err)
	}

	if _, err := ie.ParseConnectivityTypeFields([]byte{0xe0}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseConnectivityTypeFields = %v, want ErrInvalidValue", err)
	}

	if _, err := ie.ParseConnectivityTypeFields([]byte{0x01, 0x02}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseConnectivityTypeFields = %v, want ErrInvalidLength", err)
	}

	if _, err := ie.ParseConnectivityTypeFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseConnectivityTypeFields = %v, want ErrInvalidLength", err)
	}
}

func TestConnectivityTypeFieldsNil(t *testing.T) {
	t.Parallel()

	var f *ie.ConnectivityTypeFields
	if err := f.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Validate = %v, want ErrInvalidValue", err)
	}

	if f.MarshalLen() != 0 {
		t.Fatalf("MarshalLen = %d, want 0", f.MarshalLen())
	}
}
