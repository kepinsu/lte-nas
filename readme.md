# LTE NAS

A Go library for LTE NAS Information Element (IE) envelopes, inspired by
`go-gtp`. It has no runtime dependencies.

## Scope

The [`ie`](ie/) package covers the envelopes of all 129 definitions in 3GPP
TS 24.301 v18.9.0 (Release 18), clause 9.9: V/TV, LV/TLV, and LV-E/TLV-E.
Raw values remain accessible in `IE.Payload`. Typed value models are added one
IE at a time, beginning with Quality of service; NAS message presence rules and
procedure-level validation remain outside this package's scope.

- Usage-specific constructors for every tagged IE.
- Explicit constructors for positional IEs.
- `Marshal`, `MarshalTo`, `MarshalLen`, and the standard binary interfaces.
- Parsing for individual IEs and optional IE sequences.
- Validation of lengths, tags, and buffer bounds.
- Packing of two positional type-1 IEs into a shared octet.

EMM and ESM causes use the existing `MmCause` and `SmCause` enumerations.

## Tags and constructors

`Kind` is the actual IEI written on the wire. Public constants are therefore
NAS wire values, not internal definition identifiers:

```go
ie.T3412ExtendedValue == 0x5e
ie.FullNetworkName == 0x43
ie.ESMCause == 0x58
```

A usage constructor selects the IEI itself. Its name identifies the usage so a
caller cannot attach an incompatible tag to a value:

```go
timer := ie.NewT3346Value(0x21)
// timer.Kind == ie.T3346Value == 0x5f

wire, err := timer.Marshal()
// wire == []byte{0x5f, 0x01, 0x21}
```

One definition can be reused under several tags. In that case, constructors are
distinct, for example `NewT3412Value`, `NewT3423Value`, `NewT3442Value`,
`NewT3412ExtendedValue`, `NewFullNetworkName`, and `NewShortNetworkName`.

Some tags have different meanings depending on the message. Constants can
therefore be aliases for the same wire value: `ie.ESMCause` and
`ie.UENetworkCapability` both equal `0x58`, while their constructors select
different internal formats and validation rules.

## Typed IE values

The raw constructors are kept for compatibility. When an IE receives a typed
model, its implementation lives in its own Go file and provides:

- a `…Fields` structure that exposes its semantic fields;
- a `New…WithFields` constructor that returns an `*IE` and any validation
  error; and
- an `(*IE).…` accessor and `Parse…Fields` function for decoding.

For example, [`ie/quality_of_service.go`](ie/quality_of_service.go) translates
the Quality of service payload defined in TS 24.008 clause 10.5.6.5. Bit rates
are expressed in bit/s, the bit-rate extensions are selected automatically,
and the legacy raw constructor remains valid:

```go
fields := &ie.QualityOfServiceFields{
    TrafficClass:                ie.QualityOfServiceTrafficClassInteractive,
    MaximumSDUSize:              1500,
    MaximumBitRateForUplink:     ie.QualityOfServiceBitRate{BitsPerSecond: 64_000},
    MaximumBitRateForDownlink:   ie.QualityOfServiceBitRate{BitsPerSecond: 576_000},
    TransferDelay:               200 * time.Millisecond,
    GuaranteedBitRateForUplink:  ie.QualityOfServiceBitRate{BitsPerSecond: 8_000},
    GuaranteedBitRateForDownlink: ie.QualityOfServiceBitRate{BitsPerSecond: 8_000},
}
qos, err := ie.NewQualityOfServiceWithFields(fields)
// qos is a normal *ie.IE and can be marshaled in the usual way.

decoded, err := qos.QualityOfService()

// Existing callers can still supply the complete raw value.
raw := ie.NewQualityOfService(payload)
```

## Positional IEs

A positional IE has no IEI: its `Kind` is `ie.NoTag`. Constructors use the
`NewPositional` prefix to make this distinction visible at the call site.

```go
mandatory := ie.NewPositionalGPRSTimer2(0x21)
wire, err := mandatory.Marshal()
// wire == []byte{0x01, 0x21} // LV, without a tag

decoded, err := ie.ParseWithoutTag(wire, ie.NewPositionalGPRSTimer2(0))
```

`MarshalWithoutTag` also encodes a tagged instance in positional form without
changing its `Kind`. Conversely, `UnmarshalWithoutTag` resets `Kind` to `NoTag`.
Positional type-1 values are stored in the low nibble; use `MarshalHalfOctets`
and `ParseHalfOctets` for a low/high pair.

```go
wire, err := ie.MarshalHalfOctets(
    ie.NewPositionalRequestType(1),
    ie.NewPositionalPDNType(3),
)
// wire == []byte{0x31}

request, pdn, err := ie.ParseHalfOctets(
    wire,
    ie.NewPositionalRequestType(0),
    ie.NewPositionalPDNType(0),
)
```

## Tagged decoding

A tagged constructor can be used as a decoding prototype and enforces its IEI:

```go
timer := ie.NewT3346Value(0)
if err := timer.UnmarshalBinary([]byte{0x5f, 0x01, 0x21}); err != nil {
    return err
}
```

`Parse(buffer, kind)` is available when a tag has a globally unambiguous wire
envelope. For tags whose format varies between messages (`0x32`, `0x58`, `0x5b`,
and `0x5c`), `Parse` returns `ErrAmbiguousIE`; create the expected usage with its
constructor and call `UnmarshalBinary` instead.

`ParseMultiIEs` receives the tagged optional part of a message. It preserves
order and repetitions and returns IEs with their on-wire `Kind`:

```go
ies, err := ie.ParseMultiIEs(optionalPart)
// ies[0].Kind is the received IEI, for example ie.EPSBearerContextStatus (0x57).
```

Unknown non-comprehension-required tags are skipped according to the TS 24.007
EPS rules. Comprehension-required tags (`00..0F` and `7E..7F`) return
`ErrUnknownIE`; truncated IEs return `io.ErrUnexpectedEOF`. A known tag whose
envelope depends on the message returns `ErrAmbiguousIE`.

`Length` is informational: encoding always recalculates the length from
`Payload`. A decoding failure does not modify the IE, and a `MarshalTo` failure
does not modify the destination buffer. Errors can be inspected with `errors.Is`.

When a typed decoder receives a type-4 or type-6 IE whose length exceeds the
current catalogue bound, it retains the complete value so the next IE remains
aligned (TS 24.007 section 11.4.2). `Validate` and marshal methods reject that
IE: callers must not emit an extension that this release does not define.

## Generation

Format definitions remain internal. The generation sources are under
[`ie/internal/iegen/`](ie/internal/iegen/):

```text
TS 24.301 (PDF or layout-preserving text)
  + catalogue_overrides.json
  + tagged_usage_overrides.json
  → extract_catalogue.py
  → catalogue.csv + tagged_usages.csv
  → go generate ./ie
  → catalogue_generated.go, constructors_generated.go,
    constructors_generated_test.go, and tagged_generated.go
```

`extract_catalogue.py` produces both CSV files deterministically for a given
release. `tagged_usage_overrides.json` only names the few usages where one
definition is used under several tags; it contains neither formats nor wire
values. The Go generator derives formats and ambiguous tags from `catalogue.csv`
and `tagged_usages.csv`, without a third registry.

```sh
go generate ./ie
python3 -m unittest discover -s ie/internal/iegen -p 'test_*.py'
go test -race ./...
go vet ./...
```

Sources:

- [TS 24.301 v18.9.0, clause 9.9 and tables 8.2/8.3](https://www.etsi.org/deliver/etsi_ts/124300_124399/124301/18.09.00_60/ts_124301v180900p.pdf)
- [TS 24.008 v18.7.0, clause 10.5.6.5](https://www.etsi.org/deliver/etsi_ts/124000_124099/124008/18.07.00_60/ts_124008v180700p.pdf)
- [TS 24.007 v18.1.0, clauses 11.2.4–11.2.5](https://www.etsi.org/deliver/etsi_ts/124000_124099/124007/18.01.00_60/ts_124007v180100p.pdf)
- [go-gtp](https://github.com/wmnsk/go-gtp/tree/main/gtpv2/ie)
