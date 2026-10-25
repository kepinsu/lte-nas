package ie_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestLCSIndicatorFields(t *testing.T) {
	t.Parallel()

	for _, lcsInd := range []ie.LcsInd{
		ie.LcsIndNormal,
		ie.LcsIndEmergency,
	} {
		lcsInd := lcsInd
		t.Run("type", func(t *testing.T) {
			t.Parallel()
			fields := &ie.LCSIndicatorFields{Type: lcsInd}
			value, err := ie.NewLCSIndicatorWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := value.LCSIndicator()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("LCSIndicator = %+v, want %+v", decoded, fields)
			}
		})
	}

	if _, err := ie.ParseLCSIndicatorFields([]byte{0x02}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseLCSIndicatorFields = %v, want ErrInvalidValue", err)
	}

	if _, err := ie.ParseLCSIndicatorFields([]byte{0x02}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseLCSIndicatorFields = %v, want ErrInvalidValue", err)
	}

	if _, err := ie.ParseLCSIndicatorFields([]byte{0x01, 0x02}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseLCSIndicatorFields = %v, want ErrInvalidLength", err)
	}

	if _, err := ie.ParseLCSIndicatorFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseLCSIndicatorFields = %v, want ErrInvalidLength", err)
	}
}

func TestLCSIndicatorFieldsNil(t *testing.T) {
	t.Parallel()

	var f *ie.LCSIndicatorFields
	if err := f.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Validate = %v, want ErrInvalidValue", err)
	}

	if f.MarshalLen() != 0 {
		t.Fatalf("MarshalLen = %d, want 0", f.MarshalLen())
	}
}
