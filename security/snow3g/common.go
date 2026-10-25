// Snow3G field arithmetic helpers.
//
// These functions implement Galois Field GF(2^16) multiply using the AES
// primitive polynomial x^16 + x^5 + x^3 + x^2 + 1 (0x1B).
// Adapted from the ETSI/SAGE UEA2/UIA2 specification.

package snow3g

// mul64x doubles a value in GF(2^16) with reduction by polynomial c.
// If the high bit is set, XOR the shifted result with c.
func mul64x(v uint64, c uint64) uint64 {
	if v&0x8000000000000000 != 0 {
		return (v << 1) ^ c
	}
	return v << 1
}

// mul64xpow applies mul64x i times (i.e. multiplies v by 2^i in GF).
func mul64xpow(v uint64, i uint8, c uint64) uint64 {
	for i > 0 {
		v = mul64x(v, c)
		i--
	}
	return v
}

// Mul64 performs Galois Field multiplication of v and p using polynomial c.
// Each set bit in p selects the corresponding 2^k-multiple of v, XORed together.
// c is the AES Galois polynomial 0x1B for Snow3G field operations.
func Mul64(v uint64, p uint64, c uint64) uint64 {
	r := uint64(0)
	for i := 0; i < 64; i++ {
		if (p>>i)&0x1 != 0 {
			r ^= mul64xpow(v, uint8(i), c)
		}
	}
	return r
}

// Mask8bit returns a uint8 with the upper n bits set to 1 and all lower bits
// cleared. Panics if n < 0 or n > 8.
func Mask8bit(n int) uint8 {
	if n < 0 || n > 8 {
		panic("snow3g: Mask8bit n out of range [0, 8]")
	}
	return uint8(0xFF ^ ((1 << (8 - n)) - 1))
}
