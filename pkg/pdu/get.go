package pdu

import (
	"encoding/binary"

	"github.com/mathpl/typed"
)

type Get struct {
	SearchRange Range
}

func (g *Get) GetOID() OID {
	return g.SearchRange.From.GetIdentifier()
}

func (g *Get) SetOID(oid OID) {
	g.SearchRange.From.SetIdentifier(oid)
}

func (g *Get) Type() Type {
	return TypeGet
}

func (g *Get) Read(f Flags, rb *typed.ReadBuffer) error {
	g.SearchRange.Read(f, rb)
	return rb.Err()
}

func (g *Get) Write(bo binary.ByteOrder, wb *typed.WriteBuffer) error {
	g.SearchRange.Write(bo, wb)
	return wb.Err()
}
