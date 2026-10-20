package ie_test

import (
	"errors"
	"testing"

	"lte-nas/ie"
)

func TestIE_ValueAsUint8(t *testing.T) {
	t.Parallel()
	i := ie.NewPositionalPDNType(3)
	if got, err := i.ValueAsUint8(); err != nil || got != 3 {
		t.Fatalf("ValueAsUint8 = %d, %v", got, err)
	}
	for _, tc := range []struct {
		name string
		ie   *ie.IE
		err  error
	}{
		{"nil", nil, ie.ErrNilIE},
		{"empty", ie.NewRaw(ie.NoTag, ie.IEFormatOne, nil), ie.ErrInvalidLength},
		{"too long", ie.NewRaw(ie.NoTag, ie.IEFormatThree, []byte{1, 2}), ie.ErrInvalidLength},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if _, err := tc.ie.ValueAsUint8(); !errors.Is(err, tc.err) {
				t.Fatalf("ValueAsUint8 = %v", err)
			}
		})
	}
}

func TestIE_ValueAsUint16(t *testing.T) {
	t.Parallel()
	i := ie.NewPositionalShortMAC([]byte{0x12, 0x34})
	if got, err := i.ValueAsUint16(); err != nil || got != 0x1234 {
		t.Fatalf("ValueAsUint16 = %d, %v", got, err)
	}
	if _, err := ie.NewPositionalShortMAC([]byte{1}).ValueAsUint16(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("short ValueAsUint16 = %v", err)
	}
}

func TestIE_ValueAsUint32(t *testing.T) {
	t.Parallel()
	i := ie.NewNonceUE([]byte{0x12, 0x34, 0x56, 0x78})
	if got, err := i.ValueAsUint32(); err != nil || got != 0x12345678 {
		t.Fatalf("ValueAsUint32 = %d, %v", got, err)
	}
	if _, err := ie.NewNonceUE([]byte{1}).ValueAsUint32(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("short ValueAsUint32 = %v", err)
	}
}
