package ie

import (
	"fmt"
	"io"
)

// VoiceDomainPreferenceAndUEUsageSettingFields is the structured value of the
// Voice domain preference and UE's usage setting IE (TS 24.301 clause 9.9.3.44,
// see TS 24.008 clause 10.5.5.28).
type VoiceDomainPreferenceAndUEUsageSettingFields struct {
	// VoiceDomainPreference indicates the voice domain preference.
	VoiceDomainPreference bool
	// UEUsageSetting indicates the UE usage setting.
	UEUsageSetting uint8
}

// NewVoiceDomainPreferenceAndUEUsageSettingWithFields creates a tagged Voice
// domain preference and UE's usage setting IE from fields.
func NewVoiceDomainPreferenceAndUEUsageSettingWithFields(fields *VoiceDomainPreferenceAndUEUsageSettingFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewVoiceDomainPreferenceAndUEUsageSetting(payload[0]), nil
}

// VoiceDomainPreferenceAndUEUsageSetting returns the IE as structured fields.
func (i *IE) VoiceDomainPreferenceAndUEUsageSetting() (*VoiceDomainPreferenceAndUEUsageSettingFields, error) {
	if err := i.requireDefinition(VoiceDomainPreferenceAndUEUsageSettingDefinition); err != nil {
		return nil, err
	}
	return ParseVoiceDomainPreferenceAndUEUsageSettingFields(i.Payload)
}

// ParseVoiceDomainPreferenceAndUEUsageSettingFields decodes a Voice domain
// preference and UE's usage setting value without its IEI.
func ParseVoiceDomainPreferenceAndUEUsageSettingFields(payload []byte) (*VoiceDomainPreferenceAndUEUsageSettingFields, error) {
	fields := &VoiceDomainPreferenceAndUEUsageSettingFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that UEUsageSetting fits 3 bits.
func (f *VoiceDomainPreferenceAndUEUsageSettingFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: voice domain preference and UE usage setting fields are nil", ErrInvalidValue)
	}
	if f.UEUsageSetting > 0x07 {
		return fmt.Errorf("%w: UE usage setting exceeds its encoded width", ErrInvalidValue)
	}
	return nil
}

// Marshal encodes fields into a Voice domain preference and UE's usage setting
// value, without its IEI.
func (f *VoiceDomainPreferenceAndUEUsageSettingFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	var encoded byte
	if f.VoiceDomainPreference {
		encoded |= 0xC0
	}
	encoded |= byte(f.UEUsageSetting) & 0x07
	return []byte{encoded}, nil
}

// MarshalTo encodes fields into b.
func (f *VoiceDomainPreferenceAndUEUsageSettingFields) MarshalTo(b []byte) error {
	payload, err := f.Marshal()
	if err != nil {
		return err
	}
	if len(b) < len(payload) {
		return io.ErrShortBuffer
	}
	b[0] = payload[0]
	return nil
}

// MarshalLen returns the size of the Voice domain preference and UE's usage
// setting value payload.
func (f *VoiceDomainPreferenceAndUEUsageSettingFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a Voice domain preference and UE's usage setting value
// without its IEI.
func (f *VoiceDomainPreferenceAndUEUsageSettingFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: voice domain preference and UE usage setting fields are nil", ErrInvalidValue)
	}
	if len(payload) == 0 {
		return io.ErrShortBuffer
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: voice domain preference and UE usage setting has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	if payload[0]&0x38 != 0 {
		return fmt.Errorf("%w: voice domain preference and UE usage setting spare bits are nonzero", ErrInvalidValue)
	}
	*f = VoiceDomainPreferenceAndUEUsageSettingFields{
		VoiceDomainPreference: payload[0]&0xC0 == 0xC0,
		UEUsageSetting:        payload[0] & 0x07,
	}
	return nil
}
