package pdu

import (
	"encoding/binary"
	"time"

	"github.com/mathpl/typed"
)

type Response struct {
	UpTime    time.Duration
	Error     Error
	Index     uint16
	Variables Variables
}

func (r *Response) Type() Type {
	return TypeResponse
}

func (r *Response) Read(f Flags, rb *typed.ReadBuffer) error {
	bo := f.GetByteOrder()
	r.UpTime = time.Duration(rb.ReadUint32Endian(bo)*100) * time.Second
	r.Error = Error(rb.ReadUint16Endian(bo))
	r.Index = rb.ReadUint16Endian(bo)
	err := (&r.Variables).Read(f, rb)
	if err != nil {
		return err
	}
	return rb.Err()
}

func (r *Response) Write(bo binary.ByteOrder, wb *typed.WriteBuffer) error {
	uptime := uint32(r.UpTime.Seconds() / 100)
	wb.WriteUint32Endian(bo, uptime)
	wb.WriteUint16Endian(bo, uint16(r.Error))
	wb.WriteUint16Endian(bo, r.Index)
	r.Variables.Write(bo, wb)

	return wb.Err()
}
