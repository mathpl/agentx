package pdu

import (
	"encoding/binary"

	"github.com/mathpl/typed"
)

const HeaderSize = 20

type Header struct {
	Version       byte
	Type          Type
	Flags         Flags
	SessionID     uint32
	TransactionID uint32
	PacketID      uint32
	PayloadLength uint32
}

func ReadHeader(rb *typed.ReadBuffer) (*Header, error) {
	h := &Header{}
	h.Version = rb.ReadSingleByte()
	h.Type = Type(rb.ReadSingleByte())
	h.Flags = Flags(rb.ReadSingleByte())

	// reseved
	rb.ReadSingleByte()

	bo := h.Flags.GetByteOrder()

	h.SessionID = rb.ReadUint32Endian(bo)
	h.TransactionID = rb.ReadUint32Endian(bo)
	h.PacketID = rb.ReadUint32Endian(bo)
	h.PayloadLength = rb.ReadUint32Endian(bo)

	return h, rb.Err()
}

func (h *Header) Write(bo binary.ByteOrder, wb *typed.WriteBuffer) (typed.Uint32Ref, error) {
	wb.WriteSingleByte(h.Version)
	wb.WriteSingleByte(byte(h.Type))
	wb.WriteSingleByte(byte(h.Flags))
	wb.WriteSingleByte(byte(0))

	wb.WriteUint32Endian(bo, h.SessionID)
	wb.WriteUint32Endian(bo, h.TransactionID)
	wb.WriteUint32Endian(bo, h.PacketID)

	// defered
	//wb.WriteUint32(bo, h.PayloadLength)

	return wb.DeferUint32(), wb.Err()
}
