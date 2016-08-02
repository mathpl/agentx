package pdu

import (
	"encoding/binary"

	"github.com/mathpl/typed"
)

type ObjectIdentifier struct {
	Prefix         uint8
	Include        byte
	Subidentifiers []uint32
}

func (o *ObjectIdentifier) SetIdentifier(oid OID) {
	o.Subidentifiers = make([]uint32, 0)

	if len(oid) > 4 && oid[0] == 1 && oid[1] == 3 && oid[2] == 6 && oid[3] == 1 {
		o.Subidentifiers = append(o.Subidentifiers, uint32(1), uint32(3), uint32(6), uint32(1), uint32(oid[4]))
		oid = oid[5:]
	}

	o.Subidentifiers = append(o.Subidentifiers, oid...)
}

func (o *ObjectIdentifier) GetIdentifier() OID {
	var oid OID
	if o.Prefix != 0 {
		oid = append(oid, 1, 3, 6, 1, uint32(o.Prefix))
	}
	return append(oid, o.Subidentifiers...)
}

func (o *ObjectIdentifier) Read(f Flags, rb *typed.ReadBuffer) {
	count := rb.ReadSingleByte()
	o.Prefix = rb.ReadSingleByte()
	o.Include = rb.ReadSingleByte()

	// Reserved
	rb.ReadSingleByte()

	o.Subidentifiers = make([]uint32, 0, count)

	for index := byte(0); index < count; index++ {
		o.Subidentifiers = append(o.Subidentifiers, rb.ReadUint32Endian(f.GetByteOrder()))
	}
}

func (o *ObjectIdentifier) Write(bo binary.ByteOrder, wb *typed.WriteBuffer) error {
	wb.WriteSingleByte(byte(len(o.Subidentifiers)))
	wb.WriteSingleByte(o.Prefix)
	wb.WriteSingleByte(o.Include)
	wb.WriteSingleByte(0)

	for _, subidentifier := range o.Subidentifiers {
		wb.WriteUint32Endian(bo, subidentifier)
	}

	return wb.Err()
}

func (o *ObjectIdentifier) String() string {
	return o.GetIdentifier().String()
}
