package pdu

import (
	"encoding/binary"

	"github.com/mathpl/typed"
)

type AddAgentCaps struct {
	ID          ObjectIdentifier
	Description OctetString
}

func (a *AddAgentCaps) Type() Type {
	return TypeAddAgentCaps
}

func (a *AddAgentCaps) Read(f Flags, rb *typed.ReadBuffer) error {
	a.ID.Read(f, rb)
	a.Description.Read(f, rb)
	return rb.Err()
}

func (a *AddAgentCaps) Write(bo binary.ByteOrder, wb *typed.WriteBuffer) error {
	a.ID.Write(bo, wb)
	a.Description.Write(bo, wb)
	return wb.Err()
}
