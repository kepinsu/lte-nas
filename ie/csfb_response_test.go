package ie_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestCSFBResponseFields(t *testing.T) {
	t.Parallel()

	for _, csfbResp := range []ie.CsfbResponse{
		ie.CsfbResponseNotApplicable,
		ie.CsfbResponseSupported,
		ie.CsfbResponseNotSupported,
	} {
		csfbResp := csfbResp
		t.Run("type", func(t *testing.T) {
			t.Parallel()
			fields := &ie.CSFBResponseFields{Type: csfbResp}
			value, err := ie.NewCSFBResponseWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := value.CSFBResponse()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("CSFBResponse = %+v, want %+v", decoded, fields)
			}
		})
	}

	if _, err := ie.ParseCSFBResponseFields([]byte{0x03}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseCSFBResponseFields = %v, want ErrInvalidValue", err)
	}

	if _, err := ie.ParseCSFBResponseFields([]byte{0x03}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseCSFBResponseFields = %v, want ErrInvalidValue", err)
	}

	if _, err := ie.ParseCSFBResponseFields([]byte{0x00, 0x01}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseCSFBResponseFields = %v, want ErrInvalidLength", err)
	}

	if _, err := ie.ParseCSFBResponseFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseCSFBResponseFields = %v, want ErrInvalidLength", err)
	}
}

func TestCSFBResponseFieldsNil(t *testing.T) {
	t.Parallel()

	var f *ie.CSFBResponseFields
	if err := f.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Validate = %v, want ErrInvalidValue", err)
	}

	if f.MarshalLen() != 0 {
		t.Fatalf("MarshalLen = %d, want 0", f.MarshalLen())
	}
}
