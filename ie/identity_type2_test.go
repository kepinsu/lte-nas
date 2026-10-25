package ie_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestIdentityType2Fields(t *testing.T) {
	t.Parallel()

	for _, idType := range []ie.IdentityType2{
		ie.IdentityType2IMEI,
		ie.IdentityType2IMEISV,
		ie.IdentityType2TMSI,
		ie.IdentityType2GUTI,
		ie.IdentityType2MSISDN,
	} {
		idType := idType
		t.Run("type", func(t *testing.T) {
			t.Parallel()
			fields := &ie.IdentityType2Fields{Type: idType}
			value, err := ie.NewPositionalIdentityType2WithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := value.IdentityType2()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("IdentityType2 = %+v, want %+v", decoded, fields)
			}
		})
	}

	if _, err := ie.ParseIdentityType2Fields([]byte{0x00}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseIdentityType2Fields = %v, want ErrInvalidValue", err)
	}

	if _, err := ie.ParseIdentityType2Fields([]byte{0x10}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseIdentityType2Fields = %v, want ErrInvalidValue", err)
	}

	if _, err := ie.ParseIdentityType2Fields([]byte{0x01, 0x02}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseIdentityType2Fields = %v, want ErrInvalidLength", err)
	}

	if _, err := ie.ParseIdentityType2Fields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseIdentityType2Fields = %v, want ErrInvalidLength", err)
	}
}

func TestIdentityType2FieldsNil(t *testing.T) {
	t.Parallel()

	var f *ie.IdentityType2Fields
	if err := f.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Validate = %v, want ErrInvalidValue", err)
	}

	if f.MarshalLen() != 0 {
		t.Fatalf("MarshalLen = %d, want 0", f.MarshalLen())
	}
}
