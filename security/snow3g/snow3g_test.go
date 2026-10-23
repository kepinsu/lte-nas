package snow3g

import (
	"fmt"
	"testing"
)

// TestSnow3GKeystream verifies SNOW 3G keystream against reference implementation output.
func TestSnow3GKeystream(t *testing.T) {
	tests := []struct {
		name   string
		k      [4]uint32
		iv     [4]uint32
		count  int
		expect []uint32
	}{
		{
			name:   "K=0x0,IV=0x0",
			k:      [4]uint32{0, 0, 0, 0},
			iv:     [4]uint32{0, 0, 0, 0},
			count:  16,
			expect: []uint32{0xc764a037, 0xb12fc857, 0xd470c3a5, 0xe24d982c, 0xbbfd2240, 0x178825b0, 0xe3a81869, 0x53b05c86, 0xdf405096, 0xb68c45c1, 0x889faf1c, 0x12cb24d6, 0xaac6e917, 0x735da28a, 0x06dfec55, 0xd98021a4},
		},
		{
			name:   "K=0x0,IV=0x1",
			k:      [4]uint32{0, 0, 0, 0},
			iv:     [4]uint32{1, 0, 0, 0},
			count:  16,
			expect: []uint32{0xddd86e74, 0x2b23508e, 0xa28991ac, 0x8e2b9cda, 0xd8acec47, 0x20173f80, 0xfb44f570, 0xa4429efe, 0xa40e0af1, 0x6bfa7858, 0xff9ed762, 0xe5f97304, 0x581cfeda, 0xfc1050cb, 0x316032b7, 0x817d0b21},
		},
		{
			name:   "K=0xDEADBEEF...,IV=0xFEDCBA98...",
			k:      [4]uint32{0xDEADBEEF, 0xC0DEFACE, 0x12345678, 0x9ABCDEF0},
			iv:     [4]uint32{0xFEDCBA98, 0x76543210, 0xBABEDCAD, 0xFACEFEED},
			count:  16,
			expect: []uint32{0x03750020, 0x967c87bf, 0x0af328d0, 0x1e68771d, 0xe8cc2c93, 0xb6763185, 0x8202a967, 0x2815494c, 0x0e206cd3, 0x29aa1168, 0xb850c02a, 0xcee0c344, 0x16afccba, 0x9e5c21a1, 0xd9cddd93, 0x2ffed4d0},
		},
		{
			name:   "K=1,2,3,4,IV=0x12345678",
			k:      [4]uint32{1, 2, 3, 4},
			iv:     [4]uint32{0x12345678, 0, 0, 0},
			count:  16,
			expect: []uint32{0x848adc63, 0x8c7ec338, 0x0091cacc, 0x02cbf351, 0x2bda03ba, 0xd4dcba41, 0xa4726fd2, 0xa9c7924f, 0x01b614cc, 0x78856734, 0xcaa50b3f, 0x91908f62, 0x38669b3a, 0x8fdbb00e, 0xe270bb43, 0x7572e576},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var s Snow3G
			s.Initialize(tc.k, tc.iv)

			got := s.GenerateKeystream(tc.count)

			if len(got) != tc.count {
				t.Fatalf("keystream length mismatch: got %d, want %d", len(got), tc.count)
			}

			for i := range got {
				if got[i] != tc.expect[i] {
					t.Errorf("keystream[%d]: got 0x%08x, want 0x%08x", i, got[i], tc.expect[i])
				}
			}
		})
	}
}

// TestSnow3GKeystreamSelfCheck verifies that NextKey() and GenerateKeystream
// produce the same output, and that consecutive calls are deterministic.
func TestSnow3GKeystreamSelfCheck(t *testing.T) {
	k := [4]uint32{0xDEADBEEF, 0xC0DEFACE, 0x12345678, 0x9ABCDEF0}
	iv := [4]uint32{0xFEDCBA98, 0x76543210, 0xBABEDCAD, 0xFACEFEED}

	// --- Test: NextKey and GenerateKeystream agree ---
	var s1 Snow3G
	s1.Initialize(k, iv)

	var s2 Snow3G
	s2.Initialize(k, iv)

	const n = 16
	gs := s2.GenerateKeystream(n)

	for i := 0; i < n; i++ {
		nk := s1.NextKey()
		if nk != gs[i] {
			t.Errorf("NextKey(%d) = 0x%08x, GenerateKeystream[%d] = 0x%08x",
				i, nk, i, gs[i])
		}
	}

	// --- Test: determinism (re-initialise and compare) ---
	var s3 Snow3G
	s3.Initialize(k, iv)
	gs2 := s3.GenerateKeystream(n)

	for i := 0; i < n; i++ {
		if gs[i] != gs2[i] {
			t.Errorf("non-deterministic keystream[%d]: 0x%08x vs 0x%08x",
				i, gs[i], gs2[i])
		}
	}

	// --- Test: consecutive keystream is not all zeros or constant ---
	allSame := true
	first := gs[0]
	for _, v := range gs {
		if v != first {
			allSame = false
			break
		}
	}
	if allSame {
		t.Errorf("keystream is constant (all 0x%08x) — cipher not mixing", first)
	}
}

// TestSnow3GFSMRegisters verifies FSM state changes on every clock.
func TestSnow3GFSMRegisters(t *testing.T) {
	k := [4]uint32{0xAAAAAAAA, 0x55555555, 0xFFFFFFFF, 0x00000000}
	iv := [4]uint32{0x12345678, 0x9ABCDEF0, 0xDEADBEEF, 0xFACEFACE}

	var s Snow3G
	s.Initialize(k, iv)

	r1a := s.FSM.R1
	r2a := s.FSM.R2
	r3a := s.FSM.R3

	// NextKey should advance FSM
	_ = s.NextKey()

	r1b := s.FSM.R1
	r2b := s.FSM.R2
	r3b := s.FSM.R3

	// With the key above the FSM registers will not all change identically,
	// but at least one should differ.
	if r1a == r1b && r2a == r2b && r3a == r3b {
		t.Errorf("FSM registers unchanged after NextKey")
	}

	fmt.Printf("FSM before: R1=0x%08x R2=0x%08x R3=0x%08x\n", r1a, r2a, r3a)
	fmt.Printf("FSM after:  R1=0x%08x R2=0x%08x R3=0x%08x\n", r1b, r2b, r3b)
}

// TestSnow3GLFSRRegisters verifies LFSR state changes on every clock.
func TestSnow3GLFSRRegisters(t *testing.T) {
	k := [4]uint32{0xDEADBEEF, 0xC0DEAD10, 0xBADF00D, 0xCAFEBABE}
	iv := [4]uint32{0xFEEDFACE, 0xBADC0DE, 0x11111111, 0x22222222}

	var s Snow3G
	s.Initialize(k, iv)

	// Snapshot LFSR after initialization (before any keystream clocks).
	lfsrInit := s.LFSR

	// Generate one word (NextKey performs extra clocks).
	_ = s.NextKey()

	lfsrAfter := s.LFSR

	// The LFSR must have shifted — at least some words changed.
	different := false
	for i := 0; i < 16 && !different; i++ {
		switch i {
		case 0:
			different = lfsrInit.S0 != lfsrAfter.S0
		case 1:
			different = lfsrInit.S1 != lfsrAfter.S1
		case 15:
			different = lfsrInit.S15 != lfsrAfter.S15
		}
	}

	if !different {
		t.Errorf("LFSR registers unchanged after NextKey")
	}

	fmt.Printf("LFSR.S0 init=0x%08x after=0x%08x\n", lfsrInit.S0, lfsrAfter.S0)
	fmt.Printf("LFSR.S15 init=0x%08x after=0x%08x\n", lfsrInit.S15, lfsrAfter.S15)
}

// TestSnow3GSBox verifies the S-box + mulx lookup tables.
func TestSnow3GSBox(t *testing.T) {
	// Verify SR and SQ boxes have 256 entries.
	if len(srBox) != 256 {
		t.Errorf("srBox has %d entries, want 256", len(srBox))
	}
	if len(sqBox) != 256 {
		t.Errorf("sqBox has %d entries, want 256", len(sqBox))
	}

	// Known SR/SQ values match the AES S-box for low indices.
	if srBox[0] != 0x63 {
		t.Errorf("srBox[0] = 0x%02x, want 0x63", srBox[0])
	}
	if srBox[1] != 0x7C {
		t.Errorf("srBox[1] = 0x%02x, want 0x7C", srBox[1])
	}
	if sqBox[0] != 0x25 {
		t.Errorf("sqBox[0] = 0x%02x, want 0x25", sqBox[0])
	}

	// Verify mulx arithmetic: mulx(0x63, 0x1b)
	// 0x63 << 1 = 0xC6, high bit was not set so no reduction.
	if got := mulx(0x63, 0x1b); got != 0xC6 {
		t.Errorf("mulx(0x63, 0x1b) = 0x%02x, want 0xC6", got)
	}
	// With high bit set: mulx(0x80, 0x1b) = (0x00 ^ 0x1b) = 0x1b
	if got := mulx(0x80, 0x1b); got != 0x1B {
		t.Errorf("mulx(0x80, 0x1b) = 0x%02x, want 0x1B", got)
	}

	// s1(0x00000000): all bytes → srBox[0]=0x63, then s1 mixing.
	s1zero := s1(0)
	t.Logf("s1(0x00000000) = 0x%08x", s1zero)

	// s2(0x00000000): all bytes → sqBox[0]=0x25
	mulx_25 := mulx(0x25, 0x69)
	t.Logf("mulx(0x25, 0x69) = 0x%02x", mulx_25)
	s2zero := s2(0)
	t.Logf("s2(0x00000000) = 0x%08x", s2zero)
}

// TestSnow3GInitializeLFSR checks that LFSR state after init matches expectations.
func TestSnow3GInitializeLFSR(t *testing.T) {
	k := [4]uint32{0x11111111, 0x22222222, 0x33333333, 0x44444444}
	iv := [4]uint32{0xAAAA, 0xBBBB, 0xCCCC, 0xDDDD}

	var s Snow3G
	s.Initialize(k, iv)

	// After 32 warmup cycles the LFSR state is pseudo-random.
	// Verify that some registers changed from initial load values.
	wasChanged := false
	for i, v := range []uint32{
		s.LFSR.S0, s.LFSR.S1, s.LFSR.S2, s.LFSR.S3,
		s.LFSR.S4, s.LFSR.S5, s.LFSR.S6, s.LFSR.S7,
		s.LFSR.S8, s.LFSR.S9, s.LFSR.S10, s.LFSR.S11,
		s.LFSR.S12, s.LFSR.S13, s.LFSR.S14, s.LFSR.S15,
	} {
		initial := []uint32{
			k[0] ^ 0xffffffff, k[1] ^ 0xffffffff, k[2] ^ 0xffffffff, k[3] ^ 0xffffffff,
			k[0], k[1], k[2], k[3],
			k[0] ^ 0xffffffff, k[1] ^ 0xffffffff ^ iv[3], k[2] ^ 0xffffffff ^ iv[2], k[3] ^ iv[0],
			k[0] ^ iv[1], k[1], k[2], k[3],
		}
		if v != initial[i] {
			wasChanged = true
			break
		}
	}

	if !wasChanged {
		t.Error("LFSR unchanged after warmup cycles — cipher not mixing")
	}

	fmt.Printf("LFSR.S0=0x%08x S1=0x%08x S15=0x%08x\n", s.LFSR.S0, s.LFSR.S1, s.LFSR.S15)
}
