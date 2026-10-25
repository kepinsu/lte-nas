// UEA2 implements the EEA1 confidentiality algorithm (F8) over SNOW 3G,
// as defined by 3GPP TS 33.401.
//
// Adapted from the ETSI / SAGE specification of the 3GPP Confidentiality and
// Integrity Algorithms UEA2 & UIA2.
// Document 1: UEA2 and UIA2 Specification, Version 2.1 from the 16th March
// 2009, annex 4 (https://www.gsma.com/security/wp-content/uploads/2019/05/uea2uia2d1v21.pdf).

package snow3g

// UEA2 is a stream cipherer wrapping a Snow3G instance configured with the
// EEA1 key and IV derivation. The keystream position advances on every F8
// call, so one instance should be used per PDU (per COUNT^C value).
type UEA2 struct {
	snow3g *Snow3G
}

// NewUEA2 creates a UEA2 cipherer from the 128-bit ciphering key ck, the
// 32-bit COUNT^C value, the radio bearer identity and the transmission
// direction (DirUplink or DirDownlink).
func NewUEA2(ck []uint8, count uint32, bearer uint32, direction uint32) *UEA2 {
	k := [4]uint32{}
	for i := 0; i < 4; i++ {
		k[3-i] = (uint32(ck[4*i]) << 24) ^ (uint32(ck[4*i+1]) << 16) ^
			(uint32(ck[4*i+2]) << 8) ^ uint32(ck[4*i+3])
	}

	iv := [4]uint32{}
	iv[3] = count
	iv[2] = (bearer << 27) | ((direction & 0x1) << 26)
	iv[1] = iv[3]
	iv[0] = iv[2]

	s3g := &Snow3G{}
	s3g.Initialize(k, iv)

	return &UEA2{snow3g: s3g}
}

// Encrypt applies F8 to data and returns the ciphertext. It is an alias of
// F8 because XOR with the keystream is symmetric.
func (u *UEA2) Encrypt(data []byte, blength uint32) []byte {
	return u.F8(data, blength)
}

// Decrypt applies F8 to data and returns the plaintext. It is an alias of
// F8 because XOR with the keystream is symmetric.
func (u *UEA2) Decrypt(data []byte, blength uint32) []byte {
	return u.F8(data, blength)
}

// F8 applies the EEA1 ciphering function: it XORs the first blength bits of
// data with the SNOW 3G keystream derived from the configured key and IV.
// Unused trailing bits in the final byte are zeroed and the result is
// truncated to ceil(blength/8) bytes.
func (u *UEA2) F8(data []byte, blength uint32) []byte {
	zeroBits := blength & 0x7
	length := blength >> 3
	if zeroBits > 0 {
		length += 1
	}
	n := int((blength + 31) >> 5)

	ks := u.snow3g.GenerateKeystream(n)
	output := make([]byte, length)

	for i := 0; i < n; i++ {
		for j := 0; j < 4 && i*4+j < int(length); j++ {
			output[4*i+j] = data[4*i+j] ^ uint8((ks[i]>>(8*(3-j)))&0xff)
		}
	}

	// Zero the unused trailing bits of the final byte.
	if zeroBits > 0 {
		output[length-1] &= uint8(0xff) << (8 - zeroBits)
	}

	return output
}
