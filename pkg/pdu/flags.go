package pdu

import "encoding/binary"

type Flags byte

const (
	FlagInstanceRegistration Flags = 1 << 0
	FlagNewIndex             Flags = 1 << 1
	FlagAnyIndex             Flags = 1 << 2
	FlagNonDefaultContext    Flags = 1 << 3
	FlagNetworkByteOrder     Flags = 1 << 4
)

func (f Flags) GetByteOrder() binary.ByteOrder {
	var bo binary.ByteOrder
	bo = binary.LittleEndian
	if f&FlagNetworkByteOrder != 0 {
		bo = binary.BigEndian
	}
	return bo
}

func (f Flags) GetNonDefaultContext() bool {
	return (f&FlagNonDefaultContext != 0)
}

func (f Flags) GetNewIndex() bool {
	return (f&FlagNewIndex != 0)
}

func (f Flags) GetAnyIndex() bool {
	return (f&FlagAnyIndex != 0)
}
