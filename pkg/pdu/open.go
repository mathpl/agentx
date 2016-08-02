package pdu

import (
	"encoding/binary"

	"github.com/mathpl/typed"
)

type Open struct {
	Timeout     Timeout
	ID          ObjectIdentifier
	Description OctetString
}

func (o *Open) Type() Type {
	return TypeOpen
}

func (o *Open) Read(f Flags, rb *typed.ReadBuffer) error {
	o.Timeout.Read(f, rb)
	o.ID.Read(f, rb)
	o.Description.Read(f, rb)
	return rb.Err()
}

func (o *Open) Write(bo binary.ByteOrder, wb *typed.WriteBuffer) error {
	o.Timeout.Write(bo, wb)
	o.ID.Write(bo, wb)
	o.Description.Write(bo, wb)
	return wb.Err()
}
