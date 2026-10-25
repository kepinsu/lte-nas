package ie_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestPacketFlowIdentifierFields(t *testing.T) {
	t.Parallel()

	for _, pfi := range []ie.Pfi{
		ie.PfiReserved,
		ie.PfiUEIdentity,
		ie.PfiBearerIdentity,
	} {
		pfi := pfi
		t.Run("type", func(t *testing.T) {
			t.Parallel()
			fields := &ie.PacketFlowIdentifierFields{Type: pfi}
			value, err := ie.NewPositionalPacketFlowIdentifierWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := value.PacketFlowIdentifier()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("PacketFlowIdentifier = %+v, want %+v", decoded, fields)
			}
		})
	}

	if _, err := ie.ParsePacketFlowIdentifierFields([]byte{0x03}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParsePacketFlowIdentifierFields = %v, want ErrInvalidValue", err)
	}

	if _, err := ie.ParsePacketFlowIdentifierFields([]byte{0x03}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParsePacketFlowIdentifierFields = %v, want ErrInvalidValue", err)
	}

	if _, err := ie.ParsePacketFlowIdentifierFields([]byte{0x01, 0x02}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParsePacketFlowIdentifierFields = %v, want ErrInvalidLength", err)
	}

	if _, err := ie.ParsePacketFlowIdentifierFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParsePacketFlowIdentifierFields = %v, want ErrInvalidLength", err)
	}
}

func TestPacketFlowIdentifierFieldsNil(t *testing.T) {
	t.Parallel()

	var f *ie.PacketFlowIdentifierFields
	if err := f.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Validate = %v, want ErrInvalidValue", err)
	}

	if f.MarshalLen() != 0 {
		t.Fatalf("MarshalLen = %d, want 0", f.MarshalLen())
	}
}
