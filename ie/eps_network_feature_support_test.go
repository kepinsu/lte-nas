package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestEPSNetworkFeatureSupportFields(t *testing.T) {
	t.Parallel()
	fields := &ie.EPSNetworkFeatureSupportFields{
		IMSVoiceOverPSSession:                   true,
		EmergencyBearerServices:                 true,
		CSLocationServices:                      ie.CSLocationServicesSupportSupported,
		ExtendedServiceRequestForPacketServices: true,
		ControlPlaneCIoT:                        true,
		UserPlaneCIoT:                           true,
		ExtendedProtocolConfigurationOptions:    true,
		DualConnectivityWithNRRestricted:        true,
		FifteenEPSBearerContexts:                true,
		PagingIndicationForVoiceServices:        true,
		PagingTimingCollisionControl:            true,
	}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0xab, 0xa9, 0x12}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewEPSNetworkFeatureSupportWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.EPSNetworkFeatureSupport()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("EPSNetworkFeatureSupport = %+v, want %+v", decoded, fields)
	}
}

func TestEPSNetworkFeatureSupportFields_Invalid(t *testing.T) {
	t.Parallel()
	fields := &ie.EPSNetworkFeatureSupportFields{CSLocationServices: 3}
	if _, err := fields.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Marshal error = %v, want ErrInvalidValue", err)
	}
	if err := fields.MarshalTo(nil); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("MarshalTo error = %v, want ErrInvalidValue", err)
	}
	if _, err := ie.ParseEPSNetworkFeatureSupportFields([]byte{0, 0, 0x80}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseEPSNetworkFeatureSupportFields error = %v, want ErrInvalidValue", err)
	}

	fields = &ie.EPSNetworkFeatureSupportFields{}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
