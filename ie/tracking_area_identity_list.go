package ie

import (
	"encoding/binary"
	"fmt"
	"io"
)

const maximumTrackingAreaIdentityListEntries = 16

const (
	trackingAreaIdentityListNonConsecutiveTACs byte = iota
	trackingAreaIdentityListConsecutiveTACs
	trackingAreaIdentityListDifferentPLMNs
)

// TrackingAreaIdentityListFields is the structured value of the Tracking area
// identity list IE (TS 24.301 clause 9.9.3.33). It contains the expanded list
// of TAIs; the encoder selects the compact list form when adjacent entries have
// a shared PLMN.
type TrackingAreaIdentityListFields struct {
	TrackingAreas []TrackingAreaIdentityFields
}

// NewForbiddenTAIsForRoamingWithFields creates a tagged forbidden TAIs for
// roaming IE from fields.
func NewForbiddenTAIsForRoamingWithFields(fields *TrackingAreaIdentityListFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewForbiddenTAIsForRoaming(payload), nil
}

// NewForbiddenTAIsForRegionalProvisionOfServiceWithFields creates a tagged
// forbidden TAIs for regional provision of service IE from fields.
func NewForbiddenTAIsForRegionalProvisionOfServiceWithFields(fields *TrackingAreaIdentityListFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewForbiddenTAIsForRegionalProvisionOfService(payload), nil
}

// NewTAIListWithFields creates a tagged TAI list IE from fields.
func NewTAIListWithFields(fields *TrackingAreaIdentityListFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewTAIList(payload), nil
}

// NewPositionalTrackingAreaIdentityListWithFields creates a positional TAI
// list IE from fields.
func NewPositionalTrackingAreaIdentityListWithFields(fields *TrackingAreaIdentityListFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalTrackingAreaIdentityList(payload), nil
}

// TrackingAreaIdentityList returns the IE as structured fields.
func (i *IE) TrackingAreaIdentityList() (*TrackingAreaIdentityListFields, error) {
	if err := i.requireDefinition(TrackingAreaIdentityListDefinition); err != nil {
		return nil, err
	}
	return ParseTrackingAreaIdentityListFields(i.Payload)
}

// ParseTrackingAreaIdentityListFields decodes a Tracking area identity list
// value without its IEI and length octet.
func ParseTrackingAreaIdentityListFields(payload []byte) (*TrackingAreaIdentityListFields, error) {
	fields := &TrackingAreaIdentityListFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that the list has a valid 3GPP size and valid TAI values.
func (f *TrackingAreaIdentityListFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: tracking area identity list fields are nil", ErrInvalidValue)
	}
	if len(f.TrackingAreas) == 0 || len(f.TrackingAreas) > maximumTrackingAreaIdentityListEntries {
		return fmt.Errorf("%w: tracking area identity list has %d entries, want 1..%d", ErrInvalidLength, len(f.TrackingAreas), maximumTrackingAreaIdentityListEntries)
	}
	for index := range f.TrackingAreas {
		if err := f.TrackingAreas[index].Validate(); err != nil {
			return fmt.Errorf("tracking area %d: %w", index, err)
		}
	}
	return nil
}

// Marshal encodes fields into a Tracking area identity list value without its
// IEI and length octet.
func (f *TrackingAreaIdentityListFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	payload := make([]byte, 0, f.MarshalLen())
	for index := 0; index < len(f.TrackingAreas); {
		end := trackingAreasWithSamePLMN(f.TrackingAreas, index)
		if end-index > 1 {
			if trackingAreasAreConsecutive(f.TrackingAreas[index:end]) {
				payload = appendConsecutiveTAIs(payload, f.TrackingAreas[index:end])
			} else {
				payload = appendNonConsecutiveTAIs(payload, f.TrackingAreas[index:end])
			}
			index = end
			continue
		}
		end = trackingAreasWithoutAdjacentSharedPLMN(f.TrackingAreas, index)
		payload = appendDifferentPLMNTAIs(payload, f.TrackingAreas[index:end])
		index = end
	}
	return payload, nil
}

// MarshalTo encodes fields into b.
func (f *TrackingAreaIdentityListFields) MarshalTo(b []byte) error {
	payload, err := f.Marshal()
	if err != nil {
		return err
	}
	if len(b) < len(payload) {
		return io.ErrShortBuffer
	}
	copy(b, payload)
	return nil
}

// MarshalLen returns the size of the Tracking area identity list value payload
// produced by Marshal. It returns zero for a nil receiver.
func (f *TrackingAreaIdentityListFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	length := 0
	for index := 0; index < len(f.TrackingAreas); {
		end := trackingAreasWithSamePLMN(f.TrackingAreas, index)
		switch {
		case end-index == 1:
			end = trackingAreasWithoutAdjacentSharedPLMN(f.TrackingAreas, index)
			length += 1 + 5*(end-index)
		case trackingAreasAreConsecutive(f.TrackingAreas[index:end]):
			length += 6
		default:
			length += 4 + 2*(end-index)
		}
		index = end
	}
	return length
}

// UnmarshalBinary decodes a Tracking area identity list value without its IEI
// and length octet. It stores the first 16 TAIs, as specified by TS 24.301.
func (f *TrackingAreaIdentityListFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: tracking area identity list fields are nil", ErrInvalidValue)
	}
	if len(payload) < 6 || len(payload) > 96 {
		return fmt.Errorf("%w: tracking area identity list has %d value octets, want 6..96", ErrInvalidLength, len(payload))
	}
	trackingAreas := make([]TrackingAreaIdentityFields, 0, maximumTrackingAreaIdentityListEntries)
	for offset := 0; offset < len(payload); {
		if err := decodePartialTrackingAreaIdentityList(payload, &offset, &trackingAreas); err != nil {
			return err
		}
	}
	*f = TrackingAreaIdentityListFields{TrackingAreas: trackingAreas}
	return nil
}

func trackingAreasWithSamePLMN(trackingAreas []TrackingAreaIdentityFields, start int) int {
	end := start + 1
	for end < len(trackingAreas) && end-start < maximumTrackingAreaIdentityListEntries && trackingAreas[end].PLMN == trackingAreas[start].PLMN {
		end++
	}
	return end
}

func trackingAreasWithoutAdjacentSharedPLMN(trackingAreas []TrackingAreaIdentityFields, start int) int {
	end := start + 1
	for end < len(trackingAreas) && end-start < maximumTrackingAreaIdentityListEntries {
		nextEnd := trackingAreasWithSamePLMN(trackingAreas, end)
		if nextEnd-end > 1 {
			break
		}
		end = nextEnd
	}
	return end
}

func trackingAreasAreConsecutive(trackingAreas []TrackingAreaIdentityFields) bool {
	for index := 1; index < len(trackingAreas); index++ {
		if uint32(trackingAreas[index-1].TrackingAreaCode)+1 != uint32(trackingAreas[index].TrackingAreaCode) {
			return false
		}
	}
	return true
}

func appendNonConsecutiveTAIs(payload []byte, trackingAreas []TrackingAreaIdentityFields) []byte {
	payload = append(payload, byte(trackingAreaIdentityListNonConsecutiveTACs<<5)|byte(len(trackingAreas)-1))
	plmn, _ := trackingAreas[0].PLMN.marshal()
	payload = append(payload, plmn[:]...)
	for _, trackingArea := range trackingAreas {
		payload = binary.BigEndian.AppendUint16(payload, trackingArea.TrackingAreaCode)
	}
	return payload
}

func appendConsecutiveTAIs(payload []byte, trackingAreas []TrackingAreaIdentityFields) []byte {
	payload = append(payload, byte(trackingAreaIdentityListConsecutiveTACs<<5)|byte(len(trackingAreas)-1))
	plmn, _ := trackingAreas[0].PLMN.marshal()
	payload = append(payload, plmn[:]...)
	return binary.BigEndian.AppendUint16(payload, trackingAreas[0].TrackingAreaCode)
}

func appendDifferentPLMNTAIs(payload []byte, trackingAreas []TrackingAreaIdentityFields) []byte {
	payload = append(payload, byte(trackingAreaIdentityListDifferentPLMNs<<5)|byte(len(trackingAreas)-1))
	for _, trackingArea := range trackingAreas {
		plmn, _ := trackingArea.PLMN.marshal()
		payload = append(payload, plmn[:]...)
		payload = binary.BigEndian.AppendUint16(payload, trackingArea.TrackingAreaCode)
	}
	return payload
}

func decodePartialTrackingAreaIdentityList(payload []byte, offset *int, trackingAreas *[]TrackingAreaIdentityFields) error {
	if *offset >= len(payload) {
		return fmt.Errorf("%w: tracking area identity list is missing a partial-list header", ErrInvalidLength)
	}
	header := payload[*offset]
	*offset += 1
	if header&0x80 != 0 {
		return fmt.Errorf("%w: tracking area identity list spare bit is nonzero", ErrInvalidValue)
	}
	listType := (header >> 5) & 0x03
	numberOfElements := int(header&0x1f) + 1
	if numberOfElements > maximumTrackingAreaIdentityListEntries {
		numberOfElements = maximumTrackingAreaIdentityListEntries
	}
	switch listType {
	case trackingAreaIdentityListNonConsecutiveTACs:
		return decodeNonConsecutiveTAIs(payload, offset, numberOfElements, trackingAreas)
	case trackingAreaIdentityListConsecutiveTACs:
		return decodeConsecutiveTAIs(payload, offset, numberOfElements, trackingAreas)
	case trackingAreaIdentityListDifferentPLMNs:
		return decodeDifferentPLMNTAIs(payload, offset, numberOfElements, trackingAreas)
	default:
		return fmt.Errorf("%w: tracking area identity list type %d is reserved", ErrInvalidValue, listType)
	}
}

func decodeNonConsecutiveTAIs(payload []byte, offset *int, numberOfElements int, trackingAreas *[]TrackingAreaIdentityFields) error {
	if len(payload)-*offset < 3+2*numberOfElements {
		return fmt.Errorf("%w: tracking area identity list non-consecutive partial list is truncated", ErrInvalidLength)
	}
	plmn, err := parsePLMN(payload[*offset : *offset+3])
	if err != nil {
		return err
	}
	*offset += 3
	for index := 0; index < numberOfElements; index++ {
		trackingAreaCode := binary.BigEndian.Uint16(payload[*offset:])
		*offset += 2
		appendTrackingArea(trackingAreas, TrackingAreaIdentityFields{PLMN: plmn, TrackingAreaCode: trackingAreaCode})
	}
	return nil
}

func decodeConsecutiveTAIs(payload []byte, offset *int, numberOfElements int, trackingAreas *[]TrackingAreaIdentityFields) error {
	if len(payload)-*offset < 5 {
		return fmt.Errorf("%w: tracking area identity list consecutive partial list is truncated", ErrInvalidLength)
	}
	plmn, err := parsePLMN(payload[*offset : *offset+3])
	if err != nil {
		return err
	}
	*offset += 3
	firstTrackingAreaCode := binary.BigEndian.Uint16(payload[*offset:])
	*offset += 2
	if uint32(firstTrackingAreaCode)+uint32(numberOfElements-1) > 0xffff {
		return fmt.Errorf("%w: tracking area identity list consecutive TACs overflow", ErrInvalidValue)
	}
	for index := 0; index < numberOfElements; index++ {
		appendTrackingArea(trackingAreas, TrackingAreaIdentityFields{
			PLMN:             plmn,
			TrackingAreaCode: firstTrackingAreaCode + uint16(index),
		})
	}
	return nil
}

func decodeDifferentPLMNTAIs(payload []byte, offset *int, numberOfElements int, trackingAreas *[]TrackingAreaIdentityFields) error {
	if len(payload)-*offset < 5*numberOfElements {
		return fmt.Errorf("%w: tracking area identity list different-PLMN partial list is truncated", ErrInvalidLength)
	}
	for index := 0; index < numberOfElements; index++ {
		plmn, err := parsePLMN(payload[*offset : *offset+3])
		if err != nil {
			return err
		}
		*offset += 3
		trackingAreaCode := binary.BigEndian.Uint16(payload[*offset:])
		*offset += 2
		appendTrackingArea(trackingAreas, TrackingAreaIdentityFields{PLMN: plmn, TrackingAreaCode: trackingAreaCode})
	}
	return nil
}

func appendTrackingArea(trackingAreas *[]TrackingAreaIdentityFields, trackingArea TrackingAreaIdentityFields) {
	if len(*trackingAreas) < maximumTrackingAreaIdentityListEntries {
		*trackingAreas = append(*trackingAreas, trackingArea)
	}
}
