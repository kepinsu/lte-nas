package snow3g

import (
	"bytes"
	"encoding/hex"
	"strings"
	"testing"
)

// uea2Vector is one official UEA2 specification test set.
type uea2Vector struct {
	name      string
	count     uint32
	bearer    uint32
	direction uint32
	ck        string
	blength   uint32
	plaintext string
	cipher    string
}

// uea2SpecVectors are the five test sets from the UEA2/UIA2 specification
// (version 2.1, sections 4.3 to 4.7).
var uea2SpecVectors = []uea2Vector{
	{
		name:      "spec set 1",
		count:     0x72A4F20F,
		bearer:    0x0c,
		direction: DirDownlink,
		ck:        "2BD6459F82C5B300952C49104881FF48",
		blength:   798,
		plaintext: "7EC61272 743BF161 4726446A 6C38CED1 " +
			"66F6CA76 EB543004 4286346C EF130F92 " +
			"922B0345 0D3A9975 E5BD2EA0 EB55AD8E " +
			"1B199E3E C4316020 E9A1B285 E7627953 " +
			"59B7BDFD 39BEF4B2 484583D5 AFE082AE " +
			"E638BF5F D5A60619 3901A08F 4AB41AAB " +
			"9B134880",
		cipher: "8CEBA629 43DCED3A 0990B06E A1B0A2C4 " +
			"FB3CEDC7 1B369F42 BA64C1EB 6665E72A " +
			"A1C9BB0D EAA20FE8 6058B8BA EE2C2E7F " +
			"0BECCE48 B52932A5 3C9D5F93 1A3A7C53 " +
			"2259AF43 25E2A65E 3084AD5F 6A513B7B " +
			"DDC1B65F 0AA0D97A 053DB55A 88C4C4F9 " +
			"605E4140",
	},
	{
		name:      "spec set 2",
		count:     0xE28BCF7B,
		bearer:    0x18,
		direction: DirUplink,
		ck:        "EFA8B2229E720C2A7C36EA55E9605695",
		blength:   510,
		plaintext: "10111231 E060253A 43FD3F57 E37607AB " +
			"2827B599 B6B1BBDA 37A8ABCC 5A8C550D " +
			"1BFB2F49 4624FB50 367FA36C E3BC68F1 " +
			"1CF93B15 10376B02 130F812A 9FA169D8",
		cipher: "E0DA15CA 8E2554F5 E56C9468 DC6C7C12 " +
			"9C568AA5 032317E0 4E072964 6CABEFA6 " +
			"89864C41 0F24F919 E61E3DFD FAD77E56 " +
			"0DB0A9CD 36C34AE4 181490B2 9F5FA2FC",
	},
	{
		name:      "spec set 3",
		count:     0xFA556B26,
		bearer:    0x03,
		direction: DirDownlink,
		ck:        "5ACB1D644C0D51204EA5F1451010D852",
		blength:   120,
		plaintext: "AD9C441F 890B38C4 57A49D42 1407E8",
		cipher:    "BA0F3130 0334C56B 52A7497C BAC046",
	},
	{
		name:      "spec set 4",
		count:     0x398A59B4,
		bearer:    0x05,
		direction: DirDownlink,
		ck:        "D3C5D592327FB11C4035C6680AF8C6D1",
		blength:   253,
		plaintext: "981BA682 4C1BFB1A B4854720 29B71D80 " +
			"8CE33E2C C3C0B5FC 1F3DE8A6 DC66B1F0",
		cipher: "989B719C DC33CEB7 CF276A52 827CEF94 " +
			"A56C40C0 AB9D81F7 A2A9BAC6 0E11C4B0",
	},
	{
		name:      "spec set 5",
		count:     0x72A4F20F,
		bearer:    0x09,
		direction: DirUplink,
		ck:        "6090EAE04C83706EECBF652BE8E36566",
		blength:   837,
		plaintext: "40981BA6 824C1BFB 4286B299 783DAF44 " +
			"2C099F7A B0F58D5C 8E46B104 F08F01B4 " +
			"1AB48547 2029B71D 36BD1A3D 90DC3A41 " +
			"B46D5167 2AC4C966 3A2BE063 DA4BC8D2 " +
			"808CE33E 2CCCBFC6 34E1B259 060876A0 " +
			"FBB5A437 EBCC8D31 C19E4454 318745E3 " +
			"98764598 7A986F2C B0",
		cipher: "5892BBA8 8BBBCAAE AE769AA0 6B683D3A " +
			"17CC04A3 69881697 435E44FE D5FF9AF5 " +
			"7B9E890D 4D5C6470 9885D48A E40690EC " +
			"043BAAE9 705796E4 A9FF5A4B 8D8B36D7 " +
			"F3FE57CC 6CFD6CD0 05CD3852 A85E94CE " +
			"6BCD90D0 D07839CE 09733544 CA8E3508 " +
			"43248550 922AC128 18",
	},
}

func decodeHex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(strings.ReplaceAll(s, " ", ""))
	if err != nil {
		t.Fatalf("invalid hex %q: %v", s, err)
	}
	return b
}

// TestUEA2Vectors verifies F8 against the official UEA2 specification test
// sets, covering key parsing, IV construction and blength handling.
func TestUEA2Vectors(t *testing.T) {
	for _, tc := range uea2SpecVectors {
		t.Run(tc.name, func(t *testing.T) {
			ck := decodeHex(t, tc.ck)
			plaintext := decodeHex(t, tc.plaintext)
			expected := decodeHex(t, tc.cipher)

			u := NewUEA2(ck, tc.count, tc.bearer, tc.direction)
			got := u.Encrypt(plaintext, tc.blength)

			if !bytes.Equal(got, expected) {
				t.Errorf("F8 mismatch:\ngot  %x\nwant %x", got, expected)
			}
		})
	}
}

// TestUEA2RoundTrip verifies that decrypting an encrypted payload recovers
// the original plaintext for byte-aligned and unaligned bit lengths and for
// both transmission directions. Transmit and receive use separate UEA2
// instances configured with identical parameters, mirroring how a UE and an
// eNodeB each derive the keystream from COUNT^C.
func TestUEA2RoundTrip(t *testing.T) {
	ck := []byte{
		0x2B, 0xD6, 0x45, 0x9F, 0x82, 0xC5, 0xB3, 0x00,
		0x95, 0x2C, 0x49, 0x10, 0x48, 0x81, 0xFF, 0x48,
	}

	tests := []struct {
		name      string
		count     uint32
		bearer    uint32
		direction uint32
		blength   uint32
	}{
		{"byte-aligned downlink", 0x00000001, 0x01, DirDownlink, 64},
		{"byte-aligned small", 0xFA556B26, 0x03, DirDownlink, 120},
		{"unaligned downlink", 0x72A4F20F, 0x0c, DirDownlink, 798},
		{"unaligned uplink", 0xE28BCF7B, 0x18, DirUplink, 510},
		{"single bit", 0x398A59B4, 0x05, DirUplink, 1},
		{"long unaligned", 0x72A4F20F, 0x09, DirUplink, 837},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n := int((tc.blength + 7) >> 3)
			pt := make([]byte, n)
			for i := range pt {
				pt[i] = uint8(0xA5 ^ i)
			}
			// NAS payloads carry padding zeros in unused trailing bits.
			if rem := tc.blength % 8; rem != 0 {
				pt[len(pt)-1] &= Mask8bit(int(rem))
			}

			tx := NewUEA2(ck, tc.count, tc.bearer, tc.direction)
			rx := NewUEA2(ck, tc.count, tc.bearer, tc.direction)
			ct := tx.Encrypt(pt, tc.blength)
			rt := rx.Decrypt(ct, tc.blength)

			if !bytes.Equal(rt, pt) {
				t.Fatalf("round trip mismatch:\nplain      %x\nciphertext %x\nrecovered  %x", pt, ct, rt)
			}
		})
	}
}

// TestUEA2TrailingBitsZeroed verifies that F8 clears the unused trailing
// bits of the final byte when blength is not a multiple of eight.
func TestUEA2TrailingBitsZeroed(t *testing.T) {
	ck := []byte{
		0x2B, 0xD6, 0x45, 0x9F, 0x82, 0xC5, 0xB3, 0x00,
		0x95, 0x2C, 0x49, 0x10, 0x48, 0x81, 0xFF, 0x48,
	}

	tests := []struct {
		name    string
		blength uint32
	}{
		{"two trailing bits", 798},
		{"six trailing bits", 510},
		{"one trailing bit", 253},
		{"five trailing bits", 837},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n := int((tc.blength + 7) >> 3)
			pt := make([]byte, n)
			for i := range pt {
				pt[i] = 0xFF
			}

			u := NewUEA2(ck, 0x00000001, 0x01, DirDownlink)
			ct := u.Encrypt(pt, tc.blength)

			rem := tc.blength % 8
			if got := ct[len(ct)-1] & (0xFF >> rem); got != 0 {
				t.Errorf("last byte 0x%02x: low %d bits = 0x%02x, want 0x00", ct[len(ct)-1], rem, got)
			}
		})
	}
}

// TestUEA2EncryptDecryptAreAliasesOfF8 verifies that Encrypt and Decrypt
// both delegate to F8, which is symmetric under XOR. Fresh instances are
// used because each F8 call advances the keystream position.
func TestUEA2EncryptDecryptAreAliasesOfF8(t *testing.T) {
	ck := []byte{
		0xEF, 0xA8, 0xB2, 0x22, 0x9E, 0x72, 0x0C, 0x2A,
		0x7C, 0x36, 0xEA, 0x55, 0xE9, 0x60, 0x56, 0x95,
	}
	pt := []byte{0x10, 0x11, 0x12, 0x31, 0xE0, 0x60, 0x25, 0x3A}

	enc := NewUEA2(ck, 0xE28BCF7B, 0x18, DirUplink).Encrypt(pt, 64)
	f8 := NewUEA2(ck, 0xE28BCF7B, 0x18, DirUplink).F8(pt, 64)
	dec := NewUEA2(ck, 0xE28BCF7B, 0x18, DirUplink).Decrypt(pt, 64)

	if !bytes.Equal(enc, f8) {
		t.Errorf("Encrypt differs from F8:\nencrypt %x\nf8      %x", enc, f8)
	}
	if !bytes.Equal(dec, f8) {
		t.Errorf("Decrypt differs from F8:\ndecrypt %x\nf8      %x", dec, f8)
	}
}
