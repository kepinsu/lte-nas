package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestNetworkPolicyFields(t *testing.T) {
	t.Parallel()
	fields := &ie.NetworkPolicyFields{RedirPolicy: true}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x01}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewNetworkPolicyWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.NetworkPolicy()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("NetworkPolicy = %+v, want %+v", decoded, fields)
	}
}

func TestNetworkPolicyFields_Zero(t *testing.T) {
	t.Parallel()
	fields := &ie.NetworkPolicyFields{RedirPolicy: false}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x00}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
	value, err := ie.NewNetworkPolicyWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.NetworkPolicy()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("NetworkPolicy = %+v, want %+v", decoded, fields)
	}
}

func TestNetworkPolicyFields_Invalid(t *testing.T) {
	t.Parallel()
	if _, err := ie.ParseNetworkPolicyFields([]byte{0x02}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseNetworkPolicyFields error = %v, want ErrInvalidValue", err)
	}
	if _, err := ie.ParseNetworkPolicyFields([]byte{}); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("ParseNetworkPolicyFields error = %v, want io.ErrShortBuffer", err)
	}
	fields := &ie.NetworkPolicyFields{}
	if err := fields.MarshalTo(nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want io.ErrShortBuffer", err)
	}
}
