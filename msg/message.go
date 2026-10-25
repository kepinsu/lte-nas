package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/security"
)

// Message represents a NAS (Non-Access Stratum) message for EMM or ESM as
// defined in TS 24.301.
//
// Parse dispatch is not part of this interface
type Message interface {
	MarshalTo([]byte) error
	UnmarshalBinary(b []byte) error
	MarshalLen() int
	MessageType() uint8
	MessageTypeName() string
	IsRequest() bool
}

// eMMLike is the interface for EMM messages that embed an EmmHeader.
type eMMLike interface {
	Header() *EmmHeader
}

// ============================================================================
// Write helpers
// ============================================================================

// writeIE writes the tagged representation of a single IE into b at offset off,
// returning the number of bytes consumed. It skips nil IEs (returning 0 without
// advancing off) and delegates to i.MarshalTo for valid tagged IEs. Returns
// io.ErrShortBuffer if b is too short, or the wrapped error from i.MarshalTo on
// validation failure.
func writeIE(b []byte, off int, i *ie.IE) (int, error) {
	if i == nil {
		return 0, nil
	}

	// Normalize collision tags to FormatFour TLV before writing. Collision tags
	// (0x32, 0x58, 0x5b, 0x5c) must use TLV encoding in tagged-optional parts.
	// We write them as a copy to avoid mutating the caller's IE.
	if ie.IsAmbiguousTagged(i.Kind) {
		n := 2 + len(i.Payload) // tag + length byte + payload
		if len(b)-off < n {
			return 0, io.ErrShortBuffer
		}
		b[off] = byte(i.Kind)
		b[off+1] = byte(len(i.Payload))
		copy(b[off+2:], i.Payload)
		return n, nil
	}

	n := i.MarshalLen()
	if len(b)-off < n {
		return 0, io.ErrShortBuffer
	}
	if err := i.Validate(); err != nil {
		return 0, err
	}
	if err := i.MarshalTo(b[off:]); err != nil {
		return 0, err
	}
	return n, nil
}

// ============================================================================
// Parsing — dispatcher functions (Parse, parseEMM, parseESM)
// ============================================================================

// Parse dispatches the given bytes to the correct concrete message type based on
// the protocol discriminator (PD) and security header type (SHT) encoded in the
// first octet, following the go-pfcp pattern.
//
// For EMM messages (PD = 0x7) the message type is extracted from byte 1 (plain,
// SHT 0 or 0xC0) or byte 6 (integrity-protected, SHT 1..5). Invalid SHT values
// return a wrapped ie.ErrInvalidValue.
//
// For ESM messages (PD = 0x2) the message type is extracted from byte 2. The
// header is always 3 octets.
//
// Unknown message types within EMM/ESM and any non-EPS PD fall back to Generic.
func Parse(b []byte) (Message, error) {
	if len(b) < 1 {
		return nil, io.ErrUnexpectedEOF
	}

	pd := b[0] & 0x0f
	var m Message

	switch pd {
	case PdEmm:
		m = parseEMM(b)
		if m == nil {
			return nil, fmt.Errorf("%w: invalid security header type %#02x", ie.ErrInvalidValue, b[0]>>4)
		}
	case PdEsm:
		m = parseESM(b)
		if m == nil {
			return nil, io.ErrUnexpectedEOF
		}
	default:
		m = &Generic{}
	}

	if err := m.UnmarshalBinary(b); err != nil {
		return nil, fmt.Errorf("unmarshal %T: %w", m, err)
	}

	return m, nil
}

// parseEMM extracts the EMM message type from b based on the security header
// type (SHT) in the high nibble of b[0], dispatches to a concrete type, and
// returns it. Invalid SHT values produce a wrapped ie.ErrInvalidValue; too-short
// payloads produce io.ErrUnexpectedEOF.
func parseEMM(b []byte) Message {
	sht := b[0] >> 4

	if sht > 5 && sht != 0xC0 {
		return nil // invalid SHT — caller wraps in error
	}

	var msgType uint8

	if sht == SecurityHeaderTypePlain || sht == 0xC0 {
		if len(b) < 2 {
			return nil // too short for plain EMM
		}
		msgType = b[1]
	} else {
		if len(b) < 7 {
			return nil // too short for integrity-protected EMM
		}
		msgType = b[6]
	}

	switch msgType {
	case EMMAttachRequest:
		return &AttachRequest{}
	case EMMAttachAccept:
		return &AttachAccept{}
	case EMMAttachComplete:
		return &AttachComplete{}
	case EMMAttachReject:
		return &AttachReject{}
	case EMMDetachRequest:
		return &DetachRequest{}
	case EMMDetachAccept:
		return &DetachAccept{}
	case EMMTrackingAreaUpdateRequest:
		return &TrackingAreaUpdateRequest{}
	case EMMTrackingAreaUpdateAccept:
		return &TrackingAreaUpdateAccept{}
	case EMMTrackingAreaUpdateComplete:
		return &TrackingAreaUpdateComplete{}
	case EMMTrackingAreaUpdateReject:
		return &TrackingAreaUpdateReject{}
	case EMMExtendedServiceRequest:
		return &ExtendedServiceRequest{}
	case EMMControlPlaneServiceRequest:
		return &ControlPlaneServiceRequest{}
	case EMMServiceReject:
		return &ServiceReject{}
	case EMMServiceAccept:
		return &ServiceAccept{}
	case EMMGUTIReallocationCommand:
		return &GUTIReallocationCommand{}
	case EMMGUTIReallocationComplete:
		return &GUTIReallocationComplete{}
	case EMMAuthenticationRequest:
		return &AuthenticationRequest{}
	case EMMAuthenticationResponse:
		return &AuthenticationResponse{}
	case EMMAuthenticationReject:
		return &AuthenticationReject{}
	case EMMAuthenticationFailure:
		return &AuthenticationFailure{}
	case EMMIdentityRequest:
		return &IdentityRequest{}
	case EMMIdentityResponse:
		return &IdentityResponse{}
	case EMMSecurityModeCommand:
		return &SecurityModeCommand{}
	case EMMSecurityModeComplete:
		return &SecurityModeComplete{}
	case EMMSecurityModeReject:
		return &SecurityModeReject{}
	case EMMStatus:
		return &EmmStatus{}
	case EMMInformation:
		return &EmmInformation{}
	case EMMDownlinkNASTransport:
		return &DownlinkNASTransport{}
	case EMMUplinkNASTransport:
		return &UplinkNASTransport{}
	case EMMCSServiceNotification:
		return &CSServiceNotification{}
	case EMMDownlinkGenericNASTransport:
		return &DownlinkGenericNASTransport{}
	case EMMUplinkGenericNASTransport:
		return &UplinkGenericNASTransport{}
	default:
		return &Generic{}
	}
}

// parseESM extracts the ESM message type from byte 2 of b, dispatches to a
// concrete type, and returns it. The ESM header is always 3 octets.
func parseESM(b []byte) Message {
	if len(b) < 3 {
		return nil // too short for ESM header
	}

	msgType := b[2]

	switch msgType {
	case ESMActivateDefaultEPSBearerContextRequest:
		return &ActivateDefaultEPSBearerContextRequest{}
	case ESMActivateDefaultEPSBearerContextAccept:
		return &ActivateDefaultEPSBearerContextAccept{}
	case ESMActivateDefaultEPSBearerContextReject:
		return &ActivateDefaultEPSBearerContextReject{}
	case ESMActivateDedicatedEPSBearerContextRequest:
		return &ActivateDedicatedEPSBearerContextRequest{}
	case ESMActivateDedicatedEPSBearerContextAccept:
		return &ActivateDedicatedEPSBearerContextAccept{}
	case ESMActivateDedicatedEPSBearerContextReject:
		return &ActivateDedicatedEPSBearerContextReject{}
	case ESMModifyEPSBearerContextRequest:
		return &ModifyEPSBearerContextRequest{}
	case ESMModifyEPSBearerContextAccept:
		return &ModifyEPSBearerContextAccept{}
	case ESMModifyEPSBearerContextReject:
		return &ModifyEPSBearerContextReject{}
	case ESMDeactivateEPSBearerContextRequest:
		return &DeactivateEPSBearerContextRequest{}
	case ESMDeactivateEPSBearerContextAccept:
		return &DeactivateEPSBearerContextAccept{}
	case ESMPDNConnectivityRequest:
		return &PDNConnectivityRequest{}
	case ESMPDNConnectivityReject:
		return &PDNConnectivityReject{}
	case ESMPDNDisconnectRequest:
		return &PDNDisconnectRequest{}
	case ESMPDNDisconnectReject:
		return &PDNDisconnectReject{}
	case ESMBearerResourceAllocationRequest:
		return &BearerResourceAllocationRequest{}
	case ESMBearerResourceAllocationReject:
		return &BearerResourceAllocationReject{}
	case ESMBearerResourceModificationRequest:
		return &BearerResourceModificationRequest{}
	case ESMBearerResourceModificationReject:
		return &BearerResourceModificationReject{}
	case ESMInformationRequest:
		return &EsmInfoRequest{}
	case ESMInformationResponse:
		return &EsmInfoResponse{}
	case ESMNotification:
		return &Notification{}
	case ESMDummy:
		return &DummyMessage{}
	case ESMStatus:
		return &ESMStatusMsg{}
	case ESMRemoteUEReport:
		return &RemoteUEReport{}
	case ESMRemoteUEReportResponse:
		return &RemoteUEReportResponse{}
	case ESMDataTransport:
		return &EsmDataTransport{}
	default:
		return &Generic{}
	}
}

// ============================================================================
// Ciphering helpers (DecodeWithSecurity, EncodeWithSecurity)
// ============================================================================

// DecodeWithSecurity dispatches ciphered EMM bytes to the correct message type,
// decrypts the payload using ctx (SHT ≥ 0x02), verifies MAC-I, and returns the
// fully parsed Message. Plain messages (SHT = 0x00) bypass decryption and MAC
// verification, behaving like Parse() but returning nil for non-EMM protocols.
//
// Signature mirrors Parse() — returns (Message, error).
func DecodeWithSecurity(b []byte, ctx *security.EMMSecurityContext) (Message, error) {
	if len(b) < 2 {
		return nil, io.ErrUnexpectedEOF
	}

	sht := b[0] >> 4
	pd := b[0] & 0x0f

	if pd != PdEmm {
		return nil, fmt.Errorf("%w: expected EMM protocol discriminator", ie.ErrInvalidValue)
	}

	// Plain messages — no MAC-I, no ciphering.
	if sht == SecurityHeaderTypePlain || sht == 0xC0 {
		msg := parseEMM(b)
		if msg == nil {
			return nil, fmt.Errorf("unknown EMM message type 0x%02x", b[1])
		}
		if err := msg.UnmarshalBinary(b); err != nil {
			return nil, fmt.Errorf("unmarshal %T: %w", msg, err)
		}
		return msg, nil
	}

	// Integrity-protected header: 7 octets.
	if len(b) < 7 {
		return nil, io.ErrUnexpectedEOF
	}

	var macI [4]byte
	copy(macI[:], b[1:5])

	msg := parseEMM(b)
	if msg == nil {
		return nil, fmt.Errorf("unknown EMM message type 0x%02x", b[6])
	}

	// Extract plaintext payload (bytes after header).
	payload := b[7:]

	if ctx != nil {
		fresh := uint32(b[5] & 0xff) // lower 8 bits of NAS Count
		ok, err := verifyMACI(ctx, security.KNASDirectionUplink, fresh, payload, macI)
		if err != nil {
			return nil, fmt.Errorf("verify MAC-I: %w", err)
		}
		if !ok {
			return nil, fmt.Errorf("MAC-I verification failed")
		}

		// Decrypt payload if ciphering is active (SHT ≥ 0x02).
		if sht >= SecurityHeaderTypeIntegrityProtectedAndCiphered {
			decrypted, err := decryptPayload(ctx, payload)
			if err != nil {
				return nil, fmt.Errorf("decrypt: %w", err)
			}
			payload = decrypted
		}
	}

	if err := msg.UnmarshalBinary(append(b[:7], payload...)); err != nil {
		return nil, fmt.Errorf("unmarshal %T: %w", msg, err)
	}

	return msg, nil
}

// EncodeWithSecurity marshals an EMM message to its wire representation and,
// when ctx is provided and SHT indicates ciphering (≥ 0x02), encrypts the full
// payload, computes MAC-I over plaintext, and prepends a 7-octet integrity header.
// Plain messages (SHT = 0x00) are returned as-is with no modification.
func EncodeWithSecurity(m Message, ctx *security.EMMSecurityContext) ([]byte, error) {
	if m == nil {
		return nil, fmt.Errorf("message is nil")
	}

	// EMM messages embed EmmHeader; ESM messages do not.
	var sht uint8
	switch emm := m.(type) {
	case eMMLike:
		sht = emm.Header().SecurityHeaderType
	default:
		return nil, fmt.Errorf("EncodeWithSecurity requires an EMM message")
	}

	// Plain messages — just marshal.
	if sht == SecurityHeaderTypePlain || sht == 0xC0 {
		return marshalFull(m)
	}

	if ctx == nil {
		return marshalFull(m)
	}

	// Compute sizes.
	var totalHeaderLen int
	switch emm := m.(type) {
	case eMMLike:
		totalHeaderLen = emm.Header().MarshalLen()
	}
	if totalHeaderLen != 7 {
		return nil, fmt.Errorf("unsupported header length %d for ciphered message", totalHeaderLen)
	}

	// Marshal to get the plaintext payload.
	plaintext, err := marshalFull(m)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	// The payload is everything after the 7-octet header.
	payload := plaintext[7:]

	// Compute MAC-I over the PLAINTEXT payload (TS 24.301 §9.6).
	var headerSeq uint8
	switch emm := m.(type) {
	case eMMLike:
		headerSeq = emm.Header().SequenceNumber
	}
	fresh := uint32(headerSeq & 0xff)
	macI, err := ctx.ComputeMACI(ctx.DownlinkNASCount, fresh, security.KNASDirectionDownlink, payload)
	if err != nil {
		return nil, fmt.Errorf("compute MAC-I: %w", err)
	}

	// Encrypt the payload if ciphering is active.
	ciphered := payload
	if sht >= SecurityHeaderTypeIntegrityProtectedAndCiphered {
		ciphered, err = encryptPayload(ctx, payload)
		if err != nil {
			return nil, fmt.Errorf("encrypt: %w", err)
		}
	}

	// Build output: security header (with MAC-I) + ciphered/unciphered payload.
	out := make([]byte, 1+4+1+len(ciphered)) // SHT byte + MAC-I + SeqNum + payload
	out[0] = byte(sht<<4 | PdEmm)
	copy(out[1:5], macI[:])
	out[5] = headerSeq
	copy(out[6:], ciphered)

	return out, nil
}

// marshalFull returns the complete wire representation of m (header + payload).
func marshalFull(m Message) ([]byte, error) {
	n := m.MarshalLen()
	if n == 0 {
		return nil, fmt.Errorf("marshalLen returned 0 for %T", m)
	}
	buf := make([]byte, n)
	if err := m.MarshalTo(buf); err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	return buf, nil
}

// ============================================================================
// Ciphering primitives (decryptPayload, encryptPayload, verifyMACI)
// ============================================================================

// decryptPayload applies the ciphering algorithm selected in ctx to the
// given ciphertext buffer.  Returns a copy of the decrypted data or an error.
// If ctx is nil, data is returned unchanged (plain messages).
func decryptPayload(ctx *security.EMMSecurityContext, data []byte) ([]byte, error) {
	if ctx == nil || len(data) == 0 {
		return data, nil
	}

	decrypted, err := ctx.Decrypt(ctx.UplinkNASCount, 0, data)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}

	out := make([]byte, len(decrypted))
	copy(out, decrypted)
	return out, nil
}

// encryptPayload applies the ciphering algorithm selected in ctx to the given
// plaintext.  Returns a copy of the encrypted data or an error.
func encryptPayload(ctx *security.EMMSecurityContext, data []byte) ([]byte, error) {
	if ctx == nil || len(data) == 0 {
		return data, nil
	}

	encrypted, err := ctx.Encrypt(ctx.DownlinkNASCount, 0, data)
	if err != nil {
		return nil, fmt.Errorf("encrypt: %w", err)
	}

	out := make([]byte, len(encrypted))
	copy(out, encrypted)
	return out, nil
}

// verifyMACI verifies the 4-octet MAC-I against msg using ctx.IntAlgo and
// direction. Returns false when the check fails (without error).
func verifyMACI(ctx *security.EMMSecurityContext, direction uint8, fresh uint32, msg []byte, expectedMACI [4]byte) (bool, error) {
	if ctx == nil || len(msg) == 0 {
		return false, fmt.Errorf("context or msg empty")
	}

	macI, err := ctx.ComputeMACI(ctx.UplinkNASCount, fresh, direction, msg)
	if err != nil {
		return false, fmt.Errorf("compute MAC-I: %w", err)
	}

	for i := range expectedMACI {
		if macI[i] != expectedMACI[i] {
			return false, nil
		}
	}
	return true, nil
}

// encodeWithHeaderAssembles a ciphered NAS PDU from an already-encoded payload.
// It prepends a 7-octet security header (SHT + MAC-I + SeqNum) and appends
// the encrypted payload. Returns a wrapped ie.ErrInvalidValue or io.ErrShortBuffer
// on failure.
