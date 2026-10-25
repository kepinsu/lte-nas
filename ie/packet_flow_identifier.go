package ie

import (
	"fmt"
	"io"
)

// Pfi identifies the packet flow, as defined in TS 24.301 clause 9.9.4.8.
type Pfi uint8

// Packet flow identifier values.
const (
	PfiReserved       Pfi = 0x00
	PfiUEIdentity     Pfi = 0x01
	PfiBearerIdentity Pfi = 0x02
)

// PacketFlowIdentifierFields is the structured value of the positional Packet
// flow identifier IE.
type PacketFlowIdentifierFields struct {
	Type Pfi
}

// NewPositionalPacketFlowIdentifierWithFields creates a positional Packet flow
// identifier IE from fields.
func NewPositionalPacketFlowIdentifierWithFields(fields *PacketFlowIdentifierFields) (*IE, error) {
	payload, err := fields.Marshal()
	if err != nil {
		return nil, err
	}
	return NewPositionalPacketFlowIdentifier(payload[0]), nil
}

// PacketFlowIdentifier returns the IE as structured fields.
func (i *IE) PacketFlowIdentifier() (*PacketFlowIdentifierFields, error) {
	if err := i.requireDefinition(PacketFlowIdentifierDefinition); err != nil {
		return nil, err
	}
	return ParsePacketFlowIdentifierFields(i.Payload)
}

// ParsePacketFlowIdentifierFields decodes a Packet flow identifier value
// without its IEI.
func ParsePacketFlowIdentifierFields(payload []byte) (*PacketFlowIdentifierFields, error) {
	fields := &PacketFlowIdentifierFields{}
	if err := fields.UnmarshalBinary(payload); err != nil {
		return nil, err
	}
	return fields, nil
}

// Validate verifies that Type is a defined Packet flow identifier value.
func (f *PacketFlowIdentifierFields) Validate() error {
	if f == nil {
		return fmt.Errorf("%w: Packet flow identifier fields are nil", ErrInvalidValue)
	}
	if !f.Type.valid() {
		return fmt.Errorf("%w: Packet flow identifier %#02x is reserved", ErrInvalidValue, f.Type)
	}
	return nil
}

func (t Pfi) valid() bool {
	switch t {
	case PfiReserved, PfiUEIdentity, PfiBearerIdentity:
		return true
	default:
		return false
	}
}

// Marshal encodes fields into a Packet flow identifier value without its IEI.
func (f *PacketFlowIdentifierFields) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return []byte{byte(f.Type)}, nil
}

// MarshalTo encodes fields into b.
func (f *PacketFlowIdentifierFields) MarshalTo(b []byte) error {
	if err := f.Validate(); err != nil {
		return err
	}
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	b[0] = byte(f.Type)
	return nil
}

// MarshalLen returns the size of the Packet flow identifier value payload.
func (f *PacketFlowIdentifierFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1
}

// UnmarshalBinary decodes a Packet flow identifier value without its IEI.
func (f *PacketFlowIdentifierFields) UnmarshalBinary(payload []byte) error {
	if f == nil {
		return fmt.Errorf("%w: Packet flow identifier fields are nil", ErrInvalidValue)
	}
	if len(payload) != 1 {
		return fmt.Errorf("%w: Packet flow identifier has %d value octets, want 1", ErrInvalidLength, len(payload))
	}
	decoded := PacketFlowIdentifierFields{Type: Pfi(payload[0])}
	if err := decoded.Validate(); err != nil {
		return err
	}
	*f = decoded
	return nil
}
