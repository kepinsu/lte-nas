package ie_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestPLMNListFields(t *testing.T) {
	t.Parallel()
	fields := &ie.PLMNListFields{PLMNs: []ie.PLMN{
		{MCC: "208", MNC: "93"},
		{MCC: "310", MNC: "260"},
	}}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x02, 0xf8, 0x39, 0x13, 0x00, 0x62}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewPLMNListWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.PLMNList()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("PLMNList = %+v, want %+v", decoded, fields)
	}
}

func TestPLMNListFields_Invalid(t *testing.T) {
	t.Parallel()
	if _, err := (&ie.PLMNListFields{}).Marshal(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Marshal error = %v, want ErrInvalidLength", err)
	}
	if _, err := ie.ParsePLMNListFields([]byte{0x02, 0xf8}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParsePLMNListFields error = %v, want ErrInvalidLength", err)
	}
}
