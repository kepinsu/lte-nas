package ie_test

import (
	"errors"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestRadioPriorityFields(t *testing.T) {
	t.Parallel()

	for _, rp := range []ie.RadioPri{
		ie.RadioPri1,
		ie.RadioPri4,
		ie.RadioPri8,
		ie.RadioPri15,
		ie.RadioPriDefault,
	} {
		rp := rp
		t.Run("type", func(t *testing.T) {
			t.Parallel()
			fields := &ie.RadioPriorityFields{Type: rp}
			value, err := ie.NewPositionalRadioPriorityWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := value.RadioPriority()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("RadioPriority = %+v, want %+v", decoded, fields)
			}
		})
	}

	if _, err := ie.ParseRadioPriorityFields([]byte{0x10}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseRadioPriorityFields = %v, want ErrInvalidValue", err)
	}



	if _, err := ie.ParseRadioPriorityFields([]byte{0x01, 0x02}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseRadioPriorityFields = %v, want ErrInvalidLength", err)
	}

	if _, err := ie.ParseRadioPriorityFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseRadioPriorityFields = %v, want ErrInvalidLength", err)
	}
}

func TestRadioPriorityFieldsNil(t *testing.T) {
	t.Parallel()

	var f *ie.RadioPriorityFields
	if err := f.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Validate = %v, want ErrInvalidValue", err)
	}

	if f.MarshalLen() != 0 {
		t.Fatalf("MarshalLen = %d, want 0", f.MarshalLen())
	}
}
