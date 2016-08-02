package pdu

import (
	"encoding/binary"

	"github.com/mathpl/typed"
)

type Packet interface {
	TypeOwner
	Read(Flags, *typed.ReadBuffer) error
	Write(binary.ByteOrder, *typed.WriteBuffer) error
}
