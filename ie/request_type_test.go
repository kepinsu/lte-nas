package ie_test

import (
	"errors"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestRequestTypeFields(t *testing.T) {
	t.Parallel()

	for _, reqType := range []ie.RequestType{
		ie.RequestTypeInitial,
		ie.RequestTypeEmergency,
		ie.RequestTypePeriodic,
		ie.RequestTypeMobilized,
	} {
		reqType := reqType
		t.Run("type", func(t *testing.T) {
			t.Parallel()
			fields := &ie.RequestTypeFields{Type: reqType}
			value, err := ie.NewPositionalRequestTypeWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := value.RequestType()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("RequestType = %+v, want %+v", decoded, fields)
			}
		})
	}

	if _, err := ie.ParseRequestTypeFields([]byte{0x00}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseRequestTypeFields = %v, want ErrInvalidValue", err)
	}

	if _, err := ie.ParseRequestTypeFields([]byte{0x05}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseRequestTypeFields = %v, want ErrInvalidValue", err)
	}



	if _, err := ie.ParseRequestTypeFields([]byte{0x01, 0x02}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseRequestTypeFields = %v, want ErrInvalidLength", err)
	}

	if _, err := ie.ParseRequestTypeFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseRequestTypeFields = %v, want ErrInvalidLength", err)
	}
}

func TestRequestTypeFieldsNil(t *testing.T) {
	t.Parallel()

	var f *ie.RequestTypeFields
	if err := f.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Validate = %v, want ErrInvalidValue", err)
	}

	if f.MarshalLen() != 0 {
		t.Fatalf("MarshalLen = %d, want 0", f.MarshalLen())
	}
}
