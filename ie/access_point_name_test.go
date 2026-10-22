package ie_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"lte-nas/ie"
)

func TestAccessPointNameFields(t *testing.T) {
	t.Parallel()
	fields := &ie.AccessPointNameFields{Name: "internet.example"}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{8, 'i', 'n', 't', 'e', 'r', 'n', 'e', 't', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e'}
	if !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}

	value, err := ie.NewAccessPointNameWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.AccessPointName()
	if err != nil {
		t.Fatal(err)
	}
	if decoded.Name != fields.Name {
		t.Fatalf("AccessPointName = %q, want %q", decoded.Name, fields.Name)
	}
	if fields.MarshalLen() != len(want) {
		t.Fatalf("MarshalLen = %d, want %d", fields.MarshalLen(), len(want))
	}
}

func TestAccessPointNameFields_RejectInvalid(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"", ".internet", "internet.", "internet..example"} {
		if _, err := (&ie.AccessPointNameFields{Name: name}).Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
			t.Fatalf("Marshal(%q) = %v, want ErrInvalidValue", name, err)
		}
	}
	if _, err := ie.ParseAccessPointNameFields([]byte{8, 'i'}); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("ParseAccessPointNameFields = %v, want ErrInvalidValue", err)
	}

	fields := &ie.AccessPointNameFields{Name: "internet"}
	buffer := []byte{0xaa}
	if err := fields.MarshalTo(buffer); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo = %v, want io.ErrShortBuffer", err)
	}
	if buffer[0] != 0xaa {
		t.Fatalf("MarshalTo modified a short buffer: got %x, want aa", buffer)
	}
}
