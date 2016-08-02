package pdu

import (
	"encoding/binary"

	"github.com/mathpl/typed"
)

type Ranges []Range

func (r *Ranges) Read(f Flags, rb *typed.ReadBuffer) error {
	for rb.BytesRemaining() > 0 {
		ran := Range{}
		ran.Read(f, rb)
		*r = append(*r, ran)
	}
	return rb.Err()
}

func (r *Ranges) Write(bo binary.ByteOrder, wb *typed.WriteBuffer) error {
	for _, ran := range *r {
		ran.Write(bo, wb)
	}
	return wb.Err()
}
