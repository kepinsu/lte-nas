package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestVoiceDomainPreferenceAndUEUsageSettingFields(t *testing.T) {
	t.Parallel()
	fields := &ie.VoiceDomainPreferenceAndUEUsageSettingFields{VoiceDomainPreference: true, UEUsageSetting: 5}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0xc5}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewVoiceDomainPreferenceAndUEUsageSettingWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.VoiceDomainPreferenceAndUEUsageSetting()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("VoiceDomainPreferenceAndUEUsageSetting = %+v, want %+v", decoded, fields)
	}
}

func TestVoiceDomainPreferenceAndUEUsageSettingFields_Zero(t *testing.T) {
	t.Parallel()
	fields := &ie.VoiceDomainPreferenceAndUEUsageSettingFields{VoiceDomainPreference: false, UEUsageSetting: 0}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewVoiceDomainPreferenceAndUEUsageSettingWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.VoiceDomainPreferenceAndUEUsageSetting()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("VoiceDomainPreferenceAndUEUsageSetting = %+v, want %+v", decoded, fields)
	}
}

func TestVoiceDomainPreferenceAndUEUsageSettingFields_Invalid(t *testing.T) {
	t.Parallel()
	fields := &ie.VoiceDomainPreferenceAndUEUsageSettingFields{UEUsageSetting: 8}
	if _, err := fields.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Marshal error = %v, want ErrInvalidValue", err)
	}
	if _, err := ie.ParseVoiceDomainPreferenceAndUEUsageSettingFields([]byte{0x88}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseVoiceDomainPreferenceAndUEUsageSettingFields error = %v, want ErrInvalidValue", err)
	}
	if _, err := ie.ParseVoiceDomainPreferenceAndUEUsageSettingFields([]byte{}); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("ParseVoiceDomainPreferenceAndUEUsageSettingFields error = %v, want io.ErrShortBuffer", err)
	}
}
