package pdu

import (
	"encoding/binary"

	"github.com/mathpl/typed"
)

type OctetString struct {
	Text string
}

func (o *OctetString) Read(f Flags, rb *typed.ReadBuffer) error {
	l := int(rb.ReadUint32Endian(f.GetByteOrder()))
	o.Text = rb.ReadString(l)
	return rb.Err()
}

func (o *OctetString) Write(bo binary.ByteOrder, wb *typed.WriteBuffer) error {
	wb.WriteUint32Endian(bo, uint32(len(o.Text)))
	wb.WriteString(o.Text)
	return wb.Err()
}
