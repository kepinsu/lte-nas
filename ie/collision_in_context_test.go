package ie

import (
	"testing"
)

func TestCollisionTagsInContext(t *testing.T) {
	// Wire bytes: LLC SAPI (0x32, TLV format with payload [3]) + ESM Cause (0x58, TLV with payload [0x16])
	combined := []byte{0x32, 0x01, 0x03, 0x58, 0x01, 0x16}
	
	t.Logf("Combined wire: %v", combined)
	
	result, err := ParseMultiIEs(combined)
	if err != nil {
		t.Fatalf("ParseMultiIEs error: %v", err)
	}
	for i, r := range result {
		t.Logf("IE %d: kind=%#02x format=%d payload=%v len=%d", i, r.Kind, r.Format, r.Payload, r.MarshalLen())
	}
}
