package pdu

import (
	"encoding/binary"

	"github.com/mathpl/typed"
)

type AllocateIndex struct {
	Variables Variables
}

func (a *AllocateIndex) Type() Type {
	return TypeIndexAllocate
}

func (a *AllocateIndex) Read(f Flags, rb *typed.ReadBuffer) error {
	a.Variables.Read(f, rb)
	return rb.Err()
}

func (a *AllocateIndex) Write(bo binary.ByteOrder, wb *typed.WriteBuffer) error {
	a.Variables.Write(bo, wb)
	return wb.Err()
}
