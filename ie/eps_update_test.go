package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestEPSAttachResultFields(t *testing.T) {
	t.Parallel()
	for _, result := range []ie.EPSAttachResult{
		ie.EPSAttachResultEPSOnly,
		ie.EPSAttachResultCombinedEPSAndIMSI,
	} {
		result := result
		t.Run("result", func(t *testing.T) {
			t.Parallel()
			fields := &ie.EPSAttachResultFields{Result: result}
			value, err := ie.NewPositionalEPSAttachResultWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := value.EPSAttachResult()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("EPSAttachResult = %+v, want %+v", decoded, fields)
			}
		})
	}

	if _, err := (&ie.EPSAttachResultFields{}).Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Marshal = %v, want ErrInvalidValue", err)
	}
}

func TestEPSUpdateResultFields(t *testing.T) {
	t.Parallel()
	for _, result := range []ie.EPSUpdateResult{
		ie.EPSUpdateResultTrackingAreaUpdated,
		ie.EPSUpdateResultCombinedTrackingAreaAndLocationAreaUpdated,
		ie.EPSUpdateResultTrackingAreaUpdatedWithISR,
		ie.EPSUpdateResultCombinedTrackingAreaAndLocationAreaUpdatedWithISR,
	} {
		result := result
		t.Run("result", func(t *testing.T) {
			t.Parallel()
			fields := &ie.EPSUpdateResultFields{Result: result}
			payload, err := fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := ie.ParseEPSUpdateResultFields(payload)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("ParseEPSUpdateResultFields = %+v, want %+v", decoded, fields)
			}
		})
	}
}

func TestEPSUpdateTypeFields(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		typeValue ie.EPSUpdateType
		active    bool
		want      byte
	}{
		{typeValue: ie.EPSUpdateTypeTrackingArea, want: 0x00},
		{typeValue: ie.EPSUpdateTypeCombinedTrackingAreaAndLocationArea, active: true, want: 0x09},
		{typeValue: ie.EPSUpdateTypeCombinedTrackingAreaAndLocationAreaWithIMSI, want: 0x02},
		{typeValue: ie.EPSUpdateTypePeriodic, active: true, want: 0x0b},
	} {
		tc := tc
		t.Run("type", func(t *testing.T) {
			t.Parallel()
			fields := &ie.EPSUpdateTypeFields{Type: tc.typeValue, Active: tc.active}
			payload, err := fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(payload, []byte{tc.want}) {
				t.Fatalf("Marshal = %x, want %02x", payload, tc.want)
			}
			value, err := ie.NewPositionalEPSUpdateTypeWithFields(fields)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := value.EPSUpdateType()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(decoded, fields) {
				t.Fatalf("EPSUpdateType = %+v, want %+v", decoded, fields)
			}
		})
	}

	fields := &ie.EPSUpdateTypeFields{Type: ie.EPSUpdateTypeTrackingArea}
	buffer := []byte{0xaa}
	if err := fields.MarshalTo(buffer[:0]); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo = %v, want io.ErrShortBuffer", err)
	}
	if buffer[0] != 0xaa {
		t.Fatalf("MarshalTo modified short buffer: got %x, want aa", buffer)
	}
	if _, err := ie.ParseEPSUpdateTypeFields([]byte{0x04}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseEPSUpdateTypeFields = %v, want ErrInvalidValue", err)
	}
}
