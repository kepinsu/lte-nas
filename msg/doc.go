// Package msg provides the LTE NAS (Non-Access Stratum) message architecture
// defined in 3GPP TS 24.301 V18.9.0.
//
// # Architecture
//
// The msg package decodes EPC NAS messages for EPS mobility management (EMM)
// and EPS session management (ESM). Each message type is represented by a type
// that implements the Message interface. The package layout is flat: one file per
// message type in msg.
//
// # Header layout
//
// EMM messages carry a security header type (SHT, TS 24.301 Table 9.3.1) and a
// protocol discriminator (PD) in the first octet. Plain messages (SHT=0x00) have
// a 2-octet header; integrity-protected messages (SHT 0x01..0x05 or 0xC0) have a
// 7-octet header including a MAC and sequence number. See TS 24.301 clauses 9.5
// and 9.6 for details.
//
// ESM messages carry a 3-octet header with bearer identity, protocol discriminator,
// procedure transaction identity (PTI), and message type (TS 24.301 clauses 9.3.2
// and 9.4).
//
// # Message types
//
// There are 59 message types defined in TS 24.301 Tables 9.8.1 (EMM, 32 types) and
// 9.8.2 (ESM, 27 types). Constants are defined in constants.go. EMM values have the
// top two bits set to 01 (xx01xxxx); ESM values have them set to 11 (xx11xxxx).
//
// # Decoding
//
// Unknown or unparseable messages are wrapped in Generic.
// The optional part of a message (information elements) is parsed by
// ie.ParseMultiIEs, which resolves ambiguous tags per-message context.

package msg
