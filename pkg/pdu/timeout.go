package pdu

import (
	"encoding/binary"
	"time"

	"github.com/mathpl/typed"
)

type Timeout struct {
	Duration time.Duration
	Priority byte
}

func (t Timeout) Read(f Flags, rb *typed.ReadBuffer) error {
	d := rb.ReadSingleByte()
	t.Duration = time.Duration(d) * time.Second
	t.Priority = rb.ReadSingleByte()
	// Reserved
	rb.ReadBytes(2)
	return rb.Err()
}

func (t *Timeout) Write(bo binary.ByteOrder, wb *typed.WriteBuffer) error {
	wb.WriteSingleByte(byte(t.Duration.Seconds()))
	wb.WriteSingleByte(t.Priority)
	wb.WriteBytes([]byte{0, 0})
	return wb.Err()
}
