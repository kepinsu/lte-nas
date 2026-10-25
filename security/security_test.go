package security

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math"
	"strings"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

// These tests exercise the EPS NAS security API of package security against
// specification vectors:
//
//   - K_ASME -> CK/IK derivation: TS 33.401 clause 5.1 formula, evaluated by
//     an independent inline oracle (no published Annex A key-derivation
//     vector bytes are available locally).
//   - EEA1 (UEA2/F8) and EIA1 (UIA2/F9): GSMA UEA2/UIA2 specification test
//     sets, the same vectors used by security/snow3g/*_test.go. Only the
//     byte-aligned sets are asserted here because the context API passes
//     len(data)*8 bits per NAS usage; the non-aligned primitive vectors stay
//     covered at the snow3g package level.
//   - EEA2 and EIA2: TS 33.401 annex A test sets (AES-CTR / AES-CMAC), whose
//     byte-wise construction matches myaes.EEA2/myaes.EIA2 exactly.
//
// Vector values were copied read-only from known-good external test suites;
// nothing from those modules is imported.

const (
	kasmeAnchorA = "000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F"
	kasmeAnchorB = "FFEECCBBAA9988776655443322110000112233445566778899AABBCCDDEEFF00"
)

// hexBytes decodes a hex string, ignoring spaces between byte pairs.
func hexBytes(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(strings.ReplaceAll(s, " ", ""))
	if err != nil {
		t.Fatalf("invalid hex %q: %v", s, err)
	}
	return b
}

// ts33401DerivedKey evaluates the TS 33.401 clause 5.1 key derivation
// independently of deriveKey: HMAC-SHA256 keyed with K_ASME over
// algType || 0x0010 || algoID || 0x01, truncated to the first 16 octets.
// algType is 0x03 for CK (NAS encryption) and 0x02 for IK (NAS integrity).
func ts33401DerivedKey(kasme []byte, algType, algoID byte) []byte {
	mac := hmac.New(sha256.New, kasme)
	mac.Write([]byte{algType, 0x00, 0x10, algoID, 0x01})
	return mac.Sum(nil)[:16]
}

// TestDeriveCK verifies CK derivation from K_ASME for the supported
// ciphering algorithm identifiers against the TS 33.401 clause 5.1 oracle.
func TestDeriveCK(t *testing.T) {
	tests := []struct {
		name   string
		kasme  string
		algoID byte
	}{
		{"EEA0", kasmeAnchorA, 0x00},
		{"EEA1", kasmeAnchorA, 0x01},
		{"EEA2", kasmeAnchorA, 0x02},
		{"EEA1 other anchor", kasmeAnchorB, 0x01},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			kasme := hexBytes(t, tc.kasme)
			got := DeriveCK(kasme, tc.algoID)
			if len(got) != 16 {
				t.Fatalf("CK length = %d, want 16", len(got))
			}
			want := ts33401DerivedKey(kasme, 0x03, tc.algoID)
			if !bytes.Equal(got, want) {
				t.Errorf("CK mismatch:\ngot  %x\nwant %x", got, want)
			}
		})
	}
}

// TestDeriveIK verifies IK derivation from K_ASME for the supported
// integrity algorithm identifiers against the TS 33.401 clause 5.1 oracle.
func TestDeriveIK(t *testing.T) {
	tests := []struct {
		name   string
		kasme  string
		algoID byte
	}{
		{"EIA0", kasmeAnchorA, 0x00},
		{"EIA1", kasmeAnchorA, 0x01},
		{"EIA2", kasmeAnchorA, 0x02},
		{"EIA1 other anchor", kasmeAnchorB, 0x01},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			kasme := hexBytes(t, tc.kasme)
			got := DeriveIK(kasme, tc.algoID)
			if len(got) != 16 {
				t.Fatalf("IK length = %d, want 16", len(got))
			}
			want := ts33401DerivedKey(kasme, 0x02, tc.algoID)
			if !bytes.Equal(got, want) {
				t.Errorf("IK mismatch:\ngot  %x\nwant %x", got, want)
			}
		})
	}
}

// TestDerivedKeysDistinct verifies that the derived keys are pairwise
// different across algTypes and algorithm identifiers and that derivation is
// deterministic.
func TestDerivedKeysDistinct(t *testing.T) {
	a := hexBytes(t, kasmeAnchorA)
	keys := map[string][]byte{
		"CK(EEA1)": DeriveCK(a, 0x01),
		"CK(EEA2)": DeriveCK(a, 0x02),
		"IK(EIA1)": DeriveIK(a, 0x01),
		"IK(EIA2)": DeriveIK(a, 0x02),
	}
	names := []string{"CK(EEA1)", "CK(EEA2)", "IK(EIA1)", "IK(EIA2)"}
	for i := 0; i < len(names); i++ {
		for j := i + 1; j < len(names); j++ {
			if bytes.Equal(keys[names[i]], keys[names[j]]) {
				t.Errorf("%s == %s, want distinct derived keys", names[i], names[j])
			}
		}
	}
	if !bytes.Equal(DeriveCK(a, 0x01), keys["CK(EEA1)"]) {
		t.Error("DeriveCK is not deterministic")
	}
}

// TestSetupContext verifies the full EMMSecurityContext setup flow: valid
// combinations produce a complete, inactive context with correctly derived
// keys, and invalid inputs are rejected with the documented errors.
func TestSetupContext(t *testing.T) {
	success := []struct {
		name    string
		kasme   string
		ksi     uint8
		encAlgo ie.NASEncryptionAlgorithm
		intAlgo ie.NASIntegrityAlgorithm
	}{
		{"EEA1 EIA1 ksi 0", kasmeAnchorA, 0, EEA1, EIA1},
		{"EEA2 EIA2 ksi 7", kasmeAnchorA, 7, EEA2, EIA2},
		{"EEA0 EIA0 ksi 3", kasmeAnchorB, 3, EEA0, EIA0},
	}
	for _, tc := range success {
		t.Run(tc.name, func(t *testing.T) {
			kasme := hexBytes(t, tc.kasme)
			ctx, err := SetupContext(kasme, tc.ksi, tc.encAlgo, tc.intAlgo)
			if err != nil {
				t.Fatalf("SetupContext returned error: %v", err)
			}
			if !ctx.Valid || ctx.Active {
				t.Errorf("Valid=%v Active=%v, want Valid=true Active=false", ctx.Valid, ctx.Active)
			}
			if !bytes.Equal(ctx.Kasme, kasme) {
				t.Errorf("Kasme = %x, want %x", ctx.Kasme, kasme)
			}
			if !bytes.Equal(ctx.CK, DeriveCK(kasme, byte(tc.encAlgo))) {
				t.Errorf("CK = %x, want %x", ctx.CK, DeriveCK(kasme, byte(tc.encAlgo)))
			}
			if !bytes.Equal(ctx.IK, DeriveIK(kasme, byte(tc.intAlgo))) {
				t.Errorf("IK = %x, want %x", ctx.IK, DeriveIK(kasme, byte(tc.intAlgo)))
			}
			if ctx.EncAlgo != tc.encAlgo || ctx.IntAlgo != tc.intAlgo {
				t.Errorf("algorithms = (%d,%d), want (%d,%d)", ctx.EncAlgo, ctx.IntAlgo, tc.encAlgo, tc.intAlgo)
			}
			if ctx.UplinkNASCount != 0 || ctx.DownlinkNASCount != 0 {
				t.Errorf("counts = (%d,%d), want (0,0)", ctx.UplinkNASCount, ctx.DownlinkNASCount)
			}
		})
	}

	failure := []struct {
		name     string
		kasmeLen int
		ksi      uint8
		encAlgo  ie.NASEncryptionAlgorithm
		intAlgo  ie.NASIntegrityAlgorithm
		wantErr  error
	}{
		{"empty K_ASME", 0, 0, EEA1, EIA1, ErrInvalidKASME},
		{"short K_ASME", 16, 0, EEA1, EIA1, ErrInvalidKASME},
		{"long K_ASME", 33, 0, EEA1, EIA1, ErrInvalidKASME},
		{"KSI above 3 bits", 32, 8, EEA1, EIA1, ErrInvalidEKSI},
		{"KSI full byte", 32, 0xFF, EEA1, EIA1, ErrInvalidEKSI},
		{"unsupported ciphering algorithm", 32, 0, ie.NASEncryptionAlgorithm128EEA3, EIA1, ErrUnsupportedAlgo},
		{"unsupported integrity algorithm", 32, 0, EEA1, ie.NASIntegrityAlgorithm128EIA3, ErrUnsupportedAlgo},
	}
	for _, tc := range failure {
		t.Run(tc.name, func(t *testing.T) {
			var kasme []byte
			if tc.kasmeLen > 0 {
				kasme = make([]byte, tc.kasmeLen)
				for i := range kasme {
					kasme[i] = byte(i)
				}
			}
			ctx, err := SetupContext(kasme, tc.ksi, tc.encAlgo, tc.intAlgo)
			if ctx != nil {
				t.Fatalf("SetupContext returned context %v, want nil", ctx)
			}
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("error = %v, want %v", err, tc.wantErr)
			}
		})
	}
}

// TestSetupContextCopiesKASME verifies that the context stores its own copy
// of the anchor key instead of aliasing the caller's slice.
func TestSetupContextCopiesKASME(t *testing.T) {
	kasme := hexBytes(t, kasmeAnchorA)
	first := kasme[0]
	ctx, err := SetupContext(kasme, 0, EEA1, EIA1)
	if err != nil {
		t.Fatalf("SetupContext returned error: %v", err)
	}
	kasme[0] ^= 0xFF
	if ctx.Kasme[0] != first {
		t.Errorf("Kasme[0] = 0x%02x after mutating the input, want 0x%02x", ctx.Kasme[0], first)
	}
}

// TestEEA0Identity verifies that the null ciphering algorithm returns the
// payload unchanged in both directions and rejects nil payloads.
func TestEEA0Identity(t *testing.T) {
	ctx, err := SetupContext(hexBytes(t, kasmeAnchorA), 0, EEA0, EIA0)
	if err != nil {
		t.Fatalf("SetupContext returned error: %v", err)
	}

	tests := []struct {
		name   string
		count  uint32
		bearer uint8
		size   int
	}{
		{"empty payload", 0, 0, 0},
		{"single byte", 1, 5, 1},
		{"block aligned", 0x1234, 0x0C, 16},
		{"odd size", 0x72A4F20F, 0x18, 33},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pt := make([]byte, tc.size)
			for i := range pt {
				pt[i] = uint8(0xA5 ^ i)
			}
			ct, err := ctx.Encrypt(tc.count, tc.bearer, pt)
			if err != nil {
				t.Fatalf("Encrypt returned error: %v", err)
			}
			if !bytes.Equal(ct, pt) {
				t.Errorf("Encrypt altered the payload:\ngot  %x\nwant %x", ct, pt)
			}
			rt, err := ctx.Decrypt(tc.count, tc.bearer, ct)
			if err != nil {
				t.Fatalf("Decrypt returned error: %v", err)
			}
			if !bytes.Equal(rt, pt) {
				t.Errorf("Decrypt altered the payload:\ngot  %x\nwant %x", rt, pt)
			}
		})
	}

	t.Run("nil payload", func(t *testing.T) {
		if _, err := ctx.Encrypt(0, 0, nil); !errors.Is(err, ErrInvalidPayload) {
			t.Errorf("Encrypt(nil) error = %v, want %v", err, ErrInvalidPayload)
		}
		if _, err := ctx.Decrypt(0, 0, nil); !errors.Is(err, ErrInvalidPayload) {
			t.Errorf("Decrypt(nil) error = %v, want %v", err, ErrInvalidPayload)
		}
	})
}

// eea1Vector is one GSMA UEA2 specification test set usable through the
// context API (byte-aligned bit length, downlink direction).
type eea1Vector struct {
	name   string
	count  uint32
	bearer uint8
	ck     string
	plain  string
	cipher string
}

// eea1SpecVectors holds the byte-aligned UEA2 spec set 3 (specification
// version 2.1, section 4.5). Sets 1, 2, 4 and 5 use non-byte-aligned bit
// lengths that the NAS context API cannot express; they remain covered by
// snow3g.TestUEA2Vectors at the primitive level.
var eea1SpecVectors = []eea1Vector{
	{
		name:   "UEA2 spec set 3",
		count:  0xFA556B26,
		bearer: 0x03,
		ck:     "5ACB1D644C0D51204EA5F1451010D852",
		plain:  "AD9C441F 890B38C4 57A49D42 1407E8",
		cipher: "BA0F3130 0334C56B 52A7497C BAC046",
	},
}

// TestEEA1Vectors verifies EEA1 (UEA2/F8) encryption through the security
// context against the GSMA UEA2 specification vector.
func TestEEA1Vectors(t *testing.T) {
	for _, tc := range eea1SpecVectors {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &EMMSecurityContext{CK: hexBytes(t, tc.ck), EncAlgo: EEA1, Valid: true}
			plain := hexBytes(t, tc.plain)
			want := hexBytes(t, tc.cipher)

			got, err := ctx.Encrypt(tc.count, tc.bearer, plain)
			if err != nil {
				t.Fatalf("Encrypt returned error: %v", err)
			}
			if !bytes.Equal(got, want) {
				t.Errorf("ciphertext mismatch:\ngot  %x\nwant %x", got, want)
			}

			// F8 is XOR symmetric: applying the same parameters again
			// recovers the plaintext.
			back, err := ctx.Encrypt(tc.count, tc.bearer, got)
			if err != nil {
				t.Fatalf("re-Encrypt returned error: %v", err)
			}
			if !bytes.Equal(back, plain) {
				t.Errorf("re-encryption did not recover the plaintext:\ngot  %x\nwant %x", back, plain)
			}
		})
	}
}

// TestEEA1RoundTrip verifies EEA1 encrypt/decrypt round trips for several
// COUNT^C/bearer/direction/payload-size combinations, and that the two
// transmission directions produce different keystreams.
func TestEEA1RoundTrip(t *testing.T) {
	ctx := &EMMSecurityContext{
		CK:      hexBytes(t, "2BD6459F82C5B300952C49104881FF48"),
		EncAlgo: EEA1,
		Valid:   true,
	}

	tests := []struct {
		name   string
		count  uint32
		bearer uint8
		size   int
	}{
		{"small payload", 1, 1, 8},
		{"spec-like mid", 0xFA556B26, 3, 15},
		{"larger payload", 0x72A4F20F, 0x0C, 100},
		{"high count", math.MaxUint32, 0x18, 64},
		{"zero count", 0, 0, 32},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pt := make([]byte, tc.size)
			for i := range pt {
				pt[i] = uint8(0xA5 ^ i)
			}

			dl, err := ctx.Encrypt(tc.count, tc.bearer, pt)
			if err != nil {
				t.Fatalf("Encrypt returned error: %v", err)
			}
			if bytes.Equal(dl, pt) {
				t.Fatal("EEA1 left the payload unchanged")
			}
			rt, err := ctx.Encrypt(tc.count, tc.bearer, dl)
			if err != nil {
				t.Fatalf("re-Encrypt returned error: %v", err)
			}
			if !bytes.Equal(rt, pt) {
				t.Errorf("downlink round trip mismatch:\nplain %x\ncipher %x\ngot   %x", pt, dl, rt)
			}

			ul, err := ctx.Decrypt(tc.count, tc.bearer, pt)
			if err != nil {
				t.Fatalf("Decrypt returned error: %v", err)
			}
			rtUp, err := ctx.Decrypt(tc.count, tc.bearer, ul)
			if err != nil {
				t.Fatalf("re-Decrypt returned error: %v", err)
			}
			if !bytes.Equal(rtUp, pt) {
				t.Errorf("uplink round trip mismatch:\nplain %x\ncipher %x\ngot   %x", pt, ul, rtUp)
			}

			if bytes.Equal(dl, ul) {
				t.Error("downlink and uplink keystreams are identical")
			}
		})
	}
}

// eea2Vector is one TS 33.401 annex A EEA2 test set.
type eea2Vector struct {
	name      string
	count     uint32
	bearer    uint8
	direction uint8
	ck        string
	ibs       string
	obs       string
}

// eea2SpecVectors are the six TS 33.401 annex A EEA2 test sets.
var eea2SpecVectors = []eea2Vector{
	{
		name:      "annex A set 1",
		count:     0x398A59B4,
		bearer:    0x15,
		direction: KNASDirectionDownlink,
		ck:        "D3 C5 D5 92 32 7F B1 1C 40 35 C6 68 0A F8 C6 D1",
		ibs: "98 1B A6 82 4C 1B FB 1A B4 85 47 20 29 B7 1D 80" +
			"8C E3 3E 2C C3 C0 B5 FC 1F 3D E8 A6 DC 66 B1 F0",
		obs: "E9 FE D8 A6 3D 15 53 04 D7 1D F2 0B F3 E8 22 14" +
			"B2 0E D7 DA D2 F2 33 DC 3C 22 D7 BD EE ED 8E 78",
	},
	{
		name:      "annex A set 2",
		count:     0xC675A64B,
		bearer:    0x0C,
		direction: KNASDirectionDownlink,
		ck:        "2B D6 45 9F 82 C4 40 E0 95 2C 49 10 48 05 FF 48",
		ibs: "7E C6 12 72 74 3B F1 61 47 26 44 6A 6C 38 CE D1" +
			"66 F6 CA 76 EB 54 30 04 42 86 34 6C EF 13 0F 92" +
			"92 2B 03 45 0D 3A 99 75 E5 BD 2E A0 EB 55 AD 8E" +
			"1B 19 9E 3E C4 31 60 20 E9 A1 B2 85 E7 62 79 53" +
			"59 B7 BD FD 39 BE F4 B2 48 45 83 D5 AF E0 82 AE" +
			"E6 38 BF 5F D5 A6 06 19 39 01 A0 8F 4A B4 1A AB" +
			"9B 13 48",
		obs: "59 61 60 53 53 C6 4B DC A1 5B 19 5E 28 85 53 A9" +
			"10 63 25 06 D6 20 0A A7 90 C4 C8 06 C9 99 04 CF" +
			"24 45 CC 50 BB 1C F1 68 A4 96 73 73 4E 08 1B 57" +
			"E3 24 CE 52 59 C0 E7 8D 4C D9 7B 87 09 76 50 3C" +
			"09 43 F2 CB 5A E8 F0 52 C7 B7 D3 92 23 95 87 B8" +
			"95 60 86 BC AB 18 83 60 42 E2 E6 CE 42 43 2A 17" +
			"10 5C 53",
	},
	{
		name:      "annex A set 3",
		count:     0x544D49CD,
		bearer:    0x04,
		direction: KNASDirectionUplink,
		ck:        "0A 8B 6B D8 D9 B0 8B 08 D6 4E 32 D1 81 77 77 FB",
		ibs: "FD 40 A4 1D 37 0A 1F 65 74 50 95 68 7D 47 BA 1D" +
			"36 D2 34 9E 23 F6 44 39 2C 8E A9 C4 9D 40 C1 32" +
			"71 AF F2 64 D0 F2 48",
		obs: "75 75 0D 37 B4 BB A2 A4 DE DB 34 23 5B D6 8C 66" +
			"45 AC DA AC A4 81 38 A3 B0 C4 71 E2 A7 04 1A 57" +
			"64 23 D2 92 72 87 F0",
	},
	{
		name:      "annex A set 4",
		count:     0x72D8C671,
		bearer:    0x10,
		direction: KNASDirectionDownlink,
		ck:        "AA 1F 95 AE A5 33 BC B3 2E B6 3B F5 2D 8F 83 1A",
		ibs: "FB 1B 96 C5 C8 BA DF B2 E8 E8 ED FD E7 8E 57 F2" +
			"AD 81 E7 41 03 FC 43 0A 53 4D CC 37 AF CE C7 0E" +
			"15 17 BB 06 F2 72 19 DA E4 90 22 DD C4 7A 06 8D" +
			"E4 C9 49 6A 95 1A 6B 09 ED BD C8 64 C7 AD BD 74" +
			"0A C5 0C 02 2F 30 82 BA FD 22 D7 81 97 C5 D5 08" +
			"B9 77 BC A1 3F 32 E6 52 E7 4B A7 28 57 60 77 CE" +
			"62 8C 53 5E 87 DC 60 77 BA 07 D2 90 68 59 0C 8C" +
			"B5 F1 08 8E 08 2C FA 0E C9 61 30 2D 69 CF 3D",
		obs: "DF B4 40 AC B3 77 35 49 EF C0 46 28 AE B8 D8 15" +
			"62 75 23 0B DC 69 0D 94 B0 0D 8D 95 F2 8C 4B 56" +
			"30 7F 60 F4 CA 55 EB A6 61 EB BA 72 AC 80 8F A8" +
			"C4 9E 26 78 8E D0 4A 5D 60 6C B4 18 DE 74 87 8B" +
			"9A 22 F8 EF 29 59 0B C4 EB 57 C9 FA F7 C4 15 24" +
			"A8 85 B8 97 9C 42 3F 2F 8F 8E 05 92 A9 87 92 01" +
			"BE 7F F9 77 7A 16 2A B8 10 FE B3 24 BA 74 C4 C1" +
			"56 E0 4D 39 09 72 09 65 3A C3 3E 5A 5F 2D 88",
	},
	{
		name:      "annex A set 5",
		count:     0xC675A64B,
		bearer:    0x0C,
		direction: KNASDirectionDownlink,
		ck:        "96 18 AE 46 89 1F 86 57 8E EB E9 0E F7 A1 20 2E",
		ibs: "8D AA 17 B1 AE 05 05 29 C6 82 7F 28 C0 EF 6A 12" +
			"42 E9 3F 8B 31 4F B1 8A 77 F7 90 AE 04 9F ED D6" +
			"12 26 7F EC AE FC 45 01 74 D7 6D 9F 9A A7 75 5A" +
			"30 CD 90 A9 A5 87 4B F4 8E AF 70 EE A3 A6 2A 25" +
			"0A 8B 6B D8 D9 B0 8B 08 D6 4E 32 D1 81 77 77 FB" +
			"54 4D 49 CD 49 72 0E 21 9D BF 8B BE D3 39 04 E1" +
			"FD 40 A4 1D 37 0A 1F 65 74 50 95 68 7D 47 BA 1D" +
			"36 D2 34 9E 23 F6 44 39 2C 8E A9 C4 9D 40 C1 32" +
			"71 AF F2 64 D0 F2 48 41 D6 46 5F 09 96 FF 84 E6" +
			"5F C5 17 C5 3E FC 33 63 C3 84 92",
		obs: "91 9C 8C 33 D6 67 89 70 3D 05 A0 D7 CE 82 A2 AE" +
			"AC 4E E7 6C 0F 4D A0 50 33 5E 8A 84 E7 89 7B A5" +
			"DF 2F 36 BD 51 3E 3D 0C 85 78 C7 A0 FC F0 43 E0" +
			"3A A3 A3 9F BA AD 7D 15 BE 07 4F AA 5D 90 29 F7" +
			"1F B4 57 B6 47 83 47 14 B0 E1 8F 11 7F CA 10 67" +
			"79 45 09 6C 8C 5F 32 6B A8 D6 09 5E B2 9C 3E 36" +
			"CF 24 5D 16 22 AA FE 92 1F 75 66 C4 F5 D6 44 F2" +
			"F1 FC 0E C6 84 DD B2 13 49 74 76 22 E2 09 29 5D" +
			"27 FF 3F 95 62 33 71 D4 9B 14 7C 0A F4 86 17 1F" +
			"22 CD 04 B1 CB EB 26 58 22 3E 69",
	},
	{
		name:      "annex A set 6",
		count:     0xACA4F50F,
		bearer:    0x0B,
		direction: KNASDirectionUplink,
		ck:        "54 F4 E2 E0 4C 83 78 6E EC 8F B5 AB E8 E3 65 66",
		ibs: "40 98 1B A6 82 4C 1B FB 42 86 B2 99 78 3D AF 44" +
			"2C 09 9F 7A B0 F5 8D 5C 8E 46 B1 04 F0 8F 01 B4" +
			"1A B4 85 47 20 29 B7 1D 36 BD 1A 3D 90 DC 3A 41" +
			"B4 6D 51 67 2A C4 C9 66 3A 2B E0 63 DA 4B C8 D2" +
			"80 8C E3 3E 2C CC BF C6 34 E1 B2 59 06 08 76 A0" +
			"FB B5 A4 37 EB CC 8D 31 C1 9E 44 54 31 87 45 E3" +
			"FA 16 BB 11 AD AE 24 88 79 FE 52 DB 25 43 E5 3C" +
			"F4 45 D3 D8 28 CE 0B F5 C5 60 59 3D 97 27 8A 59" +
			"76 2D D0 C2 C9 CD 68 D4 49 6A 79 25 08 61 40 14" +
			"B1 3B 6A A5 11 28 C1 8C D6 A9 0B 87 97 8C 2F F1" +
			"CA BE 7D 9F 89 8A 41 1B FD B8 4F 68 F6 72 7B 14" +
			"99 CD D3 0D F0 44 3A B4 A6 66 53 33 0B CB A1 10" +
			"5E 4C EC 03 4C 73 E6 05 B4 31 0E AA AD CF D5 B0" +
			"CA 27 FF D8 9D 14 4D F4 79 27 59 42 7C 9C C1 F8" +
			"CD 8C 87 20 23 64 B8 A6 87 95 4C B0 5A 8D 4E 2D" +
			"99 E7 3D B1 60 DE B1 80 AD 08 41 E9 67 41 A5 D5" +
			"9F E4 18 9F 15 42 00 26 FE 4C D1 21 04 93 2F B3" +
			"8F 73 53 40 43 8A AF 7E CA 6F D5 CF D3 A1 95 CE" +
			"5A BE 65 27 2A F6 07 AD A1 BE 65 A6 B4 C9 C0 69" +
			"32 34 09 2C 4D 01 8F 17 56 C6 DB 9D C8 A6 D8 0B" +
			"88 81 38 61 6B 68 12 62 F9 54 D0 E7 71 17 48 78" +
			"0D 92 29 1D 86 29 99 72 DB 74 1C FA 4F 37 B8 B5" +
			"6C DB 18 A7 CA 82 18 E8 6E 4B 4B 71 6A 4D 04 37" +
			"1F BE C2 62 FC 5A D0 B3 81 9B 18 7B 97 E5 5B 1A" +
			"4D 7C 19 EE 24 C8 B4 D7 72 3C FE DF 04 5B 8A CA" +
			"E4 86 95 17 D8 0E 50 61 5D 90 35 D5 D9 C5 A4 0A" +
			"F6 02 28 0B 54 25 97 B0 CB 18 61 9E EB 35 92 57" +
			"59 D1 95 E1 00 E8 E4 AA 0C 38 A3 C2 AB E0 F3 D8" +
			"FF 04 F3 C3 3C 29 50 69 C2 36 94 B5 BB EA CD D5" +
			"42 E2 8E 8A 94 ED B9 11 9F 41 2D 05 4B E1 FA 72" +
			"00 B0 90",
		obs: "5C B7 2C 6E DC 87 8F 15 66 E1 02 53 AF C3 64 C9" +
			"FA 54 0D 91 4D B9 4C BE E2 75 D0 91 7C A6 AF 0D" +
			"77 AC B4 EF 3B BE 1A 72 2B 2E F5 BD 1D 4B 8E 2A" +
			"A5 02 4E C1 38 8A 20 1E 7B CE 79 20 AE C6 15 89" +
			"5F 76 3A 55 64 DC C4 C4 82 A2 EE 1D 8B FE CC 44" +
			"98 EC A8 3F BB 75 F9 AB 53 0E 0D AF BE DE 2F A5" +
			"89 5B 82 99 1B 62 77 C5 29 E0 F2 52 9D 7F 79 60" +
			"6B E9 67 06 29 6D ED FA 9D 74 12 B6 16 95 8C B5" +
			"63 C6 78 C0 28 25 C3 0D 0A EE 77 C4 C1 46 D2 76" +
			"54 12 42 1A 80 8D 13 CE C8 19 69 4C 75 AD 57 2E" +
			"9B 97 3D 94 8B 81 A9 33 7C 3B 2A 17 19 2E 22 C2" +
			"06 9F 7E D1 16 2A F4 4C DE A8 17 60 36 65 E8 07" +
			"CE 40 C8 E0 DD 9D 63 94 DC 6E 31 15 3F E1 95 5C" +
			"47 AF B5 1F 26 17 EE 0C 5E 3B 8E F1 AD 75 74 ED" +
			"34 3E DC 27 43 CC 94 C9 90 E1 F1 FD 26 42 53 C1" +
			"78 DE A7 39 C0 BE FE EB CD 9F 9B 76 D4 9C 10 15" +
			"C9 FE CF 50 E5 3B 8B 52 04 DB CD 3E ED 86 38 55" +
			"DA BC DC C9 4B 31 E3 18 02 15 68 85 5C 8B 9E 52" +
			"A9 81 95 7A 11 28 27 F9 78 BA 96 0F 14 47 91 1B" +
			"31 7B 55 11 FB CC 7F B1 3A C1 53 DB 74 25 11 17" +
			"E4 86 1E B9 E8 3B FF FF C4 EB 77 55 57 90 38 E5" +
			"79 24 B1 F7 8B 3E 1A D9 0B AB 2A 07 87 1B 72 DB" +
			"5E EF 96 C3 34 04 49 66 DB 0C 37 CA FD 1A 89 E5" +
			"64 6A 35 80 EB 64 65 F1 21 DC E9 CB 88 D8 5B 96" +
			"CF 23 CC CC D4 28 07 67 BE E8 EE B2 3D 86 52 46" +
			"1D B6 49 31 03 00 3B AF 89 F5 E1 82 61 EA 43 C8" +
			"4A 92 EB FF FF E4 90 9D C4 6C 51 92 F8 25 F7 70" +
			"60 0B 96 02 C5 57 B5 F8 B4 31 A7 9D 45 97 7D D9" +
			"C4 1B 86 3D A9 E1 42 E9 00 20 CF D0 74 D6 92 7B" +
			"7A B3 B6 72 5D 1A 6F 3F 98 B9 C9 DA A8 98 2A FF" +
			"06 78 28",
	},
}

// TestEEA2Vectors verifies EEA2 (AES-CTR) encryption and decryption through
// the security context against the TS 33.401 annex A vectors. Encrypt pins
// the downlink direction and Decrypt the uplink direction, so each set is
// applied through the method matching its vector direction.
func TestEEA2Vectors(t *testing.T) {
	for _, tc := range eea2SpecVectors {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &EMMSecurityContext{CK: hexBytes(t, tc.ck), EncAlgo: EEA2, Valid: true}
			ibs := hexBytes(t, tc.ibs)
			obs := hexBytes(t, tc.obs)

			apply := func(data []byte) ([]byte, error) {
				if tc.direction == KNASDirectionDownlink {
					return ctx.Encrypt(tc.count, tc.bearer, data)
				}
				return ctx.Decrypt(tc.count, tc.bearer, data)
			}

			got, err := apply(ibs)
			if err != nil {
				t.Fatalf("cipher returned error: %v", err)
			}
			if !bytes.Equal(got, obs) {
				t.Errorf("ciphertext mismatch:\ngot  %x\nwant %x", got, obs)
			}

			back, err := apply(got)
			if err != nil {
				t.Fatalf("inverse cipher returned error: %v", err)
			}
			if !bytes.Equal(back, ibs) {
				t.Errorf("decryption did not recover the input:\ngot  %x\nwant %x", back, ibs)
			}
		})
	}
}

// TestCipheringDependsOnCountAndAlgo verifies that the ciphertext changes
// with the NAS COUNT^C value for both EEA1 and EEA2 and that the two
// algorithms produce different output for identical inputs.
func TestCipheringDependsOnCountAndAlgo(t *testing.T) {
	payload := make([]byte, 32)
	for i := range payload {
		payload[i] = uint8(i)
	}
	ck := hexBytes(t, "D3C5D592327FB11C4035C6680AF8C6D1")
	e1 := &EMMSecurityContext{CK: ck, EncAlgo: EEA1, Valid: true}
	e2 := &EMMSecurityContext{CK: ck, EncAlgo: EEA2, Valid: true}
	const count = 0x00000001

	dl1, err := e1.Encrypt(count, 5, payload)
	if err != nil {
		t.Fatalf("EEA1 Encrypt returned error: %v", err)
	}
	dl1Next, err := e1.Encrypt(count+1, 5, payload)
	if err != nil {
		t.Fatalf("EEA1 Encrypt returned error: %v", err)
	}
	if bytes.Equal(dl1, dl1Next) {
		t.Error("EEA1 ciphertext is identical for consecutive COUNT^C values")
	}

	dl2, err := e2.Encrypt(count, 5, payload)
	if err != nil {
		t.Fatalf("EEA2 Encrypt returned error: %v", err)
	}
	dl2Next, err := e2.Encrypt(count+1, 5, payload)
	if err != nil {
		t.Fatalf("EEA2 Encrypt returned error: %v", err)
	}
	if bytes.Equal(dl2, dl2Next) {
		t.Error("EEA2 ciphertext is identical for consecutive COUNT^C values")
	}

	if bytes.Equal(dl1, dl2) {
		t.Error("EEA1 and EEA2 produced identical ciphertext")
	}
}

// TestEIA0MAC verifies that the null integrity algorithm returns a
// four-octet all-zero MAC-I regardless of the message parameters, that the
// zero MAC verifies, and that a non-zero MAC is rejected.
func TestEIA0MAC(t *testing.T) {
	ctx, err := SetupContext(hexBytes(t, kasmeAnchorA), 0, EEA0, EIA0)
	if err != nil {
		t.Fatalf("SetupContext returned error: %v", err)
	}

	tests := []struct {
		name      string
		count     uint32
		fresh     uint32
		direction uint8
		size      int
	}{
		{"uplink small", 0x38A6F056, 0x1F, KNASDirectionUplink, 16},
		{"downlink larger", 0x36AF6144, 0x18, KNASDirectionDownlink, 32},
		{"zero count", 0, 0, KNASDirectionUplink, 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			msg := make([]byte, tc.size)
			for i := range msg {
				msg[i] = uint8(i)
			}
			mac, err := ctx.ComputeMACI(tc.count, tc.fresh, tc.direction, msg)
			if err != nil {
				t.Fatalf("ComputeMACI returned error: %v", err)
			}
			if len(mac) != 4 {
				t.Fatalf("MAC-I length = %d, want 4", len(mac))
			}
			for i, b := range mac {
				if b != 0 {
					t.Fatalf("MAC-I[%d] = 0x%02x, want 0x00", i, b)
				}
			}

			ok, err := ctx.VerifyMACI(tc.count, tc.fresh, tc.direction, msg, mac)
			if err != nil || !ok {
				t.Errorf("VerifyMACI(zeros) = (%v, %v), want (true, nil)", ok, err)
			}
			ok, err = ctx.VerifyMACI(tc.count, tc.fresh, tc.direction, msg, []byte{0, 0, 0, 1})
			if err != nil || ok {
				t.Errorf("VerifyMACI(non-zero) = (%v, %v), want (false, nil)", ok, err)
			}
		})
	}

	t.Run("nil message", func(t *testing.T) {
		if _, err := ctx.ComputeMACI(0, 0, KNASDirectionUplink, nil); !errors.Is(err, ErrInvalidPayload) {
			t.Errorf("ComputeMACI(nil) error = %v, want %v", err, ErrInvalidPayload)
		}
	})
}

// TestComputeMACIRejectsEmptyMessage verifies that ComputeMACI and VerifyMACI
// reject an empty (non-nil) message for every supported integrity algorithm
// with ErrInvalidPayload instead of panicking inside the UIA2 F9 primitive,
// which indexes out of bounds for a zero-bit input.
func TestComputeMACIRejectsEmptyMessage(t *testing.T) {
	tests := []struct {
		name      string
		algo      ie.NASIntegrityAlgorithm
		ik        string
		direction uint8
	}{
		{"EIA0 uplink", EIA0, "", KNASDirectionUplink},
		{"EIA1 downlink", EIA1, "C736C6AAB22BFFF91E2698D2E22AD57E", KNASDirectionDownlink},
		{"EIA2 uplink", EIA2, "83FD23A244A74CF358DA3019F1722635", KNASDirectionUplink},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &EMMSecurityContext{IntAlgo: tc.algo, Valid: true}
			if tc.ik != "" {
				ctx.IK = hexBytes(t, tc.ik)
			}
			if _, err := ctx.ComputeMACI(1, 0x05, tc.direction, []byte{}); !errors.Is(err, ErrInvalidPayload) {
				t.Errorf("ComputeMACI(empty) error = %v, want %v", err, ErrInvalidPayload)
			}
		})
	}

	t.Run("VerifyMACI EIA1 empty message", func(t *testing.T) {
		ctx := &EMMSecurityContext{IK: hexBytes(t, "C736C6AAB22BFFF91E2698D2E22AD57E"), IntAlgo: EIA1, Valid: true}
		ok, err := ctx.VerifyMACI(1, 0x05, KNASDirectionDownlink, []byte{}, make([]byte, macISize))
		if ok || !errors.Is(err, ErrInvalidPayload) {
			t.Errorf("VerifyMACI(empty) = (%v, %v), want (false, %v)", ok, err, ErrInvalidPayload)
		}
	})
}

// eia1Vector is one GSMA UIA2 specification test set usable through the
// context API (byte-aligned bit length).
type eia1Vector struct {
	name      string
	count     uint32
	fresh     uint32
	direction uint8
	ik        string
	msg       string
	macI      string
}

// eia1SpecVectors hold the byte-aligned UIA2 spec sets 4 and 5 (specification
// version 2.1, sections 5.6 and 5.7). Sets 1-3 use non-byte-aligned bit
// lengths that the NAS context API cannot express, and set 6 shares its
// 2056-octet message with the EIA2 long vector in TestEIA2Vectors; all of
// them remain covered at the snow3g package level.
var eia1SpecVectors = []eia1Vector{
	{
		name:      "UIA2 spec set 4",
		count:     0x14793E41,
		fresh:     0x0397E8FD,
		direction: KNASDirectionDownlink,
		ik:        "C736C6AAB22BFFF91E2698D2E22AD57E",
		msg:       "D0A7D463DF9FB2B2 78833FA02E235AA1 72BD970C1473E129 07FB648B6599AAA0 B24A038665422B20 A499276A50427009",
		macI:      "38B554C0",
	},
	{
		name:      "UIA2 spec set 5",
		count:     0x296F393C,
		fresh:     0x6B227737,
		direction: KNASDirectionDownlink,
		ik:        "F4EBEC69E73EAF2EB2CF6AF4B3120FFD",
		msg:       "10BFFF839E0C7165 8DBB2D1707E14572 4F41C16F48BF403C 3B18E38FD5D1663B 6F6D900193E3CEA8 BB4F1B4F5BE82203 2232A78D7D75238D 5E6DAECD3B4322CF 59BC7EA84AB18811 B5BFB7BC553F4FE4 4478CE287A148799 90D18D12CA79D2C8 55149021CD5CE8CA 0371CA04FCCE143E 3D7CFEE94585B588 5CAC46068B",
		macI:      "061745AE",
	},
}

// TestEIA1Vectors verifies EIA1 (UIA2/F9) MAC-I computation through the
// security context against the GSMA UIA2 specification vectors.
func TestEIA1Vectors(t *testing.T) {
	for _, tc := range eia1SpecVectors {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &EMMSecurityContext{IK: hexBytes(t, tc.ik), IntAlgo: EIA1, Valid: true}
			msg := hexBytes(t, tc.msg)
			want := hexBytes(t, tc.macI)

			mac, err := ctx.ComputeMACI(tc.count, tc.fresh, tc.direction, msg)
			if err != nil {
				t.Fatalf("ComputeMACI returned error: %v", err)
			}
			if !bytes.Equal(mac, want) {
				t.Errorf("MAC-I mismatch:\ngot  %x\nwant %x", mac, want)
			}

			ok, err := ctx.VerifyMACI(tc.count, tc.fresh, tc.direction, msg, mac)
			if err != nil || !ok {
				t.Errorf("VerifyMACI = (%v, %v), want (true, nil)", ok, err)
			}
		})
	}
}

// eia2Vector is one TS 33.401 annex A EIA2 test set.
type eia2Vector struct {
	name      string
	count     uint32
	fresh     uint8
	direction uint8
	ik        string
	msg       string
	macI      string
}

// eia2SpecVectors are the TS 33.401 annex A EIA2 test sets 2, 5 and 8
// (short, multi-block and long messages).
var eia2SpecVectors = []eia2Vector{
	{
		name:      "annex A set 2",
		count:     0x398A59B4,
		fresh:     0x1A,
		direction: KNASDirectionDownlink,
		ik:        "D3 C5 D5 92 32 7F B1 1C 40 35 C6 68 0A F8 C6 D1",
		msg:       "48 45 83 D5 AF E0 82 AE",
		macI:      "B93787E6",
	},
	{
		name:      "annex A set 5",
		count:     0x36AF6144,
		fresh:     0x0F,
		direction: KNASDirectionDownlink,
		ik:        "83 FD 23 A2 44 A7 4C F3 58 DA 30 19 F1 72 26 35",
		msg: "35 C6 87 16 63 3C 66 FB 75 0C 26 68 65 D5 3C 11" +
			"EA 05 B1 E9 FA 49 C8 39 8D 48 E1 EF A5 90 9D 39" +
			"47 90 28 37 F5 AE 96 D5 A0 5B C8 D6 1C A8 DB EF" +
			"1B 13 A4 B4 AB FE 4F B1 00 60 45 B6 74 BB 54 72" +
			"93 04 C3 82 BE 53 A5 AF 05 55 61 76 F6 EA A2 EF" +
			"1D 05 E4 B0 83 18 1E E6 74 CD A5 A4 85 F7 4D 7A",
		macI: "E657E182",
	},
	{
		name:      "annex A set 8",
		count:     0x296F393C,
		fresh:     0x0B,
		direction: KNASDirectionDownlink,
		ik:        "B3 12 0F FD B2 CF 6A F4 E7 3E AF 2E F4 EB EC 69",
		msg: "00 00 00 00 00 00 00 00 01 01 01 01 01 01 01 01" +
			"E0 95 80 45 F3 A0 BB A4 E3 96 83 46 F0 A3 B8 A7" +
			"C0 2A 01 8A E6 40 76 52 26 B9 87 C9 13 E6 CB F0" +
			"83 57 00 16 CF 83 EF BC 61 C0 82 51 3E 21 56 1A" +
			"42 7C 00 9D 28 C2 98 EF AC E7 8E D6 D5 6C 2D 45" +
			"05 AD 03 2E 9C 04 DC 60 E7 3A 81 69 6D A6 65 C6" +
			"C4 86 03 A5 7B 45 AB 33 22 15 85 E6 8E E3 16 91" +
			"87 FB 02 39 52 86 32 DD 65 6C 80 7E A3 24 8B 7B" +
			"46 D0 02 B2 B5 C7 45 8E B8 5B 9C E9 58 79 E0 34" +
			"08 59 05 5E 3B 0A BB C3 EA CE 87 19 CA A8 02 65" +
			"C9 72 05 D5 DC 4B CC 90 2F E1 83 96 29 ED 71 32" +
			"8A 0F 04 49 F5 88 55 7E 68 98 86 0E 04 2A EC D8" +
			"4B 24 04 C2 12 C9 22 2D A5 BF 8A 89 EF 67 97 87" +
			"0C F5 07 71 A6 0F 66 A2 EE 62 85 36 57 AD DF 04" +
			"CD DE 07 FA 41 4E 11 F1 2B 4D 81 B9 B4 E8 AC 53" +
			"8E A3 06 66 68 8D 88 1F 6C 34 84 21 99 2F 31 B9" +
			"4F 88 06 ED 8F CC FF 4C 91 23 B8 96 42 52 7A D6" +
			"13 B1 09 BF 75 16 74 85 F1 26 8B F8 84 B4 CD 23" +
			"D2 9A 09 34 92 57 03 D6 34 09 8F 77 67 F1 BE 74" +
			"91 E7 08 A8 BB 94 9A 38 73 70 8A EF 4A 36 23 9E" +
			"50 CC 08 23 5C D5 ED 6B BE 57 86 68 A1 7B 58 C1" +
			"17 1D 0B 90 E8 13 A9 E4 F5 8A 89 D7 19 B1 10 42" +
			"D6 36 0B 1B 0F 52 DE B7 30 A5 8D 58 FA F4 63 15" +
			"95 4B 0A 87 26 91 47 59 77 DC 88 C0 D7 33 FE FF" +
			"54 60 0A 0C C1 D0 30 0A AA EB 94 57 2C 6E 95 B0" +
			"1A E9 0D E0 4F 1D CE 47 F8 7E 8F A7 BE BF 77 E1" +
			"DB C2 0D 6B A8 5C B9 14 3D 51 8B 28 5D FA 04 B6" +
			"98 BF 0C F7 81 9F 20 FA 7A 28 8E B0 70 3D 99 5C" +
			"59 94 0C 7C 66 DE 57 A9 B7 0F 82 37 9B 70 E2 03" +
			"1E 45 0F CF D2 18 13 26 FC D2 8D 88 23 BA AA 80" +
			"DF 6E 0F 44 35 59 64 75 39 FD 89 07 C0 FF D9 D7" +
			"9C 13 0E D8 1C 9A FD 9B 7E 84 8C 9F ED 38 44 3D" +
			"5D 38 0E 53 FB DB 8A C8 C3 D3 F0 68 76 05 4F 12" +
			"24 61 10 7D E9 2F EA 09 C6 F6 92 3A 18 8D 53 AF" +
			"E5 4A 10 F6 0E 6E 9D 5A 03 D9 96 B5 FB C8 20 F8" +
			"A6 37 11 6A 27 AD 04 B4 44 A0 93 2D D6 0F BD 12" +
			"67 1C 11 E1 C0 EC 73 E7 89 87 9F AA 3D 42 C6 4D" +
			"20 CD 12 52 74 2A 37 68 C2 5A 90 15 85 88 8E CE" +
			"E1 E6 12 D9 93 6B 40 3B 07 75 94 9A 66 CD FD 99" +
			"A2 9B 13 45 BA A8 D9 D5 40 0C 91 02 4B 0A 60 73" +
			"63 B0 13 CE 5D E9 AE 86 9D 3B 8D 95 B0 57 0B 3C" +
			"2D 39 14 22 D3 24 50 CB CF AE 96 65 22 86 E9 6D" +
			"EC 12 14 A9 34 65 27 98 0A 81 92 EA C1 C3 9A 3A" +
			"AF 6F 15 35 1D A6 BE 76 4D F8 97 72 EC 04 07 D0" +
			"6E 44 15 BE FA E7 C9 25 80 DF 9B F5 07 49 7C 8F" +
			"29 95 16 0D 4E 21 8D AA CB 02 94 4A BF 83 34 0C" +
			"E8 BE 16 86 A9 60 FA F9 0E 2D 90 C5 5C C6 47 5B" +
			"AB C3 17 1A 80 A3 63 17 49 54 95 5D 71 01 DA B1" +
			"6A E8 17 91 67 E2 14 44 B4 43 A9 EA AA 7C 91 DE" +
			"36 D1 18 C3 9D 38 9F 8D D4 46 9A 84 6C 9A 26 2B" +
			"F7 FA 18 48 7A 79 E8 DE 11 69 9E 0B 8F DF 55 7C" +
			"B4 87 19 D4 53 BA 71 30 56 10 9B 93 A2 18 C8 96" +
			"75 AC 19 5F B4 FB 06 63 9B 37 97 14 49 55 B3 C9" +
			"32 7D 1A EC 00 3D 42 EC D0 EA 98 AB F1 9F FB 4A" +
			"F3 56 1A 67 E7 7C 35 BF 15 C5 9C 24 12 DA 88 1D" +
			"B0 2B 1B FB CE BF AC 51 52 BC 99 BC 3F 1D 15 F7" +
			"71 00 1B 70 29 FE DB 02 8F 8B 85 2B C4 40 7E B8" +
			"3F 89 1C 9C A7 33 25 4F DD 1E 9E DB 56 91 9C E9" +
			"FE A2 1C 17 40 72 52 1C 18 31 9A 54 B5 D4 EF BE" +
			"BD DF 1D 8B 69 B1 CB F2 5F 48 9F CC 98 13 72 54" +
			"7C F4 1D 00 8E F0 BC A1 92 6F 93 4B 73 5E 09 0B" +
			"3B 25 1E B3 3A 36 F8 2E D9 B2 9C F4 CB 94 41 88" +
			"FA 0E 1E 38 DD 77 8F 7D 1C 9D 98 7B 28 D1 32 DF" +
			"B9 73 1F A4 F4 B4 16 93 5B E4 9D E3 05 16 AF 35" +
			"78 58 1F 2F 13 F5 61 C0 66 33 61 94 1E AB 24 9A" +
			"4B C1 23 F8 D1 5C D7 11 A9 56 A1 BF 20 FE 6E B7" +
			"8A EA 23 73 36 1D A0 42 6C 79 A5 30 C3 BB 1D E0" +
			"C9 97 22 EF 1F DE 39 AC 2B 00 A0 A8 EE 7C 80 0A" +
			"08 BC 22 64 F8 9F 4E FF E6 27 AC 2F 05 31 FB 55" +
			"4F 6D 21 D7 4C 59 0A 70 AD FA A3 90 BD FB B3 D6" +
			"8E 46 21 5C AB 18 7D 23 68 D5 A7 1F 5E BE C0 81" +
			"CD 3B 20 C0 82 DB E4 CD 2F AC A2 87 73 79 5D 6B" +
			"0C 10 20 4B 65 9A 93 9E F2 9B BE 10 88 24 36 24" +
			"42 99 27 A7 EB 57 6D D3 A0 0E A5 E0 1A F5 D4 75" +
			"83 B2 27 2C 0C 16 1A 80 65 21 A1 6F F9 B0 A7 22" +
			"C0 CF 26 B0 25 D5 83 6E 22 58 A4 F7 D4 77 3A C8" +
			"01 E4 26 3B C2 94 F4 3D EF 7F A8 70 3F 3A 41 97" +
			"46 35 25 88 76 52 B0 B2 A4 A2 A7 CF 87 F0 09 14" +
			"87 1E 25 03 91 13 C7 E1 61 8D A3 40 64 B5 7A 43" +
			"C4 63 24 9F B8 D0 5E 0F 26 F4 A6 D8 49 72 E7 A9" +
			"05 48 24 14 5F 91 29 5C DB E3 9A 6F 92 0F AC C6" +
			"59 71 2B 46 A5 4B A2 95 BB E6 A9 01 54 E9 1B 33" +
			"98 5A 2B CD 42 0A D5 C6 7E C9 AD 8E B7 AC 68 64" +
			"DB 27 2A 51 6B C9 4C 28 39 B0 A8 16 9A 6B F5 8E" +
			"1A 0C 2A DA 8C 88 3B 7B F4 97 A4 91 71 26 8E D1" +
			"5D DD 29 69 38 4E 7F F4 BF 4A AB 2E C9 EC C6 52" +
			"9C F6 29 E2 DF 0F 08 A7 7A 65 AF A1 2A A9 B5 05" +
			"DF 8B 28 7E F6 CC 91 49 3D 1C AA 39 07 6E 28 EF" +
			"1E A0 28 F5 11 8D E6 1A E0 2B B6 AE FC 33 43 A0" +
			"50 29 2F 19 9F 40 18 57 B2 BE AD 5E 6E E2 A1 F1" +
			"91 02 2F 92 78 01 6F 04 77 91 A9 D1 8D A7 D2 A6" +
			"D2 7F 2E 0E 51 C2 F6 EA 30 E8 AC 49 A0 60 4F 4C" +
			"13 54 2E 85 B6 83 81 B9 FD CF A0 CE 4B 2D 34 13" +
			"54 85 2D 36 02 45 C5 36 B6 12 AF 71 F3 E7 7C 90" +
			"95 AE 2D BD E5 04 B2 65 73 3D AB FE 10 A2 0F C7" +
			"D6 D3 2C 21 CC C7 2B 8B 34 44 AE 66 3D 65 92 2D" +
			"17 F8 2C AA 2B 86 5C D8 89 13 D2 91 A6 58 99 02" +
			"6E A1 32 84 39 72 3C 19 8C 36 B0 C3 C8 D0 85 BF" +
			"AF 8A 32 0F DE 33 4B 4A 49 19 B4 4C 2B 95 F6 E8" +
			"EC F7 33 93 F7 F0 D2 A4 0E 60 B1 D4 06 52 6B 02" +
			"2D DC 33 18 10 B1 A5 F7 C3 47 BD 53 ED 1F 10 5D" +
			"6A 0D 30 AB A4 77 E1 78 88 9A B2 EC 55 D5 58 DE" +
			"AB 26 30 20 43 36 96 2B 4D B5 B6 63 B6 90 2B 89" +
			"E8 5B 31 BC 6A F5 0F C5 0A CC B3 FB 9B 57 B6 63" +
			"29 70 31 37 8D B4 78 96 D7 FB AF 6C 60 0A DD 2C" +
			"67 F9 36 DB 03 79 86 DB 85 6E B4 9C F2 DB 3F 7D" +
			"A6 D2 36 50 E4 38 F1 88 40 41 B0 13 11 9E 4C 2A" +
			"E5 AF 37 CC CD FB 68 66 07 38 B5 8B 3C 59 D1 C0" +
			"24 84 37 47 2A BA 1F 35 CA 1F B9 0C D7 14 AA 9F" +
			"63 55 34 F4 9E 7C 5B BA 81 C2 B6 B3 6F DE E2 1C" +
			"A2 7E 34 7F 79 3D 2C E9 44 ED B2 3C 8C 9B 91 4B" +
			"E1 03 35 E3 50 FE B5 07 03 94 B7 A4 A1 5C 0C A1" +
			"20 28 35 68 B7 BF C2 54 FE 83 8B 13 7A 21 47 CE" +
			"7C 11 3A 3A 4D 65 49 9D 9E 86 B8 7D BC C7 F0 3B" +
			"BD 3A 3A B1 AA 24 3E CE 5B A9 BC F2 5F 82 83 6C" +
			"FE 47 3B 2D 83 E7 A7 20 1C D0 B9 6A 72 45 1E 86" +
			"3F 6C 3B A6 64 A6 D0 73 D1 F7 B5 ED 99 08 65 D9" +
			"78 BD 38 15 D0 60 94 FC 9A 2A BA 52 21 C2 2D 5A" +
			"B9 96 38 9E 37 21 E3 AF 5F 05 BE DD C2 87 5E 0D" +
			"FA EB 39 02 1E E2 7A 41 18 7C BB 45 EF 40 C3 E7" +
			"3B C0 39 89 F9 A3 0D 12 C5 4B A7 D2 14 1D A8 A8" +
			"75 49 3E 65 77 6E F3 5F 97 DE BC 22 86 CC 4A F9" +
			"B4 62 3E EE 90 2F 84 0C 52 F1 B8 AD 65 89 39 AE" +
			"F7 1F 3F 72 B9 EC 1D E2 15 88 BD 35 48 4E A4 44" +
			"36 34 3F F9 5E AD 6A B1 D8 AF B1 B2 A3 03 DF 1B" +
			"71 E5 3C 4A EA 6B 2E 3E 93 72 BE 0D 1B C9 97 98" +
			"B0 CE 3C C1 0D 2A 59 6D 56 5D BA 82 F8 8C E4 CF" +
			"F3 B3 3D 5D 24 E9 C0 83 11 24 BF 1A D5 4B 79 25" +
			"32 98 3D D6 C3 A8 B7 D0",
		macI: "EBD5CCB0",
	},
}

// TestEIA2Vectors verifies EIA2 (AES-CMAC) MAC-I computation through the
// security context against the TS 33.401 annex A vectors, and that a FRESH
// value wider than eight bits is rejected.
func TestEIA2Vectors(t *testing.T) {
	for _, tc := range eia2SpecVectors {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &EMMSecurityContext{IK: hexBytes(t, tc.ik), IntAlgo: EIA2, Valid: true}
			msg := hexBytes(t, tc.msg)
			want := hexBytes(t, tc.macI)

			mac, err := ctx.ComputeMACI(tc.count, uint32(tc.fresh), tc.direction, msg)
			if err != nil {
				t.Fatalf("ComputeMACI returned error: %v", err)
			}
			if !bytes.Equal(mac, want) {
				t.Errorf("MAC-I mismatch:\ngot  %x\nwant %x", mac, want)
			}

			ok, err := ctx.VerifyMACI(tc.count, uint32(tc.fresh), tc.direction, msg, mac)
			if err != nil || !ok {
				t.Errorf("VerifyMACI = (%v, %v), want (true, nil)", ok, err)
			}
		})
	}

	t.Run("fresh exceeds 8 bits", func(t *testing.T) {
		ctx := &EMMSecurityContext{IK: hexBytes(t, eia2SpecVectors[0].ik), IntAlgo: EIA2, Valid: true}
		if _, err := ctx.ComputeMACI(1, 0x100, KNASDirectionDownlink, []byte{0x01}); !errors.Is(err, ErrInvalidPayload) {
			t.Errorf("ComputeMACI(fresh=0x100) error = %v, want %v", err, ErrInvalidPayload)
		}
	})
}

// TestVerifyMACIRejectsTampered verifies that VerifyMACI accepts a valid
// MAC-I and rejects a tampered message, a tampered MAC-I and a
// wrongly-lengthed MAC-I, for both SNOW 3G and AES based algorithms and in
// both transmission directions.
func TestVerifyMACIRejectsTampered(t *testing.T) {
	msg := []byte("security mode command payload")

	tests := []struct {
		name      string
		algo      ie.NASIntegrityAlgorithm
		ik        string
		direction uint8
	}{
		{"EIA1 downlink", EIA1, "C736C6AAB22BFFF91E2698D2E22AD57E", KNASDirectionDownlink},
		{"EIA1 uplink", EIA1, "C736C6AAB22BFFF91E2698D2E22AD57E", KNASDirectionUplink},
		{"EIA2 downlink", EIA2, "83FD23A244A74CF358DA3019F1722635", KNASDirectionDownlink},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &EMMSecurityContext{IK: hexBytes(t, tc.ik), IntAlgo: tc.algo, Valid: true}
			const (
				count = 1
				fresh = 0x05
			)

			mac, err := ctx.ComputeMACI(count, fresh, tc.direction, msg)
			if err != nil {
				t.Fatalf("ComputeMACI returned error: %v", err)
			}
			ok, err := ctx.VerifyMACI(count, fresh, tc.direction, msg, mac)
			if err != nil || !ok {
				t.Fatalf("VerifyMACI(valid) = (%v, %v), want (true, nil)", ok, err)
			}

			tamperedMsg := append([]byte(nil), msg...)
			tamperedMsg[0] ^= 0x01
			ok, err = ctx.VerifyMACI(count, fresh, tc.direction, tamperedMsg, mac)
			if err != nil || ok {
				t.Errorf("VerifyMACI(tampered message) = (%v, %v), want (false, nil)", ok, err)
			}

			tamperedMac := append([]byte(nil), mac...)
			tamperedMac[3] ^= 0x80
			ok, err = ctx.VerifyMACI(count, fresh, tc.direction, msg, tamperedMac)
			if err != nil || ok {
				t.Errorf("VerifyMACI(tampered MAC-I) = (%v, %v), want (false, nil)", ok, err)
			}

			ok, err = ctx.VerifyMACI(count, fresh, tc.direction, msg, mac[:3])
			if !errors.Is(err, ErrInvalidMACI) {
				t.Errorf("VerifyMACI(short MAC-I) error = %v, want %v", err, ErrInvalidMACI)
			}
		})
	}
}

// TestGeteKSIParseeKSI verifies the eKSI octet packing/unpacking: a full
// round trip over every valid KSI/SN combination, the mask boundaries, and
// that out-of-range input bits are masked away.
func TestGeteKSIParseeKSI(t *testing.T) {
	for sn := uint8(0); sn <= 0x0F; sn++ {
		for ksi := uint8(0); ksi <= 0x07; ksi++ {
			packed := GeteKSI(ksi, sn)
			gotKsi, gotSn := ParseeKSI(packed)
			if gotKsi != ksi || gotSn != sn {
				t.Errorf("round trip failed: GeteKSI(%d,%d)=0x%02X parsed back as (%d,%d)",
					ksi, sn, packed, gotKsi, gotSn)
			}
		}
	}

	boundary := []struct {
		name    string
		packed  uint8
		wantKsi uint8
		wantSn  uint8
	}{
		{"all zero", 0x00, 0, 0},
		{"max KSI and SN", 0x7F, 7, 15},
		{"spare high bit ignored", 0x80, 0, 0},
		{"SN field only", 0xF8, 0, 15},
		{"KSI field only", 0x07, 7, 0},
	}
	for _, tc := range boundary {
		t.Run(tc.name, func(t *testing.T) {
			ksi, sn := ParseeKSI(tc.packed)
			if ksi != tc.wantKsi || sn != tc.wantSn {
				t.Errorf("ParseeKSI(0x%02X) = (%d,%d), want (%d,%d)",
					tc.packed, ksi, sn, tc.wantKsi, tc.wantSn)
			}
		})
	}

	t.Run("packing masks inputs", func(t *testing.T) {
		if got := GeteKSI(0xFF, 0xFF); got != 0x7F {
			t.Errorf("GeteKSI(0xFF,0xFF) = 0x%02X, want 0x7F", got)
		}
		if got := GeteKSI(1, 2); got != 0x11 {
			t.Errorf("GeteKSI(1,2) = 0x%02X, want 0x11", got)
		}
	})
}

// TestIncrementCounts verifies that the uplink and downlink NAS counters
// advance independently and wrap at the 32-bit boundary.
func TestIncrementCounts(t *testing.T) {
	ctx := &EMMSecurityContext{}
	for i := 0; i < 3; i++ {
		ctx.IncrementUplinkCount()
	}
	for i := 0; i < 5; i++ {
		ctx.IncrementDownlinkCount()
	}
	if ctx.UplinkNASCount != 3 || ctx.DownlinkNASCount != 5 {
		t.Errorf("counts = (%d,%d), want (3,5)", ctx.UplinkNASCount, ctx.DownlinkNASCount)
	}

	ctx.UplinkNASCount = math.MaxUint32
	ctx.IncrementUplinkCount()
	if ctx.UplinkNASCount != 0 {
		t.Errorf("uplink count = %d after wrapping MaxUint32, want 0", ctx.UplinkNASCount)
	}
	if ctx.DownlinkNASCount != 5 {
		t.Errorf("downlink count = %d, want 5 (unaffected by uplink wrap)", ctx.DownlinkNASCount)
	}
}

// TestSecurityErrors verifies the error paths of the ciphering and
// integrity routing: malformed keys and unsupported algorithm identifiers.
func TestSecurityErrors(t *testing.T) {
	payload := []byte{0x01, 0x02, 0x03}

	t.Run("EEA1 short CK", func(t *testing.T) {
		ctx := &EMMSecurityContext{CK: make([]byte, 15), EncAlgo: EEA1, Valid: true}
		if _, err := ctx.Encrypt(1, 1, payload); !errors.Is(err, ErrInvalidKey) {
			t.Errorf("Encrypt error = %v, want %v", err, ErrInvalidKey)
		}
	})

	t.Run("EEA1 nil CK", func(t *testing.T) {
		ctx := &EMMSecurityContext{EncAlgo: EEA1, Valid: true}
		if _, err := ctx.Encrypt(1, 1, payload); !errors.Is(err, ErrInvalidKey) {
			t.Errorf("Encrypt error = %v, want %v", err, ErrInvalidKey)
		}
	})

	t.Run("EEA2 short CK", func(t *testing.T) {
		ctx := &EMMSecurityContext{CK: make([]byte, 15), EncAlgo: EEA2, Valid: true}
		if _, err := ctx.Encrypt(1, 1, payload); !errors.Is(err, ErrInvalidKey) {
			t.Errorf("Encrypt error = %v, want %v", err, ErrInvalidKey)
		}
	})

	t.Run("EIA1 short IK", func(t *testing.T) {
		ctx := &EMMSecurityContext{IK: make([]byte, 15), IntAlgo: EIA1, Valid: true}
		if _, err := ctx.ComputeMACI(1, 1, KNASDirectionDownlink, payload); !errors.Is(err, ErrInvalidKey) {
			t.Errorf("ComputeMACI error = %v, want %v", err, ErrInvalidKey)
		}
	})

	t.Run("EIA2 short IK", func(t *testing.T) {
		ctx := &EMMSecurityContext{IK: make([]byte, 15), IntAlgo: EIA2, Valid: true}
		if _, err := ctx.ComputeMACI(1, 1, KNASDirectionDownlink, payload); !errors.Is(err, ErrInvalidKey) {
			t.Errorf("ComputeMACI error = %v, want %v", err, ErrInvalidKey)
		}
	})

	t.Run("unsupported ciphering algorithm", func(t *testing.T) {
		ctx := &EMMSecurityContext{CK: make([]byte, 16), EncAlgo: ie.NASEncryptionAlgorithm128EEA3, Valid: true}
		if _, err := ctx.Encrypt(1, 1, payload); !errors.Is(err, ErrUnsupportedAlgo) {
			t.Errorf("Encrypt error = %v, want %v", err, ErrUnsupportedAlgo)
		}
	})

	t.Run("unsupported integrity algorithm", func(t *testing.T) {
		ctx := &EMMSecurityContext{IK: make([]byte, 16), IntAlgo: ie.NASIntegrityAlgorithm128EIA3, Valid: true}
		if _, err := ctx.ComputeMACI(1, 1, KNASDirectionDownlink, payload); !errors.Is(err, ErrUnsupportedAlgo) {
			t.Errorf("ComputeMACI error = %v, want %v", err, ErrUnsupportedAlgo)
		}
	})
}
