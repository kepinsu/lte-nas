package msg_test

import (
	"testing"

	"github.com/kepinsu/lte-nas/msg"
)

// FuzzParse ensures that msg.Parse never panics on crafted or malformed inputs.
// The seed corpus covers valid plain and integrity-protected EMM/ESM messages,
// as well as edge cases like empty input, truncated headers, unknown message
// types and non-EPS protocol discriminators.
func FuzzParse(f *testing.F) {
	// Valid plain EMM: ATTACH REQUEST with mandatory IEs.
	f.Add([]byte{
		0x07, 0x41, 0x39, // header (plain EMM + ATTACH REQUEST)
		0x82, 0x71, 0xC3, 0x14, 0xD0, 0x04, 0xD0, 0x01, // half-octets + pos IEs
	})

	// Valid integrity-protected EMM: Security Mode Complete (SHT=1, 7-octet header).
	f.Add([]byte{
		0x17, 0x5E, 0x11, 0x22, 0x33, 0x44, 0x05, // header: PD=7+SHT=1 + SMC
	})

	// Valid ESM: ACTIVATE DEFAULT EPS BEARER CONTEXT REQUEST (3-octet header).
	f.Add([]byte{
		0x21, 0x01, 0xC1, // ESM header + PTI=1
		0x72, 0x05, 0x0D, 0x64, 0x6F, 0x31, // QoS + APN
	})

	// Valid ESM: ESM STATUS.
	f.Add([]byte{
		0x21, 0x01, 0xE8, // ESM header + ESM STATUS
		0x01, // cause
	})

	// Malformed: empty input.
	f.Add([]byte{})

	// Malformed: single byte.
	f.Add([]byte{0x07})

	// Malformed: truncated MAC (only 2 of 4 MAC bytes in protected header).
	f.Add([]byte{0x17, 0x5E, 0x11, 0x22})

	// Malformed: truncated LV (length indicator but no value).
	f.Add([]byte{0x72, 0xFF})

	// Malformed: comprehension-required tag (IE type 0x00).
	f.Add([]byte{0x00, 0x10})

	// Malformed: type-6 unknown tag (tag 0x7E).
	f.Add([]byte{0x7E, 0x01, 0xAA})

	// Malformed: ambiguous tag 0x58 (ESM cause / UE NC).
	f.Add([]byte{0x58, 0x01, 0x01})

	// Edge: valid plain EMM but unknown message type.
	f.Add([]byte{0x07, 0xFF})

	// Edge: non-EPS PD.
	f.Add([]byte{0x3F, 0x01, 0x02})

	// Edge: valid SHT=0xC0 (legacy Service Request).
	f.Add([]byte{0xC7, 0x4D})

	f.Fuzz(func(_ *testing.T, payload []byte) {
		_, _ = msg.Parse(payload)
	})
}

// ---------------------------------------------------------------------------
// Targeted fuzzers: one per message type with meaningful seed corpus.
// ---------------------------------------------------------------------------

// FuzzAuthenticationResponseUnmarshal exercises the AuthenticationResponse decoder.
func FuzzAuthenticationResponseUnmarshal(f *testing.F) {
	// Valid plain header + 5-byte RES*.
	f.Add([]byte{
		0x07, 0x53, // plain EMM + AUTHENTICATION RESPONSE
		0x84, 0x05, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE,
	})

	// Valid integrity-protected header (7 octets) + payload.
	f.Add([]byte{
		0x17, 0x53, 0x11, 0x22, 0x33, 0x44, 0x05, // header
		0x84, 0x05, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, // payload
	})

	// Empty body (only header).
	f.Add([]byte{0x07, 0x53})

	// Truncated RES*.
	f.Add([]byte{0x07, 0x53, 0x84})
	f.Add([]byte{0x07, 0x53, 0x84, 0x03})

	// Random bytes that may contain IEs.
	f.Add([]byte{0x84, 0xFF, 0xAA, 0xBB, 0xCC})
	f.Add([]byte{0x7E, 0x01, 0xAA})

	// Empty.
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, payload []byte) {
		var m msg.AuthenticationResponse
		_ = m.UnmarshalBinary(payload)
	})
}

// FuzzActivateDefaultEPSBearerContextRequestUnmarshal exercises the
// ActivateDefaultEPSBearerContextRequest decoder.
func FuzzActivateDefaultEPSBearerContextRequestUnmarshal(f *testing.F) {
	// Valid header + EPS QoS (5 bytes).
	f.Add([]byte{
		0x21, 0x01, 0xC1, // ESM header
		0x72, 0x05, 0x0D, 0xAA, 0xBB, // EPS QoS
	})

	// Header only.
	f.Add([]byte{0x21, 0x01, 0xC1})

	// Partial QoS (only length byte).
	f.Add([]byte{0x21, 0x01, 0xC1, 0xFF})

	// Valid QoS + valid LV IE (tag 0x58, length 1, value).
	f.Add([]byte{
		0x21, 0x01, 0xC1,
		0x72, 0x05, 0x0D, 0xAA, 0xBB,
		0x58, 0x01, 0x01,
	})

	// Random trailing bytes.
	f.Add([]byte{0x72, 0xFF, 0xAA})
	f.Add([]byte{0xFF, 0x58, 0x01, 0x01})

	// Empty.
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, payload []byte) {
		var m msg.ActivateDefaultEPSBearerContextRequest
		_ = m.UnmarshalBinary(payload)
	})
}

// FuzzSecurityModeCompleteUnmarshal exercises the SecurityModeComplete decoder.
func FuzzSecurityModeCompleteUnmarshal(f *testing.F) {
	// Valid plain header (no positional IEs, only tagged optional).
	f.Add([]byte{0x07, 0x5E})

	// Header + IMEISV (tagged type-4 TLV).
	f.Add([]byte{0x07, 0x5E, 0x23, 0x06, 0xAA, 0xBB, 0xCC, 0xDD})

	// Header + replayed NAS message container (tagged type-6 TLV-E).
	f.Add([]byte{0x07, 0x5E, 0x79, 0x00, 0x04, 0xAA, 0xBB, 0xCC, 0xDD})

	// Header + UE radio capability ID (tagged type-1 TV).
	f.Add([]byte{0x07, 0x5E, 0x66, 0x10})

	// Truncated IE (length says more data than available).
	f.Add([]byte{0x23, 0xFF})

	// Type-1 TV with only tag (no payload for some formats).
	f.Add([]byte{0x23})

	// Random bytes.
	f.Add([]byte{0xFF, 0xAA, 0xBB})
	f.Add([]byte{0x79, 0xFF})

	// Empty.
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, payload []byte) {
		var m msg.SecurityModeComplete
		_ = m.UnmarshalBinary(payload)
	})
}

// ---------------------------------------------------------------------------
// Pre-computed seeds used by the fuzzers. These provide a stable, small set of
// test vectors that cover the main code paths without depending on any network
// capture. They are included here as convenience constants.
// ---------------------------------------------------------------------------

// SeedPlainAttachRequest is a minimal valid plain EMM ATTACH REQUEST.
var SeedPlainAttachRequest = []byte{0x07, 0x41, 0x39, 0xC3, 0x23, 0x05}

// SeedProtectedSMC is a minimal valid integrity-protected Security Mode Command.
var SeedProtectedSMC = []byte{0x17, 0x5D, 0xAA, 0xBB, 0xCC, 0xDD, 0x01}

// SeedActivateDefaultBearer is a minimal valid ESM ACTIVATE DEFAULT EPS BEARER CONTEXT REQUEST.
var SeedActivateDefaultBearer = []byte{0x21, 0x01, 0xC1, 0x72, 0x05}

// ---------------------------------------------------------------------------
