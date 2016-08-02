package pdu

import (
	"encoding/binary"

	"github.com/mathpl/typed"
)

type Range struct {
	From ObjectIdentifier
	To   ObjectIdentifier
}

func (r *Range) Read(f Flags, rb *typed.ReadBuffer) {
	r.From.Read(f, rb)
	r.To.Read(f, rb)
}

func (r *Range) Write(bo binary.ByteOrder, wb *typed.WriteBuffer) error {
	r.From.Write(bo, wb)
	r.To.Write(bo, wb)
	return wb.Err()
}
