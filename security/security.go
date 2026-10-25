// Package security implements the EPS NAS security procedures defined by
// 3GPP TS 33.401: key derivation from the K_ASME anchor key, NAS ciphering
// (EEA0, EEA1, EEA2), NAS integrity protection (EIA0, EIA1, EIA2), eKSI
// handling and NAS COUNT^C counter management.
//
// The cryptographic primitives are provided by the sibling snow3g (SNOW 3G
// based UEA2/UIA2) and myaes (AES based EEA2/EIA2) packages.
package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"math"

	"github.com/kepinsu/lte-nas/ie"
	"github.com/kepinsu/lte-nas/security/myaes"
	"github.com/kepinsu/lte-nas/security/snow3g"
)

// Supported NAS ciphering algorithm identifiers.
const (
	// EEA0 is the null ciphering algorithm (no protection).
	EEA0 ie.NASEncryptionAlgorithm = iota
	// EEA1 is the 128-bit SNOW 3G based ciphering algorithm (UEA2, F8).
	EEA1
	// EEA2 is the 128-bit AES based ciphering algorithm.
	EEA2
)

// Supported NAS integrity algorithm identifiers.
const (
	// EIA0 is the null integrity algorithm.
	EIA0 ie.NASIntegrityAlgorithm = iota
	// EIA1 is the 128-bit SNOW 3G based integrity algorithm (UIA2, F9).
	EIA1
	// EIA2 is the 128-bit AES based integrity algorithm.
	EIA2
)

const (
	// KNASDirectionUplink is the UE-to-network transmission direction bit.
	KNASDirectionUplink = 0
	// KNASDirectionDownlink is the network-to-UE transmission direction bit.
	KNASDirectionDownlink = 1

	// KSI_BIT_MASK selects the 3 key-set-identifier bits of an eKSI octet.
	KSI_BIT_MASK = 0x07
	// SN_BIT_MASK selects the 4 serving-node bits of an eKSI octet.
	SN_BIT_MASK = 0x78

	// prfNASenc is the TS 33.401 annex A algType for CK derivation.
	prfNASenc = 0x03
	// prfNASint is the TS 33.401 annex A algType for IK derivation.
	prfNASint = 0x02

	kasmeSize = 32 // K_ASME length in octets
	keySize   = 16 // CK / IK length in octets
	macISize  = 4  // MAC-I length in octets
)

// Errors returned by the security API.
var (
	ErrInvalidKASME    = errors.New("invalid K_ASME")
	ErrInvalidEKSI     = errors.New("invalid eKSI")
	ErrInvalidKey      = errors.New("invalid NAS security key")
	ErrInvalidPayload  = errors.New("invalid NAS payload")
	ErrInvalidMACI     = errors.New("invalid MAC-I")
	ErrUnsupportedAlgo = errors.New("unsupported NAS security algorithm")
)

// EMMSecurityContext holds the EPS NAS security context established by the
// EMM security mode command procedure (TS 24.301 9.1.3.6).
type EMMSecurityContext struct {
	// Kasme is the 32-octet anchor key K_ASME.
	Kasme []byte
	// CK is the 16-octet NAS ciphering key derived from K_ASME.
	CK []byte
	// IK is the 16-octet NAS integrity key derived from K_ASME.
	IK []byte
	// eKSI is the key set identifier selected by the network.
	eKSI uint8
	// IntAlgo is the selected NAS integrity protection algorithm.
	IntAlgo ie.NASIntegrityAlgorithm
	// EncAlgo is the selected NAS ciphering algorithm.
	EncAlgo ie.NASEncryptionAlgorithm
	// UplinkNASCount is the COUNT^C value for UE-to-network NAS messages.
	UplinkNASCount uint32
	// DownlinkNASCount is the COUNT^C value for network-to-UE NAS messages.
	DownlinkNASCount uint32
	// Valid reports whether the context holds a complete key set.
	Valid bool
	// Active reports whether the context is applied to the NAS link.
	Active bool
}

// DeriveCK derives the 16-octet NAS ciphering key CK from the 32-octet
// K_ASME anchor key for the given ciphering algorithm identifier, per TS
// 33.401 annex A.2.
func DeriveCK(kasme []byte, algoID byte) []byte {
	return deriveKey(kasme, prfNASenc, algoID)
}

// DeriveIK derives the 16-octet NAS integrity key IK from the 32-octet
// K_ASME anchor key for the given integrity algorithm identifier, per TS
// 33.401 annex A.2.
func DeriveIK(kasme []byte, algoID byte) []byte {
	return deriveKey(kasme, prfNASint, algoID)
}

// SetupContext creates a validated EMMSecurityContext from the 32-octet
// K_ASME anchor key, deriving CK and IK per TS 33.401 annex A.2.
func SetupContext(kasme []byte, ksi uint8, encAlgo ie.NASEncryptionAlgorithm, intAlgo ie.NASIntegrityAlgorithm) (*EMMSecurityContext, error) {
	if len(kasme) != kasmeSize {
		return nil, fmt.Errorf("%w: expected %d octets, got %d", ErrInvalidKASME, kasmeSize, len(kasme))
	}
	if ksi > KSI_BIT_MASK {
		return nil, fmt.Errorf("%w: KSI exceeds 3 bits (%d)", ErrInvalidEKSI, ksi)
	}
	if !supportedEncAlgo(encAlgo) {
		return nil, fmt.Errorf("%w: ciphering algorithm %d", ErrUnsupportedAlgo, encAlgo)
	}
	if !supportedIntAlgo(intAlgo) {
		return nil, fmt.Errorf("%w: integrity algorithm %d", ErrUnsupportedAlgo, intAlgo)
	}

	return &EMMSecurityContext{
		Kasme:   append([]byte(nil), kasme...),
		CK:      DeriveCK(kasme, byte(encAlgo)),
		IK:      DeriveIK(kasme, byte(intAlgo)),
		eKSI:    ksi,
		IntAlgo: intAlgo,
		EncAlgo: encAlgo,
		Valid:   true,
	}, nil
}

// Encrypt applies the selected NAS ciphering algorithm to plaintext with the
// downlink COUNT^C value and returns the ciphertext.
func (c *EMMSecurityContext) Encrypt(nasCount uint32, bearer uint8, plaintext []byte) ([]byte, error) {
	return c.cipher(nasCount, bearer, plaintext, KNASDirectionDownlink)
}

// Decrypt applies the selected NAS ciphering algorithm to ciphertext with the
// uplink COUNT^C value and returns the plaintext.
func (c *EMMSecurityContext) Decrypt(nasCount uint32, bearer uint8, ciphertext []byte) ([]byte, error) {
	return c.cipher(nasCount, bearer, ciphertext, KNASDirectionUplink)
}

// ComputeMACI computes the 4-octet MAC-I over msg with the selected NAS
// integrity algorithm.
func (c *EMMSecurityContext) ComputeMACI(nasCount uint32, fresh uint32, direction uint8, msg []byte) ([]byte, error) {
	if len(msg) == 0 {
		return nil, fmt.Errorf("%w: message is empty", ErrInvalidPayload)
	}
	switch c.IntAlgo {
	case EIA0:
		return make([]byte, macISize), nil
	case EIA1:
		if err := c.checkIntegrityKey(); err != nil {
			return nil, err
		}
		return snow3g.NewUIA2(c.IK, nasCount, fresh, direction).F9(msg, uint64(len(msg)*8)), nil
	case EIA2:
		if err := c.checkIntegrityKey(); err != nil {
			return nil, err
		}
		if fresh > math.MaxUint8 {
			return nil, fmt.Errorf("%w: FRESH exceeds 8 bits (%d)", ErrInvalidPayload, fresh)
		}
		return myaes.EIA2(c.IK, nasCount, uint8(fresh), direction, msg)
	default:
		return nil, fmt.Errorf("%w: integrity algorithm %d", ErrUnsupportedAlgo, c.IntAlgo)
	}
}

// VerifyMACI recomputes the MAC-I over msg and compares it with receivedMACI
// in constant time, reporting whether the integrity check passed.
func (c *EMMSecurityContext) VerifyMACI(nasCount uint32, fresh uint32, direction uint8, msg []byte, receivedMACI []byte) (bool, error) {
	if len(receivedMACI) != macISize {
		return false, fmt.Errorf("%w: expected %d octets, got %d", ErrInvalidMACI, macISize, len(receivedMACI))
	}
	computed, err := c.ComputeMACI(nasCount, fresh, direction, msg)
	if err != nil {
		return false, err
	}
	return hmac.Equal(computed, receivedMACI), nil
}

// GeteKSI packs the 3-bit key set identifier and the 4-bit serving node
// identifier into a single eKSI octet: KSI occupies bits 0-2 and SN bits 3-6.
func GeteKSI(ksi, sn uint8) uint8 {
	return (sn&0x0F)<<3 | (ksi & KSI_BIT_MASK)
}

// ParseeKSI unpacks the KSI (bits 0-2) and SN (bits 3-6) from an eKSI octet.
func ParseeKSI(eKSI uint8) (ksi, sn uint8) {
	return eKSI & KSI_BIT_MASK, (eKSI & SN_BIT_MASK) >> 3
}

// IncrementUplinkCount advances the UE-to-network NAS COUNT^C value.
func (c *EMMSecurityContext) IncrementUplinkCount() {
	c.UplinkNASCount++
}

// IncrementDownlinkCount advances the network-to-UE NAS COUNT^C value.
func (c *EMMSecurityContext) IncrementDownlinkCount() {
	c.DownlinkNASCount++
}

// deriveKey evaluates the TS 33.401 annex A key derivation function:
// HMAC-SHA256(K, algType || length || algoID || lenAlgo) truncated to the
// first 16 octets, with length fixed to 0x0010 and lenAlgo to 0x01.
func deriveKey(kasme []byte, algType, algoID byte) []byte {
	p := []byte{algType, 0x00, 0x10, algoID, 0x01}
	mac := hmac.New(sha256.New, kasme)
	mac.Write(p)
	return mac.Sum(nil)[:keySize]
}

// cipher routes the selected ciphering algorithm for the given direction.
func (c *EMMSecurityContext) cipher(nasCount uint32, bearer uint8, data []byte, direction uint8) ([]byte, error) {
	if data == nil {
		return nil, fmt.Errorf("%w: payload is nil", ErrInvalidPayload)
	}
	switch c.EncAlgo {
	case EEA0:
		return data, nil
	case EEA1:
		if err := c.checkCipheringKey(); err != nil {
			return nil, err
		}
		return snow3g.NewUEA2(c.CK, nasCount, uint32(bearer), uint32(direction)).F8(data, uint32(len(data)*8)), nil
	case EEA2:
		if err := c.checkCipheringKey(); err != nil {
			return nil, err
		}
		return myaes.EEA2(c.CK, nasCount, bearer, direction, data)
	default:
		return nil, fmt.Errorf("%w: ciphering algorithm %d", ErrUnsupportedAlgo, c.EncAlgo)
	}
}

// checkCipheringKey verifies that CK is present and 16 octets long.
func (c *EMMSecurityContext) checkCipheringKey() error {
	if len(c.CK) != keySize {
		return fmt.Errorf("%w: CK must be %d octets, got %d", ErrInvalidKey, keySize, len(c.CK))
	}
	return nil
}

// checkIntegrityKey verifies that IK is present and 16 octets long.
func (c *EMMSecurityContext) checkIntegrityKey() error {
	if len(c.IK) != keySize {
		return fmt.Errorf("%w: IK must be %d octets, got %d", ErrInvalidKey, keySize, len(c.IK))
	}
	return nil
}

// supportedEncAlgo reports whether the ciphering algorithm is implemented.
func supportedEncAlgo(a ie.NASEncryptionAlgorithm) bool {
	return a >= EEA0 && a <= EEA2
}

// supportedIntAlgo reports whether the integrity algorithm is implemented.
func supportedIntAlgo(a ie.NASIntegrityAlgorithm) bool {
	return a >= EIA0 && a <= EIA2
}
