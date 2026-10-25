package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/security"
)

// ===========================================================================
// Fully-ciphered EMM messages: tagged-optional payload is decrypted in bulk
// when param != nil.  No positional mandatory IEs are ciphered — they sit in
// plaintext at the start of the payload (after the 7-octet header).
// ===========================================================================

// --- SecurityModeComplete (0x5E) -------------------------------------------
// No positional mandatory IEs — all tagged optional.  SHT plain (0x00) normally,
// but ciphered decoding handles SHT != 0x00 (integrity + ciphering).

// DecodeWithSecurity decodes a SECURITY MODE COMPLETE message whose tagged
// optional payload was fully ciphered.  param decrypts the entire IEs block
// using ctx.UplinkNASCount; when nil the payload is treated as plaintext.
func (msg *SecurityModeComplete) DecodeWithSecurity(param *security.EMMSecurityContext, b []byte) error {
	if msg == nil {
		return fmt.Errorf("%w: SecurityModeComplete is nil", ie.ErrNilIE)
	}

	var tmp SecurityModeComplete

	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	if len(pay) == 0 {
		return nil // only optional IE, ok if empty
	}

	// Decrypt the tagged-optional payload.
	if param != nil {
		pay, err = decryptPayload(param, pay)
		if err != nil {
			return fmt.Errorf("decipher security mode complete payload: %w", err)
		}
	}

	tagged, err := ie.ParseMultiIEs(pay)
	if err != nil {
		return fmt.Errorf("optional: %w", err)
	}

	for _, i := range tagged {
		switch i.Kind {
		case ie.MobileIdentity:
			if tmp.IMEISV == nil {
				tmp.IMEISV = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.ReplayedNASMessageContainer:
			if tmp.ReplayedNASMessageContainer == nil {
				tmp.ReplayedNASMessageContainer = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.UERadioCapabilityID:
			if tmp.UERadioCapabilityID == nil {
				tmp.UERadioCapabilityID = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.UERadioCapabilityIDDeletionIndication:
			if tmp.IEs == nil {
				tmp.IEs = append(tmp.IEs, i)
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		case ie.UECoarseLocationInformation:
			if tmp.UECoarseLocationInformation == nil {
				tmp.UECoarseLocationInformation = i
			} else {
				tmp.IEs = append(tmp.IEs, i)
			}
		default:
			tmp.IEs = append(tmp.IEs, i)
		}
	}

	*msg = tmp
	return nil
}

// EncodeWithSecurity encodes a ciphered SECURITY MODE COMPLETE message with a
// 7-octet security header (SHT + MAC-I + SeqNum) followed by the encrypted
// payload.  Returns a complete NAS PDU on success.
func (msg *SecurityModeComplete) EncodeWithSecurity(param *security.EMMSecurityContext) ([]byte, error) {
	if msg == nil {
		return nil, fmt.Errorf("%w: SecurityModeComplete is nil", ie.ErrNilIE)
	}

	payload := make([]byte, msg.MarshalLen())
	off := 0

	if msg.IMEISV != nil {
		n := msg.IMEISV.MarshalLen()
		if _, err := writeIE(payload[off:], 0, msg.IMEISV); err != nil {
			return nil, fmt.Errorf("IMEISV: %w", err)
		}
		off += n
	}
	if msg.ReplayedNASMessageContainer != nil {
		n := msg.ReplayedNASMessageContainer.MarshalLen()
		if _, err := writeIE(payload[off:], 0, msg.ReplayedNASMessageContainer); err != nil {
			return nil, fmt.Errorf("replayed NAS message container: %w", err)
		}
		off += n
	}
	if msg.UERadioCapabilityID != nil {
		n := msg.UERadioCapabilityID.MarshalLen()
		if _, err := writeIE(payload[off:], 0, msg.UERadioCapabilityID); err != nil {
			return nil, fmt.Errorf("UE radio capability ID: %w", err)
		}
		off += n
	}
	if msg.UECoarseLocationInformation != nil {
		n := msg.UECoarseLocationInformation.MarshalLen()
		if _, err := writeIE(payload[off:], 0, msg.UECoarseLocationInformation); err != nil {
			return nil, fmt.Errorf("UE coarse location information: %w", err)
		}
		off += n
	}

	for _, i := range msg.IEs {
		n := i.MarshalLen()
		if _, err := writeIE(payload[off:], 0, i); err != nil {
			return nil, fmt.Errorf("extra IE %d: %w", i.Kind, err)
		}
		off += n
	}

	ciphered, err := encryptPayload(param, payload)
	if err != nil {
		return nil, fmt.Errorf("encrypt: %w", err)
	}

	fresh := uint32(param.DownlinkNASCount & 0xff)
	macI, err := param.ComputeMACI(param.DownlinkNASCount, fresh, security.KNASDirectionDownlink, payload)
	if err != nil {
		return nil, fmt.Errorf("compute MAC-I: %w", err)
	}

	sht := byte(SecurityHeaderTypeIntegrityProtectedAndCiphered)
	pd := PdEmm
	hdr := make([]byte, 7)
	hdr[0] = byte(sht<<4 | pd)
	copy(hdr[1:5], macI[:])
	hdr[5] = byte(param.DownlinkNASCount & 0xff)

	return append(hdr, ciphered...), nil
}

// --- AuthenticationResponse (0x53) ----------------------------------------
// One positional mandatory IE (AuthenticationResponseParameter, LV type-4).
// No tagged optional IEs defined per spec table.

// DecodeWithSecurity decodes an AUTHENTICATION RESPONSE message whose payload
// (positional AuthenticationResponseParameter + any trailing unrecognized IEs)
// was fully ciphered.
func (msg *AuthenticationResponse) DecodeWithSecurity(param *security.EMMSecurityContext, b []byte) error {
	if msg == nil {
		return fmt.Errorf("%w: AuthenticationResponse is nil", ie.ErrNilIE)
	}

	var tmp AuthenticationResponse

	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	// Decrypt the payload.
	if param != nil {
		pay, err = decryptPayload(param, pay)
		if err != nil {
			return fmt.Errorf("decipher authentication response payload: %w", err)
		}
	}

	pos := 0

	tmp.AuthenticationResponseParameter, err = ie.ParseWithoutTag(
		pay[pos:], ie.NewPositionalAuthenticationResponseParameter(make([]byte, 5)),
	)
	if err != nil {
		return fmt.Errorf("AuthenticationResponseParameter: %w", err)
	}
	pos += tmp.AuthenticationResponseParameter.MarshalLenWithoutTag()

	// Parse trailing unrecognized IEs as tagged optional.
	if pos < len(pay) {
		tagged, err := ie.ParseMultiIEs(pay[pos:])
		if err != nil {
			return fmt.Errorf("optional trailing IEs: %w", err)
		}

		for _, i := range tagged {
			tmp.IEs = append(tmp.IEs, i)
		}
	}

	*msg = tmp
	return nil
}

// EncodeWithSecurity encodes a ciphered AUTHENTICATION RESPONSE message.
func (msg *AuthenticationResponse) EncodeWithSecurity(param *security.EMMSecurityContext) ([]byte, error) {
	if msg == nil {
		return nil, fmt.Errorf("%w: AuthenticationResponse is nil", ie.ErrNilIE)
	}

	payload := make([]byte, msg.MarshalLen())
	off := 0

	if msg.AuthenticationResponseParameter != nil {
		n := msg.AuthenticationResponseParameter.MarshalLenWithoutTag()
		if err := msg.AuthenticationResponseParameter.MarshalToWithoutTag(payload[off:]); err != nil {
			return nil, fmt.Errorf("AuthenticationResponseParameter: %w", err)
		}
		off += n
	}

	for _, i := range msg.IEs {
		n := i.MarshalLen()
		if _, err := writeIE(payload[off:], 0, i); err != nil {
			return nil, fmt.Errorf("extra IE %d: %w", i.Kind, err)
		}
		off += n
	}

	ciphered, err := encryptPayload(param, payload)
	if err != nil {
		return nil, fmt.Errorf("encrypt: %w", err)
	}

	fresh := uint32(param.DownlinkNASCount & 0xff)
	macI, err := param.ComputeMACI(param.DownlinkNASCount, fresh, security.KNASDirectionDownlink, payload)
	if err != nil {
		return nil, fmt.Errorf("compute MAC-I: %w", err)
	}

	sht := byte(SecurityHeaderTypeIntegrityProtectedAndCiphered)
	pd := PdEmm
	hdr := make([]byte, 7)
	hdr[0] = byte(sht<<4 | pd)
	copy(hdr[1:5], macI[:])
	hdr[5] = byte(param.DownlinkNASCount & 0xff)

	return append(hdr, ciphered...), nil
}
