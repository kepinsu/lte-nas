package msg

import (
	"fmt"
	"io"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/security"
)

// SecurityModeCommand represents a SECURITY MODE COMMAND message as defined in
// TS 24.301 clause 8.2.20 (Table 8.2.20.1, direction network to UE).
//
// The positional mandatory part encodes selected NAS security algorithms (1 octet,
// type-3), a spare half octet (type-1) and the replayed UE security capabilities
// (LV type-4). All remaining fields are tagged optional IEs in the order shown in
// TS 24.301 Table 8.2.20.1. This message is integrity-protected with a security
// header type of 0x03 (integrity protected with new EPS security context) or
// 0x04 (integrity protected and ciphered with new EPS security context), both
// yielding a 7-octet header per TS 24.301 clause 9.3.1 NOTE 1 and NOTE 2.
type SecurityModeCommand struct {
	EmmHeader
	SelectedNASSecurityAlgorithms *ie.IE // positional type-3 (9.9.3.23)
	SpareHalfOctet                *ie.IE // positional type-1 (9.9.2.9)
	ReplayedUESecurityCapability  *ie.IE // positional LV type-4 (9.9.3.36)
	// Tagged optional IEs resolved via prototypes in the spec table.
	IMEISVRequest                          *ie.IE // tagged optional TV type-1, 0xc0 (9.9.3.18)
	NonceUE                                *ie.IE // tagged optional TV type-3, 0x55 (9.9.3.25)
	NonceMME                               *ie.IE // tagged optional TV type-3, 0x56 (9.9.3.25)
	HashMME                                *ie.IE // tagged optional TLV, 0x4f (9.9.3.50)
	ReplayedUEAdditionalSecurityCapability *ie.IE // tagged optional TLV, 0x6f (9.9.3.53)
	UERadioCapabilityIDRequest             *ie.IE // tagged optional TLV, 0x37 (9.9.3.59)
	UEInformationRequest                   *ie.IE // tagged optional TV type-1, 0xd0 (9.9.3.71)
	IEs                                    []*ie.IE
}

// NewSecurityModeCommand creates a SECURITY MODE COMMAND message with the
// mandatory positional IEs from TS 24.301 clause 8.2.20, followed by optional
// tagged IEs classified by their Kind. Fields for recognized tags are resolved
// first (first occurrence wins); duplicates and unrecognized IEs are appended
// to the IEs field.
func NewSecurityModeCommand(selectedNASSecurityAlgorithms, spareHalfOctet,
	replayedUESecurityCapability *ie.IE, ies ...*ie.IE,
) *SecurityModeCommand {
	if selectedNASSecurityAlgorithms == nil || spareHalfOctet == nil ||
		replayedUESecurityCapability == nil {
		panic("SecurityModeCommand positional IEs must not be nil")
	}

	m := &SecurityModeCommand{
		SelectedNASSecurityAlgorithms: selectedNASSecurityAlgorithms,
		SpareHalfOctet:                spareHalfOctet,
		ReplayedUESecurityCapability:  replayedUESecurityCapability,
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Kind {
		case ie.IMEISVRequest:
			if m.IMEISVRequest == nil {
				m.IMEISVRequest = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.NonceUE:
			if m.NonceUE == nil {
				m.NonceUE = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.NonceMME:
			if m.NonceMME == nil {
				m.NonceMME = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.HashMME:
			if m.HashMME == nil {
				m.HashMME = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.UERadioCapabilityIDRequest:
			if m.UERadioCapabilityIDRequest == nil {
				m.UERadioCapabilityIDRequest = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.UEInformationRequest:
			if m.UEInformationRequest == nil {
				m.UEInformationRequest = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		case ie.UEAdditionalSecurityCapability:
			if m.ReplayedUEAdditionalSecurityCapability == nil {
				m.ReplayedUEAdditionalSecurityCapability = i
			} else {
				m.IEs = append(m.IEs, i)
			}
		default:
			m.IEs = append(m.IEs, i)
		}
	}

	return m
}

// MarshalLen returns the wire size of the SECURITY MODE COMMAND message, or 0 for
// a nil receiver. The result includes header, all positional mandatory IEs and
// all named tagged IEs in spec-table order.
func (m *SecurityModeCommand) MarshalLen() int {
	if m == nil {
		return 0
	}

	n := m.EmmHeader.MarshalLen()
	if n == 0 {
		return 0
	}

	if m.SelectedNASSecurityAlgorithms != nil {
		n += m.SelectedNASSecurityAlgorithms.MarshalLenWithoutTag()
	}
	if m.SpareHalfOctet != nil {
		n++
	}
	if m.ReplayedUESecurityCapability != nil {
		n += m.ReplayedUESecurityCapability.MarshalLenWithoutTag()
	}

	if m.IMEISVRequest != nil {
		n += m.IMEISVRequest.MarshalLen()
	}
	if m.NonceUE != nil {
		n += m.NonceUE.MarshalLen()
	}
	if m.NonceMME != nil {
		n += m.NonceMME.MarshalLen()
	}
	if m.HashMME != nil {
		n += m.HashMME.MarshalLen()
	}
	if m.ReplayedUEAdditionalSecurityCapability != nil {
		n += m.ReplayedUEAdditionalSecurityCapability.MarshalLen()
	}
	if m.UERadioCapabilityIDRequest != nil {
		n += m.UERadioCapabilityIDRequest.MarshalLen()
	}
	if m.UEInformationRequest != nil {
		n += m.UEInformationRequest.MarshalLen()
	}

	for _, i := range m.IEs {
		n += i.MarshalLen()
	}

	return n
}

// MarshalTo encodes the SECURITY MODE COMMAND message into b. The nil receiver
// returns a wrapped ie.ErrNilIE; io.ErrShortBuffer without writing if b is too
// short. Header, positional mandatory IEs (no tags) and tagged optional IEs in
// spec-table order follow TS 24.301 clause 8.2.20 table order.
func (m *SecurityModeCommand) MarshalTo(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: SecurityModeCommand is nil", ie.ErrNilIE)
	}

	hlen := m.EmmHeader.MarshalLen()
	if hlen == 0 {
		return fmt.Errorf("%w: EmmHeader MarshalLen returned 0", ie.ErrNilIE)
	}

	if len(b) < hlen {
		return io.ErrShortBuffer
	}

	if err := m.EmmHeader.MarshalTo(b); err != nil {
		return fmt.Errorf("header: %w", err)
	}

	b = b[hlen:]
	n := 0

	// Positional mandatory IEs without IEI tags.
	if m.SelectedNASSecurityAlgorithms != nil {
		if err := m.SelectedNASSecurityAlgorithms.MarshalToWithoutTag(b[n:]); err != nil {
			return fmt.Errorf("selected NAS security algorithms: %w", err)
		}
		n += m.SelectedNASSecurityAlgorithms.MarshalLenWithoutTag()
	}
	if m.SpareHalfOctet != nil {
		if err := m.SpareHalfOctet.MarshalToWithoutTag(b[n:]); err != nil {
			return fmt.Errorf("spare half octet: %w", err)
		}
		n++
	}
	if m.ReplayedUESecurityCapability != nil {
		n0 := n
		if err := m.ReplayedUESecurityCapability.MarshalToWithoutTag(b[n:]); err != nil {
			return fmt.Errorf("replayed UE security capability: %w", err)
		}
		n += m.ReplayedUESecurityCapability.MarshalLenWithoutTag()
		_ = n0
	}

	// Named tagged IEs in spec-table order.
	if m.IMEISVRequest != nil {
		w, err := writeIE(b, n, m.IMEISVRequest)
		if err != nil {
			return fmt.Errorf("IMEISV request: %w", err)
		}
		n += w
	}
	if m.NonceUE != nil {
		w, err := writeIE(b, n, m.NonceUE)
		if err != nil {
			return fmt.Errorf("nonce UE: %w", err)
		}
		n += w
	}
	if m.NonceMME != nil {
		w, err := writeIE(b, n, m.NonceMME)
		if err != nil {
			return fmt.Errorf("nonce MME: %w", err)
		}
		n += w
	}
	if m.HashMME != nil {
		w, err := writeIE(b, n, m.HashMME)
		if err != nil {
			return fmt.Errorf("HashMME: %w", err)
		}
		n += w
	}
	if m.ReplayedUEAdditionalSecurityCapability != nil {
		w, err := writeIE(b, n, m.ReplayedUEAdditionalSecurityCapability)
		if err != nil {
			return fmt.Errorf("replayed UE additional security capability: %w", err)
		}
		n += w
	}
	if m.UERadioCapabilityIDRequest != nil {
		w, err := writeIE(b, n, m.UERadioCapabilityIDRequest)
		if err != nil {
			return fmt.Errorf("UE radio capability ID request: %w", err)
		}
		n += w
	}
	if m.UEInformationRequest != nil {
		w, err := writeIE(b, n, m.UEInformationRequest)
		if err != nil {
			return fmt.Errorf("UE information request: %w", err)
		}
		n += w
	}

	// Extra tagged IEs.
	for _, i := range m.IEs {
		w, err := writeIE(b, n, i)
		if err != nil {
			return fmt.Errorf("extra IE: %w", err)
		}
		n += w
	}

	return nil
}

// UnmarshalBinary decodes a SECURITY MODE COMMAND message from b. The receiver
// is left unchanged on error. Header, positional mandatory IEs via
// ParseWithoutTag with NewPositional... prototypes, and tagged optional IEs
// resolved via prototypes for all recognized tags.
func (m *SecurityModeCommand) UnmarshalBinary(b []byte) error {
	if m == nil {
		return fmt.Errorf("%w: SecurityModeCommand is nil", ie.ErrNilIE)
	}

	var tmp SecurityModeCommand

	// Decode header (7 octets for SHT 0x03/0x04).
	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	pos := 0

	if len(pay) == 0 {
		return io.ErrUnexpectedEOF
	}

	// Positional mandatory: selected NAS security algorithms (1 byte, type-3).
	tmp.SelectedNASSecurityAlgorithms, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalNASSecurityAlgorithms(0))
	if err != nil {
		return fmt.Errorf("selected NAS security algorithms: %w", err)
	}
	pos += tmp.SelectedNASSecurityAlgorithms.MarshalLenWithoutTag()

	// Spare half octet (1 byte, type-1).
	if len(pay[pos:]) == 0 {
		return io.ErrUnexpectedEOF
	}
	tmp.SpareHalfOctet, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalSpareHalfOctet(0))
	if err != nil {
		return fmt.Errorf("spare half octet: %w", err)
	}
	pos++

	// Replayed UE security capabilities (LV, type-4).
	if len(pay[pos:]) == 0 {
		return io.ErrUnexpectedEOF
	}
	tmp.ReplayedUESecurityCapability, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalUESecurityCapability([]byte{0}))
	if err != nil {
		return fmt.Errorf("replayed UE security capability: %w", err)
	}
	pos += tmp.ReplayedUESecurityCapability.MarshalLenWithoutTag()

	// Tagged optional IEs — prototypes resolve all recognized tags.
	if len(pay[pos:]) > 0 {
		tagged, err := ie.ParseMultiIEs(pay[pos:])
		if err != nil {
			return fmt.Errorf("optional: %w", err)
		}

		for _, i := range tagged {
			switch i.Kind {
			case ie.IMEISVRequest:
				if tmp.IMEISVRequest == nil {
					tmp.IMEISVRequest = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case ie.NonceUE:
				if tmp.NonceUE == nil {
					tmp.NonceUE = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case ie.NonceMME:
				if tmp.NonceMME == nil {
					tmp.NonceMME = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case ie.HashMME:
				if tmp.HashMME == nil {
					tmp.HashMME = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case ie.UERadioCapabilityIDRequest:
				if tmp.UERadioCapabilityIDRequest == nil {
					tmp.UERadioCapabilityIDRequest = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case ie.UEInformationRequest:
				if tmp.UEInformationRequest == nil {
					tmp.UEInformationRequest = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case ie.UEAdditionalSecurityCapability:
				if tmp.ReplayedUEAdditionalSecurityCapability == nil {
					tmp.ReplayedUEAdditionalSecurityCapability = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			default:
				tmp.IEs = append(tmp.IEs, i)
			}
		}
	}

	*m = tmp
	return nil
}

// MessageType returns the SECURITY MODE COMMAND message type constant (0x5D).
func (m *SecurityModeCommand) MessageType() uint8 {
	return EMMSecurityModeCommand
}

// MessageTypeName returns the descriptive name from Table 9.8.1.
func (m *SecurityModeCommand) MessageTypeName() string {
	return "Security mode command"
}

// IsRequest returns false — SECURITY MODE COMMAND is sent network to UE.
func (m *SecurityModeCommand) IsRequest() bool {
	return false
}

// SecurityModeCommand implements the Message interface.
var _ Message = (*SecurityModeCommand)(nil)

// DecodeWithSecurity decodes a SECURITY MODE COMMAND message with per-IE
// partial ciphering (SHT 0x03/0x04).  When param is not nil, it decrypts
// the HashMME tagged-optional IE in place using ctx.EncAlgo + DownlinkNASCount.
// Positional mandatory IEs (SelectedNASSecurityAlgorithms, SpareHalfOctet,
// ReplayedUESecurityCapability) are never ciphered.
func (msg *SecurityModeCommand) DecodeWithSecurity(param *security.EMMSecurityContext, b []byte) error {
	if msg == nil {
		return fmt.Errorf("%w: SecurityModeCommand is nil", ie.ErrNilIE)
	}

	var tmp SecurityModeCommand

	pay, err := tmp.EmmHeader.UnmarshalBinary(b)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	pos := 0

	tmp.SelectedNASSecurityAlgorithms, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalNASSecurityAlgorithms(0))
	if err != nil {
		return fmt.Errorf("selected NAS security algorithms: %w", err)
	}
	pos += tmp.SelectedNASSecurityAlgorithms.MarshalLenWithoutTag()

	if len(pay[pos:]) == 0 {
		return io.ErrUnexpectedEOF
	}
	tmp.SpareHalfOctet, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalSpareHalfOctet(0))
	if err != nil {
		return fmt.Errorf("spare half octet: %w", err)
	}
	pos++

	if len(pay[pos:]) == 0 {
		return io.ErrUnexpectedEOF
	}
	tmp.ReplayedUESecurityCapability, err = ie.ParseWithoutTag(pay[pos:], ie.NewPositionalUESecurityCapability([]byte{0}))
	if err != nil {
		return fmt.Errorf("replayed UE security capability: %w", err)
	}
	pos += tmp.ReplayedUESecurityCapability.MarshalLenWithoutTag()

	// Parse tagged-optional IEs from payload[pos:].
	if len(pay) > pos {
		taggedPayload := pay[pos:]

		// Decrypt the tagged-optional payload if cipher parameters are provided.
		if param != nil {
			taggedPayload, err = decryptPayload(param, taggedPayload)
			if err != nil {
				return fmt.Errorf("decipher command payload: %w", err)
			}
		}

		tagged, err := ie.ParseMultiIEs(taggedPayload)
		if err != nil {
			return fmt.Errorf("optional: %w", err)
		}

		for _, i := range tagged {
			switch i.Kind {
			case ie.IMEISVRequest:
				if tmp.IMEISVRequest == nil {
					tmp.IMEISVRequest = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case ie.NonceUE:
				if tmp.NonceUE == nil {
					tmp.NonceUE = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case ie.NonceMME:
				if tmp.NonceMME == nil {
					tmp.NonceMME = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case ie.HashMME:
				if tmp.HashMME == nil {
					tmp.HashMME = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case ie.UERadioCapabilityIDRequest:
				if tmp.UERadioCapabilityIDRequest == nil {
					tmp.UERadioCapabilityIDRequest = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case ie.UEInformationRequest:
				if tmp.UEInformationRequest == nil {
					tmp.UEInformationRequest = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			case ie.UEAdditionalSecurityCapability:
				if tmp.ReplayedUEAdditionalSecurityCapability == nil {
					tmp.ReplayedUEAdditionalSecurityCapability = i
				} else {
					tmp.IEs = append(tmp.IEs, i)
				}
			default:
				tmp.IEs = append(tmp.IEs, i)
			}
		}
	}

	*msg = tmp
	return nil
}

// EncodeWithSecurity encodes a SECURITY MODE COMMAND message with integrity-
// protected and optionally ciphered tagged-optional IEs.  When param is not nil,
// the HashMME value (if present) is encrypted in place.  Returns a complete
// ciphered NAS PDU on success.
func (msg *SecurityModeCommand) EncodeWithSecurity(param *security.EMMSecurityContext) ([]byte, error) {
	if msg == nil {
		return nil, fmt.Errorf("%w: SecurityModeCommand is nil", ie.ErrNilIE)
	}

	if param == nil {
		// No ciphering: encode as plain and let the caller add the security header.
		return msg.encodePayload()
	}

	// Step 1: Calculate payload size (positional + tagged, HashMME ciphered).
	totalSize := msg.MarshalLen()

	// If HashMME is ciphered, add 1 extra byte (cipher overhead is typically 0
	// for stream ciphers, but account for it in the calculation).
	payload := make([]byte, totalSize)

	// Step 2: Encode positional mandatory IEs (plaintext).
	off := 0
	if msg.SelectedNASSecurityAlgorithms != nil {
		if err := msg.SelectedNASSecurityAlgorithms.MarshalToWithoutTag(payload[off:]); err != nil {
			return nil, fmt.Errorf("selected NAS security algorithms: %w", err)
		}
		off += msg.SelectedNASSecurityAlgorithms.MarshalLenWithoutTag()
	}
	if msg.SpareHalfOctet != nil {
		if err := msg.SpareHalfOctet.MarshalToWithoutTag(payload[off:]); err != nil {
			return nil, fmt.Errorf("spare half octet: %w", err)
		}
		off++
	}
	if msg.ReplayedUESecurityCapability != nil {
		if err := msg.ReplayedUESecurityCapability.MarshalToWithoutTag(payload[off:]); err != nil {
			return nil, fmt.Errorf("replayed UE security capability: %w", err)
		}
		off += msg.ReplayedUESecurityCapability.MarshalLenWithoutTag()
	}

	// Step 3: Encode tagged-optional IEs, encrypting HashMME in place.
	if msg.HashMME != nil {
		// Encrypt the HashMME value in place.
		cipheredValue, err := encryptPayload(param, msg.HashMME.Payload)
		if err != nil {
			return nil, fmt.Errorf("cipher HashMME: %w", err)
		}
		hashMME := *msg.HashMME
		hashMME.Payload = cipheredValue
		n := hashMME.MarshalLen()
		if _, err := writeIE(payload[off:], 0, &hashMME); err != nil {
			return nil, fmt.Errorf("HashMME: %w", err)
		}
		off += n
	}

	// Encode other named tagged IEs (plaintext).
	if msg.IMEISVRequest != nil {
		n := msg.IMEISVRequest.MarshalLen()
		if _, err := writeIE(payload[off:], 0, msg.IMEISVRequest); err != nil {
			return nil, fmt.Errorf("IMEISV request: %w", err)
		}
		off += n
	}
	if msg.NonceUE != nil {
		n := msg.NonceUE.MarshalLen()
		if _, err := writeIE(payload[off:], 0, msg.NonceUE); err != nil {
			return nil, fmt.Errorf("nonce UE: %w", err)
		}
		off += n
	}
	if msg.NonceMME != nil {
		n := msg.NonceMME.MarshalLen()
		if _, err := writeIE(payload[off:], 0, msg.NonceMME); err != nil {
			return nil, fmt.Errorf("nonce MME: %w", err)
		}
		off += n
	}
	if msg.ReplayedUEAdditionalSecurityCapability != nil {
		n := msg.ReplayedUEAdditionalSecurityCapability.MarshalLen()
		if _, err := writeIE(payload[off:], 0, msg.ReplayedUEAdditionalSecurityCapability); err != nil {
			return nil, fmt.Errorf("replayed UE additional security capability: %w", err)
		}
		off += n
	}
	if msg.UERadioCapabilityIDRequest != nil {
		n := msg.UERadioCapabilityIDRequest.MarshalLen()
		if _, err := writeIE(payload[off:], 0, msg.UERadioCapabilityIDRequest); err != nil {
			return nil, fmt.Errorf("UE radio capability ID request: %w", err)
		}
		off += n
	}
	if msg.UEInformationRequest != nil {
		n := msg.UEInformationRequest.MarshalLen()
		if _, err := writeIE(payload[off:], 0, msg.UEInformationRequest); err != nil {
			return nil, fmt.Errorf("UE information request: %w", err)
		}
		off += n
	}

	// Encode extra tagged IEs (plaintext).
	for _, i := range msg.IEs {
		n := i.MarshalLen()
		if _, err := writeIE(payload[off:], 0, i); err != nil {
			return nil, fmt.Errorf("extra IE %d: %w", i.Kind, err)
		}
		off += n
	}

	// Step 4: Compute MAC-I over the full payload (including ciphered HashMME).
	fresh := uint32(param.DownlinkNASCount & 0xff)
	macI, err := param.ComputeMACI(param.DownlinkNASCount, fresh, security.KNASDirectionDownlink, payload)
	if err != nil {
		return nil, fmt.Errorf("compute MAC-I: %w", err)
	}

	// Step 5: Encrypt the entire payload (mandatory + ciphered tagged).
	ciphered, err := encryptPayload(param, payload)
	if err != nil {
		return nil, fmt.Errorf("encrypt: %w", err)
	}

	// Step 6: Build header + encrypted payload.
	sht := byte(SecurityHeaderTypeIntegrityProtectedWithNewEPSContext) // 0x03
	pd := PdEmm
	hdr := make([]byte, 7)
	hdr[0] = byte(sht<<4 | pd)
	copy(hdr[1:5], macI[:])
	hdr[5] = byte(param.DownlinkNASCount & 0xff)

	return append(hdr, ciphered...), nil
}

// encodePayload returns the plaintext payload for non-ciphered encoding.
func (msg *SecurityModeCommand) encodePayload() ([]byte, error) {
	if msg == nil {
		return nil, fmt.Errorf("%w: SecurityModeCommand is nil", ie.ErrNilIE)
	}

	payload := make([]byte, msg.MarshalLen())
	off := 0

	if msg.SelectedNASSecurityAlgorithms != nil {
		if err := msg.SelectedNASSecurityAlgorithms.MarshalToWithoutTag(payload[off:]); err != nil {
			return nil, fmt.Errorf("selected NAS security algorithms: %w", err)
		}
		off += msg.SelectedNASSecurityAlgorithms.MarshalLenWithoutTag()
	}
	if msg.SpareHalfOctet != nil {
		if err := msg.SpareHalfOctet.MarshalToWithoutTag(payload[off:]); err != nil {
			return nil, fmt.Errorf("spare half octet: %w", err)
		}
		off++
	}
	if msg.ReplayedUESecurityCapability != nil {
		if err := msg.ReplayedUESecurityCapability.MarshalToWithoutTag(payload[off:]); err != nil {
			return nil, fmt.Errorf("replayed UE security capability: %w", err)
		}
		off += msg.ReplayedUESecurityCapability.MarshalLenWithoutTag()
	}

	// Named tagged IEs in spec-table order.
	if msg.IMEISVRequest != nil {
		n := msg.IMEISVRequest.MarshalLen()
		if _, err := writeIE(payload[off:], 0, msg.IMEISVRequest); err != nil {
			return nil, fmt.Errorf("IMEISV request: %w", err)
		}
		off += n
	}
	if msg.NonceUE != nil {
		n := msg.NonceUE.MarshalLen()
		if _, err := writeIE(payload[off:], 0, msg.NonceUE); err != nil {
			return nil, fmt.Errorf("nonce UE: %w", err)
		}
		off += n
	}
	if msg.NonceMME != nil {
		n := msg.NonceMME.MarshalLen()
		if _, err := writeIE(payload[off:], 0, msg.NonceMME); err != nil {
			return nil, fmt.Errorf("nonce MME: %w", err)
		}
		off += n
	}
	if msg.HashMME != nil {
		n := msg.HashMME.MarshalLen()
		if _, err := writeIE(payload[off:], 0, msg.HashMME); err != nil {
			return nil, fmt.Errorf("HashMME: %w", err)
		}
		off += n
	}
	if msg.ReplayedUEAdditionalSecurityCapability != nil {
		n := msg.ReplayedUEAdditionalSecurityCapability.MarshalLen()
		if _, err := writeIE(payload[off:], 0, msg.ReplayedUEAdditionalSecurityCapability); err != nil {
			return nil, fmt.Errorf("replayed UE additional security capability: %w", err)
		}
		off += n
	}
	if msg.UERadioCapabilityIDRequest != nil {
		n := msg.UERadioCapabilityIDRequest.MarshalLen()
		if _, err := writeIE(payload[off:], 0, msg.UERadioCapabilityIDRequest); err != nil {
			return nil, fmt.Errorf("UE radio capability ID request: %w", err)
		}
		off += n
	}
	if msg.UEInformationRequest != nil {
		n := msg.UEInformationRequest.MarshalLen()
		if _, err := writeIE(payload[off:], 0, msg.UEInformationRequest); err != nil {
			return nil, fmt.Errorf("UE information request: %w", err)
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

	return payload, nil
}

// Header returns the EMM security header embedded in SecurityModeCommand.
func (m *SecurityModeCommand) Header() *EmmHeader {
	return &m.EmmHeader
}
