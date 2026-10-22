package ie_test

import (
	"bytes"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestEPSQualityOfServiceFields(t *testing.T) {
	t.Parallel()
	uplinkMaximum := ie.QualityOfServiceBitRate{BitsPerSecond: 64_000}
	downlinkMaximum := ie.QualityOfServiceBitRate{BitsPerSecond: 100_000_000}
	uplinkGuaranteed := ie.QualityOfServiceBitRate{BitsPerSecond: 8_000}
	downlinkGuaranteed := ie.QualityOfServiceBitRate{BitsPerSecond: 300_000_000}
	fields := &ie.EPSQualityOfServiceFields{
		QCI:                          9,
		MaximumBitRateForUplink:      &uplinkMaximum,
		MaximumBitRateForDownlink:    &downlinkMaximum,
		GuaranteedBitRateForUplink:   &uplinkGuaranteed,
		GuaranteedBitRateForDownlink: &downlinkGuaranteed,
	}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x09, 0x40, 0xfe, 0x08, 0xfe, 0x00, 0x9e, 0x00, 0xfa, 0x00, 0x00, 0x00, 0x0b}
	if !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewEPSQualityOfServiceWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.EPSQualityOfService()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("EPSQualityOfService = %+v, want %+v", decoded, fields)
	}
}

func TestEPSQualityOfServiceFields_WithoutBitRates(t *testing.T) {
	t.Parallel()
	fields := &ie.EPSQualityOfServiceFields{QCI: 0}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(payload, []byte{0}) {
		t.Fatalf("Marshal = %x, want 00", payload)
	}
}
