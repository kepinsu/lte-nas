package ie

import (
	"testing"
)

func TestDRXWireFormat(t *testing.T) {
	drx := NewDRXParameter([]byte{0x78, 0x01})
	t.Logf("DRX: Kind=%#02x Format=%d Length=%d", drx.Kind, drx.Format, drx.Length)
	t.Logf("Payload: %v, len=%d", drx.Payload, len(drx.Payload))
	
	wire := make([]byte, drx.MarshalLen())
	_ = drx.MarshalTo(wire)
	t.Logf("Wire bytes: %v", wire)
	
	// Test ParseMultiIEs on just DRX wire bytes
	result, err := ParseMultiIEs(wire)
	if err != nil {
		t.Logf("ParseMultiIEs error: %v", err)
	} else {
		for _, r := range result {
			t.Logf("Parsed: kind=%#02x format=%d payload=%v", r.Kind, r.Format, r.Payload)
		}
	}
}
