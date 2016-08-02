package pdu

import (
	"encoding/binary"

	"github.com/mathpl/typed"
)

type Close struct {
	Reason Reason
}

func (c *Close) Type() Type {
	return TypeClose
}

func (c *Close) Read(f Flags, rb *typed.ReadBuffer) error {
	c.Reason = Reason(rb.ReadSingleByte())
	rb.ReadBytes(3)
	return rb.Err()
}

func (c *Close) Write(bo binary.ByteOrder, wb *typed.WriteBuffer) error {
	wb.WriteSingleByte(byte(c.Reason))
	wb.WriteBytes([]byte{0, 0, 0})
	return wb.Err()
}
