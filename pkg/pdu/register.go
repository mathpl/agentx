package pdu

import (
	"encoding/binary"

	"github.com/mathpl/typed"
)

type Register struct {
	Timeout Timeout
	Subtree ObjectIdentifier
}

func (r *Register) Type() Type {
	return TypeRegister
}

func (r *Register) Read(f Flags, rb *typed.ReadBuffer) error {
	r.Timeout.Read(f, rb)
	r.Subtree.Read(f, rb)
	return rb.Err()
}

func (r *Register) Write(bo binary.ByteOrder, wb *typed.WriteBuffer) error {
	r.Timeout.Write(bo, wb)
	r.Subtree.Write(bo, wb)
	return wb.Err()
}
