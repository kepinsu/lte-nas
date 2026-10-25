package ie_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestLLCServiceAccessPointIdentifierFields(t *testing.T) {
	t.Parallel()

	for _, llcSapi := range []ie.LlcSapi{
		ie.LlcSapiReserved,
		ie.LlcSapiLLC,
	} {
		llcSapi := llcSapi
		t.Run("type", func(t *testing.T) {
			t.Parallel()
			fields := &ie.LLCServiceAccessPointIdentifierFields{Type: llcSapi}
			value, err := ie.NewLLCServiceAccessPointIdentifierWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := value.LLCServiceAccessPointIdentifier()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("LLCServiceAccessPointIdentifier = %+v, want %+v", decoded, fields)
			}
		})
	}

	if _, err := ie.ParseLLCServiceAccessPointIdentifierFields([]byte{0x02}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseLLCServiceAccessPointIdentifierFields = %v, want ErrInvalidValue", err)
	}

	if _, err := ie.ParseLLCServiceAccessPointIdentifierFields([]byte{0x02}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseLLCServiceAccessPointIdentifierFields = %v, want ErrInvalidValue", err)
	}

	if _, err := ie.ParseLLCServiceAccessPointIdentifierFields([]byte{0x01, 0x02}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseLLCServiceAccessPointIdentifierFields = %v, want ErrInvalidLength", err)
	}

	if _, err := ie.ParseLLCServiceAccessPointIdentifierFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseLLCServiceAccessPointIdentifierFields = %v, want ErrInvalidLength", err)
	}
}

func TestLLCServiceAccessPointIdentifierFieldsNil(t *testing.T) {
	t.Parallel()

	var f *ie.LLCServiceAccessPointIdentifierFields
	if err := f.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Validate = %v, want ErrInvalidValue", err)
	}

	if f.MarshalLen() != 0 {
		t.Fatalf("MarshalLen = %d, want 0", f.MarshalLen())
	}
}
