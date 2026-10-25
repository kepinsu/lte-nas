package ie

import (
	"testing"
)

func TestCollisionDebugSkip(t *testing.T) {
	// [0x32, 0x01, 0x05] = LLC SAPI in TLV: tag=0x32, length=1, payload=[5]
	// [0x58, 0x01, 0x26] = ESM Cause in TLV: tag=0x58, length=1, payload=[0x26]
	b := []byte{0x32, 0x01, 0x05, 0x58, 0x01, 0x26}
	
	// Step 1: Check what skipUnknownIE returns for collision tag 0x32
	t.Logf("b[offset=0:] = %v, len=%d", b[0:], len(b))
	n, err := skipUnknownIE(b)
	t.Logf("skipUnknownIE: n=%d, err=%v", n, err)
	
	// Step 2: If offset advances to 2, check the next tag
	if err == nil && n > 0 {
		remaining := b[n:]
		t.Logf("b[offset=%d:] = %v, len=%d", n, remaining, len(remaining))
		
		tag2 := remaining[0]
		t.Logf("tag2 = %d (0x%02x), Kind=%d, is collision: %v", tag2, tag2, Kind(tag2), _ambiguousTaggedFormats[Kind(tag2)] != struct{}{})
		
		n2, err2 := skipUnknownIE(remaining)
		t.Logf("skipUnknownIE on remaining: n=%d, err=%v", n2, err2)
	}
	
	// Step 3: Full ParseMultiIEs trace
	t.Log("--- Full ParseMultiIEs ---")
	result, err := ParseMultiIEs(b)
	if err != nil {
		t.Fatalf("ParseMultiIEs error: %v", err)
	}
	for i, r := range result {
		t.Logf("IE %d: kind=%#02x format=%d payload=%v len=%d", i, r.Kind, r.Format, r.Payload, r.MarshalLen())
	}
}
