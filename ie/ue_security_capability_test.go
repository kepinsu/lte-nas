package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestUESecurityCapabilityFields(t *testing.T) {
	t.Parallel()
	fields := &ie.UESecurityCapabilityFields{
		EPSCipheringAlgorithms: []ie.NASEncryptionAlgorithm{
			ie.NASEncryptionAlgorithmEEA0,
			ie.NASEncryptionAlgorithm128EEA2,
		},
		EPSIntegrityAlgorithms: []ie.NASIntegrityAlgorithm{
			ie.NASIntegrityAlgorithm128EIA1,
			ie.NASIntegrityAlgorithm128EIA2,
		},
		EPSUserPlaneIntegrityProtection: true,
		UMTSCipheringAlgorithms: []ie.UMTSCipheringAlgorithm{
			ie.UMTSCipheringAlgorithmUEA1,
		},
		UMTSIntegrityAlgorithms: []ie.UMTSIntegrityAlgorithm{
			ie.UMTSIntegrityAlgorithmUIA2,
		},
		GPRSCipheringAlgorithms: []ie.GPRSCipheringAlgorithm{
			ie.GPRSCipheringAlgorithmGEA3,
		},
	}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0xa0, 0x61, 0x40, 0x40, 0x20}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewPositionalUESecurityCapabilityWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.UESecurityCapability()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("UESecurityCapability = %+v, want %+v", decoded, fields)
	}
}

func TestUESecurityCapabilityFields_GPRSOnly(t *testing.T) {
	t.Parallel()
	fields := &ie.UESecurityCapabilityFields{
		GPRSCipheringAlgorithms: []ie.GPRSCipheringAlgorithm{ie.GPRSCipheringAlgorithmGEA1},
	}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00, 0x00, 0x00, 0x00, 0x80}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
}

func TestUESecurityCapabilityFields_Invalid(t *testing.T) {
	t.Parallel()
	fields := &ie.UESecurityCapabilityFields{
		EPSIntegrityAlgorithms: []ie.NASIntegrityAlgorithm{ie.NASIntegrityAlgorithmEIA7},
	}
	if _, err := fields.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Marshal error = %v, want ErrInvalidValue", err)
	}
	if err := fields.MarshalTo(nil); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("MarshalTo error = %v, want ErrInvalidValue", err)
	}
	if _, err := ie.ParseUESecurityCapabilityFields([]byte{0x00, 0x00, 0x00}); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("ParseUESecurityCapabilityFields error = %v, want ErrInvalidLength", err)
	}

	fields = &ie.UESecurityCapabilityFields{}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
