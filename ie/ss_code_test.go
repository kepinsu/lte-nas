package ie_test

import (
	"errors"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestSSCodeFields(t *testing.T) {
	t.Parallel()

	for _, ssCode := range []ie.SsCode{
		ie.SsCodeValid,
		ie.SsCodeNotApplicable,
	} {
		ssCode := ssCode
		t.Run("type", func(t *testing.T) {
			t.Parallel()
			fields := &ie.SSCodeFields{Type: ssCode}
			value, err := ie.NewSSCodeWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := value.SSCode()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("SSCode = %+v, want %+v", decoded, fields)
			}
		})
	}

	if _, err := ie.ParseSSCodeFields([]byte{0x00}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseSSCodeFields = %v, want ErrInvalidValue", err)
	}

	if _, err := ie.ParseSSCodeFields([]byte{0x03}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseSSCodeFields = %v, want ErrInvalidValue", err)
	}



	if _, err := ie.ParseSSCodeFields([]byte{0x01, 0x02}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseSSCodeFields = %v, want ErrInvalidLength", err)
	}

	if _, err := ie.ParseSSCodeFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseSSCodeFields = %v, want ErrInvalidLength", err)
	}
}

func TestSSCodeFieldsNil(t *testing.T) {
	t.Parallel()

	var f *ie.SSCodeFields
	if err := f.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Validate = %v, want ErrInvalidValue", err)
	}

	if f.MarshalLen() != 0 {
		t.Fatalf("MarshalLen = %d, want 0", f.MarshalLen())
	}
}
