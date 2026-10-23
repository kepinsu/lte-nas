package myaes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
)

// EEA2 implements 4G EPS NAS ciphering algorithm (AES-CTR mode).
//
// Counter IV layout per 3GPP TS 33.401 §6.3.2:
//
//	bytes[0..3]   = count (big-endian 32-bit)
//	byte[4]       = (bearer << 3) | (direction << 2)
//	bytes[5..15]  = 0x00
func EEA2(key []byte, count uint32, bearer uint8, direction uint8, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	counter := make([]byte, aes.BlockSize)
	binary.BigEndian.PutUint32(counter, count)
	counter[4] = (bearer << 3) | (direction << 2)

	stream := cipher.NewCTR(block, counter)
	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, plaintext)

	return ciphertext, nil
}

// EIA2 implements 4G EPS NAS integrity algorithm (AES-CMAC mode).
//
// CMAC input per 3GPP TS 33.401 §6.3.3:
//
//	bytes[0..3]   = count (big-endian 32-bit)
//	byte[4]       = (fresh << 3) | (direction << 2)
//	bytes[5..7]   = padding / additional context
//	bytes[8..]    = message
//
// Returns the first 4 bytes of the CMAC output (MIC truncated to 32 bits).
func EIA2(key []byte, count uint32, fresh uint8, direction uint8, msg []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Build CMAC input buffer: header(8 bytes) + message
	macInput := make([]byte, 8+len(msg))
	binary.BigEndian.PutUint32(macInput, count)
	macInput[4] = (fresh << 3) | (direction << 2)
	copy(macInput[8:], msg)

	return cmacSum(macInput, block), nil
}

// ---------------------------------------------------------------------------
// AES-CMAC (RFC 4493 = OMAC1) — pure Go, stdlib-only implementation
// ---------------------------------------------------------------------------

// cmacSum computes the AES-CMAC over data using the provided 128-bit block cipher.
func cmacSum(data []byte, block cipher.Block) []byte {
	const blockSize = 16

	// Step 1: Derive sub-keys K1 and K2 from L = Encrypt(K, 0^blockSize)
	l := make([]byte, blockSize)
	block.Encrypt(l, l)

	k1 := leftShift128(l)
	if msbSet(l) {
		xorWithR128(k1)
	}

	k2 := leftShift128(k1)
	if msbSet(k1) {
		xorWithR128(k2)
	}

	// Step 2: Calculate number of blocks (including padded block if needed)
	n := (len(data) + blockSize - 1) / blockSize // ceil(len/blockSize)
	if n == 0 {
		n = 1 // Empty input → one padded block
	}

	flag := len(data) != 0 && len(data)%blockSize == 0 // true = complete last block, false = needs padding

	// Step 4: Prepare M_last (last block XORed with K1 or padded+K2)
	var M_last = make([]byte, blockSize)
	if flag {
		// Complete block: M_n XOR K1
		copy(M_last, data[(n-1)*blockSize:])
		for j := range M_last {
			M_last[j] ^= k1[j]
		}
	} else {
		// Partial block: pad(remainder) XOR K2
		copy(M_last, data[(n-1)*blockSize:])
		M_last[len(data)%blockSize] = 0x80 // pad with 10^i
		for j := range M_last {
			M_last[j] ^= k2[j]
		}
	}

	// Step 5: Initialize CBC-MAC state X to all zeros
	x := make([]byte, blockSize)

	// Step 6: CBC-MAC on blocks M_1 ... M_{n-1} (the complete/remainder blocks before M_last)
	for i := 0; i < n-1; i++ {
		xorBlock(x, data[i*blockSize:(i+1)*blockSize])
		block.Encrypt(x, x)
	}

	// Step 6 continued: Y = M_last XOR X; T = AES(Y)
	xorBlock(x, M_last)
	block.Encrypt(x, x)

	return x[:4]
}

// leftShift128 performs OMAC1-style left-shift on a 16-byte array (GF(2^128) × x).
// Carry propagates from the most significant position to the least significant.
func leftShift128(src []byte) (dst []byte) {
	dst = make([]byte, len(src))
	carry := byte(0)
	for i := len(src) - 1; i >= 0; i-- {
		dst[i] = (src[i] << 1) | carry
		carry = src[i] >> 7
	}
	return dst
}

// msbSet returns true if the most-significant bit of src is set.
func msbSet(src []byte) bool { return src[0]&0x80 != 0 }

// xorWithR128 XORs the last byte of dst with the MSB constant of R128 (0x87).
func xorWithR128(dst []byte) { dst[len(dst)-1] ^= 0x87 }

// xorBlock XORs src into dst in-place (both must be blockSize bytes).
func xorBlock(dst, src []byte) {
	for i := range dst {
		dst[i] ^= src[i]
	}
}
