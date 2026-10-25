package myaes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"testing"
)

func TestEEA2_RoundTrip(t *testing.T) {
	key := make([]byte, aes.BlockSize) // 128-bit key (all zeros for deterministic testing)
	count := uint32(1)
	bearer := uint8(5)
	direction := uint8(0)

	plaintext := []byte("Hello, LTE NAS world! This is a ciphering test message.")

	ciphertext, err := EEA2(key, count, bearer, direction, plaintext)
	if err != nil {
		t.Fatalf("EEA2 encrypt failed: %v", err)
	}

	if len(ciphertext) == 0 {
		t.Fatal("EEA2 produced empty ciphertext")
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Fatal("EEA2 ciphertext equals plaintext — XOR not applied")
	}

	decrypted, err := EEA2(key, count, bearer, direction, ciphertext)
	if err != nil {
		t.Fatalf("EEA2 decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("round-trip failed: got %q, want %q", decrypted, plaintext)
	}
}

func TestEEA2_EmptyPlaintext(t *testing.T) {
	key := make([]byte, aes.BlockSize)
	result, err := EEA2(key, 0, 5, 0, []byte{})
	if err != nil {
		t.Fatalf("EEA2 empty plaintext failed: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty output, got %d bytes", len(result))
	}
}

func TestEEA2_DifferentCounts(t *testing.T) {
	key := []byte("0123456789ABCDEF") // exactly 16 bytes
	msg := []byte("ABCD")

	ct1, err := EEA2(key, 100, 5, 0, msg)
	if err != nil {
		t.Fatalf("EEA2 failed: %v", err)
	}

	ct2, err := EEA2(key, 101, 5, 0, msg)
	if err != nil {
		t.Fatalf("EEA2 failed: %v", err)
	}

	if bytes.Equal(ct1, ct2) {
		t.Fatal("different counts should produce different ciphertext")
	}
}

func TestEEA2_DifferentBearers(t *testing.T) {
	key := []byte("0123456789ABCDEF")
	msg := []byte("ABCD")

	ct1, err := EEA2(key, 1, 5, 0, msg)
	if err != nil {
		t.Fatalf("EEA2 failed: %v", err)
	}

	ct2, err := EEA2(key, 1, 6, 0, msg)
	if err != nil {
		t.Fatalf("EEA2 failed: %v", err)
	}

	if bytes.Equal(ct1, ct2) {
		t.Fatal("different bearers should produce different ciphertext")
	}
}

func TestEEA2_DifferentDirections(t *testing.T) {
	key := []byte("0123456789ABCDEF")
	msg := []byte("ABCD")

	ct0, err := EEA2(key, 1, 5, 0, msg)
	if err != nil {
		t.Fatalf("EEA2 failed: %v", err)
	}

	ct1, err := EEA2(key, 1, 5, 1, msg)
	if err != nil {
		t.Fatalf("EEA2 failed: %v", err)
	}

	if bytes.Equal(ct0, ct1) {
		t.Fatal("different directions should produce different ciphertext")
	}
}

func TestEEA2_MultipleBlocks(t *testing.T) {
	key := []byte("0123456789ABCDEF")
	// 32 bytes — exactly 2 AES blocks
	plaintext := make([]byte, 32)
	for i := range plaintext {
		plaintext[i] = byte(i)
	}

	ciphertext, err := EEA2(key, 42, 0, 1, plaintext)
	if err != nil {
		t.Fatalf("EEA2 failed: %v", err)
	}

	decrypted, err := EEA2(key, 42, 0, 1, ciphertext)
	if err != nil {
		t.Fatalf("EEA2 decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("multi-block round-trip failed")
	}
}

func TestEIA2_OutputSize(t *testing.T) {
	key := []byte("0123456789ABCDEF")
	fresh := uint8(0)
	msg := []byte("test message for integrity check")

	mac, err := EIA2(key, 1, fresh, 0, msg)
	if err != nil {
		t.Fatalf("EIA2 failed: %v", err)
	}

	if len(mac) != 4 {
		t.Fatalf("expected MAC length 4, got %d", len(mac))
	}
}

func TestEIA2_DifferentMessages(t *testing.T) {
	key := []byte("0123456789ABCDEF")

	mac1, err := EIA2(key, 1, 0, 0, []byte("message A"))
	if err != nil {
		t.Fatalf("EIA2 failed: %v", err)
	}

	mac2, err := EIA2(key, 1, 0, 0, []byte("message B"))
	if err != nil {
		t.Fatalf("EIA2 failed: %v", err)
	}

	if bytes.Equal(mac1, mac2) {
		t.Fatal("different messages should produce different MACs")
	}
}

func TestEIA2_DifferentCounts(t *testing.T) {
	key := []byte("0123456789ABCDEF")

	mac1, err := EIA2(key, 1, 0, 0, []byte("same message"))
	if err != nil {
		t.Fatalf("EIA2 failed: %v", err)
	}

	mac2, err := EIA2(key, 2, 0, 0, []byte("same message"))
	if err != nil {
		t.Fatalf("EIA2 failed: %v", err)
	}

	if bytes.Equal(mac1, mac2) {
		t.Fatal("different counts should produce different MACs")
	}
}

func TestEIA2_DifferentDirections(t *testing.T) {
	key := []byte("0123456789ABCDEF")

	mac0, err := EIA2(key, 1, 0, 0, []byte("direction test"))
	if err != nil {
		t.Fatalf("EIA2 (dir=0) failed: %v", err)
	}

	mac1, err := EIA2(key, 1, 0, 1, []byte("direction test"))
	if err != nil {
		t.Fatalf("EIA2 (dir=1) failed: %v", err)
	}

	if bytes.Equal(mac0, mac1) {
		t.Fatal("different directions should produce different MACs")
	}
}

func TestEIA2_DifferentFresh(t *testing.T) {
	key := []byte("0123456789ABCDEF")

	mac0, err := EIA2(key, 1, 0, 0, []byte("fresh test"))
	if err != nil {
		t.Fatalf("EIA2 (fresh=0) failed: %v", err)
	}

	mac1, err := EIA2(key, 1, 1, 0, []byte("fresh test"))
	if err != nil {
		t.Fatalf("EIA2 (fresh=1) failed: %v", err)
	}

	if bytes.Equal(mac0, mac1) {
		t.Fatal("different fresh values should produce different MACs")
	}
}

func TestEIA2_EmptyMessage(t *testing.T) {
	key := []byte("0123456789ABCDEF")

	mac, err := EIA2(key, 1, 0, 0, []byte{})
	if err != nil {
		t.Fatalf("EIA2 empty message failed: %v", err)
	}

	if len(mac) != 4 {
		t.Fatalf("expected MAC length 4, got %d", len(mac))
	}

	if bytes.Equal(mac, []byte{0, 0, 0, 0}) {
		t.Fatal("MAC is all zeros — possible bug in CMAC")
	}
}

func TestEIA2_MultipleBlocks(t *testing.T) {
	key := []byte("0123456789ABCDEF")
	// 24 bytes — exceeds one block, requires multi-block CMAC
	msg := make([]byte, 24)
	for i := range msg {
		msg[i] = byte(i + 1)
	}

	mac, err := EIA2(key, 10, 5, 0, msg)
	if err != nil {
		t.Fatalf("EIA2 failed: %v", err)
	}

	if len(mac) != 4 {
		t.Fatalf("expected MAC length 4, got %d", len(mac))
	}
}

func TestCMAC_SameKeySameInput(t *testing.T) {
	key := []byte("0123456789ABCDEF")
	msg := []byte("determinism check for CMAC")

	mac1, err := EIA2(key, 7, 0, 0, msg)
	if err != nil {
		t.Fatalf("EIA2 failed: %v", err)
	}

	mac2, err := EIA2(key, 7, 0, 0, msg)
	if err != nil {
		t.Fatalf("EIA2 failed: %v", err)
	}

	if !bytes.Equal(mac1, mac2) {
		t.Fatal("deterministic inputs should produce same MAC")
	}
}

// RFC 4493 Section 4 test vectors for AES-CMAC (raw, without headers)
func TestCMAC_RFC4493_Empty(t *testing.T) {
	key := []byte{
		0x2b, 0x7e, 0x15, 0x16, 0x28, 0xae, 0xd2, 0xa6,
		0xab, 0xf7, 0x15, 0x88, 0x09, 0xcf, 0x4f, 0x3c,
	}
	msg := []byte{}

	// RFC 4493: bb1d6929 e9593728 7fa37d12 9b756746
	expectedMAC := []byte{0xbb, 0x1d, 0x69, 0x29}

	block := mustNewCipher(key)
	mac := cmacSum(msg, block)
	if len(mac) != 4 {
		t.Fatalf("expected MAC length 4, got %d", len(mac))
	}
	if !bytes.Equal(mac, expectedMAC) {
		t.Errorf("CMAC(empty): got %x, want %x", mac, expectedMAC)
	}
}

func TestCMAC_RFC4493_OneBlock(t *testing.T) {
	key := []byte{
		0x2b, 0x7e, 0x15, 0x16, 0x28, 0xae, 0xd2, 0xa6,
		0xab, 0xf7, 0x15, 0x88, 0x09, 0xcf, 0x4f, 0x3c,
	}
	msg := []byte{
		0x6b, 0xc1, 0xbe, 0xe2, 0x2e, 0x40, 0x9f, 0x96,
		0xe9, 0x3d, 0x7e, 0x11, 0x73, 0x93, 0x17, 0x2a,
	}

	// RFC 4493: 070a16b4 6b4d4144 f79bdd9d d04a287c
	expectedMAC := []byte{0x07, 0x0a, 0x16, 0xb4}

	block := mustNewCipher(key)
	mac := cmacSum(msg, block)
	if len(mac) != 4 {
		t.Fatalf("expected MAC length 4, got %d", len(mac))
	}
	if !bytes.Equal(mac, expectedMAC) {
		t.Errorf("CMAC(1-block): got %x, want %x", mac, expectedMAC)
	}
}

func TestCMAC_RFC4493_TwoBlock(t *testing.T) {
	key := []byte{
		0x2b, 0x7e, 0x15, 0x16, 0x28, 0xae, 0xd2, 0xa6,
		0xab, 0xf7, 0x15, 0x88, 0x09, 0xcf, 0x4f, 0x3c,
	}
	msg := []byte{
		0x6b, 0xc1, 0xbe, 0xe2, 0x2e, 0x40, 0x9f, 0x96,
		0xe9, 0x3d, 0x7e, 0x11, 0x73, 0x93, 0x17, 0x2a,
		0xae, 0x2d, 0x8a, 0x57, 0x1e, 0x03, 0xac, 0x9c,
		0x9e, 0xb7, 0x6f, 0xac, 0x45, 0xaf, 0x8e, 0x51,
		0x30, 0xc8, 0x1c, 0x46, 0xa3, 0x5c, 0xe4, 0x11,
	}

	// RFC 4493: dfa66747 de9ae630 30ca3261 1497c827
	expectedMAC := []byte{0xdf, 0xa6, 0x67, 0x47}

	block := mustNewCipher(key)
	mac := cmacSum(msg, block)
	if len(mac) != 4 {
		t.Fatalf("expected MAC length 4, got %d", len(mac))
	}
	if !bytes.Equal(mac, expectedMAC) {
		t.Errorf("CMAC(2-blocks): got %x, want %x", mac, expectedMAC)
	}
}

func TestCMAC_RFC4493_FourBlock(t *testing.T) {
	key := []byte{
		0x2b, 0x7e, 0x15, 0x16, 0x28, 0xae, 0xd2, 0xa6,
		0xab, 0xf7, 0x15, 0x88, 0x09, 0xcf, 0x4f, 0x3c,
	}
	msg := []byte{
		0x6b, 0xc1, 0xbe, 0xe2, 0x2e, 0x40, 0x9f, 0x96,
		0xe9, 0x3d, 0x7e, 0x11, 0x73, 0x93, 0x17, 0x2a,
		0xae, 0x2d, 0x8a, 0x57, 0x1e, 0x03, 0xac, 0x9c,
		0x9e, 0xb7, 0x6f, 0xac, 0x45, 0xaf, 0x8e, 0x51,
		0x30, 0xc8, 0x1c, 0x46, 0xa3, 0x5c, 0xe4, 0x11,
		0xe5, 0xfb, 0xc1, 0x19, 0x1a, 0x0a, 0x52, 0xef,
		0xf6, 0x9f, 0x24, 0x45, 0xdf, 0x4f, 0x9b, 0x17,
		0xad, 0x2b, 0x41, 0x7b, 0xe6, 0x6c, 0x37, 0x10,
	}

	// RFC 4493: 51f0bebf 7e3b9d92 fc497417 79363cfe
	expectedMAC := []byte{0x51, 0xf0, 0xbe, 0xbf}

	block := mustNewCipher(key)
	mac := cmacSum(msg, block)
	if len(mac) != 4 {
		t.Fatalf("expected MAC length 4, got %d", len(mac))
	}
	if !bytes.Equal(mac, expectedMAC) {
		t.Errorf("CMAC(4-blocks): got %x, want %x", mac, expectedMAC)
	}
}

func mustNewCipher(key []byte) cipher.Block {
	b, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return b
}
