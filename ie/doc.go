// Package ie encodes and decodes LTE NAS information element envelopes defined
// in 3GPP TS 24.301 V18.9.0, clause 9.9, using the categories of TS 24.007.
//
// Every public Kind constant is an actual on-wire IEI. A tagged constructor,
// such as NewT3412ExtendedValue, sets IE.Kind to its expected tag. A positional
// constructor has a NewPositional prefix and returns an IE with Kind == NoTag.
// This keeps wire tags visible to callers without asking them to select a tag
// that might be incompatible with an IE value.
//
// Parse decodes an IEI with a globally unambiguous wire envelope. When a tag has
// different formats in different messages, create the matching tagged IE with
// its generated constructor and call UnmarshalBinary. ParseWithoutTag accepts a
// positional constructor result as its prototype. ParseMultiIEs scans only tagged
// optional IEs and returns their on-wire Kind; message-specific meaning and all
// positional IE parsing belong to the message parser.
//
// Unknown non-comprehension-required IEs are skipped according to TS 24.007
// EPS rules; unknown comprehension-required IEs return ErrUnknownIE. Colliding
// known wire formats return ErrAmbiguousIE rather than guessing a boundary.
// Typed decoding retains type-4 and type-6 extension octets beyond the
// currently-known value length, while Validate and Marshal reject them.
//
// Positional type-1 IEs share octets. MarshalHalfOctets and ParseHalfOctets handle
// a low/high pair. MarshalWithoutTag on a single type-1 IE uses one storage byte
// with its value in the low nibble, not a complete on-wire octet assignment.
//
// Payloads are owned copies and remain accessible, in the style of go-gtp.
// Compound values default to their TS-defined value octets; typed helpers are
// introduced one IE at a time without removing the raw constructors. For example,
// QualityOfServiceFields and IE.QualityOfService provide the semantic view of the
// Quality of service IE. Envelope validation is not a complete message or procedure
// conformance check.
package ie
