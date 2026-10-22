package ie_test

import (
	"errors"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestNBS1DRXParameterFields(t *testing.T) {
	t.Parallel()

	for _, nbs1drx := range []ie.Nbs1DrxParam{
		ie.Nbs1DrxParamSupported,
		ie.Nbs1DrxParamNotApplicable,
	} {
		nbs1drx := nbs1drx
		t.Run("type", func(t *testing.T) {
			t.Parallel()
			fields := &ie.NBS1DRXParameterFields{Type: nbs1drx}
			value, err := ie.NewNBS1DRXParameterWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := value.NBS1DRXParameter()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("NBS1DRXParameter = %+v, want %+v", decoded, fields)
			}
		})
	}

	if _, err := ie.ParseNBS1DRXParameterFields([]byte{0x00}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseNBS1DRXParameterFields = %v, want ErrInvalidValue", err)
	}

	if _, err := ie.ParseNBS1DRXParameterFields([]byte{0x03}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseNBS1DRXParameterFields = %v, want ErrInvalidValue", err)
	}



	if _, err := ie.ParseNBS1DRXParameterFields([]byte{0x01, 0x02}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseNBS1DRXParameterFields = %v, want ErrInvalidLength", err)
	}

	if _, err := ie.ParseNBS1DRXParameterFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseNBS1DRXParameterFields = %v, want ErrInvalidLength", err)
	}
}

func TestNBS1DRXParameterFieldsNil(t *testing.T) {
	t.Parallel()

	var f *ie.NBS1DRXParameterFields
	if err := f.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Validate = %v, want ErrInvalidValue", err)
	}

	if f.MarshalLen() != 0 {
		t.Fatalf("MarshalLen = %d, want 0", f.MarshalLen())
	}
}
