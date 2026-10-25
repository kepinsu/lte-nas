package ie_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestGenericMessageContainerTypeFields(t *testing.T) {
	t.Parallel()

	for _, gmct := range []ie.GenericMessageContainerType{
		ie.GenericMessageContainerTypeESM,
	} {
		gmct := gmct
		t.Run("type", func(t *testing.T) {
			t.Parallel()
			fields := &ie.GenericMessageContainerTypeFields{Type: gmct}
			value, err := ie.NewPositionalGenericMessageContainerTypeWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := value.GenericMessageContainerType()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("GenericMessageContainerType = %+v, want %+v", decoded, fields)
			}
		})
	}

	if _, err := ie.ParseGenericMessageContainerTypeFields([]byte{0x00}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseGenericMessageContainerTypeFields = %v, want ErrInvalidValue", err)
	}

	if _, err := ie.ParseGenericMessageContainerTypeFields([]byte{0x02}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseGenericMessageContainerTypeFields = %v, want ErrInvalidValue", err)
	}

	if _, err := ie.ParseGenericMessageContainerTypeFields([]byte{0x01, 0x02}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseGenericMessageContainerTypeFields = %v, want ErrInvalidLength", err)
	}

	if _, err := ie.ParseGenericMessageContainerTypeFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseGenericMessageContainerTypeFields = %v, want ErrInvalidLength", err)
	}
}

func TestGenericMessageContainerTypeFieldsNil(t *testing.T) {
	t.Parallel()

	var f *ie.GenericMessageContainerTypeFields
	if err := f.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Validate = %v, want ErrInvalidValue", err)
	}

	if f.MarshalLen() != 0 {
		t.Fatalf("MarshalLen = %d, want 0", f.MarshalLen())
	}
}
