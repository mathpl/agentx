package pdu

import (
	"encoding/binary"

	"github.com/mathpl/typed"
)

type Context struct {
	Data []byte
}

func (c *Context) Read(f Flags, rb *typed.ReadBuffer) error {
	l := int(rb.ReadUint32Endian(f.GetByteOrder()))
	c.Data = rb.ReadBytes(l)
	return rb.Err()
}

func (c *Context) Write(bo binary.ByteOrder, wb *typed.WriteBuffer) error {
	wb.WriteBytes(c.Data)
	return wb.Err()
}
