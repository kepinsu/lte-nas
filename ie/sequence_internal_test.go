package ie

import (
	"errors"
	"io"
	"testing"
)

func TestSkipUnknownIE(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		wire []byte
		want int
		err  error
	}{
		{"empty input", nil, 0, io.ErrUnexpectedEOF},
		{"one-octet type1 or type2", []byte{0xff, 0xaa}, 1, nil},
		{"TLV", []byte{0x12, 1, 0xff, 0xaa}, 3, nil},
		{"empty TLV", []byte{0x12, 0, 0xaa}, 2, nil},
		{"TLVE", []byte{0x78, 0, 1, 0xff, 0xaa}, 4, nil},
		{"7C is not comprehension required in Rel18", []byte{0x7c, 0, 0}, 3, nil},
		{"7D is not comprehension required in Rel18", []byte{0x7d, 0, 0}, 3, nil},
		{"7E is comprehension required", []byte{0x7e}, 0, ErrUnknownIE},
		{"7F is comprehension required", []byte{0x7f}, 0, ErrUnknownIE},
		{"0F is comprehension required", []byte{0x0f}, 0, ErrUnknownIE},
		{"truncated TLV length", []byte{0x12}, 0, io.ErrUnexpectedEOF},
		{"truncated TLV value", []byte{0x12, 2, 0}, 0, io.ErrUnexpectedEOF},
		{"truncated TLVE length", []byte{0x78, 0}, 0, io.ErrUnexpectedEOF},
		{"truncated TLVE value", []byte{0x78, 0, 2, 0}, 0, io.ErrUnexpectedEOF},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := skipUnknownIE(tc.wire)
			if got != tc.want || !errors.Is(err, tc.err) {
				t.Fatalf("skipUnknownIE = %d, %v; want %d, %v", got, err, tc.want, tc.err)
			}
		})
	}
}

func TestSkipUnknownIEComprehensionRanges(t *testing.T) {
	t.Parallel()
	for tag := 0; tag <= 255; tag++ {
		n, err := skipUnknownIE([]byte{byte(tag), 0, 0})
		if tag <= 0x0f || tag == 0x7e || tag == 0x7f {
			if !errors.Is(err, ErrUnknownIE) || n != 0 {
				t.Fatalf("comprehension-required tag %02x: %d, %v", tag, n, err)
			}
			continue
		}
		want := 2
		if tag >= 0x80 {
			want = 1
		} else if tag >= 0x78 {
			want = 3
		}
		if err != nil || n != want {
			t.Fatalf("non-comprehension-required tag %02x: %d, %v; want %d", tag, n, err, want)
		}
	}
}
