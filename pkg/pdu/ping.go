package pdu

import (
	"encoding/binary"

	"github.com/mathpl/typed"
)

type Ping struct {
}

func (r *Ping) Type() Type {
	return TypePing
}

func (r *Ping) Read(f Flags, rb *typed.ReadBuffer) error {
	return rb.Err()
}

func (r *Ping) Write(bo binary.ByteOrder, wb *typed.WriteBuffer) error {
	return wb.Err()
}
