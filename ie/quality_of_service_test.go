package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/kepinsu/lte-nas/ie"
)

func TestQualityOfServiceFields_Marshal(t *testing.T) {
	t.Parallel()
	fields := qualityOfServiceFields()

	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x1b, 0x42, 0x06, 0x73, 0x96, 0x40, 0x80, 0x63, 0x42, 0x08, 0x08, 0x11}
	if !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	if got := fields.MarshalLen(); got != len(want) {
		t.Fatalf("MarshalLen = %d, want %d", got, len(want))
	}

	decoded, err := ie.ParseQualityOfServiceFields(payload)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("ParseQualityOfServiceFields = %+v, want %+v", decoded, fields)
	}
}

func TestQualityOfServiceFields_MarshalExtendedBitRates(t *testing.T) {
	t.Parallel()
	fields := &ie.QualityOfServiceFields{
		MaximumBitRateForUplink:      ie.QualityOfServiceBitRate{BitsPerSecond: 100_000_000},
		MaximumBitRateForDownlink:    ie.QualityOfServiceBitRate{BitsPerSecond: 300_000_000},
		GuaranteedBitRateForUplink:   ie.QualityOfServiceBitRate{BitsPerSecond: 9_000_000},
		GuaranteedBitRateForDownlink: ie.QualityOfServiceBitRate{BitsPerSecond: 17_000_000},
	}

	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0xfe, 0xfe, 0x00, 0x00, 0xfe, 0xfe, 0x00,
		0xfa, 0x4b, 0x9e, 0x04, 0x0b, 0x00,
	}
	if !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	if got := fields.MarshalLen(); got != len(want) {
		t.Fatalf("MarshalLen = %d, want %d", got, len(want))
	}

	decoded, err := ie.ParseQualityOfServiceFields(payload)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("ParseQualityOfServiceFields = %+v, want %+v", decoded, fields)
	}
}

func TestQualityOfServiceFields_MarshalUplinkExtendedBitRates(t *testing.T) {
	t.Parallel()
	fields := &ie.QualityOfServiceFields{
		MaximumBitRateForUplink: ie.QualityOfServiceBitRate{BitsPerSecond: 100_000_000},
	}

	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0xfe, 0xff, 0x00, 0x00, 0xff, 0xff, 0x00, 0x00, 0x00, 0x9e, 0x00}
	if !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
}

func TestQualityOfServiceFields_MarshalUplinkExtended2BitRates(t *testing.T) {
	t.Parallel()
	fields := &ie.QualityOfServiceFields{
		MaximumBitRateForUplink: ie.QualityOfServiceBitRate{BitsPerSecond: 300_000_000},
	}

	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0xfe, 0xff, 0x00, 0x00, 0xff, 0xff, 0x00,
		0x00, 0x00, 0xfa, 0x00, 0x00, 0x00, 0x0b, 0x00,
	}
	if !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
}

func TestQualityOfServiceFields_BitRateCoding(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name          string
		rate          ie.QualityOfServiceBitRate
		base          byte
		extended      byte
		extended2     byte
		payloadLength int
	}{
		{name: "subscribed", rate: ie.QualityOfServiceBitRate{IsSubscribed: true}, base: 0x00, payloadLength: 12},
		{name: "zero", rate: ie.QualityOfServiceBitRate{}, base: 0xff, payloadLength: 12},
		{name: "1 kbit/s", rate: ie.QualityOfServiceBitRate{BitsPerSecond: 1_000}, base: 0x01, payloadLength: 12},
		{name: "63 kbit/s", rate: ie.QualityOfServiceBitRate{BitsPerSecond: 63_000}, base: 0x3f, payloadLength: 12},
		{name: "64 kbit/s", rate: ie.QualityOfServiceBitRate{BitsPerSecond: 64_000}, base: 0x40, payloadLength: 12},
		{name: "568 kbit/s", rate: ie.QualityOfServiceBitRate{BitsPerSecond: 568_000}, base: 0x7f, payloadLength: 12},
		{name: "576 kbit/s", rate: ie.QualityOfServiceBitRate{BitsPerSecond: 576_000}, base: 0x80, payloadLength: 12},
		{name: "8.64 Mbit/s", rate: ie.QualityOfServiceBitRate{BitsPerSecond: 8_640_000}, base: 0xfe, payloadLength: 12},
		{name: "8.7 Mbit/s", rate: ie.QualityOfServiceBitRate{BitsPerSecond: 8_700_000}, base: 0xfe, extended: 0x01, payloadLength: 14},
		{name: "16 Mbit/s", rate: ie.QualityOfServiceBitRate{BitsPerSecond: 16_000_000}, base: 0xfe, extended: 0x4a, payloadLength: 14},
		{name: "17 Mbit/s", rate: ie.QualityOfServiceBitRate{BitsPerSecond: 17_000_000}, base: 0xfe, extended: 0x4b, payloadLength: 14},
		{name: "128 Mbit/s", rate: ie.QualityOfServiceBitRate{BitsPerSecond: 128_000_000}, base: 0xfe, extended: 0xba, payloadLength: 14},
		{name: "130 Mbit/s", rate: ie.QualityOfServiceBitRate{BitsPerSecond: 130_000_000}, base: 0xfe, extended: 0xbb, payloadLength: 14},
		{name: "256 Mbit/s", rate: ie.QualityOfServiceBitRate{BitsPerSecond: 256_000_000}, base: 0xfe, extended: 0xfa, payloadLength: 14},
		{name: "260 Mbit/s", rate: ie.QualityOfServiceBitRate{BitsPerSecond: 260_000_000}, base: 0xfe, extended: 0xfa, extended2: 0x01, payloadLength: 18},
		{name: "500 Mbit/s", rate: ie.QualityOfServiceBitRate{BitsPerSecond: 500_000_000}, base: 0xfe, extended: 0xfa, extended2: 0x3d, payloadLength: 18},
		{name: "510 Mbit/s", rate: ie.QualityOfServiceBitRate{BitsPerSecond: 510_000_000}, base: 0xfe, extended: 0xfa, extended2: 0x3e, payloadLength: 18},
		{name: "1.5 Gbit/s", rate: ie.QualityOfServiceBitRate{BitsPerSecond: 1_500_000_000}, base: 0xfe, extended: 0xfa, extended2: 0xa1, payloadLength: 18},
		{name: "1.6 Gbit/s", rate: ie.QualityOfServiceBitRate{BitsPerSecond: 1_600_000_000}, base: 0xfe, extended: 0xfa, extended2: 0xa2, payloadLength: 18},
		{name: "10 Gbit/s", rate: ie.QualityOfServiceBitRate{BitsPerSecond: 10_000_000_000}, base: 0xfe, extended: 0xfa, extended2: 0xf6, payloadLength: 18},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fields := &ie.QualityOfServiceFields{MaximumBitRateForDownlink: tc.rate}
			payload, err := fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if got := len(payload); got != tc.payloadLength {
				t.Fatalf("payload length = %d, want %d", got, tc.payloadLength)
			}
			if got := payload[6]; got != tc.base {
				t.Fatalf("base bit rate = %#02x, want %#02x", got, tc.base)
			}
			if len(payload) >= 14 && payload[12] != tc.extended {
				t.Fatalf("extended bit rate = %#02x, want %#02x", payload[12], tc.extended)
			}
			if len(payload) >= 18 && payload[16] != tc.extended2 {
				t.Fatalf("extended-2 bit rate = %#02x, want %#02x", payload[16], tc.extended2)
			}

			decoded, err := ie.ParseQualityOfServiceFields(payload)
			if err != nil {
				t.Fatal(err)
			}
			if got := decoded.MaximumBitRateForDownlink; got != tc.rate {
				t.Fatalf("maximum bit rate = %+v, want %+v", got, tc.rate)
			}
		})
	}
}

func TestQualityOfServiceFields_MaximumSDUSizeAndTransferDelayCoding(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name           string
		maximumSDUSize uint16
		transferDelay  time.Duration
		maximumSDUCode byte
		transferCode   byte
	}{
		{name: "subscribed", maximumSDUSize: 0, transferDelay: 0, maximumSDUCode: 0x00, transferCode: 0x00},
		{name: "10 octets and 10 ms", maximumSDUSize: 10, transferDelay: 10 * time.Millisecond, maximumSDUCode: 0x01, transferCode: 0x01},
		{name: "1500 octets and 150 ms", maximumSDUSize: 1500, transferDelay: 150 * time.Millisecond, maximumSDUCode: 0x96, transferCode: 0x0f},
		{name: "1502 octets and 200 ms", maximumSDUSize: 1502, transferDelay: 200 * time.Millisecond, maximumSDUCode: 0x97, transferCode: 0x10},
		{name: "1510 octets and 950 ms", maximumSDUSize: 1510, transferDelay: 950 * time.Millisecond, maximumSDUCode: 0x98, transferCode: 0x1f},
		{name: "1520 octets and 1 second", maximumSDUSize: 1520, transferDelay: time.Second, maximumSDUCode: 0x99, transferCode: 0x20},
		{name: "4 seconds", maximumSDUSize: 1500, transferDelay: 4 * time.Second, maximumSDUCode: 0x96, transferCode: 0x3e},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fields := &ie.QualityOfServiceFields{
				MaximumSDUSize: tc.maximumSDUSize,
				TransferDelay:  tc.transferDelay,
			}
			payload, err := fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if got := payload[4]; got != tc.maximumSDUCode {
				t.Fatalf("maximum SDU size = %#02x, want %#02x", got, tc.maximumSDUCode)
			}
			if got := payload[8] >> 2; got != tc.transferCode {
				t.Fatalf("transfer delay = %#02x, want %#02x", got, tc.transferCode)
			}

			decoded, err := ie.ParseQualityOfServiceFields(payload)
			if err != nil {
				t.Fatal(err)
			}
			if decoded.MaximumSDUSize != tc.maximumSDUSize || decoded.TransferDelay != tc.transferDelay {
				t.Fatalf("ParseQualityOfServiceFields = %+v, want maximumSDUSize=%d transferDelay=%s", decoded, tc.maximumSDUSize, tc.transferDelay)
			}
		})
	}
}

func TestNewQualityOfServiceWithFields(t *testing.T) {
	t.Parallel()
	fields := qualityOfServiceFields()
	if err := fields.Validate(); err != nil {
		t.Fatal(err)
	}

	value, err := ie.NewQualityOfServiceWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	if value.Kind != ie.QualityOfService {
		t.Fatalf("Kind = %#02x, want %#02x", value.Kind, ie.QualityOfService)
	}

	wire, err := value.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	wantWire := []byte{0x30, 0x0c, 0x1b, 0x42, 0x06, 0x73, 0x96, 0x40, 0x80, 0x63, 0x42, 0x08, 0x08, 0x11}
	if !bytes.Equal(wire, wantWire) {
		t.Fatalf("Marshal = %x, want %x", wire, wantWire)
	}

	decoded, err := value.QualityOfService()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("QualityOfService = %+v, want %+v", decoded, fields)
	}

	raw := ie.NewQualityOfService(value.Payload)
	decoded, err = raw.QualityOfService()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("raw QualityOfService = %+v, want %+v", decoded, fields)
	}
}

func TestNewPositionalQualityOfServiceWithFields(t *testing.T) {
	t.Parallel()
	value, err := ie.NewPositionalQualityOfServiceWithFields(qualityOfServiceFields())
	if err != nil {
		t.Fatal(err)
	}
	if value.Kind != ie.NoTag {
		t.Fatalf("Kind = %#02x, want NoTag", value.Kind)
	}
	if _, err := value.QualityOfService(); err != nil {
		t.Fatal(err)
	}
}

func TestQualityOfServiceFields_RejectInvalidValues(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		fields *ie.QualityOfServiceFields
	}{
		{
			name: "unrepresentable bit rate",
			fields: &ie.QualityOfServiceFields{
				MaximumBitRateForUplink: ie.QualityOfServiceBitRate{BitsPerSecond: 65_000},
			},
		},
		{
			name: "subscribed bit rate with value",
			fields: &ie.QualityOfServiceFields{
				MaximumBitRateForUplink: ie.QualityOfServiceBitRate{BitsPerSecond: 64_000, IsSubscribed: true},
			},
		},
		{
			name: "unrepresentable maximum SDU size",
			fields: &ie.QualityOfServiceFields{
				MaximumSDUSize: 1501,
			},
		},
		{
			name: "unrepresentable transfer delay",
			fields: &ie.QualityOfServiceFields{
				TransferDelay: 160 * time.Millisecond,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if _, err := tc.fields.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
				t.Fatalf("Marshal = %v, want ErrInvalidValue", err)
			}
		})
	}
}

func TestQualityOfServiceFields_UnmarshalBinary(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		payload []byte
		err     error
	}{
		{
			name:    "incomplete extended pair",
			payload: append(make([]byte, 12), 0),
			err:     ie.ErrInvalidLength,
		},
		{
			name:    "nonzero spare bits",
			payload: []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			err:     ie.ErrInvalidValue,
		},
		{
			name:    "reserved maximum SDU size",
			payload: []byte{0x00, 0x00, 0x00, 0x00, 0x9a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			err:     ie.ErrInvalidValue,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fields := *qualityOfServiceFields()
			before := fields
			if err := fields.UnmarshalBinary(tc.payload); !errors.Is(err, tc.err) {
				t.Fatalf("UnmarshalBinary = %v, want %v", err, tc.err)
			}
			if !reflect.DeepEqual(fields, before) {
				t.Fatalf("UnmarshalBinary modified fields: got %+v, want %+v", fields, before)
			}
		})
	}
}

func TestQualityOfServiceFields_MarshalTo(t *testing.T) {
	t.Parallel()
	fields := qualityOfServiceFields()
	destination := bytes.Repeat([]byte{0xaa}, fields.MarshalLen()-1)
	before := bytes.Clone(destination)

	if err := fields.MarshalTo(destination); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo = %v, want io.ErrShortBuffer", err)
	}
	if !bytes.Equal(destination, before) {
		t.Fatalf("MarshalTo modified destination: got %x, want %x", destination, before)
	}

	destination = bytes.Repeat([]byte{0xaa}, fields.MarshalLen()+1)
	if err := fields.MarshalTo(destination); err != nil {
		t.Fatal(err)
	}
	want := append([]byte{0x1b, 0x42, 0x06, 0x73, 0x96, 0x40, 0x80, 0x63, 0x42, 0x08, 0x08, 0x11}, 0xaa)
	if !bytes.Equal(destination, want) {
		t.Fatalf("MarshalTo = %x, want %x", destination, want)
	}
}

func TestIE_QualityOfServiceRejectsAnotherIE(t *testing.T) {
	t.Parallel()
	if _, err := ie.NewT3346Value(0).QualityOfService(); !errors.Is(err, ie.ErrUnexpectedIE) {
		t.Fatalf("QualityOfService = %v, want ErrUnexpectedIE", err)
	}
}

func qualityOfServiceFields() *ie.QualityOfServiceFields {
	return &ie.QualityOfServiceFields{
		DelayClass:                   ie.QualityOfServiceDelayClassThree,
		ReliabilityClass:             ie.QualityOfServiceReliabilityClassAcknowledgedRLC,
		PeakThroughputClass:          ie.QualityOfServicePeakThroughputClass8000OctetsPerSecond,
		PrecedenceClass:              ie.QualityOfServicePrecedenceClassNormal,
		MeanThroughputClass:          ie.QualityOfServiceMeanThroughputClass5000OctetsPerHour,
		TrafficClass:                 ie.QualityOfServiceTrafficClassInteractive,
		DeliveryOrder:                ie.QualityOfServiceDeliveryOrderNotRequired,
		DeliveryOfErroneousSDUs:      ie.QualityOfServiceDeliveryOfErroneousSDUsNotAllowed,
		MaximumSDUSize:               1500,
		MaximumBitRateForUplink:      ie.QualityOfServiceBitRate{BitsPerSecond: 64_000},
		MaximumBitRateForDownlink:    ie.QualityOfServiceBitRate{BitsPerSecond: 576_000},
		ResidualBitErrorRate:         ie.QualityOfServiceResidualBitErrorRateOneIn10000,
		SDUErrorRatio:                ie.QualityOfServiceSDUErrorRatioOneIn1000,
		TransferDelay:                200 * time.Millisecond,
		TrafficHandlingPriority:      ie.QualityOfServiceTrafficHandlingPriorityTwo,
		GuaranteedBitRateForUplink:   ie.QualityOfServiceBitRate{BitsPerSecond: 8_000},
		GuaranteedBitRateForDownlink: ie.QualityOfServiceBitRate{BitsPerSecond: 8_000},
		SignallingIndication:         true,
		SourceStatisticsDescriptor:   ie.QualityOfServiceSourceStatisticsDescriptorSpeech,
	}
}
