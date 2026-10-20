package ie

import "fmt"

// PLMN identifies a public land mobile network using its three-digit mobile
// country code and two- or three-digit mobile network code. It is encoded as
// specified by ITU-T E.212 and referenced throughout TS 24.301.
type PLMN struct {
	MCC string
	MNC string
}

func (p PLMN) validate() error {
	if len(p.MCC) != 3 || (len(p.MNC) != 2 && len(p.MNC) != 3) {
		return fmt.Errorf("%w: PLMN requires a three-digit MCC and a two- or three-digit MNC", ErrInvalidValue)
	}
	for _, digit := range p.MCC + p.MNC {
		if digit < '0' || digit > '9' {
			return fmt.Errorf("%w: PLMN contains a non-decimal digit", ErrInvalidValue)
		}
	}
	return nil
}

func (p PLMN) marshal() ([3]byte, error) {
	if err := p.validate(); err != nil {
		return [3]byte{}, err
	}
	mcc1 := p.MCC[0] - '0'
	mcc2 := p.MCC[1] - '0'
	mcc3 := p.MCC[2] - '0'
	mnc1 := p.MNC[0] - '0'
	mnc2 := p.MNC[1] - '0'
	mnc3 := byte(0x0f)
	if len(p.MNC) == 3 {
		mnc3 = p.MNC[2] - '0'
	}
	return [3]byte{
		mcc2<<4 | mcc1,
		mnc3<<4 | mcc3,
		mnc2<<4 | mnc1,
	}, nil
}

func parsePLMN(payload []byte) (PLMN, error) {
	if len(payload) != 3 {
		return PLMN{}, fmt.Errorf("%w: PLMN has %d octets, want 3", ErrInvalidLength, len(payload))
	}
	mcc1 := payload[0] & 0x0f
	mcc2 := payload[0] >> 4
	mcc3 := payload[1] & 0x0f
	mnc3 := payload[1] >> 4
	mnc1 := payload[2] & 0x0f
	mnc2 := payload[2] >> 4
	for _, digit := range []byte{mcc1, mcc2, mcc3, mnc1, mnc2} {
		if digit > 9 {
			return PLMN{}, fmt.Errorf("%w: PLMN contains an invalid BCD digit", ErrInvalidValue)
		}
	}
	if mnc3 > 9 && mnc3 != 0x0f {
		return PLMN{}, fmt.Errorf("%w: PLMN contains an invalid MNC BCD digit", ErrInvalidValue)
	}
	plmn := PLMN{
		MCC: string([]byte{'0' + mcc1, '0' + mcc2, '0' + mcc3}),
		MNC: string([]byte{'0' + mnc1, '0' + mnc2}),
	}
	if mnc3 != 0x0f {
		plmn.MNC += string([]byte{'0' + mnc3})
	}
	return plmn, nil
}
