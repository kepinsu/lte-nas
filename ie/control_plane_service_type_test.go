package ie_test

import (
	"errors"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestControlPlaneServiceTypeFields(t *testing.T) {
	t.Parallel()

	for _, cpst := range []ie.ControlPlaneServiceType{
		ie.ControlPlaneServiceTypeNASRegistration,
		ie.ControlPlaneServiceTypeNASServices,
	} {
		cpst := cpst
		t.Run("type", func(t *testing.T) {
			t.Parallel()
			fields := &ie.ControlPlaneServiceTypeFields{Type: cpst}
			value, err := ie.NewPositionalControlPlaneServiceTypeWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := value.ControlPlaneServiceType()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("ControlPlaneServiceType = %+v, want %+v", decoded, fields)
			}
		})
	}

	if _, err := ie.ParseControlPlaneServiceTypeFields([]byte{0x00}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseControlPlaneServiceTypeFields = %v, want ErrInvalidValue", err)
	}

	if _, err := ie.ParseControlPlaneServiceTypeFields([]byte{0x03}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseControlPlaneServiceTypeFields = %v, want ErrInvalidValue", err)
	}



	if _, err := ie.ParseControlPlaneServiceTypeFields([]byte{0x01, 0x02}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseControlPlaneServiceTypeFields = %v, want ErrInvalidLength", err)
	}

	if _, err := ie.ParseControlPlaneServiceTypeFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseControlPlaneServiceTypeFields = %v, want ErrInvalidLength", err)
	}
}

func TestControlPlaneServiceTypeFieldsNil(t *testing.T) {
	t.Parallel()

	var f *ie.ControlPlaneServiceTypeFields
	if err := f.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Validate = %v, want ErrInvalidValue", err)
	}

	if f.MarshalLen() != 0 {
		t.Fatalf("MarshalLen = %d, want 0", f.MarshalLen())
	}
}
