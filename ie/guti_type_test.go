package ie_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestGUTITypeFields(t *testing.T) {
	t.Parallel()

	for _, gutiType := range []ie.GutIType{
		ie.GutiTypeEPS,
		ie.GutITypeNonEPS,
	} {
		gutiType := gutiType
		t.Run("type", func(t *testing.T) {
			t.Parallel()
			fields := &ie.GUTITypeFields{Type: gutiType}
			value, err := ie.NewGUTITypeWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := value.GUTIType()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("GUTIType = %+v, want %+v", decoded, fields)
			}
		})
	}

	if _, err := ie.ParseGUTITypeFields([]byte{0x02}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseGUTITypeFields = %v, want ErrInvalidValue", err)
	}

	if _, err := ie.ParseGUTITypeFields([]byte{0x02}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseGUTITypeFields = %v, want ErrInvalidValue", err)
	}

	if _, err := ie.ParseGUTITypeFields([]byte{0x00, 0x01}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseGUTITypeFields = %v, want ErrInvalidLength", err)
	}

	if _, err := ie.ParseGUTITypeFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseGUTITypeFields = %v, want ErrInvalidLength", err)
	}
}

func TestGUTITypeFieldsNil(t *testing.T) {
	t.Parallel()

	var f *ie.GUTITypeFields
	if err := f.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Validate = %v, want ErrInvalidValue", err)
	}

	if f.MarshalLen() != 0 {
		t.Fatalf("MarshalLen = %d, want 0", f.MarshalLen())
	}
}
