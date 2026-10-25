package snow3g

import (
	"testing"
)

func TestMul64(t *testing.T) {
	const poly = 0x1B

	tests := []struct {
		name   string
		v, p   uint64
		expect uint64
	}{
		{"identity low", 0x1234, 1, 0x1234},
		{"identity high", 0xDEADBEEFCAFEBABE, 1, 0xDEADBEEFCAFEBABE},
		{"identity one", 1, 1, 1},
		{"zero v", 0, 0xABCD, 0},
		{"zero p", 0xABCDEF1234567890, 0, 0},
		{"both zero", 0, 0, 0},

		// mul64(x, 2) doubles x in GF with reduction by poly when needed
		{"double no overflow", 0x10, 2, 0x20},

		// When v has bit 63 set, doubling causes reduction:
		// mul64x(0x8000000000000000, 0x1B) = ((v<<1)^0x1B) overflow -> 0 ^ 0x1B = 0x1B
		{"double high-bit v", 0x8000000000000000, 2, 0x1B},

		// Polynomial multiplication: Mul64(x, x, poly) should produce x*x in the field
		{"p=x with v=1", 0x8000000000000000, 1, 0x8000000000000000},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Mul64(tc.v, tc.p, poly)
			if got != tc.expect {
				t.Errorf("Mul64(0x%X, 0x%X, poly) = 0x%X, want 0x%X",
					tc.v, tc.p, got, tc.expect)
			}
		})
	}
}

func TestMul64_PureFunction(t *testing.T) {
	const poly = 0x1B
	// Call the same inputs multiple times — result must be identical.
	got1 := Mul64(0xABCDEF1234567890, 0xFEDCBA0987654321, poly)
	got2 := Mul64(0xABCDEF1234567890, 0xFEDCBA0987654321, poly)
	got3 := Mul64(0xABCDEF1234567890, 0xFEDCBA0987654321, poly)

	if got1 != got2 || got2 != got3 {
		t.Errorf("Mul64 is not deterministic: 0x%X, 0x%X, 0x%X", got1, got2, got3)
	}
}

func TestMask8bit(t *testing.T) {
	tests := []struct {
		n      int
		expect uint8
	}{
		{0, 0x00}, // no bits set (upper 0)
		{1, 0x80}, // upper 1 bit
		{4, 0xF0}, // upper 4 bits
		{7, 0xFE}, // upper 7 bits
		{8, 0xFF}, // all 8 bits
	}

	for _, tc := range tests {
		t.Run(string(rune(tc.n)), func(t *testing.T) {
			got := Mask8bit(tc.n)
			if got != tc.expect {
				t.Errorf("Mask8bit(%d) = 0x%02X, want 0x%02X", tc.n, got, tc.expect)
			}
		})
	}

	t.Run("panic negative", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Mask8bit(-1) did not panic")
			}
		}()
		Mask8bit(-1)
	})

	t.Run("panic >8", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Mask8bit(9) did not panic")
			}
		}()
		Mask8bit(9)
	})
}

func TestMul64_xPowConsistency(t *testing.T) {
	// mul64xpow(v, i, c) == Mul64(v, 1<<i, c) for single-bit p
	const poly = 0x1B

	cases := []uint64{
		0, 1, 2, 7, 0xFF, 0x8000000000000000,
	}

	for _, v := range cases {
		for i := uint8(0); i < 64; i++ {
			p := uint64(1) << i
			gotPow := mul64xpow(v, i, poly)
			gotMul := Mul64(v, p, poly)
			if gotPow != gotMul {
				t.Errorf("mul64xpow(0x%X, %d) != Mul64(0x%X, 1<<%d): gotPow=0x%X, gotMul=0x%X",
					v, i, v, i, gotPow, gotMul)
			}
		}
	}
}
