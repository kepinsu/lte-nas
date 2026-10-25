# LTE NAS

A Go library for LTE Non-Access Stratum (NAS) protocol, implementing 3GPP TS 24.301 v18.9.0 (Release 18). It covers Information Element envelopes, NAS message parsing/encoding, and EMM security (integrity + ciphering). **Zero external dependencies.**

## Packages

| Package | Description |
|---------|-------------|
| [`ie`](#ie--information-elements) | All 129+ 3GPP TS 24.301 §9.9 IEs — V/TV, LV/TLV, LV-E/TLV-E |
| [`msg`](#msg--nas-messages) | EMM and ESM message types with dispatch, security headers, ciphering |
| [`security`](#security--emm-security) | EMM security context: integrity (UIA1/UIA2), ciphering (EEA1/EEA2), MAC-I |

---

## `ie` — Information Elements

The [`ie`](ie/) package provides raw and typed models for every IE in TS 24.301 v18.9.0 §9.9. Raw values are always accessible via `IE.Payload`.

### Wire formats

- **V / TV** — variable length with 1-octet length (format 3)
- **LV / TLV** — variable length with 1-octet length (format 4)
- **LV-E / TLV-E** — variable length with 2-octet BE length (format 6)

### Tags and constructors

`ie.Kind` is the wire IEI value. Usage constructors prevent cross-message misuse:

```go
timer := ie.NewT3346Value(0x21)
// timer.Kind == ie.T3346Value == 0x5f

wire, _ := timer.Marshal()
// wire == []byte{0x5f, 0x01, 0x21}
```

One definition may be reused under several tags. Each usage has its own constructor:
`NewT3412Value`, `NewT3423Value`, `NewT3412ExtendedValue`, `NewFullNetworkName`, `NewShortNetworkName`, etc.

### Typed value models

When an IE receives a typed model, it lives in its own file with:

- a `…Fields` struct exposing semantic fields;
- a `New…WithFields` constructor returning `*IE` and validation error;
- an `(*IE).…` accessor for decoding.

Example — QoS ([`ie/quality_of_service.go`](ie/quality_of_service.go)):

```go
fields := &ie.QualityOfServiceFields{
    TrafficClass:         ie.QualityOfServiceTrafficClassInteractive,
    MaximumSDUSize:       1500,
    MaximumBitRateForUplink: ie.QualityOfServiceBitRate{BitsPerSecond: 64_000},
    MaximumBitRateForDownlink: ie.QualityOfServiceBitRate{BitsPerSecond: 576_000},
    TransferDelay:       200 * time.Millisecond,
    GuaranteedBitRateForUplink: ie.QualityOfServiceBitRate{BitsPerSecond: 8_000},
    GuaranteedBitRateForDownlink: ie.QualityOfServiceBitRate{BitsPerSecond: 8_000},
}
qos, _ := ie.NewQualityOfServiceWithFields(fields)

// Decode back.
decoded, _ := qos.QualityOfService()
```

GPRS timers expose a single `time.Duration`:

```go
timer, _ := ie.NewT3346ValueWithFields(&ie.GPRSTimer2Fields{
    Duration: 10 * time.Minute,
})
```

### Positional IEs

Positional IEs have no IEI (`Kind == ie.NoTag`). Constructors use the `NewPositional` prefix:

```go
pos := ie.NewPositionalGPRSTimer2(0x21)
wire, _ := pos.Marshal()
// wire == []byte{0x01, 0x21} // LV without tag
```

Two positional type-1 IEs pack into a single octet:

```go
wire, _ := ie.MarshalHalfOctets(
    ie.NewPositionalRequestType(1),
    ie.NewPositionalPDNType(3),
)
// wire == []byte{0x31} // request=1, pdn=3
```

### Parsing

Parse individual IEs with `Parse` (globally unambiguous tags) or `UnmarshalBinary` (message-context tags like `0x32`, `0x58`, `0x5b`, `0x5c`):

```go
timer := ie.NewT3346Value(0)
timer.UnmarshalBinary([]byte{0x5f, 0x01, 0x21})
```

Parse tagged optional IE sequences:

```go
ies, err := ie.ParseMultiIEs(optionalPart)
// ies[0].Kind == ie.EPSBearerContextStatus (0x57)
```

Unknown non-comprehension-required tags pass through; comprehension-required tags (`0x00–0x0F`, `0x7E–0x7F`) return `ErrUnknownIE`.

### Code generation

Format definitions live in [`ie/internal/iegen/`](ie/internal/iegen/):

```
TS 24.301 v18.9.0 (PDF/text)
  + catalogue_overrides.json
  + tagged_usage_overrides.json
  → extract_catalogue.py
  → catalogue.csv + tagged_usages.csv
  → go generate ./ie
  → catalogue_generated.go, constructors_generated.go, …
```

```sh
go generate ./ie
python3 -m unittest discover -s ie/internal/iegen -p 'test_*.py'
go test -race ./...
```

---

## `msg` — NAS Messages

The [`msg`](msg/) package implements all EMM and ESM message types from TS 24.301 §8.2, with dispatch and EMM security (integrity + ciphering).

### Message interface

Every concrete message implements:

```go
type Message interface {
    MarshalTo([]byte) error
    UnmarshalBinary(b []byte) error
    MarshalLen() int
    MessageType() uint8
    MessageTypeName() string
    IsRequest() bool
    Header() *EmmHeader          // EMM messages only
}
```

EMM messages embed `EmmHeader` (security header type, PTI, message type, sequence number). ESM messages do not — use the `EMMLike` interface to access it.

### Plain parsing (no security)

```go
msg, err := msg.Parse(du)
// Returns concrete *AttachRequest, *AuthenticationRequest, *SecurityModeCommand, etc.

esm, err := msg.Parse(esmDU)
// Returns *PDNConnectivityRequest, *ModifyEPSBearerContextRequest, etc.
```

`Parse` dispatches based on the protocol discriminator (PD = `0x7` for EMM, `0x2` for ESM) and security header type (SHT).

### Security-aware dispatch

Two public functions handle integrity-protected and ciphered messages:

```go
// Decode a ciphered EMM message, verify MAC-I, decrypt payload, dispatch.
msg, err := msg.DecodeWithSecurity(cipheredDU, securityCtx)

// Marshal an EMM message, compute MAC-I over plaintext, encrypt payload, build 7-octet header.
wire, err := msg.EncodeWithSecurity(msg, securityCtx)
```

**Security Header Types (TS 24.301 §9.3.1):**

| SHT | Meaning | Cipher? | MAC-I? |
|-----|---------|---------|--------|
| `0x00` | Plain | No | No |
| `0x01` | Integrity protected | No | Yes |
| `0x02` | Integrity + ciphered | Yes | Yes |
| `0x03` | Integrity, new context | No | Yes |
| `0x04` | Integrity + ciphered, new context | Yes | Yes |
| `0x05` | Integrity + partial cipher | Partial | Yes |

- **SHT `0x00`** (plain): bypasses MAC-I and decryption, same as `Parse()`.
- **SHT `0x01` / `0x03`** (integrity only): MAC-I computed/verified, payload stays plaintext.
- **SHT `0x02` / `0x04`** (full ciphering): MAC-I over plaintext, payload encrypted before encoding / decrypted after decoding.
- **SHT `0x05`** (partial cipher): handled by message-specific code.

### Per-message security methods

Some messages provide dedicated `DecodeWithSecurity` / `EncodeWithSecurity` methods:

- `*SecurityModeCommand` — partially ciphered (MME network ID hashed into payload)
- `*SecurityModeComplete` — fully ciphered tagged-optional payload
- `*AuthenticationResponse` — fully ciphered tagged-optional payload

All other EMM messages use the generic dispatcher. ESM messages have no ciphering methods (they travel inside `DownlinkNASTransport` / `UplinkNASTransport`).

### Message catalogue (EMM)

| Msg type | Constant | Dir |
|----------|----------|-----|
| Attach Request | `EMMAttachRequest` (0x41) | UE→network |
| Attach Accept | `EMMAttachAccept` (0x42) | network→UE |
| Attach Complete | `EMMAttachComplete` (0x43) | UE→network |
| Attach Reject | `EMMAttachReject` (0x44) | network→UE |
| GUTI Reallocation Command | `EMMGUTIReallocationCommand` (0x52) | network→UE |
| GUTI Reallocation Complete | `EMMGUTIReallocationComplete` (0x53) | UE→network |
| Identity Request | `EMMIdentityRequest` (0x54) | network→UE |
| Identity Response | `EMMIdentityResponse` (0x55) | UE→network |
| Security Mode Command | `EMMSecurityModeCommand` (0x5a) | network→UE |
| Security Mode Complete | `EMMSecurityModeComplete` (0x5b) | UE→network |
| Security Mode Reject | `EMMSecurityModeReject` (0x5c) | UE→network |
| Configuration Update Command | `EMMConfigurationUpdateCommand` (0x5d) | network→UE |
| Configuration Update Complete | `EMMConfigurationUpdateComplete` (0x5e) | UE→network |
| Paging | `EMMPaging` (0x61) | network→UE |
| EMM Information | `EMMInformation` (0x71) | network→UE |
| Registration Request | `EMMRegistrationRequest` (0x72) | UE→network |
| Registration Accept | `EMMRegistrationAccept` (0x73) | network→UE |
| Registration Complete | `EMMRegistrationComplete` (0x74) | UE→network |
| Registration Reject | `EMMRegistrationReject` (0x75) | network→UE |
| Network Restart Notification | `EMMNetworkRestartNotification` (0x76) | network→UE |
| Network Restart Acknowledge | `EMMNetworkRestartAcknowledge` (0x77) | UE→network |
| Downlink NAS Transport | `EMMDownlinkNASTransport` (0x78) | network→UE |
| Uplink NAS Transport | `EMMUplinkNASTransport` (0x79) | UE→network |
| Paging Enhancements | `EMMPagingEnhancements` (0x7a) | network→UE |
| Network Information | `EMMNetworkInformation` (0x7b) | network→UE |
| Status (MM) | `EMMStatus` (0x7f) | both |
| Status (SM) | `ESMStatus` (0x7f) | both |

### Message catalogue (ESM)

| Msg type | Constant | Dir |
|----------|----------|-----|
| Activate Default EPS Bearer Context Request | `ESMActivateDefaultEPSBearerContextRequest` (0x45) | network→UE |
| Activate Default EPS Bearer Context Accept | `ESMActivateDefaultEPSBearerContextAccept` (0x46) | UE→network |
| Activate Default EPS Bearer Context Reject | `ESMActivateDefaultEPSBearerContextReject` (0x47) | UE→network |
| Activate Dedicated EPS Bearer Context Request | `ESMActivateDedicatedEPSBearerContextRequest` (0x48) | network→UE |
| Activate Dedicated EPS Bearer Context Accept | `ESMActivateDedicatedEPSBearerContextAccept` (0x49) | UE→network |
| Activate Dedicated EPS Bearer Context Reject | `ESMActivateDedicatedEPSBearerContextReject` (0x4a) | UE→network |
| Modify EPS Bearer Context Request | `ESMModifyEPSBearerContextRequest` (0x4b) | network→UE |
| Modify EPS Bearer Context Accept | `ESMModifyEPSBearerContextAccept` (0x4c) | UE→network |
| Modify EPS Bearer Context Reject | `ESMModifyEPSBearerContextReject` (0x4d) | UE→network |
| Deactivate EPS Bearer Context Request | `ESMDeactivateEPSBearerContextRequest` (0x4e) | network→UE |
| Deactivate EPS Bearer Context Accept | `ESMDeactivateEPSBearerContextAccept` (0x4f) | UE→network |
| Downlink EPS Session Management | `ESMDownlinkSessionManagement` (0x78) | network→UE |
| Uplink EPS Session Management | `ESMUplinkSessionManagement` (0x79) | UE→network |
| PDN Connectivity Request | `ESMPDNConnectivityRequest` (0x81) | UE→network |
| PDN Connectivity Accept | `ESMPDNConnectivityAccept` (0x82) | network→UE |
| PDN Connectivity Reject | `ESMPDNConnectivityReject` (0x83) | network→UE |
| PDN Disconnect Request | `ESMPDNDisconnectRequest` (0x84) | UE→network |
| PDN Disconnect Accept | `ESMPDNDisconnectAccept` (0x85) | network→UE |
| Bearer Resource Modification Request | `ESMBearerResourceModificationRequest` (0x87) | UE→network |
| Bearer Resource Allocation Request | `ESMBearerResourceAllocationRequest` (0x86) | network→UE |
| Bearer Resource Allocation Reject | `ESMBearerResourceAllocationReject` (0x88) | network→UE |
| ESM Information Request | `ESMInformationRequest` (0x89) | network→UE |
| ESM Information Response | `ESMInformationResponse` (0x8a) | UE→network |
| ESM Notification | `ESMNotification` (0x8b) | network→UE |
| ESM Status | `ESMStatus` (0x8c) | both |

---

## `security` — EMM Security

The [`security`](security/) package implements the 3GPP TS 33.401 EMM security functions: integrity (UIA1, UIA2), ciphering (EEA1, EEA2), and MAC-I.

### Security context

```go
type EMMSecurityContext struct {
    Algorithm  uint8       // KSI-EAI, EEA, UIA
    CipherKey  [16]byte    // CK (Snow3G: 4+28, AES: 16)
    IntegrityKey [16]byte  // IK (Snow3G: 4+28, AES: 16)
    UplinkNASCount   NASCount
    DownlinkNASCount NASCount
}
```

### Available algorithms

| Alg ID | Ciphering (EEA) | Integrity (UIA) |
|--------|-----------------|-----------------|
| `0x00` | EEA0 (null) | UIA1 (Snow3G-Uea1/Kit-1 / UEA2/AES) |
| `0x01` | EEA1 (Snow3G) | UIA2 (AES) |
| `0x02` | EEA2 (AES) | — |
| `0x03` | EEA3 (ZUC) | UIA3 (ZUC) |

### Core API

```go
// Encrypt: cipher the plaintext payload in Downlink direction.
encrypted, err := ctx.Encrypt(count, algo, plaintext)

// Decrypt: decipher the ciphertext payload in Uplink direction.
decrypted, err := ctx.Decrypt(count, algo, ciphertext)

// ComputeMACI: compute the 4-octet MAC-I over msg.
macI, err := ctx.ComputeMACI(count, freshCount, direction, msg)

// Verify: compare MAC-I bytes.
if !bytes.Equal(computedMACI, expectedMACI) { … }
```

### Algorithm backends

| Package | Implements | Key size |
|---------|-----------|----------|
| [`security/snow3g`](security/snow3g) | EEA1 / UIA1 (Snow3G) | 128-bit |
| [`security/myaes`](security/myaes) | EEA2 / UIA2 (AES-CTR) | 128-bit |

---

## Test coverage

```
$ go test ./... -cover
ok  github.com/kepinsu/lte-nas/ie            coverage: 73.1% of statements
ok  github.com/kepinsu/lte-nas/msg           coverage: 75.2% of statements
ok  github.com/kepinsu/lte-nas/security      coverage: 100.0% of statements
ok  github.com/kepinsu/lte-nas/security/myaes   coverage: 97.9% of statements
ok  github.com/kepinsu/lte-nas/security/snow3g  coverage: 100.0% of statements
```

---

## Build

```sh
go test -race ./...
go vet ./...
go generate ./ie    # (only when spec version changes)
```

---

## Standards references

- [TS 24.301 v18.9.0](https://www.etsi.org/deliver/etsi_ts/124300_124399/124301/18.09.00_60/) — EPS mobility management & session management
- [TS 24.007 v18.1.0](https://www.etsi.org/deliver/etsi_ts/124000_124099/124007/18.01.00_60/) — NAS protocol general aspects
- [TS 33.401 v18.2.0](https://www.etsi.org/deliver/etsi_ts/133300_133399/133401/18.02.00_60/) — 3GPP security architecture
