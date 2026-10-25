package snow3g

import (
	"fmt"
	"testing"
)

// TestSnow3GLFSRTiming traces LFSR state after each warmup cycle to find the bug.
func TestSnow3GLFSRTiming(t *testing.T) {
	k := [4]uint32{0, 0, 0, 0}
	iv := [4]uint32{0, 0, 0, 0}

	var s Snow3G
	// Manual initialize to see initial state
	s.LFSR.S15 = k[3] ^ iv[0]
	s.LFSR.S14 = k[2]
	s.LFSR.S13 = k[1]
	s.LFSR.S12 = k[0] ^ iv[1]
	s.LFSR.S11 = k[3] ^ 0xffffffff
	s.LFSR.S10 = k[2] ^ 0xffffffff ^ iv[2]
	s.LFSR.S9 = k[1] ^ 0xffffffff ^ iv[3]
	s.LFSR.S8 = k[0] ^ 0xffffffff
	s.LFSR.S7 = k[3]
	s.LFSR.S6 = k[2]
	s.LFSR.S5 = k[1]
	s.LFSR.S4 = k[0]
	s.LFSR.S3 = k[3] ^ 0xffffffff
	s.LFSR.S2 = k[2] ^ 0xffffffff
	s.LFSR.S1 = k[1] ^ 0xffffffff
	s.LFSR.S0 = k[0] ^ 0xffffffff
	s.FSM.R1 = 0
	s.FSM.R2 = 0
	s.FSM.R3 = 0

	fmt.Println("=== Pre-warmup state ===")
	for i := 0; i < 16; i++ {
		fmt.Printf("S%d=0x%08x ", i, getLFSRValue(&s, uint8(i)))
		if i == 7 || i == 15 {
			fmt.Println()
		}
	}

	// Trace warmup cycle 0,1,2,...31
	for i := 0; i < 32; i++ {
		f := s.clockFSM()
		s.clockLFSRInitializationMode(f)

		if i < 5 {
			fmt.Printf("\nAfter warmup cycle %d: FSM.F=0x%08x, S15=0x%08x\n", i, f, s.LFSR.S15)
			fmt.Printf("  S0=0x%08x S1=0x%08x S5=0x%08x S11=0x%08x\n",
				s.LFSR.S0, s.LFSR.S1, s.LFSR.S5, s.LFSR.S11)
		}
	}

	fmt.Println("\n=== Post-warmup state ===")
	for i := 0; i < 16; i++ {
		fmt.Printf("S%d=0x%08x ", i, getLFSRValue(&s, uint8(i)))
		if i == 7 || i == 15 {
			fmt.Println()
		}
	}

	// Now check: does Initialize() give the same state?
	var s2 Snow3G
	s2.Initialize(k, iv)
	for i := 0; i < 16; i++ {
		v1 := getLFSRValue(&s, uint8(i))
		v2 := getLFSRValue(&s2, uint8(i))
		if v1 != v2 {
			t.Errorf("S%d mismatch after manual vs Initialize(): 0x%08x vs 0x%08x",
				i, v1, v2)
		}
	}

	fmt.Println("\n=== Manual and Initialize() match: ", len(func() []uint32 {
		var result []uint32
		for i := 0; i < 16; i++ {
			if getLFSRValue(&s, uint8(i)) != getLFSRValue(&s2, uint8(i)) {
				result = append(result, 0)
			}
		}
		return result
	}()) == 0, " ===")

	t.Logf("timing trace complete")
}
