package ie_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestLinkedEPSBearerIdentityFields(t *testing.T) {
	t.Parallel()

	for _, lbeID := range []ie.LinkedEPSBearerIdentity{
		ie.LinkedEPSBearerIdentity1,
		ie.LinkedEPSBearerIdentity3,
		ie.LinkedEPSBearerIdentity5,
		ie.LinkedEPSBearerIdentity7,
	} {
		lbeID := lbeID
		t.Run("type", func(t *testing.T) {
			t.Parallel()
			fields := &ie.LinkedEPSBearerIdentityFields{Type: lbeID}
			value, err := ie.NewPositionalLinkedEPSBearerIdentityWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := value.LinkedEPSBearerIdentity()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("LinkedEPSBearerIdentity = %+v, want %+v", decoded, fields)
			}
		})
	}

	if _, err := ie.ParseLinkedEPSBearerIdentityFields([]byte{0x00}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseLinkedEPSBearerIdentityFields = %v, want ErrInvalidValue", err)
	}

	if _, err := ie.ParseLinkedEPSBearerIdentityFields([]byte{0x08}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseLinkedEPSBearerIdentityFields = %v, want ErrInvalidValue", err)
	}

	if _, err := ie.ParseLinkedEPSBearerIdentityFields([]byte{0x01, 0x02}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseLinkedEPSBearerIdentityFields = %v, want ErrInvalidLength", err)
	}

	if _, err := ie.ParseLinkedEPSBearerIdentityFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseLinkedEPSBearerIdentityFields = %v, want ErrInvalidLength", err)
	}
}

func TestLinkedEPSBearerIdentityFieldsNil(t *testing.T) {
	t.Parallel()

	var f *ie.LinkedEPSBearerIdentityFields
	if err := f.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Validate = %v, want ErrInvalidValue", err)
	}

	if f.MarshalLen() != 0 {
		t.Fatalf("MarshalLen = %d, want 0", f.MarshalLen())
	}
}
