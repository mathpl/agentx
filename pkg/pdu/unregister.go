package pdu

import (
	"encoding/binary"

	"github.com/mathpl/typed"
)

type Unregister struct {
	Timeout Timeout
	Subtree ObjectIdentifier
}

func (u *Unregister) Type() Type {
	return TypeUnregister
}

func (u *Unregister) Read(f Flags, rb *typed.ReadBuffer) error {
	u.Timeout.Read(f, rb)
	u.Subtree.Read(f, rb)
	return rb.Err()
}

func (u *Unregister) Write(bo binary.ByteOrder, wb *typed.WriteBuffer) error {
	u.Timeout.Write(bo, wb)
	u.Subtree.Write(bo, wb)
	return wb.Err()
}
