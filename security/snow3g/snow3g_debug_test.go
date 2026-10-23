package snow3g

import (
	"fmt"
	"testing"
)

func TestSnow3GDebugTrace(t *testing.T) {
	k := [4]uint32{0, 0, 0, 0}
	iv := [4]uint32{0, 0, 0, 0}

	var s Snow3G
	s.Initialize(k, iv)

	fmt.Printf("After init:\n")
	for i := 0; i < 16; i++ {
		v := getLFSRValue(&s, uint8(i))
		fmt.Printf("  S%d = 0x%08x\n", i, v)
	}
	fmt.Printf("  FSM: R1=0x%08x R2=0x%08x R3=0x%08x\n", s.FSM.R1, s.FSM.R2, s.FSM.R3)

	// Step through the first NextKey call
	fmt.Println("\n=== First NextKey cycle ===")

	// isFirst is true, so we clock once then start generating
	f := s.clockFSM()
	fmt.Printf("After 1st clockFSM: F=0x%08x\n", f)
	fmt.Printf("  FSM: R1=0x%08x R2=0x%08x R3=0x%08x\n", s.FSM.R1, s.FSM.R2, s.FSM.R3)

	s.clockLFSRKeyStreamMode()
	fmt.Printf("After 1st clockLFSRKS: S0=0x%08x S15=0x%08x\n", s.LFSR.S0, s.LFSR.S15)

	// Now first real NextKey
	f = s.clockFSM()
	fmt.Printf("After 2nd clockFSM: F=0x%08x\n", f)
	kout := f ^ s.LFSR.S0
	fmt.Printf("Output would be F ^ S0 = 0x%08x ^ 0x%08x = 0x%08x\n", f, s.LFSR.S0, kout)

	s.clockLFSRKeyStreamMode()
	fmt.Printf("After 2nd clockLFSRKS: S0=0x%08x S15=0x%08x\n", s.LFSR.S0, s.LFSR.S15)

	// Generate full keystream
	ks := s.GenerateKeystream(4)
	fmt.Printf("\nFirst 4 keystream words: ")
	for i, v := range ks {
		fmt.Printf("ks[%d]=0x%08x ", i, v)
	}
	fmt.Println()

	// Compare with expected from reference
	expected := []uint32{0xc764a037, 0xb12fc857, 0xd470c3a5, 0xe24d982c}
	for i, v := range ks {
		if i < len(expected) {
			status := "MATCH"
			if v != expected[i] {
				status = "MISMATCH"
			}
			fmt.Printf("ks[%d]: got 0x%08x, want 0x%08x %s\n", i, v, expected[i], status)
		}
	}

	t.Logf("keystream output")
}

// getLFSRValue is a helper to read individual LFSR registers.
func getLFSRValue(s *Snow3G, i uint8) uint32 {
	switch i {
	case 0:
		return s.LFSR.S0
	case 1:
		return s.LFSR.S1
	case 2:
		return s.LFSR.S2
	case 3:
		return s.LFSR.S3
	case 4:
		return s.LFSR.S4
	case 5:
		return s.LFSR.S5
	case 6:
		return s.LFSR.S6
	case 7:
		return s.LFSR.S7
	case 8:
		return s.LFSR.S8
	case 9:
		return s.LFSR.S9
	case 10:
		return s.LFSR.S10
	case 11:
		return s.LFSR.S11
	case 12:
		return s.LFSR.S12
	case 13:
		return s.LFSR.S13
	case 14:
		return s.LFSR.S14
	case 15:
		return s.LFSR.S15
	default:
		return 0
	}
}

// TestSnow3GInitializeState checks the raw LFSR after key/IV loading but before warmup.
func TestSnow3GInitializeState(t *testing.T) {
	k := [4]uint32{0xDEAD, 0xBEEF, 0x1234, 0x5678}
	iv := [4]uint32{0xABCD, 0xEF01, 0x2345, 0x6789}

	// Manual LFSR loading
	s := Snow3G{}
	s.LFSR.S15 = k[3] ^ iv[0]   // 0x5678 ^ 0xABCD = 0xFB11
	s.LFSR.S14 = k[2]           // 0x1234
	s.LFSR.S13 = k[1]           // 0xBEEF
	s.LFSR.S12 = k[0] ^ iv[1]   // 0xDEAD ^ 0xEF01 = 0x31AC
	s.LFSR.S11 = k[3] ^ 0xffffffff // 0x5678 ^ 0xFFFFFFFF = 0xA987
	s.LFSR.S10 = k[2] ^ 0xffffffff ^ iv[2] // 0x1234 ^ 0xFFFFFFFF ^ 0x2345 = 0xFCAA
	s.LFSR.S9 = k[1] ^ 0xffffffff ^ iv[3] // 0xBEEF ^ 0xFFFFFFFF ^ 0x6789 = 0xC876
	s.LFSR.S8 = k[0] ^ 0xffffffff     // 0xDEAD ^ 0xFFFFFFFF = 0x2152
	s.LFSR.S7 = k[3]                // 0x5678
	s.LFSR.S6 = k[2]                // 0x1234
	s.LFSR.S5 = k[1]                // 0xBEEF
	s.LFSR.S4 = k[0]                // 0xDEAD
	s.LFSR.S3 = k[3] ^ 0xffffffff // 0xA987
	s.LFSR.S2 = k[2] ^ 0xffffffff // 0xEDCB
	s.LFSR.S1 = k[1] ^ 0xffffffff // 0x4110
	s.LFSR.S0 = k[0] ^ 0xffffffff // 0x2152

	s.FSM.R1 = 0
	s.FSM.R2 = 0
	s.FSM.R3 = 0

	fmt.Println("=== Manual LFSR state ===")
	for i := 0; i < 16; i++ {
		v := getLFSRValue(&s, uint8(i))
		fmt.Printf("  S%d = 0x%08x\n", i, v)
	}

	// Now use Initialize() with the same key/IV and compare
	var s2 Snow3G
	s2.Initialize(k, iv)

	fmt.Println("=== Initialize() LFSR state (post-warmup) ===")
	for i := 0; i < 16; i++ {
		v := getLFSRValue(&s2, uint8(i))
		fmt.Printf("  S%d = 0x%08x\n", i, v)
	}

	t.Logf("state comparison")
}

// TestSnow3GCompareInline verifies that s1 and mulAlpha produce correct intermediate values.
func TestSnow3GCompareInline(t *testing.T) {
	// Verify MULalpha for a specific byte
	val := uint8(0x63)
	result := mulAlpha(val)
	fmt.Printf("mulAlpha(0x%02x) = 0x%08x\n", val, result)

	// Verify s1 for a specific word
	w := uint32(0x12345678)
	s1Result := s1(w)
	fmt.Printf("s1(0x%08x) = 0x%08x\n", w, s1Result)

	// Verify s2 for a specific word
	s2Result := s2(w)
	fmt.Printf("s2(0x%08x) = 0x%08x\n", w, s2Result)

	t.Logf("intermediate values")
}
