package pdu

import (
	"encoding/binary"

	"github.com/mathpl/typed"
)

type DeallocateIndex struct {
	Variables Variables
}

func (d *DeallocateIndex) Type() Type {
	return TypeIndexDeallocate
}

func (d *DeallocateIndex) Read(f Flags, rb *typed.ReadBuffer) error {
	d.Variables.Read(f, rb)
	return rb.Err()
}

func (d *DeallocateIndex) Write(bo binary.ByteOrder, wb *typed.WriteBuffer) error {
	d.Variables.Write(bo, wb)
	return wb.Err()
}
