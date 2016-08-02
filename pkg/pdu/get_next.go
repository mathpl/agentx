package pdu

import (
	"encoding/binary"

	"github.com/mathpl/typed"
)

type GetNext struct {
	SearchRanges Ranges
}

func (g *GetNext) Type() Type {
	return TypeGetNext
}

func (g *GetNext) Read(f Flags, rb *typed.ReadBuffer) error {
	g.SearchRanges.Read(f, rb)
	return rb.Err()
}

func (g *GetNext) Write(bo binary.ByteOrder, wb *typed.WriteBuffer) error {
	g.SearchRanges.Write(bo, wb)
	return wb.Err()
}
